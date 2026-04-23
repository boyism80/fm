"use strict";

const { redisCacheKey } = require("../redis-cache-key");
const messages = require("../../protobuf/fminternal/internal_service_pb");

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

class PartyService {
    constructor(
        internalContext,
        appConfiguration,
        partyRepository,
        partyMemberRepository,
        characterRealtimeStateRepository,
        characterRepository,
        unifiedRepository,
        sessionRepository,
        rabbitmqService
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
        this._mqDirectReady = false;
    }

    async _ensurePartyMqExchange() {
        if (this._mqDirectReady) {
            return;
        }
        await this.rabbitmqService.assertDirectExchange(AMQ_DIRECT_EXCHANGE);
        this._mqDirectReady = true;
    }

    async _publishToPartyRoutes(eventType, worldId, partyId, revision, extraPayload = {}) {
        await this._ensurePartyMqExchange();
        const routingKey = `fm.${worldId}.all.party`;
        return this.rabbitmqService.publish(AMQ_DIRECT_EXCHANGE, routingKey, eventType, {
            event_id: extraPayload.event_id ?? `${Date.now()}-${Math.random().toString(36).slice(2, 8)}`,
            party_id: partyId,
            revision: Number(revision),
            occurred_at: new Date().toISOString(),
            ...extraPayload,
        });
    }

    async _publishToPartyGameChannel(worldId, channelId, eventType, payload = {}) {
        await this._ensurePartyMqExchange();
        const routingKey = `fm.${worldId}.${channelId}.party`;
        return this.rabbitmqService.publish(AMQ_DIRECT_EXCHANGE, routingKey, eventType, {
            occurred_at: new Date().toISOString(),
            ...payload,
        });
    }

    _invitePendingKey(worldId, characterId) {
        return redisCacheKey(`w${worldId}:party:invite_pending:${characterId}`);
    }

    _assertWorld(worldId) {
        const wid = String(worldId);
        if (!this.app.postgresql.worlds[wid]) {
            const err = new Error(`Unknown world_id: ${worldId}`);
            err.code = "UNKNOWN_WORLD";
            throw err;
        }
    }

    _assertCharacterId(characterId) {
        const n = Number(characterId);
        if (!Number.isInteger(n) || n <= 0 || n > 0xffffffff) {
            const err = new Error("character_id must be a positive uint32");
            err.code = "INVALID_CHARACTER_ID";
            throw err;
        }
    }

    _assertPartyId(partyId) {
        const n = Number(partyId);
        if (!Number.isInteger(n) || n < 0 || n > Number.MAX_SAFE_INTEGER) {
            const err = new Error("party_id must be a non-negative integer");
            err.code = "INVALID_PAYLOAD";
            throw err;
        }
    }

    _assertName(name) {
        if (typeof name !== "string" || name.length <= 0 || name.length > 32) {
            const err = new Error("character_name must be non-empty and <= 32 chars");
            err.code = "INVALID_PAYLOAD";
            throw err;
        }
    }

    _assertUInt16(value, fieldName) {
        const n = Number(value);
        if (!Number.isInteger(n) || n < 0 || n > 0xffff) {
            const err = new Error(`${fieldName} must be uint16`);
            err.code = "INVALID_PAYLOAD";
            throw err;
        }
    }

    
    async _nextPartyId(txClient, worldId) {
        this._assertWorld(worldId);
        const res = await txClient.query("SELECT nextval('party_id_seq') AS id");
        const id = res.rows?.[0]?.id;
        if (id == null) {
            const err = new Error("nextval('party_id_seq') returned no id");
            err.code = "PARTY_ID_SEQ_ERROR";
            throw err;
        }
        return Number(id);
    }

    async _publishPartyEvent(eventType, worldId, partyId, revision, extraPayload = {}) {
        await this._publishToPartyRoutes(eventType, worldId, partyId, revision, {
            world_id: Number(worldId),
            ...extraPayload,
        });
    }

    _computePartyUiChannelIndex(sess) {
        const ch = sess?.gameServer?.channelId;
        const n = Number(ch);
        if (ch == null || !Number.isFinite(n) || n < 0) {
            return -2;
        }
        return n;
    }

    async _sessionChannelIndex(worldId, characterId) {
        const sess = await this.sessionRepo.getCharacterSession(worldId, characterId);
        return this._computePartyUiChannelIndex(sess);
    }

