import { CHARACTER_MODEL, CHARACTER_PERSISTED } from "../character-persisted";
import { bindingsFromProtoList } from "../key-layout-io";
import { INVENTORY_MODEL, INVENTORY_PERSISTED } from "../inventory-persisted";
import { SKILL_MODEL, SKILL_PERSISTED } from "../skill-persisted";
import { BUFF_MODEL, BUFF_PERSISTED } from "../buff-persisted";
import { grpcMapper } from "../mappers";
import type { CharacterOverviewListItem, CharacterOverviewService } from "../../services/character-overview-service";
import type { CharacterService } from "../../services/character-service";
import type { GrpcCall, GrpcCallback, GrpcErrorHandler } from "./types";
import type { CharacterModel } from "../../repos/character-repository";
import type { InventoryModel } from "../../repos/inventory-repository";
import type { SkillModel } from "../../repos/skill-repository";
import type { BuffModel } from "../../repos/buff-repository";
import type {
    BuffPersisted,
    CharacterPersisted,
    CharacterOverview,
    InventoryPersisted,
    CheckCharacterNameReply,
    CheckCharacterNameRequest,
    CreateCharacterReply,
    CreateCharacterRequest,
    DeleteCharacterReply,
    DeleteCharacterRequest,
    GetCharacterListReply,
    GetCharacterListRequest,
    SaveCharacterReply,
    SaveCharacterRequest,
    SaveCharactersReply,
    SaveCharactersRequest,
    SkillPersisted,
} from "../../protobuf/generated/fminternal/internal_service";

type CharacterListResult = Awaited<ReturnType<CharacterOverviewService["getCharacterList"]>>;
type CharacterListItem = CharacterListResult["characters"][number];

function protoMapToObject(protoMap: Record<number, number> | undefined): Record<string, number> {
    const obj: Record<string, number> = {};
    if (protoMap) {
        for (const [key, value] of Object.entries(protoMap)) {
            obj[key] = value;
        }
    }
    return obj;
}

function makeCharacterOverview(model: CharacterListItem): CharacterOverview {
    return {
        characterId: model.characterId,
        name: model.name,
        gender: model.gender,
        skinColor: model.skinColor,
        face: model.face,
        hair: model.hair,
        level: model.level,
        classId: model.classId,
        mapId: model.mapId,
        spawnPoint: model.spawnPoint,
        accountId: model.accountId ?? 0,
        worldId: model.worldId ?? 0,
        rank: model.rank ?? 0,
        rankDiff: model.rankDiff ?? 0,
        classRank: model.classRank ?? 0,
        classRankDiff: model.classRankDiff ?? 0,
        baseLooks: (model.baseLooks ?? {}) as Record<number, number>,
        overlays: (model.overlays ?? {}) as Record<number, number>,
    };
}

export function createCharacterHandlers(
    characterService: CharacterService,
    characterOverviewService: CharacterOverviewService,
    grpcError: GrpcErrorHandler
) {
    return {
        async saveCharacter(call: GrpcCall<SaveCharacterRequest>, callback: GrpcCallback<SaveCharacterReply>) {
            try {
                const msgChar = call.request.character;
                if (!msgChar) {
                    throw Object.assign(new Error("character is required"), { code: "INVALID_PAYLOAD" });
                }
                const persisted = grpcMapper.map<CharacterPersisted, CharacterModel>(
                    msgChar,
                    CHARACTER_PERSISTED,
                    CHARACTER_MODEL
                );
                const baseLooks = protoMapToObject(call.request.baseLooks);
                const overlays = protoMapToObject(call.request.overlays);
                await characterService.saveCharacter(persisted, baseLooks, overlays);
                callback(null, { ok: true });
            } catch (err) {
                grpcError(err, callback);
            }
        },
        async saveCharacters(call: GrpcCall<SaveCharactersRequest>, callback: GrpcCallback<SaveCharactersReply>) {
            try {
                const entries = call.request.entries.map((entry) => {
                    const msgChar = entry.character;
                    if (!msgChar) {
                        throw Object.assign(new Error("character is required in entry"), { code: "INVALID_PAYLOAD" });
                    }
                    return {
                        persisted: grpcMapper.map<CharacterPersisted, CharacterModel>(
                            msgChar,
                            CHARACTER_PERSISTED,
                            CHARACTER_MODEL
                        ),
                        baseLooks: protoMapToObject(entry.baseLooks),
                        overlays: protoMapToObject(entry.overlays),
                        inventory: entry.inventory.map((inventory) =>
                            grpcMapper.map<InventoryPersisted, InventoryModel>(
                                inventory,
                                INVENTORY_PERSISTED,
                                INVENTORY_MODEL
                            )
                        ),
                        skills: entry.skills.map((skill) =>
                            grpcMapper.map<SkillPersisted, SkillModel>(
                                skill,
                                SKILL_PERSISTED,
                                SKILL_MODEL
                            )
                        ),
                        buffs: (entry.buffs ?? []).map((buff) =>
                            grpcMapper.map<BuffPersisted, BuffModel>(
                                buff,
                                BUFF_PERSISTED,
                                BUFF_MODEL
                            )
                        ),
                        keyLayout: bindingsFromProtoList(entry.keyLayout),
                    };
                });
                await characterService.saveCharacters(entries);
                callback(null, { ok: true });
            } catch (err) {
                grpcError(err, callback);
            }
        },
        async getCharacterList(call: GrpcCall<GetCharacterListRequest>, callback: GrpcCallback<GetCharacterListReply>) {
            try {
                const result = await characterOverviewService.getCharacterList(call.request.accountId, call.request.worldId);
                callback(null, { characters: result.characters.map(makeCharacterOverview), slotCount: result.slotCount });
            } catch (err) {
                grpcError(err, callback);
            }
        },
        async checkCharacterName(call: GrpcCall<CheckCharacterNameRequest>, callback: GrpcCallback<CheckCharacterNameReply>) {
            try {
                const result = await characterOverviewService.checkCharacterName(call.request.name);
                callback(null, { exists: result.exists });
            } catch (err) {
                grpcError(err, callback);
            }
        },
        async createCharacter(call: GrpcCall<CreateCharacterRequest>, callback: GrpcCallback<CreateCharacterReply>) {
            try {
                const result = await characterService.createCharacter(call.request.accountId, call.request.worldId, call.request);
                callback(null, {
                    success: result.success,
                    errorMsg: result.errorMsg ?? "",
                    character: result.success && result.character ? makeCharacterOverview(result.character) : undefined,
                });
            } catch (err) {
                grpcError(err, callback);
            }
        },
        async deleteCharacter(call: GrpcCall<DeleteCharacterRequest>, callback: GrpcCallback<DeleteCharacterReply>) {
            try {
                const result = await characterService.deleteCharacter(call.request.accountId, call.request.characterId);
                callback(null, { success: result.success });
            } catch (err) {
                grpcError(err, callback);
            }
        },
    };
}
