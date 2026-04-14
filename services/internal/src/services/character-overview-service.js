"use strict";

const DEFAULT_MAP_ID = 10000;
const DEFAULT_SPAWN = 3;
const { EQUIP_SLOT, LOOK_SLOT } = require("../constants/equipment-slots");

class CharacterOverviewService {
    constructor(internalContext, unifiedRepository, characterRepository, characterOverviewRepository, accountRepository, inventoryRepository) {
        this.ctx = internalContext;
        this.unifiedRepo = unifiedRepository;
        this.characterRepo = characterRepository;
        this.overviewRepo = characterOverviewRepository;
        this.accountRepo = accountRepository;
        this.inventoryRepo = inventoryRepository;
    }

    _worldId() {
        return this.ctx.appConfiguration.app.world_id;
    }

    async getCharacterList(accountId, worldId) {
        const wid = worldId ?? this._worldId();
        const characters = await this.overviewRepo.getByAccount(wid, accountId);

        const account = await this.accountRepo.getById(wid, accountId);
        const slotCount = account?.characterSlotCount ?? 6;

        return { characters, slotCount };
    }

    async checkCharacterName(name) {
        const entry = await this.unifiedRepo.findCharacterNameEntry(name);
        return { exists: entry !== null };
    }

    async createCharacter(accountId, worldId, params) {
        const { name, face, hair, skinColor, topItemId, bottomItemId, shoesItemId, weaponItemId } = params;
        const wid = worldId ?? this._worldId();

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

        const characterId = nameEntry.character_id;
        const pool = this.ctx.getPgDataPool(wid, characterId);
        const account = await this.accountRepo.getById(wid, accountId);
        const role = account?.role ?? 0;

        await pool.query(
            `INSERT INTO characters
               (id, account_id, world_id, name, gender, skin_color, face, hair,
                level, class_id, role, str, dex, int_stat, luk, hp, max_hp, mp, max_mp,
                ability_point, exp, map_id, spawn_point, pos_x, pos_y, stance,
                meso, skill_point, deleted, created_at, updated_at)
             VALUES
               ($1,$2,$3,$4,$5,$6,$7,$8,1,0,$9,12,5,4,4,50,50,5,5,0,0,$10,$11,0,0,0,0,0,false,NOW(),NOW())`,
            [characterId, accountId, wid, name, 0, skinColor, face, hair, role, DEFAULT_MAP_ID, DEFAULT_SPAWN]
        );

        const equips = [
            { itemId: topItemId, slot: EQUIP_SLOT.TOP, lookSlot: LOOK_SLOT.TOP, offset: 1n },
            { itemId: bottomItemId, slot: EQUIP_SLOT.BOTTOM, lookSlot: LOOK_SLOT.BOTTOM, offset: 2n },
            { itemId: shoesItemId, slot: EQUIP_SLOT.SHOES, lookSlot: LOOK_SLOT.SHOES, offset: 3n },
            { itemId: weaponItemId, slot: EQUIP_SLOT.WEAPON, lookSlot: LOOK_SLOT.WEAPON, offset: 4n },
        ].filter((e) => Number(e.itemId) > 0);

        if (equips.length) {
            const baseUniqueId = BigInt(characterId) * 100n;
            await this.inventoryRepo.saveAll(
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
        await this.overviewRepo.upsert(wid, overview);

        return { success: true, errorMsg: "", character: overview };
    }

    async deleteCharacter(accountId, characterId) {
        const worldId = this._worldId();
        const pool = this.ctx.getPgDataPool(worldId, characterId);

        const { rowCount } = await pool.query(
            "UPDATE characters SET deleted = true, updated_at = NOW() WHERE id = $1 AND account_id = $2 AND deleted = false",
            [characterId, accountId]
        );
        if (rowCount === 0) {
            return { success: false };
        }

        await this.unifiedRepo.softDeleteCharacterName(characterId);
        await this.overviewRepo.remove(worldId, accountId, characterId);

        const redisKey = this.characterRepo.getRedisKey(worldId, characterId);
        const { client } = this.ctx.getRedisDataAccess(worldId, characterId);
        await client.del(redisKey).catch(() => {});

        return { success: true };
    }
}

module.exports = { CharacterOverviewService };
