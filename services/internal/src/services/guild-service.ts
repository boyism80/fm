import {
    GuildErrorCode,
    GuildMemberRank,
    type GuildMember,
    type GuildLogo as GuildLogoMessage,
    type Guild as GuildMessage,
} from "../protobuf/generated/fminternal/internal_service";
import type { PoolClient } from "pg";
import { AppConfiguration } from "../config/app-configuration";
import { InternalContext } from "../context/internal-context";
import { UnifiedRepository } from "../repos/unified-repository";
import { GuildRepository } from "../repos/guild-repository";
import { GuildMemberRepository } from "../repos/guild-member-repository";
import { CharacterRealtimeStateRepository } from "../repos/character-realtime-state-repository";
import { CharacterRepository } from "../repos/character-repository";
import { SessionRepository } from "../repos/session-repository";
import type { CharacterSession } from "../repos/session-repository";
import type { GuildModel } from "../repos/guild-repository";
import type { GuildMemberModel } from "../repos/guild-member-repository";
import { RabbitMQService } from "./rabbitmq-service";
import { DistributedLockService } from "./distributed-lock-service";
import {
    DEFAULT_GUILD_LOGO,
    DEFAULT_GUILD_RANK_TITLES,
    type GuildLogo,
    type GuildRankTitles,
} from "../types/guild-json";

const messages = { GuildErrorCode };

const DEFAULT_GUILD_CAPACITY = 10;
const MIN_GUILD_NAME_LEN = 3;
const MAX_GUILD_NAME_LEN = 12;
const MAX_GUILD_NOTICE_LEN = 100;

const AMQ_DIRECT_EXCHANGE = "amq.direct";

