import type { PoolClient } from "pg";
import type { InternalContext } from "../context/internal-context";
import type { CharacterRepository } from "../repos/character-repository";
import type { UnifiedRepository } from "../repos/unified-repository";
import type { SessionRepository } from "../repos/session-repository";
import type { SocialService } from "./social-service";
import type { CharacterModel } from "../types/repository-models";
import type {
    FamilyMember,
    FamilyReply,
    FamilyRequest,
    FamilyState,
} from "../protobuf/generated/fminternal/internal_service";
import {
    FamilyAction,
    FamilyEvent,
    FamilyFeature,
    FamilyResult,
} from "../protobuf/generated/fminternal/internal_service";

type Relative = {
    character_id: number;
    senior_id: number;
    name: string;
    level: number;
    class_id: number;
    reputation: number;
    total_reputation: number;
    notice: string;
    progress: string;
    buff_until: Date | null;
    buff_feature: FamilyFeature;
};

type FamilyMessage = {
    id: number;
    payload: Record<string, unknown>;
};

const TREE_LOCK = "SELECT pg_advisory_xact_lock(281, $1)";
const TREES = `
    WITH RECURSIVE up AS (
        SELECT character_id, senior_id FROM social_family WHERE character_id = ANY($1::int[])
        UNION
        SELECT f.character_id, f.senior_id FROM social_family f JOIN up ON f.character_id = up.senior_id
    ), down AS (
        SELECT f.* FROM social_family f JOIN up ON f.character_id = up.character_id AND up.senior_id = 0
        UNION
        SELECT f.* FROM social_family f JOIN down ON f.senior_id = down.character_id
    )
    SELECT * FROM down`;
const COSTS = [300, 500, 700, 800, 1000, 1200, 3000];
const TODAY = "(NOW() AT TIME ZONE 'Asia/Seoul')::date";
const UNITY_MEMBERS = 7;

export function familyRoot(rows: Map<number, Relative>, id: number): number {
    const seen = new Set<number>();
    while (rows.get(id)?.senior_id) {
        if (seen.has(id)) {
            throw new Error("Family cycle");
        }
        seen.add(id);
        id = rows.get(id)!.senior_id;
    }
    return id;
}

export function checkRecruit(
    rows: Map<number, Relative>,
    senior: number,
    junior: number,
    seniorLevel: number,
    juniorLevel: number,
): FamilyResult {
    if (senior === junior || juniorLevel < 10) {
        return FamilyResult.FAMILY_RESULT_CANNOT_BE_JUNIOR;
    }
    if (seniorLevel <= juniorLevel) {
        return FamilyResult.FAMILY_RESULT_JUNIOR_LEVEL_TOO_HIGH;
    }
    if (seniorLevel - juniorLevel > 20) {
        return FamilyResult.FAMILY_RESULT_LEVEL_GAP_TOO_LARGE;
    }
    if (rows.get(junior)?.senior_id) {
        return FamilyResult.FAMILY_RESULT_ALREADY_JUNIOR;
    }
    if ([...rows.values()].filter(r => r.senior_id === senior).length >= 2) {
        return FamilyResult.FAMILY_RESULT_CANNOT_BE_JUNIOR;
    }
    if (familyRoot(rows, senior) === junior) {
        return FamilyResult.FAMILY_RESULT_SAME_FAMILY;
    }
    return FamilyResult.FAMILY_RESULT_OK;
}

function familyFailure(status = FamilyResult.FAMILY_RESULT_FAILED): FamilyReply {
    return { status, family: undefined, targetId: 0, targetName: "", accepted: false };
}

export class FamilyService {
    constructor(
        private readonly internalContext: InternalContext,
        private readonly characterRepository: CharacterRepository,
        private readonly unifiedRepository: UnifiedRepository,
        private readonly sessionRepository: SessionRepository,
        private readonly socialService: SocialService,
    ) {}

    // Whole trees that contain any of the ids, read in one statement so a
    // concurrent merge or split cannot show half a tree.
    private async trees(tx: PoolClient, ids: number[]) {
        const result = await tx.query<Relative>(TREES, [ids]);
        return new Map(result.rows.map(r => [r.character_id, r]));
    }

    // Creating the row must not lock an existing one: a tree operation that
    // holds the tree lock may be waiting to update it.
    private async ensure(tx: PoolClient, character: CharacterModel) {
        await tx.query(
            `INSERT INTO social_family (character_id, name, level, class_id, progress) VALUES ($1, $2, $3, $4, $5)
            ON CONFLICT (character_id) DO NOTHING`,
            [character.characterId, character.name, character.level, character.classId, character.familyProgress ?? "0"],
        );
    }

