"use strict";

const MAX_NAME_LEN = 32;

class CharacterService {
    constructor(characterRepository, characterOverviewRepository, appConfiguration) {
        this.repo = characterRepository;
        this.overviewRepo = characterOverviewRepository;
        this.app = appConfiguration;
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
        return this.repo.getById(worldId, characterId);
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
            await this.repo.saveAll(worldId, models);

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
                await this.overviewRepo.upsert(persisted.worldId, overview);
            }
        }
    }
}

module.exports = { CharacterService };
