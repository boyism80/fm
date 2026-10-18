import type { AccountRepository } from "@fm/shared/repos/account-repository";
import type { AccountWorldRepository } from "@fm/shared/repos/account-world-repository";
import type { CharacterOverviewModel, CharacterOverviewRepository } from "@fm/shared/repos/character-overview-repository";
import type { UnifiedRepository } from "@fm/shared/repos/unified-repository";
import type { WzService } from "./wz-service";
import type { DistributedLockService } from "./distributed-lock-service";
import type { CharacterRank, RankingService } from "./ranking-service";

export type CharacterOverviewListItem = CharacterOverviewModel & CharacterRank;

const NO_RANK: CharacterRank = { rank: 0, rankDiff: 0, classRank: 0, classRankDiff: 0 };

export class CharacterOverviewService {
    private readonly unifiedRepo: UnifiedRepository;
    private readonly overviewRepo: CharacterOverviewRepository;
    private readonly accountRepo: AccountRepository;
    private readonly accountWorldRepo: AccountWorldRepository;
    private readonly wzService: WzService;
    private readonly distributedLockService: DistributedLockService;
    private readonly rankingService: RankingService;

    constructor(
        unifiedRepository: UnifiedRepository,
        characterOverviewRepository: CharacterOverviewRepository,
        accountRepository: AccountRepository,
        accountWorldRepository: AccountWorldRepository,
        wzService: WzService,
        distributedLockService: DistributedLockService,
        rankingService: RankingService
    ) {
        this.unifiedRepo = unifiedRepository;
        this.overviewRepo = characterOverviewRepository;
        this.accountRepo = accountRepository;
        this.accountWorldRepo = accountWorldRepository;
        this.wzService = wzService;
        this.distributedLockService = distributedLockService;
        this.rankingService = rankingService;
    }

    async getCharacterList(accountId: number, worldId: number): Promise<{ characters: CharacterOverviewListItem[]; slotCount: number; hasSecondPassword: boolean }> {
        const overviewMap = await this.overviewRepo.getAll(worldId, String(accountId));
        const overviews = [...overviewMap.values()];
        const ranks = await this.rankingService.characterRanks(overviews.map((overview) => overview.characterId));
        const characters = overviews.map((overview) => ({ ...overview, ...(ranks.get(overview.characterId) ?? NO_RANK) }));
        const account = await this.accountRepo.get(accountId);
        const slotCount = await this.accountWorldRepo.findCharacterSlotCount(worldId, accountId);
        return { characters, slotCount, hasSecondPassword: account?.secondPasswordHash != null };
    }

    /** A free name is reserved for the account until it creates the character or its login session ends. */
    async checkCharacterName(name: string, accountId: number, worldId: number) {
        // The client only knows "taken"; an invalid name is reported the same way.
        if (this.wzService.isValidName(name) === false) {
            return { exists: true };
        }
        if (accountId <= 0) {
            return { exists: true };
        }
        // Hold the same lock as createCharacter so a concurrent name check cannot
        // delete the pending reservation while creation is in progress.
        await using _lock = await this.distributedLockService.acquireUnifiedLock(`create_character:a${accountId}`);
        const characterId = await this.unifiedRepo.reserveCharacterName(name, accountId, worldId);
        return { exists: characterId === null };
    }
}