    private async rename(tx: PoolClient, character: CharacterModel) {
        await tx.query(
            "UPDATE social_family SET name=$2 WHERE character_id=$1 AND name<>$2",
            [character.characterId, character.name],
        );
    }

    // Only work that holds a tree's root lock may change that tree, so roots
    // are locked in id order and the transaction restarts, releasing them, if
    // a root moved before its lock was taken.
    private async withTrees<T>(
        world: number,
        members: CharacterModel[],
        others: number[],
        fn: (tx: PoolClient, rows: Map<number, Relative>) => Promise<T>,
    ): Promise<T> {
        const ids = [...members.map(m => m.characterId), ...others];
        const roots = (rows: Map<number, Relative>) => [...new Set(ids.filter(id => rows.has(id)).map(id => familyRoot(rows, id)))]
            .sort((a, b) => a - b);
        for (;;) {
            const result = await this.internalContext.withPgGlobalTransaction(world, async tx => {
                for (const member of members) {
                    await this.ensure(tx, member);
                }
                const locked = roots(await this.trees(tx, ids));
                for (const root of locked) {
                    await tx.query(TREE_LOCK, [root]);
                }
                for (const member of members) {
                    await this.rename(tx, member);
                }

                const rows = await this.trees(tx, ids);
                if (roots(rows).join() !== locked.join()) {
                    return undefined;
                }
                return { value: await fn(tx, rows) };
            });
            if (result) {
                return result.value;
            }
        }
    }

    private async state(tx: PoolClient, world: number, rows: Map<number, Relative>, id: number): Promise<FamilyState> {
        const me = rows.get(id)!;
        const root = familyRoot(rows, id);
        const leader = rows.get(root)!;
        const children = (parent: number) => [...rows.values()]
            .filter(r => r.senior_id === parent)
            .sort((a, b) => a.character_id - b.character_id);
        const descendants = (parent: number): number => {
            const pending = children(parent);
            const seen = new Set<number>();
            while (pending.length) {
                const child = pending.pop()!;
                if (seen.has(child.character_id)) {
                    throw new Error("Family cycle");
                }
                seen.add(child.character_id);
                pending.push(...children(child.character_id));
            }
            return seen.size;
        };

        // Client expects the leader first, then two generations above and below the selected character.
        const visible: Relative[] = [leader];
        const add = (row: Relative | undefined) => {
            if (row && !visible.includes(row)) {
                visible.push(row);
            }
        };
        const senior = rows.get(me.senior_id);
        add(senior ? rows.get(senior.senior_id) : undefined);
        add(senior);
        add(me);
        if (senior) {
            for (const sibling of children(senior.character_id)) {
                add(sibling);
            }
        }
        for (const junior of children(id)) {
            add(junior);
        }
        for (const junior of children(id)) {
            for (const grand of children(junior.character_id)) {
                add(grand);
            }
        }

        const members: FamilyMember[] = [];
        for (const row of visible) {
            members.push({
                characterId: row.character_id,
                seniorId: row.senior_id,
                name: row.name,
                level: row.level,
                classId: row.class_id,
                reputation: row.reputation,
                totalReputation: row.total_reputation,
                channelId: (await this.sessionRepository.findChannel(world, row.character_id)) ?? -1,
                descendants: descendants(row.character_id),
            });
        }

        const uses = await tx.query<{ feature: FamilyFeature }>(
            `SELECT feature FROM social_family_uses WHERE character_id=$1 AND day=${TODAY} ORDER BY feature`,
            [id],
        );
        const inFamily = me.senior_id !== 0 || children(id).length > 0;
        return {
            characterId: id,
            seniorId: me.senior_id,
            juniors: children(id).map(r => r.character_id),
            leaderId: inFamily ? root : 0,
            leaderName: inFamily ? leader.name : "",
            notice: leader.notice,
            reputation: me.reputation,
            totalReputation: me.total_reputation,
            usedFeatures: uses.rows.map(r => r.feature),
            members,
            buffUntil: String(me.buff_until?.getTime() ?? 0),
            buffFeature: me.buff_feature,
        };
    }

