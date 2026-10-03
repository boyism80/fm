import {
    Alliance,
    Guild,
    AllianceErrorCode,
    GuildMemberRank,
    type Alliance as AllianceMessage,
    type Guild as GuildMessage,
} from "../protobuf/generated/fminternal/internal_service";
import type { PoolClient } from "pg";
import { AppConfiguration } from "../config/app-configuration";
import { InternalContext } from "../context/internal-context";
import { AllianceRepository } from "../repos/alliance-repository";
import { GuildRepository } from "../repos/guild-repository";
import { GuildMemberRepository } from "../repos/guild-member-repository";
import { CharacterRealtimeStateRepository } from "../repos/character-realtime-state-repository";
import type { AllianceModel } from "../repos/alliance-repository";
import { RabbitMQService } from "./rabbitmq-service";
import { DistributedLockService } from "./distributed-lock-service";
import { DistributedLockGuard, DistributedLockMultiGuard } from "../system/distributed-lock";
import { GuildService } from "./guild-service";
import {
    ALLIANCE_CAPACITY_MAX,
    DEFAULT_ALLIANCE_CAPACITY,
    DEFAULT_ALLIANCE_RANK_TITLES,
    type AllianceRankTitles,
} from "../types/alliance-json";

const messages = { AllianceErrorCode };

const AMQ_DIRECT_EXCHANGE = "amq.direct";

const ALLIANCE_EVT = {
    CREATED: "created",
    DISBANDED: "disbanded",
    GUILD_LEFT: "guild_left",
    GUILD_ADDED: "guild_added",
    CAPACITY_CHANGED: "capacity_changed",
    RANK_TITLES_CHANGED: "rank_titles_changed",
    MEMBER_RANK_CHANGED: "member_rank_changed",
    LEADER_CHANGED: "leader_changed",
    NOTICE_CHANGED: "notice_changed",
} as const;

const MAX_ALLIANCE_NOTICE_LEN = 100;
const MIN_ALLIANCE_NAME_LEN = 3;
const MAX_ALLIANCE_NAME_LEN = 12;

export type CreateAllianceResult = {
    ok: boolean;
    code?: AllianceErrorCode;
    allianceId?: number;
    revision?: number;
    alliance?: AllianceMessage;
};

export type GetAllianceResult = { found: boolean; alliance?: AllianceMessage };

export type DisbandAllianceResult = {
    ok: boolean;
    code?: AllianceErrorCode;
    allianceId?: number;
    revision?: number;
};

export type LeaveAllianceResult = {
    ok: boolean;
    code?: AllianceErrorCode;
    allianceId?: number;
    revision?: number;
    disbanded?: boolean;
    removedGuildId?: number;
    alliance?: AllianceMessage;
};

export type ExpelAllianceGuildResult = {
    ok: boolean;
    code?: AllianceErrorCode;
    allianceId?: number;
    revision?: number;
    disbanded?: boolean;
    removedGuildId?: number;
    alliance?: AllianceMessage;
};

export type AcceptAllianceInviteResult = {
    ok: boolean;
    code?: AllianceErrorCode;
    allianceId?: number;
    revision?: number;
    alliance?: AllianceMessage;
};

export type IncreaseAllianceCapacityResult = {
    ok: boolean;
    code?: AllianceErrorCode;
    allianceId?: number;
    revision?: number;
    alliance?: AllianceMessage;
};

export type ChangeAllianceRankTitlesResult = {
    ok: boolean;
    code?: AllianceErrorCode;
    allianceId?: number;
    revision?: number;
    alliance?: AllianceMessage;
};

export type ChangeAllianceMemberRankResult = {
    ok: boolean;
    code?: AllianceErrorCode;
    allianceId?: number;
    revision?: number;
    targetCharacterId?: number;
    newAllianceRank?: number;
    alliance?: AllianceMessage;
};

export type ChangeAllianceLeaderResult = {
    ok: boolean;
    code?: AllianceErrorCode;
    allianceId?: number;
    revision?: number;
    oldLeaderCharacterId?: number;
    newLeaderCharacterId?: number;
    alliance?: AllianceMessage;
};

export type ChangeAllianceNoticeResult = {
    ok: boolean;
    code?: AllianceErrorCode;
    allianceId?: number;
    revision?: number;
    alliance?: AllianceMessage;
};

export class AllianceService {
    private readonly ctx: InternalContext;
    private readonly app: AppConfiguration;
    private readonly allianceRepo: AllianceRepository;
    private readonly guildRepo: GuildRepository;
    private readonly guildMemberRepo: GuildMemberRepository;
    private readonly characterRealtimeStateRepo: CharacterRealtimeStateRepository;
    private readonly rabbitmqService: RabbitMQService;
    private readonly distributedLockService: DistributedLockService;
    private readonly guildService: GuildService;

    constructor(
        internalContext: InternalContext,
        appConfiguration: AppConfiguration,
        allianceRepository: AllianceRepository,
        guildRepository: GuildRepository,
        guildMemberRepository: GuildMemberRepository,
        characterRealtimeStateRepository: CharacterRealtimeStateRepository,
        rabbitmqService: RabbitMQService,
        distributedLockService: DistributedLockService,
        guildService: GuildService
    ) {
        this.ctx = internalContext;
        this.app = appConfiguration;
        this.allianceRepo = allianceRepository;
        this.guildRepo = guildRepository;
        this.guildMemberRepo = guildMemberRepository;
        this.characterRealtimeStateRepo = characterRealtimeStateRepository;
        this.rabbitmqService = rabbitmqService;
        this.distributedLockService = distributedLockService;
        this.guildService = guildService;
    }

    private assertWorld(worldId: number) {
        const wid = String(worldId);
        if (!this.app.postgresql.worlds[wid]) {
            const err = new Error(`Unknown world_id: ${worldId}`) as Error & { code?: string };
            err.code = "UNKNOWN_WORLD";
            throw err;
        }
    }

    private assertCharacterId(characterId: number) {
        const n = characterId;
        if (!Number.isInteger(n) || n <= 0 || n > 0xffffffff) {
            const err = new Error("character_id must be a positive uint32") as Error & { code?: string };
            err.code = "INVALID_CHARACTER_ID";
            throw err;
        }
    }

    private validateAllianceName(name: string): AllianceErrorCode | null {
        if (typeof name !== "string") {
            return messages.AllianceErrorCode.ALLIANCE_ERROR_ALLIANCE_NAME_INVALID;
        }
        const trimmed = name.trim();
        if (trimmed.length < MIN_ALLIANCE_NAME_LEN || trimmed.length > MAX_ALLIANCE_NAME_LEN) {
            return messages.AllianceErrorCode.ALLIANCE_ERROR_ALLIANCE_NAME_INVALID;
        }
        return null;
    }

    private async publishToAllianceRoutes(
        eventType: string,
        worldId: number,
        allianceId: number,
        revision: number,
        extraPayload: Record<string, unknown> = {}
    ) {
        await this.rabbitmqService.assertDirectExchange(AMQ_DIRECT_EXCHANGE);
        const routingKey = `fm.${worldId}.all.alliance`;
        return this.rabbitmqService.publish(AMQ_DIRECT_EXCHANGE, routingKey, eventType, {
            event_id: extraPayload.event_id ?? `${Date.now()}-${Math.random().toString(36).slice(2, 8)}`,
            world_id: worldId,
            alliance_id: allianceId,
            revision,
            occurred_at: new Date().toISOString(),
            ...extraPayload,
        });
    }

    private async allianceToPb(worldId: number, alliance: AllianceModel, fallbackGuildIds?: number[]): Promise<AllianceMessage> {
        let guildIds = alliance.guildIds;
        if (guildIds.length === 0 && fallbackGuildIds != null && fallbackGuildIds.length > 0) {
            guildIds = fallbackGuildIds;
        }
        const guilds: GuildMessage[] = [];
        for (const guildId of guildIds) {
            const loaded = await this.guildService.getGuild(worldId, guildId);
            if (loaded.guild) {
                guilds.push(await this.guildService.guildToPb(worldId, loaded.guild, loaded.members ?? []));
            }
        }
        return {
            worldId,
            allianceId: alliance.allianceId,
            name: alliance.name,
            leaderCharacterId: alliance.leaderCharacterId,
            revision: alliance.revision,
            capacity: alliance.capacity,
            notice: alliance.notice,
            rankTitles: [...alliance.rankTitles],
            guildIds,
            guilds,
        };
    }

    async getAlliance(worldId: number, allianceId: number): Promise<GetAllianceResult> {
        this.assertWorld(worldId);
        if (!Number.isInteger(allianceId) || allianceId < 1) {
            return { found: false };
        }
        const alliance = await this.allianceRepo.get(worldId, allianceId);
        if (!alliance) {
            return { found: false };
        }
        return { found: true, alliance: await this.allianceToPb(worldId, alliance) };
    }

