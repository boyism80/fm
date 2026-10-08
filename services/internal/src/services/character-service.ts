const MAX_NAME_LEN = 32;
const DEFAULT_MAP_ID = 10000;
const DEFAULT_SPAWN = 3;
import { EQUIP_SLOT, LOOK_SLOT } from "../constants/equipment-slots";
import { getDefaultKeyLayoutBindings } from "../constants/default-key-layout-bindings";
import { bindingsToJsonString, jsonStringToBindings } from "../grpc/key-layout-io";
import type { KeyLayoutBindingModel } from "../grpc/key-layout-io";
import { AccountSessionState, GuildErrorCode, PartyErrorCode } from "../protobuf/generated/fminternal/internal_service";
import type { CreateCharacterRequest } from "../protobuf/generated/fminternal/internal_service";
import { AppConfiguration } from "../config/app-configuration";
import { AccountRepository } from "../repos/account-repository";
import { CharacterOverviewRepository } from "../repos/character-overview-repository";
import { CharacterRepository } from "../repos/character-repository";
import type { CharacterModel } from "../repos/character-repository";
import { InventoryRepository } from "../repos/inventory-repository";
import type { InventoryModel } from "../repos/inventory-repository";
import { KeyLayoutRepository } from "../repos/key-layout-repository";
import { StorageRepository } from "../repos/storage-repository";
import type { StorageModel } from "../repos/storage-repository";
import { StorageItemRepository } from "../repos/storage-item-repository";
import { SkillRepository } from "../repos/skill-repository";
import type { SkillModel } from "../repos/skill-repository";
import { BuffRepository } from "../repos/buff-repository";
import type { BuffModel } from "../repos/buff-repository";
import { QuestRepository } from "../repos/quest-repository";
import type { QuestModel } from "../repos/quest-repository";
import { SavedLocationRepository } from "../repos/saved-location-repository";
import type { SavedLocationModel } from "../repos/saved-location-repository";
import { CharacterRecordRepository } from "../repos/character-record-repository";
import type { RecordModel } from "../repos/character-record-repository";
import { AccountRecordRepository } from "../repos/account-record-repository";
import { MonsterBookRepository } from "../repos/monster-book-repository";
import type { MonsterBookCardModel } from "../repos/monster-book-repository";
import { CharacterBuddyRepository } from "../repos/character-buddy-repository";
import { CharacterRealtimeStateRepository } from "../repos/character-realtime-state-repository";
import { SessionRepository } from "../repos/session-repository";
import { UnifiedRepository } from "../repos/unified-repository";
import { WzService } from "./wz-service";
import { DistributedLockService } from "./distributed-lock-service";
import { GuildService } from "./guild-service";
import { PartyService } from "./party-service";

type CharacterInput = {
    worldId: number;
    characterId: number;
    accountId: number;
    name: string;
    gender: number;
    skinColor: number;
    face: number;
    hair: number;
    level: number;
    classId: number;
    mapId: number;
    spawnPoint: number;
    role?: number;
    str?: number;
    dex?: number;
    intStat?: number;
    luk?: number;
    hp?: number;
    maxHp?: number;
    mp?: number;
    maxMp?: number;
    abilityPoint?: number;
    exp?: number;
    positionX?: number;
    positionY?: number;
    stance?: number;
    meso?: number;
    skillPoint?: number;
    population?: number;
    hpApUsed?: number;
    petHpItem?: number;
    petMpItem?: number;
    summonedPet?: number;
    slotLimits?: number[];
    monsterBookCover?: number;
    teleportStones?: number[];
    vipTeleportStones?: number[];
    hidden?: boolean;
};

export type CharacterRowModel = CharacterModel;

export type SaveCharacterEntry = {
    character: CharacterInput;
    baseLooks?: Record<string, number>;
    overlays?: Record<string, number>;
    inventory?: InventoryModel[];
    skills?: SkillModel[];
    buffs?: BuffModel[];
    quests?: QuestModel[];
    savedLocations?: SavedLocationModel[];
    records?: RecordModel[];
    accountRecords?: RecordModel[];
    monsterBook?: MonsterBookCardModel[];
    keyLayout?: KeyLayoutBindingModel[];
    storage?: { storage: StorageModel; items: InventoryModel[] };
};

