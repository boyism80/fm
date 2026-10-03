import type { InternalContext } from "../context/internal-context";
import type { AccountRepository } from "../repos/account-repository";
import type { CharacterOverviewModel, CharacterOverviewRepository } from "../repos/character-overview-repository";
import type { UnifiedRepository } from "../repos/unified-repository";
import type { WzService } from "./wz-service";
import type { DistributedLockService } from "./distributed-lock-service";

export type CharacterOverviewListItem = CharacterOverviewModel;

export class CharacterOverviewService {
    private readonly ctx: InternalContext;
    private readonly unifiedRepo: UnifiedRepository;
    private readonly overviewRepo: CharacterOverviewRepository;
    private readonly accountRepo: AccountRepository;
    private readonly wzService: WzService;
    private readonly distributedLockService: DistributedLockService;

    constructor(
        internalContext: InternalContext,
        unifiedRepository: UnifiedRepository,
        characterOverviewRepository: CharacterOverviewRepository,
        accountRepository: AccountRepository,
        wzService: WzService,
        distributedLockService: DistributedLockService
    ) {
        this.ctx = internalContext;
        this.unifiedRepo = unifiedRepository;
        this.overviewRepo = characterOverviewRepository;
        this.accountRepo = accountRepository;
        this.wzService = wzService;
        this.distributedLockService = distributedLockService;
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

    /** A free name is reserved for the account until it creates the character or its login session ends. */
    async checkCharacterName(name: string, accountId: number) {
        // The client only knows "taken"; a forbidden name is reported the same way.
        if (this.wzService.isForbiddenName(name)) {
            return { exists: true };
        }
        if (accountId <= 0) {
            return { exists: true };
        }
        // Hold the same lock as createCharacter so a concurrent name check cannot
        // delete the pending reservation while creation is in progress.
        await using _lock = await this.distributedLockService.acquireWorldDataLock(
            this.worldId(),
            `create_character:a${accountId}`
        );
        const characterId = await this.unifiedRepo.reserveCharacterName(name, accountId, this.worldId());
        return { exists: characterId === null };
    }
}
