import {
    Guild,
    GuildErrorCode,
    GuildMemberRank,
    type GuildMember,
    type GuildLogo as GuildLogoMessage,
    type Guild as GuildMessage,
} from "@fm/shared/protobuf/generated/fminternal/internal_service";
import type { PoolClient } from "pg";
import { AppConfiguration } from "@fm/shared/config/app-configuration";
import { InternalContext } from "@fm/shared/context/internal-context";
import { UnifiedRepository } from "@fm/shared/repos/unified-repository";
import {
    GuildChangeRequestRepository,
    GuildCreateRequestRepository,
    GuildRepository,
} from "@fm/shared/repos/guild-repository";
import { RequestStatus } from "@fm/shared/repos/request-repository";
import { GuildMemberRepository } from "@fm/shared/repos/guild-member-repository";
import { GuildBulletinBoardRepository } from "@fm/shared/repos/guild-bulletin-board-repository";
import { CharacterRealtimeStateRepository } from "@fm/shared/repos/character-realtime-state-repository";
import { SessionRepository } from "@fm/shared/repos/session-repository";
import type { GuildModel } from "@fm/shared/repos/guild-repository";
import { guildWeek } from "@fm/shared/guild-week";
import type { GuildMemberModel } from "@fm/shared/repos/guild-member-repository";
import { RabbitMQService } from "./rabbitmq-service";
import { DistributedLockService } from "./distributed-lock-service";
import { WzService } from "./wz-service";
import { redisCacheKey } from "@fm/shared/redis-cache-key";
import {
    DEFAULT_GUILD_LOGO,
    DEFAULT_GUILD_RANK_TITLES,
    type GuildLogo,
    type GuildRankTitles,
} from "@fm/shared/types/guild-json";

export type {
    ListGuildBulletinBoardThreadsResult,
    ShowGuildBulletinBoardThreadResult,
    CreateGuildBulletinBoardThreadResult,
    UpdateGuildBulletinBoardThreadResult,
    DeleteGuildBulletinBoardThreadResult,
    CreateGuildBulletinBoardReplyResult,
    DeleteGuildBulletinBoardReplyResult,
    GuildMembershipResult,
} from "./guild-bulletin-board-service";

export type {
    CreateAllianceResult,
    GetAllianceResult,
    DisbandAllianceResult,
    LeaveAllianceResult,
    ExpelAllianceGuildResult,
    AcceptAllianceInviteResult,
    IncreaseAllianceCapacityResult,
    ChangeAllianceRankTitlesResult,
    ChangeAllianceMemberRankResult,
    ChangeAllianceLeaderResult,
    ChangeAllianceNoticeResult,
} from "./alliance-service";

const messages = { GuildErrorCode };

const DEFAULT_GUILD_CAPACITY = 10;
const MAX_GUILD_NOTICE_LEN = 100;
const GUILD_CAPACITY_STEP = 5;
const GUILD_CAPACITY_STANDARD_MAX = 100;
const GUILD_CAPACITY_EXTENDED_MAX = 200;
const GUILD_CAPACITY_EXTENDED_GP_COST = 2000;
const GUILD_GP_MAX = 2147483647;
const GUILD_RANKING_SIZE = 50;
const GUILD_RANKING_TTL_SEC = 60;
const DEFAULT_ALLIANCE_MEMBER_RANK = 5;

const EVT = {
    CREATED: "created",
    MEMBER_JOINED: "member_joined",
    MEMBER_LEFT: "member_left",
    RANK_TITLES_CHANGED: "rank_titles_changed",
    MEMBER_RANK_CHANGED: "member_rank_changed",
    EMBLEM_CHANGED: "emblem_changed",
    NOTICE_CHANGED: "notice_changed",
    CAPACITY_CHANGED: "capacity_changed",
    GP_CHANGED: "gp_changed",
    MEMBER_ONLINE_CHANGED: "member_online_changed",
    DISBANDED: "disbanded",
    MESSAGE: "message",
} as const;

export type CreateGuildResult = {
    ok: boolean;
    code?: GuildErrorCode;
    guildId?: number;
    guild?: GuildMessage;
};

