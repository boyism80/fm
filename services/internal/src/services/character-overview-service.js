"use strict";

class CharacterOverviewService {
    constructor(internalContext, unifiedRepository, characterOverviewRepository, accountRepository) {
        this.ctx = internalContext;
        this.unifiedRepo = unifiedRepository;
        this.overviewRepo = characterOverviewRepository;
        this.accountRepo = accountRepository;
    }

    _worldId() {
        return this.ctx.appConfiguration.app.world_id;
    }

    async getCharacterList(accountId, worldId) {
        const wid = worldId ?? this._worldId();
        const overviewMap = await this.overviewRepo.getAll(wid, accountId);
        const characters = [...overviewMap.values()];

        const account = await this.accountRepo.get(wid, accountId);
        const slotCount = account?.characterSlotCount ?? 6;

        return { characters, slotCount };
    }

    async checkCharacterName(name) {
        const entry = await this.unifiedRepo.findCharacterNameEntry(name);
        return { exists: entry !== null };
    }
}

module.exports = { CharacterOverviewService };