    private async assignAllianceRanks(
        worldId: number,
        guildId: number,
        masterAllianceRank: number,
        dataTx: PoolClient
    ) {
        const members = [...(await this.guildMemberRepo.getAll(worldId, String(guildId), { txClient: dataTx })).values()];
        if (members.length === 0) {
            return;
        }
        const updated = members.map((member) => {
            const allianceRank =
                member.guildRank === GuildMemberRank.GUILD_MEMBER_RANK_MASTER
                    ? masterAllianceRank
                    : 3;
            return {
                ...member,
                allianceRank,
            };
        });
        await this.guildMemberRepo.setAll(worldId, updated, { txClient: dataTx });
    }

    private async rollbackCreateAlliance(worldId: number, allianceId: number, guildIds: number[]) {
        await this.ctx.withPgDataTransaction(worldId, allianceId, async (dataTx: PoolClient) => {
            for (const guildId of guildIds) {
                const guild = await this.guildRepo.get(worldId, guildId, { txClient: dataTx });
                if (guild) {
                    await this.guildRepo.set(
                        worldId,
                        { ...guild, allianceId: null, revision: guild.revision + 1 },
                        { txClient: dataTx }
                    );
                }
                const members = [...(await this.guildMemberRepo.getAll(worldId, String(guildId), { txClient: dataTx })).values()];
                if (members.length > 0) {
                    const cleared = members.map((member) => ({
                        ...member,
                        allianceRank: null,
                    }));
                    await this.guildMemberRepo.setAll(worldId, cleared, { txClient: dataTx });
                }
            }
            await this.allianceRepo.delete({ worldId, allianceId }, { txClient: dataTx });
        }).catch(() => {});
        await this.allianceRepo.invalidateCache(worldId, allianceId).catch(() => {});
        for (const guildId of guildIds) {
            await this.guildRepo.invalidateCache(worldId, guildId).catch(() => {});
            await this.guildMemberRepo.invalidateCache(worldId, String(guildId)).catch(() => {});
        }
    }

    async createAlliance(
        worldId: number,
        allianceName: string,
        leaderCharacterId: number,
        partnerCharacterId: number
    ): Promise<CreateAllianceResult> {
        this.assertWorld(worldId);
        this.assertCharacterId(leaderCharacterId);
        this.assertCharacterId(partnerCharacterId);

        if (leaderCharacterId === partnerCharacterId) {
            return { ok: false, code: messages.AllianceErrorCode.ALLIANCE_ERROR_PARTNER_INVALID };
        }

        const nameError = this.validateAllianceName(allianceName);
        if (nameError != null) {
            return { ok: false, code: nameError };
        }
        const trimmedName = allianceName.trim();

        await using _leaderLock = await this.distributedLockService.acquireWorldDataLock(
            worldId,
            `character_realtime:${leaderCharacterId}`,
        );
        await using _partnerLock = await this.distributedLockService.acquireWorldDataLock(
            worldId,
            `character_realtime:${partnerCharacterId}`,
        );

        const leaderState = await this.characterRealtimeStateRepo.get(worldId, leaderCharacterId);
        const partnerState = await this.characterRealtimeStateRepo.get(worldId, partnerCharacterId);
        const guildId = leaderState?.guildId;
        const partnerGuildId = partnerState?.guildId;
        if (guildId == null || partnerGuildId == null) {
            return { ok: false, code: messages.AllianceErrorCode.ALLIANCE_ERROR_GUILD_NOT_FOUND };
        }
        if (guildId === partnerGuildId) {
            return { ok: false, code: messages.AllianceErrorCode.ALLIANCE_ERROR_PARTNER_INVALID };
        }

        await using _guildLock1 = await this.distributedLockService.acquireWorldDataLock(worldId, `guild:${guildId}`);
        await using _guildLock2 = await this.distributedLockService.acquireWorldDataLock(
            worldId,
            `guild:${partnerGuildId}`,
        );

        const guild1 = await this.guildRepo.get(worldId, guildId);
        const guild2 = await this.guildRepo.get(worldId, partnerGuildId);
        if (!guild1 || !guild2) {
            return { ok: false, code: messages.AllianceErrorCode.ALLIANCE_ERROR_GUILD_NOT_FOUND };
        }
        if (guild1.allianceId != null || guild2.allianceId != null) {
            return { ok: false, code: messages.AllianceErrorCode.ALLIANCE_ERROR_GUILD_ALREADY_IN_ALLIANCE };
        }
        if (guild1.leaderCharacterId !== leaderCharacterId) {
            return { ok: false, code: messages.AllianceErrorCode.ALLIANCE_ERROR_NOT_GUILD_MASTER };
        }

        const guild1Members = [...(await this.guildMemberRepo.getAll(worldId, String(guildId))).values()];
        const leaderMember = guild1Members.find((m) => m.characterId === leaderCharacterId);
        if (!leaderMember || leaderMember.guildRank !== GuildMemberRank.GUILD_MEMBER_RANK_MASTER) {
            return { ok: false, code: messages.AllianceErrorCode.ALLIANCE_ERROR_NOT_GUILD_MASTER };
        }

        const guild2Members = [...(await this.guildMemberRepo.getAll(worldId, String(partnerGuildId))).values()];
        if (guild2.leaderCharacterId !== partnerCharacterId) {
            return { ok: false, code: messages.AllianceErrorCode.ALLIANCE_ERROR_PARTNER_INVALID };
        }
        const partnerMember = guild2Members.find((m) => m.characterId === partnerCharacterId);
        if (!partnerMember || partnerMember.guildRank !== GuildMemberRank.GUILD_MEMBER_RANK_MASTER) {
            return { ok: false, code: messages.AllianceErrorCode.ALLIANCE_ERROR_PARTNER_INVALID };
        }

        const initialGuildIds: number[] = [guildId, partnerGuildId];
        let savedAlliance: AllianceModel;
        let allianceId: number | null = null;

        try {
            const newAllianceId = await this.ctx.withPgGlobalTransaction(worldId, async (txClient: PoolClient) => {
                const res = await txClient.query("SELECT nextval('alliance_id_seq') AS id");
                const id = res.rows?.[0]?.id;
                if (id == null) {
                    const err = new Error("nextval('alliance_id_seq') returned no id") as Error & { code?: string };
                    err.code = "ALLIANCE_ID_SEQ_ERROR";
                    throw err;
                }
                const n = Number(id);
                if (!Number.isInteger(n) || n < 1) {
                    const err = new Error("alliance_id_seq returned invalid id") as Error & { code?: string };
                    err.code = "ALLIANCE_ID_SEQ_ERROR";
                    throw err;
                }
                return n;
            });
            allianceId = newAllianceId;

            savedAlliance = await this.ctx.withPgDataTransaction(worldId, newAllianceId, async (dataTx: PoolClient) => {
                return this.allianceRepo.set(
                    worldId,
                    {
                        worldId,
                        allianceId: newAllianceId,
                        name: trimmedName,
                        leaderCharacterId,
                        guildIds: initialGuildIds,
                        rankTitles: [...DEFAULT_ALLIANCE_RANK_TITLES],
                        capacity: DEFAULT_ALLIANCE_CAPACITY,
                        notice: "",
                        revision: 1,
                    },
                    { txClient: dataTx }
                );
            });

            await this.ctx.withPgDataTransaction(worldId, guildId, async (dataTx: PoolClient) => {
                await this.guildRepo.set(
                    worldId,
                    { ...guild1, allianceId: newAllianceId, revision: guild1.revision + 1 },
                    { txClient: dataTx }
                );
                await this.assignAllianceRanks(worldId, guildId, 1, dataTx);
            });

            await this.ctx.withPgDataTransaction(worldId, partnerGuildId, async (dataTx: PoolClient) => {
                await this.guildRepo.set(
                    worldId,
                    { ...guild2, allianceId: newAllianceId, revision: guild2.revision + 1 },
                    { txClient: dataTx }
                );
                await this.assignAllianceRanks(worldId, partnerGuildId, 2, dataTx);
            });
        } catch (err: unknown) {
            const code = (err as { code?: string })?.code;
            if (code === "23505") {
                return { ok: false, code: messages.AllianceErrorCode.ALLIANCE_ERROR_ALLIANCE_NAME_TAKEN };
            }
            if (allianceId != null) {
                await this.rollbackCreateAlliance(worldId, allianceId, [guildId, partnerGuildId]);
            }
            throw err;
        }

        const createdAllianceId = savedAlliance.allianceId;

        await this.allianceRepo.invalidateCache(worldId, createdAllianceId);
        await this.guildRepo.invalidateCache(worldId, guildId);
        await this.guildRepo.invalidateCache(worldId, partnerGuildId);
        await this.guildMemberRepo.invalidateCache(worldId, String(guildId));
        await this.guildMemberRepo.invalidateCache(worldId, String(partnerGuildId));

        const allianceMessage = await this.allianceToPb(worldId, savedAlliance, initialGuildIds);
        const wire = Alliance.encode(allianceMessage).finish();
        await this.publishToAllianceRoutes(ALLIANCE_EVT.CREATED, worldId, createdAllianceId, savedAlliance.revision, {
            alliance_name: trimmedName,
            guild_ids: [guildId, partnerGuildId],
            leader_character_id: leaderCharacterId,
            alliance_pb: Buffer.from(wire).toString("base64"),
        });

        return {
            ok: true,
            allianceId: createdAllianceId,
            revision: savedAlliance.revision,
            alliance: allianceMessage,
        };
    }

