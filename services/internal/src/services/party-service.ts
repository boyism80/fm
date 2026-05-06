import { redisCacheKey } from "../redis-cache-key";
import { Party, type Party as PartyMessage, type PartyDoor, type PartyMember, PartyErrorCode } from "../protobuf/generated/fminternal/internal_service";
import type { PoolClient } from "pg";
import type { PartyModel } from "../repos/party-repository";
import type { PartyMemberModel } from "../repos/party-member-repository";
import type { CharacterSession } from "../repos/session-repository";
import { AppConfiguration } from "../config/app-configuration";
import { InternalContext } from "../context/internal-context";
import { PartyRepository } from "../repos/party-repository";
import { PartyMemberRepository } from "../repos/party-member-repository";
import { CharacterRealtimeStateRepository } from "../repos/character-realtime-state-repository";
import { CharacterRepository } from "../repos/character-repository";
import { UnifiedRepository } from "../repos/unified-repository";
import { SessionRepository } from "../repos/session-repository";
import { RabbitMQService } from "./rabbitmq-service";
const messages = { PartyErrorCode };

const MAX_PARTY_MEMBERS = 6;
const PARTY_STATE_ACTIVE = "ACTIVE";
const EVT = {
    CREATED: "created",
    MEMBER_JOINED: "member_joined",
    MEMBER_LEFT: "member_left",
    DISBANDED: "disbanded",
    LEADER_CHANGED: "leader_changed",
    PARTY_SYNC: "party_sync",
    LOG_ONOFF: "log_onoff",
};

const INVITE_PENDING_TTL_SEC = 300;
const AMQ_DIRECT_EXCHANGE = "amq.direct";

export type PartyMutationResult = { ok: boolean; code?: number; partyId?: number; revision?: number };
export type InvitePartyResult = PartyMutationResult & { targetCharacterId?: number; targetChannelId?: number };
export type DenyPartyResult = { ok: boolean; code?: number };
export type LeavePartyResult = PartyMutationResult & {
    disbanded?: boolean;
    leaderChanged?: boolean;
    oldLeaderCharacterId?: number;
    newLeaderCharacterId?: number;
    realtimeStateCharacterIds?: number[];
};
export type ExpelPartyResult = PartyMutationResult & { disbanded?: boolean; realtimeStateCharacterIds?: number[] };
export type BroadcastMultiChatResult = { ok: boolean; code?: number; deliveredCount?: number };
export type GetPartyResult = { found: boolean; party?: PartyModel; members?: PartyMemberModel[] };

export class PartyService {
    private readonly ctx: InternalContext;
    private readonly app: AppConfiguration;
    private readonly partyRepo: PartyRepository;
    private readonly partyMemberRepo: PartyMemberRepository;
    private readonly characterRealtimeStateRepo: CharacterRealtimeStateRepository;
    private readonly characterRepo: CharacterRepository;
    private readonly unifiedRepo: UnifiedRepository;
    private readonly sessionRepo: SessionRepository;
    private readonly rabbitmqService: RabbitMQService;
    private mqDirectReady: boolean;

    constructor(
        internalContext: InternalContext,
        appConfiguration: AppConfiguration,
        partyRepository: PartyRepository,
        partyMemberRepository: PartyMemberRepository,
        characterRealtimeStateRepository: CharacterRealtimeStateRepository,
        characterRepository: CharacterRepository,
        unifiedRepository: UnifiedRepository,
        sessionRepository: SessionRepository,
        rabbitmqService: RabbitMQService
    ) {
        this.ctx = internalContext;
        this.app = appConfiguration;
        this.partyRepo = partyRepository;
        this.partyMemberRepo = partyMemberRepository;
        this.characterRealtimeStateRepo = characterRealtimeStateRepository;
        this.characterRepo = characterRepository;
        this.unifiedRepo = unifiedRepository;
        this.sessionRepo = sessionRepository;
        this.rabbitmqService = rabbitmqService;
        this.mqDirectReady = false;
    }

    private async ensurePartyMqExchange() {
        if (this.mqDirectReady) {
            return;
        }
        await this.rabbitmqService.assertDirectExchange(AMQ_DIRECT_EXCHANGE);
        this.mqDirectReady = true;
    }

    private async publishToPartyRoutes(
        eventType: string,
        worldId: number,
        partyId: number,
        revision: number,
        extraPayload: Record<string, unknown> = {}
    ) {
        await this.ensurePartyMqExchange();
        const routingKey = `fm.${worldId}.all.party`;
        return this.rabbitmqService.publish(AMQ_DIRECT_EXCHANGE, routingKey, eventType, {
            event_id: extraPayload.event_id ?? `${Date.now()}-${Math.random().toString(36).slice(2, 8)}`,
            party_id: partyId,
            revision,
            occurred_at: new Date().toISOString(),
            ...extraPayload,
        });
    }

