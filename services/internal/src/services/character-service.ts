const DEFAULT_MAP_ID = 10000;
const DEFAULT_SPAWN = 3;
import { EQUIP_SLOT, LOOK_SLOT } from "../constants/equipment-slots";
import { getDefaultKeyLayoutBindings } from "../constants/default-key-layout-bindings";
import { bindingsToJsonString, jsonStringToBindings, makeKeyLayoutProtoList } from "../grpc/key-layout-io";
import { CHARACTER_MODEL, CHARACTER_PROTO } from "../grpc/character-proto";
import { BUDDY_ENTRY, BUDDY_LIST_ENTRY } from "../grpc/buddy-proto";
import { INVENTORY_MODEL, INVENTORY_PROTO } from "../grpc/inventory-proto";
import { SKILL_MODEL, SKILL_PROTO } from "../grpc/skill-proto";
import { BUFF_MODEL, BUFF_PROTO } from "../grpc/buff-proto";
import { QUEST_MODEL, QUEST_PROTO } from "../grpc/quest-proto";
import { SAVED_LOCATION_MODEL, SAVED_LOCATION_PROTO } from "../grpc/saved-location-proto";
import { grpcMapper } from "../grpc/mappers";
import type { KeyLayoutBindingModel } from "../grpc/key-layout-io";
import {
    AccountSessionState,
    GuildErrorCode,
    PartyErrorCode,
} from "@fm/shared/protobuf/generated/fminternal/internal_service";
import type {
    BuddyEntry,
    Buff,
    Character,
    CreateCharacterRequest,
    EnterGameReply,
    Inventory,
    Quest,
    SavedLocation,
    Skill,
} from "@fm/shared/protobuf/generated/fminternal/internal_service";
import { AppConfiguration } from "@fm/shared/config/app-configuration";
import { AccountRepository } from "@fm/shared/repos/account-repository";
import { AccountWorldRepository } from "@fm/shared/repos/account-world-repository";
import { CharacterOverviewRepository } from "@fm/shared/repos/character-overview-repository";
import type { CharacterOverviewModel } from "@fm/shared/repos/character-overview-repository";
import { CharacterRepository } from "@fm/shared/repos/character-repository";
import type { CharacterModel } from "@fm/shared/repos/character-repository";
import { InventoryRepository } from "@fm/shared/repos/inventory-repository";
import type { InventoryModel } from "@fm/shared/repos/inventory-repository";
import { KeyLayoutRepository } from "@fm/shared/repos/key-layout-repository";
import { SkillMacroRepository } from "@fm/shared/repos/skill-macro-repository";
import type { SkillMacroModel } from "@fm/shared/repos/skill-macro-repository";
import { MountRepository } from "@fm/shared/repos/mount-repository";
import type { MountModel } from "@fm/shared/repos/mount-repository";
import { UnsettledRequestRepository } from "@fm/shared/repos/unsettled-request-repository";
import type { UnsettledRequestModel } from "@fm/shared/repos/unsettled-request-repository";
import { StorageRepository } from "@fm/shared/repos/storage-repository";
import type { StorageModel } from "@fm/shared/repos/storage-repository";
import { StorageItemRepository } from "@fm/shared/repos/storage-item-repository";
import { SkillRepository } from "@fm/shared/repos/skill-repository";
import type { SkillModel } from "@fm/shared/repos/skill-repository";
import { BuffRepository } from "@fm/shared/repos/buff-repository";
import type { BuffModel } from "@fm/shared/repos/buff-repository";
import { QuestRepository } from "@fm/shared/repos/quest-repository";
import type { QuestModel } from "@fm/shared/repos/quest-repository";
import { SavedLocationRepository } from "@fm/shared/repos/saved-location-repository";
import type { SavedLocationModel } from "@fm/shared/repos/saved-location-repository";
import { CharacterRecordRepository } from "@fm/shared/repos/character-record-repository";
import type { RecordModel } from "@fm/shared/repos/character-record-repository";
import { AccountRecordRepository } from "@fm/shared/repos/account-record-repository";
import { MonsterBookRepository } from "@fm/shared/repos/monster-book-repository";
import type { MonsterBookCardModel } from "@fm/shared/repos/monster-book-repository";
import { CharacterBuddyRepository } from "@fm/shared/repos/character-buddy-repository";
import { CharacterRealtimeStateRepository } from "@fm/shared/repos/character-realtime-state-repository";
import { SessionRepository } from "@fm/shared/repos/session-repository";
import { ShopRepository } from "@fm/shared/repos/shop-repository";
import { UnifiedRepository } from "@fm/shared/repos/unified-repository";
import { WzService } from "./wz-service";
import { DistributedLockService } from "./distributed-lock-service";
import { GuildService } from "./guild-service";
import type { FamilyService } from "./family-service";
import type { BuddyService, BuddyListEntry } from "./buddy-service";
import type { MarriageService } from "./marriage-service";
import type { CashShopService } from "./cash-shop-service";
import type { SocialService } from "./social-service";
import { PartyService } from "./party-service";
import type { InternalContext } from "@fm/shared/context/internal-context";

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
    summonedPet?: string;
    slotLimits?: number[];
    monsterBookCover?: number;
    teleportStones?: number[];
    vipTeleportStones?: number[];
    hidden?: boolean;
    savedAtUnixMs?: number;
    expReachedUnixMs?: number;
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
    skillMacros?: SkillMacroModel[];
    mount?: Pick<MountModel, "level" | "exp" | "fatigue">;
    storage?: { storage: StorageModel; items: InventoryModel[] };
    claimedShops?: number[];
    unsettledRequests?: Omit<UnsettledRequestModel, "characterId">[];
    cashWeapon?: number;
};