export class CharacterService {
    private readonly repo: CharacterRepository;
    private readonly overviewRepo: CharacterOverviewRepository;
    private readonly accountRepo: AccountRepository;
    private readonly unifiedRepo: UnifiedRepository;
    private readonly inventoryRepo: InventoryRepository;
    private readonly skillRepo: SkillRepository;
    private readonly buffRepo: BuffRepository;
    private readonly questRepo: QuestRepository;
    private readonly savedLocationRepo: SavedLocationRepository;
    private readonly characterRecordRepo: CharacterRecordRepository;
    private readonly accountRecordRepo: AccountRecordRepository;
    private readonly monsterBookRepo: MonsterBookRepository;
    private readonly keyLayoutRepo: KeyLayoutRepository;
    private readonly storageRepo: StorageRepository;
    private readonly storageItemRepo: StorageItemRepository;
    private readonly buddyRepo: CharacterBuddyRepository;
    private readonly realtimeStateRepo: CharacterRealtimeStateRepository;
    private readonly sessionRepo: SessionRepository;
    private readonly app: AppConfiguration;
    private readonly wzService: WzService;
    private readonly distributedLockService: DistributedLockService;
    private readonly guildService: GuildService;
    private readonly partyService: PartyService;

    constructor(
        characterRepository: CharacterRepository,
        characterOverviewRepository: CharacterOverviewRepository,
        accountRepository: AccountRepository,
        unifiedRepository: UnifiedRepository,
        inventoryRepository: InventoryRepository,
        skillRepository: SkillRepository,
        buffRepository: BuffRepository,
        questRepository: QuestRepository,
        savedLocationRepository: SavedLocationRepository,
        characterRecordRepository: CharacterRecordRepository,
        accountRecordRepository: AccountRecordRepository,
        monsterBookRepository: MonsterBookRepository,
        keyLayoutRepository: KeyLayoutRepository,
        storageRepository: StorageRepository,
        storageItemRepository: StorageItemRepository,
        characterBuddyRepository: CharacterBuddyRepository,
        characterRealtimeStateRepository: CharacterRealtimeStateRepository,
        sessionRepository: SessionRepository,
        appConfiguration: AppConfiguration,
        wzService: WzService,
        distributedLockService: DistributedLockService,
        guildService: GuildService,
        partyService: PartyService
    ) {
        this.repo = characterRepository;
        this.overviewRepo = characterOverviewRepository;
        this.accountRepo = accountRepository;
        this.unifiedRepo = unifiedRepository;
        this.inventoryRepo = inventoryRepository;
        this.skillRepo = skillRepository;
        this.buffRepo = buffRepository;
        this.questRepo = questRepository;
        this.savedLocationRepo = savedLocationRepository;
        this.characterRecordRepo = characterRecordRepository;
        this.accountRecordRepo = accountRecordRepository;
        this.monsterBookRepo = monsterBookRepository;
        this.keyLayoutRepo = keyLayoutRepository;
        this.storageRepo = storageRepository;
        this.storageItemRepo = storageItemRepository;
        this.buddyRepo = characterBuddyRepository;
        this.realtimeStateRepo = characterRealtimeStateRepository;
        this.sessionRepo = sessionRepository;
        this.app = appConfiguration;
        this.wzService = wzService;
        this.distributedLockService = distributedLockService;
        this.guildService = guildService;
        this.partyService = partyService;
    }

    private worldId() {
        return this.app.app.world_id;
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
        if (!Number.isInteger(characterId) || characterId <= 0 || characterId > 0xffffffff) {
            const err = new Error("character_id must be a positive uint32") as Error & { code?: string };
            err.code = "INVALID_CHARACTER_ID";
            throw err;
        }
    }

    private assertAccountId(accountId: number) {
        if (!Number.isInteger(accountId) || accountId <= 0 || accountId > 0xffffffff) {
            const err = new Error("account_id must be a positive uint32") as Error & { code?: string };
            err.code = "INVALID_PAYLOAD";
            throw err;
        }
    }

    private validateCharacter(p: { name?: string }) {
        if (typeof p.name !== "string" || p.name.length > MAX_NAME_LEN) {
            const err = new Error(`name must be a string of length <= ${MAX_NAME_LEN}`) as Error & { code?: string };
            err.code = "INVALID_PAYLOAD";
            throw err;
        }
    }

    async getCharacter(worldId: number, characterId: number): Promise<CharacterRowModel | null> {
        this.assertWorld(worldId);
        this.assertCharacterId(characterId);
        const row = await this.repo.get(worldId, characterId);
        return row as CharacterRowModel | null;
    }

