import { CHARACTER_MODEL, CHARACTER_PROTO } from "../character-proto";
import { bindingsFromProtoList } from "../key-layout-io";
import { INVENTORY_MODEL, INVENTORY_PROTO } from "../inventory-proto";
import { SKILL_MODEL, SKILL_PROTO } from "../skill-proto";
import { BUFF_MODEL, BUFF_PROTO } from "../buff-proto";
import { QUEST_MODEL, QUEST_PROTO } from "../quest-proto";
import { SAVED_LOCATION_MODEL, SAVED_LOCATION_PROTO } from "../saved-location-proto";
import { grpcMapper } from "../mappers";
import type { CharacterOverviewListItem, CharacterOverviewService } from "../../services/character-overview-service";
import type { CharacterService, SaveCharacterEntry } from "../../services/character-service";
import type { GrpcCall, GrpcCallback, GrpcErrorHandler } from "./types";
import type { CharacterModel } from "../../repos/character-repository";
import type { InventoryModel } from "../../repos/inventory-repository";
import type { SkillModel } from "../../repos/skill-repository";
import type { BuffModel } from "../../repos/buff-repository";
import type { QuestModel } from "../../repos/quest-repository";
import type { SavedLocationModel } from "../../repos/saved-location-repository";
import type {
    Buff,
    Character,
    CharacterOverview,
    CharacterSaveEntry,
    Inventory,
    CheckCharacterNameReply,
    CheckCharacterNameRequest,
    CreateCharacterReply,
    CreateCharacterRequest,
    DeleteCharacterReply,
    DeleteCharacterRequest,
    GetCharacterListReply,
    GetCharacterListRequest,
    LoadStorageReply,
    LoadStorageRequest,
    SaveCharacterReply,
    SaveCharacterRequest,
    SaveCharactersReply,
    SaveCharactersRequest,
    Quest,
    SavedLocation,
    Skill,
} from "../../protobuf/generated/fminternal/internal_service";
import { Controller, Method } from "../grpc-method-decorator";

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

