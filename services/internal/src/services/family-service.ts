import type { PoolClient } from "pg";
import type { InternalContext } from "@fm/shared/context/internal-context";
import type { CharacterRepository } from "@fm/shared/repos/character-repository";
import type { FamilyMoveRequestRepository } from "@fm/shared/repos/family-request-repository";
import { RequestStatus } from "@fm/shared/repos/request-repository";
import type { UnifiedRepository } from "@fm/shared/repos/unified-repository";
import type { SessionRepository } from "@fm/shared/repos/session-repository";
import type { SocialService } from "./social-service";
import type { CharacterModel } from "@fm/shared/types/repository-models";
import type {
    FamilyMember,
    FamilyReply,
    FinishFamilyMoveRequest,
    GetFamilyReply,
    GetFamilyRequest,
    JoinFamilyRequest,
    SeparateFamilyRequest,
    SetFamilyPreceptRequest,
    UseFamilyFeatureRequest,
} from "@fm/shared/protobuf/generated/fminternal/internal_service";
import { Family, FamilyFeature, FamilyResult } from "@fm/shared/protobuf/generated/fminternal/internal_service";

type Relative = {
    character_id: number;
    leader_id: number | null;
    senior_id: number;
    name: string;
    level: number;
    class_id: number;
    reputation: number;
    total_reputation: number;
    buff_until: Date | null;
    buff_feature: FamilyFeature;
};

const RELATIVES = `
    SELECT * FROM social_family
    WHERE character_id = ANY($1::int[])
       OR leader_id IN (SELECT leader_id FROM social_family WHERE character_id = ANY($1::int[]) AND leader_id IS NOT NULL)
    ORDER BY character_id`;
const COSTS = [300, 500, 700, 800, 1000, 1200, 3000];
const TODAY = "(NOW() AT TIME ZONE 'Asia/Seoul')::date";
const UNITY_MEMBERS = 7;
const MAX_JUNIORS = 2;

export class FamilyService {
    constructor(
        private readonly internalContext: InternalContext,
        private readonly characterRepository: CharacterRepository,
        private readonly familyMoveRequestRepository: FamilyMoveRequestRepository,
        private readonly unifiedRepository: UnifiedRepository,
        private readonly sessionRepository: SessionRepository,
        private readonly socialService: SocialService,
    ) {}

    // Everyone in the families of the ids, plus the ids themselves when they have no family.
    private async relatives(tx: PoolClient, ids: number[]) {
        const result = await tx.query<Relative>(RELATIVES, [ids]);
        return new Map(result.rows.map(r => [r.character_id, r]));
    }

    // The member and everyone below it.
    private subtree(rows: Map<number, Relative>, id: number): number[] {
        const ids = [id];
        for (let i = 0; i < ids.length; i++) {
            for (const row of rows.values()) {
                if (row.senior_id === ids[i]) {
                    ids.push(row.character_id);
                }
            }
        }
        return ids;
    }

    // Channels are looked up after the global connection is returned; the lookup reads the data shards.
    private async families(world: number, ids: number[]): Promise<Family[]> {
        const { rows, precepts, uses } = await this.internalContext.withPgGlobalTransaction(world, async tx => {
            const rows = await this.relatives(tx, ids);
            const leaders = [...new Set([...rows.values()].map(r => r.leader_id).filter(id => id !== null))];
            const precepts = await tx.query<{ leader_id: number; precept: string }>(
                "SELECT leader_id, precept FROM families WHERE leader_id = ANY($1::int[])",
                [leaders],
            );
            const uses = await tx.query<{ character_id: number; feature: FamilyFeature }>(
                `SELECT character_id, feature FROM social_family_uses WHERE character_id = ANY($1::int[]) AND day=${TODAY} ORDER BY feature`,
                [[...rows.keys()]],
            );
            return { rows, precepts, uses };
        });

        const families = new Map<number, Family>();
        for (const row of rows.values()) {
            const key = row.leader_id ?? row.character_id;
            let family = families.get(key);
            if (!family) {
                family = {
                    leaderId: row.leader_id ?? undefined,
                    precept: precepts.rows.find(p => p.leader_id === row.leader_id)?.precept ?? "",
                    members: [],
                };
                families.set(key, family);
            }
            family.members.push({
                characterId: row.character_id,
                seniorId: row.senior_id,
                name: row.name,
                level: row.level,
                classId: row.class_id,
                reputation: row.reputation,
                totalReputation: row.total_reputation,
                channelId: (await this.sessionRepository.findChannel(world, row.character_id)) ?? -1,
                buffUntil: String(row.buff_until?.getTime() ?? 0),
                buffFeature: row.buff_feature,
                usedFeatures: uses.rows.filter(u => u.character_id === row.character_id).map(u => u.feature),
            });
        }
        return [...families.values()];
    }

