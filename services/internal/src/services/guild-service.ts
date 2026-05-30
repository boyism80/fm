import {
    GuildErrorCode,
    GuildMemberRank,
    type GuildMember,
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
import {
    DEFAULT_GUILD_LOGO,
    DEFAULT_GUILD_RANK_TITLES,
} from "../types/guild-json";

const messages = { GuildErrorCode };

const DEFAULT_GUILD_CAPACITY = 10;
const MIN_GUILD_NAME_LEN = 3;
const MAX_GUILD_NAME_LEN = 12;

const AMQ_DIRECT_EXCHANGE = "amq.direct";

const EVT = {
    CREATED: "created",
} as const;

export type CreateGuildResult = {
    ok: boolean;
    code?: GuildErrorCode;
    guildId?: number;
    revision?: number;
    guild?: GuildMessage;
};

export type GetGuildResult = { found: boolean; guild?: GuildModel; members?: GuildMemberModel[] };

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

    constructor(
        internalContext: InternalContext,
        appConfiguration: AppConfiguration,
        unifiedRepository: UnifiedRepository,
        guildRepository: GuildRepository,
        guildMemberRepository: GuildMemberRepository,
        characterRealtimeStateRepository: CharacterRealtimeStateRepository,
        characterRepository: CharacterRepository,
        sessionRepository: SessionRepository,
        rabbitmqService: RabbitMQService
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

        const existingName = await this.unifiedRepo.findGuildNameEntry(normalizedName);
        if (existingName) {
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_GUILD_NAME_TAKEN };
        }

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
