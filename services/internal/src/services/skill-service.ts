import type { SkillModel, SkillRepository } from "../repos/skill-repository";

export class SkillService {
    private readonly repo: SkillRepository;

    constructor(skillRepository: SkillRepository) {
        this.repo = skillRepository;
    }

    async getSkills(worldId: number, characterId: number): Promise<SkillModel[]> {
        const map = await this.repo.getAll(worldId, String(characterId));
        return [...map.values()] as SkillModel[];
    }

    async saveSkills(worldId: number, models: SkillModel[]) {
        return this.repo.setAll(worldId, models);
    }
}