function characterOverview(
    character: CharacterInput,
    baseLooks: Record<string, number>,
    overlays: Record<string, number>,
    cashWeapon: number,
): CharacterOverviewModel {
    return {
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
        baseLooks,
        overlays,
        cashWeapon,
        str: character.str ?? 0,
        dex: character.dex ?? 0,
        intStat: character.intStat ?? 0,
        luk: character.luk ?? 0,
        hp: character.hp ?? 0,
        maxHp: character.maxHp ?? 0,
        mp: character.mp ?? 0,
        maxMp: character.maxMp ?? 0,
        abilityPoint: character.abilityPoint ?? 0,
        skillPoint: character.skillPoint ?? 0,
        exp: character.exp ?? 0,
        population: character.population ?? 0,
    };
}

export class CharacterService {
    private readonly repo: CharacterRepository;
    private readonly overviewRepo: CharacterOverviewRepository;
    private readonly accountRepo: AccountRepository;
    private readonly accountWorldRepo: AccountWorldRepository;
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
    private readonly skillMacroRepo: SkillMacroRepository;
    private readonly mountRepo: MountRepository;
    private readonly unsettledRequestRepo: UnsettledRequestRepository;
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
    private readonly shopRepo: ShopRepository;
    private readonly familyService: FamilyService;
    private readonly buddyService: BuddyService;
    private readonly marriageService: MarriageService;
    private readonly cashShopService: CashShopService;
    private readonly socialService: SocialService;
    private readonly ctx: InternalContext;