    private async family(world: number, id: number) {
        return (await this.families(world, [id])).find(f => f.members.some(m => m.characterId === id));
    }

    private alone(character: CharacterModel): Family {
        return {
            leaderId: undefined,
            precept: "",
            members: [{
                characterId: character.characterId,
                seniorId: 0,
                name: character.name,
                level: character.level,
                classId: character.classId,
                reputation: 0,
                totalReputation: 0,
                channelId: -1,
                buffUntil: "0",
                buffFeature: FamilyFeature.FAMILY_FEATURE_TELEPORT,
                usedFeatures: [],
            }],
        };
    }

    // Sends each family to its online members, one message per channel.
    private async publish(world: number, ids: number[]): Promise<Family[]> {
        const families = await this.families(world, ids);
        for (const family of families) {
            const channels = new Map<number, number[]>();
            for (const member of family.members) {
                if (member.channelId < 0) {
                    continue;
                }
                channels.set(member.channelId, [...(channels.get(member.channelId) ?? []), member.characterId]);
            }
            for (const [channel, members] of channels) {
                await this.socialService
                    .publish(world, channel, "family", { character_ids: members, family: Family.toJSON(family) })
                    .catch(err => console.error("[family] publish failed", err));
            }
        }
        return families;
    }

    private async reply(world: number, status: FamilyResult, id: number, ...others: number[]): Promise<FamilyReply> {
        if (status !== FamilyResult.FAMILY_RESULT_OK) {
            return { status, family: undefined };
        }

        const families = await this.publish(world, [id, ...others.filter(n => n !== 0)]);
        return { status, family: families.find(f => f.members.some(m => m.characterId === id)) };
    }

    // A family changes only while its families row is locked. Families are locked
    // in leader order before the member rows; the transaction restarts if a member
    // changed family in between.
    private async withFamilies<T>(
        world: number,
        ids: number[],
        fn: (tx: PoolClient, rows: Map<number, Relative>) => Promise<T>,
    ): Promise<T> {
        const sorted = [...new Set(ids)].sort((a, b) => a - b);
        const memberships = async (tx: PoolClient, lock: boolean) => {
            const result = await tx.query<{ leader_id: number | null }>(
                `SELECT leader_id FROM social_family WHERE character_id = ANY($1::int[]) ORDER BY character_id${lock ? " FOR UPDATE" : ""}`,
                [sorted],
            );
            return result.rows.map(r => r.leader_id);
        };
        for (;;) {
            const result = await this.internalContext.withPgGlobalTransaction(world, async tx => {
                const before = await memberships(tx, false);
                await tx.query(
                    "SELECT leader_id FROM families WHERE leader_id = ANY($1::int[]) ORDER BY leader_id FOR UPDATE",
                    [before.filter(id => id !== null)],
                );
                if ((await memberships(tx, true)).join() !== before.join()) {
                    return undefined;
                }

                return { value: await fn(tx, await this.relatives(tx, sorted)) };
            });
            if (result) {
                return result.value;
            }
        }
    }

    async load(world: number, character: CharacterModel): Promise<Family> {
        const family = await this.family(world, character.characterId);
        return family ?? this.alone(character);
    }

    async get(req: GetFamilyRequest): Promise<GetFamilyReply> {
        const entry = await this.unifiedRepository.findCharacterNameEntry(req.name);
        if (entry?.world_id !== req.worldId) {
            return { found: false, characterId: 0, family: undefined };
        }

        const world = req.worldId;
        const id = entry.character_id;
        const family = await this.family(world, id);
        if (family) {
            return { found: true, characterId: id, family };
        }
        const character = await this.characterRepository.get(world, id);
        if (!character) {
            return { found: false, characterId: 0, family: undefined };
        }
        return { found: true, characterId: id, family: this.alone(character) };
    }

