"use strict";

const { fillMessageFromPersisted: fillCharacterMessage } = require("../character-persisted");
const { makeInventoryMessage } = require("../inventory-persisted");
const { makeSkillMessage } = require("../skill-persisted");

function createSessionHandlers(
    characterService,
    inventoryService,
    skillService,
    sessionService,
    internalConfig,
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
                reply.setErrorCode(transition.code ?? "");
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

                const [inventoryList, skillList] = await Promise.all([
                    inventoryService.getInventory(worldId, characterId),
                    skillService.getSkills(worldId, characterId),
                ]);
                reply.setFound(true);
                const ch = new messages.CharacterPersisted();
                fillCharacterMessage(ch, row);
                reply.setCharacter(ch);
                reply.setInventoryList(inventoryList.map(makeInventoryMessage));
                reply.setSkillsList(skillList.map(makeSkillMessage));
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
                reply.setErrorCode(result.code ?? "");
                callback(null, reply);
            } catch (err) {
                grpcError(err, callback);
            }
        },

        async logoutSession(call, callback) {
            try {
                const result = await sessionService.logout(
                    call.request.getWorldId(),
                    call.request.getAccountId()
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
