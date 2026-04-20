"use strict";

const { fillMessageFromPersisted: fillCharacterMessage } = require("../character-persisted");
const { makeKeyLayoutProtoList } = require("../key-layout-io");
const { makeInventoryMessage } = require("../inventory-persisted");
const { makeSkillMessage } = require("../skill-persisted");

function createSessionHandlers(
    characterService,
    inventoryRepository,
    skillService,
    sessionService,
    internalConfig,
    characterRealtimeStateRepository,
    messages,
    grpcError
) {
    return {
        async beginGameTransition(call, callback) {
            try {
                const worldId = call.request.getWorldId();
                const accountId = call.request.getAccountId();
                const characterId = call.request.getCharacterId();
                const row = await characterService.getCharacter(worldId, characterId);
                if (!row) {
                    throw Object.assign(new Error(`character not found: ${characterId}`), { code: "INVALID_CHARACTER_ID" });
                }
                if (Number(row.accountId) !== Number(accountId)) {
                    throw Object.assign(new Error(`account mismatch for character ${characterId}`), { code: "INVALID_PAYLOAD" });
                }
                const transition = await sessionService.beginTransition(
                    worldId,
                    row.accountId,
                    row.characterId,
                    row.name
                );
                const reply = new messages.BeginGameTransitionReply();
                reply.setOk(Boolean(transition.ok));
                reply.setErrorCode(transition.code ?? messages.SessionErrorCode.SESSION_UNKNOWN);
                callback(null, reply);
            } catch (err) {
                grpcError(err, callback);
            }
        },

        async enterGame(call, callback) {
            try {
                const worldId = call.request.getWorldId();
                const characterId = call.request.getCharacterId();
                const channelId = call.request.getChannelId();
                const worldCfg = internalConfig.game_servers?.worlds?.[String(worldId)];
                if (!worldCfg) {
                    throw Object.assign(new Error(`Unknown world_id for game attach: ${worldId}`), { code: "UNKNOWN_WORLD" });
                }
                const hasChannel = (worldCfg.channels ?? []).some((ch) => Number(ch.channel_id) === Number(channelId));
                if (!hasChannel) {
                    throw Object.assign(new Error(`Unknown channel_id for world ${worldId}: ${channelId}`), { code: "UNKNOWN_CHANNEL" });
                }

                const row = await characterService.getCharacter(worldId, characterId);
                const reply = new messages.EnterGameReply();
                if (!row) {
                    reply.setFound(false);
                    callback(null, reply);
                    return;
                }

                const attach = await sessionService.attachGameSession(worldId, row.accountId, channelId);
                if (!attach.ok) {
                    throw new Error(`attach game session failed: ${attach.code}`);
                }

                const [inventoryList, skillList, keyLayoutBindings] = await Promise.all([
                    inventoryRepository.getAll(worldId, characterId).then((map) => [...map.values()]),
                    skillService.getSkills(worldId, characterId),
                    characterService.getKeyLayoutBindings(worldId, characterId),
                ]);

                const realtime = await characterRealtimeStateRepository.get(worldId, characterId);
                const guildId = realtime?.guildId ?? 0;

                reply.setFound(true);
                const ch = new messages.CharacterPersisted();
                fillCharacterMessage(ch, row);
                reply.setCharacter(ch);
                reply.setInventoryList(inventoryList.map(makeInventoryMessage));
                reply.setSkillsList(skillList.map(makeSkillMessage));
                reply.setKeyLayoutList(makeKeyLayoutProtoList(messages, keyLayoutBindings));
                if (realtime != null && realtime.partyId != null) {
                    reply.setPartyId(Number(realtime.partyId));
                }
                reply.setGuildId(guildId);
                callback(null, reply);
            } catch (err) {
                grpcError(err, callback);
            }
        },

        async refreshSession(call, callback) {
            try {
                const result = await sessionService.refresh(
                    call.request.getWorldId(),
                    call.request.getAccountId()
                );
                const reply = new messages.RefreshSessionReply();
                reply.setOk(Boolean(result.ok));
                reply.setErrorCode(result.code ?? messages.SessionErrorCode.SESSION_UNKNOWN);
                callback(null, reply);
            } catch (err) {
                grpcError(err, callback);
            }
        },

        async logoutSession(call, callback) {
            try {
                const result = await sessionService.logout(
                    call.request.getWorldId(),
                    call.request.getAccountId(),
                    {
                        disconnectSource: call.request.getDisconnectSource(),
                        transferDisconnect: call.request.getTransferDisconnect(),
                    }
                );
                const reply = new messages.LogoutSessionReply();
                reply.setOk(Boolean(result.ok));
                callback(null, reply);
            } catch (err) {
                grpcError(err, callback);
            }
        },
    };
}

module.exports = { createSessionHandlers };