    private async dissolveAllianceInTransaction(
        worldId: number,
        allianceId: number,
        dataTx: PoolClient
    ): Promise<
        | {
              ok: true;
              allianceId: number;
              revision: number;
              guildIds: number[];
              memberCharacterIds: number[];
          }
        | { ok: false; code: AllianceErrorCode }
    > {
        const lockedAlliance = await this.allianceRepo.get(worldId, allianceId, { txClient: dataTx });
        if (!lockedAlliance) {
            return { ok: false, code: messages.AllianceErrorCode.ALLIANCE_ERROR_ALLIANCE_NOT_FOUND };
        }

        const memberCharacterIds: number[] = [];
        const nextRevision = lockedAlliance.revision + 1;

        for (const gid of lockedAlliance.guildIds) {
            if (!Number.isInteger(gid) || gid < 0) {
                continue;
            }
            const g = await this.guildRepo.get(worldId, gid, { txClient: dataTx });
            if (!g) {
                continue;
            }
            await this.guildRepo.set(
                worldId,
                { ...g, allianceId: null, revision: g.revision + 1 },
                { txClient: dataTx }
            );
            const members = [...(await this.guildMemberRepo.getAll(worldId, String(gid), { txClient: dataTx })).values()];
            if (members.length > 0) {
                const cleared = members.map((member) => ({
                    ...member,
                    allianceRank: null,
                }));
                await this.guildMemberRepo.setAll(worldId, cleared, { txClient: dataTx });
                for (const member of cleared) {
                    memberCharacterIds.push(member.characterId);
                }
            }
        }

        await this.allianceRepo.set(
            worldId,
            {
                ...lockedAlliance,
                revision: nextRevision,
                disbandedAt: new Date(),
            },
            { txClient: dataTx }
        );

        return {
            ok: true,
            allianceId,
            revision: nextRevision,
            guildIds: lockedAlliance.guildIds,
            memberCharacterIds,
        };
    }

    async acceptAllianceInvite(
        worldId: number,
        characterId: number,
        allianceId: number,
        guildId: number
    ): Promise<AcceptAllianceInviteResult> {
        this.assertWorld(worldId);
        this.assertCharacterId(characterId);
        if (!Number.isInteger(allianceId) || allianceId < 1) {
            return { ok: false, code: messages.AllianceErrorCode.ALLIANCE_ERROR_ALLIANCE_NOT_FOUND };
        }
        if (!Number.isInteger(guildId) || guildId < 1) {
            return { ok: false, code: messages.AllianceErrorCode.ALLIANCE_ERROR_GUILD_NOT_FOUND };
        }

        await using _characterRealtimeLock = await this.distributedLockService.acquireWorldDataLock(
            worldId,
            `character_realtime:${characterId}`,
        );

        const state = await this.characterRealtimeStateRepo.get(worldId, characterId);
        if (state?.guildId == null || state.guildId !== guildId) {
            return { ok: false, code: messages.AllianceErrorCode.ALLIANCE_ERROR_GUILD_NOT_FOUND };
        }

        await using _guildLock = await this.distributedLockService.acquireWorldDataLock(
            worldId,
            `guild:${guildId}`,
        );

        const guild = await this.guildRepo.get(worldId, guildId);
        if (!guild) {
            return { ok: false, code: messages.AllianceErrorCode.ALLIANCE_ERROR_GUILD_NOT_FOUND };
        }
        if (guild.allianceId != null && guild.allianceId >= 1) {
            return { ok: false, code: messages.AllianceErrorCode.ALLIANCE_ERROR_GUILD_ALREADY_IN_ALLIANCE };
        }
        if (guild.leaderCharacterId !== characterId) {
            return { ok: false, code: messages.AllianceErrorCode.ALLIANCE_ERROR_NOT_GUILD_MASTER };
        }

        const members = [...(await this.guildMemberRepo.getAll(worldId, String(guildId))).values()];
        const requesterMember = members.find((m) => m.characterId === characterId);
        if (!requesterMember || requesterMember.guildRank !== GuildMemberRank.GUILD_MEMBER_RANK_MASTER) {
            return { ok: false, code: messages.AllianceErrorCode.ALLIANCE_ERROR_NOT_GUILD_MASTER };
        }

        await using _allianceLock = await this.distributedLockService.acquireWorldDataLock(
            worldId,
            `alliance:${allianceId}`,
        );

        const alliance = await this.allianceRepo.get(worldId, allianceId);
        if (!alliance) {
            return { ok: false, code: messages.AllianceErrorCode.ALLIANCE_ERROR_ALLIANCE_NOT_FOUND };
        }

        const otherGuildIds = [...alliance.guildIds].sort((a, b) => a - b);
        await using _otherGuildLocks = await this.distributedLockService.acquireWorldDataLocks(
            worldId,
            otherGuildIds.map((id) => `guild:${id}`),
        );

        const txResult = await this.ctx.withPgDataTransaction(worldId, allianceId, async (dataTx: PoolClient) => {
            const lockedAlliance = await this.allianceRepo.get(worldId, allianceId, { txClient: dataTx });
            if (!lockedAlliance) {
                return { ok: false as const, code: messages.AllianceErrorCode.ALLIANCE_ERROR_ALLIANCE_NOT_FOUND };
            }
            if (lockedAlliance.guildIds.includes(guildId)) {
                return { ok: false as const, code: messages.AllianceErrorCode.ALLIANCE_ERROR_GUILD_ALREADY_IN_ALLIANCE };
            }
            if (lockedAlliance.guildIds.length >= lockedAlliance.capacity) {
                return { ok: false as const, code: messages.AllianceErrorCode.ALLIANCE_ERROR_CAPACITY_FULL };
            }

            const lockedGuild = await this.guildRepo.get(worldId, guildId, { txClient: dataTx });
            if (!lockedGuild || lockedGuild.allianceId != null) {
                return { ok: false as const, code: messages.AllianceErrorCode.ALLIANCE_ERROR_GUILD_ALREADY_IN_ALLIANCE };
            }

            const nextRevision = lockedAlliance.revision + 1;
            const nextGuildIds = [...lockedAlliance.guildIds, guildId];

            await this.guildRepo.set(
                worldId,
                { ...lockedGuild, allianceId, revision: lockedGuild.revision + 1 },
                { txClient: dataTx }
            );
            await this.assignAllianceRanks(worldId, guildId, 3, dataTx);

            const savedAlliance = await this.allianceRepo.set(
                worldId,
                {
                    ...lockedAlliance,
                    guildIds: nextGuildIds,
                    revision: nextRevision,
                },
                { txClient: dataTx }
            );

            return {
                ok: true as const,
                allianceId,
                revision: nextRevision,
                savedAlliance,
            };
        });

        if (!txResult.ok || txResult.allianceId == null || txResult.revision == null) {
            return txResult;
        }

        await this.allianceRepo.invalidateCache(worldId, txResult.allianceId);
        await this.guildRepo.invalidateCache(worldId, guildId);
        await this.guildMemberRepo.invalidateCache(worldId, String(guildId));
        for (const gid of otherGuildIds) {
            await this.guildRepo.invalidateCache(worldId, gid);
            await this.guildMemberRepo.invalidateCache(worldId, String(gid));
        }

        const allianceMessage = await this.allianceToPb(worldId, txResult.savedAlliance);

        const wireAlliance = Alliance.encode(allianceMessage).finish();
        await this.publishToAllianceRoutes(
            ALLIANCE_EVT.GUILD_ADDED,
            worldId,
            txResult.allianceId,
            txResult.revision,
            {
                added_guild_id: guildId,
                alliance_pb: Buffer.from(wireAlliance).toString("base64"),
            }
        );

        return {
            ok: true,
            allianceId: txResult.allianceId,
            revision: txResult.revision,
            alliance: allianceMessage,
        };
    }