    _sortPartyMemberModels(memberModels, leaderCharacterId) {
        const arr = Array.isArray(memberModels) ? [...memberModels] : [...memberModels.values()];
        const leaderId = Number(leaderCharacterId ?? 0);
        return arr.sort((a, b) => {
            const aLead = Number(a.characterId) === leaderId ? 0 : 1;
            const bLead = Number(b.characterId) === leaderId ? 0 : 1;
            if (aLead !== bLead) {
                return aLead - bLead;
            }
            const ta =
                a.joinedAt instanceof Date
                    ? a.joinedAt.getTime()
                    : new Date(a.joinedAt || 0).getTime();
            const tb =
                b.joinedAt instanceof Date
                    ? b.joinedAt.getTime()
                    : new Date(b.joinedAt || 0).getTime();
            if (ta !== tb) {
                return ta - tb;
            }
            return Number(a.characterId) - Number(b.characterId);
        });
    }

    _partyToPb(worldId, party, memberModels) {
        const list = this._sortPartyMemberModels(memberModels, party.leaderCharacterId);
        const p = new messages.Party();
        p.setWorldId(Number(worldId));
        p.setPartyId(party.partyId);
        p.setLeaderCharacterId(party.leaderCharacterId);
        p.setRevision(Number(party.revision));
        p.setState(String(party.state ?? ""));
        p.setMembersList(
            list.map((m) => {
                const mm = new messages.PartyMember();
                mm.setWorldId(Number(worldId));
                mm.setCharacterId(m.characterId);
                mm.setCharacterName(String(m.characterName ?? ""));
                mm.setLevel(Number(m.level ?? 0));
                mm.setClassId(Number(m.classId ?? 0));
                mm.setRole(String(m.role ?? ""));
                mm.setMapId(Number(m.mapId ?? 0));
                mm.setChannelIndex(Number(m.channelIndex ?? -2));
                const d = m.door;
                if (
                    d &&
                    typeof d === "object" &&
                    Number.isFinite(Number(d.town)) &&
                    Number.isFinite(Number(d.target)) &&
                    Number.isFinite(Number(d.x)) &&
                    Number.isFinite(Number(d.y))
                ) {
                    const pd = new messages.PartyDoor();
                    pd.setTown(Number(d.town));
                    pd.setTarget(Number(d.target));
                    pd.setX(Number(d.x));
                    pd.setY(Number(d.y));
                    mm.setDoor(pd);
                } else {
                    mm.clearDoor();
                }
                return mm;
            })
        );
        return p;
    }

    _doorJsonFromPayload(doorPayload) {
        if (doorPayload == null) {
            return null;
        }
        const town = Number(doorPayload.town);
        const target = Number(doorPayload.target);
        const x = Number(doorPayload.x);
        const y = Number(doorPayload.y);
        if (
            Number.isFinite(town) &&
            Number.isFinite(target) &&
            Number.isFinite(x) &&
            Number.isFinite(y)
        ) {
            return { town, target, x, y };
        }
        return null;
    }

    async updatePartyMember(member) {
        if (!member || member.worldId == null) {
            return { ok: false, code: messages.PartyErrorCode.UNKNOWN };
        }
        const worldId = Number(member.worldId);
        this._assertWorld(worldId);
        const characterId = Number(member.characterId);
        this._assertCharacterId(characterId);
        this._assertUInt16(member.level, "level");
        this._assertUInt16(member.classId, "class_id");
        const doorJson = this._doorJsonFromPayload(member.door);

        const result = await this.ctx.withPgGlobalTransaction(worldId, async (txClient) => {
            const state = await this.characterRealtimeStateRepo.get(worldId, characterId, { txClient });
            if (state?.partyId == null) {
                return { ok: false, code: messages.PartyErrorCode.NOT_IN_PARTY };
            }
            const partyId = Number(state.partyId);
            if (!Number.isInteger(partyId) || partyId < 0) {
                return { ok: false, code: messages.PartyErrorCode.NOT_IN_PARTY };
            }
            const party = await this.partyRepo.get(worldId, partyId, { txClient });
            if (!party || party.state !== PARTY_STATE_ACTIVE) {
                return { ok: false, code: messages.PartyErrorCode.PARTY_NOT_FOUND };
            }
            const members = await this.partyMemberRepo.getAll(worldId, partyId, { txClient });
            const self = members.get(String(characterId));
            if (!self) {
                return { ok: false, code: messages.PartyErrorCode.NOT_IN_PARTY };
            }
            const name = String(member.characterName ?? "").trim();
            if (!name) {
                return { ok: false, code: messages.PartyErrorCode.CHARACTER_NOT_FOUND };
            }
            await this.partyMemberRepo.set(
                worldId,
                {
                    ...self,
                    characterName: name,
                    level: Number(member.level),
                    classId: Number(member.classId),
                    mapId: Number(member.mapId ?? 0),
                    door: doorJson,
                },
                { txClient }
            );
            const nextRevision = Number(party.revision) + 1;
            const updatedParty = await this.partyRepo.set(
                worldId,
                {
                    ...party,
                    revision: nextRevision,
                },
                { txClient }
            );
            return { ok: true, party: updatedParty, triggerCharacterId: characterId };
        });
        if (!result.ok) {
            return result;
        }

        await this.partyRepo.evictCache(worldId, result.party.partyId);
        await this.partyMemberRepo.evictGroupCache(worldId, result.party.partyId);

        const membersMap = await this.partyMemberRepo.getAll(worldId, result.party.partyId);
        const partyPb = this._partyToPb(worldId, result.party, [...membersMap.values()]);
        const wire = partyPb.serializeBinary();
        await this._publishPartyEvent(EVT.PARTY_SYNC, worldId, result.party.partyId, result.party.revision, {
            trigger_character_id: Number(result.triggerCharacterId),
            party_pb: Buffer.from(wire).toString("base64"),
        });

        return { ok: true, partyId: result.party.partyId, revision: result.party.revision };
    }

