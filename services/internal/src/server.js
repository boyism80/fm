"use strict";

const grpc = require("@grpc/grpc-js");
const { InternalService } = require("../protobuf/fminternal/ping_grpc_pb.js");
const messages = require("../protobuf/fminternal/ping_pb.js");
const { createAppContainer } = require("./container");
const { autoMigrateAllIfEnabled } = require("./auto-migrate");
const { fillMessageFromPersisted: fillCharacterMessage, persistedFromMessage: characterFromMessage } = require("./grpc/character-persisted");
const { makeInventoryMessage, persistedFromMessage: inventoryFromMessage } = require("./grpc/inventory-persisted");
const { makeSkillMessage } = require("./grpc/skill-persisted");

const INVALID_CODES = new Set([
    "UNKNOWN_WORLD", "INVALID_CHARACTER_ID", "INVALID_OWNER_ID",
    "INVALID_UNIQUE_ID", "INVALID_PAYLOAD",
]);

function grpcError(err, callback) {
    const code = err && INVALID_CODES.has(err.code)
        ? grpc.status.INVALID_ARGUMENT
        : grpc.status.INTERNAL;
    callback({ code, message: err.message || String(err) });
}

function _protoMapToObject(protoMap) {
    const obj = {};
    if (protoMap) {
        protoMap.forEach((value, key) => {
            obj[key] = value;
        });
    }
    return obj;
}

