import type { BuffModel, BuffRepository } from "../repos/buff-repository";

export class BuffService {
    private readonly repo: BuffRepository;

    constructor(buffRepository: BuffRepository) {
        this.repo = buffRepository;
    }

    async getBuffs(worldId: number, characterId: number): Promise<BuffModel[]> {
        const map = await this.repo.getAll(worldId, String(characterId));
        return [...map.values()] as BuffModel[];
    }
}