    private async lockRequesterAlliance(worldId: number, characterId: number, otherCharacterIds: number[] = []) {
        const guards: DistributedLockGuard[] = [];
        const fail = async (code: AllianceErrorCode) => {
            await new DistributedLockMultiGuard(guards).release();
            return { ok: false as const, code };
        };
        try {
            const characterIds = [...new Set([characterId, ...otherCharacterIds])].sort((a, b) => a - b);
            for (const id of characterIds) {
                guards.push(await this.distributedLockService.acquireWorldDataLock(worldId, `character_realtime:${id}`));
            }

            const state = await this.characterRealtimeStateRepo.get(worldId, characterId);
            const guildId = state?.guildId;
            if (guildId == null || !Number.isInteger(guildId) || guildId < 1) {
                return await fail(messages.AllianceErrorCode.ALLIANCE_ERROR_GUILD_NOT_FOUND);
            }
            guards.push(await this.distributedLockService.acquireWorldDataLock(worldId, `guild:${guildId}`));
            const guild = await this.guildRepo.get(worldId, guildId);
            if (!guild) {
                return await fail(messages.AllianceErrorCode.ALLIANCE_ERROR_GUILD_NOT_FOUND);
            }
            const allianceId = guild.allianceId;
            if (allianceId == null || allianceId < 1) {
                return await fail(messages.AllianceErrorCode.ALLIANCE_ERROR_NOT_IN_ALLIANCE);
            }

            guards.push(await this.distributedLockService.acquireWorldDataLock(worldId, `alliance:${allianceId}`));
            const alliance = await this.allianceRepo.get(worldId, allianceId);
            if (!alliance) {
                return await fail(messages.AllianceErrorCode.ALLIANCE_ERROR_ALLIANCE_NOT_FOUND);
            }

            const members = [...(await this.guildMemberRepo.getAll(worldId, String(guildId))).values()];
            return {
                ok: true as const,
                guildId,
                guild,
                allianceId,
                alliance,
                requester: members.find((m) => m.characterId === characterId),
                locks: new DistributedLockMultiGuard(guards),
            };
        } catch (err) {
            await new DistributedLockMultiGuard(guards).release();
            throw err;
        }
    }

    async increaseAllianceCapacity(
        worldId: number,
        characterId: number
    ): Promise<IncreaseAllianceCapacityResult> {
        this.assertWorld(worldId);
        this.assertCharacterId(characterId);

        const locked = await this.lockRequesterAlliance(worldId, characterId);
        if (!locked.ok) {
            return locked;
        }
        await using _locks = locked.locks;
        const { guild, allianceId, alliance, requester: requesterMember } = locked;
        if (guild.leaderCharacterId !== characterId) {
            return { ok: false, code: messages.AllianceErrorCode.ALLIANCE_ERROR_NOT_GUILD_MASTER };
        }
        if (alliance.leaderCharacterId !== characterId) {
            return { ok: false, code: messages.AllianceErrorCode.ALLIANCE_ERROR_NOT_ALLIANCE_LEADER };
        }
        if (!requesterMember || requesterMember.guildRank !== GuildMemberRank.GUILD_MEMBER_RANK_MASTER) {
            return { ok: false, code: messages.AllianceErrorCode.ALLIANCE_ERROR_NOT_GUILD_MASTER };
        }
        if (requesterMember.allianceRank !== 1) {
            return { ok: false, code: messages.AllianceErrorCode.ALLIANCE_ERROR_NOT_ALLIANCE_LEADER };
        }

        const txResult = await this.ctx.withPgDataTransaction(worldId, allianceId, async (dataTx: PoolClient) => {
            const lockedAlliance = await this.allianceRepo.get(worldId, allianceId, { txClient: dataTx });
            if (!lockedAlliance) {
                return { ok: false as const, code: messages.AllianceErrorCode.ALLIANCE_ERROR_ALLIANCE_NOT_FOUND };
            }
            if (lockedAlliance.leaderCharacterId !== characterId) {
                return { ok: false as const, code: messages.AllianceErrorCode.ALLIANCE_ERROR_NOT_ALLIANCE_LEADER };
            }
            if (lockedAlliance.capacity >= ALLIANCE_CAPACITY_MAX) {
                return { ok: false as const, code: messages.AllianceErrorCode.ALLIANCE_ERROR_CAPACITY_MAX };
            }

            const nextRevision = lockedAlliance.revision + 1;
            const savedAlliance = await this.allianceRepo.set(
                worldId,
                {
                    ...lockedAlliance,
                    capacity: lockedAlliance.capacity + 1,
                    revision: nextRevision,
                },
                { txClient: dataTx }
            );

            return {
                ok: true as const,
                allianceId,
                revision: nextRevision,
                savedAlliance,
            };
        });

        if (!txResult.ok || txResult.allianceId == null || txResult.revision == null) {
            return txResult;
        }

        await this.allianceRepo.invalidateCache(worldId, txResult.allianceId);

        const allianceMessage = await this.allianceToPb(worldId, txResult.savedAlliance);
        const wireAlliance = Alliance.encode(allianceMessage).finish();
        await this.publishToAllianceRoutes(
            ALLIANCE_EVT.CAPACITY_CHANGED,
            worldId,
            txResult.allianceId,
            txResult.revision,
            {
                alliance_pb: Buffer.from(wireAlliance).toString("base64"),
            }
        );

        return {
            ok: true,
            allianceId: txResult.allianceId,
            revision: txResult.revision,
            alliance: allianceMessage,
        };
    }

    private normalizeAllianceRankTitles(rankTitles: string[] | null | undefined): AllianceRankTitles | null {
        if (!rankTitles || rankTitles.length !== 5) {
            return null;
        }
        return rankTitles.map((entry) => String(entry ?? "")) as AllianceRankTitles;
    }

    private canChangeAllianceMemberRank(allianceRank: number | null | undefined) {
        return allianceRank != null && allianceRank >= 1 && allianceRank <= 2;
    }

    async changeAllianceRankTitles(
        worldId: number,
        characterId: number,
        rankTitles: string[]
    ): Promise<ChangeAllianceRankTitlesResult> {
        this.assertWorld(worldId);
        this.assertCharacterId(characterId);

        const normalizedRankTitles = this.normalizeAllianceRankTitles(rankTitles);
        if (!normalizedRankTitles) {
            return { ok: false, code: messages.AllianceErrorCode.ALLIANCE_ERROR_INVALID_RANK_TITLES };
        }

        const locked = await this.lockRequesterAlliance(worldId, characterId);
        if (!locked.ok) {
            return locked;
        }
        await using _locks = locked.locks;
        const { allianceId, alliance, requester: requesterMember } = locked;
        if (alliance.leaderCharacterId !== characterId) {
            return { ok: false, code: messages.AllianceErrorCode.ALLIANCE_ERROR_NOT_ALLIANCE_LEADER };
        }
        if (!requesterMember || requesterMember.guildRank !== GuildMemberRank.GUILD_MEMBER_RANK_MASTER) {
            return { ok: false, code: messages.AllianceErrorCode.ALLIANCE_ERROR_NOT_GUILD_MASTER };
        }
        if (requesterMember.allianceRank !== 1) {
            return { ok: false, code: messages.AllianceErrorCode.ALLIANCE_ERROR_NOT_ALLIANCE_LEADER };
        }

        const txResult = await this.ctx.withPgDataTransaction(worldId, allianceId, async (dataTx: PoolClient) => {
            const lockedAlliance = await this.allianceRepo.get(worldId, allianceId, { txClient: dataTx });
            if (!lockedAlliance) {
                return { ok: false as const, code: messages.AllianceErrorCode.ALLIANCE_ERROR_ALLIANCE_NOT_FOUND };
            }
            if (lockedAlliance.leaderCharacterId !== characterId) {
                return { ok: false as const, code: messages.AllianceErrorCode.ALLIANCE_ERROR_NOT_ALLIANCE_LEADER };
            }

            const nextRevision = lockedAlliance.revision + 1;
            const savedAlliance = await this.allianceRepo.set(
                worldId,
                {
                    ...lockedAlliance,
                    rankTitles: normalizedRankTitles,
                    revision: nextRevision,
                },
                { txClient: dataTx }
            );

            return {
                ok: true as const,
                allianceId,
                revision: nextRevision,
                savedAlliance,
            };
        });

        if (!txResult.ok || txResult.allianceId == null || txResult.revision == null) {
            return txResult;
        }

        await this.allianceRepo.invalidateCache(worldId, txResult.allianceId);

        const allianceMessage = await this.allianceToPb(worldId, txResult.savedAlliance);
        const wireAlliance = Alliance.encode(allianceMessage).finish();
        await this.publishToAllianceRoutes(
            ALLIANCE_EVT.RANK_TITLES_CHANGED,
            worldId,
            txResult.allianceId,
            txResult.revision,
            {
                alliance_pb: Buffer.from(wireAlliance).toString("base64"),
            }
        );

        return {
            ok: true,
            allianceId: txResult.allianceId,
            revision: txResult.revision,
            alliance: allianceMessage,
        };
    }