function _makeCharacterOverview(model) {
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

async function main() {
    const container = createAppContainer();
    const appConfiguration = container.resolve("appConfiguration");
    const internalConfig = appConfiguration.raw;
    await autoMigrateAllIfEnabled(internalConfig);
    const internalContext = container.resolve("internalContext");
    const characterService = container.resolve("characterService");
    const inventoryService = container.resolve("inventoryService");
    const skillService = container.resolve("skillService");
    const accountService = container.resolve("accountService");
    const characterOverviewService = container.resolve("characterOverviewService");

    const wid = String(internalConfig.app.world_id);
    const pgw = internalConfig.postgresql.worlds[wid];
    const rgw = internalConfig.redis.worlds[wid];
    const pgUnified = internalConfig.postgresql.unified;
    console.log(
        `fm internal: loaded ${internalConfig.configPath} | grpc ${internalConfig.grpc.host}:${internalConfig.grpc.port} | world ${wid} | ` +
            `pg global ${pgw.global.host}:${pgw.global.port}/${pgw.global.database} | pg_data_shards ${pgw.data.length} | ` +
            `redis global ${rgw.global.host}:${rgw.global.port} | redis_data_shards ${rgw.data.length} | ` +
            `character_cache_ttl_s ${appConfiguration.getCharacterCacheTtlSeconds()} | item_cache_ttl_s ${appConfiguration.getItemCacheTtlSeconds()}` +
            (pgUnified ? ` | pg_unified ${pgUnified.host}:${pgUnified.port}/${pgUnified.database}` : "")
    );

    const server = new grpc.Server();

    server.addService(InternalService, {
        ping(_call, callback) {
            const reply = new messages.PingReply();
            reply.setMessage("pong");
            callback(null, reply);
        },

        async getCharacter(call, callback) {
            try {
                const worldId     = call.request.getWorldId();
                const characterId = call.request.getCharacterId();
                const [row, inventoryList, skillList] = await Promise.all([
                    characterService.getCharacter(worldId, characterId),
                    inventoryService.getInventory(worldId, characterId),
                    skillService.getSkills(worldId, characterId),
                ]);
                const reply = new messages.GetCharacterReply();
                if (!row) {
                    reply.setFound(false);
                } else {
                    reply.setFound(true);
                    const ch = new messages.CharacterPersisted();
                    fillCharacterMessage(ch, row);
                    reply.setCharacter(ch);
                    reply.setInventoryList(inventoryList.map(makeInventoryMessage));
                    reply.setSkillsList(skillList.map(makeSkillMessage));
                }
                callback(null, reply);
            } catch (err) {
                grpcError(err, callback);
            }
        },

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
                        persisted:  characterFromMessage(msgChar),
                        baseLooks:  _protoMapToObject(entry.getBaseLooksMap()),
                        overlays:   _protoMapToObject(entry.getOverlaysMap()),
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

        async getInventory(call, callback) {
            try {
                const items = await inventoryService.getInventory(
                    call.request.getWorldId(),
                    call.request.getOwnerId()
                );
                const reply = new messages.GetInventoryReply();
                reply.setItemsList(items.map(makeInventoryMessage));
                callback(null, reply);
            } catch (err) {
                grpcError(err, callback);
            }
        },

        async saveInventory(call, callback) {
            try {
                const worldId = call.request.getWorldId();
                const msgItems = call.request.getItemsList();
                if (msgItems.length) {
                    await inventoryService.saveInventory(worldId, msgItems.map(inventoryFromMessage));
                }
                const reply = new messages.SaveInventoryReply();
                reply.setOk(true);
                callback(null, reply);
            } catch (err) {
                grpcError(err, callback);
            }
        },

        async deleteInventory(call, callback) {
            try {
                await inventoryService.deleteInventory(
                    call.request.getWorldId(),
                    call.request.getOwnerId(),
                    call.request.getUniqueIdsList().map(String)
                );
                const reply = new messages.DeleteInventoryReply();
                reply.setOk(true);
                callback(null, reply);
            } catch (err) {
                grpcError(err, callback);
            }
        },

        async loginAccount(call, callback) {
            try {
                const result = await accountService.loginAccount(
                    call.request.getLoginId(),
                    call.request.getPassword(),
                    call.request.getMacAddress(),
                    call.request.getIpAddress(),
                    call.request.getInitialRole()
                );
                const reply = new messages.LoginAccountReply();
                reply.setStatus(result.status);
                reply.setAccountId(result.accountId ?? 0);
                reply.setGender(result.gender ?? 0);
                reply.setIsChatBlocked(result.isChatBlocked ?? false);
                const blockedUntil = result.chatBlockedUntil
                    ? new Date(result.chatBlockedUntil).getTime()
                    : 0;
                reply.setChatBlockedUntilUnixMs(blockedUntil);
                reply.setCharacterSlotCount(result.characterSlotCount ?? 6);
                reply.setRole(result.role ?? 0);
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
                reply.setCharactersList(result.characters.map(_makeCharacterOverview));
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
                const result = await characterOverviewService.createCharacter(
                    call.request.getAccountId(),
                    call.request.getWorldId(),
                    {
                        name:         call.request.getName(),
                        face:         call.request.getFace(),
                        hair:         call.request.getHair(),
                        skinColor:    call.request.getSkinColor(),
                        topItemId:    call.request.getTopItemId(),
                        bottomItemId: call.request.getBottomItemId(),
                        shoesItemId:  call.request.getShoesItemId(),
                        weaponItemId: call.request.getWeaponItemId(),
                    }
                );
                const reply = new messages.CreateCharacterReply();
                reply.setSuccess(result.success);
                reply.setErrorMsg(result.errorMsg ?? "");
                if (result.success && result.character) {
                    reply.setCharacter(_makeCharacterOverview(result.character));
                }
                callback(null, reply);
            } catch (err) {
                grpcError(err, callback);
            }
        },

        async deleteCharacter(call, callback) {
            try {
                const result = await characterOverviewService.deleteCharacter(
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
    });

    const addr = `${internalConfig.grpc.host}:${internalConfig.grpc.port}`;

    server.bindAsync(addr, grpc.ServerCredentials.createInsecure(), (err, boundPort) => {
        if (err) {
            console.error(err);
            process.exit(1);
        }
        server.start();
        console.log(`fm internal gRPC listening on ${addr} (bound port ${boundPort})`);
    });

    async function shutdown() {
        server.tryShutdown(async () => {
            await internalContext.close().catch(() => {});
            process.exit(0);
        });
        setTimeout(() => {
            internalContext.close().catch(() => {});
            process.exit(1);
        }, 10_000).unref();
    }

    process.on("SIGINT", shutdown);
    process.on("SIGTERM", shutdown);
}

main().catch((err) => {
    console.error("[startup] fatal:", err);
    process.exit(1);
});