    async applyMemberChannelIndex(worldId, characterId, channelIndex) {
        this._assertWorld(worldId);
        this._assertCharacterId(characterId);
        const ch = Number(channelIndex);
        if (!Number.isInteger(ch) || ch < -2) {
            return;
        }

        const result = await this.ctx.withPgGlobalTransaction(worldId, async (txClient) => {
            const state = await this.characterRealtimeStateRepo.get(worldId, characterId, { txClient });
            if (state?.partyId == null) {
                return { skip: true };
            }
            const partyId = Number(state.partyId);
            if (!Number.isInteger(partyId) || partyId < 0) {
                return { skip: true };
            }
            const party = await this.partyRepo.get(worldId, partyId, { txClient });
            if (!party || party.state !== PARTY_STATE_ACTIVE) {
                return { skip: true };
            }
            const members = await this.partyMemberRepo.getAll(worldId, partyId, { txClient });
            const self = members.get(String(characterId));
            if (!self) {
                return { skip: true };
            }
            if (Number(self.channelIndex) === ch) {
                return { skip: true };
            }
            await this.partyMemberRepo.set(worldId, { ...self, channelIndex: ch }, { txClient });
            const nextRevision = Number(party.revision) + 1;
            const updatedParty = await this.partyRepo.set(
                worldId,
                {
                    ...party,
                    revision: nextRevision,
                },
                { txClient }
            );
            return { skip: false, party: updatedParty };
        });

        if (result.skip) {
            return;
        }

        await this.partyRepo.evictCache(worldId, result.party.partyId);
        await this.partyMemberRepo.evictGroupCache(worldId, result.party.partyId);

        const membersMap = await this.partyMemberRepo.getAll(worldId, result.party.partyId);
        const partyPb = this._partyToPb(worldId, result.party, [...membersMap.values()]);
        const wire = partyPb.serializeBinary();
        await this._publishPartyEvent(EVT.LOG_ONOFF, worldId, result.party.partyId, result.party.revision, {
            character_id: Number(characterId),
            party_pb: Buffer.from(wire).toString("base64"),
        });
    }

