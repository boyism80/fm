const MAX_NAME_LEN = 32;
const DEFAULT_MAP_ID = 10000;
const DEFAULT_SPAWN = 3;
import { EQUIP_SLOT, LOOK_SLOT } from "../constants/equipment-slots";
import { getDefaultKeyLayoutBindings } from "../constants/default-key-layout-bindings";
import { bindingsToJsonString, jsonStringToBindings } from "../grpc/key-layout-io";
import type { KeyLayoutBindingModel } from "../grpc/key-layout-io";
import type { CreateCharacterRequest } from "../protobuf/generated/fminternal/internal_service";
import { AppConfiguration } from "../config/app-configuration";
import { AccountRepository } from "../repos/account-repository";
import { CharacterOverviewRepository } from "../repos/character-overview-repository";
import { CharacterRepository } from "../repos/character-repository";
import type { CharacterModel } from "../repos/character-repository";
import { InventoryRepository } from "../repos/inventory-repository";
import type { InventoryModel } from "../repos/inventory-repository";
import { KeyLayoutRepository } from "../repos/key-layout-repository";
import { SkillRepository } from "../repos/skill-repository";
import type { SkillModel } from "../repos/skill-repository";
import { BuffRepository } from "../repos/buff-repository";
import type { BuffModel } from "../repos/buff-repository";
import { QuestRepository } from "../repos/quest-repository";
import type { QuestModel } from "../repos/quest-repository";
import { CharacterBuddyRepository } from "../repos/character-buddy-repository";
import { CharacterRealtimeStateRepository } from "../repos/character-realtime-state-repository";
import { UnifiedRepository } from "../repos/unified-repository";
import { WzService } from "./wz-service";
import { DistributedLockService } from "./distributed-lock-service";

type CharacterPersistedInput = {
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
    hidden?: boolean;
};

export type CharacterRowModel = CharacterModel;