    async changeAllianceNotice(
        worldId: number,
        characterId: number,
        notice: string
    ): Promise<ChangeAllianceNoticeResult> {
        this.assertWorld(worldId);
        this.assertCharacterId(characterId);

        if (notice.length > MAX_ALLIANCE_NOTICE_LEN) {
            return { ok: false, code: messages.AllianceErrorCode.ALLIANCE_ERROR_INVALID_NOTICE };
        }

        const locked = await this.lockRequesterAlliance(worldId, characterId);
        if (!locked.ok) {
            return locked;
        }
        await using _locks = locked.locks;
        const { guildId, guild, allianceId, alliance, requester: requesterMember } = locked;
        if (!requesterMember || requesterMember.guildRank !== GuildMemberRank.GUILD_MEMBER_RANK_MASTER) {
            return { ok: false, code: messages.AllianceErrorCode.ALLIANCE_ERROR_NOT_GUILD_MASTER };
        }
        if (requesterMember.allianceRank == null || requesterMember.allianceRank > 2) {
            return { ok: false, code: messages.AllianceErrorCode.ALLIANCE_ERROR_NOT_AUTHORIZED };
        }

        const txResult = await this.ctx.withPgDataTransaction(worldId, allianceId, async (dataTx: PoolClient) => {
            const lockedAlliance = await this.allianceRepo.get(worldId, allianceId, { txClient: dataTx });
            if (!lockedAlliance) {
                return { ok: false as const, code: messages.AllianceErrorCode.ALLIANCE_ERROR_ALLIANCE_NOT_FOUND };
            }

            const nextRevision = lockedAlliance.revision + 1;
            const savedAlliance = await this.allianceRepo.set(
                worldId,
                {
                    ...lockedAlliance,
                    notice,
                    revision: nextRevision,
                },
                { txClient: dataTx }
            );

            return {
                ok: true as const,
                allianceId,
                revision: nextRevision,
                savedAlliance,
            };
        });

        if (!txResult.ok || txResult.allianceId == null || txResult.revision == null) {
            return txResult;
        }

        await this.allianceRepo.invalidateCache(worldId, txResult.allianceId);

        const allianceMessage = await this.allianceToPb(worldId, txResult.savedAlliance);
        const wireAlliance = Alliance.encode(allianceMessage).finish();
        await this.publishToAllianceRoutes(
            ALLIANCE_EVT.NOTICE_CHANGED,
            worldId,
            txResult.allianceId,
            txResult.revision,
            {
                notice,
                alliance_pb: Buffer.from(wireAlliance).toString("base64"),
            }
        );

        return {
            ok: true,
            allianceId: txResult.allianceId,
            revision: txResult.revision,
            alliance: allianceMessage,
        };
    }

    async changeAllianceLeader(
        worldId: number,
        characterId: number,
        newLeaderCharacterId: number
    ): Promise<ChangeAllianceLeaderResult> {
        this.assertWorld(worldId);
        this.assertCharacterId(characterId);
        this.assertCharacterId(newLeaderCharacterId);

        if (newLeaderCharacterId === characterId) {
            return { ok: false, code: messages.AllianceErrorCode.ALLIANCE_ERROR_INVALID_LEADER_CANDIDATE };
        }

        const locked = await this.lockRequesterAlliance(worldId, characterId, [newLeaderCharacterId]);
        if (!locked.ok) {
            return locked;
        }
        await using _locks = locked.locks;
        const { guildId: requesterGuildId, allianceId, alliance, requester: requesterMember } = locked;
        if (alliance.leaderCharacterId !== characterId) {
            return { ok: false, code: messages.AllianceErrorCode.ALLIANCE_ERROR_NOT_ALLIANCE_LEADER };
        }
        if (!requesterMember || requesterMember.guildRank !== GuildMemberRank.GUILD_MEMBER_RANK_MASTER) {
            return { ok: false, code: messages.AllianceErrorCode.ALLIANCE_ERROR_NOT_GUILD_MASTER };
        }
        if (requesterMember.allianceRank !== 1) {
            return { ok: false, code: messages.AllianceErrorCode.ALLIANCE_ERROR_NOT_ALLIANCE_LEADER };
        }

        const otherGuildIds = alliance.guildIds
            .filter((id) => id !== requesterGuildId)
            .sort((a, b) => a - b);
        await using _otherGuildLocks = await this.distributedLockService.acquireWorldDataLocks(
            worldId,
            otherGuildIds.map((id) => `guild:${id}`),
        );

        const txResult = await this.ctx.withPgDataTransaction(worldId, allianceId, async (dataTx: PoolClient) => {
            const lockedAlliance = await this.allianceRepo.get(worldId, allianceId, { txClient: dataTx });
            if (!lockedAlliance) {
                return { ok: false as const, code: messages.AllianceErrorCode.ALLIANCE_ERROR_ALLIANCE_NOT_FOUND };
            }
            if (lockedAlliance.leaderCharacterId !== characterId) {
                return { ok: false as const, code: messages.AllianceErrorCode.ALLIANCE_ERROR_NOT_ALLIANCE_LEADER };
            }

            const oldLeaderCharacterId = lockedAlliance.leaderCharacterId;
            let oldLeaderGuildId: number | null = null;
            let newLeaderGuildId: number | null = null;

            for (const gid of lockedAlliance.guildIds) {
                const members = [
                    ...(await this.guildMemberRepo.getAll(worldId, String(gid), { txClient: dataTx })).values(),
                ];
                for (const member of members) {
                    if (member.characterId === oldLeaderCharacterId) {
                        oldLeaderGuildId = gid;
                    }
                    if (member.characterId === newLeaderCharacterId) {
                        newLeaderGuildId = gid;
                    }
                }
            }

            if (oldLeaderGuildId == null || newLeaderGuildId == null) {
                return { ok: false as const, code: messages.AllianceErrorCode.ALLIANCE_ERROR_INVALID_LEADER_CANDIDATE };
            }
            if (oldLeaderGuildId === newLeaderGuildId) {
                return { ok: false as const, code: messages.AllianceErrorCode.ALLIANCE_ERROR_INVALID_LEADER_CANDIDATE };
            }

            const newLeaderMembers = [
                ...(await this.guildMemberRepo.getAll(worldId, String(newLeaderGuildId), { txClient: dataTx })).values(),
            ];
            const newLeaderMember = newLeaderMembers.find((m) => m.characterId === newLeaderCharacterId);
            if (
                !newLeaderMember
                || newLeaderMember.guildRank !== GuildMemberRank.GUILD_MEMBER_RANK_MASTER
                || newLeaderMember.allianceRank !== 2
            ) {
                return { ok: false as const, code: messages.AllianceErrorCode.ALLIANCE_ERROR_INVALID_LEADER_CANDIDATE };
            }

            const oldLeaderMembers = [
                ...(await this.guildMemberRepo.getAll(worldId, String(oldLeaderGuildId), { txClient: dataTx })).values(),
            ];
            const oldLeaderMember = oldLeaderMembers.find((m) => m.characterId === oldLeaderCharacterId);
            if (
                !oldLeaderMember
                || oldLeaderMember.guildRank !== GuildMemberRank.GUILD_MEMBER_RANK_MASTER
                || oldLeaderMember.allianceRank !== 1
            ) {
                return { ok: false as const, code: messages.AllianceErrorCode.ALLIANCE_ERROR_INVALID_LEADER_CANDIDATE };
            }

            await this.guildMemberRepo.set(
                worldId,
                { ...oldLeaderMember, allianceRank: 2 },
                { txClient: dataTx }
            );
            await this.guildMemberRepo.set(
                worldId,
                { ...newLeaderMember, allianceRank: 1 },
                { txClient: dataTx }
            );

            const nextGuildIds = [...lockedAlliance.guildIds];
            const leaderGuildAt0 = nextGuildIds[0];
            const newLeaderIndex = nextGuildIds.indexOf(newLeaderGuildId);
            if (newLeaderIndex < 1 || leaderGuildAt0 !== oldLeaderGuildId) {
                return { ok: false as const, code: messages.AllianceErrorCode.ALLIANCE_ERROR_INVALID_LEADER_CANDIDATE };
            }
            nextGuildIds[0] = newLeaderGuildId;
            nextGuildIds[newLeaderIndex] = leaderGuildAt0;

            const nextRevision = lockedAlliance.revision + 1;
            const savedAlliance = await this.allianceRepo.set(
                worldId,
                {
                    ...lockedAlliance,
                    leaderCharacterId: newLeaderCharacterId,
                    guildIds: nextGuildIds,
                    revision: nextRevision,
                },
                { txClient: dataTx }
            );

            const oldGuild = await this.guildRepo.get(worldId, oldLeaderGuildId, { txClient: dataTx });
            if (oldGuild) {
                await this.guildRepo.set(
                    worldId,
                    { ...oldGuild, revision: oldGuild.revision + 1 },
                    { txClient: dataTx }
                );
            }
            const newGuild = await this.guildRepo.get(worldId, newLeaderGuildId, { txClient: dataTx });
            if (newGuild) {
                await this.guildRepo.set(
                    worldId,
                    { ...newGuild, revision: newGuild.revision + 1 },
                    { txClient: dataTx }
                );
            }

            return {
                ok: true as const,
                allianceId,
                revision: nextRevision,
                oldLeaderCharacterId,
                newLeaderCharacterId,
                savedAlliance,
                oldLeaderGuildId,
                newLeaderGuildId,
            };
        });

        if (!txResult.ok || txResult.allianceId == null || txResult.revision == null) {
            return txResult;
        }

        await this.allianceRepo.invalidateCache(worldId, txResult.allianceId);
        await this.guildRepo.invalidateCache(worldId, txResult.oldLeaderGuildId);
        await this.guildMemberRepo.invalidateCache(worldId, String(txResult.oldLeaderGuildId));
        await this.guildRepo.invalidateCache(worldId, txResult.newLeaderGuildId);
        await this.guildMemberRepo.invalidateCache(worldId, String(txResult.newLeaderGuildId));
        for (const gid of otherGuildIds) {
            if (gid === txResult.oldLeaderGuildId || gid === txResult.newLeaderGuildId) {
                continue;
            }
            await this.guildRepo.invalidateCache(worldId, gid);
            await this.guildMemberRepo.invalidateCache(worldId, String(gid));
        }

        const allianceMessage = await this.allianceToPb(worldId, txResult.savedAlliance);
        const wireAlliance = Alliance.encode(allianceMessage).finish();
        await this.publishToAllianceRoutes(
            ALLIANCE_EVT.LEADER_CHANGED,
            worldId,
            txResult.allianceId,
            txResult.revision,
            {
                old_leader_character_id: txResult.oldLeaderCharacterId,
                new_leader_character_id: txResult.newLeaderCharacterId,
                alliance_pb: Buffer.from(wireAlliance).toString("base64"),
            }
        );

        return {
            ok: true,
            allianceId: txResult.allianceId,
            revision: txResult.revision,
            oldLeaderCharacterId: txResult.oldLeaderCharacterId,
            newLeaderCharacterId: txResult.newLeaderCharacterId,
            alliance: allianceMessage,
        };
    }

