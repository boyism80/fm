import type { InternalContext } from "../context/internal-context";
import type { AccountRepository } from "../repos/account-repository";
import type { CharacterOverviewModel, CharacterOverviewRepository } from "../repos/character-overview-repository";
import type { UnifiedRepository } from "../repos/unified-repository";

export type CharacterOverviewListItem = CharacterOverviewModel;

export class CharacterOverviewService {
    private readonly ctx: InternalContext;
    private readonly unifiedRepo: UnifiedRepository;
    private readonly overviewRepo: CharacterOverviewRepository;
    private readonly accountRepo: AccountRepository;

    constructor(
        internalContext: InternalContext,
        unifiedRepository: UnifiedRepository,
        characterOverviewRepository: CharacterOverviewRepository,
        accountRepository: AccountRepository
    ) {
        this.ctx = internalContext;
        this.unifiedRepo = unifiedRepository;
        this.overviewRepo = characterOverviewRepository;
        this.accountRepo = accountRepository;
    }

    private worldId() {
        return this.ctx.appConfiguration.app.world_id;
    }

    async getCharacterList(accountId: number, worldId?: number): Promise<{ characters: CharacterOverviewListItem[]; slotCount: number }> {
        const wid = worldId ?? this.worldId();
        const overviewMap = await this.overviewRepo.getAll(wid, String(accountId));
        const characters = [...overviewMap.values()] as CharacterOverviewListItem[];
        const account = await this.accountRepo.get(wid, accountId);
        const slotCount = account?.characterSlotCount ?? 6;
        return { characters, slotCount };
    }

    async checkCharacterName(name: string) {
        const entry = await this.unifiedRepo.findCharacterNameEntry(name);
        return { exists: entry !== null };
    }
}