    async getKeyLayoutBindings(worldId: number, characterId: number) {
        this.assertWorld(worldId);
        this.assertCharacterId(characterId);
        const keyLayout = await this.keyLayoutRepo.get(worldId, characterId);
        const keyLayoutJson = keyLayout?.keyLayoutJson;
        const keyLayoutJsonString = typeof keyLayoutJson === "string"
            ? keyLayoutJson
            : JSON.stringify(keyLayoutJson ?? {});
        return jsonStringToBindings(keyLayoutJsonString);
    }

    async loadStorage(worldId: number, accountId: number) {
        this.assertWorld(worldId);
        this.assertAccountId(accountId);
        const storage = await this.storageRepo.get(worldId, accountId);
        const items = await this.storageItemRepo.getAll(worldId, String(accountId));
        return {
            storage: storage ?? { accountId, worldId, slots: 4, meso: 0 },
            items: [...items.values()],
        };
    }

    async saveCharacter(character: CharacterInput, baseLooks?: Record<string, number>, overlays?: Record<string, number>) {
        return this.saveCharacters([{ character, baseLooks, overlays }]);
    }

    async saveCharacters(entries: SaveCharacterEntry[]) {
        if (!entries || entries.length === 0) {
            return;
        }

        const byWorld = new Map<number, SaveCharacterEntry[]>();
        for (const { character, baseLooks, overlays, inventory, skills, buffs, quests, savedLocations, records, accountRecords, monsterBook, keyLayout, storage } of entries) {
            this.assertWorld(character.worldId);
            this.assertCharacterId(character.characterId);
            this.assertAccountId(character.accountId);
            this.validateCharacter(character);
            if (!character.mapId || character.mapId <= 0) {
                continue;
            }
            const wid = character.worldId;
            if (!byWorld.has(wid)) {
                byWorld.set(wid, []);
            }
            byWorld.get(wid)?.push({ character, baseLooks, overlays, inventory, skills, buffs, quests, savedLocations, records, accountRecords, monsterBook, keyLayout, storage });
        }

        for (const [worldId, group] of byWorld) {
            const lockKeys = group
                .map(({ character }) => character.characterId)
                .sort((a, b) => a - b)
                .map((characterId) => `character:${characterId}`);
            await using _characterLocks = await this.distributedLockService.acquireWorldDataLocks(worldId, lockKeys);

            const models = group.map(({ character }) => character);
            await this.repo.setAll(worldId, models);

            for (const { character, baseLooks, overlays, inventory, skills, buffs, quests, savedLocations, records, accountRecords, monsterBook, keyLayout, storage } of group) {
                if (!character.accountId) {
                    continue;
                }
                const overview = {
                    characterId: character.characterId,
                    accountId: character.accountId,
                    worldId: character.worldId,
                    name: character.name,
                    gender: character.gender,
                    skinColor: character.skinColor,
                    face: character.face,
                    hair: character.hair,
                    level: character.level,
                    classId: character.classId,
                    mapId: character.mapId,
                    spawnPoint: character.spawnPoint,
                    rank: 0,
                    rankDiff: 0,
                    classRank: 0,
                    classRankDiff: 0,
                    baseLooks: baseLooks ?? {},
                    overlays: overlays ?? {},
                };
                await this.overviewRepo.set(character.worldId, overview);

                if (inventory !== undefined) {
                    await this.inventoryRepo.replaceBySnapshot(character.worldId, String(character.characterId), inventory);
                }
                if (skills !== undefined) {
                    await this.skillRepo.replaceBySnapshot(character.worldId, String(character.characterId), skills.map((m) => ({ ...m, characterId: character.characterId })));
                }
                if (buffs !== undefined) {
                    await this.buffRepo.replaceBySnapshot(character.worldId, String(character.characterId), buffs.map((m) => ({ ...m, characterId: character.characterId })));
                }
                if (quests !== undefined) {
                    await this.questRepo.replaceBySnapshot(character.worldId, String(character.characterId), quests.map((m) => ({ ...m, characterId: character.characterId })));
                }
                if (savedLocations !== undefined) {
                    await this.savedLocationRepo.replaceBySnapshot(character.worldId, String(character.characterId), savedLocations.map((m) => ({ ...m, characterId: character.characterId })));
                }
                if (records !== undefined) {
                    await this.characterRecordRepo.replaceBySnapshot(character.worldId, String(character.characterId), records.map((m) => ({ ...m, ownerId: character.characterId })));
                }
                if (accountRecords !== undefined) {
                    await this.accountRecordRepo.replaceBySnapshot(character.worldId, String(character.accountId), accountRecords.map((m) => ({ ...m, ownerId: character.accountId })));
                }
                if (monsterBook !== undefined) {
                    await this.monsterBookRepo.replaceBySnapshot(character.worldId, String(character.characterId), monsterBook.map((m) => ({ ...m, characterId: character.characterId })));
                }
                if (keyLayout !== undefined) {
                    await this.keyLayoutRepo.set(character.worldId, {
                        characterId: character.characterId,
                        worldId: character.worldId,
                        keyLayoutJson: bindingsToJsonString(keyLayout),
                    });
                }
                if (storage !== undefined) {
                    await this.storageRepo.set(character.worldId, { ...storage.storage, accountId: character.accountId, worldId: character.worldId });
                    await this.storageItemRepo.replaceBySnapshot(character.worldId, String(character.accountId), storage.items);
                }
            }
        }
    }