    async changeAllianceMemberRank(
        worldId: number,
        requesterCharacterId: number,
        targetCharacterId: number,
        promote: boolean
    ): Promise<ChangeAllianceMemberRankResult> {
        this.assertWorld(worldId);
        this.assertCharacterId(requesterCharacterId);
        this.assertCharacterId(targetCharacterId);

        const locked = await this.lockRequesterAlliance(worldId, requesterCharacterId, [targetCharacterId]);
        if (!locked.ok) {
            return locked;
        }
        await using _locks = locked.locks;
        const { guildId: requesterGuildId, allianceId, alliance } = locked;

        const otherGuildIds = alliance.guildIds
            .filter((id) => id !== requesterGuildId)
            .sort((a, b) => a - b);
        await using _otherGuildLocks = await this.distributedLockService.acquireWorldDataLocks(
            worldId,
            otherGuildIds.map((id) => `guild:${id}`),
        );

        const txResult = await this.ctx.withPgDataTransaction(worldId, allianceId, async (dataTx: PoolClient) => {
            const requesterGuildLocked = await this.guildRepo.get(worldId, requesterGuildId, { txClient: dataTx });
            if (!requesterGuildLocked || requesterGuildLocked.allianceId !== allianceId) {
                return { ok: false as const, code: messages.AllianceErrorCode.ALLIANCE_ERROR_NOT_IN_ALLIANCE };
            }

            const requesterMembers = [
                ...(await this.guildMemberRepo.getAll(worldId, String(requesterGuildId), { txClient: dataTx })).values(),
            ];
            const requesterMember = requesterMembers.find((m) => m.characterId === requesterCharacterId);
            if (!requesterMember || requesterMember.guildRank !== GuildMemberRank.GUILD_MEMBER_RANK_MASTER) {
                return { ok: false as const, code: messages.AllianceErrorCode.ALLIANCE_ERROR_NOT_GUILD_MASTER };
            }
            if (!this.canChangeAllianceMemberRank(requesterMember.allianceRank)) {
                return { ok: false as const, code: messages.AllianceErrorCode.ALLIANCE_ERROR_NOT_AUTHORIZED };
            }

            let targetGuildId: number | null = null;
            let targetMember: (typeof requesterMembers)[number] | undefined;
            for (const gid of alliance.guildIds) {
                const members = [
                    ...(await this.guildMemberRepo.getAll(worldId, String(gid), { txClient: dataTx })).values(),
                ];
                const found = members.find((m) => m.characterId === targetCharacterId);
                if (found) {
                    targetGuildId = gid;
                    targetMember = found;
                    break;
                }
            }
            if (targetGuildId == null || !targetMember) {
                return { ok: false as const, code: messages.AllianceErrorCode.ALLIANCE_ERROR_TARGET_NOT_IN_ALLIANCE };
            }

            const currentRank = targetMember.allianceRank;
            if (currentRank == null || currentRank <= 2) {
                return { ok: false as const, code: messages.AllianceErrorCode.ALLIANCE_ERROR_INVALID_MEMBER_RANK };
            }

            const newAllianceRank = promote ? currentRank - 1 : currentRank + 1;
            if (newAllianceRank < 2 || newAllianceRank > 5) {
                return { ok: false as const, code: messages.AllianceErrorCode.ALLIANCE_ERROR_INVALID_MEMBER_RANK };
            }
            if (newAllianceRank === 2 && requesterMember.allianceRank !== 1) {
                return { ok: false as const, code: messages.AllianceErrorCode.ALLIANCE_ERROR_NOT_AUTHORIZED };
            }

            await this.guildMemberRepo.set(
                worldId,
                { ...targetMember, allianceRank: newAllianceRank },
                { txClient: dataTx }
            );

            const targetGuild = await this.guildRepo.get(worldId, targetGuildId, { txClient: dataTx });
            if (targetGuild) {
                await this.guildRepo.set(
                    worldId,
                    { ...targetGuild, revision: targetGuild.revision + 1 },
                    { txClient: dataTx }
                );
            }

            const lockedAlliance = await this.allianceRepo.get(worldId, allianceId, { txClient: dataTx });
            if (!lockedAlliance) {
                return { ok: false as const, code: messages.AllianceErrorCode.ALLIANCE_ERROR_ALLIANCE_NOT_FOUND };
            }

            const nextRevision = lockedAlliance.revision + 1;
            const savedAlliance = await this.allianceRepo.set(
                worldId,
                { ...lockedAlliance, revision: nextRevision },
                { txClient: dataTx }
            );

            return {
                ok: true as const,
                allianceId,
                revision: nextRevision,
                targetCharacterId,
                newAllianceRank,
                savedAlliance,
                targetGuildId,
            };
        });

        if (!txResult.ok || txResult.allianceId == null || txResult.revision == null) {
            return txResult;
        }

        await this.allianceRepo.invalidateCache(worldId, txResult.allianceId);
        const guildIdsToEvict = new Set<number>([requesterGuildId]);
        if (txResult.targetGuildId != null) {
            guildIdsToEvict.add(txResult.targetGuildId);
        }
        for (const gid of guildIdsToEvict) {
            await this.guildRepo.invalidateCache(worldId, gid);
            await this.guildMemberRepo.invalidateCache(worldId, String(gid));
        }

        const allianceMessage = await this.allianceToPb(worldId, txResult.savedAlliance);
        const wireAlliance = Alliance.encode(allianceMessage).finish();
        await this.publishToAllianceRoutes(
            ALLIANCE_EVT.MEMBER_RANK_CHANGED,
            worldId,
            txResult.allianceId,
            txResult.revision,
            {
                character_id: txResult.targetCharacterId,
                new_alliance_rank: txResult.newAllianceRank,
                alliance_pb: Buffer.from(wireAlliance).toString("base64"),
            }
        );

        return {
            ok: true,
            allianceId: txResult.allianceId,
            revision: txResult.revision,
            targetCharacterId: txResult.targetCharacterId,
            newAllianceRank: txResult.newAllianceRank,
            alliance: allianceMessage,
        };
    }