    // Reputation goes to whoever was senior when it was earned, even if they have
    // left since. Rows are updated in id order so two saves cannot deadlock.
    // Saves of one junior run one at a time, or two of them could both pay the same credit delta.
    async progress(world: number, character: CharacterModel) {
        const reputation = character.familyReputation ?? {};
        const ids = [character.characterId, ...Object.keys(reputation).map(Number)]
            .filter((id, i, all) => all.indexOf(id) === i)
            .sort((a, b) => a - b);
        const receivers = await this.internalContext.withPgGlobalTransaction(world, async tx => {
            await tx.query("SELECT pg_advisory_xact_lock($1, $2)", [0x66616d69, character.characterId]);
            const credited: number[] = [];
            for (const id of ids) {
                if (id === character.characterId) {
                    await tx.query(
                        `UPDATE social_family SET name=$2, level=$3, class_id=$4, progress=GREATEST(progress, $5)
                        WHERE character_id=$1 AND (name<>$2 OR level<>$3 OR class_id<>$4 OR progress<$5::bigint)`,
                        [id, character.name, character.level, character.classId, character.familyProgress ?? "0"],
                    );
                    continue;
                }

                const credit = await tx.query<{ reputation: string }>(
                    "SELECT reputation FROM social_family_credit WHERE character_id=$1 AND receiver_id=$2",
                    [character.characterId, id],
                );
                const delta = BigInt(reputation[id]!) - BigInt(credit.rows[0]?.reputation ?? "0");
                if (delta <= 0n) {
                    continue;
                }

                const updated = await tx.query(
                    `UPDATE social_family SET
                        reputation=LEAST(2147483647, reputation::bigint+$2::bigint),
                        total_reputation=LEAST(2147483647, total_reputation::bigint+$2::bigint)
                    WHERE character_id=$1`,
                    [id, String(delta)],
                );
                if (!updated.rowCount) {
                    continue;
                }

                await tx.query(
                    `INSERT INTO social_family_credit VALUES ($1, $2, $3)
                    ON CONFLICT (character_id, receiver_id) DO UPDATE SET reputation=EXCLUDED.reputation`,
                    [character.characterId, id, reputation[id]],
                );
                credited.push(id);
            }
            return credited;
        });

        if (receivers.length > 0) {
            await this.publish(world, receivers);
        }
    }

    async presence(world: number, characterId: number) {
        await this.publish(world, [characterId]);
    }

    private checkJoin(rows: Map<number, Relative>, senior: FamilyMember, junior: FamilyMember): FamilyResult {
        if (senior.characterId === junior.characterId || junior.level < 10) {
            return FamilyResult.FAMILY_RESULT_CANNOT_BE_JUNIOR;
        }
        if (senior.level <= junior.level) {
            return FamilyResult.FAMILY_RESULT_JUNIOR_LEVEL_TOO_HIGH;
        }
        if (senior.level - junior.level > 20) {
            return FamilyResult.FAMILY_RESULT_LEVEL_GAP_TOO_LARGE;
        }
        if (rows.get(junior.characterId)?.senior_id) {
            return FamilyResult.FAMILY_RESULT_ALREADY_JUNIOR;
        }
        if ([...rows.values()].filter(r => r.senior_id === senior.characterId).length >= MAX_JUNIORS) {
            return FamilyResult.FAMILY_RESULT_CANNOT_BE_JUNIOR;
        }
        if (rows.get(senior.characterId)?.leader_id === junior.characterId) {
            return FamilyResult.FAMILY_RESULT_SAME_FAMILY;
        }
        return FamilyResult.FAMILY_RESULT_OK;
    }

