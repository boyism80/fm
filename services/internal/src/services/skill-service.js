"use strict";

class SkillService {
    constructor(skillRepository) {
        this.repo = skillRepository;
    }

    async getSkills(worldId, characterId) {
        return this.repo.getByCharacter(worldId, characterId);
    }

    async saveSkills(worldId, models) {
        return this.repo.saveAll(worldId, models);
    }
}

module.exports = { SkillService };