type SaveCharacterEntry = {
    persisted: CharacterPersistedInput;
    baseLooks?: Record<string, number>;
    overlays?: Record<string, number>;
    inventory?: InventoryModel[];
    skills?: SkillModel[];
    buffs?: BuffModel[];
    quests?: QuestModel[];
    keyLayout?: KeyLayoutBindingModel[];
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
    private readonly keyLayoutRepo: KeyLayoutRepository;
    private readonly buddyRepo: CharacterBuddyRepository;
    private readonly realtimeStateRepo: CharacterRealtimeStateRepository;
    private readonly app: AppConfiguration;
    private readonly wzService: WzService;
    private readonly distributedLockService: DistributedLockService;

    constructor(
        characterRepository: CharacterRepository,
        characterOverviewRepository: CharacterOverviewRepository,
        accountRepository: AccountRepository,
        unifiedRepository: UnifiedRepository,
        inventoryRepository: InventoryRepository,
        skillRepository: SkillRepository,
        buffRepository: BuffRepository,
        questRepository: QuestRepository,
        keyLayoutRepository: KeyLayoutRepository,
        characterBuddyRepository: CharacterBuddyRepository,
        characterRealtimeStateRepository: CharacterRealtimeStateRepository,
        appConfiguration: AppConfiguration,
        wzService: WzService,
        distributedLockService: DistributedLockService
    ) {
        this.repo = characterRepository;
        this.overviewRepo = characterOverviewRepository;
        this.accountRepo = accountRepository;
        this.unifiedRepo = unifiedRepository;
        this.inventoryRepo = inventoryRepository;
        this.skillRepo = skillRepository;
        this.buffRepo = buffRepository;
        this.questRepo = questRepository;
        this.keyLayoutRepo = keyLayoutRepository;
        this.buddyRepo = characterBuddyRepository;
        this.realtimeStateRepo = characterRealtimeStateRepository;
        this.app = appConfiguration;
        this.wzService = wzService;
        this.distributedLockService = distributedLockService;
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

    private validatePersisted(p: { name?: string }) {
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

    async saveCharacter(persisted: CharacterPersistedInput, baseLooks?: Record<string, number>, overlays?: Record<string, number>) {
        return this.saveCharacters([{ persisted, baseLooks, overlays }]);
    }

    async saveCharacters(entries: SaveCharacterEntry[]) {
        if (!entries || entries.length === 0) {
            return;
        }

        const byWorld = new Map<number, SaveCharacterEntry[]>();
        for (const { persisted, baseLooks, overlays, inventory, skills, buffs, quests, keyLayout } of entries) {
            this.assertWorld(persisted.worldId);
            this.assertCharacterId(persisted.characterId);
            this.assertAccountId(persisted.accountId);
            this.validatePersisted(persisted);
            const wid = persisted.worldId;
            if (!byWorld.has(wid)) {
                byWorld.set(wid, []);
            }
            byWorld.get(wid)?.push({ persisted, baseLooks, overlays, inventory, skills, buffs, quests, keyLayout });
        }

        for (const [worldId, group] of byWorld) {
            const lockKeys = group
                .map(({ persisted }) => persisted.characterId)
                .sort((a, b) => a - b)
                .map((characterId) => `character:${characterId}`);
            await using _characterLocks = await this.distributedLockService.acquireWorldDataLocks(worldId, lockKeys);

            const models = group.map(({ persisted }) => persisted);
            await this.repo.setAll(worldId, models);

            for (const { persisted, baseLooks, overlays, inventory, skills, buffs, quests, keyLayout } of group) {
                if (!persisted.accountId) {
                    continue;
                }
                const overview = {
                    characterId: persisted.characterId,
                    accountId: persisted.accountId,
                    worldId: persisted.worldId,
                    name: persisted.name,
                    gender: persisted.gender,
                    skinColor: persisted.skinColor,
                    face: persisted.face,
                    hair: persisted.hair,
                    level: persisted.level,
                    classId: persisted.classId,
                    mapId: persisted.mapId,
                    spawnPoint: persisted.spawnPoint,
                    rank: 0,
                    rankDiff: 0,
                    classRank: 0,
                    classRankDiff: 0,
                    baseLooks: baseLooks ?? {},
                    overlays: overlays ?? {},
                };
                await this.overviewRepo.set(persisted.worldId, overview);

                if (inventory !== undefined) {
                    await this.inventoryRepo.replaceBySnapshot(persisted.worldId, persisted.characterId, inventory);
                }
                if (skills !== undefined) {
                    await this.skillRepo.replaceBySnapshot(persisted.worldId, persisted.characterId, skills);
                }
                if (buffs !== undefined) {
                    await this.buffRepo.replaceBySnapshot(persisted.worldId, persisted.characterId, buffs);
                }
                if (quests !== undefined) {
                    await this.questRepo.replaceBySnapshot(persisted.worldId, persisted.characterId, quests);
                }
                if (keyLayout !== undefined) {
                    await this.keyLayoutRepo.set(persisted.worldId, {
                        characterId: persisted.characterId,
                        worldId: persisted.worldId,
                        keyLayoutJson: bindingsToJsonString(keyLayout),
                    });
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
        this.validatePersisted({ name });

        await using _createCharacterLock = await this.distributedLockService.acquireWorldDataLock(wid, `create_character:a${accountId}`);

        const overviewMap = await this.overviewRepo.getAll(wid, String(accountId));
        const account = await this.accountRepo.get(wid, accountId);
        const slotCount = account?.characterSlotCount ?? 6;
        if (overviewMap.size >= slotCount) {
            return { success: false, errorMsg: "캐릭터 슬롯이 부족합니다." };
        }

        let nameEntry;
        try {
            nameEntry = await this.unifiedRepo.reserveCharacterName(name, accountId, wid);
        } catch (err: unknown) {
            const code = (err as { code?: string })?.code;
            if (code === "23505") {
                return { success: false, errorMsg: "이미 사용 중인 이름입니다." };
            }
            throw err;
        }

        const characterId = nameEntry.character_id;
        const role = account?.role ?? 0;

        const persisted = {
            characterId,
            accountId,
            worldId: wid,
            name,
            gender: 0,
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
            hidden: false,
        };
        await this.repo.set(wid, persisted);
        await this.keyLayoutRepo.set(wid, {
            characterId,
            worldId: wid,
            keyLayoutJson: bindingsToJsonString(getDefaultKeyLayoutBindings()),
        });

        const equips = [
            { itemId: topItemId, slot: EQUIP_SLOT.TOP, lookSlot: LOOK_SLOT.TOP },
            { itemId: bottomItemId, slot: EQUIP_SLOT.BOTTOM, lookSlot: LOOK_SLOT.BOTTOM },
            { itemId: shoesItemId, slot: EQUIP_SLOT.SHOES, lookSlot: LOOK_SLOT.SHOES },
            { itemId: weaponItemId, slot: EQUIP_SLOT.WEAPON, lookSlot: LOOK_SLOT.WEAPON },
        ].filter((e) => e.itemId > 0);

        if (equips.length) {
            const enhances = await Promise.all(equips.map((e) => this.wzService.getEnhanceChance(e.itemId)));
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

        const baseLooks: Record<string, number> = {};
        for (const e of equips) {
            baseLooks[e.lookSlot] = e.itemId;
        }

        const overview = {
            characterId,
            accountId,
            worldId: wid,
            name,
            gender: 0,
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
        await this.overviewRepo.set(wid, overview);

        return { success: true, errorMsg: "", character: overview };
    }

    async deleteCharacter(accountId: number, characterId: number) {
        const worldId = this.worldId();
        this.assertWorld(worldId);
        this.assertAccountId(accountId);
        this.assertCharacterId(characterId);

        await using _characterLock = await this.distributedLockService.acquireWorldDataLock(
            worldId,
            `character:${characterId}`,
        );

        const character = await this.repo.get(worldId, characterId);
        if (!character || character.accountId !== accountId) {
            return { success: false };
        }

        await this.inventoryRepo.replaceBySnapshot(worldId, characterId, []);
        await this.skillRepo.replaceBySnapshot(worldId, characterId, []);
        await this.buffRepo.replaceBySnapshot(worldId, characterId, []);
        await this.questRepo.replaceBySnapshot(worldId, characterId, []);
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