    async leaveAlliance(worldId: number, characterId: number): Promise<LeaveAllianceResult> {
        this.assertWorld(worldId);
        this.assertCharacterId(characterId);

        const locked = await this.lockRequesterAlliance(worldId, characterId);
        if (!locked.ok) {
            return locked;
        }
        await using _locks = locked.locks;
        const { guildId, guild, allianceId, alliance, requester: requesterMember } = locked;
        if (!requesterMember || requesterMember.guildRank !== GuildMemberRank.GUILD_MEMBER_RANK_MASTER) {
            return { ok: false, code: messages.AllianceErrorCode.ALLIANCE_ERROR_NOT_GUILD_MASTER };
        }
        if (requesterMember.allianceRank == null || requesterMember.allianceRank > 2) {
            return { ok: false, code: messages.AllianceErrorCode.ALLIANCE_ERROR_UNKNOWN };
        }

        const leaderGuildId = alliance.guildIds[0];
        if (leaderGuildId === guildId) {
            const guildIds = [...alliance.guildIds];
            const guildLocks = guildIds
                .filter((id) => id !== guildId)
                .sort((a, b) => a - b)
                .map((id) => `guild:${id}`);
            await using _otherGuildLocks = await this.distributedLockService.acquireWorldDataLocks(
                worldId,
                guildLocks,
            );

            const result = await this.ctx.withPgDataTransaction(worldId, allianceId, async (dataTx: PoolClient) => {
                return this.dissolveAllianceInTransaction(worldId, allianceId, dataTx);
            });

            if (!result.ok || result.allianceId == null || result.revision == null) {
                return result;
            }

            await this.allianceRepo.invalidateCache(worldId, result.allianceId);
            for (const gid of result.guildIds ?? []) {
                await this.guildRepo.invalidateCache(worldId, gid);
                await this.guildMemberRepo.invalidateCache(worldId, String(gid));
            }

            await this.publishToAllianceRoutes(ALLIANCE_EVT.DISBANDED, worldId, result.allianceId, result.revision, {
                guild_ids: result.guildIds ?? [],
                member_character_ids: result.memberCharacterIds ?? [],
            });

            return {
                ok: true,
                allianceId: result.allianceId,
                revision: result.revision,
                disbanded: true,
                removedGuildId: guildId,
            };
        }

        const otherGuildIds = alliance.guildIds.filter((id) => id !== guildId);
        await using _otherGuildLocks = await this.distributedLockService.acquireWorldDataLocks(
            worldId,
            otherGuildIds.sort((a, b) => a - b).map((id) => `guild:${id}`),
        );

        const txResult = await this.ctx.withPgDataTransaction(worldId, allianceId, async (dataTx: PoolClient) => {
            const lockedAlliance = await this.allianceRepo.get(worldId, allianceId, { txClient: dataTx });
            if (!lockedAlliance) {
                return { ok: false as const, code: messages.AllianceErrorCode.ALLIANCE_ERROR_ALLIANCE_NOT_FOUND };
            }
            if (!lockedAlliance.guildIds.includes(guildId)) {
                return { ok: false as const, code: messages.AllianceErrorCode.ALLIANCE_ERROR_NOT_IN_ALLIANCE };
            }

            const lockedGuild = await this.guildRepo.get(worldId, guildId, { txClient: dataTx });
            if (!lockedGuild || lockedGuild.allianceId == null || lockedGuild.allianceId !== allianceId) {
                return { ok: false as const, code: messages.AllianceErrorCode.ALLIANCE_ERROR_NOT_IN_ALLIANCE };
            }

            const nextRevision = lockedAlliance.revision + 1;
            const nextGuildIds = lockedAlliance.guildIds.filter((id) => id !== guildId);

            await this.guildRepo.set(
                worldId,
                { ...lockedGuild, allianceId: null, revision: lockedGuild.revision + 1 },
                { txClient: dataTx }
            );
            const members = [...(await this.guildMemberRepo.getAll(worldId, String(guildId), { txClient: dataTx })).values()];
            if (members.length > 0) {
                const cleared = members.map((member) => ({
                    ...member,
                    allianceRank: null,
                }));
                await this.guildMemberRepo.setAll(worldId, cleared, { txClient: dataTx });
            }

            const savedAlliance = await this.allianceRepo.set(
                worldId,
                {
                    ...lockedAlliance,
                    guildIds: nextGuildIds,
                    revision: nextRevision,
                },
                { txClient: dataTx }
            );

            return {
                ok: true as const,
                allianceId,
                revision: nextRevision,
                removedGuildId: guildId,
                savedAlliance,
            };
        });

        if (!txResult.ok || txResult.allianceId == null || txResult.revision == null || txResult.removedGuildId == null) {
            return txResult;
        }

        await this.allianceRepo.invalidateCache(worldId, txResult.allianceId);
        await this.guildRepo.invalidateCache(worldId, guildId);
        await this.guildMemberRepo.invalidateCache(worldId, String(guildId));
        for (const gid of otherGuildIds) {
            await this.guildRepo.invalidateCache(worldId, gid);
            await this.guildMemberRepo.invalidateCache(worldId, String(gid));
        }

        const allianceMessage = await this.allianceToPb(worldId, txResult.savedAlliance);
        const removedGuildLoaded = await this.guildService.getGuild(worldId, guildId);
        const removedGuildMessage =
            removedGuildLoaded.guild != null
                ? await this.guildService.guildToPb(worldId, removedGuildLoaded.guild, removedGuildLoaded.members ?? [])
                : undefined;

        const wireAlliance = Alliance.encode(allianceMessage).finish();
        const extraPayload: Record<string, unknown> = {
            removed_guild_id: txResult.removedGuildId,
            expelled: false,
            alliance_pb: Buffer.from(wireAlliance).toString("base64"),
        };
        if (removedGuildMessage != null) {
            const wireRemoved = Guild.encode(removedGuildMessage).finish();
            extraPayload.removed_guild_pb = Buffer.from(wireRemoved).toString("base64");
        }

        await this.publishToAllianceRoutes(
            ALLIANCE_EVT.GUILD_LEFT,
            worldId,
            txResult.allianceId,
            txResult.revision,
            extraPayload
        );

        return {
            ok: true,
            allianceId: txResult.allianceId,
            revision: txResult.revision,
            disbanded: false,
            removedGuildId: txResult.removedGuildId,
            alliance: allianceMessage,
        };
    }