    async expireMoves(world: number) {
        const payers = await this.internalContext.withPgGlobalTransaction(world, async tx => {
            const moves = await tx.query<{ operation_id: string; payer_id: number; feature: number }>(
                "UPDATE social_family_moves SET status='refunded' WHERE status='pending' AND expires_at<=NOW() RETURNING *",
            );
            for (const move of moves.rows) {
                await tx.query("DELETE FROM social_family_uses WHERE operation_id=$1", [move.operation_id]);
                await tx.query(
                    "UPDATE social_family SET reputation=LEAST(2147483647, reputation::bigint+$2) WHERE character_id=$1",
                    [move.payer_id, COSTS[move.feature]],
                );
            }
            return moves.rows.map(move => move.payer_id);
        });

        for (const id of payers) {
            await this.notify(world, id, { event: FamilyEvent.FAMILY_EVENT_REFRESH }).catch(() => {});
        }
    }

    async load(world: number, character: CharacterModel) {
        return this.internalContext.withPgGlobalTransaction(world, async tx => {
            await this.ensure(tx, character);
            await this.rename(tx, character);
            return this.state(tx, world, await this.trees(tx, [character.characterId]), character.characterId);
        });
    }

    // Receivers are locked through their own trees: a senior may have left
    // since the game server earned that reputation.
    async progress(world: number, character: CharacterModel) {
        const receivers = Object.keys(character.familyReputation ?? {}).map(Number);
        const changed = await this.withTrees(world, [character], receivers, async (tx, rows) => {
            await tx.query(
                "UPDATE social_family SET level=$2, class_id=$3, progress=GREATEST(progress, $4) WHERE character_id=$1",
                [character.characterId, character.level, character.classId, character.familyProgress ?? "0"],
            );

            const ids: number[] = [];
            for (const [receiver, reputation] of Object.entries(character.familyReputation ?? {})) {
                const receiverId = Number(receiver);
                if (!rows.has(receiverId) || receiverId === character.characterId) {
                    continue;
                }

                const credit = await tx.query<{ reputation: string }>(
                    "SELECT reputation FROM social_family_credit WHERE character_id=$1 AND receiver_id=$2",
                    [character.characterId, receiverId],
                );
                const delta = BigInt(reputation) - BigInt(credit.rows[0]?.reputation ?? "0");
                if (delta <= 0n) {
                    continue;
                }

                await tx.query(
                    `INSERT INTO social_family_credit VALUES ($1, $2, $3)
                    ON CONFLICT (character_id, receiver_id) DO UPDATE SET reputation=EXCLUDED.reputation`,
                    [character.characterId, receiverId, reputation],
                );
                await tx.query(
                    `UPDATE social_family SET
                        reputation=LEAST(2147483647, reputation::bigint+$2::bigint),
                        total_reputation=LEAST(2147483647, total_reputation::bigint+$2::bigint)
                    WHERE character_id=$1`,
                    [receiverId, String(delta)],
                );
                ids.push(receiverId);
            }
            return ids;
        });

        for (const id of changed) {
            await this.notify(world, id, { event: FamilyEvent.FAMILY_EVENT_REFRESH }).catch(() => {});
        }
    }

    private async notify(world: number, id: number, payload: Record<string, unknown>) {
        const channel = await this.sessionRepository.findChannel(world, id);
        if (channel == null) {
            return;
        }
        await this.socialService.publish(world, channel, "family", { character_id: id, ...payload });
    }

    async presence(world: number, characterId: number) {
        const ids = await this.internalContext.withPgGlobalTransaction(world, async tx => {
            const rows = await this.trees(tx, [characterId]);
            if (!rows.has(characterId)) {
                return [];
            }

            const root = familyRoot(rows, characterId);
            return [...rows.values()]
                .filter(row => familyRoot(rows, row.character_id) === root)
                .map(row => row.character_id);
        });

        for (const id of ids) {
            await this.notify(world, id, { event: FamilyEvent.FAMILY_EVENT_REFRESH }).catch(() => {});
        }
    }