export function makeSaveCharacterEntry(entry: CharacterSaveEntry): SaveCharacterEntry {
    const msgChar = entry.character;
    if (!msgChar) {
        throw Object.assign(new Error("character is required in entry"), { code: "INVALID_PAYLOAD" });
    }
    return {
        character: grpcMapper.map<Character, CharacterModel>(
            msgChar,
            CHARACTER_PROTO,
            CHARACTER_MODEL
        ),
        baseLooks: protoMapToObject(entry.baseLooks),
        overlays: protoMapToObject(entry.overlays),
        inventory: entry.inventory.map((inventory) =>
            grpcMapper.map<Inventory, InventoryModel>(
                inventory,
                INVENTORY_PROTO,
                INVENTORY_MODEL
            )
        ),
        skills: entry.skills.map((skill) =>
            grpcMapper.map<Skill, SkillModel>(
                skill,
                SKILL_PROTO,
                SKILL_MODEL
            )
        ),
        buffs: (entry.buffs ?? []).map((buff) =>
            grpcMapper.map<Buff, BuffModel>(
                buff,
                BUFF_PROTO,
                BUFF_MODEL
            )
        ),
        quests: (entry.quests ?? []).map((quest) =>
            grpcMapper.map<Quest, QuestModel>(
                quest,
                QUEST_PROTO,
                QUEST_MODEL
            )
        ),
        savedLocations: (entry.savedLocations ?? []).map((loc) =>
            grpcMapper.map<SavedLocation, SavedLocationModel>(
                loc,
                SAVED_LOCATION_PROTO,
                SAVED_LOCATION_MODEL
            )
        ),
        keyLayout: bindingsFromProtoList(entry.keyLayout),
        storage: entry.storage
            ? {
                storage: {
                    accountId: entry.storage.accountId,
                    worldId: entry.storage.worldId,
                    slots: entry.storage.slots,
                    meso: entry.storage.meso,
                },
                items: entry.storage.items.map((item) =>
                    grpcMapper.map<Inventory, InventoryModel>(
                        item,
                        INVENTORY_PROTO,
                        INVENTORY_MODEL
                    )
                ),
            }
            : undefined,
    };
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

@Controller("characterController")
export class CharacterGrpcController {
    private readonly characterService: CharacterService;
    private readonly characterOverviewService: CharacterOverviewService;
    private readonly grpcError: GrpcErrorHandler;

    constructor(
        characterService: CharacterService,
        characterOverviewService: CharacterOverviewService,
        grpcError: GrpcErrorHandler
    ) {
        this.characterService = characterService;
        this.characterOverviewService = characterOverviewService;
        this.grpcError = grpcError;
    }

    @Method("saveCharacter")
    async saveCharacter(call: GrpcCall<SaveCharacterRequest>, callback: GrpcCallback<SaveCharacterReply>) {
        try {
            const msgChar = call.request.character;
            if (!msgChar) {
                throw Object.assign(new Error("character is required"), { code: "INVALID_PAYLOAD" });
            }
            const character = grpcMapper.map<Character, CharacterModel>(
                msgChar,
                CHARACTER_PROTO,
                CHARACTER_MODEL
            );
            const baseLooks = protoMapToObject(call.request.baseLooks);
            const overlays = protoMapToObject(call.request.overlays);
            await this.characterService.saveCharacter(character, baseLooks, overlays);
            callback(null, { ok: true });
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("saveCharacters")
    async saveCharacters(call: GrpcCall<SaveCharactersRequest>, callback: GrpcCallback<SaveCharactersReply>) {
        try {
            const entries = call.request.entries.map(makeSaveCharacterEntry);
            await this.characterService.saveCharacters(entries);
            callback(null, { ok: true });
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("loadStorage")
    async loadStorage(call: GrpcCall<LoadStorageRequest>, callback: GrpcCallback<LoadStorageReply>) {
        try {
            const result = await this.characterService.loadStorage(call.request.worldId, call.request.accountId);
            callback(null, {
                storage: {
                    accountId: result.storage.accountId,
                    worldId: result.storage.worldId,
                    slots: result.storage.slots,
                    meso: result.storage.meso,
                    items: result.items.map((item) =>
                        grpcMapper.map<InventoryModel, Inventory>(
                            item,
                            INVENTORY_MODEL,
                            INVENTORY_PROTO
                        )
                    ),
                },
            });
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("getCharacterList")
    async getCharacterList(call: GrpcCall<GetCharacterListRequest>, callback: GrpcCallback<GetCharacterListReply>) {
        try {
            const result = await this.characterOverviewService.getCharacterList(call.request.accountId, call.request.worldId);
            callback(null, { characters: result.characters.map(makeCharacterOverview), slotCount: result.slotCount });
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("checkCharacterName")
    async checkCharacterName(call: GrpcCall<CheckCharacterNameRequest>, callback: GrpcCallback<CheckCharacterNameReply>) {
        try {
            const result = await this.characterOverviewService.checkCharacterName(call.request.name, call.request.accountId);
            callback(null, { exists: result.exists });
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("createCharacter")
    async createCharacter(call: GrpcCall<CreateCharacterRequest>, callback: GrpcCallback<CreateCharacterReply>) {
        try {
            const result = await this.characterService.createCharacter(call.request.accountId, call.request.worldId, call.request);
            callback(null, {
                success: result.success,
                errorMsg: result.errorMsg ?? "",
                character: result.character ? makeCharacterOverview(result.character) : undefined,
            });
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("deleteCharacter")
    async deleteCharacter(call: GrpcCall<DeleteCharacterRequest>, callback: GrpcCallback<DeleteCharacterReply>) {
        try {
            const result = await this.characterService.deleteCharacter(call.request.accountId, call.request.characterId);
            callback(null, { success: result.success });
        } catch (err) {
            this.grpcError(err, callback);
        }
    }
}