    async createParty(worldId, leader) {
        this._assertWorld(worldId);
        if (!leader || Number(leader.worldId) !== Number(worldId)) {
            return { ok: false, code: messages.PartyErrorCode.UNKNOWN };
        }
        const leaderCharacterId = Number(leader.characterId);
        this._assertCharacterId(leaderCharacterId);
        const name = String(leader.characterName ?? "").trim();
        if (!name) {
            return { ok: false, code: messages.PartyErrorCode.CHARACTER_NOT_FOUND };
        }
        this._assertUInt16(leader.level, "level");
        this._assertUInt16(leader.classId, "class_id");
        const doorJson = this._doorJsonFromPayload(leader.door);
        let leaderChannelIndex = leader.channelIndex;
        if (
            leaderChannelIndex == null ||
            !Number.isFinite(Number(leaderChannelIndex)) ||
            Number(leaderChannelIndex) < -2
        ) {
            leaderChannelIndex = await this._sessionChannelIndex(worldId, leaderCharacterId);
        } else {
            leaderChannelIndex = Number(leaderChannelIndex);
        }

        const result = await this.ctx.withPgGlobalTransaction(worldId, async (txClient) => {
            const state = await this.characterRealtimeStateRepo.get(worldId, leaderCharacterId, { txClient });
            if (state?.partyId != null) {
                return { ok: false, code: messages.PartyErrorCode.ALREADY_IN_PARTY };
            }

            const partyId = await this._nextPartyId(txClient, worldId);
            const party = await this.partyRepo.set(worldId, {
                worldId,
                partyId,
                leaderCharacterId,
                state: PARTY_STATE_ACTIVE,
                revision: 1,
            }, { txClient });
            await this.partyMemberRepo.set(worldId, {
                worldId,
                partyId,
                characterId: leaderCharacterId,
                characterName: name,
                level: Number(leader.level),
                classId: Number(leader.classId),
                role: "LEADER",
                mapId: Number(leader.mapId ?? 0),
                channelIndex: leaderChannelIndex,
                door: doorJson,
            }, { txClient });
            await this.characterRealtimeStateRepo.set(worldId, {
                worldId,
                characterId: leaderCharacterId,
                partyId,
                guildId: state?.guildId ?? null,
            }, { txClient });
            return { ok: true, partyId: party.partyId, revision: party.revision };
        });
        if (!result.ok) {
            return result;
        }
        await this.partyRepo.evictCache(worldId, result.partyId);
        await this.partyMemberRepo.evictGroupCache(worldId, result.partyId);
        await this.characterRealtimeStateRepo.evictCache(worldId, leaderCharacterId);

        await this._publishPartyEvent(EVT.CREATED, worldId, result.partyId, result.revision, {
            leader_character_id: Number(leaderCharacterId),
        });

        return result;
    }

    async inviteParty(worldId, inviterCharacterId, targetCharacterName) {
        this._assertWorld(worldId);
        this._assertCharacterId(inviterCharacterId);
        const trimmed = String(targetCharacterName ?? "").trim();
        this._assertName(trimmed);

        const inviterState = await this.characterRealtimeStateRepo.get(worldId, inviterCharacterId);
        if (inviterState?.partyId == null) {
            return { ok: false, code: messages.PartyErrorCode.INVITER_NOT_IN_PARTY };
        }
        const partyId = Number(inviterState.partyId);
        if (!Number.isInteger(partyId) || partyId < 0) {
            return { ok: false, code: messages.PartyErrorCode.INVITER_NOT_IN_PARTY };
        }
        const party = await this.partyRepo.get(worldId, partyId);
        if (!party || party.state !== PARTY_STATE_ACTIVE) {
            return { ok: false, code: messages.PartyErrorCode.PARTY_NOT_FOUND };
        }
        const members = await this.partyMemberRepo.getAll(worldId, partyId);
        if (!members.has(String(inviterCharacterId))) {
            return { ok: false, code: messages.PartyErrorCode.INVITER_NOT_IN_PARTY };
        }
        if (members.size >= MAX_PARTY_MEMBERS) {
            return { ok: false, code: messages.PartyErrorCode.PARTY_FULL };
        }

        const entry = await this.unifiedRepo.findCharacterNameEntry(trimmed);
        if (!entry || Number(entry.world_id) !== Number(worldId)) {
            return { ok: false, code: messages.PartyErrorCode.CHARACTER_NOT_FOUND };
        }
        const targetCharacterId = Number(entry.character_id);
        if (targetCharacterId === Number(inviterCharacterId)) {
            return { ok: false, code: messages.PartyErrorCode.CANNOT_INVITE_SELF };
        }

        const targetRt = await this.characterRealtimeStateRepo.get(worldId, targetCharacterId);
        if (targetRt?.partyId != null) {
            return { ok: false, code: messages.PartyErrorCode.TARGET_ALREADY_IN_PARTY };
        }

        const sess = await this.sessionRepo.getCharacterSession(worldId, targetCharacterId);
        const ch = sess?.gameServer?.channelId;
        const online =
            sess?.state === "ONLINE" && Boolean(sess?.gameServer?.connected);
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
        const inviteKey = this._invitePendingKey(worldId, targetCharacterId);
        await client.set(inviteKey, String(partyId), "EX", INVITE_PENDING_TTL_SEC);

        await this._publishToPartyGameChannel(worldId, Number(ch), "party_invite", {
            world_id: Number(worldId),
            party_id: partyId,
            target_character_id: targetCharacterId,
            inviter_name: inviterName,
            party_search: false,
        });

        return {
            ok: true,
            partyId,
            targetCharacterId,
            targetChannelId: Number(ch),
        };
    }