export type GetGuildResult = { found: boolean; guild?: GuildModel; members?: GuildMemberModel[] };

export type AcceptGuildInviteResult = {
    ok: boolean;
    code?: GuildErrorCode;
    guild?: GuildMessage;
};

export type LeaveGuildResult = {
    ok: boolean;
    code?: GuildErrorCode;
    guildId?: number;
};

export type ExpelGuildResult = {
    ok: boolean;
    code?: GuildErrorCode;
    guildId?: number;
};

export type ChangeGuildRankTitlesResult = {
    ok: boolean;
    code?: GuildErrorCode;
    guildId?: number;
};

export type ChangeGuildMemberRankResult = {
    ok: boolean;
    code?: GuildErrorCode;
    guildId?: number;
};

export type ChangeGuildEmblemResult = {
    ok: boolean;
    code?: GuildErrorCode;
    guildId?: number;
};

export type ChangeGuildNoticeResult = {
    ok: boolean;
    code?: GuildErrorCode;
    guildId?: number;
};

export type DisbandGuildResult = {
    ok: boolean;
    code?: GuildErrorCode;
    guildId?: number;
    allianceId?: number | null;
};

export type IncreaseGuildCapacityResult = {
    ok: boolean;
    code?: GuildErrorCode;
    guildId?: number;
    capacity?: number;
    gp?: number;
};

export type GainGuildGPResult = {
    ok: boolean;
    code?: GuildErrorCode;
    guildId?: number;
    gp?: number;
};

export type GuildRankingEntryResult = {
    guildId: number;
    name: string;
    gp: number;
    logo: GuildLogo;
};

export class GuildService {
    private readonly ctx: InternalContext;
    private readonly app: AppConfiguration;
    private readonly unifiedRepo: UnifiedRepository;
    private readonly guildRepo: GuildRepository;
    private readonly guildCreateRequestRepo: GuildCreateRequestRepository;
    private readonly guildChangeRequestRepo: GuildChangeRequestRepository;
    private readonly guildMemberRepo: GuildMemberRepository;
    private readonly guildBulletinBoardRepo: GuildBulletinBoardRepository;
    private readonly characterRealtimeStateRepo: CharacterRealtimeStateRepository;
    private readonly sessionRepo: SessionRepository;
    private readonly rabbitmqService: RabbitMQService;
    private readonly distributedLockService: DistributedLockService;
    private readonly wzService: WzService;