const EVT = {
    CREATED: "created",
    MEMBER_JOINED: "member_joined",
    MEMBER_LEFT: "member_left",
    RANK_TITLES_CHANGED: "rank_titles_changed",
    MEMBER_RANK_CHANGED: "member_rank_changed",
    EMBLEM_CHANGED: "emblem_changed",
    NOTICE_CHANGED: "notice_changed",
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

export class GuildService {
    private readonly ctx: InternalContext;
    private readonly app: AppConfiguration;
    private readonly unifiedRepo: UnifiedRepository;
    private readonly guildRepo: GuildRepository;
    private readonly guildMemberRepo: GuildMemberRepository;
    private readonly characterRealtimeStateRepo: CharacterRealtimeStateRepository;
    private readonly characterRepo: CharacterRepository;
    private readonly sessionRepo: SessionRepository;
    private readonly rabbitmqService: RabbitMQService;
    private readonly distributedLockService: DistributedLockService;

    constructor(
        internalContext: InternalContext,
        appConfiguration: AppConfiguration,
        unifiedRepository: UnifiedRepository,
        guildRepository: GuildRepository,
        guildMemberRepository: GuildMemberRepository,
        characterRealtimeStateRepository: CharacterRealtimeStateRepository,
        characterRepository: CharacterRepository,
        sessionRepository: SessionRepository,
        rabbitmqService: RabbitMQService,
        distributedLockService: DistributedLockService
    ) {
        this.ctx = internalContext;
        this.app = appConfiguration;
        this.unifiedRepo = unifiedRepository;
        this.guildRepo = guildRepository;
        this.guildMemberRepo = guildMemberRepository;
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

        await using _characterRealtimeLock = await this.distributedLockService.acquireWorldGlobalLock(
            worldId,
            `character_realtime:${leaderCharacterId}`,
        );

        const leaderState = await this.characterRealtimeStateRepo.get(worldId, leaderCharacterId);
        if (leaderState?.guildId != null && leaderState.guildId > 0) {
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

            await this.ctx.withPgGlobalTransaction(worldId, async (globalTx: PoolClient) => {
                const state = await this.characterRealtimeStateRepo.get(worldId, leaderCharacterId, { txClient: globalTx });
                if (state?.guildId != null && state.guildId > 0) {
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
                    { txClient: globalTx }
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

        await using _characterRealtimeLock = await this.distributedLockService.acquireWorldGlobalLock(
            worldId,
            `character_realtime:${characterId}`,
        );

        await using _guildLock = await this.distributedLockService.acquireWorldGlobalLock(
            worldId,
            `guild:${guildId}`,
        );

        const result = await this.ctx.withPgGlobalTransaction(worldId, async (txClient: PoolClient) => {
            const guild = await this.guildRepo.get(worldId, guildId, { txClient });
            if (!guild) {
                return { ok: false as const, code: messages.GuildErrorCode.GUILD_ERROR_GUILD_NOT_FOUND };
            }

            const members = await this.guildMemberRepo.getAll(worldId, String(guildId), { txClient });
            if (members.size >= guild.capacity) {
                return { ok: false as const, code: messages.GuildErrorCode.GUILD_ERROR_GUILD_FULL };
            }

            const state = await this.characterRealtimeStateRepo.get(worldId, characterId, { txClient });
            if (state?.guildId != null && state.guildId > 0) {
                return { ok: false as const, code: messages.GuildErrorCode.GUILD_ERROR_ALREADY_IN_GUILD };
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

            await this.characterRealtimeStateRepo.set(
                worldId,
                {
                    worldId,
                    characterId,
                    partyId: state?.partyId ?? null,
                    guildId,
                    buddyCapacity: state?.buddyCapacity,
                },
                { txClient }
            );

            return { ok: true as const, revision: updatedGuild.revision };
        });

        if (!result.ok) {
            return { ok: false, code: result.code ?? messages.GuildErrorCode.GUILD_ERROR_UNKNOWN };
        }

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

        await using _characterRealtimeLock = await this.distributedLockService.acquireWorldGlobalLock(
            worldId,
            `character_realtime:${characterId}`,
        );

        const lockedState = await this.characterRealtimeStateRepo.get(worldId, characterId);
        if (lockedState?.guildId == null || lockedState.guildId <= 0) {
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_NOT_IN_GUILD };
        }
        const lockedGuildId = lockedState.guildId;
        if (!Number.isInteger(lockedGuildId) || lockedGuildId <= 0) {
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_NOT_IN_GUILD };
        }

        await using _guildLock = await this.distributedLockService.acquireWorldGlobalLock(
            worldId,
            `guild:${lockedGuildId}`,
        );

        const result = await this.ctx.withPgGlobalTransaction(worldId, async (txClient: PoolClient) => {
            const state = await this.characterRealtimeStateRepo.get(worldId, characterId, { txClient });
            if (state?.guildId == null || state.guildId <= 0) {
                return { ok: false as const, code: messages.GuildErrorCode.GUILD_ERROR_NOT_IN_GUILD };
            }
            const guildId = state.guildId;
            if (!Number.isInteger(guildId) || guildId <= 0) {
                return { ok: false as const, code: messages.GuildErrorCode.GUILD_ERROR_NOT_IN_GUILD };
            }

            const guild = await this.guildRepo.get(worldId, guildId, { txClient });
            if (!guild) {
                await this.characterRealtimeStateRepo.set(
                    worldId,
                    { worldId, characterId, partyId: state.partyId ?? null, guildId: null, buddyCapacity: state.buddyCapacity },
                    { txClient }
                );
                return { ok: true as const, guildId, revision: 0 };
            }

            if (guild.leaderCharacterId === characterId) {
                return { ok: false as const, code: messages.GuildErrorCode.GUILD_ERROR_LEADER_CANNOT_LEAVE };
            }

            const members = await this.guildMemberRepo.getAll(worldId, String(guildId), { txClient });
            if (!members.has(String(characterId))) {
                return { ok: false as const, code: messages.GuildErrorCode.GUILD_ERROR_NOT_IN_GUILD };
            }

            await this.guildMemberRepo.del(worldId, String(guildId), characterId, { txClient });
            await this.characterRealtimeStateRepo.set(
                worldId,
                { worldId, characterId, partyId: state.partyId ?? null, guildId: null, buddyCapacity: state.buddyCapacity },
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

        await using _firstCharacterRealtimeLock = await this.distributedLockService.acquireWorldGlobalLock(
            worldId,
            `character_realtime:${firstCharacterId}`,
        );

        await using _secondCharacterRealtimeLock = await this.distributedLockService.acquireWorldGlobalLock(
            worldId,
            `character_realtime:${secondCharacterId}`,
        );

        const requesterLockedState = await this.characterRealtimeStateRepo.get(worldId, requesterCharacterId);
        if (requesterLockedState?.guildId == null || requesterLockedState.guildId <= 0) {
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_NOT_IN_GUILD };
        }
        const lockedGuildId = requesterLockedState.guildId;
        if (!Number.isInteger(lockedGuildId) || lockedGuildId <= 0) {
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_NOT_IN_GUILD };
        }

        await using _guildLock = await this.distributedLockService.acquireWorldGlobalLock(
            worldId,
            `guild:${lockedGuildId}`,
        );

        const result = await this.ctx.withPgGlobalTransaction(worldId, async (txClient: PoolClient) => {
            const requesterState = await this.characterRealtimeStateRepo.get(worldId, requesterCharacterId, { txClient });
            if (requesterState?.guildId == null || requesterState.guildId <= 0) {
                return { ok: false as const, code: messages.GuildErrorCode.GUILD_ERROR_NOT_IN_GUILD };
            }
            const guildId = requesterState.guildId;
            if (!Number.isInteger(guildId) || guildId <= 0) {
                return { ok: false as const, code: messages.GuildErrorCode.GUILD_ERROR_NOT_IN_GUILD };
            }

            const targetState = await this.characterRealtimeStateRepo.get(worldId, targetCharacterId, { txClient });
            if (targetState?.guildId !== guildId) {
                return { ok: false as const, code: messages.GuildErrorCode.GUILD_ERROR_TARGET_NOT_IN_GUILD };
            }

            const guild = await this.guildRepo.get(worldId, guildId, { txClient });
            if (!guild) {
                return { ok: false as const, code: messages.GuildErrorCode.GUILD_ERROR_GUILD_NOT_FOUND };
            }

            const members = await this.guildMemberRepo.getAll(worldId, String(guildId), { txClient });
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

            await this.guildMemberRepo.del(worldId, String(guildId), targetCharacterId, { txClient });
            await this.characterRealtimeStateRepo.set(
                worldId,
                {
                    worldId,
                    characterId: targetCharacterId,
                    partyId: targetState.partyId ?? null,
                    guildId: null,
                    buddyCapacity: targetState.buddyCapacity,
                },
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

        await using _characterRealtimeLock = await this.distributedLockService.acquireWorldGlobalLock(
            worldId,
            `character_realtime:${characterId}`,
        );

        const lockedState = await this.characterRealtimeStateRepo.get(worldId, characterId);
        if (lockedState?.guildId == null || lockedState.guildId <= 0) {
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_NOT_IN_GUILD };
        }
        const lockedGuildId = lockedState.guildId;
        if (!Number.isInteger(lockedGuildId) || lockedGuildId <= 0) {
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_NOT_IN_GUILD };
        }

        await using _guildLock = await this.distributedLockService.acquireWorldGlobalLock(
            worldId,
            `guild:${lockedGuildId}`,
        );

        const result = await this.ctx.withPgGlobalTransaction(worldId, async (txClient: PoolClient) => {
            const state = await this.characterRealtimeStateRepo.get(worldId, characterId, { txClient });
            if (state?.guildId == null || state.guildId <= 0) {
                return { ok: false as const, code: messages.GuildErrorCode.GUILD_ERROR_NOT_IN_GUILD };
            }
            const guildId = state.guildId;
            if (!Number.isInteger(guildId) || guildId <= 0) {
                return { ok: false as const, code: messages.GuildErrorCode.GUILD_ERROR_NOT_IN_GUILD };
            }

            const guild = await this.guildRepo.get(worldId, guildId, { txClient });
            if (!guild) {
                return { ok: false as const, code: messages.GuildErrorCode.GUILD_ERROR_GUILD_NOT_FOUND };
            }
            if (guild.leaderCharacterId !== characterId) {
                return { ok: false as const, code: messages.GuildErrorCode.GUILD_ERROR_NOT_AUTHORIZED };
            }

            const members = await this.guildMemberRepo.getAll(worldId, String(guildId), { txClient });
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

        await using _firstCharacterRealtimeLock = await this.distributedLockService.acquireWorldGlobalLock(
            worldId,
            `character_realtime:${firstCharacterId}`,
        );

        await using _secondCharacterRealtimeLock = await this.distributedLockService.acquireWorldGlobalLock(
            worldId,
            `character_realtime:${secondCharacterId}`,
        );

        const requesterLockedState = await this.characterRealtimeStateRepo.get(worldId, requesterCharacterId);
        if (requesterLockedState?.guildId == null || requesterLockedState.guildId <= 0) {
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_NOT_IN_GUILD };
        }
        const lockedGuildId = requesterLockedState.guildId;
        if (!Number.isInteger(lockedGuildId) || lockedGuildId <= 0) {
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_NOT_IN_GUILD };
        }

        await using _guildLock = await this.distributedLockService.acquireWorldGlobalLock(
            worldId,
            `guild:${lockedGuildId}`,
        );

        const result = await this.ctx.withPgGlobalTransaction(worldId, async (txClient: PoolClient) => {
            const requesterState = await this.characterRealtimeStateRepo.get(worldId, requesterCharacterId, { txClient });
            if (requesterState?.guildId == null || requesterState.guildId <= 0) {
                return { ok: false as const, code: messages.GuildErrorCode.GUILD_ERROR_NOT_IN_GUILD };
            }
            const guildId = requesterState.guildId;
            if (!Number.isInteger(guildId) || guildId <= 0) {
                return { ok: false as const, code: messages.GuildErrorCode.GUILD_ERROR_NOT_IN_GUILD };
            }

            const targetState = await this.characterRealtimeStateRepo.get(worldId, targetCharacterId, { txClient });
            if (targetState?.guildId !== guildId) {
                return { ok: false as const, code: messages.GuildErrorCode.GUILD_ERROR_TARGET_NOT_IN_GUILD };
            }

            const guild = await this.guildRepo.get(worldId, guildId, { txClient });
            if (!guild) {
                return { ok: false as const, code: messages.GuildErrorCode.GUILD_ERROR_GUILD_NOT_FOUND };
            }

            const members = await this.guildMemberRepo.getAll(worldId, String(guildId), { txClient });
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

        await using _characterRealtimeLock = await this.distributedLockService.acquireWorldGlobalLock(
            worldId,
            `character_realtime:${characterId}`,
        );

        const lockedState = await this.characterRealtimeStateRepo.get(worldId, characterId);
        if (lockedState?.guildId == null || lockedState.guildId <= 0) {
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_NOT_IN_GUILD };
        }
        const lockedGuildId = lockedState.guildId;
        if (!Number.isInteger(lockedGuildId) || lockedGuildId <= 0) {
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_NOT_IN_GUILD };
        }

        await using _guildLock = await this.distributedLockService.acquireWorldGlobalLock(
            worldId,
            `guild:${lockedGuildId}`,
        );

        const result = await this.ctx.withPgGlobalTransaction(worldId, async (txClient: PoolClient) => {
            const state = await this.characterRealtimeStateRepo.get(worldId, characterId, { txClient });
            if (state?.guildId == null || state.guildId <= 0) {
                return { ok: false as const, code: messages.GuildErrorCode.GUILD_ERROR_NOT_IN_GUILD };
            }
            const guildId = state.guildId;
            if (!Number.isInteger(guildId) || guildId <= 0) {
                return { ok: false as const, code: messages.GuildErrorCode.GUILD_ERROR_NOT_IN_GUILD };
            }

            const guild = await this.guildRepo.get(worldId, guildId, { txClient });
            if (!guild) {
                return { ok: false as const, code: messages.GuildErrorCode.GUILD_ERROR_GUILD_NOT_FOUND };
            }
            if (guild.leaderCharacterId !== characterId) {
                return { ok: false as const, code: messages.GuildErrorCode.GUILD_ERROR_NOT_AUTHORIZED };
            }

            const members = await this.guildMemberRepo.getAll(worldId, String(guildId), { txClient });
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

        await using _characterRealtimeLock = await this.distributedLockService.acquireWorldGlobalLock(
            worldId,
            `character_realtime:${characterId}`,
        );

        const lockedState = await this.characterRealtimeStateRepo.get(worldId, characterId);
        if (lockedState?.guildId == null || lockedState.guildId <= 0) {
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_NOT_IN_GUILD };
        }
        const lockedGuildId = lockedState.guildId;
        if (!Number.isInteger(lockedGuildId) || lockedGuildId <= 0) {
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_NOT_IN_GUILD };
        }

        await using _guildLock = await this.distributedLockService.acquireWorldGlobalLock(
            worldId,
            `guild:${lockedGuildId}`,
        );

        const result = await this.ctx.withPgGlobalTransaction(worldId, async (txClient: PoolClient) => {
            const state = await this.characterRealtimeStateRepo.get(worldId, characterId, { txClient });
            if (state?.guildId == null || state.guildId <= 0) {
                return { ok: false as const, code: messages.GuildErrorCode.GUILD_ERROR_NOT_IN_GUILD };
            }
            const guildId = state.guildId;
            if (!Number.isInteger(guildId) || guildId <= 0) {
                return { ok: false as const, code: messages.GuildErrorCode.GUILD_ERROR_NOT_IN_GUILD };
            }

            const guild = await this.guildRepo.get(worldId, guildId, { txClient });
            if (!guild) {
                return { ok: false as const, code: messages.GuildErrorCode.GUILD_ERROR_GUILD_NOT_FOUND };
            }

            const members = await this.guildMemberRepo.getAll(worldId, String(guildId), { txClient });
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

    async disbandGuild(worldId: number, characterId: number): Promise<DisbandGuildResult> {
        this.assertWorld(worldId);
        this.assertCharacterId(characterId);

        await using _characterRealtimeLock = await this.distributedLockService.acquireWorldGlobalLock(
            worldId,
            `character_realtime:${characterId}`,
        );

        const lockedState = await this.characterRealtimeStateRepo.get(worldId, characterId);
        if (lockedState?.guildId == null || lockedState.guildId <= 0) {
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_NOT_IN_GUILD };
        }
        const lockedGuildId = lockedState.guildId;
        if (!Number.isInteger(lockedGuildId) || lockedGuildId <= 0) {
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_NOT_IN_GUILD };
        }

        await using _guildLock = await this.distributedLockService.acquireWorldGlobalLock(
            worldId,
            `guild:${lockedGuildId}`,
        );

        const membersMap = await this.guildMemberRepo.getAll(worldId, String(lockedGuildId));
        const otherMemberKeyRests = [...membersMap.values()]
            .map((member) => member.characterId)
            .filter((memberCharacterId) => memberCharacterId !== characterId)
            .sort((a, b) => a - b)
            .map((memberCharacterId) => `character_realtime:${memberCharacterId}`);

        await using _otherMemberLocks = await this.distributedLockService.acquireWorldGlobalLocks(
            worldId,
            otherMemberKeyRests,
        );

        const result = await this.ctx.withPgGlobalTransaction(worldId, async (txClient: PoolClient) => {
            const state = await this.characterRealtimeStateRepo.get(worldId, characterId, { txClient });
            if (state?.guildId == null || state.guildId <= 0) {
                return { ok: false as const, code: messages.GuildErrorCode.GUILD_ERROR_NOT_IN_GUILD };
            }
            const guildId = state.guildId;
            if (!Number.isInteger(guildId) || guildId <= 0) {
                return { ok: false as const, code: messages.GuildErrorCode.GUILD_ERROR_NOT_IN_GUILD };
            }

            const guild = await this.guildRepo.get(worldId, guildId, { txClient });
            if (!guild) {
                return { ok: false as const, code: messages.GuildErrorCode.GUILD_ERROR_GUILD_NOT_FOUND };
            }

            const members = await this.guildMemberRepo.getAll(worldId, String(guildId), { txClient });
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
            await this.guildMemberRepo.deleteAllForGuild(worldId, guildId, { txClient });

            for (const memberCharacterId of memberCharacterIds) {
                const memberState = await this.characterRealtimeStateRepo.get(worldId, memberCharacterId, { txClient });
                if (!memberState) {
                    continue;
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
            }

            return {
                ok: true as const,
                guildId,
                revision: nextRevision,
                memberCharacterIds,
            };
        });

        if (!result.ok || result.guildId == null || result.revision == null) {
            return result;
        }

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
        if (state?.guildId == null || state.guildId <= 0) {
            return;
        }
        const guildId = state.guildId;
        if (!Number.isInteger(guildId) || guildId <= 0) {
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
        if (!Number.isInteger(n) || n <= 0 || n > 0xffffffff) {
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
                return {
                    worldId,
                    characterId: m.characterId,
                    characterName: m.characterName,
                    level: m.level,
                    classId: m.classId,
                    rank: m.guildRank,
                    channelIndex,
                };
            })
        );
        return {
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
}