    async denyParty(worldId, deniedCharacterId, inviterCharacterName, action) {
        this._assertWorld(worldId);
        this._assertCharacterId(deniedCharacterId);
        const inviterName = String(inviterCharacterName ?? "").trim();
        this._assertName(inviterName);

        const { client } = this.ctx.getRedisGlobalAccess(worldId);
        const inviteKey = this._invitePendingKey(worldId, deniedCharacterId);
        const pending = await client.get(inviteKey);
        if (pending == null) {
            return { ok: false, code: messages.PartyErrorCode.INVITE_EXPIRED_OR_INVALID };
        }
        const pendingPartyId = Number(pending);
        if (!Number.isInteger(pendingPartyId) || pendingPartyId < 0) {
            return { ok: false, code: messages.PartyErrorCode.INVITE_EXPIRED_OR_INVALID };
        }

        const inviterEntry = await this.unifiedRepo.findCharacterNameEntry(inviterName);
        if (!inviterEntry || Number(inviterEntry.world_id) !== Number(worldId)) {
            await client.del(inviteKey);
            return { ok: false, code: messages.PartyErrorCode.CHARACTER_NOT_FOUND };
        }
        const inviterCharacterId = Number(inviterEntry.character_id);
        const inviterSession = await this.sessionRepo.getCharacterSession(worldId, inviterCharacterId);
        const inviterChRaw = inviterSession?.gameServer?.channelId;
        const inviterChannelID = Number(inviterChRaw);
        const inviterOnline =
            inviterSession?.state === "ONLINE" && Boolean(inviterSession?.gameServer?.connected);
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
        await this._publishToPartyGameChannel(worldId, inviterChannelID, "party_invite_denied", {
            world_id: Number(worldId),
            inviter_character_id: inviterCharacterId,
            denied_character_name: deniedName,
            action: Number(action ?? 0),
        });
        return { ok: true };
    }

    async joinParty(worldId, partyId, member, skipInvitePendingCheck = false) {
        this._assertWorld(worldId);
        this._assertPartyId(partyId);
        if (!member || Number(member.worldId) !== Number(worldId)) {
            return { ok: false, code: messages.PartyErrorCode.UNKNOWN };
        }
        const characterId = Number(member.characterId);
        this._assertCharacterId(characterId);
        const name = String(member.characterName ?? "").trim();
        this._assertName(name);
        this._assertUInt16(member.level, "level");
        this._assertUInt16(member.classId, "class_id");
        const doorJson = this._doorJsonFromPayload(member.door);
        let joinChannelIndex = member.channelIndex;
        if (
            joinChannelIndex == null ||
            !Number.isFinite(Number(joinChannelIndex)) ||
            Number(joinChannelIndex) < -2
        ) {
            joinChannelIndex = await this._sessionChannelIndex(worldId, characterId);
        } else {
            joinChannelIndex = Number(joinChannelIndex);
        }

        const { client } = this.ctx.getRedisGlobalAccess(worldId);
        const inviteKey = this._invitePendingKey(worldId, characterId);
        if (!skipInvitePendingCheck) {
            const pending = await client.get(inviteKey);
            if (pending == null || Number(pending) !== Number(partyId)) {
                return { ok: false, code: messages.PartyErrorCode.INVITE_EXPIRED_OR_INVALID };
            }
        }

        const result = await this.ctx.withPgGlobalTransaction(worldId, async (txClient) => {
            const party = await this.partyRepo.get(worldId, partyId, { txClient });
            if (!party || party.state !== PARTY_STATE_ACTIVE) {
                return { ok: false, code: messages.PartyErrorCode.PARTY_NOT_FOUND };
            }

            const members = await this.partyMemberRepo.getAll(worldId, partyId, { txClient });
            if (members.size >= MAX_PARTY_MEMBERS) {
                return { ok: false, code: messages.PartyErrorCode.PARTY_FULL };
            }

            const state = await this.characterRealtimeStateRepo.get(worldId, characterId, { txClient });
            if (state?.partyId != null) {
                return { ok: false, code: messages.PartyErrorCode.ALREADY_IN_PARTY };
            }

            await this.partyMemberRepo.set(worldId, {
                worldId,
                partyId,
                characterId,
                characterName: name,
                level: Number(member.level),
                classId: Number(member.classId),
                role: "MEMBER",
                mapId: Number(member.mapId ?? 0),
                channelIndex: joinChannelIndex,
                door: doorJson,
            }, { txClient });
            const nextRevision = Number(party.revision) + 1;
            const updatedParty = await this.partyRepo.set(worldId, {
                ...party,
                revision: nextRevision,
            }, { txClient });
            await this.characterRealtimeStateRepo.set(worldId, {
                worldId,
                characterId,
                partyId,
                guildId: state?.guildId ?? null,
            }, { txClient });
            return { ok: true, partyId: updatedParty.partyId, revision: updatedParty.revision };
        });
        if (!result.ok) {
            return result;
        }
        await client.del(inviteKey);
        await this.partyRepo.evictCache(worldId, partyId);
        await this.partyMemberRepo.evictGroupCache(worldId, partyId);
        await this.characterRealtimeStateRepo.evictCache(worldId, characterId);

        await this._publishPartyEvent(EVT.MEMBER_JOINED, worldId, result.partyId, result.revision, {
            character_id: Number(characterId),
        });

        return result;
    }