    async createCharacter(
        accountId: number,
        worldId: number | undefined,
        params: Pick<CreateCharacterRequest, "name" | "face" | "hair" | "skinColor" | "topItemId" | "bottomItemId" | "shoesItemId" | "weaponItemId">
    ) {
        const { name, face, hair, skinColor, topItemId, bottomItemId, shoesItemId, weaponItemId } = params;
        const wid = worldId ?? this.worldId();
        this.assertWorld(wid);
        this.assertAccountId(accountId);
        this.validateCharacter({ name });
        const session = await this.sessionRepo.getAccountSession(wid, accountId);
        if (session?.state !== AccountSessionState.ACCOUNT_SESSION_STATE_LOGIN) {
            return { success: false, errorMsg: "로그인 상태가 아닙니다." };
        }
        if (this.wzService.isForbiddenName(name)) {
            return { success: false, errorMsg: "사용할 수 없는 이름입니다." };
        }

        await using _createCharacterLock = await this.distributedLockService.acquireWorldDataLock(wid, `create_character:a${accountId}`);

        const overviewMap = await this.overviewRepo.getAll(wid, String(accountId));
        const account = await this.accountRepo.get(wid, accountId);
        const slotCount = account?.characterSlotCount ?? 6;
        if (overviewMap.size >= slotCount) {
            return { success: false, errorMsg: "캐릭터 슬롯이 부족합니다." };
        }
        const gender = account?.gender ?? 0;
        if (this.wzService.canMakeCharacter(gender, [face, hair, topItemId, bottomItemId, shoesItemId, weaponItemId]) === false) {
            return { success: false, errorMsg: "선택할 수 없는 외형입니다." };
        }

        const characterId = await this.unifiedRepo.reserveCharacterName(name, accountId, wid);
        if (characterId === null) {
            return { success: false, errorMsg: "이미 사용 중인 이름입니다." };
        }
        const role = account?.role ?? 0;

        const character = {
            characterId,
            accountId,
            worldId: wid,
            name,
            gender,
            skinColor,
            face,
            hair,
            level: 1,
            classId: 0,
            role,
            str: 12,
            dex: 5,
            intStat: 4,
            luk: 4,
            hp: 50,
            maxHp: 50,
            mp: 5,
            maxMp: 5,
            abilityPoint: 0,
            exp: 0,
            mapId: DEFAULT_MAP_ID,
            spawnPoint: DEFAULT_SPAWN,
            positionX: 0,
            positionY: 0,
            stance: 0,
            meso: 0,
            skillPoint: 0,
            population: 0,
            hpApUsed: 0,
            hidden: false,
        };
        const equips = [
            { itemId: topItemId, slot: EQUIP_SLOT.TOP, lookSlot: LOOK_SLOT.TOP },
            { itemId: bottomItemId, slot: EQUIP_SLOT.BOTTOM, lookSlot: LOOK_SLOT.BOTTOM },
            { itemId: shoesItemId, slot: EQUIP_SLOT.SHOES, lookSlot: LOOK_SLOT.SHOES },
            { itemId: weaponItemId, slot: EQUIP_SLOT.WEAPON, lookSlot: LOOK_SLOT.WEAPON },
        ].filter((e) => e.itemId > 0);

        const baseLooks: Record<string, number> = {};
        for (const e of equips) {
            baseLooks[e.lookSlot] = e.itemId;
        }

        const overview = {
            characterId,
            accountId,
            worldId: wid,
            name,
            gender,
            skinColor,
            face,
            hair,
            level: 1,
            classId: 0,
            mapId: DEFAULT_MAP_ID,
            spawnPoint: DEFAULT_SPAWN,
            rank: 0,
            rankDiff: 0,
            classRank: 0,
            classRankDiff: 0,
            baseLooks,
            overlays: {},
        };

        try {
            await this.repo.set(wid, character);
            await this.keyLayoutRepo.set(wid, {
                characterId,
                worldId: wid,
                keyLayoutJson: bindingsToJsonString(getDefaultKeyLayoutBindings()),
            });

            if (equips.length) {
                const enhances = await Promise.all(equips.map((e) => this.wzService.getEnhanceChance(e.itemId)));
                await this.inventoryRepo.replaceBySnapshot(wid, String(characterId), []);
                await this.inventoryRepo.setAll(
                    wid,
                    equips.map((e, idx) => ({
                        uniqueId: null,
                        ownerId: characterId,
                        inventoryType: 1,
                        itemId: e.itemId,
                        slot: e.slot,
                        count: 1,
                        expiration: null,
                        enhanceChance: Number.isFinite(enhances[idx]) ? Math.max(0, enhances[idx] as number) : 0,
                        enhanceCount: 0,
                        flag: 0,
                        skillBonus: 0,
                        ownerName: null,
                    }))
                );
            }

            await this.overviewRepo.set(wid, overview);
            await this.unifiedRepo.confirmCharacterName(name);
        } catch (err) {
            // Remove any data written for this characterId so the reserved slot
            // stays clean and the name reservation can be reclaimed later.
            await Promise.allSettled([
                this.repo.delete({ worldId: wid, characterId, accountId }),
                this.keyLayoutRepo.delete({ worldId: wid, characterId }),
                this.inventoryRepo.replaceBySnapshot(wid, String(characterId), []),
                this.overviewRepo.delete({ worldId: wid, accountId, characterId }),
                this.unifiedRepo.releaseCharacterNameReservation(accountId),
            ]);
            throw err;
        }

        return { success: true, errorMsg: "", character: overview };
    }