    private async publishToPartyGameChannel(worldId: number, channelId: number, eventType: string, payload: Record<string, unknown> = {}) {
        await this.ensurePartyMqExchange();
        const routingKey = `fm.${worldId}.${channelId}.party`;
        return this.rabbitmqService.publish(AMQ_DIRECT_EXCHANGE, routingKey, eventType, {
            occurred_at: new Date().toISOString(),
            ...payload,
        });
    }

    private invitePendingKey(worldId: number, characterId: number) {
        return redisCacheKey(`w${worldId}:party:invite_pending:${characterId}`);
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

    private assertPartyId(partyId: number) {
        const n = partyId;
        if (!Number.isInteger(n) || n < 0 || n > Number.MAX_SAFE_INTEGER) {
            const err = new Error("party_id must be a non-negative integer") as Error & { code?: string };
            err.code = "INVALID_PAYLOAD";
            throw err;
        }
    }

    private assertName(name: string) {
        if (typeof name !== "string" || name.length <= 0 || name.length > 32) {
            const err = new Error("character_name must be non-empty and <= 32 chars") as Error & { code?: string };
            err.code = "INVALID_PAYLOAD";
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

    private async nextPartyId(txClient: PoolClient, worldId: number) {
        this.assertWorld(worldId);
        const res = await txClient.query("SELECT nextval('party_id_seq') AS id");
        const id = res.rows?.[0]?.id;
        if (id == null) {
            const err = new Error("nextval('party_id_seq') returned no id") as Error & { code?: string };
            err.code = "PARTY_ID_SEQ_ERROR";
            throw err;
        }
        return Number(id);
    }

    private async publishPartyEvent(
        eventType: string,
        worldId: number,
        partyId: number,
        revision: number,
        extraPayload: Record<string, unknown> = {}
    ) {
        await this.publishToPartyRoutes(eventType, worldId, partyId, revision, {
            world_id: worldId,
            ...extraPayload,
        });
    }

    private computePartyUiChannelIndex(sess: CharacterSession | null) {
        const ch = sess?.gameServer?.channelId;
        if (ch == null || !Number.isFinite(ch) || ch < 0) {
            return -2;
        }
        return ch;
    }

    private async sessionChannelIndex(worldId: number, characterId: number) {
        const sess = await this.sessionRepo.getCharacterSession(worldId, characterId);
        return this.computePartyUiChannelIndex(sess);
    }

    private sortPartyMemberModels(memberModels: PartyMemberModel[] | Map<string, PartyMemberModel>, leaderCharacterId: number) {
        const arr = Array.isArray(memberModels) ? [...memberModels] : [...memberModels.values()];
        const leaderId = leaderCharacterId;
        return arr.sort((a, b) => {
            const aLead = a.characterId === leaderId ? 0 : 1;
            const bLead = b.characterId === leaderId ? 0 : 1;
            if (aLead !== bLead) {
                return aLead - bLead;
            }
            const ta = a.joinedAt instanceof Date ? a.joinedAt.getTime() : new Date(a.joinedAt || 0).getTime();
            const tb = b.joinedAt instanceof Date ? b.joinedAt.getTime() : new Date(b.joinedAt || 0).getTime();
            if (ta !== tb) {
                return ta - tb;
            }
            return a.characterId - b.characterId;
        });
    }

    private partyToPb(worldId: number, party: PartyModel, memberModels: PartyMemberModel[] | Map<string, PartyMemberModel>): PartyMessage {
        const list = this.sortPartyMemberModels(memberModels, party.leaderCharacterId);
        return {
            worldId,
            partyId: party.partyId,
            leaderCharacterId: party.leaderCharacterId,
            revision: party.revision,
            state: party.state,
            members: list.map((m): PartyMember => {
                const member: PartyMember = {
                    worldId,
                    characterId: m.characterId,
                    characterName: m.characterName,
                    level: m.level,
                    classId: m.classId,
                    role: m.role,
                    mapId: m.mapId,
                    channelIndex: m.channelIndex,
                    door: undefined,
                };
                const d = m.door;
                if (d && Number.isFinite(d.town) && Number.isFinite(d.target) && Number.isFinite(d.x) && Number.isFinite(d.y)) {
                    const door: PartyDoor = {
                        town: d.town,
                        target: d.target,
                        x: d.x,
                        y: d.y,
                    };
                    member.door = door;
                }
                return member;
            }),
        };
    }

    private doorJsonFromPayload(doorPayload: PartyDoor | null | undefined): PartyDoor | null {
        if (doorPayload == null) {
            return null;
        }
        const town = Number(doorPayload.town);
        const target = Number(doorPayload.target);
        const x = Number(doorPayload.x);
        const y = Number(doorPayload.y);
        if (Number.isFinite(town) && Number.isFinite(target) && Number.isFinite(x) && Number.isFinite(y)) {
            return { town, target, x, y };
        }
        return null;
    }

    async updatePartyMember(member: PartyMember | null | undefined): Promise<PartyMutationResult> {
        if (!member || member.worldId == null) {
            return { ok: false, code: messages.PartyErrorCode.UNKNOWN };
        }
        const worldId = member.worldId;
        this.assertWorld(worldId);
        const characterId = member.characterId;
        this.assertCharacterId(characterId);
        const memberLevel = member.level;
        const memberClassId = member.classId;
        this.assertUInt16(memberLevel, "level");
        this.assertUInt16(memberClassId, "class_id");
        const doorJson = this.doorJsonFromPayload(member.door);

        const result = await this.ctx.withPgGlobalTransaction(worldId, async (txClient: PoolClient) => {
            const state = await this.characterRealtimeStateRepo.get(worldId, characterId, { txClient });
            if (state?.partyId == null) {
                return { ok: false, code: messages.PartyErrorCode.NOT_IN_PARTY };
            }
            const partyId = state.partyId;
            if (!Number.isInteger(partyId) || partyId < 0) {
                return { ok: false, code: messages.PartyErrorCode.NOT_IN_PARTY };
            }
            const party = await this.partyRepo.get(worldId, partyId, { txClient });
            if (!party || party.state !== PARTY_STATE_ACTIVE) {
                return { ok: false, code: messages.PartyErrorCode.PARTY_NOT_FOUND };
            }
            const members = await this.partyMemberRepo.getAll(worldId, String(partyId), { txClient });
            const self = members.get(String(characterId));
            if (!self) {
                return { ok: false, code: messages.PartyErrorCode.NOT_IN_PARTY };
            }
            const name = member.characterName.trim();
            if (!name) {
                return { ok: false, code: messages.PartyErrorCode.CHARACTER_NOT_FOUND };
            }
            await this.partyMemberRepo.set(
                worldId,
                { ...self, characterName: name, level: memberLevel, classId: memberClassId, mapId: member.mapId ?? 0, door: doorJson },
                { txClient }
            );
            const nextRevision = party.revision + 1;
            const updatedParty = await this.partyRepo.set(worldId, { ...party, revision: nextRevision }, { txClient });
            return { ok: true, party: updatedParty, triggerCharacterId: characterId };
        });
        if (!result.ok || !result.party) {
            return result;
        }

        await this.partyRepo.evictCache(worldId, result.party.partyId);
        await this.partyMemberRepo.evictGroupCache(worldId, String(result.party.partyId));
        const membersMap = await this.partyMemberRepo.getAll(worldId, String(result.party.partyId));
        const partyPb = this.partyToPb(worldId, result.party, [...membersMap.values()]);
        const wire = Party.encode(partyPb).finish();
        await this.publishPartyEvent(EVT.PARTY_SYNC, worldId, result.party.partyId, result.party.revision, {
            trigger_character_id: result.triggerCharacterId,
            party_pb: Buffer.from(wire).toString("base64"),
        });
        return { ok: true, partyId: result.party.partyId, revision: result.party.revision };
    }

    async applyMemberChannelIndex(worldId: number, characterId: number, channelIndex: number) {
        this.assertWorld(worldId);
        this.assertCharacterId(characterId);
        const ch = channelIndex;
        if (!Number.isInteger(ch) || ch < -2) {
            return;
        }

        const result = await this.ctx.withPgGlobalTransaction(worldId, async (txClient: PoolClient) => {
            const state = await this.characterRealtimeStateRepo.get(worldId, characterId, { txClient });
            if (state?.partyId == null) {
                return { skip: true };
            }
            const partyId = state.partyId;
            if (!Number.isInteger(partyId) || partyId < 0) {
                return { skip: true };
            }
            const party = await this.partyRepo.get(worldId, partyId, { txClient });
            if (!party || party.state !== PARTY_STATE_ACTIVE) {
                return { skip: true };
            }
            const members = await this.partyMemberRepo.getAll(worldId, String(partyId), { txClient });
            const self = members.get(String(characterId));
            if (!self) {
                return { skip: true };
            }
            if (self.channelIndex === ch) {
                return { skip: true };
            }
            await this.partyMemberRepo.set(worldId, { ...self, channelIndex: ch }, { txClient });
            const nextRevision = party.revision + 1;
            const updatedParty = await this.partyRepo.set(worldId, { ...party, revision: nextRevision }, { txClient });
            return { skip: false, party: updatedParty };
        });
        if (result.skip || !result.party) {
            return;
        }

        await this.partyRepo.evictCache(worldId, result.party.partyId);
        await this.partyMemberRepo.evictGroupCache(worldId, String(result.party.partyId));
        const membersMap = await this.partyMemberRepo.getAll(worldId, String(result.party.partyId));
        const partyPb = this.partyToPb(worldId, result.party, [...membersMap.values()]);
        const wire = Party.encode(partyPb).finish();
        await this.publishPartyEvent(EVT.LOG_ONOFF, worldId, result.party.partyId, result.party.revision, {
            character_id: characterId,
            party_pb: Buffer.from(wire).toString("base64"),
        });
    }

    async createParty(worldId: number, leader: PartyMember | null | undefined): Promise<PartyMutationResult> {
        this.assertWorld(worldId);
        if (!leader || leader.worldId !== worldId) {
            return { ok: false, code: messages.PartyErrorCode.UNKNOWN };
        }
        const leaderCharacterId = leader.characterId;
        this.assertCharacterId(leaderCharacterId);
        const name = leader.characterName.trim();
        if (!name) {
            return { ok: false, code: messages.PartyErrorCode.CHARACTER_NOT_FOUND };
        }
        const leaderLevel = leader.level;
        const leaderClassId = leader.classId;
        this.assertUInt16(leaderLevel, "level");
        this.assertUInt16(leaderClassId, "class_id");
        const doorJson = this.doorJsonFromPayload(leader.door);
        let leaderChannelIndex = leader.channelIndex;
        if (leaderChannelIndex == null || !Number.isFinite(leaderChannelIndex) || leaderChannelIndex < -2) {
            leaderChannelIndex = await this.sessionChannelIndex(worldId, leaderCharacterId);
        } else {
            leaderChannelIndex = leaderChannelIndex;
        }

        const result = await this.ctx.withPgGlobalTransaction(worldId, async (txClient: PoolClient) => {
            const state = await this.characterRealtimeStateRepo.get(worldId, leaderCharacterId, { txClient });
            if (state?.partyId != null) {
                return { ok: false, code: messages.PartyErrorCode.ALREADY_IN_PARTY };
            }

            const partyId = await this.nextPartyId(txClient, worldId);
            const party = await this.partyRepo.set(worldId, { worldId, partyId, leaderCharacterId, state: PARTY_STATE_ACTIVE, revision: 1 }, { txClient });
            await this.partyMemberRepo.set(worldId, {
                worldId, partyId, characterId: leaderCharacterId, characterName: name, level: leaderLevel,
                classId: leaderClassId, role: "LEADER", mapId: leader.mapId ?? 0, channelIndex: leaderChannelIndex, door: doorJson,
            }, { txClient });
            await this.characterRealtimeStateRepo.set(worldId, { worldId, characterId: leaderCharacterId, partyId, guildId: state?.guildId ?? null }, { txClient });
            return { ok: true, partyId: party.partyId, revision: party.revision };
        });
        if (!result.ok || result.partyId == null || result.revision == null) {
            return result;
        }
        await this.partyRepo.evictCache(worldId, result.partyId);
        await this.partyMemberRepo.evictGroupCache(worldId, String(result.partyId));
        await this.characterRealtimeStateRepo.evictCache(worldId, leaderCharacterId);
        await this.publishPartyEvent(EVT.CREATED, worldId, result.partyId, result.revision, { leader_character_id: leaderCharacterId });
        return result;
    }

    async inviteParty(worldId: number, inviterCharacterId: number, targetCharacterName: string): Promise<InvitePartyResult> {
        this.assertWorld(worldId);
        this.assertCharacterId(inviterCharacterId);
        const trimmed = targetCharacterName.trim();
        this.assertName(trimmed);

        const inviterState = await this.characterRealtimeStateRepo.get(worldId, inviterCharacterId);
        if (inviterState?.partyId == null) {
            return { ok: false, code: messages.PartyErrorCode.INVITER_NOT_IN_PARTY };
        }
        const partyId = inviterState.partyId;
        if (!Number.isInteger(partyId) || partyId < 0) {
            return { ok: false, code: messages.PartyErrorCode.INVITER_NOT_IN_PARTY };
        }
        const party = await this.partyRepo.get(worldId, partyId);
        if (!party || party.state !== PARTY_STATE_ACTIVE) {
            return { ok: false, code: messages.PartyErrorCode.PARTY_NOT_FOUND };
        }
        const members = await this.partyMemberRepo.getAll(worldId, String(partyId));
        if (!members.has(String(inviterCharacterId))) {
            return { ok: false, code: messages.PartyErrorCode.INVITER_NOT_IN_PARTY };
        }
        if (members.size >= MAX_PARTY_MEMBERS) {
            return { ok: false, code: messages.PartyErrorCode.PARTY_FULL };
        }

        const entry = await this.unifiedRepo.findCharacterNameEntry(trimmed);
        if (!entry || entry.world_id !== worldId) {
            return { ok: false, code: messages.PartyErrorCode.CHARACTER_NOT_FOUND };
        }
        const targetCharacterId = entry.character_id;
        if (targetCharacterId === inviterCharacterId) {
            return { ok: false, code: messages.PartyErrorCode.CANNOT_INVITE_SELF };
        }

        const targetRt = await this.characterRealtimeStateRepo.get(worldId, targetCharacterId);
        if (targetRt?.partyId != null) {
            return { ok: false, code: messages.PartyErrorCode.TARGET_ALREADY_IN_PARTY };
        }

        const sess = await this.sessionRepo.getCharacterSession(worldId, targetCharacterId);
        const ch = sess?.gameServer?.channelId;
        const online = sess?.state === "ONLINE" && sess?.gameServer?.connected === true;
        const chNum = Number(ch);
        if (!online || ch == null || !Number.isFinite(chNum) || chNum < 0) {
            return { ok: false, code: messages.PartyErrorCode.TARGET_OFFLINE };
        }

        const inviterRow = await this.characterRepo.get(worldId, inviterCharacterId);
        const inviterName = inviterRow?.name ? String(inviterRow.name) : "";
        if (!inviterName) {
            return { ok: false, code: messages.PartyErrorCode.CHARACTER_NOT_FOUND };
        }

        const { client } = this.ctx.getRedisGlobalAccess(worldId);
        const inviteKey = this.invitePendingKey(worldId, targetCharacterId);
        await client.set(inviteKey, String(partyId), "EX", INVITE_PENDING_TTL_SEC);

        await this.publishToPartyGameChannel(worldId, chNum, "party_invite", {
            world_id: worldId, party_id: partyId, target_character_id: targetCharacterId, inviter_name: inviterName, party_search: false,
        });
        return { ok: true, partyId, targetCharacterId, targetChannelId: chNum };
    }

    async denyParty(worldId: number, deniedCharacterId: number, inviterCharacterName: string, action: number): Promise<DenyPartyResult> {
        this.assertWorld(worldId);
        this.assertCharacterId(deniedCharacterId);
        const inviterName = inviterCharacterName.trim();
        this.assertName(inviterName);

        const { client } = this.ctx.getRedisGlobalAccess(worldId);
        const inviteKey = this.invitePendingKey(worldId, deniedCharacterId);
        const pending = await client.get(inviteKey);
        if (pending == null) {
            return { ok: false, code: messages.PartyErrorCode.INVITE_EXPIRED_OR_INVALID };
        }
        const pendingPartyId = Number(pending);
        if (!Number.isInteger(pendingPartyId) || pendingPartyId < 0) {
            return { ok: false, code: messages.PartyErrorCode.INVITE_EXPIRED_OR_INVALID };
        }

        const inviterEntry = await this.unifiedRepo.findCharacterNameEntry(inviterName);
        if (!inviterEntry || inviterEntry.world_id !== worldId) {
            await client.del(inviteKey);
            return { ok: false, code: messages.PartyErrorCode.CHARACTER_NOT_FOUND };
        }
        const inviterCharacterId = inviterEntry.character_id;
        const inviterSession = await this.sessionRepo.getCharacterSession(worldId, inviterCharacterId);
        const inviterChRaw = inviterSession?.gameServer?.channelId;
        const inviterChannelID = Number(inviterChRaw);
        const inviterOnline = inviterSession?.state === "ONLINE" && inviterSession?.gameServer?.connected === true;
        if (!inviterOnline || inviterChRaw == null || !Number.isFinite(inviterChannelID) || inviterChannelID < 0) {
            await client.del(inviteKey);
            return { ok: false, code: messages.PartyErrorCode.TARGET_OFFLINE };
        }

        const deniedRow = await this.characterRepo.get(worldId, deniedCharacterId);
        const deniedName = deniedRow?.name ? String(deniedRow.name) : "";
        if (!deniedName) {
            await client.del(inviteKey);
            return { ok: false, code: messages.PartyErrorCode.CHARACTER_NOT_FOUND };
        }

        await client.del(inviteKey);
        await this.publishToPartyGameChannel(worldId, inviterChannelID, "party_invite_denied", {
            world_id: worldId,
            inviter_character_id: inviterCharacterId,
            denied_character_name: deniedName,
            action,
        });
        return { ok: true };
    }

    async joinParty(worldId: number, partyId: number, member: PartyMember | null | undefined, skipInvitePendingCheck = false): Promise<PartyMutationResult> {
        this.assertWorld(worldId);
        this.assertPartyId(partyId);
        if (!member || member.worldId !== worldId) {
            return { ok: false, code: messages.PartyErrorCode.UNKNOWN };
        }
        const characterId = member.characterId;
        this.assertCharacterId(characterId);
        const name = member.characterName.trim();
        this.assertName(name);
        const memberLevel = member.level;
        const memberClassId = member.classId;
        this.assertUInt16(memberLevel, "level");
        this.assertUInt16(memberClassId, "class_id");
        const doorJson = this.doorJsonFromPayload(member.door);
        let joinChannelIndex = member.channelIndex;
        if (joinChannelIndex == null || !Number.isFinite(joinChannelIndex) || joinChannelIndex < -2) {
            joinChannelIndex = await this.sessionChannelIndex(worldId, characterId);
        } else {
            joinChannelIndex = joinChannelIndex;
        }

        const { client } = this.ctx.getRedisGlobalAccess(worldId);
        const inviteKey = this.invitePendingKey(worldId, characterId);
        if (!skipInvitePendingCheck) {
            const pending = await client.get(inviteKey);
            if (pending == null || Number(pending) !== partyId) {
                return { ok: false, code: messages.PartyErrorCode.INVITE_EXPIRED_OR_INVALID };
            }
        }

        const result = await this.ctx.withPgGlobalTransaction(worldId, async (txClient: PoolClient) => {
            const party = await this.partyRepo.get(worldId, partyId, { txClient });
            if (!party || party.state !== PARTY_STATE_ACTIVE) {
                return { ok: false, code: messages.PartyErrorCode.PARTY_NOT_FOUND };
            }
            const members = await this.partyMemberRepo.getAll(worldId, String(partyId), { txClient });
            if (members.size >= MAX_PARTY_MEMBERS) {
                return { ok: false, code: messages.PartyErrorCode.PARTY_FULL };
            }
            const state = await this.characterRealtimeStateRepo.get(worldId, characterId, { txClient });
            if (state?.partyId != null) {
                return { ok: false, code: messages.PartyErrorCode.ALREADY_IN_PARTY };
            }

            await this.partyMemberRepo.set(worldId, {
                worldId, partyId, characterId, characterName: name, level: memberLevel, classId: memberClassId, role: "MEMBER",
                mapId: member.mapId ?? 0, channelIndex: joinChannelIndex, door: doorJson,
            }, { txClient });
            const nextRevision = party.revision + 1;
            const updatedParty = await this.partyRepo.set(worldId, { ...party, revision: nextRevision }, { txClient });
            await this.characterRealtimeStateRepo.set(worldId, { worldId, characterId, partyId, guildId: state?.guildId ?? null }, { txClient });
            return { ok: true, partyId: updatedParty.partyId, revision: updatedParty.revision };
        });
        if (!result.ok || result.partyId == null || result.revision == null) {
            return result;
        }
        await client.del(inviteKey);
        await this.partyRepo.evictCache(worldId, partyId);
        await this.partyMemberRepo.evictGroupCache(worldId, String(partyId));
        await this.characterRealtimeStateRepo.evictCache(worldId, characterId);
        await this.publishPartyEvent(EVT.MEMBER_JOINED, worldId, result.partyId, result.revision, { character_id: characterId });
        return result;
    }

    async leaveParty(worldId: number, characterId: number): Promise<LeavePartyResult> {
        this.assertWorld(worldId);
        this.assertCharacterId(characterId);

        const result = await this.ctx.withPgGlobalTransaction(worldId, async (txClient: PoolClient) => {
            const state = await this.characterRealtimeStateRepo.get(worldId, characterId, { txClient });
            if (state?.partyId == null) {
                return { ok: false, code: messages.PartyErrorCode.NOT_IN_PARTY };
            }
            const partyId = state.partyId;
            if (!Number.isInteger(partyId) || partyId < 0) {
                return { ok: false, code: messages.PartyErrorCode.NOT_IN_PARTY };
            }
            const party = await this.partyRepo.get(worldId, partyId, { txClient });
            if (!party) {
                await this.characterRealtimeStateRepo.set(worldId, { worldId, characterId, partyId: null, guildId: state.guildId ?? null }, { txClient });
                return { ok: true, partyId, revision: 0, disbanded: false };
            }

            await this.partyMemberRepo.del(worldId, String(partyId), characterId, { txClient });
            await this.characterRealtimeStateRepo.set(worldId, { worldId, characterId, partyId: null, guildId: state.guildId ?? null }, { txClient });

            const remaining = await this.partyMemberRepo.getAll(worldId, String(partyId), { txClient });
            if (remaining.size === 0) {
                await this.partyRepo.delete({ worldId, partyId }, { txClient });
                return { ok: true, partyId, revision: party.revision + 1, disbanded: true, realtimeStateCharacterIds: [characterId] };
            }

            const remainingMembers = [...remaining.values()];
            if (party.leaderCharacterId === characterId && remainingMembers.length === 1) {
                const onlyRemaining = remainingMembers[0];
                if (!onlyRemaining || onlyRemaining.characterId == null) {
                    return { ok: false, code: messages.PartyErrorCode.UNKNOWN };
                }
                const disbandCharacterIds = [characterId, onlyRemaining.characterId];
                for (const memberCharacterId of disbandCharacterIds) {
                    const memberState = await this.characterRealtimeStateRepo.get(worldId, memberCharacterId, { txClient });
                    await this.characterRealtimeStateRepo.set(worldId, { worldId, characterId: memberCharacterId, partyId: null, guildId: memberState?.guildId ?? null }, { txClient });
                }
                await this.partyMemberRepo.del(worldId, String(partyId), onlyRemaining.characterId, { txClient });
                await this.partyRepo.delete({ worldId, partyId }, { txClient });
                return { ok: true, partyId, revision: party.revision + 1, disbanded: true, realtimeStateCharacterIds: disbandCharacterIds };
            }

            const previousLeaderCharacterId = party.leaderCharacterId;
            let nextLeaderCharacterId = previousLeaderCharacterId;
            if (party.leaderCharacterId === characterId) {
                const sorted = remainingMembers.sort((a, b) => b.level - a.level);
                const topMember = sorted[0];
                if (!topMember || topMember.characterId == null) {
                    return { ok: false, code: messages.PartyErrorCode.UNKNOWN };
                }
                nextLeaderCharacterId = topMember.characterId;
                await this.partyMemberRepo.set(worldId, { ...topMember, role: "LEADER" }, { txClient });
            }
            const nextRevision = party.revision + 1;
            const updatedParty = await this.partyRepo.set(worldId, { ...party, leaderCharacterId: nextLeaderCharacterId, revision: nextRevision }, { txClient });
            return {
                ok: true,
                partyId: updatedParty.partyId,
                revision: updatedParty.revision,
                disbanded: false,
                leaderChanged: previousLeaderCharacterId !== nextLeaderCharacterId,
                oldLeaderCharacterId: previousLeaderCharacterId,
                newLeaderCharacterId: nextLeaderCharacterId,
                realtimeStateCharacterIds: [characterId],
            };
        });
        if (!result.ok || result.partyId == null || result.revision == null) {
            return result;
        }
        await this.partyRepo.evictCache(worldId, result.partyId);
        await this.partyMemberRepo.evictGroupCache(worldId, String(result.partyId));
        for (const affectedCharacterId of result.realtimeStateCharacterIds ?? [characterId]) {
            await this.characterRealtimeStateRepo.evictCache(worldId, affectedCharacterId);
        }

        if (result.disbanded) {
            await this.publishPartyEvent(EVT.DISBANDED, worldId, result.partyId, result.revision, { character_id: characterId });
        } else {
            await this.publishPartyEvent(EVT.MEMBER_LEFT, worldId, result.partyId, result.revision, {
                character_id: characterId,
                leader_changed: result.leaderChanged ?? false,
                old_leader_character_id: result.oldLeaderCharacterId ?? 0,
                new_leader_character_id: result.newLeaderCharacterId ?? 0,
            });
        }
        return result;
    }

    async expelParty(worldId: number, requesterCharacterId: number, targetCharacterId: number): Promise<ExpelPartyResult> {
        this.assertWorld(worldId);
        this.assertCharacterId(requesterCharacterId);
        this.assertCharacterId(targetCharacterId);
        if (requesterCharacterId === targetCharacterId) {
            return { ok: false, code: messages.PartyErrorCode.CANNOT_EXPEL_SELF };
        }

        const result = await this.ctx.withPgGlobalTransaction(worldId, async (txClient: PoolClient) => {
            const requesterState = await this.characterRealtimeStateRepo.get(worldId, requesterCharacterId, { txClient });
            if (requesterState?.partyId == null) {
                return { ok: false, code: messages.PartyErrorCode.NOT_IN_PARTY };
            }
            const partyId = requesterState.partyId;
            if (!Number.isInteger(partyId) || partyId < 0) {
                return { ok: false, code: messages.PartyErrorCode.NOT_IN_PARTY };
            }
            const party = await this.partyRepo.get(worldId, partyId, { txClient });
            if (!party || party.state !== PARTY_STATE_ACTIVE) {
                return { ok: false, code: messages.PartyErrorCode.PARTY_NOT_FOUND };
            }
            if (party.leaderCharacterId !== requesterCharacterId) {
                return { ok: false, code: messages.PartyErrorCode.NOT_PARTY_LEADER };
            }

            const members = await this.partyMemberRepo.getAll(worldId, String(partyId), { txClient });
            const targetMember = members.get(String(targetCharacterId));
            if (!targetMember) {
                return { ok: false, code: messages.PartyErrorCode.TARGET_NOT_IN_PARTY };
            }

            const targetState = await this.characterRealtimeStateRepo.get(worldId, targetCharacterId, { txClient });
            await this.partyMemberRepo.del(worldId, String(partyId), targetCharacterId, { txClient });
            await this.characterRealtimeStateRepo.set(worldId, { worldId, characterId: targetCharacterId, partyId: null, guildId: targetState?.guildId ?? null }, { txClient });

            const remaining = await this.partyMemberRepo.getAll(worldId, String(partyId), { txClient });
            if (remaining.size === 0) {
                await this.partyRepo.delete({ worldId, partyId }, { txClient });
                return { ok: true, partyId, revision: party.revision + 1, disbanded: true, realtimeStateCharacterIds: [targetCharacterId] };
            }

            const nextRevision = party.revision + 1;
            const updatedParty = await this.partyRepo.set(worldId, { ...party, revision: nextRevision }, { txClient });
            return { ok: true, partyId: updatedParty.partyId, revision: updatedParty.revision, disbanded: false, realtimeStateCharacterIds: [targetCharacterId] };
        });
        if (!result.ok || result.partyId == null || result.revision == null) {
            return result;
        }

        await this.partyRepo.evictCache(worldId, result.partyId);
        await this.partyMemberRepo.evictGroupCache(worldId, String(result.partyId));
        for (const affectedCharacterId of result.realtimeStateCharacterIds ?? [targetCharacterId]) {
            await this.characterRealtimeStateRepo.evictCache(worldId, affectedCharacterId);
        }

        if (result.disbanded) {
            await this.publishPartyEvent(EVT.DISBANDED, worldId, result.partyId, result.revision, {
                character_id: targetCharacterId,
                expelled_by_character_id: requesterCharacterId,
            });
        } else {
            await this.publishPartyEvent(EVT.MEMBER_LEFT, worldId, result.partyId, result.revision, {
                character_id: targetCharacterId,
                expelled_by_character_id: requesterCharacterId,
            });
        }
        return result;
    }

    async changePartyLeader(worldId: number, partyId: number, requesterCharacterId: number, newLeaderCharacterId: number): Promise<PartyMutationResult> {
        this.assertWorld(worldId);
        this.assertPartyId(partyId);
        this.assertCharacterId(requesterCharacterId);
        this.assertCharacterId(newLeaderCharacterId);

        const newLeaderSession = await this.sessionRepo.getCharacterSession(worldId, newLeaderCharacterId);
        const newLeaderOnline = newLeaderSession?.state === "ONLINE" && newLeaderSession?.gameServer?.connected === true;
        if (!newLeaderOnline) {
            return { ok: false, code: messages.PartyErrorCode.TARGET_OFFLINE };
        }

        const result = await this.ctx.withPgGlobalTransaction(worldId, async (txClient: PoolClient) => {
            const party = await this.partyRepo.get(worldId, partyId, { txClient });
            if (!party || party.state !== PARTY_STATE_ACTIVE) {
                return { ok: false, code: messages.PartyErrorCode.PARTY_NOT_FOUND };
            }
            if (party.leaderCharacterId !== requesterCharacterId) {
                return { ok: false, code: messages.PartyErrorCode.NOT_PARTY_LEADER };
            }
            if (requesterCharacterId === newLeaderCharacterId) {
                return { ok: false, code: messages.PartyErrorCode.TARGET_ALREADY_LEADER };
            }

            const members = await this.partyMemberRepo.getAll(worldId, String(partyId), { txClient });
            const newLeader = members.get(String(newLeaderCharacterId));
            if (!newLeader) {
                return { ok: false, code: messages.PartyErrorCode.TARGET_NOT_IN_PARTY };
            }

            await this.partyMemberRepo.set(worldId, { ...newLeader, role: "LEADER" }, { txClient });
            const requester = members.get(String(requesterCharacterId));
            if (requester) await this.partyMemberRepo.set(worldId, { ...requester, role: "MEMBER" }, { txClient });
            const nextRevision = party.revision + 1;
            const updatedParty = await this.partyRepo.set(worldId, { ...party, leaderCharacterId: newLeaderCharacterId, revision: nextRevision }, { txClient });
            return { ok: true, partyId: updatedParty.partyId, revision: updatedParty.revision };
        });
        if (!result.ok || result.partyId == null || result.revision == null) {
            return result;
        }
        await this.partyRepo.evictCache(worldId, partyId);
        await this.partyMemberRepo.evictGroupCache(worldId, String(partyId));
        await this.publishPartyEvent(EVT.LEADER_CHANGED, worldId, result.partyId, result.revision, {
            old_leader_character_id: requesterCharacterId,
            new_leader_character_id: newLeaderCharacterId,
        });
        return result;
    }

    async broadcastMultiChat(worldId: number, memberId: number, senderCharacterId: number, chatMode: number, senderName: string, message: string): Promise<BroadcastMultiChatResult> {
        this.assertWorld(worldId);
        this.assertPartyId(memberId);
        this.assertCharacterId(senderCharacterId);
        const trimmedMsg = message;
        if (trimmedMsg.length <= 0 || trimmedMsg.length > 500) {
            return { ok: false, code: messages.PartyErrorCode.UNKNOWN };
        }
        const name = senderName.trim();
        this.assertName(name);
        const mode = chatMode;
        if (!Number.isInteger(mode) || mode < 0 || mode > 255) {
            return { ok: false, code: messages.PartyErrorCode.UNKNOWN };
        }
        await this.publishToPartyRoutes("multi_chat", worldId, memberId, 0, {
            world_id: worldId,
            member_id: memberId,
            sender_character_id: senderCharacterId,
            chat_mode: mode,
            sender_name: name,
            message: trimmedMsg,
        });
        return { ok: true, deliveredCount: 1 };
    }

    async getParty(worldId: number, partyId: number): Promise<GetPartyResult> {
        this.assertWorld(worldId);
        this.assertPartyId(partyId);
        const party = await this.partyRepo.get(worldId, partyId);
        if (!party) {
            return { found: false };
        }
        const members = await this.partyMemberRepo.getAll(worldId, String(partyId));
        return { found: true, party, members: this.sortPartyMemberModels(members, party.leaderCharacterId) };
    }
}
