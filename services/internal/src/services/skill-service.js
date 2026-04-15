"use strict";

class SkillService {
    constructor(skillRepository) {
        this.repo = skillRepository;
    }

    async getSkills(worldId, characterId) {
        const map = await this.repo.getAll(worldId, characterId);
        return [...map.values()];
    }

    async saveSkills(worldId, models) {
        return this.repo.setAll(worldId, models);
    }
}

module.exports = { SkillService };