    async leaveParty(worldId, characterId) {
        this._assertWorld(worldId);
        this._assertCharacterId(characterId);

        const result = await this.ctx.withPgGlobalTransaction(worldId, async (txClient) => {
            const state = await this.characterRealtimeStateRepo.get(worldId, characterId, { txClient });
            if (state?.partyId == null) {
                return { ok: false, code: messages.PartyErrorCode.NOT_IN_PARTY };
            }
            const partyId = state.partyId;
            const party = await this.partyRepo.get(worldId, partyId, { txClient });
            if (!party) {
                await this.characterRealtimeStateRepo.set(worldId, {
                    worldId,
                    characterId,
                    partyId: null,
                    guildId: state.guildId ?? null,
                }, { txClient });
                return { ok: true, partyId, revision: 0, disbanded: false };
            }

            await this.partyMemberRepo.del(worldId, partyId, characterId, { txClient });
            await this.characterRealtimeStateRepo.set(worldId, {
                worldId,
                characterId,
                partyId: null,
                guildId: state.guildId ?? null,
            }, { txClient });

            const remaining = await this.partyMemberRepo.getAll(worldId, partyId, { txClient });
            if (remaining.size === 0) {
                await this.partyRepo.delete({ worldId, partyId }, { txClient });
                return {
                    ok: true,
                    partyId,
                    revision: Number(party.revision) + 1,
                    disbanded: true,
                    realtimeStateCharacterIds: [Number(characterId)],
                };
            }

            const remainingMembers = [...remaining.values()];
            if (Number(party.leaderCharacterId) === Number(characterId) && remainingMembers.length === 1) {
                const disbandCharacterIds = [Number(characterId), Number(remainingMembers[0].characterId)];
                for (const memberCharacterId of disbandCharacterIds) {
                    const memberState = await this.characterRealtimeStateRepo.get(worldId, memberCharacterId, { txClient });
                    await this.characterRealtimeStateRepo.set(worldId, {
                        worldId,
                        characterId: memberCharacterId,
                        partyId: null,
                        guildId: memberState?.guildId ?? null,
                    }, { txClient });
                }
                await this.partyMemberRepo.del(worldId, partyId, remainingMembers[0].characterId, { txClient });
                await this.partyRepo.delete({ worldId, partyId }, { txClient });
                return {
                    ok: true,
                    partyId,
                    revision: Number(party.revision) + 1,
                    disbanded: true,
                    realtimeStateCharacterIds: disbandCharacterIds,
                };
            }

            const previousLeaderCharacterId = Number(party.leaderCharacterId);
            let nextLeaderCharacterId = previousLeaderCharacterId;
            if (Number(party.leaderCharacterId) === Number(characterId)) {
                const sorted = remainingMembers.sort((a, b) => Number(b.level) - Number(a.level));
                nextLeaderCharacterId = sorted[0].characterId;
                await this.partyMemberRepo.set(worldId, { ...sorted[0], role: "LEADER" }, { txClient });
            }
            const nextRevision = Number(party.revision) + 1;
            const updatedParty = await this.partyRepo.set(worldId, {
                ...party,
                leaderCharacterId: nextLeaderCharacterId,
                revision: nextRevision,
            }, { txClient });
            return {
                ok: true,
                partyId: updatedParty.partyId,
                revision: updatedParty.revision,
                disbanded: false,
                leaderChanged: previousLeaderCharacterId !== Number(nextLeaderCharacterId),
                oldLeaderCharacterId: previousLeaderCharacterId,
                newLeaderCharacterId: Number(nextLeaderCharacterId),
                realtimeStateCharacterIds: [Number(characterId)],
            };
        });
        if (!result.ok) {
            return result;
        }
        await this.partyRepo.evictCache(worldId, result.partyId);
        await this.partyMemberRepo.evictGroupCache(worldId, result.partyId);
        for (const affectedCharacterId of result.realtimeStateCharacterIds ?? [Number(characterId)]) {
            await this.characterRealtimeStateRepo.evictCache(worldId, affectedCharacterId);
        }

        if (result.disbanded) {
            await this._publishPartyEvent(EVT.DISBANDED, worldId, result.partyId, result.revision, {
                character_id: Number(characterId),
            });
        } else {
            await this._publishPartyEvent(EVT.MEMBER_LEFT, worldId, result.partyId, result.revision, {
                character_id: Number(characterId),
                leader_changed: Boolean(result.leaderChanged),
                old_leader_character_id: Number(result.oldLeaderCharacterId ?? 0),
                new_leader_character_id: Number(result.newLeaderCharacterId ?? 0),
            });
        }

        return result;
    }