    constructor(
        internalContext: InternalContext,
        appConfiguration: AppConfiguration,
        unifiedRepository: UnifiedRepository,
        guildRepository: GuildRepository,
        guildCreateRequestRepository: GuildCreateRequestRepository,
        guildChangeRequestRepository: GuildChangeRequestRepository,
        guildMemberRepository: GuildMemberRepository,
        guildBulletinBoardRepository: GuildBulletinBoardRepository,
        characterRealtimeStateRepository: CharacterRealtimeStateRepository,
        sessionRepository: SessionRepository,
        rabbitmqService: RabbitMQService,
        distributedLockService: DistributedLockService,
        wzService: WzService,
    ) {
        this.ctx = internalContext;
        this.app = appConfiguration;
        this.unifiedRepo = unifiedRepository;
        this.guildRepo = guildRepository;
        this.guildCreateRequestRepo = guildCreateRequestRepository;
        this.guildChangeRequestRepo = guildChangeRequestRepository;
        this.guildMemberRepo = guildMemberRepository;
        this.guildBulletinBoardRepo = guildBulletinBoardRepository;
        this.characterRealtimeStateRepo = characterRealtimeStateRepository;
        this.sessionRepo = sessionRepository;
        this.rabbitmqService = rabbitmqService;
        this.distributedLockService = distributedLockService;
        this.wzService = wzService;
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

    private async publishToGuildRoutes(
        eventType: string,
        worldId: number,
        guildId: number,
        extraPayload: Record<string, unknown> = {},
    ) {
        const routingKey = `fm.${worldId}.all.guild`;
        return this.rabbitmqService.publish(routingKey, eventType, {
            world_id: worldId,
            guild_id: guildId,
            ...extraPayload,
        });
    }

    async guildSnapshot(worldId: number, guildId: number): Promise<string | null> {
        const loaded = await this.getGuild(worldId, guildId);
        if (!loaded.found || !loaded.guild) {
            return null;
        }
        const guildPb = await this.guildToPb(worldId, loaded.guild, loaded.members ?? []);
        return Buffer.from(Guild.encode(guildPb).finish()).toString("base64");
    }

    private async publishGuildChange(
        eventType: string,
        worldId: number,
        guildId: number,
        extraPayload: Record<string, unknown> = {},
    ) {
        const guildPb = await this.guildSnapshot(worldId, guildId);
        if (guildPb == null) {
            return;
        }
        await this.publishToGuildRoutes(eventType, worldId, guildId, { ...extraPayload, guild_pb: guildPb });
    }

    private async rollbackCreateGuild(worldId: number, guildId: number, leaderCharacterId: number) {
        await this.ctx
            .withPgDataTransaction(worldId, guildId, async (dataTx: PoolClient) => {
                await this.guildMemberRepo.deleteAllForGuild(worldId, guildId, { txClient: dataTx });
                await this.guildBulletinBoardRepo.deleteAllForGuild(worldId, guildId, { txClient: dataTx });
                await this.guildRepo.delete({ worldId, guildId }, { txClient: dataTx });
            })
            .catch(() => {});
        await this.unifiedRepo.deleteGuildName(guildId).catch(() => {});
        await this.guildRepo.invalidateCache(worldId, guildId).catch(() => {});
        await this.guildMemberRepo.invalidateCache(worldId, String(guildId)).catch(() => {});
        await this.characterRealtimeStateRepo.invalidateCache(worldId, leaderCharacterId).catch(() => {});
    }

    async createGuild(
        worldId: number,
        guildName: string,
        leader: GuildMember | null | undefined,
        requestId: string,
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

        if (this.wzService.isValidName(guildName) === false) {
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_GUILD_NAME_INVALID };
        }

        await using _characterRealtimeLock = await this.distributedLockService.acquireWorldDataLock(
            worldId,
            `character_realtime:${leaderCharacterId}`,
        );

        const recorded =
            requestId === "0" ? null : await this.guildCreateRequestRepo.find(worldId, leaderCharacterId, requestId);
        if (recorded?.status === RequestStatus.Applied) {
            const recordedGuildId = recorded.outcome.guild_id;
            const loaded = await this.getGuild(worldId, recordedGuildId);
            if (!loaded?.guild) {
                return { ok: true, guildId: recordedGuildId };
            }
            return {
                ok: true,
                guildId: recordedGuildId,
                guild: await this.guildToPb(worldId, loaded.guild, loaded.members ?? []),
            };
        }

        const leaderState = await this.characterRealtimeStateRepo.get(worldId, leaderCharacterId);
        if (leaderState?.guildId != null) {
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_ALREADY_IN_GUILD };
        }