    async expelAllianceGuild(
        worldId: number,
        characterId: number,
        targetGuildId: number,
        clientAllianceId?: number
    ): Promise<ExpelAllianceGuildResult> {
        this.assertWorld(worldId);
        this.assertCharacterId(characterId);
        if (!Number.isInteger(targetGuildId) || targetGuildId < 1) {
            return { ok: false, code: messages.AllianceErrorCode.ALLIANCE_ERROR_GUILD_NOT_FOUND };
        }

        const locked = await this.lockRequesterAlliance(worldId, characterId);
        if (!locked.ok) {
            return locked;
        }
        await using _locks = locked.locks;
        const { guildId: requesterGuildId, allianceId, alliance, requester: requesterMember } = locked;
        if (targetGuildId === requesterGuildId) {
            return { ok: false, code: messages.AllianceErrorCode.ALLIANCE_ERROR_CANNOT_EXPEL_OWN_GUILD };
        }
        if (clientAllianceId != null && clientAllianceId > 0 && clientAllianceId !== allianceId) {
            return { ok: false, code: messages.AllianceErrorCode.ALLIANCE_ERROR_ALLIANCE_NOT_FOUND };
        }
        if (!alliance.guildIds.includes(targetGuildId)) {
            return { ok: false, code: messages.AllianceErrorCode.ALLIANCE_ERROR_TARGET_NOT_IN_ALLIANCE };
        }
        if (!requesterMember || requesterMember.guildRank !== GuildMemberRank.GUILD_MEMBER_RANK_MASTER) {
            return { ok: false, code: messages.AllianceErrorCode.ALLIANCE_ERROR_NOT_GUILD_MASTER };
        }
        if (requesterMember.allianceRank !== 1) {
            return { ok: false, code: messages.AllianceErrorCode.ALLIANCE_ERROR_NOT_ALLIANCE_LEADER };
        }

        const leaderGuildId = alliance.guildIds[0];
        if (leaderGuildId === targetGuildId) {
            const guildIds = [...alliance.guildIds];
            const guildLocks = guildIds
                .filter((id) => id !== requesterGuildId)
                .sort((a, b) => a - b)
                .map((id) => `guild:${id}`);
            await using _otherGuildLocks = await this.distributedLockService.acquireWorldDataLocks(
                worldId,
                guildLocks,
            );

            const result = await this.ctx.withPgDataTransaction(worldId, allianceId, async (dataTx: PoolClient) => {
                return this.dissolveAllianceInTransaction(worldId, allianceId, dataTx);
            });

            if (!result.ok || result.allianceId == null || result.revision == null) {
                return result;
            }

            await this.allianceRepo.invalidateCache(worldId, result.allianceId);
            for (const gid of result.guildIds ?? []) {
                await this.guildRepo.invalidateCache(worldId, gid);
                await this.guildMemberRepo.invalidateCache(worldId, String(gid));
            }

            await this.publishToAllianceRoutes(ALLIANCE_EVT.DISBANDED, worldId, result.allianceId, result.revision, {
                guild_ids: result.guildIds ?? [],
                member_character_ids: result.memberCharacterIds ?? [],
            });

            return {
                ok: true,
                allianceId: result.allianceId,
                revision: result.revision,
                disbanded: true,
                removedGuildId: targetGuildId,
            };
        }

        const otherGuildIds = alliance.guildIds.filter((id) => id !== requesterGuildId);
        await using _otherGuildLocks = await this.distributedLockService.acquireWorldDataLocks(
            worldId,
            otherGuildIds.sort((a, b) => a - b).map((id) => `guild:${id}`),
        );

        const txResult = await this.ctx.withPgDataTransaction(worldId, allianceId, async (dataTx: PoolClient) => {
            const lockedAlliance = await this.allianceRepo.get(worldId, allianceId, { txClient: dataTx });
            if (!lockedAlliance) {
                return { ok: false as const, code: messages.AllianceErrorCode.ALLIANCE_ERROR_ALLIANCE_NOT_FOUND };
            }
            if (!lockedAlliance.guildIds.includes(targetGuildId)) {
                return { ok: false as const, code: messages.AllianceErrorCode.ALLIANCE_ERROR_TARGET_NOT_IN_ALLIANCE };
            }

            const lockedGuild = await this.guildRepo.get(worldId, targetGuildId, { txClient: dataTx });
            if (!lockedGuild || lockedGuild.allianceId == null || lockedGuild.allianceId !== allianceId) {
                return { ok: false as const, code: messages.AllianceErrorCode.ALLIANCE_ERROR_TARGET_NOT_IN_ALLIANCE };
            }

            const nextRevision = lockedAlliance.revision + 1;
            const nextGuildIds = lockedAlliance.guildIds.filter((id) => id !== targetGuildId);

            await this.guildRepo.set(
                worldId,
                { ...lockedGuild, allianceId: null, revision: lockedGuild.revision + 1 },
                { txClient: dataTx }
            );
            const members = [...(await this.guildMemberRepo.getAll(worldId, String(targetGuildId), { txClient: dataTx })).values()];
            if (members.length > 0) {
                const cleared = members.map((member) => ({
                    ...member,
                    allianceRank: null,
                }));
                await this.guildMemberRepo.setAll(worldId, cleared, { txClient: dataTx });
            }

            const savedAlliance = await this.allianceRepo.set(
                worldId,
                {
                    ...lockedAlliance,
                    guildIds: nextGuildIds,
                    revision: nextRevision,
                },
                { txClient: dataTx }
            );

            return {
                ok: true as const,
                allianceId,
                revision: nextRevision,
                removedGuildId: targetGuildId,
                savedAlliance,
            };
        });

        if (!txResult.ok || txResult.allianceId == null || txResult.revision == null || txResult.removedGuildId == null) {
            return txResult;
        }

        await this.allianceRepo.invalidateCache(worldId, txResult.allianceId);
        await this.guildRepo.invalidateCache(worldId, targetGuildId);
        await this.guildMemberRepo.invalidateCache(worldId, String(targetGuildId));
        for (const gid of otherGuildIds) {
            if (gid === targetGuildId) {
                continue;
            }
            await this.guildRepo.invalidateCache(worldId, gid);
            await this.guildMemberRepo.invalidateCache(worldId, String(gid));
        }

        const allianceMessage = await this.allianceToPb(worldId, txResult.savedAlliance);
        const removedGuildLoaded = await this.guildService.getGuild(worldId, targetGuildId);
        const removedGuildMessage =
            removedGuildLoaded.guild != null
                ? await this.guildService.guildToPb(worldId, removedGuildLoaded.guild, removedGuildLoaded.members ?? [])
                : undefined;

        const wireAlliance = Alliance.encode(allianceMessage).finish();
        const extraPayload: Record<string, unknown> = {
            removed_guild_id: txResult.removedGuildId,
            expelled: true,
            alliance_pb: Buffer.from(wireAlliance).toString("base64"),
        };
        if (removedGuildMessage != null) {
            const wireRemoved = Guild.encode(removedGuildMessage).finish();
            extraPayload.removed_guild_pb = Buffer.from(wireRemoved).toString("base64");
        }

        await this.publishToAllianceRoutes(
            ALLIANCE_EVT.GUILD_LEFT,
            worldId,
            txResult.allianceId,
            txResult.revision,
            extraPayload
        );

        return {
            ok: true,
            allianceId: txResult.allianceId,
            revision: txResult.revision,
            disbanded: false,
            removedGuildId: txResult.removedGuildId,
            alliance: allianceMessage,
        };
    }

    async disbandAlliance(worldId: number, characterId: number): Promise<DisbandAllianceResult> {
        this.assertWorld(worldId);
        this.assertCharacterId(characterId);

        const locked = await this.lockRequesterAlliance(worldId, characterId);
        if (!locked.ok) {
            return locked;
        }
        await using _locks = locked.locks;
        const { guildId, guild, allianceId, alliance } = locked;
        if (alliance.leaderCharacterId !== characterId) {
            return { ok: false, code: messages.AllianceErrorCode.ALLIANCE_ERROR_NOT_ALLIANCE_LEADER };
        }
        if (guild.leaderCharacterId !== characterId) {
            return { ok: false, code: messages.AllianceErrorCode.ALLIANCE_ERROR_NOT_GUILD_MASTER };
        }

        const guildIds = [...alliance.guildIds];
        const guildLocks = guildIds
            .filter((id) => id !== guildId)
            .sort((a, b) => a - b)
            .map((id) => `guild:${id}`);
        await using _otherGuildLocks = await this.distributedLockService.acquireWorldDataLocks(
            worldId,
            guildLocks,
        );

        const requesterMembers = [...(await this.guildMemberRepo.getAll(worldId, String(guildId))).values()];
        const requesterMember = requesterMembers.find((m) => m.characterId === characterId);
        if (!requesterMember || requesterMember.guildRank !== GuildMemberRank.GUILD_MEMBER_RANK_MASTER) {
            return { ok: false, code: messages.AllianceErrorCode.ALLIANCE_ERROR_NOT_GUILD_MASTER };
        }
        if (requesterMember.allianceRank !== 1) {
            return { ok: false, code: messages.AllianceErrorCode.ALLIANCE_ERROR_NOT_ALLIANCE_LEADER };
        }

        const result = await this.ctx.withPgDataTransaction(worldId, allianceId, async (dataTx: PoolClient) => {
            const lockedAlliance = await this.allianceRepo.get(worldId, allianceId, { txClient: dataTx });
            if (!lockedAlliance) {
                return { ok: false as const, code: messages.AllianceErrorCode.ALLIANCE_ERROR_ALLIANCE_NOT_FOUND };
            }
            if (lockedAlliance.leaderCharacterId !== characterId) {
                return { ok: false as const, code: messages.AllianceErrorCode.ALLIANCE_ERROR_NOT_ALLIANCE_LEADER };
            }

            return this.dissolveAllianceInTransaction(worldId, allianceId, dataTx);
        });

        if (!result.ok || result.allianceId == null || result.revision == null) {
            return result;
        }

        await this.allianceRepo.invalidateCache(worldId, result.allianceId);
        for (const gid of result.guildIds ?? []) {
            await this.guildRepo.invalidateCache(worldId, gid);
            await this.guildMemberRepo.invalidateCache(worldId, String(gid));
        }

        await this.publishToAllianceRoutes(ALLIANCE_EVT.DISBANDED, worldId, result.allianceId, result.revision, {
            guild_ids: result.guildIds ?? [],
            member_character_ids: result.memberCharacterIds ?? [],
        });

        return {
            ok: true,
            allianceId: result.allianceId,
            revision: result.revision,
        };
    }

    async broadcastAllianceMultiChat(
        worldId: number,
        allianceId: number,
        senderCharacterId: number,
        senderName: string,
        message: string
    ): Promise<{ ok: boolean; deliveredCount?: number }> {
        this.assertWorld(worldId);
        this.assertCharacterId(senderCharacterId);
        if (!Number.isInteger(allianceId) || allianceId < 1) {
            return { ok: false };
        }
        const trimmedMsg = (message ?? "").trim();
        if (trimmedMsg.length <= 0 || trimmedMsg.length > 500) {
            return { ok: false };
        }
        const name = (senderName ?? "").trim();
        if (!name) {
            return { ok: false };
        }
        const realtime = await this.characterRealtimeStateRepo.get(worldId, senderCharacterId);
        const senderGuildId = realtime?.guildId;
        if (senderGuildId == null || senderGuildId < 1) {
            return { ok: false };
        }
        const alliance = await this.allianceRepo.get(worldId, allianceId);
        if (!alliance) {
            return { ok: false };
        }
        if (!alliance.guildIds.includes(senderGuildId)) {
            return { ok: false };
        }
        const guildMembers = await this.guildMemberRepo.getAll(worldId, String(senderGuildId));
        if (!guildMembers.has(String(senderCharacterId))) {
            return { ok: false };
        }
        await this.publishToAllianceRoutes("multi_chat", worldId, allianceId, alliance.revision, {
            alliance_id: allianceId,
            sender_character_id: senderCharacterId,
            chat_mode: 3,
            sender_name: name,
            message: trimmedMsg,
        });
        return { ok: true, deliveredCount: 1 };
    }
}
