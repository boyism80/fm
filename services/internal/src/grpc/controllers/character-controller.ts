import { CHARACTER_MODEL, CHARACTER_PROTO } from "../character-proto";
import { bindingsFromProtoList } from "../key-layout-io";
import { INVENTORY_MODEL, INVENTORY_PROTO } from "../inventory-proto";
import { SKILL_MODEL, SKILL_PROTO } from "../skill-proto";
import { BUFF_MODEL, BUFF_PROTO } from "../buff-proto";
import { QUEST_MODEL, QUEST_PROTO } from "../quest-proto";
import { SAVED_LOCATION_MODEL, SAVED_LOCATION_PROTO } from "../saved-location-proto";
import { grpcMapper } from "../mappers";
import type { CharacterOverviewService } from "../../services/character-overview-service";
import type { CharacterRank } from "../../services/ranking-service";
import type { CharacterService, SaveCharacterEntry } from "../../services/character-service";
import type { GrpcCall, GrpcCallback, GrpcErrorHandler } from "./types";
import type { CharacterModel } from "@fm/shared/repos/character-repository";
import type { CharacterOverviewModel } from "@fm/shared/repos/character-overview-repository";
import type { InventoryModel } from "@fm/shared/repos/inventory-repository";
import type { SkillModel } from "@fm/shared/repos/skill-repository";
import type { BuffModel } from "@fm/shared/repos/buff-repository";
import type { QuestModel } from "@fm/shared/repos/quest-repository";
import type { SavedLocationModel } from "@fm/shared/repos/saved-location-repository";
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
    SaveCharacterReply,
    SaveCharacterRequest,
    SaveCharactersReply,
    SaveCharactersRequest,
    Quest,
    SavedLocation,
    Skill,
} from "@fm/shared/protobuf/generated/fminternal/internal_service";
import { Controller, Method } from "../grpc-method-decorator";

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
        records: (entry.records ?? []).map((record) => ({
            ...record,
            ownerId: msgChar.characterId,
        })),
        accountRecords: (entry.account?.records ?? []).map((record) => ({
            ...record,
            ownerId: msgChar.accountId,
        })),
        monsterBook: (entry.monsterBook ?? []).map((card) => ({
            characterId: entry.character?.characterId ?? 0,
            cardId: card.cardId,
            count: card.count,
        })),
        keyLayout: bindingsFromProtoList(entry.keyLayout),
        skillMacros: (entry.skillMacros ?? []).map((macro) => ({
            characterId: msgChar.characterId,
            position: macro.position,
            name: macro.name,
            shout: macro.shout,
            skill1: macro.skill1,
            skill2: macro.skill2,
            skill3: macro.skill3,
        })),
        mount: entry.mount
            ? {
                level: entry.mount.level,
                exp: entry.mount.exp,
                fatigue: entry.mount.fatigue,
            }
            : undefined,
        storage: entry.account?.storage
            ? {
                storage: {
                    accountId: entry.account.storage.accountId,
                    worldId: entry.account.storage.worldId,
                    slots: entry.account.storage.slots,
                    meso: entry.account.storage.meso,
                },
                items: entry.account.storage.items.map((item) =>
                    grpcMapper.map<Inventory, InventoryModel>(
                        item,
                        INVENTORY_PROTO,
                        INVENTORY_MODEL
                    )
                ),
            }
            : undefined,
        claimedShops: entry.claimedShops ?? [],
        unsettledRequests: (entry.unsettledRequests ?? []).map((request) => ({
            id: request.id,
            type: request.type,
            target: request.target,
            detail: request.detail,
        })),
    };
}

function makeCharacterOverview(model: CharacterOverviewModel & Partial<CharacterRank>): CharacterOverview {
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

    @Method("getCharacterList")
    async getCharacterList(call: GrpcCall<GetCharacterListRequest>, callback: GrpcCallback<GetCharacterListReply>) {
        try {
            const result = await this.characterOverviewService.getCharacterList(call.request.accountId, call.request.worldId);
            callback(null, {
                characters: result.characters.map(makeCharacterOverview),
                slotCount: result.slotCount,
                hasSecondPassword: result.hasSecondPassword,
            });
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("checkCharacterName")
    async checkCharacterName(call: GrpcCall<CheckCharacterNameRequest>, callback: GrpcCallback<CheckCharacterNameReply>) {
        try {
            const result = await this.characterOverviewService.checkCharacterName(call.request.name, call.request.accountId, call.request.worldId);
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
            const result = await this.characterService.deleteCharacter(call.request.accountId, call.request.characterId, call.request.worldId);
            callback(null, {
                success: result.success,
                secondPasswordSet: result.secondPasswordSet ?? false,
                unsettled: result.unsettled ?? false,
            });
        } catch (err) {
            this.grpcError(err, callback);
        }
    }
}