    async expelParty(worldId, requesterCharacterId, targetCharacterId) {
        this._assertWorld(worldId);
        this._assertCharacterId(requesterCharacterId);
        this._assertCharacterId(targetCharacterId);
        if (Number(requesterCharacterId) === Number(targetCharacterId)) {
            return { ok: false, code: messages.PartyErrorCode.CANNOT_EXPEL_SELF };
        }

        const result = await this.ctx.withPgGlobalTransaction(worldId, async (txClient) => {
            const requesterState = await this.characterRealtimeStateRepo.get(worldId, requesterCharacterId, { txClient });
            if (requesterState?.partyId == null) {
                return { ok: false, code: messages.PartyErrorCode.NOT_IN_PARTY };
            }
            const partyId = requesterState.partyId;
            const party = await this.partyRepo.get(worldId, partyId, { txClient });
            if (!party || party.state !== PARTY_STATE_ACTIVE) {
                return { ok: false, code: messages.PartyErrorCode.PARTY_NOT_FOUND };
            }
            if (Number(party.leaderCharacterId) !== Number(requesterCharacterId)) {
                return { ok: false, code: messages.PartyErrorCode.NOT_PARTY_LEADER };
            }

            const members = await this.partyMemberRepo.getAll(worldId, partyId, { txClient });
            const targetMember = members.get(String(targetCharacterId));
            if (!targetMember) {
                return { ok: false, code: messages.PartyErrorCode.TARGET_NOT_IN_PARTY };
            }

            const targetState = await this.characterRealtimeStateRepo.get(worldId, targetCharacterId, { txClient });
            await this.partyMemberRepo.del(worldId, partyId, targetCharacterId, { txClient });
            await this.characterRealtimeStateRepo.set(worldId, {
                worldId,
                characterId: targetCharacterId,
                partyId: null,
                guildId: targetState?.guildId ?? null,
            }, { txClient });

            const remaining = await this.partyMemberRepo.getAll(worldId, partyId, { txClient });
            if (remaining.size === 0) {
                await this.partyRepo.delete({ worldId, partyId }, { txClient });
                return {
                    ok: true,
                    partyId,
                    revision: Number(party.revision) + 1,
                    disbanded: true,
                    realtimeStateCharacterIds: [Number(targetCharacterId)],
                };
            }

            const nextRevision = Number(party.revision) + 1;
            const updatedParty = await this.partyRepo.set(worldId, {
                ...party,
                revision: nextRevision,
            }, { txClient });
            return {
                ok: true,
                partyId: updatedParty.partyId,
                revision: updatedParty.revision,
                disbanded: false,
                realtimeStateCharacterIds: [Number(targetCharacterId)],
            };
        });
        if (!result.ok) {
            return result;
        }

        await this.partyRepo.evictCache(worldId, result.partyId);
        await this.partyMemberRepo.evictGroupCache(worldId, result.partyId);
        for (const affectedCharacterId of result.realtimeStateCharacterIds ?? [Number(targetCharacterId)]) {
            await this.characterRealtimeStateRepo.evictCache(worldId, affectedCharacterId);
        }

        if (result.disbanded) {
            await this._publishPartyEvent(EVT.DISBANDED, worldId, result.partyId, result.revision, {
                character_id: Number(targetCharacterId),
                expelled_by_character_id: Number(requesterCharacterId),
            });
        } else {
            await this._publishPartyEvent(EVT.MEMBER_LEFT, worldId, result.partyId, result.revision, {
                character_id: Number(targetCharacterId),
                expelled_by_character_id: Number(requesterCharacterId),
            });
        }

        return result;
    }