        let guildId: number;
        try {
            const nameEntry = await this.unifiedRepo.reserveGuildName(guildName, worldId);
            guildId = Number(nameEntry.guild_id);
        } catch (err: unknown) {
            const code = (err as { code?: string })?.code;
            if (code === "23505") {
                return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_GUILD_NAME_TAKEN };
            }
            throw err;
        }

        let savedGuild: GuildModel;
        try {
            savedGuild = await this.ctx.withPgDataTransaction(worldId, guildId, async (dataTx: PoolClient) => {
                const savedGuild = await this.guildRepo.set(
                    worldId,
                    {
                        worldId,
                        guildId,
                        name: guildName,
                        leaderCharacterId,
                        gp: 0,
                        weeklyGp: 0,
                        gpWeek: 0,
                        capacity: DEFAULT_GUILD_CAPACITY,
                        notice: "",
                        logo: DEFAULT_GUILD_LOGO,
                        rankTitles: DEFAULT_GUILD_RANK_TITLES,
                    },
                    { txClient: dataTx },
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
                    { txClient: dataTx },
                );
                return savedGuild;
            });

            await this.ctx.withPgDataTransaction(worldId, leaderCharacterId, async (dataTx: PoolClient) => {
                const state = await this.characterRealtimeStateRepo.get(worldId, leaderCharacterId, {
                    txClient: dataTx,
                });
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
                    { txClient: dataTx },
                );
                if (requestId !== "0") {
                    await this.guildCreateRequestRepo.insert(
                        worldId,
                        leaderCharacterId,
                        requestId,
                        { guild_id: guildId },
                        dataTx,
                    );
                }
            });
        } catch (err: unknown) {
            await this.rollbackCreateGuild(worldId, guildId, leaderCharacterId);
            if ((err as { guildRollback?: boolean }).guildRollback) {
                return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_ALREADY_IN_GUILD };
            }
            throw err;
        }

        await this.guildRepo.invalidateCache(worldId, guildId);
        await this.guildMemberRepo.invalidateCache(worldId, String(guildId));
        await this.characterRealtimeStateRepo.invalidateCache(worldId, leaderCharacterId);

        await this.publishGuildChange(EVT.CREATED, worldId, guildId, {
            leader_character_id: leaderCharacterId,
            guild_name: guildName,
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
        const guildMessage = await this.guildToPb(worldId, savedGuild, [leaderMember]);

        return { ok: true, guildId, guild: guildMessage };
    }

    async acceptGuildInvite(
        worldId: number,
        guildId: number,
        member: GuildMember | null | undefined,
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

        await using _guildLock = await this.distributedLockService.acquireWorldDataLock(worldId, `guild:${guildId}`);

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
                    allianceRank: guild.allianceId != null ? DEFAULT_ALLIANCE_MEMBER_RANK : null,
                },
                { txClient },
            );

            const updatedGuild = await this.guildRepo.set(worldId, guild, { txClient });

            return { ok: true as const };
        });

        if (!result.ok) {
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
                { txClient },
            );
        });

        await this.guildRepo.invalidateCache(worldId, guildId);
        await this.guildMemberRepo.invalidateCache(worldId, String(guildId));
        await this.characterRealtimeStateRepo.invalidateCache(worldId, characterId);

        await this.publishGuildChange(EVT.MEMBER_JOINED, worldId, guildId, {
            character_id: characterId,
        });

        const loaded = await this.getGuild(worldId, guildId);
        if (!loaded.found || !loaded.guild || !loaded.members) {
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_UNKNOWN };
        }

        const guildMessage = await this.guildToPb(worldId, loaded.guild, loaded.members);
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

        const guildResult = await this.ctx.withPgDataTransaction(
            worldId,
            lockedGuildId,
            async (txClient: PoolClient) => {
                const guild = await this.guildRepo.get(worldId, lockedGuildId, { txClient });
                if (!guild) {
                    return { ok: true as const, guildId: lockedGuildId, clearRealtime: true as const };
                }

                if (guild.leaderCharacterId === characterId) {
                    return { ok: false as const, code: messages.GuildErrorCode.GUILD_ERROR_LEADER_CANNOT_LEAVE };
                }

                const members = await this.guildMemberRepo.getAll(worldId, String(lockedGuildId), { txClient });
                if (!members.has(String(characterId))) {
                    return { ok: false as const, code: messages.GuildErrorCode.GUILD_ERROR_NOT_IN_GUILD };
                }

                await this.guildMemberRepo.del(worldId, String(lockedGuildId), characterId, { txClient });

                const updatedGuild = await this.guildRepo.set(worldId, guild, { txClient });
                return { ok: true as const, guildId: updatedGuild.guildId, clearRealtime: true as const };
            },
        );

        if (!guildResult.ok || guildResult.guildId == null) {
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
                    { txClient },
                );
            });
        }
        const result = guildResult;

        await this.guildRepo.invalidateCache(worldId, result.guildId);
        await this.guildMemberRepo.invalidateCache(worldId, String(result.guildId));
        await this.characterRealtimeStateRepo.invalidateCache(worldId, characterId);

        await this.publishGuildChange(EVT.MEMBER_LEFT, worldId, result.guildId, {
            character_id: characterId,
        });

        return { ok: true, guildId: result.guildId };
    }

    private canExpelGuildMembers(rank: GuildMemberRank | undefined) {
        return rank === GuildMemberRank.GUILD_MEMBER_RANK_MASTER || rank === GuildMemberRank.GUILD_MEMBER_RANK_JUNIOR;
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
        targetCharacterId: number,
    ): Promise<ExpelGuildResult> {
        this.assertWorld(worldId);
        this.assertCharacterId(requesterCharacterId);
        this.assertCharacterId(targetCharacterId);
        if (requesterCharacterId === targetCharacterId) {
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_CANNOT_EXPEL_TARGET };
        }

        const firstCharacterId = requesterCharacterId < targetCharacterId ? requesterCharacterId : targetCharacterId;
        const secondCharacterId = requesterCharacterId < targetCharacterId ? targetCharacterId : requesterCharacterId;

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

        const guildResult = await this.ctx.withPgDataTransaction(
            worldId,
            lockedGuildId,
            async (txClient: PoolClient) => {
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

                const updatedGuild = await this.guildRepo.set(worldId, guild, { txClient });
                return { ok: true as const, guildId: updatedGuild.guildId };
            },
        );

        if (!guildResult.ok || guildResult.guildId == null) {
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
                { txClient },
            );
        });
        const result = guildResult;

        await this.guildRepo.invalidateCache(worldId, result.guildId);
        await this.guildMemberRepo.invalidateCache(worldId, String(result.guildId));
        await this.characterRealtimeStateRepo.invalidateCache(worldId, targetCharacterId);

        await this.publishGuildChange(EVT.MEMBER_LEFT, worldId, result.guildId, {
            character_id: targetCharacterId,
            expelled: true,
        });

        if ((await this.sessionRepo.findChannel(worldId, targetCharacterId)) == null) {
            // TODO: game server -> internal sendNote RPC (promise) when offline member is expelled
        }

        return { ok: true, guildId: result.guildId };
    }

    private normalizeRankTitles(rankTitles: string[] | null | undefined): GuildRankTitles | null {
        if (!rankTitles || rankTitles.length !== 5) {
            return null;
        }
        const titles = rankTitles.map((entry) => String(entry ?? ""));
        if (this.wzService.isValidRankTitles(titles) === false) {
            return null;
        }
        return titles as GuildRankTitles;
    }

    async changeGuildRankTitles(
        worldId: number,
        characterId: number,
        rankTitles: string[],
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

            const updatedGuild = await this.guildRepo.set(
                worldId,
                { ...guild, rankTitles: normalizedRankTitles },
                { txClient },
            );
            return { ok: true as const, guildId: updatedGuild.guildId };
        });

        if (!result.ok || result.guildId == null) {
            return result;
        }

        await this.guildRepo.invalidateCache(worldId, result.guildId);

        await this.publishGuildChange(EVT.RANK_TITLES_CHANGED, worldId, result.guildId, {});

        return { ok: true, guildId: result.guildId };
    }

    private canChangeMemberRank(rank: GuildMemberRank | undefined) {
        return rank === GuildMemberRank.GUILD_MEMBER_RANK_MASTER || rank === GuildMemberRank.GUILD_MEMBER_RANK_JUNIOR;
    }

    private parseAssignableMemberRank(newRank: GuildMemberRank): GuildMemberRank | null {
        if (
            newRank === GuildMemberRank.GUILD_MEMBER_RANK_JUNIOR ||
            newRank === GuildMemberRank.GUILD_MEMBER_RANK_SENIOR ||
            newRank === GuildMemberRank.GUILD_MEMBER_RANK_MEMBER ||
            newRank === GuildMemberRank.GUILD_MEMBER_RANK_NEW
        ) {
            return newRank;
        }
        return null;
    }

    async changeGuildMemberRank(
        worldId: number,
        requesterCharacterId: number,
        targetCharacterId: number,
        newRank: GuildMemberRank,
    ): Promise<ChangeGuildMemberRankResult> {
        this.assertWorld(worldId);
        this.assertCharacterId(requesterCharacterId);
        this.assertCharacterId(targetCharacterId);

        const parsedRank = this.parseAssignableMemberRank(newRank);
        if (!parsedRank) {
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_INVALID_MEMBER_RANK };
        }

        const firstCharacterId = requesterCharacterId < targetCharacterId ? requesterCharacterId : targetCharacterId;
        const secondCharacterId = requesterCharacterId < targetCharacterId ? targetCharacterId : requesterCharacterId;

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
            if (
                parsedRank === GuildMemberRank.GUILD_MEMBER_RANK_JUNIOR &&
                requesterMember.guildRank !== GuildMemberRank.GUILD_MEMBER_RANK_MASTER
            ) {
                return { ok: false as const, code: messages.GuildErrorCode.GUILD_ERROR_NOT_AUTHORIZED };
            }

            const targetMember = members.get(String(targetCharacterId));
            if (!targetMember) {
                return { ok: false as const, code: messages.GuildErrorCode.GUILD_ERROR_TARGET_NOT_IN_GUILD };
            }
            if (
                this.guildMemberRankValue(requesterMember.guildRank) >=
                this.guildMemberRankValue(targetMember.guildRank)
            ) {
                return { ok: false as const, code: messages.GuildErrorCode.GUILD_ERROR_NOT_AUTHORIZED };
            }

            await this.guildMemberRepo.set(worldId, { ...targetMember, guildRank: parsedRank }, { txClient });

            const updatedGuild = await this.guildRepo.set(worldId, guild, { txClient });
            return { ok: true as const, guildId: updatedGuild.guildId };
        });

        if (!result.ok || result.guildId == null) {
            return result;
        }

        await this.guildRepo.invalidateCache(worldId, result.guildId);
        await this.guildMemberRepo.invalidateCache(worldId, String(result.guildId));

        await this.publishGuildChange(EVT.MEMBER_RANK_CHANGED, worldId, result.guildId, {
            character_id: targetCharacterId,
        });

        return { ok: true, guildId: result.guildId };
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
        logo: GuildLogoMessage | null | undefined,
        requestId: string,
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

        if (requestId !== "0" && (await this.guildChangeRequestRepo.find(worldId, lockedGuildId, requestId)) != null) {
            return { ok: true, guildId: lockedGuildId };
        }

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

            const updatedGuild = await this.guildRepo.set(worldId, { ...guild, logo: normalizedLogo }, { txClient });
            if (requestId !== "0") {
                await this.guildChangeRequestRepo.insert(worldId, lockedGuildId, requestId, {}, txClient);
            }
            return { ok: true as const, guildId: updatedGuild.guildId };
        });

        if (!result.ok || result.guildId == null) {
            return result;
        }

        await this.guildRepo.invalidateCache(worldId, result.guildId);

        await this.publishGuildChange(EVT.EMBLEM_CHANGED, worldId, result.guildId, {});

        return { ok: true, guildId: result.guildId };
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
        notice: string | null | undefined,
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

            const updatedGuild = await this.guildRepo.set(
                worldId,
                { ...guild, notice: normalizedNotice },
                { txClient },
            );
            return { ok: true as const, guildId: updatedGuild.guildId };
        });

        if (!result.ok || result.guildId == null) {
            return result;
        }

        await this.guildRepo.invalidateCache(worldId, result.guildId);

        await this.publishGuildChange(EVT.NOTICE_CHANGED, worldId, result.guildId, {});

        return { ok: true, guildId: result.guildId };
    }

    async increaseGuildCapacity(
        worldId: number,
        characterId: number,
        extendedCap: boolean,
        requestId: string,
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

        if (requestId !== "0" && (await this.guildChangeRequestRepo.find(worldId, lockedGuildId, requestId)) != null) {
            const guild = await this.guildRepo.get(worldId, lockedGuildId);
            return { ok: true, guildId: lockedGuildId, capacity: guild?.capacity ?? 0, gp: guild?.gp ?? 0 };
        }

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

            const updatedGuild = await this.guildRepo.set(
                worldId,
                {
                    ...guild,
                    capacity: guild.capacity + GUILD_CAPACITY_STEP,
                    gp: nextGP,
                },
                { txClient },
            );
            if (requestId !== "0") {
                await this.guildChangeRequestRepo.insert(worldId, lockedGuildId, requestId, {}, txClient);
            }
            return {
                ok: true as const,
                guildId: updatedGuild.guildId,
                capacity: updatedGuild.capacity,
                gp: updatedGuild.gp,
            };
        });

        if (!result.ok || result.guildId == null) {
            return result;
        }

        await this.guildRepo.invalidateCache(worldId, result.guildId);

        await this.publishGuildChange(EVT.CAPACITY_CHANGED, worldId, result.guildId, {
            capacity: result.capacity,
            gp: result.gp,
            gp_amount: extendedCap ? -GUILD_CAPACITY_EXTENDED_GP_COST : 0,
        });

        return {
            ok: true,
            guildId: result.guildId,
            capacity: result.capacity,
            gp: result.gp,
        };
    }

    async gainGuildGP(worldId: number, guildId: number, amount: number, requestId: string): Promise<GainGuildGPResult> {
        this.assertWorld(worldId);
        if (!Number.isInteger(guildId) || guildId < 1) {
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_GUILD_NOT_FOUND };
        }
        if (!Number.isInteger(amount) || amount === 0) {
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_UNKNOWN };
        }

        await using _guildLock = await this.distributedLockService.acquireWorldDataLock(worldId, `guild:${guildId}`);

        if (requestId !== "0" && (await this.guildChangeRequestRepo.find(worldId, guildId, requestId)) != null) {
            const guild = await this.guildRepo.get(worldId, guildId);
            return { ok: true, guildId, gp: guild?.gp ?? 0 };
        }

        const result = await this.ctx.withPgDataTransaction(worldId, guildId, async (txClient: PoolClient) => {
            const guild = await this.guildRepo.get(worldId, guildId, { txClient });
            if (!guild) {
                return { ok: false as const, code: messages.GuildErrorCode.GUILD_ERROR_GUILD_NOT_FOUND };
            }
            if (requestId !== "0") {
                await this.guildChangeRequestRepo.insert(worldId, guildId, requestId, {}, txClient);
            }

            const nextGP = Math.min(GUILD_GP_MAX, Math.max(0, guild.gp + amount));
            if (nextGP === guild.gp) {
                return { ok: true as const, guildId, gp: guild.gp, changed: 0 };
            }

            // Only earned GP counts toward the weekly ranking; spending GP leaves the week's total alone.
            const week = guildWeek(Date.now());
            const weeklyGp = guild.gpWeek === week ? guild.weeklyGp : 0;
            const updatedGuild = await this.guildRepo.set(
                worldId,
                {
                    ...guild,
                    gp: nextGP,
                    weeklyGp: Math.min(GUILD_GP_MAX, weeklyGp + Math.max(0, nextGP - guild.gp)),
                    gpWeek: week,
                },
                { txClient },
            );
            return {
                ok: true as const,
                guildId,
                gp: updatedGuild.gp,
                changed: nextGP - guild.gp,
            };
        });

        if (!result.ok) {
            return result;
        }
        if (result.changed === 0) {
            return { ok: true, guildId: result.guildId, gp: result.gp };
        }

        await this.guildRepo.invalidateCache(worldId, guildId);

        await this.publishGuildChange(EVT.GP_CHANGED, worldId, guildId, {
            gp: result.gp,
            amount: result.changed,
        });

        return { ok: true, guildId: result.guildId, gp: result.gp };
    }

    async sendGuildMessage(worldId: number, guildId: number, messageType: number, message: string): Promise<boolean> {
        this.assertWorld(worldId);
        this.assertGuildId(guildId);
        const trimmed = (message ?? "").trim();
        if (trimmed.length <= 0 || trimmed.length > 500) {
            return false;
        }
        const guild = await this.guildRepo.get(worldId, guildId);
        if (!guild) {
            return false;
        }

        await this.publishToGuildRoutes(EVT.MESSAGE, worldId, guildId, {
            message_type: messageType,
            message: trimmed,
        });
        return true;
    }

    async getGuildRanking(worldId: number): Promise<GuildRankingEntryResult[]> {
        this.assertWorld(worldId);
        const { client } = this.ctx.getRedisGlobalAccess(worldId);
        const key = redisCacheKey(`w${worldId}:guild-ranking`);
        const cached = await client.get(key);
        if (cached != null) {
            return JSON.parse(cached) as GuildRankingEntryResult[];
        }

        const guilds = await this.guildRepo.getTopByGP(worldId, GUILD_RANKING_SIZE);
        const entries = guilds.map((guild) => ({
            guildId: guild.guildId,
            name: guild.name,
            gp: guild.gp,
            logo: guild.logo,
        }));
        await client.set(key, JSON.stringify(entries), "EX", GUILD_RANKING_TTL_SEC);
        return entries;
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

        const guildResult = await this.ctx.withPgDataTransaction(
            worldId,
            lockedGuildId,
            async (txClient: PoolClient) => {
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
                await this.guildRepo.set(worldId, { ...guild, disbandedAt: new Date() }, { txClient });
                await this.guildMemberRepo.deleteAllForGuild(worldId, lockedGuildId, { txClient });

                return {
                    ok: true as const,
                    guildId: lockedGuildId,
                    memberCharacterIds,
                    allianceId: guild.allianceId,
                };
            },
        );

        if (!guildResult.ok || guildResult.guildId == null) {
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
                    { txClient },
                );
            });
        }
        const result = guildResult;

        await this.guildBulletinBoardRepo.deleteAllForGuild(worldId, result.guildId);

        await this.guildRepo.invalidateCache(worldId, result.guildId);
        await this.guildMemberRepo.invalidateCache(worldId, String(result.guildId));
        for (const memberCharacterId of result.memberCharacterIds ?? []) {
            await this.characterRealtimeStateRepo.invalidateCache(worldId, memberCharacterId);
        }
        await this.unifiedRepo.deleteGuildName(result.guildId).catch(() => {});

        await this.publishToGuildRoutes(EVT.DISBANDED, worldId, result.guildId, {
            requester_character_id: characterId,
            member_character_ids: result.memberCharacterIds ?? [],
        });

        return { ok: true, guildId: result.guildId, allianceId: result.allianceId };
    }

    async publishMemberOnline(worldId: number, characterId: number, online: boolean) {
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

        const channelIndex = online ? ((await this.sessionRepo.findChannel(worldId, characterId)) ?? -2) : -2;
        await this.publishToGuildRoutes(EVT.MEMBER_ONLINE_CHANGED, worldId, guildId, {
            character_id: characterId,
            online,
            channel_index: channelIndex,
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

    private sortGuildMemberModels(
        memberModels: GuildMemberModel[] | Map<string, GuildMemberModel>,
        leaderCharacterId: number,
    ) {
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

    async guildToPb(
        worldId: number,
        guild: GuildModel,
        memberModels: GuildMemberModel[] | Map<string, GuildMemberModel>,
    ): Promise<GuildMessage> {
        const list = this.sortGuildMemberModels(memberModels, guild.leaderCharacterId);
        const members = await Promise.all(
            list.map(async (m): Promise<GuildMember> => {
                const channelIndex = (await this.sessionRepo.findChannel(worldId, m.characterId)) ?? -2;
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
            }),
        );
        const guildPb: GuildMessage = {
            worldId,
            guildId: guild.guildId,
            name: guild.name,
            leaderCharacterId: guild.leaderCharacterId,
            revision: guild.revision ?? 0,
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

    async broadcastGuildMultiChat(
        worldId: number,
        guildId: number,
        senderCharacterId: number,
        senderName: string,
        message: string,
    ): Promise<{ ok: boolean; recipientCount?: number }> {
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
        await this.publishToGuildRoutes("multi_chat", worldId, guildId, {
            guild_id: guildId,
            sender_character_id: senderCharacterId,
            chat_mode: 2,
            sender_name: name,
            message: trimmedMsg,
        });
        return { ok: true, recipientCount: 1 };
    }
}