    async family(req: FamilyRequest): Promise<FamilyReply> {
        const world = req.worldId;
        const id = req.characterId;
        const me = await this.characterRepository.get(world, id);
        const settling = req.action === FamilyAction.FAMILY_ACTION_REFUND || req.action === FamilyAction.FAMILY_ACTION_MOVE_COMPLETE;
        if (!me || (!settling && (await this.sessionRepository.findChannel(world, id)) !== req.channelId)) {
            return familyFailure();
        }

        let targetId = req.targetId;
        const named = [
            FamilyAction.FAMILY_ACTION_PEDIGREE,
            FamilyAction.FAMILY_ACTION_INVITE,
            FamilyAction.FAMILY_ACTION_USE,
            FamilyAction.FAMILY_ACTION_SUMMON,
        ];
        if (req.name && named.includes(req.action)) {
            const entry = await this.unifiedRepository.findCharacterNameEntry(req.name);
            if (entry?.world_id !== world) {
                return familyFailure(FamilyResult.FAMILY_RESULT_WRONG_NAME);
            }
            targetId = entry.character_id;
        }

        const target = targetId ? await this.characterRepository.get(world, targetId) : null;
        const events: FamilyMessage[] = [];
        const members = target ? [me, target] : [me];
        const reply = await this.withTrees(world, members, [], async (tx, rows) => {
            const own = rows.get(id)!;
            switch (req.action) {
                case FamilyAction.FAMILY_ACTION_PEDIGREE:
                case FamilyAction.FAMILY_ACTION_OPEN:
                    break;
                case FamilyAction.FAMILY_ACTION_MOVE_COMPLETE:
                case FamilyAction.FAMILY_ACTION_REFUND: {
                    if (!req.operationId) {
                        return familyFailure();
                    }

                    const complete = req.action === FamilyAction.FAMILY_ACTION_MOVE_COMPLETE;
                    const move = await tx.query<{ payer_id: number; feature: number }>(
                        `UPDATE social_family_moves SET status=$3
                        WHERE operation_id=$1 AND character_id=$2 AND status='pending' RETURNING *`,
                        [req.operationId, id, complete ? "completed" : "refunded"],
                    );
                    if (!complete) {
                        // Cancellation may arrive before the timed-out reservation RPC.
                        await tx.query(
                            `INSERT INTO social_family_moves (operation_id, character_id, payer_id, feature, status)
                            VALUES ($1, $2, 0, 0, 'refunded') ON CONFLICT DO NOTHING`,
                            [req.operationId, id],
                        );
                        for (const reservation of move.rows) {
                            await tx.query("DELETE FROM social_family_uses WHERE operation_id=$1", [req.operationId]);
                            await tx.query(
                                "UPDATE social_family SET reputation=LEAST(2147483647, reputation::bigint+$2) WHERE character_id=$1",
                                [reservation.payer_id, COSTS[reservation.feature]],
                            );
                            events.push({ id: reservation.payer_id, payload: { event: FamilyEvent.FAMILY_EVENT_REFRESH } });
                        }
                    }
                    break;
                }
                case FamilyAction.FAMILY_ACTION_INVITE:
                case FamilyAction.FAMILY_ACTION_ACCEPT: {
                    const inviting = req.action === FamilyAction.FAMILY_ACTION_INVITE;
                    const seniorId = inviting ? id : targetId;
                    const juniorId = inviting ? targetId : id;
                    const seniorLevel = inviting ? req.level : req.targetLevel;
                    const juniorLevel = inviting ? req.targetLevel : req.level;
                    const senior = rows.get(seniorId);
                    const junior = rows.get(juniorId);
                    if (!senior || !junior) {
                        return familyFailure(FamilyResult.FAMILY_RESULT_WRONG_NAME);
                    }
                    if ((await this.sessionRepository.findChannel(world, targetId)) !== req.channelId) {
                        return familyFailure(FamilyResult.FAMILY_RESULT_NOT_SAME_MAP);
                    }
                    const recruit = checkRecruit(rows, seniorId, juniorId, seniorLevel, juniorLevel);
                    if (recruit !== FamilyResult.FAMILY_RESULT_OK) {
                        return familyFailure(recruit);
                    }

                    if (inviting) {
                        await tx.query(
                            `INSERT INTO social_family_invites VALUES ($1, $2, $3, NOW()+INTERVAL '60 seconds')
                            ON CONFLICT (receiver_id) DO UPDATE SET sender_id=$2, channel_id=$3, expires_at=EXCLUDED.expires_at`,
                            [juniorId, seniorId, req.channelId],
                        );
                        events.push({ id: juniorId, payload: { event: FamilyEvent.FAMILY_EVENT_INVITE, sender_id: seniorId, name: senior.name } });
                        break;
                    }

                    const invite = await tx.query(
                        `DELETE FROM social_family_invites
                        WHERE receiver_id=$1 AND sender_id=$2 AND channel_id=$3 AND expires_at>NOW() RETURNING *`,
                        [id, targetId, req.channelId],
                    );
                    if (!invite.rowCount || req.name !== senior.name) {
                        return familyFailure();
                    }

                    if (req.accepted) {
                        await tx.query("UPDATE social_family SET senior_id=$2 WHERE character_id=$1", [id, targetId]);
                    }
                    events.push({ id: targetId, payload: { event: FamilyEvent.FAMILY_EVENT_INVITE_RESULT, name: me.name, accepted: req.accepted } });
                    if (req.accepted) {
                        events.push({ id, payload: { event: FamilyEvent.FAMILY_EVENT_JOINED, name: senior.name } });
                    }
                    break;
                }
                case FamilyAction.FAMILY_ACTION_REMOVE_JUNIOR:
                case FamilyAction.FAMILY_ACTION_REMOVE_SENIOR: {
                    const leaving = req.action === FamilyAction.FAMILY_ACTION_REMOVE_SENIOR;
                    const juniorId = leaving ? id : targetId;
                    if (leaving ? !own.senior_id : rows.get(juniorId)?.senior_id !== id) {
                        return familyFailure();
                    }

                    await tx.query("UPDATE social_family SET senior_id=0 WHERE character_id=$1", [juniorId]);
                    break;
                }
                case FamilyAction.FAMILY_ACTION_PRECEPT: {
                    if (familyRoot(rows, id) !== id || req.name.length > 200 || /[\x00-\x1f\x7f]/u.test(req.name)) {
                        return familyFailure();
                    }

                    await tx.query("UPDATE social_family SET notice=$2 WHERE character_id=$1", [id, req.name]);
                    break;
                }
                case FamilyAction.FAMILY_ACTION_USE:
                case FamilyAction.FAMILY_ACTION_SUMMON: {
                    const feature = req.action === FamilyAction.FAMILY_ACTION_SUMMON ? FamilyFeature.FAMILY_FEATURE_SUMMON : req.feature;
                    if (feature < 0 || feature >= COSTS.length) {
                        return familyFailure();
                    }
                    if (!own.senior_id && ![...rows.values()].some(r => r.senior_id === id)) {
                        return familyFailure();
                    }

                    let payer = id;
                    if (feature === FamilyFeature.FAMILY_FEATURE_TELEPORT || feature === FamilyFeature.FAMILY_FEATURE_SUMMON) {
                        const sameFamily = target && targetId !== id && familyRoot(rows, id) === familyRoot(rows, targetId);
                        if (!sameFamily) {
                            return familyFailure(FamilyResult.FAMILY_RESULT_NOT_SAME_FAMILY);
                        }
                        if ((await this.sessionRepository.findChannel(world, targetId)) !== req.channelId) {
                            return familyFailure();
                        }

                        if (req.action === FamilyAction.FAMILY_ACTION_USE && feature === FamilyFeature.FAMILY_FEATURE_SUMMON) {
                            if (own.reputation < COSTS[feature]!) {
                                return familyFailure();
                            }

                            const used = await tx.query(
                                `SELECT 1 FROM social_family_uses WHERE character_id=$1 AND feature=$2 AND day=${TODAY}`,
                                [id, feature],
                            );
                            if (used.rowCount) {
                                return familyFailure();
                            }

                            await tx.query("DELETE FROM social_family_summons WHERE expires_at<=NOW()");
                            const summon = await tx.query(
                                `INSERT INTO social_family_summons VALUES ($1, $2, $3, NOW()+INTERVAL '60 seconds')
                                ON CONFLICT (receiver_id) DO NOTHING RETURNING *`,
                                [targetId, id, req.channelId],
                            );
                            if (!summon.rowCount) {
                                return familyFailure(FamilyResult.FAMILY_RESULT_SUMMON_PENDING);
                            }

                            events.push({ id: targetId, payload: { event: FamilyEvent.FAMILY_EVENT_SUMMON, name: me.name, map_name: req.mapName } });
                            return { ...familyFailure(FamilyResult.FAMILY_RESULT_OK), targetId, targetName: target!.name };
                        }

                        if (req.action === FamilyAction.FAMILY_ACTION_SUMMON) {
                            const pending = await tx.query(
                                `DELETE FROM social_family_summons
                                WHERE receiver_id=$1 AND sender_id=$2 AND channel_id=$3 AND expires_at>NOW() RETURNING *`,
                                [id, targetId, req.channelId],
                            );
                            if (!pending.rowCount) {
                                return familyFailure(FamilyResult.FAMILY_RESULT_SUMMON_FAILED);
                            }
                            if (!req.accepted) {
                                return { ...familyFailure(FamilyResult.FAMILY_RESULT_OK), accepted: false };
                            }
                            payer = targetId;
                        }
                    }

                    if (rows.get(payer)!.reputation < COSTS[feature]!) {
                        return familyFailure();
                    }

                    const use = await tx.query(
                        `INSERT INTO social_family_uses VALUES ($1, $2, ${TODAY}, $3) ON CONFLICT DO NOTHING RETURNING *`,
                        [payer, feature, req.operationId],
                    );
                    if (!use.rowCount) {
                        return familyFailure();
                    }

                    let beneficiaries = [id];
                    if (feature === FamilyFeature.FAMILY_FEATURE_UNITY) {
                        const juniors = [...rows.values()]
                            .filter(r => r.senior_id === id || rows.get(r.senior_id)?.senior_id === id)
                            .map(r => r.character_id);
                        const online: number[] = [];
                        for (const member of [id, ...juniors]) {
                            if ((await this.sessionRepository.findChannel(world, member)) != null) {
                                online.push(member);
                            }
                        }
                        if (online.length < UNITY_MEMBERS) {
                            await tx.query(
                                `DELETE FROM social_family_uses WHERE character_id=$1 AND feature=$2 AND day=${TODAY}`,
                                [payer, feature],
                            );
                            return familyFailure();
                        }
                        beneficiaries = online;
                    }

                    if (feature === FamilyFeature.FAMILY_FEATURE_TELEPORT || feature === FamilyFeature.FAMILY_FEATURE_SUMMON) {
                        const move = await tx.query(
                            `INSERT INTO social_family_moves (operation_id, character_id, payer_id, feature)
                            VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING RETURNING *`,
                            [req.operationId, id, payer, feature],
                        );
                        if (!move.rowCount) {
                            await tx.query("DELETE FROM social_family_uses WHERE operation_id=$1", [req.operationId]);
                            return familyFailure();
                        }
                    }

                    await tx.query(
                        "UPDATE social_family SET reputation=reputation-$2 WHERE character_id=$1",
                        [payer, COSTS[feature]],
                    );
                    if (payer !== id) {
                        events.push({ id: payer, payload: { event: FamilyEvent.FAMILY_EVENT_REFRESH } });
                    }
                    if (feature !== FamilyFeature.FAMILY_FEATURE_TELEPORT && feature !== FamilyFeature.FAMILY_FEATURE_SUMMON) {
                        const short = feature === FamilyFeature.FAMILY_FEATURE_DROP_15_MINUTES || feature === FamilyFeature.FAMILY_FEATURE_EXP_15_MINUTES;
                        await tx.query(
                            `UPDATE social_family SET buff_feature=$2, buff_until=NOW()+($3::int*INTERVAL '1 minute')
                            WHERE character_id=ANY($1::int[])`,
                            [beneficiaries, feature, short ? 15 : 30],
                        );
                        for (const member of beneficiaries) {
                            events.push({ id: member, payload: { event: FamilyEvent.FAMILY_EVENT_REFRESH } });
                        }
                    }
                    break;
                }
                default:
                    return familyFailure();
            }

            rows = await this.trees(tx, [id, targetId, own.senior_id].filter(n => n !== 0));
            const regrouped = [
                FamilyAction.FAMILY_ACTION_REMOVE_JUNIOR,
                FamilyAction.FAMILY_ACTION_REMOVE_SENIOR,
                FamilyAction.FAMILY_ACTION_ACCEPT,
                FamilyAction.FAMILY_ACTION_PRECEPT,
            ];
            if (regrouped.includes(req.action)) {
                // Includes both halves after splitting and both trees after merging.
                const affected = new Set([
                    familyRoot(rows, id),
                    targetId && familyRoot(rows, targetId),
                    own.senior_id && familyRoot(rows, own.senior_id),
                ]);
                for (const row of rows.values()) {
                    if (affected.has(familyRoot(rows, row.character_id))) {
                        events.push({ id: row.character_id, payload: { event: FamilyEvent.FAMILY_EVENT_REFRESH } });
                    }
                }
            }

            const shown = req.action === FamilyAction.FAMILY_ACTION_PEDIGREE && target ? targetId : id;
            return {
                status: FamilyResult.FAMILY_RESULT_OK,
                family: await this.state(tx, world, rows, shown),
                targetId,
                targetName: target?.name ?? "",
                accepted: req.accepted,
            };
        });

        for (const event of events) {
            await this.notify(world, event.id, event.payload).catch(() => {});
        }
        return reply;
    }
}