    async changePartyLeader(worldId, partyId, requesterCharacterId, newLeaderCharacterId) {
        this._assertWorld(worldId);
        this._assertPartyId(partyId);
        this._assertCharacterId(requesterCharacterId);
        this._assertCharacterId(newLeaderCharacterId);

        const newLeaderSession = await this.sessionRepo.getCharacterSession(
            worldId,
            newLeaderCharacterId
        );
        const newLeaderOnline =
            newLeaderSession?.state === "ONLINE" &&
            Boolean(newLeaderSession?.gameServer?.connected);
        if (!newLeaderOnline) {
            return { ok: false, code: messages.PartyErrorCode.TARGET_OFFLINE };
        }

        const result = await this.ctx.withPgGlobalTransaction(worldId, async (txClient) => {
            const party = await this.partyRepo.get(worldId, partyId, { txClient });
            if (!party || party.state !== PARTY_STATE_ACTIVE) {
                return { ok: false, code: messages.PartyErrorCode.PARTY_NOT_FOUND };
            }
            if (Number(party.leaderCharacterId) !== Number(requesterCharacterId)) {
                return { ok: false, code: messages.PartyErrorCode.NOT_PARTY_LEADER };
            }
            if (Number(requesterCharacterId) === Number(newLeaderCharacterId)) {
                return { ok: false, code: messages.PartyErrorCode.TARGET_ALREADY_LEADER };
            }

            const members = await this.partyMemberRepo.getAll(worldId, partyId, { txClient });
            const newLeader = members.get(String(newLeaderCharacterId));
            if (!newLeader) {
                return { ok: false, code: messages.PartyErrorCode.TARGET_NOT_IN_PARTY };
            }

            await this.partyMemberRepo.set(worldId, { ...newLeader, role: "LEADER" }, { txClient });
            const requester = members.get(String(requesterCharacterId));
            if (requester) {
                await this.partyMemberRepo.set(worldId, { ...requester, role: "MEMBER" }, { txClient });
            }
            const nextRevision = Number(party.revision) + 1;
            const updatedParty = await this.partyRepo.set(worldId, {
                ...party,
                leaderCharacterId: newLeaderCharacterId,
                revision: nextRevision,
            }, { txClient });
            return { ok: true, partyId: updatedParty.partyId, revision: updatedParty.revision };
        });
        if (!result.ok) {
            return result;
        }
        await this.partyRepo.evictCache(worldId, partyId);
        await this.partyMemberRepo.evictGroupCache(worldId, partyId);

        await this._publishPartyEvent(EVT.LEADER_CHANGED, worldId, result.partyId, result.revision, {
            old_leader_character_id: Number(requesterCharacterId),
            new_leader_character_id: Number(newLeaderCharacterId),
        });

        return result;
    }

    async broadcastMultiChat(worldId, memberId, senderCharacterId, chatMode, senderName, message) {
        this._assertWorld(worldId);
        this._assertPartyId(memberId);
        this._assertCharacterId(senderCharacterId);
        const trimmedMsg = String(message ?? "");
        if (trimmedMsg.length <= 0 || trimmedMsg.length > 500) {
            return { ok: false, code: messages.PartyErrorCode.UNKNOWN };
        }
        const name = String(senderName ?? "").trim();
        this._assertName(name);
        const mode = Number(chatMode);
        if (!Number.isInteger(mode) || mode < 0 || mode > 255) {
            return { ok: false, code: messages.PartyErrorCode.UNKNOWN };
        }
        await this._publishToPartyRoutes("multi_chat", worldId, Number(memberId), 0, {
            world_id: Number(worldId),
            member_id: Number(memberId),
            sender_character_id: Number(senderCharacterId),
            chat_mode: mode,
            sender_name: name,
            message: trimmedMsg,
        });
        return { ok: true, deliveredCount: 1 };
    }

    async getParty(worldId, partyId) {
        this._assertWorld(worldId);
        this._assertPartyId(partyId);
        const party = await this.partyRepo.get(worldId, partyId);
        if (!party) {
            return { found: false };
        }
        const members = await this.partyMemberRepo.getAll(worldId, partyId);
        return {
            found: true,
            party,
            members: this._sortPartyMemberModels(members, party.leaderCharacterId),
        };
    }
}

module.exports = { PartyService };