    async join(req: JoinFamilyRequest): Promise<FamilyReply> {
        const world = req.worldId;
        const senior = req.senior;
        const junior = req.junior;
        if (!senior || !junior) {
            return { status: FamilyResult.FAMILY_RESULT_FAILED, family: undefined };
        }

        // Rows are created outside the family transaction so that creating one
        // never waits on a row another family transaction holds.
        await this.internalContext.withPgGlobalTransaction(world, async tx => {
            for (const member of [senior, junior]) {
                await tx.query(
                    `INSERT INTO social_family (character_id, name, level, class_id) VALUES ($1, $2, $3, $4)
                    ON CONFLICT (character_id) DO NOTHING`,
                    [member.characterId, member.name, member.level, member.classId],
                );
            }
        });

        const status = await this.withFamilies(world, [senior.characterId, junior.characterId], async (tx, rows) => {
            const result = this.checkJoin(rows, senior, junior);
            if (result !== FamilyResult.FAMILY_RESULT_OK) {
                return result;
            }

            const juniorLeader = rows.get(junior.characterId)!.leader_id;
            let leader = rows.get(senior.characterId)!.leader_id;
            if (leader === null) {
                leader = senior.characterId;
                await tx.query("INSERT INTO families (leader_id) VALUES ($1)", [leader]);
            }
            await tx.query("UPDATE social_family SET senior_id=$2 WHERE character_id=$1", [junior.characterId, senior.characterId]);
            await tx.query(
                "UPDATE social_family SET leader_id=$1 WHERE character_id = ANY($2::int[]) OR leader_id=$3",
                [leader, [senior.characterId, junior.characterId], juniorLeader],
            );
            if (juniorLeader !== null) {
                await tx.query("DELETE FROM families WHERE leader_id=$1", [juniorLeader]);
            }
            return FamilyResult.FAMILY_RESULT_OK;
        });
        return this.reply(world, status, junior.characterId);
    }

    async separate(req: SeparateFamilyRequest): Promise<FamilyReply> {
        const world = req.worldId;
        const id = req.characterId;
        const juniorId = req.juniorId;
        let seniorId = 0;
        const status = await this.withFamilies(world, [id, juniorId], async (tx, rows) => {
            const junior = rows.get(juniorId);
            if (!junior || junior.senior_id === 0 || (id !== juniorId && junior.senior_id !== id)) {
                return FamilyResult.FAMILY_RESULT_FAILED;
            }

            seniorId = junior.senior_id;
            await tx.query("UPDATE social_family SET senior_id=0 WHERE character_id=$1", [juniorId]);
            junior.senior_id = 0;

            const split = this.subtree(rows, juniorId);
            const staying = [...rows.values()].filter(r => r.leader_id === junior.leader_id && !split.includes(r.character_id));
            let leader: number | null = null;
            if (split.length > 1) {
                leader = juniorId;
                await tx.query("INSERT INTO families (leader_id) VALUES ($1)", [leader]);
            }
            await tx.query("UPDATE social_family SET leader_id=$1 WHERE character_id = ANY($2::int[])", [leader, split]);
            if (staying.length === 1) {
                await tx.query("UPDATE social_family SET leader_id=NULL WHERE character_id=$1", [staying[0]!.character_id]);
                await tx.query("DELETE FROM families WHERE leader_id=$1", [junior.leader_id]);
            }
            return FamilyResult.FAMILY_RESULT_OK;
        });
        return this.reply(world, status, id, juniorId, seniorId);
    }

    async setPrecept(req: SetFamilyPreceptRequest): Promise<FamilyReply> {
        const world = req.worldId;
        const id = req.characterId;
        const status = await this.withFamilies(world, [id], async (tx, rows) => {
            const me = rows.get(id);
            if (!me || me.leader_id !== id || req.precept.length > 200 || /[\x00-\x1f\x7f]/u.test(req.precept)) {
                return FamilyResult.FAMILY_RESULT_FAILED;
            }

            await tx.query("UPDATE families SET precept=$2 WHERE leader_id=$1", [id, req.precept]);
            return FamilyResult.FAMILY_RESULT_OK;
        });
        return this.reply(world, status, id);
    }

