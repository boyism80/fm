import {
    Alliance,
    Guild,
    AllianceErrorCode,
    GuildErrorCode,
    GuildMemberRank,
    type Alliance as AllianceMessage,
    type GuildBulletinBoardReplyEntry,
    type GuildBulletinBoardThreadDetail,
    type GuildBulletinBoardThreadEntry,
    type GuildMember,
    type GuildLogo as GuildLogoMessage,
    type Guild as GuildMessage,
} from "../protobuf/generated/fminternal/internal_service";
import type { PoolClient } from "pg";
import { AppConfiguration } from "../config/app-configuration";
import { InternalContext } from "../context/internal-context";
import { UnifiedRepository } from "../repos/unified-repository";
import { AllianceRepository } from "../repos/alliance-repository";
import { GuildRepository } from "../repos/guild-repository";
import { GuildMemberRepository } from "../repos/guild-member-repository";
import { GuildBulletinBoardRepository } from "../repos/guild-bulletin-board-repository";
import { CharacterRealtimeStateRepository } from "../repos/character-realtime-state-repository";
import { CharacterRepository } from "../repos/character-repository";
import { SessionRepository } from "../repos/session-repository";
import type { CharacterSession } from "../repos/session-repository";
import type { AllianceModel } from "../repos/alliance-repository";
import type { GuildModel } from "../repos/guild-repository";
import type { GuildMemberModel } from "../repos/guild-member-repository";
import type { GuildBulletinBoardThreadModel } from "../types/repository-models";
import { RabbitMQService } from "./rabbitmq-service";
import { DistributedLockService } from "./distributed-lock-service";
import { redisCacheKey } from "../redis-cache-key";
import {
    ALLIANCE_CAPACITY_MAX,
    ALLIANCE_INCREASE_CAPACITY_MESO_COST,
    DEFAULT_ALLIANCE_CAPACITY,
    DEFAULT_ALLIANCE_RANK_TITLES,
    type AllianceRankTitles,
} from "../types/alliance-json";
import {
    DEFAULT_GUILD_LOGO,
    DEFAULT_GUILD_RANK_TITLES,
    type GuildLogo,
    type GuildRankTitles,
} from "../types/guild-json";

const messages = { GuildErrorCode, AllianceErrorCode };

const DEFAULT_GUILD_CAPACITY = 10;
const MIN_GUILD_NAME_LEN = 3;
const MAX_GUILD_NAME_LEN = 12;
const MAX_GUILD_NOTICE_LEN = 100;
const MAX_ALLIANCE_NOTICE_LEN = 100;
const GUILD_CAPACITY_STEP = 5;
const GUILD_CAPACITY_STANDARD_MAX = 100;
const GUILD_CAPACITY_EXTENDED_MAX = 200;
const GUILD_CAPACITY_EXTENDED_GP_COST = 2000;
const MAX_GUILD_BULLETIN_TITLE_LEN = 25;
const MAX_GUILD_BULLETIN_BODY_LEN = 600;
const MAX_GUILD_BULLETIN_REPLY_LEN = 25;
const GUILD_BULLETIN_THREADS_PER_PAGE = 10;
const GUILD_BULLETIN_ICON_CASH_MIN = 0x64;
const GUILD_BULLETIN_ICON_CASH_MAX = 0x6a;
const GUILD_BULLETIN_COOLDOWN_SEC = 60;
const MIN_ALLIANCE_NAME_LEN = 3;
const MAX_ALLIANCE_NAME_LEN = 12;

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

const EVT = {
    CREATED: "created",
    MEMBER_JOINED: "member_joined",
    MEMBER_LEFT: "member_left",
    RANK_TITLES_CHANGED: "rank_titles_changed",
    MEMBER_RANK_CHANGED: "member_rank_changed",
    EMBLEM_CHANGED: "emblem_changed",
    NOTICE_CHANGED: "notice_changed",
    CAPACITY_CHANGED: "capacity_changed",
    MEMBER_ONLINE_CHANGED: "member_online_changed",
    DISBANDED: "disbanded",
} as const;

export type CreateGuildResult = {
    ok: boolean;
    code?: GuildErrorCode;
    guildId?: number;
    revision?: number;
    guild?: GuildMessage;
};

export type GetGuildResult = { found: boolean; guild?: GuildModel; members?: GuildMemberModel[] };

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

export type AcceptGuildInviteResult = {
    ok: boolean;
    code?: GuildErrorCode;
    guild?: GuildMessage;
};

export type LeaveGuildResult = {
    ok: boolean;
    code?: GuildErrorCode;
    guildId?: number;
    revision?: number;
};

export type ExpelGuildResult = {
    ok: boolean;
    code?: GuildErrorCode;
    guildId?: number;
    revision?: number;
};

export type ChangeGuildRankTitlesResult = {
    ok: boolean;
    code?: GuildErrorCode;
    guildId?: number;
    revision?: number;
};

export type ChangeGuildMemberRankResult = {
    ok: boolean;
    code?: GuildErrorCode;
    guildId?: number;
    revision?: number;
};

export type ChangeGuildEmblemResult = {
    ok: boolean;
    code?: GuildErrorCode;
    guildId?: number;
    revision?: number;
};

export type ChangeGuildNoticeResult = {
    ok: boolean;
    code?: GuildErrorCode;
    guildId?: number;
    revision?: number;
};

export type DisbandGuildResult = {
    ok: boolean;
    code?: GuildErrorCode;
    guildId?: number;
    revision?: number;
};

export type IncreaseGuildCapacityResult = {
    ok: boolean;
    code?: GuildErrorCode;
    guildId?: number;
    revision?: number;
    capacity?: number;
    gp?: number;
};

export type ListGuildBulletinBoardThreadsResult = {
    ok: boolean;
    code?: GuildErrorCode;
    threads?: GuildBulletinBoardThreadEntry[];
    listStart?: number;
    threadCount?: number;
    notice?: GuildBulletinBoardThreadEntry;
};

export type ShowGuildBulletinBoardThreadResult = {
    ok: boolean;
    code?: GuildErrorCode;
    thread?: GuildBulletinBoardThreadDetail;
};

export type CreateGuildBulletinBoardThreadResult = {
    ok: boolean;
    code?: GuildErrorCode;
    thread?: GuildBulletinBoardThreadDetail;
    threads?: GuildBulletinBoardThreadEntry[];
    listStart?: number;
    threadCount?: number;
    notice?: GuildBulletinBoardThreadEntry;
};

export type UpdateGuildBulletinBoardThreadResult = {
    ok: boolean;
    code?: GuildErrorCode;
    thread?: GuildBulletinBoardThreadDetail;
};

export type DeleteGuildBulletinBoardThreadResult = {
    ok: boolean;
    code?: GuildErrorCode;
};

export type CreateGuildBulletinBoardReplyResult = {
    ok: boolean;
    code?: GuildErrorCode;
    thread?: GuildBulletinBoardThreadDetail;
};

export type DeleteGuildBulletinBoardReplyResult = {
    ok: boolean;
    code?: GuildErrorCode;
    thread?: GuildBulletinBoardThreadDetail;
};

export type GuildMembershipResult =
    | { code: GuildErrorCode.GUILD_ERROR_NONE; guildId: number; member: GuildMemberModel }
    | { code: Exclude<GuildErrorCode, GuildErrorCode.GUILD_ERROR_NONE> };

export class GuildService {
    private readonly ctx: InternalContext;
    private readonly app: AppConfiguration;
    private readonly unifiedRepo: UnifiedRepository;
    private readonly allianceRepo: AllianceRepository;
    private readonly guildRepo: GuildRepository;
    private readonly guildMemberRepo: GuildMemberRepository;
    private readonly guildBulletinBoardRepo: GuildBulletinBoardRepository;
    private readonly characterRealtimeStateRepo: CharacterRealtimeStateRepository;
    private readonly characterRepo: CharacterRepository;
    private readonly sessionRepo: SessionRepository;
    private readonly rabbitmqService: RabbitMQService;
    private readonly distributedLockService: DistributedLockService;

    constructor(
        internalContext: InternalContext,
        appConfiguration: AppConfiguration,
        unifiedRepository: UnifiedRepository,
        allianceRepository: AllianceRepository,
        guildRepository: GuildRepository,
        guildMemberRepository: GuildMemberRepository,
        guildBulletinBoardRepository: GuildBulletinBoardRepository,
        characterRealtimeStateRepository: CharacterRealtimeStateRepository,
        characterRepository: CharacterRepository,
        sessionRepository: SessionRepository,
        rabbitmqService: RabbitMQService,
        distributedLockService: DistributedLockService
    ) {
        this.ctx = internalContext;
        this.app = appConfiguration;
        this.unifiedRepo = unifiedRepository;
        this.allianceRepo = allianceRepository;
        this.guildRepo = guildRepository;
        this.guildMemberRepo = guildMemberRepository;
        this.guildBulletinBoardRepo = guildBulletinBoardRepository;
        this.characterRealtimeStateRepo = characterRealtimeStateRepository;
        this.characterRepo = characterRepository;
        this.sessionRepo = sessionRepository;
        this.rabbitmqService = rabbitmqService;
        this.distributedLockService = distributedLockService;
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

    private assertUInt16(value: number, fieldName: string) {
        const n = value;
        if (!Number.isInteger(n) || n < 0 || n > 0xffff) {
            const err = new Error(`${fieldName} must be uint16`) as Error & { code?: string };
            err.code = "INVALID_PAYLOAD";
            throw err;
        }
    }

    private normalizeGuildName(name: string) {
        return name.trim();
    }

    private validateGuildName(name: string): GuildErrorCode | null {
        if (typeof name !== "string") {
            return messages.GuildErrorCode.GUILD_ERROR_GUILD_NAME_INVALID;
        }
        const trimmed = this.normalizeGuildName(name);
        if (trimmed.length < MIN_GUILD_NAME_LEN || trimmed.length > MAX_GUILD_NAME_LEN) {
            return messages.GuildErrorCode.GUILD_ERROR_GUILD_NAME_INVALID;
        }
        return null;
    }

    private async publishToGuildRoutes(
        eventType: string,
        worldId: number,
        guildId: number,
        revision: number,
        extraPayload: Record<string, unknown> = {}
    ) {
        await this.rabbitmqService.assertDirectExchange(AMQ_DIRECT_EXCHANGE);
        const routingKey = `fm.${worldId}.all.guild`;
        return this.rabbitmqService.publish(AMQ_DIRECT_EXCHANGE, routingKey, eventType, {
            event_id: extraPayload.event_id ?? `${Date.now()}-${Math.random().toString(36).slice(2, 8)}`,
            world_id: worldId,
            guild_id: guildId,
            revision,
            occurred_at: new Date().toISOString(),
            ...extraPayload,
        });
    }

    private async rollbackCreateGuild(worldId: number, guildId: number, leaderCharacterId: number) {
        await this.ctx.withPgDataTransaction(worldId, guildId, async (dataTx: PoolClient) => {
            await this.guildMemberRepo.deleteAllForGuild(worldId, guildId, { txClient: dataTx });
            await this.guildBulletinBoardRepo.deleteAllForGuild(worldId, guildId, { txClient: dataTx });
            await this.guildRepo.delete({ worldId, guildId }, { txClient: dataTx });
        }).catch(() => {});
        await this.unifiedRepo.deleteGuildName(guildId).catch(() => {});
        await this.guildRepo.evictCache(worldId, guildId).catch(() => {});
        await this.guildMemberRepo.evictGroupCache(worldId, String(guildId)).catch(() => {});
        await this.characterRealtimeStateRepo.evictCache(worldId, leaderCharacterId).catch(() => {});
    }

    async createGuild(
        worldId: number,
        guildName: string,
        leader: GuildMember | null | undefined
    ): Promise<CreateGuildResult> {
        this.assertWorld(worldId);
        if (!leader || leader.worldId !== worldId) {
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_LEADER_MISMATCH };
        }

        const leaderCharacterId = leader.characterId;
        this.assertCharacterId(leaderCharacterId);

        const leaderName = leader.characterName.trim();
        if (!leaderName) {
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_CHARACTER_NOT_FOUND };
        }

