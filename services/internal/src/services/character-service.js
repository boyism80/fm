"use strict";

const MAX_NAME_LEN = 32;
const DEFAULT_MAP_ID = 10000;
const DEFAULT_SPAWN = 3;
const { EQUIP_SLOT, LOOK_SLOT } = require("../constants/equipment-slots");

class CharacterService {
    constructor(
        characterRepository,
        characterOverviewRepository,
        accountRepository,
        unifiedRepository,
        inventoryRepository,
        appConfiguration
    ) {
        this.repo = characterRepository;
        this.overviewRepo = characterOverviewRepository;
        this.accountRepo = accountRepository;
        this.unifiedRepo = unifiedRepository;
        this.inventoryRepo = inventoryRepository;
        this.app = appConfiguration;
    }

    _worldId() {
        return this.app.app.world_id;
    }

    _assertWorld(worldId) {
        const wid = String(worldId);
        if (!this.app.postgresql.worlds[wid]) {
            const err = new Error(`Unknown world_id: ${worldId}`);
            err.code = "UNKNOWN_WORLD";
            throw err;
        }
    }

    _assertCharacterId(characterId) {
        const n = Number(characterId);
        if (!Number.isInteger(n) || n <= 0 || n > 0xffffffff) {
            const err = new Error("character_id must be a positive uint32");
            err.code = "INVALID_CHARACTER_ID";
            throw err;
        }
    }

    _assertAccountId(accountId) {
        const n = Number(accountId);
        if (!Number.isInteger(n) || n <= 0 || n > 0xffffffff) {
            const err = new Error("account_id must be a positive uint32");
            err.code = "INVALID_PAYLOAD";
            throw err;
        }
    }

    _validatePersisted(p) {
        if (typeof p.name !== "string" || p.name.length > MAX_NAME_LEN) {
            const err = new Error(`name must be a string of length <= ${MAX_NAME_LEN}`);
            err.code = "INVALID_PAYLOAD";
            throw err;
        }
    }

    async getCharacter(worldId, characterId) {
        this._assertWorld(worldId);
        this._assertCharacterId(characterId);
        return this.repo.get(worldId, characterId);
    }

    async saveCharacter(persisted, baseLooks, overlays) {
        return this.saveCharacters([{ persisted, baseLooks, overlays }]);
    }

    async saveCharacters(entries) {
        if (!entries || entries.length === 0) return;

        const byWorld = new Map();
        for (const { persisted, baseLooks, overlays } of entries) {
            this._assertWorld(persisted.worldId);
            this._assertCharacterId(persisted.characterId);
            this._assertAccountId(persisted.accountId);
            this._validatePersisted(persisted);
            const wid = persisted.worldId;
            if (!byWorld.has(wid)) byWorld.set(wid, []);
            byWorld.get(wid).push({ persisted, baseLooks, overlays });
        }

        for (const [worldId, group] of byWorld) {
            const models = group.map(({ persisted }) => persisted);
            await this.repo.setAll(worldId, models);

            for (const { persisted, baseLooks, overlays } of group) {
                if (!persisted.accountId) continue;
                const overview = {
                    characterId:   persisted.characterId,
                    accountId:     persisted.accountId,
                    worldId:       persisted.worldId,
                    name:          persisted.name,
                    gender:        persisted.gender,
                    skinColor:     persisted.skinColor,
                    face:          persisted.face,
                    hair:          persisted.hair,
                    level:         persisted.level,
                    classId:       persisted.classId,
                    mapId:         persisted.mapId,
                    spawnPoint:    persisted.spawnPoint,
                    rank:          0,
                    rankDiff:      0,
                    classRank:     0,
                    classRankDiff: 0,
                    baseLooks:     baseLooks ?? {},
                    overlays:      overlays ?? {},
                };
                await this.overviewRepo.set(persisted.worldId, overview);
            }
        }
    }

    async createCharacter(accountId, worldId, params) {
        const { name, face, hair, skinColor, topItemId, bottomItemId, shoesItemId, weaponItemId } = params;
        const wid = worldId ?? this._worldId();
        this._assertWorld(wid);
        this._assertAccountId(accountId);
        this._validatePersisted({ name });

        const existing = await this.unifiedRepo.findCharacterNameEntry(name);
        if (existing) {
            return { success: false, errorMsg: "이미 사용 중인 이름입니다." };
        }

        let nameEntry;
        try {
            nameEntry = await this.unifiedRepo.reserveCharacterName(name, accountId, wid);
        } catch (err) {
            if (err.code === "23505") {
                return { success: false, errorMsg: "이미 사용 중인 이름입니다." };
            }
            throw err;
        }

        const characterId = Number(nameEntry.character_id);
        const account = await this.accountRepo.get(wid, accountId);
        const role = account?.role ?? 0;
        const persisted = {
            characterId,
            accountId: Number(accountId),
            worldId: Number(wid),
            name,
            gender: 0,
            skinColor: Number(skinColor),
            face: Number(face),
            hair: Number(hair),
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
        };
        await this.repo.set(wid, persisted);

        const equips = [
            { itemId: topItemId, slot: EQUIP_SLOT.TOP, lookSlot: LOOK_SLOT.TOP, offset: 1n },
            { itemId: bottomItemId, slot: EQUIP_SLOT.BOTTOM, lookSlot: LOOK_SLOT.BOTTOM, offset: 2n },
            { itemId: shoesItemId, slot: EQUIP_SLOT.SHOES, lookSlot: LOOK_SLOT.SHOES, offset: 3n },
            { itemId: weaponItemId, slot: EQUIP_SLOT.WEAPON, lookSlot: LOOK_SLOT.WEAPON, offset: 4n },
        ].filter((e) => Number(e.itemId) > 0);

        if (equips.length) {
            const baseUniqueId = BigInt(characterId) * 100n;
            await this.inventoryRepo.setAll(
                wid,
                equips.map((e) => ({
                    uniqueId: String(baseUniqueId + e.offset),
                    ownerId: characterId,
                    itemId: Number(e.itemId),
                    slot: e.slot,
                    count: 1,
                    expiration: null,
                    enchantChance: 0,
                    flag: 0,
                    skillBonus: 0,
                    ownerName: null,
                }))
            );
        }

        const baseLooks = {};
        for (const e of equips) {
            baseLooks[e.lookSlot] = Number(e.itemId);
        }

        const overview = {
            characterId,
            accountId: Number(accountId),
            worldId: Number(wid),
            name,
            gender: 0,
            skinColor: Number(skinColor),
            face: Number(face),
            hair: Number(hair),
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

    async deleteCharacter(accountId, characterId) {
        const worldId = this._worldId();
        this._assertAccountId(accountId);
        this._assertCharacterId(characterId);

        const ok = await this.repo.delete({ worldId, characterId, accountId });
        if (!ok) {
            return { success: false };
        }

        await this.unifiedRepo.deleteCharacterName(characterId);
        await this.overviewRepo.delete({ worldId, accountId, characterId });
        return { success: true };
    }
}

module.exports = { CharacterService };