    // The payer's row is locked, so its uses of the day cannot race. A move
    // reservation is written first: a refund that arrived earlier left its
    // tombstone there.
    async use(req: UseFamilyFeatureRequest): Promise<FamilyReply> {
        const world = req.worldId;
        const id = req.characterId;
        const feature = req.feature;
        const moving = feature === FamilyFeature.FAMILY_FEATURE_TELEPORT || feature === FamilyFeature.FAMILY_FEATURE_SUMMON;
        const payer = feature === FamilyFeature.FAMILY_FEATURE_SUMMON ? req.targetId : id;
        const status = await this.withFamilies(world, moving ? [id, req.targetId] : [id], async (tx, rows) => {
            const me = rows.get(id);
            const paying = rows.get(payer);
            if (feature < 0 || feature >= COSTS.length || !me || !paying || paying.leader_id === null || paying.reputation < COSTS[feature]!) {
                return FamilyResult.FAMILY_RESULT_FAILED;
            }
            if (moving && (req.targetId === id || rows.get(req.targetId)?.leader_id !== me.leader_id)) {
                return FamilyResult.FAMILY_RESULT_NOT_SAME_FAMILY;
            }
            if (moving && req.requestId === "0") {
                return FamilyResult.FAMILY_RESULT_FAILED;
            }

            const used = await tx.query(
                `SELECT 1 FROM social_family_uses WHERE character_id=$1 AND feature=$2 AND day=${TODAY}`,
                [payer, feature],
            );
            if (used.rowCount) {
                return FamilyResult.FAMILY_RESULT_FAILED;
            }

            let beneficiaries = [id];
            if (feature === FamilyFeature.FAMILY_FEATURE_UNITY) {
                const juniors = [...rows.values()]
                    .filter(r => r.senior_id === id || rows.get(r.senior_id)?.senior_id === id)
                    .map(r => r.character_id);
                beneficiaries = [];
                for (const member of [id, ...juniors]) {
                    if ((await this.sessionRepository.findChannel(world, member)) != null) {
                        beneficiaries.push(member);
                    }
                }
                if (beneficiaries.length < UNITY_MEMBERS) {
                    return FamilyResult.FAMILY_RESULT_FAILED;
                }
            }

            const operation = moving ? req.requestId : "";
            if (moving) {
                const reserved = await this.familyMoveRequestRepository.reserve(world, id, operation, { payer_id: payer, feature }, tx);
                if (!reserved) {
                    return FamilyResult.FAMILY_RESULT_FAILED;
                }
            }
            await tx.query(`INSERT INTO social_family_uses VALUES ($1, $2, ${TODAY}, $3)`, [payer, feature, operation]);
            await tx.query("UPDATE social_family SET reputation=reputation-$2 WHERE character_id=$1", [payer, COSTS[feature]]);
            if (!moving) {
                const short = feature === FamilyFeature.FAMILY_FEATURE_DROP_15_MINUTES || feature === FamilyFeature.FAMILY_FEATURE_EXP_15_MINUTES;
                await tx.query(
                    `UPDATE social_family SET buff_feature=$2, buff_until=NOW()+($3::int*INTERVAL '1 minute')
                    WHERE character_id = ANY($1::int[])`,
                    [beneficiaries, feature, short ? 15 : 30],
                );
            }
            return FamilyResult.FAMILY_RESULT_OK;
        });
        return this.reply(world, status, id, payer);
    }

    // A refund writes its tombstone before looking for the reservation, so a
    // reservation that arrives later is rejected.
    async finishMove(req: FinishFamilyMoveRequest): Promise<FamilyReply> {
        const world = req.worldId;
        const id = req.characterId;
        if (req.requestId === "0") {
            return { status: FamilyResult.FAMILY_RESULT_FAILED, family: undefined };
        }

        if (req.complete) {
            await this.internalContext.withPgGlobalTransaction(world, async tx => {
                await this.familyMoveRequestRepository.transition(world, id, req.requestId, RequestStatus.Applied, RequestStatus.Completed, tx);
            });
            return { status: FamilyResult.FAMILY_RESULT_OK, family: undefined };
        }

        const payer = await this.internalContext.withPgGlobalTransaction(world, async tx => {
            await this.familyMoveRequestRepository.abort(world, id, req.requestId, tx);
            const reservation = await this.familyMoveRequestRepository.transition(world, id, req.requestId, RequestStatus.Applied, RequestStatus.Aborted, tx);
            if (!reservation) {
                return 0;
            }
            const payer = reservation.outcome.payer_id;
            await tx.query("DELETE FROM social_family_uses WHERE operation_id=$1", [req.requestId]);
            await tx.query(
                "UPDATE social_family SET reputation=LEAST(2147483647, reputation::bigint+$2) WHERE character_id=$1",
                [payer, COSTS[reservation.outcome.feature]],
            );
            return payer;
        });
        return this.reply(world, FamilyResult.FAMILY_RESULT_OK, id, payer);
    }
}