        const leaderLevel = leader.level;
        const leaderClassId = leader.classId;
        this.assertUInt16(leaderLevel, "level");
        this.assertUInt16(leaderClassId, "class_id");

        const nameError = this.validateGuildName(guildName);
        if (nameError != null) {
            return { ok: false, code: nameError };
        }

        const normalizedName = this.normalizeGuildName(guildName);

        await using _characterRealtimeLock = await this.distributedLockService.acquireWorldDataLock(
            worldId,
            `character_realtime:${leaderCharacterId}`,
        );

        const leaderState = await this.characterRealtimeStateRepo.get(worldId, leaderCharacterId);
        if (leaderState?.guildId != null) {
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_ALREADY_IN_GUILD };
        }

        let guildId: number;
        try {
            const nameEntry = await this.unifiedRepo.reserveGuildName(normalizedName, worldId);
            guildId = Number(nameEntry.guild_id);
        } catch (err: unknown) {
            const code = (err as { code?: string })?.code;
            if (code === "23505") {
                return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_GUILD_NAME_TAKEN };
            }
            throw err;
        }

        let revision = 1;
        let savedGuild: GuildModel;
        try {
            savedGuild = await this.ctx.withPgDataTransaction(worldId, guildId, async (dataTx: PoolClient) => {
                const savedGuild = await this.guildRepo.set(
                    worldId,
                    {
                        worldId,
                        guildId,
                        name: normalizedName,
                        leaderCharacterId,
                        gp: 0,
                        capacity: DEFAULT_GUILD_CAPACITY,
                        notice: "",
                        logo: DEFAULT_GUILD_LOGO,
                        rankTitles: DEFAULT_GUILD_RANK_TITLES,
                        revision: 1,
                    },
                    { txClient: dataTx }
                );
                await this.guildMemberRepo.set(
                    worldId,
                    {
                        worldId,
                        guildId,
                        characterId: leaderCharacterId,
                        characterName: leaderName,
                        level: leaderLevel,
                        classId: leaderClassId,
                        guildRank: GuildMemberRank.GUILD_MEMBER_RANK_MASTER,
                    },
                    { txClient: dataTx }
                );
                return savedGuild;
            });

            revision = savedGuild.revision;

            await this.ctx.withPgDataTransaction(worldId, leaderCharacterId, async (dataTx: PoolClient) => {
                const state = await this.characterRealtimeStateRepo.get(worldId, leaderCharacterId, { txClient: dataTx });
                if (state?.guildId != null) {
                    const err = new Error("leader already in guild") as Error & { guildRollback?: boolean };
                    err.guildRollback = true;
                    throw err;
                }
                await this.characterRealtimeStateRepo.set(
                    worldId,
                    {
                        worldId,
                        characterId: leaderCharacterId,
                        partyId: state?.partyId ?? null,
                        guildId,
                        buddyCapacity: state?.buddyCapacity,
                    },
                    { txClient: dataTx }
                );
            });
        } catch (err: unknown) {
            await this.rollbackCreateGuild(worldId, guildId, leaderCharacterId);
            if ((err as { guildRollback?: boolean }).guildRollback) {
                return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_ALREADY_IN_GUILD };
            }
            throw err;
        }

        await this.guildRepo.evictCache(worldId, guildId);
        await this.guildMemberRepo.evictGroupCache(worldId, String(guildId));
        await this.characterRealtimeStateRepo.evictCache(worldId, leaderCharacterId);

        await this.publishToGuildRoutes(EVT.CREATED, worldId, guildId, revision, {
            leader_character_id: leaderCharacterId,
            guild_name: normalizedName,
        });

        const leaderMember: GuildMemberModel = {
            worldId,
            guildId,
            characterId: leaderCharacterId,
            characterName: leaderName,
            level: leaderLevel,
            classId: leaderClassId,
            guildRank: GuildMemberRank.GUILD_MEMBER_RANK_MASTER,
        };
        const guildMessage = await this.buildGuildMessage(worldId, savedGuild, [leaderMember]);

        return { ok: true, guildId, revision, guild: guildMessage };
    }

    async acceptGuildInvite(
        worldId: number,
        guildId: number,
        member: GuildMember | null | undefined
    ): Promise<AcceptGuildInviteResult> {
        this.assertWorld(worldId);
        this.assertGuildId(guildId);
        if (!member || member.worldId !== worldId) {
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_UNKNOWN };
        }

        const characterId = member.characterId;
        this.assertCharacterId(characterId);

        const name = member.characterName.trim();
        if (!name) {
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_CHARACTER_NOT_FOUND };
        }

        const memberLevel = member.level;
        const memberClassId = member.classId;
        this.assertUInt16(memberLevel, "level");
        this.assertUInt16(memberClassId, "class_id");

        await using _characterRealtimeLock = await this.distributedLockService.acquireWorldDataLock(
            worldId,
            `character_realtime:${characterId}`,
        );

        await using _guildLock = await this.distributedLockService.acquireWorldDataLock(
            worldId,
            `guild:${guildId}`,
        );

        const memberState = await this.characterRealtimeStateRepo.get(worldId, characterId);
        if (memberState?.guildId != null) {
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_ALREADY_IN_GUILD };
        }

        const result = await this.ctx.withPgDataTransaction(worldId, guildId, async (txClient: PoolClient) => {
            const guild = await this.guildRepo.get(worldId, guildId, { txClient });
            if (!guild) {
                return { ok: false as const, code: messages.GuildErrorCode.GUILD_ERROR_GUILD_NOT_FOUND };
            }

            const members = await this.guildMemberRepo.getAll(worldId, String(guildId), { txClient });
            if (members.size >= guild.capacity) {
                return { ok: false as const, code: messages.GuildErrorCode.GUILD_ERROR_GUILD_FULL };
            }

            if (members.has(String(characterId))) {
                return { ok: false as const, code: messages.GuildErrorCode.GUILD_ERROR_ALREADY_IN_GUILD };
            }

            await this.guildMemberRepo.set(
                worldId,
                {
                    worldId,
                    guildId,
                    characterId,
                    characterName: name,
                    level: memberLevel,
                    classId: memberClassId,
                    guildRank: GuildMemberRank.GUILD_MEMBER_RANK_NEW,
                },
                { txClient }
            );

            const nextRevision = guild.revision + 1;
            const updatedGuild = await this.guildRepo.set(
                worldId,
                { ...guild, revision: nextRevision },
                { txClient }
            );

            return { ok: true as const, revision: updatedGuild.revision };
        });

        if (!result.ok || result.revision == null) {
            return { ok: false, code: result.code ?? messages.GuildErrorCode.GUILD_ERROR_UNKNOWN };
        }

        await this.ctx.withPgDataTransaction(worldId, characterId, async (txClient: PoolClient) => {
            await this.characterRealtimeStateRepo.set(
                worldId,
                {
                    worldId,
                    characterId,
                    partyId: memberState?.partyId ?? null,
                    guildId,
                    buddyCapacity: memberState?.buddyCapacity,
                },
                { txClient }
            );
        });

        await this.guildRepo.evictCache(worldId, guildId);
        await this.guildMemberRepo.evictGroupCache(worldId, String(guildId));
        await this.characterRealtimeStateRepo.evictCache(worldId, characterId);

        await this.publishToGuildRoutes(EVT.MEMBER_JOINED, worldId, guildId, result.revision, {
            character_id: characterId,
        });

        const loaded = await this.getGuild(worldId, guildId);
        if (!loaded.found || !loaded.guild || !loaded.members) {
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_UNKNOWN };
        }

        const guildMessage = await this.buildGuildMessage(worldId, loaded.guild, loaded.members);
        return { ok: true, guild: guildMessage };
    }

    async leaveGuild(worldId: number, characterId: number): Promise<LeaveGuildResult> {
        this.assertWorld(worldId);
        this.assertCharacterId(characterId);

        await using _characterRealtimeLock = await this.distributedLockService.acquireWorldDataLock(
            worldId,
            `character_realtime:${characterId}`,
        );

        const lockedState = await this.characterRealtimeStateRepo.get(worldId, characterId);
        if (lockedState?.guildId == null) {
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_NOT_IN_GUILD };
        }
        const lockedGuildId = lockedState.guildId;
        if (!Number.isInteger(lockedGuildId) || lockedGuildId < 1) {
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_NOT_IN_GUILD };
        }

        await using _guildLock = await this.distributedLockService.acquireWorldDataLock(
            worldId,
            `guild:${lockedGuildId}`,
        );

        const guildResult = await this.ctx.withPgDataTransaction(worldId, lockedGuildId, async (txClient: PoolClient) => {
            const guild = await this.guildRepo.get(worldId, lockedGuildId, { txClient });
            if (!guild) {
                return { ok: true as const, guildId: lockedGuildId, revision: 0, clearRealtime: true as const };
            }

            if (guild.leaderCharacterId === characterId) {
                return { ok: false as const, code: messages.GuildErrorCode.GUILD_ERROR_LEADER_CANNOT_LEAVE };
            }

            const members = await this.guildMemberRepo.getAll(worldId, String(lockedGuildId), { txClient });
            if (!members.has(String(characterId))) {
                return { ok: false as const, code: messages.GuildErrorCode.GUILD_ERROR_NOT_IN_GUILD };
            }

            await this.guildMemberRepo.del(worldId, String(lockedGuildId), characterId, { txClient });

            const nextRevision = guild.revision + 1;
            const updatedGuild = await this.guildRepo.set(worldId, { ...guild, revision: nextRevision }, { txClient });
            return { ok: true as const, guildId: updatedGuild.guildId, revision: updatedGuild.revision, clearRealtime: true as const };
        });

        if (!guildResult.ok || guildResult.guildId == null || guildResult.revision == null) {
            return guildResult;
        }

        if (guildResult.clearRealtime) {
            await this.ctx.withPgDataTransaction(worldId, characterId, async (txClient: PoolClient) => {
                await this.characterRealtimeStateRepo.set(
                    worldId,
                    {
                        worldId,
                        characterId,
                        partyId: lockedState.partyId ?? null,
                        guildId: null,
                        buddyCapacity: lockedState.buddyCapacity,
                    },
                    { txClient }
                );
            });
        }
        const result = guildResult;

        await this.guildRepo.evictCache(worldId, result.guildId);
        await this.guildMemberRepo.evictGroupCache(worldId, String(result.guildId));
        await this.characterRealtimeStateRepo.evictCache(worldId, characterId);

        await this.publishToGuildRoutes(EVT.MEMBER_LEFT, worldId, result.guildId, result.revision, {
            character_id: characterId,
        });

        return { ok: true, guildId: result.guildId, revision: result.revision };
    }

    private canExpelGuildMembers(rank: GuildMemberRank | undefined) {
        return rank === GuildMemberRank.GUILD_MEMBER_RANK_MASTER
            || rank === GuildMemberRank.GUILD_MEMBER_RANK_JUNIOR;
    }

    private guildMemberRankValue(rank: GuildMemberRank | undefined) {
        if (rank == null || rank === GuildMemberRank.GUILD_MEMBER_RANK_UNSPECIFIED) {
            return GuildMemberRank.GUILD_MEMBER_RANK_MEMBER;
        }
        return rank;
    }

    async expelGuild(
        worldId: number,
        requesterCharacterId: number,
        targetCharacterId: number
    ): Promise<ExpelGuildResult> {
        this.assertWorld(worldId);
        this.assertCharacterId(requesterCharacterId);
        this.assertCharacterId(targetCharacterId);
        if (requesterCharacterId === targetCharacterId) {
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_CANNOT_EXPEL_TARGET };
        }

        const firstCharacterId = requesterCharacterId < targetCharacterId
            ? requesterCharacterId
            : targetCharacterId;
        const secondCharacterId = requesterCharacterId < targetCharacterId
            ? targetCharacterId
            : requesterCharacterId;

        await using _firstCharacterRealtimeLock = await this.distributedLockService.acquireWorldDataLock(
            worldId,
            `character_realtime:${firstCharacterId}`,
        );

        await using _secondCharacterRealtimeLock = await this.distributedLockService.acquireWorldDataLock(
            worldId,
            `character_realtime:${secondCharacterId}`,
        );

        const requesterLockedState = await this.characterRealtimeStateRepo.get(worldId, requesterCharacterId);
        if (requesterLockedState?.guildId == null) {
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_NOT_IN_GUILD };
        }
        const lockedGuildId = requesterLockedState.guildId;
        if (!Number.isInteger(lockedGuildId) || lockedGuildId < 1) {
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_NOT_IN_GUILD };
        }

        await using _guildLock = await this.distributedLockService.acquireWorldDataLock(
            worldId,
            `guild:${lockedGuildId}`,
        );

        const targetLockedState = await this.characterRealtimeStateRepo.get(worldId, targetCharacterId);
        if (targetLockedState?.guildId !== lockedGuildId) {
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_TARGET_NOT_IN_GUILD };
        }

        const guildResult = await this.ctx.withPgDataTransaction(worldId, lockedGuildId, async (txClient: PoolClient) => {
            const guild = await this.guildRepo.get(worldId, lockedGuildId, { txClient });
            if (!guild) {
                return { ok: false as const, code: messages.GuildErrorCode.GUILD_ERROR_GUILD_NOT_FOUND };
            }

            const members = await this.guildMemberRepo.getAll(worldId, String(lockedGuildId), { txClient });
            const requesterMember = members.get(String(requesterCharacterId));
            if (!requesterMember) {
                return { ok: false as const, code: messages.GuildErrorCode.GUILD_ERROR_NOT_IN_GUILD };
            }
            if (!this.canExpelGuildMembers(requesterMember.guildRank)) {
                return { ok: false as const, code: messages.GuildErrorCode.GUILD_ERROR_NOT_AUTHORIZED };
            }

            const targetMember = members.get(String(targetCharacterId));
            if (!targetMember) {
                return { ok: false as const, code: messages.GuildErrorCode.GUILD_ERROR_TARGET_NOT_IN_GUILD };
            }

            const requesterRank = this.guildMemberRankValue(requesterMember.guildRank);
            const targetRank = this.guildMemberRankValue(targetMember.guildRank);
            if (requesterRank >= targetRank) {
                return { ok: false as const, code: messages.GuildErrorCode.GUILD_ERROR_CANNOT_EXPEL_TARGET };
            }

            await this.guildMemberRepo.del(worldId, String(lockedGuildId), targetCharacterId, { txClient });

            const nextRevision = guild.revision + 1;
            const updatedGuild = await this.guildRepo.set(worldId, { ...guild, revision: nextRevision }, { txClient });
            return { ok: true as const, guildId: updatedGuild.guildId, revision: updatedGuild.revision };
        });

        if (!guildResult.ok || guildResult.guildId == null || guildResult.revision == null) {
            return guildResult;
        }

        await this.ctx.withPgDataTransaction(worldId, targetCharacterId, async (txClient: PoolClient) => {
            await this.characterRealtimeStateRepo.set(
                worldId,
                {
                    worldId,
                    characterId: targetCharacterId,
                    partyId: targetLockedState.partyId ?? null,
                    guildId: null,
                    buddyCapacity: targetLockedState.buddyCapacity,
                },
                { txClient }
            );
        });
        const result = guildResult;

        await this.guildRepo.evictCache(worldId, result.guildId);
        await this.guildMemberRepo.evictGroupCache(worldId, String(result.guildId));
        await this.characterRealtimeStateRepo.evictCache(worldId, targetCharacterId);

        await this.publishToGuildRoutes(EVT.MEMBER_LEFT, worldId, result.guildId, result.revision, {
            character_id: targetCharacterId,
            expelled: true,
        });

        const targetChannelIndex = await this.sessionChannelIndex(worldId, targetCharacterId);
        if (targetChannelIndex < 0) {
            // TODO: game server -> internal sendNote RPC (promise) when offline member is expelled
        }

        return { ok: true, guildId: result.guildId, revision: result.revision };
    }

    private normalizeRankTitles(rankTitles: string[] | null | undefined): GuildRankTitles | null {
        if (!rankTitles || rankTitles.length !== 5) {
            return null;
        }
        return rankTitles.map((entry) => String(entry ?? "")) as GuildRankTitles;
    }

    async changeGuildRankTitles(
        worldId: number,
        characterId: number,
        rankTitles: string[]
    ): Promise<ChangeGuildRankTitlesResult> {
        this.assertWorld(worldId);
        this.assertCharacterId(characterId);

        const normalizedRankTitles = this.normalizeRankTitles(rankTitles);
        if (!normalizedRankTitles) {
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_INVALID_RANK_TITLES };
        }

        await using _characterRealtimeLock = await this.distributedLockService.acquireWorldDataLock(
            worldId,
            `character_realtime:${characterId}`,
        );

        const lockedState = await this.characterRealtimeStateRepo.get(worldId, characterId);
        if (lockedState?.guildId == null) {
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_NOT_IN_GUILD };
        }
        const lockedGuildId = lockedState.guildId;
        if (!Number.isInteger(lockedGuildId) || lockedGuildId < 1) {
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_NOT_IN_GUILD };
        }

        await using _guildLock = await this.distributedLockService.acquireWorldDataLock(
            worldId,
            `guild:${lockedGuildId}`,
        );

        const result = await this.ctx.withPgDataTransaction(worldId, lockedGuildId, async (txClient: PoolClient) => {
            const guild = await this.guildRepo.get(worldId, lockedGuildId, { txClient });
            if (!guild) {
                return { ok: false as const, code: messages.GuildErrorCode.GUILD_ERROR_GUILD_NOT_FOUND };
            }
            if (guild.leaderCharacterId !== characterId) {
                return { ok: false as const, code: messages.GuildErrorCode.GUILD_ERROR_NOT_AUTHORIZED };
            }

            const members = await this.guildMemberRepo.getAll(worldId, String(lockedGuildId), { txClient });
            const requesterMember = members.get(String(characterId));
            if (!requesterMember) {
                return { ok: false as const, code: messages.GuildErrorCode.GUILD_ERROR_NOT_IN_GUILD };
            }
            if (requesterMember.guildRank !== GuildMemberRank.GUILD_MEMBER_RANK_MASTER) {
                return { ok: false as const, code: messages.GuildErrorCode.GUILD_ERROR_NOT_AUTHORIZED };
            }

            const nextRevision = guild.revision + 1;
            const updatedGuild = await this.guildRepo.set(
                worldId,
                { ...guild, rankTitles: normalizedRankTitles, revision: nextRevision },
                { txClient }
            );
            return { ok: true as const, guildId: updatedGuild.guildId, revision: updatedGuild.revision };
        });

        if (!result.ok || result.guildId == null || result.revision == null) {
            return result;
        }

        await this.guildRepo.evictCache(worldId, result.guildId);

        await this.publishToGuildRoutes(EVT.RANK_TITLES_CHANGED, worldId, result.guildId, result.revision, {});

        return { ok: true, guildId: result.guildId, revision: result.revision };
    }

    private canChangeMemberRank(rank: GuildMemberRank | undefined) {
        return rank === GuildMemberRank.GUILD_MEMBER_RANK_MASTER
            || rank === GuildMemberRank.GUILD_MEMBER_RANK_JUNIOR;
    }

    private parseAssignableMemberRank(newRank: GuildMemberRank): GuildMemberRank | null {
        if (newRank === GuildMemberRank.GUILD_MEMBER_RANK_JUNIOR
            || newRank === GuildMemberRank.GUILD_MEMBER_RANK_SENIOR
            || newRank === GuildMemberRank.GUILD_MEMBER_RANK_MEMBER
            || newRank === GuildMemberRank.GUILD_MEMBER_RANK_NEW) {
            return newRank;
        }
        return null;
    }

    async changeGuildMemberRank(
        worldId: number,
        requesterCharacterId: number,
        targetCharacterId: number,
        newRank: GuildMemberRank
    ): Promise<ChangeGuildMemberRankResult> {
        this.assertWorld(worldId);
        this.assertCharacterId(requesterCharacterId);
        this.assertCharacterId(targetCharacterId);

        const parsedRank = this.parseAssignableMemberRank(newRank);
        if (!parsedRank) {
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_INVALID_MEMBER_RANK };
        }

        const firstCharacterId = requesterCharacterId < targetCharacterId
            ? requesterCharacterId
            : targetCharacterId;
        const secondCharacterId = requesterCharacterId < targetCharacterId
            ? targetCharacterId
            : requesterCharacterId;

        await using _firstCharacterRealtimeLock = await this.distributedLockService.acquireWorldDataLock(
            worldId,
            `character_realtime:${firstCharacterId}`,
        );

        await using _secondCharacterRealtimeLock = await this.distributedLockService.acquireWorldDataLock(
            worldId,
            `character_realtime:${secondCharacterId}`,
        );

        const requesterLockedState = await this.characterRealtimeStateRepo.get(worldId, requesterCharacterId);
        if (requesterLockedState?.guildId == null) {
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_NOT_IN_GUILD };
        }
        const lockedGuildId = requesterLockedState.guildId;
        if (!Number.isInteger(lockedGuildId) || lockedGuildId < 1) {
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_NOT_IN_GUILD };
        }

        await using _guildLock = await this.distributedLockService.acquireWorldDataLock(
            worldId,
            `guild:${lockedGuildId}`,
        );

        const targetLockedState = await this.characterRealtimeStateRepo.get(worldId, targetCharacterId);
        if (targetLockedState?.guildId !== lockedGuildId) {
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_TARGET_NOT_IN_GUILD };
        }

        const result = await this.ctx.withPgDataTransaction(worldId, lockedGuildId, async (txClient: PoolClient) => {
            const guild = await this.guildRepo.get(worldId, lockedGuildId, { txClient });
            if (!guild) {
                return { ok: false as const, code: messages.GuildErrorCode.GUILD_ERROR_GUILD_NOT_FOUND };
            }

            const members = await this.guildMemberRepo.getAll(worldId, String(lockedGuildId), { txClient });
            const requesterMember = members.get(String(requesterCharacterId));
            if (!requesterMember) {
                return { ok: false as const, code: messages.GuildErrorCode.GUILD_ERROR_NOT_IN_GUILD };
            }
            if (!this.canChangeMemberRank(requesterMember.guildRank)) {
                return { ok: false as const, code: messages.GuildErrorCode.GUILD_ERROR_NOT_AUTHORIZED };
            }
            if (parsedRank === GuildMemberRank.GUILD_MEMBER_RANK_JUNIOR
                && requesterMember.guildRank !== GuildMemberRank.GUILD_MEMBER_RANK_MASTER) {
                return { ok: false as const, code: messages.GuildErrorCode.GUILD_ERROR_NOT_AUTHORIZED };
            }

            const targetMember = members.get(String(targetCharacterId));
            if (!targetMember) {
                return { ok: false as const, code: messages.GuildErrorCode.GUILD_ERROR_TARGET_NOT_IN_GUILD };
            }

            await this.guildMemberRepo.set(
                worldId,
                { ...targetMember, guildRank: parsedRank },
                { txClient }
            );

            const nextRevision = guild.revision + 1;
            const updatedGuild = await this.guildRepo.set(worldId, { ...guild, revision: nextRevision }, { txClient });
            return { ok: true as const, guildId: updatedGuild.guildId, revision: updatedGuild.revision };
        });

        if (!result.ok || result.guildId == null || result.revision == null) {
            return result;
        }

        await this.guildRepo.evictCache(worldId, result.guildId);
        await this.guildMemberRepo.evictGroupCache(worldId, String(result.guildId));

        await this.publishToGuildRoutes(EVT.MEMBER_RANK_CHANGED, worldId, result.guildId, result.revision, {
            character_id: targetCharacterId,
        });

        return { ok: true, guildId: result.guildId, revision: result.revision };
    }

    private normalizeGuildLogo(logo: GuildLogoMessage | null | undefined): GuildLogo | null {
        if (!logo) {
            return null;
        }
        this.assertUInt16(logo.logo, "logo");
        this.assertUInt16(logo.logoBg, "logo_bg");
        const logoColor = logo.logoColor;
        const logoBgColor = logo.logoBgColor;
        if (!Number.isInteger(logoColor) || logoColor < 0 || logoColor > 0xff) {
            return null;
        }
        if (!Number.isInteger(logoBgColor) || logoBgColor < 0 || logoBgColor > 0xff) {
            return null;
        }
        return {
            logo: logo.logo,
            logoColor,
            logoBG: logo.logoBg,
            logoBGColor: logoBgColor,
        };
    }

    async changeGuildEmblem(
        worldId: number,
        characterId: number,
        logo: GuildLogoMessage | null | undefined
    ): Promise<ChangeGuildEmblemResult> {
        this.assertWorld(worldId);
        this.assertCharacterId(characterId);

        const normalizedLogo = this.normalizeGuildLogo(logo);
        if (!normalizedLogo) {
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_UNKNOWN };
        }

        await using _characterRealtimeLock = await this.distributedLockService.acquireWorldDataLock(
            worldId,
            `character_realtime:${characterId}`,
        );

        const lockedState = await this.characterRealtimeStateRepo.get(worldId, characterId);
        if (lockedState?.guildId == null) {
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_NOT_IN_GUILD };
        }
        const lockedGuildId = lockedState.guildId;
        if (!Number.isInteger(lockedGuildId) || lockedGuildId < 1) {
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_NOT_IN_GUILD };
        }

        await using _guildLock = await this.distributedLockService.acquireWorldDataLock(
            worldId,
            `guild:${lockedGuildId}`,
        );

        const result = await this.ctx.withPgDataTransaction(worldId, lockedGuildId, async (txClient: PoolClient) => {
            const guild = await this.guildRepo.get(worldId, lockedGuildId, { txClient });
            if (!guild) {
                return { ok: false as const, code: messages.GuildErrorCode.GUILD_ERROR_GUILD_NOT_FOUND };
            }
            if (guild.leaderCharacterId !== characterId) {
                return { ok: false as const, code: messages.GuildErrorCode.GUILD_ERROR_NOT_AUTHORIZED };
            }

            const members = await this.guildMemberRepo.getAll(worldId, String(lockedGuildId), { txClient });
            const requesterMember = members.get(String(characterId));
            if (!requesterMember) {
                return { ok: false as const, code: messages.GuildErrorCode.GUILD_ERROR_NOT_IN_GUILD };
            }
            if (requesterMember.guildRank !== GuildMemberRank.GUILD_MEMBER_RANK_MASTER) {
                return { ok: false as const, code: messages.GuildErrorCode.GUILD_ERROR_NOT_AUTHORIZED };
            }

            const nextRevision = guild.revision + 1;
            const updatedGuild = await this.guildRepo.set(
                worldId,
                { ...guild, logo: normalizedLogo, revision: nextRevision },
                { txClient }
            );
            return { ok: true as const, guildId: updatedGuild.guildId, revision: updatedGuild.revision };
        });

        if (!result.ok || result.guildId == null || result.revision == null) {
            return result;
        }

        await this.guildRepo.evictCache(worldId, result.guildId);

        await this.publishToGuildRoutes(EVT.EMBLEM_CHANGED, worldId, result.guildId, result.revision, {});

        return { ok: true, guildId: result.guildId, revision: result.revision };
    }

    private normalizeGuildNotice(notice: string | null | undefined): string | null {
        if (notice == null) {
            return null;
        }
        if (notice.length > MAX_GUILD_NOTICE_LEN) {
            return null;
        }
        return notice;
    }

    async changeGuildNotice(
        worldId: number,
        characterId: number,
        notice: string | null | undefined
    ): Promise<ChangeGuildNoticeResult> {
        this.assertWorld(worldId);
        this.assertCharacterId(characterId);

        const normalizedNotice = this.normalizeGuildNotice(notice);
        if (normalizedNotice == null) {
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_INVALID_NOTICE };
        }

        await using _characterRealtimeLock = await this.distributedLockService.acquireWorldDataLock(
            worldId,
            `character_realtime:${characterId}`,
        );

        const lockedState = await this.characterRealtimeStateRepo.get(worldId, characterId);
        if (lockedState?.guildId == null) {
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_NOT_IN_GUILD };
        }
        const lockedGuildId = lockedState.guildId;
        if (!Number.isInteger(lockedGuildId) || lockedGuildId < 1) {
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_NOT_IN_GUILD };
        }

        await using _guildLock = await this.distributedLockService.acquireWorldDataLock(
            worldId,
            `guild:${lockedGuildId}`,
        );

        const result = await this.ctx.withPgDataTransaction(worldId, lockedGuildId, async (txClient: PoolClient) => {
            const guild = await this.guildRepo.get(worldId, lockedGuildId, { txClient });
            if (!guild) {
                return { ok: false as const, code: messages.GuildErrorCode.GUILD_ERROR_GUILD_NOT_FOUND };
            }

            const members = await this.guildMemberRepo.getAll(worldId, String(lockedGuildId), { txClient });
            const requesterMember = members.get(String(characterId));
            if (!requesterMember) {
                return { ok: false as const, code: messages.GuildErrorCode.GUILD_ERROR_NOT_IN_GUILD };
            }
            if (!this.canChangeMemberRank(requesterMember.guildRank)) {
                return { ok: false as const, code: messages.GuildErrorCode.GUILD_ERROR_NOT_AUTHORIZED };
            }

            const nextRevision = guild.revision + 1;
            const updatedGuild = await this.guildRepo.set(
                worldId,
                { ...guild, notice: normalizedNotice, revision: nextRevision },
                { txClient }
            );
            return { ok: true as const, guildId: updatedGuild.guildId, revision: updatedGuild.revision };
        });

        if (!result.ok || result.guildId == null || result.revision == null) {
            return result;
        }

        await this.guildRepo.evictCache(worldId, result.guildId);

        await this.publishToGuildRoutes(EVT.NOTICE_CHANGED, worldId, result.guildId, result.revision, {});

        return { ok: true, guildId: result.guildId, revision: result.revision };
    }

    async increaseGuildCapacity(
        worldId: number,
        characterId: number,
        extendedCap: boolean
    ): Promise<IncreaseGuildCapacityResult> {
        this.assertWorld(worldId);
        this.assertCharacterId(characterId);

        await using _characterRealtimeLock = await this.distributedLockService.acquireWorldDataLock(
            worldId,
            `character_realtime:${characterId}`,
        );

        const lockedState = await this.characterRealtimeStateRepo.get(worldId, characterId);
        if (lockedState?.guildId == null) {
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_NOT_IN_GUILD };
        }
        const lockedGuildId = lockedState.guildId;
        if (!Number.isInteger(lockedGuildId) || lockedGuildId < 1) {
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_NOT_IN_GUILD };
        }

        await using _guildLock = await this.distributedLockService.acquireWorldDataLock(
            worldId,
            `guild:${lockedGuildId}`,
        );

        const result = await this.ctx.withPgDataTransaction(worldId, lockedGuildId, async (txClient: PoolClient) => {
            const guild = await this.guildRepo.get(worldId, lockedGuildId, { txClient });
            if (!guild) {
                return { ok: false as const, code: messages.GuildErrorCode.GUILD_ERROR_GUILD_NOT_FOUND };
            }

            const members = await this.guildMemberRepo.getAll(worldId, String(lockedGuildId), { txClient });
            const requesterMember = members.get(String(characterId));
            if (!requesterMember) {
                return { ok: false as const, code: messages.GuildErrorCode.GUILD_ERROR_NOT_IN_GUILD };
            }
            if (requesterMember.guildRank !== GuildMemberRank.GUILD_MEMBER_RANK_MASTER) {
                return { ok: false as const, code: messages.GuildErrorCode.GUILD_ERROR_NOT_AUTHORIZED };
            }

            const maxCapacity = extendedCap ? GUILD_CAPACITY_EXTENDED_MAX : GUILD_CAPACITY_STANDARD_MAX;
            if (guild.capacity + GUILD_CAPACITY_STEP > maxCapacity) {
                return { ok: false as const, code: messages.GuildErrorCode.GUILD_ERROR_CAPACITY_REACHED };
            }

            let nextGP = guild.gp;
            if (extendedCap) {
                if (nextGP < GUILD_CAPACITY_EXTENDED_GP_COST) {
                    return { ok: false as const, code: messages.GuildErrorCode.GUILD_ERROR_INSUFFICIENT_GUILD_GP };
                }
                nextGP -= GUILD_CAPACITY_EXTENDED_GP_COST;
            }

            const nextRevision = guild.revision + 1;
            const updatedGuild = await this.guildRepo.set(
                worldId,
                {
                    ...guild,
                    capacity: guild.capacity + GUILD_CAPACITY_STEP,
                    gp: nextGP,
                    revision: nextRevision,
                },
                { txClient }
            );
            return {
                ok: true as const,
                guildId: updatedGuild.guildId,
                revision: updatedGuild.revision,
                capacity: updatedGuild.capacity,
                gp: updatedGuild.gp,
            };
        });

        if (!result.ok || result.guildId == null || result.revision == null) {
            return result;
        }

        await this.guildRepo.evictCache(worldId, result.guildId);

        await this.publishToGuildRoutes(EVT.CAPACITY_CHANGED, worldId, result.guildId, result.revision, {
            capacity: result.capacity,
            gp: result.gp,
        });

        return {
            ok: true,
            guildId: result.guildId,
            revision: result.revision,
            capacity: result.capacity,
            gp: result.gp,
        };
    }

    async disbandGuild(worldId: number, characterId: number): Promise<DisbandGuildResult> {
        this.assertWorld(worldId);
        this.assertCharacterId(characterId);

        await using _characterRealtimeLock = await this.distributedLockService.acquireWorldDataLock(
            worldId,
            `character_realtime:${characterId}`,
        );

        const lockedState = await this.characterRealtimeStateRepo.get(worldId, characterId);
        if (lockedState?.guildId == null) {
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_NOT_IN_GUILD };
        }
        const lockedGuildId = lockedState.guildId;
        if (!Number.isInteger(lockedGuildId) || lockedGuildId < 1) {
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_NOT_IN_GUILD };
        }

        await using _guildLock = await this.distributedLockService.acquireWorldDataLock(
            worldId,
            `guild:${lockedGuildId}`,
        );

        const membersMap = await this.guildMemberRepo.getAll(worldId, String(lockedGuildId));
        const otherMemberKeyRests = [...membersMap.values()]
            .map((member) => member.characterId)
            .filter((memberCharacterId) => memberCharacterId !== characterId)
            .sort((a, b) => a - b)
            .map((memberCharacterId) => `character_realtime:${memberCharacterId}`);

        await using _otherMemberLocks = await this.distributedLockService.acquireWorldDataLocks(
            worldId,
            otherMemberKeyRests,
        );

        const guildResult = await this.ctx.withPgDataTransaction(worldId, lockedGuildId, async (txClient: PoolClient) => {
            const guild = await this.guildRepo.get(worldId, lockedGuildId, { txClient });
            if (!guild) {
                return { ok: false as const, code: messages.GuildErrorCode.GUILD_ERROR_GUILD_NOT_FOUND };
            }

            const members = await this.guildMemberRepo.getAll(worldId, String(lockedGuildId), { txClient });
            const requesterMember = members.get(String(characterId));
            if (!requesterMember || requesterMember.guildRank !== GuildMemberRank.GUILD_MEMBER_RANK_MASTER) {
                return { ok: false as const, code: messages.GuildErrorCode.GUILD_ERROR_NOT_AUTHORIZED };
            }

            const memberCharacterIds = [...members.values()].map((m) => m.characterId);
            const nextRevision = guild.revision + 1;
            await this.guildRepo.set(
                worldId,
                { ...guild, revision: nextRevision, disbandedAt: new Date() },
                { txClient }
            );
            await this.guildMemberRepo.deleteAllForGuild(worldId, lockedGuildId, { txClient });

            return {
                ok: true as const,
                guildId: lockedGuildId,
                revision: nextRevision,
                memberCharacterIds,
            };
        });

        if (!guildResult.ok || guildResult.guildId == null || guildResult.revision == null) {
            return guildResult;
        }

        for (const memberCharacterId of guildResult.memberCharacterIds ?? []) {
            await this.ctx.withPgDataTransaction(worldId, memberCharacterId, async (txClient: PoolClient) => {
                const memberState = await this.characterRealtimeStateRepo.get(worldId, memberCharacterId, { txClient });
                if (!memberState) {
                    return;
                }
                await this.characterRealtimeStateRepo.set(
                    worldId,
                    {
                        worldId,
                        characterId: memberCharacterId,
                        partyId: memberState.partyId ?? null,
                        guildId: null,
                        buddyCapacity: memberState.buddyCapacity,
                    },
                    { txClient }
                );
            });
        }
        const result = guildResult;

        await this.guildBulletinBoardRepo.deleteAllForGuild(worldId, result.guildId);

        await this.guildRepo.evictCache(worldId, result.guildId);
        await this.guildMemberRepo.evictGroupCache(worldId, String(result.guildId));
        for (const memberCharacterId of result.memberCharacterIds ?? []) {
            await this.characterRealtimeStateRepo.evictCache(worldId, memberCharacterId);
        }
        await this.unifiedRepo.deleteGuildName(result.guildId).catch(() => {});

        await this.publishToGuildRoutes(EVT.DISBANDED, worldId, result.guildId, result.revision, {
            requester_character_id: characterId,
            member_character_ids: result.memberCharacterIds ?? [],
        });

        return { ok: true, guildId: result.guildId, revision: result.revision };
    }

    async applyMemberOnlineState(worldId: number, characterId: number, online: boolean) {
        this.assertWorld(worldId);
        this.assertCharacterId(characterId);

        const state = await this.characterRealtimeStateRepo.get(worldId, characterId);
        if (state?.guildId == null) {
            return;
        }
        const guildId = state.guildId;
        if (!Number.isInteger(guildId) || guildId < 1) {
            return;
        }

        const guild = await this.guildRepo.get(worldId, guildId);
        if (!guild) {
            return;
        }
        const membersMap = await this.guildMemberRepo.getAll(worldId, String(guildId));
        if (!membersMap.has(String(characterId))) {
            return;
        }

        await this.publishToGuildRoutes(EVT.MEMBER_ONLINE_CHANGED, worldId, guildId, guild.revision, {
            character_id: characterId,
            online,
        });
    }

    private assertGuildId(guildId: number) {
        const n = guildId;
        if (!Number.isInteger(n) || n < 1 || n > 0xffffffff) {
            const err = new Error("guild_id must be a positive uint32") as Error & { code?: string };
            err.code = "INVALID_GUILD_ID";
            throw err;
        }
    }

    private computeGuildMemberChannelIndex(sess: CharacterSession | null) {
        const ch = sess?.gameServer?.channelId;
        if (ch == null || !Number.isFinite(ch) || ch < 0) {
            return -2;
        }
        return ch;
    }

    private async getCharacterSession(worldId: number, characterId: number): Promise<CharacterSession | null> {
        const row = await this.characterRepo.get(worldId, characterId);
        if (!row) {
            return null;
        }
        return this.sessionRepo.getCharacterSessionByName(worldId, row.name);
    }

    private async sessionChannelIndex(worldId: number, characterId: number) {
        const sess = await this.getCharacterSession(worldId, characterId);
        return this.computeGuildMemberChannelIndex(sess);
    }

    private sortGuildMemberModels(memberModels: GuildMemberModel[] | Map<string, GuildMemberModel>, leaderCharacterId: number) {
        const arr = Array.isArray(memberModels) ? [...memberModels] : [...memberModels.values()];
        const leaderId = leaderCharacterId;
        return arr.sort((a, b) => {
            const aLead = a.characterId === leaderId ? 0 : 1;
            const bLead = b.characterId === leaderId ? 0 : 1;
            if (aLead !== bLead) {
                return aLead - bLead;
            }
            const rankA = a.guildRank ?? GuildMemberRank.GUILD_MEMBER_RANK_MEMBER;
            const rankB = b.guildRank ?? GuildMemberRank.GUILD_MEMBER_RANK_MEMBER;
            if (rankA !== rankB) {
                return rankA - rankB;
            }
            const ta = a.joinedAt instanceof Date ? a.joinedAt.getTime() : new Date(a.joinedAt || 0).getTime();
            const tb = b.joinedAt instanceof Date ? b.joinedAt.getTime() : new Date(b.joinedAt || 0).getTime();
            if (ta !== tb) {
                return ta - tb;
            }
            return a.characterId - b.characterId;
        });
    }

    private async guildToPb(
        worldId: number,
        guild: GuildModel,
        memberModels: GuildMemberModel[] | Map<string, GuildMemberModel>
    ): Promise<GuildMessage> {
        const list = this.sortGuildMemberModels(memberModels, guild.leaderCharacterId);
        const members = await Promise.all(
            list.map(async (m): Promise<GuildMember> => {
                const channelIndex = await this.sessionChannelIndex(worldId, m.characterId);
                const memberPb: GuildMember = {
                    worldId,
                    characterId: m.characterId,
                    characterName: m.characterName,
                    level: m.level,
                    classId: m.classId,
                    rank: m.guildRank,
                    channelIndex,
                };
                if (m.allianceRank != null) {
                    memberPb.allianceRank = m.allianceRank;
                }
                return memberPb;
            })
        );
        const guildPb: GuildMessage = {
            worldId,
            guildId: guild.guildId,
            name: guild.name,
            leaderCharacterId: guild.leaderCharacterId,
            revision: guild.revision,
            gp: guild.gp,
            capacity: guild.capacity,
            notice: guild.notice,
            logo: {
                logo: guild.logo.logo,
                logoColor: guild.logo.logoColor,
                logoBg: guild.logo.logoBG,
                logoBgColor: guild.logo.logoBGColor,
            },
            rankTitles: [...guild.rankTitles],
            members,
        };
        if (guild.allianceId != null) {
            guildPb.allianceId = guild.allianceId;
        }
        return guildPb;
    }

    async buildGuildMessage(worldId: number, guild: GuildModel, memberModels: GuildMemberModel[]) {
        return this.guildToPb(worldId, guild, memberModels);
    }

    async getGuild(worldId: number, guildId: number): Promise<GetGuildResult> {
        this.assertWorld(worldId);
        this.assertGuildId(guildId);
        const guild = await this.guildRepo.get(worldId, guildId);
        if (!guild) {
            return { found: false };
        }
        const members = await this.guildMemberRepo.getAll(worldId, String(guildId));
        return { found: true, guild, members: this.sortGuildMemberModels(members, guild.leaderCharacterId) };
    }

    private truncateBulletinField(value: string, maxLen: number): string {
        if (value.length <= maxLen) {
            return value;
        }
        return value.substring(0, maxLen);
    }

    private isValidBulletinIcon(icon: number): boolean {
        if (icon >= 0 && icon <= 2) {
            return true;
        }
        if (icon >= GUILD_BULLETIN_ICON_CASH_MIN && icon <= GUILD_BULLETIN_ICON_CASH_MAX) {
            return true;
        }
        return false;
    }

    private canModifyBulletinContent(
        posterCharacterId: number,
        characterId: number,
        guildRank: GuildMemberRank | undefined
    ): boolean {
        if (posterCharacterId === characterId) {
            return true;
        }
        return this.canChangeMemberRank(guildRank);
    }

    private bulletinThreadToEntry(thread: GuildBulletinBoardThreadModel): GuildBulletinBoardThreadEntry {
        return {
            localThreadId: thread.localThreadId,
            posterCharacterId: thread.posterCharacterId,
            title: thread.title,
            timestampUnixMs: thread.createdAt.getTime(),
            icon: thread.icon,
            replyCount: thread.replyCount ?? 0,
        };
    }

    private async loadGuildBulletinBoardThreadListPage(
        worldId: number,
        guildId: number,
        listStart: number
    ): Promise<{
        threads: GuildBulletinBoardThreadEntry[];
        threadCount: number;
        notice?: GuildBulletinBoardThreadEntry;
        listStart: number;
    }> {
        const offset = Math.max(0, listStart);
        const [threadCount, pageRows, noticeRow] = await Promise.all([
            this.guildBulletinBoardRepo.countThreads(worldId, guildId),
            this.guildBulletinBoardRepo.listThreads(
                worldId,
                guildId,
                offset,
                GUILD_BULLETIN_THREADS_PER_PAGE
            ),
            this.guildBulletinBoardRepo.getNoticeThread(worldId, guildId),
        ]);
        const threads = pageRows.map((t) => this.bulletinThreadToEntry(t));
        const notice = noticeRow ? this.bulletinThreadToEntry(noticeRow) : undefined;
        return {
            threads,
            threadCount,
            notice,
            listStart: offset,
        };
    }

    private async loadBulletinThreadDetail(
        worldId: number,
        guildId: number,
        localThreadId: number
    ): Promise<GuildBulletinBoardThreadDetail | null> {
        const loaded = await this.guildBulletinBoardRepo.getThreadWithReplies(worldId, guildId, localThreadId);
        if (!loaded) {
            return null;
        }
        const replies: GuildBulletinBoardReplyEntry[] = loaded.replies.map((r) => {
            return {
                replyId: r.replyId,
                posterCharacterId: r.posterCharacterId,
                timestampUnixMs: r.createdAt.getTime(),
                content: r.content,
            };
        });
        return {
            localThreadId: loaded.thread.localThreadId,
            posterCharacterId: loaded.thread.posterCharacterId,
            timestampUnixMs: loaded.thread.createdAt.getTime(),
            title: loaded.thread.title,
            body: loaded.thread.body,
            icon: loaded.thread.icon,
            replies,
        };
    }

    private guildBulletinBoardCooldownKey(worldId: number, characterId: number) {
        return redisCacheKey(`w${worldId}:guild-bbs-cooldown:${characterId}`);
    }

    private async setBulletinBoardCooldown(worldId: number, characterId: number) {
        const { client } = this.ctx.getRedisGlobalAccess(worldId);
        const key = this.guildBulletinBoardCooldownKey(worldId, characterId);
        const acquired = await client.set(key, "1", "EX", GUILD_BULLETIN_COOLDOWN_SEC, "NX");
        return acquired === "OK";
    }

    private async unsetBulletinBoardCooldown(worldId: number, characterId: number) {
        const { client } = this.ctx.getRedisGlobalAccess(worldId);
        const key = this.guildBulletinBoardCooldownKey(worldId, characterId);
        await client.del(key);
    }

    private async assertGuildMembership(
        worldId: number,
        characterId: number
    ): Promise<GuildMembershipResult> {
        const state = await this.characterRealtimeStateRepo.get(worldId, characterId);
        if (state?.guildId == null) {
            return { code: messages.GuildErrorCode.GUILD_ERROR_NOT_IN_GUILD };
        }
        const guildId = state.guildId;
        if (!Number.isInteger(guildId) || guildId < 1) {
            return { code: messages.GuildErrorCode.GUILD_ERROR_NOT_IN_GUILD };
        }
        const guild = await this.guildRepo.get(worldId, guildId);
        if (!guild) {
            return { code: messages.GuildErrorCode.GUILD_ERROR_GUILD_NOT_FOUND };
        }
        const member = await this.guildMemberRepo.getItem(worldId, String(guildId), characterId);
        if (!member) {
            return { code: messages.GuildErrorCode.GUILD_ERROR_NOT_IN_GUILD };
        }
        return {
            code: messages.GuildErrorCode.GUILD_ERROR_NONE,
            guildId,
            member,
        };
    }

    async listGuildBulletinBoardThreads(
        worldId: number,
        characterId: number,
        page: number
    ): Promise<ListGuildBulletinBoardThreadsResult> {
        this.assertWorld(worldId);
        this.assertCharacterId(characterId);

        const membership = await this.assertGuildMembership(worldId, characterId);
        if (membership.code !== messages.GuildErrorCode.GUILD_ERROR_NONE) {
            return { ok: false, code: membership.code };
        }

        const listStart = Math.max(0, page) * GUILD_BULLETIN_THREADS_PER_PAGE;
        const listPage = await this.loadGuildBulletinBoardThreadListPage(
            worldId,
            membership.guildId,
            listStart
        );
        return {
            ok: true,
            threads: listPage.threads,
            listStart: listPage.listStart,
            threadCount: listPage.threadCount,
            notice: listPage.notice,
        };
    }

    async showGuildBulletinBoardThread(
        worldId: number,
        characterId: number,
        localThreadId: number
    ): Promise<ShowGuildBulletinBoardThreadResult> {
        this.assertWorld(worldId);
        this.assertCharacterId(characterId);

        const membership = await this.assertGuildMembership(worldId, characterId);
        if (membership.code !== messages.GuildErrorCode.GUILD_ERROR_NONE) {
            return { ok: false, code: membership.code };
        }

        const thread = await this.loadBulletinThreadDetail(worldId, membership.guildId, localThreadId);
        if (!thread) {
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_BULLETIN_THREAD_NOT_FOUND };
        }
        return { ok: true, thread };
    }

    async createGuildBulletinBoardThread(
        worldId: number,
        characterId: number,
        notice: boolean,
        title: string,
        body: string,
        icon: number
    ): Promise<CreateGuildBulletinBoardThreadResult> {
        this.assertWorld(worldId);
        this.assertCharacterId(characterId);

        const membership = await this.assertGuildMembership(worldId, characterId);
        if (membership.code !== messages.GuildErrorCode.GUILD_ERROR_NONE) {
            return { ok: false, code: membership.code };
        }

        if (!this.isValidBulletinIcon(icon)) {
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_INVALID_BULLETIN_CONTENT };
        }
        const normalizedTitle = this.truncateBulletinField(title, MAX_GUILD_BULLETIN_TITLE_LEN);
        const normalizedBody = this.truncateBulletinField(body, MAX_GUILD_BULLETIN_BODY_LEN);

        await using _guildLock = await this.distributedLockService.acquireWorldDataLock(
            worldId,
            `guild:${membership.guildId}`,
        );

        const isCooldown = await this.setBulletinBoardCooldown(worldId, characterId);
        if (!isCooldown) {
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_BULLETIN_COOLDOWN };
        }

        try {
            const created = await this.guildBulletinBoardRepo.createThread(worldId, {
                guildId: membership.guildId,
                posterCharacterId: characterId,
                title: normalizedTitle,
                body: normalizedBody,
                icon,
                isNotice: notice,
            });
            const thread = await this.loadBulletinThreadDetail(worldId, membership.guildId, created.localThreadId);
            if (!thread) {
                return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_UNKNOWN };
            }
            const page = await this.loadGuildBulletinBoardThreadListPage(worldId, membership.guildId, 0);
            return {
                ok: true,
                thread,
                threads: page.threads,
                listStart: page.listStart,
                threadCount: page.threadCount,
                notice: page.notice,
            };
        } catch {
            await this.unsetBulletinBoardCooldown(worldId, characterId);
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_UNKNOWN };
        }
    }

    async updateGuildBulletinBoardThread(
        worldId: number,
        characterId: number,
        localThreadId: number,
        title: string,
        body: string,
        icon: number
    ): Promise<UpdateGuildBulletinBoardThreadResult> {
        this.assertWorld(worldId);
        this.assertCharacterId(characterId);

        const membership = await this.assertGuildMembership(worldId, characterId);
        if (membership.code !== messages.GuildErrorCode.GUILD_ERROR_NONE) {
            return { ok: false, code: membership.code };
        }

        if (!this.isValidBulletinIcon(icon)) {
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_INVALID_BULLETIN_CONTENT };
        }
        const normalizedTitle = this.truncateBulletinField(title, MAX_GUILD_BULLETIN_TITLE_LEN);
        const normalizedBody = this.truncateBulletinField(body, MAX_GUILD_BULLETIN_BODY_LEN);

        const existing = await this.guildBulletinBoardRepo.getThread(worldId, membership.guildId, localThreadId);
        if (!existing) {
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_BULLETIN_THREAD_NOT_FOUND };
        }
        if (!this.canModifyBulletinContent(existing.posterCharacterId, characterId, membership.member.guildRank)) {
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_NOT_AUTHORIZED };
        }
        existing.title = normalizedTitle;
        existing.body = normalizedBody;
        existing.icon = icon;
        existing.updatedAt = new Date();
        const updated = await this.guildBulletinBoardRepo.updateThread(worldId, existing);
        if (!updated) {
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_BULLETIN_THREAD_NOT_FOUND };
        }
        const thread = await this.loadBulletinThreadDetail(worldId, membership.guildId, localThreadId);
        if (!thread) {
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_BULLETIN_THREAD_NOT_FOUND };
        }
        return { ok: true, thread };
    }

    async deleteGuildBulletinBoardThread(
        worldId: number,
        characterId: number,
        localThreadId: number
    ): Promise<DeleteGuildBulletinBoardThreadResult> {
        this.assertWorld(worldId);
        this.assertCharacterId(characterId);

        const membership = await this.assertGuildMembership(worldId, characterId);
        if (membership.code !== messages.GuildErrorCode.GUILD_ERROR_NONE) {
            return { ok: false, code: membership.code };
        }

        const existing = await this.guildBulletinBoardRepo.getThread(worldId, membership.guildId, localThreadId);
        if (!existing) {
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_BULLETIN_THREAD_NOT_FOUND };
        }
        if (!this.canModifyBulletinContent(existing.posterCharacterId, characterId, membership.member.guildRank)) {
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_NOT_AUTHORIZED };
        }
        await this.guildBulletinBoardRepo.deleteThread(worldId, membership.guildId, localThreadId);
        return { ok: true };
    }

    async createGuildBulletinBoardReply(
        worldId: number,
        characterId: number,
        localThreadId: number,
        content: string
    ): Promise<CreateGuildBulletinBoardReplyResult> {
        this.assertWorld(worldId);
        this.assertCharacterId(characterId);

        const membership = await this.assertGuildMembership(worldId, characterId);
        if (membership.code !== messages.GuildErrorCode.GUILD_ERROR_NONE) {
            return { ok: false, code: membership.code };
        }

        const normalizedContent = this.truncateBulletinField(content, MAX_GUILD_BULLETIN_REPLY_LEN);

        await using _guildLock = await this.distributedLockService.acquireWorldDataLock(
            worldId,
            `guild:${membership.guildId}`,
        );

        const isCooldown = await this.setBulletinBoardCooldown(worldId, characterId);
        if (!isCooldown) {
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_BULLETIN_COOLDOWN };
        }

        const created = await this.guildBulletinBoardRepo.createReply(worldId, {
            guildId: membership.guildId,
            localThreadId,
            posterCharacterId: characterId,
            content: normalizedContent,
        });
        if (!created) {
            await this.unsetBulletinBoardCooldown(worldId, characterId);
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_BULLETIN_THREAD_NOT_FOUND };
        }
        const thread = await this.loadBulletinThreadDetail(worldId, membership.guildId, localThreadId);
        if (!thread) {
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_BULLETIN_THREAD_NOT_FOUND };
        }
        return { ok: true, thread };
    }

    async deleteGuildBulletinBoardReply(
        worldId: number,
        characterId: number,
        localThreadId: number,
        replyId: number
    ): Promise<DeleteGuildBulletinBoardReplyResult> {
        this.assertWorld(worldId);
        this.assertCharacterId(characterId);

        const membership = await this.assertGuildMembership(worldId, characterId);
        if (membership.code !== messages.GuildErrorCode.GUILD_ERROR_NONE) {
            return { ok: false, code: membership.code };
        }

        const target = await this.guildBulletinBoardRepo.getReply(
            worldId,
            membership.guildId,
            localThreadId,
            replyId
        );
        if (!target) {
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_BULLETIN_REPLY_NOT_FOUND };
        }
        if (!this.canModifyBulletinContent(target.posterCharacterId, characterId, membership.member.guildRank)) {
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_NOT_AUTHORIZED };
        }
        const deleted = await this.guildBulletinBoardRepo.deleteReply(
            worldId,
            membership.guildId,
            localThreadId,
            replyId
        );
        if (!deleted) {
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_BULLETIN_REPLY_NOT_FOUND };
        }
        const thread = await this.loadBulletinThreadDetail(worldId, membership.guildId, localThreadId);
        if (!thread) {
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_BULLETIN_THREAD_NOT_FOUND };
        }
        return { ok: true, thread };
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
            const loaded = await this.getGuild(worldId, guildId);
            if (loaded.guild) {
                guilds.push(await this.buildGuildMessage(worldId, loaded.guild, loaded.members ?? []));
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
        await this.allianceRepo.evictCache(worldId, allianceId).catch(() => {});
        for (const guildId of guildIds) {
            await this.guildRepo.evictCache(worldId, guildId).catch(() => {});
            await this.guildMemberRepo.evictGroupCache(worldId, String(guildId)).catch(() => {});
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

        await this.allianceRepo.evictCache(worldId, createdAllianceId);
        await this.guildRepo.evictCache(worldId, guildId);
        await this.guildRepo.evictCache(worldId, partnerGuildId);
        await this.guildMemberRepo.evictGroupCache(worldId, String(guildId));
        await this.guildMemberRepo.evictGroupCache(worldId, String(partnerGuildId));

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

        await this.allianceRepo.evictCache(worldId, txResult.allianceId);
        await this.guildRepo.evictCache(worldId, guildId);
        await this.guildMemberRepo.evictGroupCache(worldId, String(guildId));
        for (const gid of otherGuildIds) {
            await this.guildRepo.evictCache(worldId, gid);
            await this.guildMemberRepo.evictGroupCache(worldId, String(gid));
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

    async increaseAllianceCapacity(
        worldId: number,
        characterId: number
    ): Promise<IncreaseAllianceCapacityResult> {
        this.assertWorld(worldId);
        this.assertCharacterId(characterId);

        await using _characterRealtimeLock = await this.distributedLockService.acquireWorldDataLock(
            worldId,
            `character_realtime:${characterId}`,
        );

        const state = await this.characterRealtimeStateRepo.get(worldId, characterId);
        if (state?.guildId == null) {
            return { ok: false, code: messages.AllianceErrorCode.ALLIANCE_ERROR_GUILD_NOT_FOUND };
        }
        const guildId = state.guildId;
        if (!Number.isInteger(guildId) || guildId < 1) {
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
        const allianceId = guild.allianceId;
        if (allianceId == null || allianceId < 1) {
            return { ok: false, code: messages.AllianceErrorCode.ALLIANCE_ERROR_NOT_IN_ALLIANCE };
        }
        if (guild.leaderCharacterId !== characterId) {
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
        if (alliance.leaderCharacterId !== characterId) {
            return { ok: false, code: messages.AllianceErrorCode.ALLIANCE_ERROR_NOT_ALLIANCE_LEADER };
        }

        const requesterMembers = [...(await this.guildMemberRepo.getAll(worldId, String(guildId))).values()];
        const requesterMember = requesterMembers.find((m) => m.characterId === characterId);
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

        await this.allianceRepo.evictCache(worldId, txResult.allianceId);

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

        await using _characterRealtimeLock = await this.distributedLockService.acquireWorldDataLock(
            worldId,
            `character_realtime:${characterId}`,
        );

        const state = await this.characterRealtimeStateRepo.get(worldId, characterId);
        if (state?.guildId == null) {
            return { ok: false, code: messages.AllianceErrorCode.ALLIANCE_ERROR_GUILD_NOT_FOUND };
        }
        const guildId = state.guildId;
        if (!Number.isInteger(guildId) || guildId < 1) {
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
        const allianceId = guild.allianceId;
        if (allianceId == null || allianceId < 1) {
            return { ok: false, code: messages.AllianceErrorCode.ALLIANCE_ERROR_NOT_IN_ALLIANCE };
        }

        await using _allianceLock = await this.distributedLockService.acquireWorldDataLock(
            worldId,
            `alliance:${allianceId}`,
        );

        const alliance = await this.allianceRepo.get(worldId, allianceId);
        if (!alliance) {
            return { ok: false, code: messages.AllianceErrorCode.ALLIANCE_ERROR_ALLIANCE_NOT_FOUND };
        }
        if (alliance.leaderCharacterId !== characterId) {
            return { ok: false, code: messages.AllianceErrorCode.ALLIANCE_ERROR_NOT_ALLIANCE_LEADER };
        }

        const requesterMembers = [...(await this.guildMemberRepo.getAll(worldId, String(guildId))).values()];
        const requesterMember = requesterMembers.find((m) => m.characterId === characterId);
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

        await this.allianceRepo.evictCache(worldId, txResult.allianceId);

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

        await using _characterRealtimeLock = await this.distributedLockService.acquireWorldDataLock(
            worldId,
            `character_realtime:${characterId}`,
        );

        const state = await this.characterRealtimeStateRepo.get(worldId, characterId);
        if (state?.guildId == null) {
            return { ok: false, code: messages.AllianceErrorCode.ALLIANCE_ERROR_GUILD_NOT_FOUND };
        }
        const guildId = state.guildId;
        if (!Number.isInteger(guildId) || guildId < 1) {
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
        const allianceId = guild.allianceId;
        if (allianceId == null || allianceId < 1) {
            return { ok: false, code: messages.AllianceErrorCode.ALLIANCE_ERROR_NOT_IN_ALLIANCE };
        }

        await using _allianceLock = await this.distributedLockService.acquireWorldDataLock(
            worldId,
            `alliance:${allianceId}`,
        );

        const alliance = await this.allianceRepo.get(worldId, allianceId);
        if (!alliance) {
            return { ok: false, code: messages.AllianceErrorCode.ALLIANCE_ERROR_ALLIANCE_NOT_FOUND };
        }

        const requesterMembers = [...(await this.guildMemberRepo.getAll(worldId, String(guildId))).values()];
        const requesterMember = requesterMembers.find((m) => m.characterId === characterId);
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

        await this.allianceRepo.evictCache(worldId, txResult.allianceId);

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

        await using _requesterCharacterLock = await this.distributedLockService.acquireWorldDataLock(
            worldId,
            `character_realtime:${characterId}`,
        );
        await using _newLeaderCharacterLock = await this.distributedLockService.acquireWorldDataLock(
            worldId,
            `character_realtime:${newLeaderCharacterId}`,
        );

        const requesterState = await this.characterRealtimeStateRepo.get(worldId, characterId);
        if (requesterState?.guildId == null) {
            return { ok: false, code: messages.AllianceErrorCode.ALLIANCE_ERROR_GUILD_NOT_FOUND };
        }
        const requesterGuildId = requesterState.guildId;
        if (!Number.isInteger(requesterGuildId) || requesterGuildId < 1) {
            return { ok: false, code: messages.AllianceErrorCode.ALLIANCE_ERROR_GUILD_NOT_FOUND };
        }

        await using _requesterGuildLock = await this.distributedLockService.acquireWorldDataLock(
            worldId,
            `guild:${requesterGuildId}`,
        );

        const requesterGuild = await this.guildRepo.get(worldId, requesterGuildId);
        if (!requesterGuild) {
            return { ok: false, code: messages.AllianceErrorCode.ALLIANCE_ERROR_GUILD_NOT_FOUND };
        }
        const allianceId = requesterGuild.allianceId;
        if (allianceId == null || allianceId < 1) {
            return { ok: false, code: messages.AllianceErrorCode.ALLIANCE_ERROR_NOT_IN_ALLIANCE };
        }

        await using _allianceLock = await this.distributedLockService.acquireWorldDataLock(
            worldId,
            `alliance:${allianceId}`,
        );

        const alliance = await this.allianceRepo.get(worldId, allianceId);
        if (!alliance) {
            return { ok: false, code: messages.AllianceErrorCode.ALLIANCE_ERROR_ALLIANCE_NOT_FOUND };
        }
        if (alliance.leaderCharacterId !== characterId) {
            return { ok: false, code: messages.AllianceErrorCode.ALLIANCE_ERROR_NOT_ALLIANCE_LEADER };
        }

        const requesterMembers = [...(await this.guildMemberRepo.getAll(worldId, String(requesterGuildId))).values()];
        const requesterMember = requesterMembers.find((m) => m.characterId === characterId);
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

        await this.allianceRepo.evictCache(worldId, txResult.allianceId);
        await this.guildRepo.evictCache(worldId, txResult.oldLeaderGuildId);
        await this.guildMemberRepo.evictGroupCache(worldId, String(txResult.oldLeaderGuildId));
        await this.guildRepo.evictCache(worldId, txResult.newLeaderGuildId);
        await this.guildMemberRepo.evictGroupCache(worldId, String(txResult.newLeaderGuildId));
        for (const gid of otherGuildIds) {
            if (gid === txResult.oldLeaderGuildId || gid === txResult.newLeaderGuildId) {
                continue;
            }
            await this.guildRepo.evictCache(worldId, gid);
            await this.guildMemberRepo.evictGroupCache(worldId, String(gid));
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

        const firstCharacterId = requesterCharacterId < targetCharacterId
            ? requesterCharacterId
            : targetCharacterId;
        const secondCharacterId = requesterCharacterId < targetCharacterId
            ? targetCharacterId
            : requesterCharacterId;

        await using _firstCharacterRealtimeLock = await this.distributedLockService.acquireWorldDataLock(
            worldId,
            `character_realtime:${firstCharacterId}`,
        );

        await using _secondCharacterRealtimeLock = await this.distributedLockService.acquireWorldDataLock(
            worldId,
            `character_realtime:${secondCharacterId}`,
        );

        const requesterState = await this.characterRealtimeStateRepo.get(worldId, requesterCharacterId);
        if (requesterState?.guildId == null) {
            return { ok: false, code: messages.AllianceErrorCode.ALLIANCE_ERROR_GUILD_NOT_FOUND };
        }
        const requesterGuildId = requesterState.guildId;
        if (!Number.isInteger(requesterGuildId) || requesterGuildId < 1) {
            return { ok: false, code: messages.AllianceErrorCode.ALLIANCE_ERROR_GUILD_NOT_FOUND };
        }

        await using _requesterGuildLock = await this.distributedLockService.acquireWorldDataLock(
            worldId,
            `guild:${requesterGuildId}`,
        );

        const requesterGuild = await this.guildRepo.get(worldId, requesterGuildId);
        if (!requesterGuild) {
            return { ok: false, code: messages.AllianceErrorCode.ALLIANCE_ERROR_GUILD_NOT_FOUND };
        }
        const allianceId = requesterGuild.allianceId;
        if (allianceId == null || allianceId < 1) {
            return { ok: false, code: messages.AllianceErrorCode.ALLIANCE_ERROR_NOT_IN_ALLIANCE };
        }

        await using _allianceLock = await this.distributedLockService.acquireWorldDataLock(
            worldId,
            `alliance:${allianceId}`,
        );

        const alliance = await this.allianceRepo.get(worldId, allianceId);
        if (!alliance) {
            return { ok: false, code: messages.AllianceErrorCode.ALLIANCE_ERROR_ALLIANCE_NOT_FOUND };
        }

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

        await this.allianceRepo.evictCache(worldId, txResult.allianceId);
        const guildIdsToEvict = new Set<number>([requesterGuildId]);
        if (txResult.targetGuildId != null) {
            guildIdsToEvict.add(txResult.targetGuildId);
        }
        for (const gid of guildIdsToEvict) {
            await this.guildRepo.evictCache(worldId, gid);
            await this.guildMemberRepo.evictGroupCache(worldId, String(gid));
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

        await using _characterRealtimeLock = await this.distributedLockService.acquireWorldDataLock(
            worldId,
            `character_realtime:${characterId}`,
        );

        const state = await this.characterRealtimeStateRepo.get(worldId, characterId);
        if (state?.guildId == null) {
            return { ok: false, code: messages.AllianceErrorCode.ALLIANCE_ERROR_GUILD_NOT_FOUND };
        }
        const guildId = state.guildId;
        if (!Number.isInteger(guildId) || guildId < 1) {
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
        const allianceId = guild.allianceId;
        if (allianceId == null || allianceId < 1) {
            return { ok: false, code: messages.AllianceErrorCode.ALLIANCE_ERROR_NOT_IN_ALLIANCE };
        }

        await using _allianceLock = await this.distributedLockService.acquireWorldDataLock(
            worldId,
            `alliance:${allianceId}`,
        );

        const alliance = await this.allianceRepo.get(worldId, allianceId);
        if (!alliance) {
            return { ok: false, code: messages.AllianceErrorCode.ALLIANCE_ERROR_ALLIANCE_NOT_FOUND };
        }

        const requesterMembers = [...(await this.guildMemberRepo.getAll(worldId, String(guildId))).values()];
        const requesterMember = requesterMembers.find((m) => m.characterId === characterId);
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

            await this.allianceRepo.evictCache(worldId, result.allianceId);
            for (const gid of result.guildIds ?? []) {
                await this.guildRepo.evictCache(worldId, gid);
                await this.guildMemberRepo.evictGroupCache(worldId, String(gid));
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

        await this.allianceRepo.evictCache(worldId, txResult.allianceId);
        await this.guildRepo.evictCache(worldId, guildId);
        await this.guildMemberRepo.evictGroupCache(worldId, String(guildId));
        for (const gid of otherGuildIds) {
            await this.guildRepo.evictCache(worldId, gid);
            await this.guildMemberRepo.evictGroupCache(worldId, String(gid));
        }

        const allianceMessage = await this.allianceToPb(worldId, txResult.savedAlliance);
        const removedGuildLoaded = await this.getGuild(worldId, guildId);
        const removedGuildMessage =
            removedGuildLoaded.guild != null
                ? await this.buildGuildMessage(worldId, removedGuildLoaded.guild, removedGuildLoaded.members ?? [])
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

        await using _characterRealtimeLock = await this.distributedLockService.acquireWorldDataLock(
            worldId,
            `character_realtime:${characterId}`,
        );

        const state = await this.characterRealtimeStateRepo.get(worldId, characterId);
        if (state?.guildId == null) {
            return { ok: false, code: messages.AllianceErrorCode.ALLIANCE_ERROR_GUILD_NOT_FOUND };
        }
        const requesterGuildId = state.guildId;
        if (!Number.isInteger(requesterGuildId) || requesterGuildId < 1) {
            return { ok: false, code: messages.AllianceErrorCode.ALLIANCE_ERROR_GUILD_NOT_FOUND };
        }
        if (targetGuildId === requesterGuildId) {
            return { ok: false, code: messages.AllianceErrorCode.ALLIANCE_ERROR_CANNOT_EXPEL_OWN_GUILD };
        }

        await using _requesterGuildLock = await this.distributedLockService.acquireWorldDataLock(
            worldId,
            `guild:${requesterGuildId}`,
        );

        const requesterGuild = await this.guildRepo.get(worldId, requesterGuildId);
        if (!requesterGuild) {
            return { ok: false, code: messages.AllianceErrorCode.ALLIANCE_ERROR_GUILD_NOT_FOUND };
        }
        const allianceId = requesterGuild.allianceId;
        if (allianceId == null || allianceId < 1) {
            return { ok: false, code: messages.AllianceErrorCode.ALLIANCE_ERROR_NOT_IN_ALLIANCE };
        }
        if (clientAllianceId != null && clientAllianceId > 0 && clientAllianceId !== allianceId) {
            return { ok: false, code: messages.AllianceErrorCode.ALLIANCE_ERROR_ALLIANCE_NOT_FOUND };
        }

        await using _allianceLock = await this.distributedLockService.acquireWorldDataLock(
            worldId,
            `alliance:${allianceId}`,
        );

        const alliance = await this.allianceRepo.get(worldId, allianceId);
        if (!alliance) {
            return { ok: false, code: messages.AllianceErrorCode.ALLIANCE_ERROR_ALLIANCE_NOT_FOUND };
        }
        if (!alliance.guildIds.includes(targetGuildId)) {
            return { ok: false, code: messages.AllianceErrorCode.ALLIANCE_ERROR_TARGET_NOT_IN_ALLIANCE };
        }

        const requesterMembers = [...(await this.guildMemberRepo.getAll(worldId, String(requesterGuildId))).values()];
        const requesterMember = requesterMembers.find((m) => m.characterId === characterId);
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

            await this.allianceRepo.evictCache(worldId, result.allianceId);
            for (const gid of result.guildIds ?? []) {
                await this.guildRepo.evictCache(worldId, gid);
                await this.guildMemberRepo.evictGroupCache(worldId, String(gid));
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

        await this.allianceRepo.evictCache(worldId, txResult.allianceId);
        await this.guildRepo.evictCache(worldId, targetGuildId);
        await this.guildMemberRepo.evictGroupCache(worldId, String(targetGuildId));
        for (const gid of otherGuildIds) {
            if (gid === targetGuildId) {
                continue;
            }
            await this.guildRepo.evictCache(worldId, gid);
            await this.guildMemberRepo.evictGroupCache(worldId, String(gid));
        }

        const allianceMessage = await this.allianceToPb(worldId, txResult.savedAlliance);
        const removedGuildLoaded = await this.getGuild(worldId, targetGuildId);
        const removedGuildMessage =
            removedGuildLoaded.guild != null
                ? await this.buildGuildMessage(worldId, removedGuildLoaded.guild, removedGuildLoaded.members ?? [])
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

        await using _characterRealtimeLock = await this.distributedLockService.acquireWorldDataLock(
            worldId,
            `character_realtime:${characterId}`,
        );

        const state = await this.characterRealtimeStateRepo.get(worldId, characterId);
        if (state?.guildId == null) {
            return { ok: false, code: messages.AllianceErrorCode.ALLIANCE_ERROR_GUILD_NOT_FOUND };
        }
        const guildId = state.guildId;
        if (!Number.isInteger(guildId) || guildId < 1) {
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
        const allianceId = guild.allianceId;
        if (allianceId == null || allianceId < 1) {
            return { ok: false, code: messages.AllianceErrorCode.ALLIANCE_ERROR_NOT_IN_ALLIANCE };
        }

        await using _allianceLock = await this.distributedLockService.acquireWorldDataLock(
            worldId,
            `alliance:${allianceId}`,
        );

        const alliance = await this.allianceRepo.get(worldId, allianceId);
        if (!alliance) {
            return { ok: false, code: messages.AllianceErrorCode.ALLIANCE_ERROR_ALLIANCE_NOT_FOUND };
        }
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

        await this.allianceRepo.evictCache(worldId, result.allianceId);
        for (const gid of result.guildIds ?? []) {
            await this.guildRepo.evictCache(worldId, gid);
            await this.guildMemberRepo.evictGroupCache(worldId, String(gid));
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

    async broadcastGuildMultiChat(
        worldId: number,
        guildId: number,
        senderCharacterId: number,
        senderName: string,
        message: string
    ): Promise<{ ok: boolean; deliveredCount?: number }> {
        this.assertWorld(worldId);
        this.assertGuildId(guildId);
        this.assertCharacterId(senderCharacterId);
        const trimmedMsg = (message ?? "").trim();
        if (trimmedMsg.length <= 0 || trimmedMsg.length > 500) {
            return { ok: false };
        }
        const name = (senderName ?? "").trim();
        if (!name) {
            return { ok: false };
        }
        const loaded = await this.getGuild(worldId, guildId);
        if (!loaded.guild) {
            return { ok: false };
        }
        const members = loaded.members ?? [];
        const sender = members.find((m) => m.characterId === senderCharacterId);
        if (!sender) {
            return { ok: false };
        }
        const guild = loaded.guild;
        await this.publishToGuildRoutes("multi_chat", worldId, guildId, guild.revision, {
            guild_id: guildId,
            sender_character_id: senderCharacterId,
            chat_mode: 2,
            sender_name: name,
            message: trimmedMsg,
        });
        return { ok: true, deliveredCount: 1 };
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