    constructor(
        characterRepository: CharacterRepository,
        characterOverviewRepository: CharacterOverviewRepository,
        accountRepository: AccountRepository,
        accountWorldRepository: AccountWorldRepository,
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
        skillMacroRepository: SkillMacroRepository,
        mountRepository: MountRepository,
        unsettledRequestRepository: UnsettledRequestRepository,
        storageRepository: StorageRepository,
        storageItemRepository: StorageItemRepository,
        characterBuddyRepository: CharacterBuddyRepository,
        characterRealtimeStateRepository: CharacterRealtimeStateRepository,
        sessionRepository: SessionRepository,
        appConfiguration: AppConfiguration,
        wzService: WzService,
        distributedLockService: DistributedLockService,
        guildService: GuildService,
        partyService: PartyService,
        shopRepository: ShopRepository,
        familyService: FamilyService,
        buddyService: BuddyService,
        marriageService: MarriageService,
        cashShopService: CashShopService,
        socialService: SocialService,
        internalContext: InternalContext,
    ) {
        this.repo = characterRepository;
        this.overviewRepo = characterOverviewRepository;
        this.accountRepo = accountRepository;
        this.accountWorldRepo = accountWorldRepository;
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
        this.skillMacroRepo = skillMacroRepository;
        this.mountRepo = mountRepository;
        this.unsettledRequestRepo = unsettledRequestRepository;
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
        this.shopRepo = shopRepository;
        this.familyService = familyService;
        this.buddyService = buddyService;
        this.marriageService = marriageService;
        this.cashShopService = cashShopService;
        this.socialService = socialService;
        this.ctx = internalContext;
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

    // PostgreSQL text cannot hold NUL; one rejected write would stop the batch halfway.
    private assertNoNul(value: unknown) {
        if (typeof value === "string") {
            if (value.includes("\u0000")) {
                const err = new Error("text must not contain NUL") as Error & { code?: string };
                err.code = "INVALID_PAYLOAD";
                throw err;
            }
            return;
        }
        if (value === null || typeof value !== "object") {
            return;
        }
        for (const child of Object.values(value)) {
            this.assertNoNul(child);
        }
    }

    private validateCharacter(p: { name?: string }) {
        if (typeof p.name !== "string") {
            const err = new Error("name must be a string") as Error & { code?: string };
            err.code = "INVALID_PAYLOAD";
            throw err;
        }
    }

    async getCharacter(worldId: number, characterId: number): Promise<CharacterRowModel | null> {
        this.assertWorld(worldId);
        this.assertCharacterId(characterId);
        return (await this.repo.get(worldId, characterId)) as CharacterRowModel | null;
    }

    async getKeyLayoutBindings(worldId: number, characterId: number) {
        this.assertWorld(worldId);
        this.assertCharacterId(characterId);
        const keyLayout = await this.keyLayoutRepo.get(worldId, characterId);
        const keyLayoutJson = keyLayout?.keyLayoutJson;
        const keyLayoutJsonString =
            typeof keyLayoutJson === "string" ? keyLayoutJson : JSON.stringify(keyLayoutJson ?? {});
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

    async saveCharacter(
        character: CharacterInput,
        baseLooks?: Record<string, number>,
        overlays?: Record<string, number>,
        cashWeapon?: number,
    ) {
        return this.saveCharacters([{ character, baseLooks, overlays, cashWeapon }]);
    }

    async load(worldId: number, row: CharacterModel): Promise<EnterGameReply> {
        const characterId = row.characterId;
        const [
            inventoryList,
            skillList,
            buffList,
            questList,
            savedLocationList,
            keyLayoutBindings,
            skillMacros,
            buddyPack,
            marriage,
            cashWishlist,
            rings,
            monsterBook,
            records,
            accountRecords,
            mount,
            storage,
            unsettledRequests,
        ] = await Promise.all([
            this.inventoryRepo.getAll(worldId, String(characterId)).then((map) => [...map.values()]),
            this.skillRepo.getAll(worldId, String(characterId)).then((map) => [...map.values()]),
            this.buffRepo.getAll(worldId, String(characterId)).then((map) => [...map.values()]),
            this.questRepo.getAll(worldId, String(characterId)).then((map) => [...map.values()]),
            this.savedLocationRepo.getAll(worldId, String(characterId)).then((map) => [...map.values()]),
            this.getKeyLayoutBindings(worldId, characterId),
            this.skillMacroRepo
                .getAll(worldId, String(characterId))
                .then((map) => [...map.values()].sort((a, b) => a.position - b.position)),
            this.buddyService.getAll(worldId, characterId),
            this.marriageService.loadMarriage(worldId, characterId),
            this.cashShopService.getWishlist(worldId, characterId),
            this.cashShopService.getRings(worldId, characterId),
            this.monsterBookRepo.getAll(worldId, String(characterId)).then((map) => [...map.values()]),
            this.characterRecordRepo.getAll(worldId, String(characterId)).then((map) => [...map.values()]),
            this.accountRecordRepo.getAll(worldId, String(row.accountId)).then((map) => [...map.values()]),
            this.mountRepo.get(worldId, characterId),
            this.loadStorage(worldId, row.accountId),
            this.unsettledRequestRepo.getAll(worldId, String(characterId)).then((map) => [...map.values()]),
        ]);
        const realtime = await this.realtimeStateRepo.get(worldId, characterId);
        const family = await this.familyService.load(worldId, row);
        const notes = await this.socialService.listNotes(worldId, characterId);
        return {
            family,
            notes,
            found: true,
            character: grpcMapper.map<CharacterModel, Character>(row, CHARACTER_MODEL, CHARACTER_PROTO),
            inventory: inventoryList.map((inventory) =>
                grpcMapper.map<InventoryModel, Inventory>(inventory, INVENTORY_MODEL, INVENTORY_PROTO),
            ),
            skills: skillList.map((skill) => grpcMapper.map<SkillModel, Skill>(skill, SKILL_MODEL, SKILL_PROTO)),
            buffs: buffList.map((buff) => grpcMapper.map<BuffModel, Buff>(buff, BUFF_MODEL, BUFF_PROTO)),
            quests: questList.map((quest) => grpcMapper.map<QuestModel, Quest>(quest, QUEST_MODEL, QUEST_PROTO)),
            savedLocations: savedLocationList.map((loc) =>
                grpcMapper.map<SavedLocationModel, SavedLocation>(loc, SAVED_LOCATION_MODEL, SAVED_LOCATION_PROTO),
            ),
            keyLayout: makeKeyLayoutProtoList(keyLayoutBindings),
            skillMacros,
            partyId: realtime?.partyId != null ? realtime.partyId : undefined,
            guildId: realtime?.guildId != null ? realtime.guildId : undefined,
            buddies: buddyPack.buddies.map((buddy) =>
                grpcMapper.map<BuddyListEntry, BuddyEntry>(buddy, BUDDY_LIST_ENTRY, BUDDY_ENTRY),
            ),
            buddyCapacity: buddyPack.capacity >>> 0,
            debuffs: [],
            marriage: marriage.marriage,
            cashWishlist: cashWishlist.filter((sn) => sn !== 0),
            rings,
            monsterBook: monsterBook.map((card) => ({ cardId: card.cardId >>> 0, count: card.count >>> 0 })),
            records: records.map(({ key, value, text, expiresAtUnixMs, recordedAtUnixMs }) => ({
                key,
                value,
                text,
                expiresAtUnixMs,
                recordedAtUnixMs,
            })),
            account: {
                storage: {
                    accountId: storage.storage.accountId,
                    worldId: storage.storage.worldId,
                    slots: storage.storage.slots,
                    meso: storage.storage.meso,
                    items: storage.items.map((item) =>
                        grpcMapper.map<InventoryModel, Inventory>(item, INVENTORY_MODEL, INVENTORY_PROTO),
                    ),
                },
                records: accountRecords.map(({ key, value, text, expiresAtUnixMs, recordedAtUnixMs }) => ({
                    key,
                    value,
                    text,
                    expiresAtUnixMs,
                    recordedAtUnixMs,
                })),
            },
            mount: {
                level: mount?.level ?? 1,
                exp: mount?.exp ?? 0,
                fatigue: mount?.fatigue ?? 0,
            },
            unsettledRequests: unsettledRequests.map(({ id, type, target, detail }) => ({ id, type, target, detail })),
        };
    }

    async saveCharacters(entries: SaveCharacterEntry[]) {
        if (!entries || entries.length === 0) {
            return;
        }
        const subject = JSON.stringify(
            entries.map(({ character }) => ({
                worldId: character.worldId,
                characterId: character.characterId,
                accountId: character.accountId,
            })),
        );
        return this.ctx.withDatabaseOperation("saveCharacters", subject, () => this.persistCharacters(entries));
    }

    private async persistCharacters(entries: SaveCharacterEntry[]) {
        for (const entry of entries) {
            this.assertNoNul(entry);
        }

        const byWorld = new Map<number, SaveCharacterEntry[]>();
        for (const {
            character,
            baseLooks,
            overlays,
            inventory,
            skills,
            buffs,
            quests,
            savedLocations,
            records,
            accountRecords,
            monsterBook,
            keyLayout,
            skillMacros,
            mount,
            storage,
            claimedShops,
            unsettledRequests,
            cashWeapon,
        } of entries) {
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
            byWorld.get(wid)?.push({
                character,
                baseLooks,
                overlays,
                inventory,
                skills,
                buffs,
                quests,
                savedLocations,
                records,
                accountRecords,
                monsterBook,
                keyLayout,
                skillMacros,
                mount,
                storage,
                claimedShops,
                unsettledRequests,
                cashWeapon,
            });
        }

        for (const [worldId, group] of byWorld) {
            const lockKeys = group
                .map(({ character }) => character.characterId)
                .sort((a, b) => a - b)
                .map((characterId) => `character:${characterId}`);
            await using _characterLocks = await this.distributedLockService.acquireWorldDataLocks(worldId, lockKeys);

            const saved = await this.repo.getMany(
                worldId,
                group.map(({ character }) => character.characterId),
            );
            const latest = group.filter(
                ({ character }) =>
                    (character.savedAtUnixMs ?? 0) >= (saved.get(character.characterId)?.savedAtUnixMs ?? 0),
            );
            if (latest.length === 0) {
                continue;
            }

            const models = latest.map(({ character }) => character);
            await this.repo.setAll(worldId, models);
            for (const { character } of latest) {
                await this.familyService.progress(worldId, character);
            }

            for (const {
                character,
                baseLooks,
                overlays,
                inventory,
                skills,
                buffs,
                quests,
                savedLocations,
                records,
                accountRecords,
                monsterBook,
                keyLayout,
                skillMacros,
                mount,
                storage,
                claimedShops,
                unsettledRequests,
                cashWeapon,
            } of latest) {
                if (!character.accountId) {
                    continue;
                }
                await this.overviewRepo.set(
                    character.worldId,
                    characterOverview(character, baseLooks ?? {}, overlays ?? {}, cashWeapon ?? 0),
                );

                if (inventory !== undefined) {
                    await this.inventoryRepo.replaceBySnapshot(
                        character.worldId,
                        String(character.characterId),
                        inventory,
                    );
                }
                if (skills !== undefined) {
                    await this.skillRepo.replaceBySnapshot(
                        character.worldId,
                        String(character.characterId),
                        skills.map((m) => ({ ...m, characterId: character.characterId })),
                    );
                }
                if (buffs !== undefined) {
                    await this.buffRepo.replaceBySnapshot(
                        character.worldId,
                        String(character.characterId),
                        buffs.map((m) => ({ ...m, characterId: character.characterId })),
                    );
                }
                if (quests !== undefined) {
                    await this.questRepo.replaceBySnapshot(
                        character.worldId,
                        String(character.characterId),
                        quests.map((m) => ({ ...m, characterId: character.characterId })),
                    );
                }
                if (savedLocations !== undefined) {
                    await this.savedLocationRepo.replaceBySnapshot(
                        character.worldId,
                        String(character.characterId),
                        savedLocations.map((m) => ({ ...m, characterId: character.characterId })),
                    );
                }
                if (records !== undefined) {
                    await this.characterRecordRepo.replaceBySnapshot(
                        character.worldId,
                        String(character.characterId),
                        records.map((m) => ({ ...m, ownerId: character.characterId })),
                    );
                }
                if (accountRecords !== undefined) {
                    await this.accountRecordRepo.replaceBySnapshot(
                        character.worldId,
                        String(character.accountId),
                        accountRecords.map((m) => ({ ...m, ownerId: character.accountId })),
                    );
                }
                if (monsterBook !== undefined) {
                    await this.monsterBookRepo.replaceBySnapshot(
                        character.worldId,
                        String(character.characterId),
                        monsterBook.map((m) => ({ ...m, characterId: character.characterId })),
                    );
                }
                if (keyLayout !== undefined) {
                    await this.keyLayoutRepo.set(character.worldId, {
                        characterId: character.characterId,
                        worldId: character.worldId,
                        keyLayoutJson: bindingsToJsonString(keyLayout),
                    });
                }
                if (skillMacros !== undefined) {
                    await this.skillMacroRepo.replaceBySnapshot(
                        character.worldId,
                        String(character.characterId),
                        skillMacros.map((macro) => ({ ...macro, characterId: character.characterId })),
                    );
                }
                if (mount !== undefined) {
                    await this.mountRepo.set(character.worldId, {
                        ...mount,
                        characterId: character.characterId,
                        worldId: character.worldId,
                    });
                }
                if (storage !== undefined) {
                    await this.storageRepo.set(character.worldId, {
                        ...storage.storage,
                        accountId: character.accountId,
                        worldId: character.worldId,
                    });
                    await this.storageItemRepo.replaceBySnapshot(
                        character.worldId,
                        String(character.accountId),
                        storage.items,
                    );
                }
                if (unsettledRequests !== undefined) {
                    await this.unsettledRequestRepo.replaceBySnapshot(
                        character.worldId,
                        String(character.characterId),
                        unsettledRequests.map((m) => ({ ...m, characterId: character.characterId })),
                    );
                }
                if (claimedShops !== undefined && claimedShops.length > 0) {
                    await this.shopRepo.claim(
                        character.worldId,
                        character.accountId,
                        character.characterId,
                        claimedShops,
                    );
                }
            }
        }
    }

    async createCharacter(
        accountId: number,
        worldId: number,
        params: Pick<
            CreateCharacterRequest,
            "name" | "face" | "hair" | "skinColor" | "topItemId" | "bottomItemId" | "shoesItemId" | "weaponItemId"
        >,
    ) {
        const { name, face, hair, skinColor, topItemId, bottomItemId, shoesItemId, weaponItemId } = params;
        this.assertWorld(worldId);
        this.assertAccountId(accountId);
        this.validateCharacter({ name });
        const session = await this.sessionRepo.getAccountSession(accountId);
        if (session?.state !== AccountSessionState.ACCOUNT_SESSION_STATE_LOGIN) {
            return { success: false, errorMsg: "로그인 상태가 아닙니다." };
        }
        if (this.wzService.isValidName(name) === false) {
            return { success: false, errorMsg: "사용할 수 없는 이름입니다." };
        }

        await using _createCharacterLock = await this.distributedLockService.acquireUnifiedLock(
            `create_character:a${accountId}`,
        );

        const overviewMap = await this.overviewRepo.getAll(worldId, String(accountId));
        const account = await this.accountRepo.get(accountId);
        const slotCount = await this.accountWorldRepo.findCharacterSlotCount(worldId, accountId);
        if (overviewMap.size >= slotCount) {
            return { success: false, errorMsg: "캐릭터 슬롯이 부족합니다." };
        }
        const gender = account?.gender ?? 0;
        if (
            this.wzService.canMakeCharacter(gender, [
                face,
                hair,
                topItemId,
                bottomItemId,
                shoesItemId,
                weaponItemId,
            ]) === false
        ) {
            return { success: false, errorMsg: "선택할 수 없는 외형입니다." };
        }

        const characterId = await this.unifiedRepo.reserveCharacterName(name, accountId, worldId);
        if (characterId === null) {
            return { success: false, errorMsg: "이미 사용 중인 이름입니다." };
        }
        const role = account?.role ?? 0;

        const character = {
            characterId,
            accountId,
            worldId,
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

        const overview = characterOverview(character, baseLooks, {}, 0);

        try {
            await this.repo.set(worldId, character);
            await this.keyLayoutRepo.set(worldId, {
                characterId,
                worldId,
                keyLayoutJson: bindingsToJsonString(getDefaultKeyLayoutBindings()),
            });
            await this.mountRepo.set(worldId, { characterId, worldId, level: 1, exp: 0, fatigue: 0 });

            if (equips.length) {
                const enhances = await Promise.all(equips.map((e) => this.wzService.getEnhanceChance(e.itemId)));
                await this.inventoryRepo.replaceBySnapshot(worldId, String(characterId), []);
                await this.inventoryRepo.setAll(
                    worldId,
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
                    })),
                );
            }

            await this.overviewRepo.set(worldId, overview);
            if ((await this.unifiedRepo.confirmCharacterName(characterId)) === false) {
                throw new Error(`name reservation for ${name} was released during creation`);
            }
        } catch (err) {
            await Promise.allSettled([
                this.repo.delete({ worldId, characterId, accountId }),
                this.keyLayoutRepo.delete({ worldId, characterId }),
                this.mountRepo.delete({ worldId, characterId }),
                this.inventoryRepo.replaceBySnapshot(worldId, String(characterId), []),
                this.overviewRepo.delete({ worldId, accountId, characterId }),
                this.unifiedRepo.releaseCharacterNameReservation(accountId),
            ]);
            throw err;
        }

        return { success: true, errorMsg: "", character: overview };
    }

