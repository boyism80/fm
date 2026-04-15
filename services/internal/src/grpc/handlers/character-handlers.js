"use strict";

const { persistedFromMessage: characterFromMessage } = require("../character-persisted");
const { persistedFromMessage: inventoryFromMessage } = require("../inventory-persisted");
const { persistedFromMessage: skillFromMessage } = require("../skill-persisted");

function _protoMapToObject(protoMap) {
    const obj = {};
    if (protoMap) {
        protoMap.forEach((value, key) => {
            obj[key] = value;
        });
    }
    return obj;
}

function _makeCharacterOverview(messages, model) {
    const msg = new messages.CharacterOverview();
    msg.setCharacterId(model.characterId);
    msg.setName(model.name);
    msg.setGender(model.gender);
    msg.setSkinColor(model.skinColor);
    msg.setFace(model.face);
    msg.setHair(model.hair);
    msg.setLevel(model.level);
    msg.setClassId(model.classId);
    msg.setMapId(model.mapId);
    msg.setSpawnPoint(model.spawnPoint);
    msg.setAccountId(model.accountId ?? 0);
    msg.setWorldId(model.worldId ?? 0);
    msg.setRank(model.rank ?? 0);
    msg.setRankDiff(model.rankDiff ?? 0);
    msg.setClassRank(model.classRank ?? 0);
    msg.setClassRankDiff(model.classRankDiff ?? 0);
    if (model.baseLooks) {
        const m = msg.getBaseLooksMap();
        for (const [k, v] of Object.entries(model.baseLooks)) {
            m.set(Number(k), Number(v));
        }
    }
    if (model.overlays) {
        const m = msg.getOverlaysMap();
        for (const [k, v] of Object.entries(model.overlays)) {
            m.set(Number(k), Number(v));
        }
    }
    return msg;
}

function createCharacterHandlers(
    characterService,
    characterOverviewService,
    messages,
    grpcError
) {
    return {
        async saveCharacter(call, callback) {
            try {
                const msgChar = call.request.getCharacter();
                if (!msgChar) {
                    throw Object.assign(new Error("character is required"), { code: "INVALID_PAYLOAD" });
                }
                const persisted = characterFromMessage(msgChar);
                const baseLooks = _protoMapToObject(call.request.getBaseLooksMap());
                const overlays = _protoMapToObject(call.request.getOverlaysMap());
                await characterService.saveCharacter(persisted, baseLooks, overlays);
                const reply = new messages.SaveCharacterReply();
                reply.setOk(true);
                callback(null, reply);
            } catch (err) {
                grpcError(err, callback);
            }
        },

        async saveCharacters(call, callback) {
            try {
                const entriesList = call.request.getEntriesList();
                const entries = entriesList.map((entry) => {
                    const msgChar = entry.getCharacter();
                    if (!msgChar) {
                        throw Object.assign(new Error("character is required in entry"), { code: "INVALID_PAYLOAD" });
                    }
                    return {
                        persisted: characterFromMessage(msgChar),
                        baseLooks: _protoMapToObject(entry.getBaseLooksMap()),
                        overlays: _protoMapToObject(entry.getOverlaysMap()),
                        inventory: entry.getInventoryList().map(inventoryFromMessage),
                        skills: entry.getSkillsList().map(skillFromMessage),
                    };
                });
                await characterService.saveCharacters(entries);
                const reply = new messages.SaveCharactersReply();
                reply.setOk(true);
                callback(null, reply);
            } catch (err) {
                grpcError(err, callback);
            }
        },

        async getCharacterList(call, callback) {
            try {
                const result = await characterOverviewService.getCharacterList(
                    call.request.getAccountId(),
                    call.request.getWorldId()
                );
                const reply = new messages.GetCharacterListReply();
                reply.setCharactersList(result.characters.map((model) => _makeCharacterOverview(messages, model)));
                reply.setSlotCount(result.slotCount);
                callback(null, reply);
            } catch (err) {
                grpcError(err, callback);
            }
        },

        async checkCharacterName(call, callback) {
            try {
                const result = await characterOverviewService.checkCharacterName(
                    call.request.getName()
                );
                const reply = new messages.CheckCharacterNameReply();
                reply.setExists(result.exists);
                callback(null, reply);
            } catch (err) {
                grpcError(err, callback);
            }
        },

        async createCharacter(call, callback) {
            try {
                const result = await characterService.createCharacter(
                    call.request.getAccountId(),
                    call.request.getWorldId(),
                    {
                        name: call.request.getName(),
                        face: call.request.getFace(),
                        hair: call.request.getHair(),
                        skinColor: call.request.getSkinColor(),
                        topItemId: call.request.getTopItemId(),
                        bottomItemId: call.request.getBottomItemId(),
                        shoesItemId: call.request.getShoesItemId(),
                        weaponItemId: call.request.getWeaponItemId(),
                    }
                );
                const reply = new messages.CreateCharacterReply();
                reply.setSuccess(result.success);
                reply.setErrorMsg(result.errorMsg ?? "");
                if (result.success && result.character) {
                    reply.setCharacter(_makeCharacterOverview(messages, result.character));
                }
                callback(null, reply);
            } catch (err) {
                grpcError(err, callback);
            }
        },

        async deleteCharacter(call, callback) {
            try {
                const result = await characterService.deleteCharacter(
                    call.request.getAccountId(),
                    call.request.getCharacterId()
                );
                const reply = new messages.DeleteCharacterReply();
                reply.setSuccess(result.success);
                callback(null, reply);
            } catch (err) {
                grpcError(err, callback);
            }
        },
    };
}

module.exports = { createCharacterHandlers };