    async deleteCharacter(accountId: number, characterId: number) {
        const worldId = this.worldId();
        this.assertWorld(worldId);
        this.assertAccountId(accountId);
        this.assertCharacterId(characterId);

        const session = await this.sessionRepo.getAccountSession(worldId, accountId);
        if (session?.state !== AccountSessionState.ACCOUNT_SESSION_STATE_LOGIN) {
            return { success: false };
        }

        await using _characterLock = await this.distributedLockService.acquireWorldDataLock(
            worldId,
            `character:${characterId}`,
        );

        const character = await this.repo.get(worldId, characterId);
        if (!character || character.accountId !== accountId) {
            return { success: false };
        }

        const guildLeft = await this.guildService.leaveGuild(worldId, characterId);
        if (!guildLeft.ok && guildLeft.code !== GuildErrorCode.GUILD_ERROR_NOT_IN_GUILD) {
            return { success: false };
        }
        const partyLeft = await this.partyService.leaveParty(worldId, characterId);
        if (!partyLeft.ok && partyLeft.code !== PartyErrorCode.NOT_IN_PARTY) {
            return { success: false };
        }

        await this.inventoryRepo.replaceBySnapshot(worldId, String(characterId), []);
        await this.skillRepo.replaceBySnapshot(worldId, String(characterId), []);
        await this.buffRepo.replaceBySnapshot(worldId, String(characterId), []);
        await this.questRepo.replaceBySnapshot(worldId, String(characterId), []);
        await this.savedLocationRepo.replaceBySnapshot(worldId, String(characterId), []);
        await this.characterRecordRepo.replaceBySnapshot(worldId, String(characterId), []);
        await this.monsterBookRepo.replaceBySnapshot(worldId, String(characterId), []);
        await this.buddyRepo.deleteAllForOwner(worldId, characterId);
        await this.buddyRepo.deleteAllReferencingBuddy(worldId, characterId);
        await this.realtimeStateRepo.delete({ worldId, characterId });
        await this.keyLayoutRepo.delete({ worldId, characterId });

        const characterDeleteModel = { worldId, characterId, accountId };
        const ok = await this.repo.delete(characterDeleteModel);
        if (!ok) {
            return { success: false };
        }

        await this.unifiedRepo.deleteCharacterName(characterId);
        await this.overviewRepo.delete({ worldId, accountId, characterId });
        return { success: true };
    }
}