    async deleteCharacter(
        accountId: number,
        characterId: number,
        worldId: number,
    ): Promise<{ success: boolean; secondPasswordSet?: boolean; unsettled?: boolean }> {
        this.assertWorld(worldId);
        this.assertAccountId(accountId);
        this.assertCharacterId(characterId);

        const session = await this.sessionRepo.getAccountSession(accountId);
        if (session?.state !== AccountSessionState.ACCOUNT_SESSION_STATE_LOGIN) {
            return { success: false };
        }

        const account = await this.accountRepo.get(accountId);
        if (!account) {
            return { success: false };
        }
        if (account.secondPasswordHash) {
            return { success: false, secondPasswordSet: true };
        }

        await using _characterLock = await this.distributedLockService.acquireWorldDataLock(
            worldId,
            `character:${characterId}`,
        );

        const character = await this.repo.get(worldId, characterId);
        if (!character || character.accountId !== accountId) {
            return { success: false };
        }
        const unsettledRequests = await this.unsettledRequestRepo.getAll(worldId, String(characterId));
        if (unsettledRequests.size > 0) {
            return { success: false, unsettled: true };
        }
        // Only the character that left them can claim store bank items, and they block the account's next shop.
        const shops = await this.shopRepo.findByAccount(worldId, accountId);
        if (shops.some((shop) => Number(shop.character_id) === characterId)) {
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
        await this.skillMacroRepo.replaceBySnapshot(worldId, String(characterId), []);
        await this.buddyRepo.deleteAllForOwner(worldId, characterId);
        await this.buddyRepo.deleteAllReferencingBuddy(worldId, characterId);
        await this.realtimeStateRepo.delete({ worldId, characterId });
        await this.keyLayoutRepo.delete({ worldId, characterId });
        await this.mountRepo.delete({ worldId, characterId });

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
