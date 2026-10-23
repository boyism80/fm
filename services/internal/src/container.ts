import { fileURLToPath } from "node:url";
import * as awilix from "awilix";
import { loadConfig } from "@fm/shared/config";
import { AppConfiguration } from "@fm/shared/config/app-configuration";
import { InternalContext } from "@fm/shared/context/internal-context";
import { CharacterRepository } from "@fm/shared/repos/character-repository";
import { InventoryRepository } from "@fm/shared/repos/inventory-repository";
import { AccountRepository } from "@fm/shared/repos/account-repository";
import { AccountWorldRepository } from "@fm/shared/repos/account-world-repository";
import { UnifiedRepository } from "@fm/shared/repos/unified-repository";
import { CharacterOverviewRepository } from "@fm/shared/repos/character-overview-repository";
import { SkillRepository } from "@fm/shared/repos/skill-repository";
import { BuffRepository } from "@fm/shared/repos/buff-repository";
import { QuestRepository } from "@fm/shared/repos/quest-repository";
import { SavedLocationRepository } from "@fm/shared/repos/saved-location-repository";
import { CharacterRecordRepository } from "@fm/shared/repos/character-record-repository";
import { AccountRecordRepository } from "@fm/shared/repos/account-record-repository";
import { MonsterBookRepository } from "@fm/shared/repos/monster-book-repository";
import { SessionRepository } from "@fm/shared/repos/session-repository";
import { KeyLayoutRepository } from "@fm/shared/repos/key-layout-repository";
import { SkillMacroRepository } from "@fm/shared/repos/skill-macro-repository";
import { MountRepository } from "@fm/shared/repos/mount-repository";
import { UnsettledRequestRepository } from "@fm/shared/repos/unsettled-request-repository";
import { StorageRepository } from "@fm/shared/repos/storage-repository";
import { StorageItemRepository } from "@fm/shared/repos/storage-item-repository";
import {
    ParcelClaimRequestRepository,
    ParcelRepository,
    ParcelReturnRequestRepository,
    ParcelSendRequestRepository,
} from "@fm/shared/repos/parcel-repository";
import { ShopOpenRequestRepository, ShopRepository } from "@fm/shared/repos/shop-repository";
import {
    CashBalanceRequestRepository,
    CashRequestRepository,
    CashShopRepository,
} from "@fm/shared/repos/cash-shop-repository";
import {
    MarriageRepository,
    MarriageRequestRepository,
    WeddingGiftClaimRequestRepository,
    WeddingGiftSendRequestRepository,
} from "@fm/shared/repos/marriage-repository";
import { PartyRepository } from "@fm/shared/repos/party-repository";
import { PartyMemberRepository } from "@fm/shared/repos/party-member-repository";
import {
    AllianceChangeRequestRepository,
    AllianceCreateRequestRepository,
    AllianceRepository,
} from "@fm/shared/repos/alliance-repository";
import {
    GuildChangeRequestRepository,
    GuildCreateRequestRepository,
    GuildRepository,
} from "@fm/shared/repos/guild-repository";
import { GuildMemberRepository } from "@fm/shared/repos/guild-member-repository";
import { GuildBulletinBoardRepository } from "@fm/shared/repos/guild-bulletin-board-repository";
import { CharacterBuddyRepository } from "@fm/shared/repos/character-buddy-repository";
import { CharacterRealtimeStateRepository } from "@fm/shared/repos/character-realtime-state-repository";
import { NoteRemoveRequestRepository } from "@fm/shared/repos/note-request-repository";
import { FamilyMoveRequestRepository } from "@fm/shared/repos/family-request-repository";
import { PlayerNpcRepository } from "@fm/shared/repos/player-npc-repository";
import { RankingRepository } from "@fm/shared/repos/ranking-repository";
import { MedalRankingRepository } from "@fm/shared/repos/medal-ranking-repository";
import { CharacterService } from "./services/character-service";
import { WzService } from "./services/wz-service";
import { BuffService } from "./services/buff-service";
import { AccountService } from "./services/account-service";
import { CharacterOverviewService } from "./services/character-overview-service";
import { SessionService } from "./services/session-service";
import { RabbitMQService } from "./services/rabbitmq-service";
import { PartyService } from "./services/party-service";
import { GuildService } from "./services/guild-service";
import { GuildBulletinBoardService } from "./services/guild-bulletin-board-service";
import { AllianceService } from "./services/alliance-service";
import { BuddyService } from "./services/buddy-service";
import { DistributedLock } from "@fm/shared/system/distributed-lock";
import { DistributedLockService } from "./services/distributed-lock-service";
import { ServerTimeService } from "./services/server-time-service";
import { OperationLogService } from "./services/operation-log-service";
import { ParcelService } from "./services/parcel-service";
import { ShopSearchRanking, ShopService } from "./services/shop-service";
import { CashShopService } from "./services/cash-shop-service";
import { FamilyService } from "./services/family-service";
import { SocialService } from "./services/social-service";
import { MarriageService } from "./services/marriage-service";
import { PlayerNpcService } from "./services/player-npc-service";
import { RankingService } from "./services/ranking-service";
import { MedalService } from "./services/medal-service";

export function createAppContainer() {
    const internalConfig = loadConfig(
        process.env.FM_INTERNAL_CONFIG ?? fileURLToPath(new URL("../config.yaml", import.meta.url)),
    );
    const container = awilix.createContainer({
        injectionMode: awilix.InjectionMode.CLASSIC,
    });
    container.register({
        internalConfig: awilix.asValue(internalConfig),
        appConfiguration: awilix.asClass(AppConfiguration).singleton(),
        internalContext: awilix.asClass(InternalContext).singleton(),
        distributedLock: awilix.asClass(DistributedLock).singleton(),
        distributedLockService: awilix.asClass(DistributedLockService).singleton(),
        characterRepository: awilix.asClass(CharacterRepository).scoped(),
        inventoryRepository: awilix.asClass(InventoryRepository).scoped(),
        accountRepository: awilix.asClass(AccountRepository).scoped(),
        accountWorldRepository: awilix.asClass(AccountWorldRepository).scoped(),
        unifiedRepository: awilix.asClass(UnifiedRepository).scoped(),
        characterOverviewRepository: awilix.asClass(CharacterOverviewRepository).scoped(),
        skillRepository: awilix.asClass(SkillRepository).scoped(),
        buffRepository: awilix.asClass(BuffRepository).scoped(),
        questRepository: awilix.asClass(QuestRepository).scoped(),
        savedLocationRepository: awilix.asClass(SavedLocationRepository).scoped(),
        characterRecordRepository: awilix.asClass(CharacterRecordRepository).scoped(),
        accountRecordRepository: awilix.asClass(AccountRecordRepository).scoped(),
        monsterBookRepository: awilix.asClass(MonsterBookRepository).scoped(),
        sessionRepository: awilix.asClass(SessionRepository).scoped(),
        keyLayoutRepository: awilix.asClass(KeyLayoutRepository).scoped(),
        skillMacroRepository: awilix.asClass(SkillMacroRepository).scoped(),
        mountRepository: awilix.asClass(MountRepository).scoped(),
        unsettledRequestRepository: awilix.asClass(UnsettledRequestRepository).scoped(),
        storageRepository: awilix.asClass(StorageRepository).scoped(),
        storageItemRepository: awilix.asClass(StorageItemRepository).scoped(),
        parcelRepository: awilix.asClass(ParcelRepository).scoped(),
        parcelSendRequestRepository: awilix.asClass(ParcelSendRequestRepository).scoped(),
        parcelReturnRequestRepository: awilix.asClass(ParcelReturnRequestRepository).scoped(),
        parcelClaimRequestRepository: awilix.asClass(ParcelClaimRequestRepository).scoped(),
        shopRepository: awilix.asClass(ShopRepository).scoped(),
        shopOpenRequestRepository: awilix.asClass(ShopOpenRequestRepository).scoped(),
        cashShopRepository: awilix.asClass(CashShopRepository).scoped(),
        cashRequestRepository: awilix.asClass(CashRequestRepository).scoped(),
        cashBalanceRequestRepository: awilix.asClass(CashBalanceRequestRepository).scoped(),
        marriageRepository: awilix.asClass(MarriageRepository).scoped(),
        marriageRequestRepository: awilix.asClass(MarriageRequestRepository).scoped(),
        weddingGiftSendRequestRepository: awilix.asClass(WeddingGiftSendRequestRepository).scoped(),
        weddingGiftClaimRequestRepository: awilix.asClass(WeddingGiftClaimRequestRepository).scoped(),
        partyRepository: awilix.asClass(PartyRepository).scoped(),
        partyMemberRepository: awilix.asClass(PartyMemberRepository).scoped(),
        allianceRepository: awilix.asClass(AllianceRepository).scoped(),
        allianceCreateRequestRepository: awilix.asClass(AllianceCreateRequestRepository).scoped(),
        allianceChangeRequestRepository: awilix.asClass(AllianceChangeRequestRepository).scoped(),
        guildRepository: awilix.asClass(GuildRepository).scoped(),
        guildCreateRequestRepository: awilix.asClass(GuildCreateRequestRepository).scoped(),
        guildChangeRequestRepository: awilix.asClass(GuildChangeRequestRepository).scoped(),
        guildMemberRepository: awilix.asClass(GuildMemberRepository).scoped(),
        guildBulletinBoardRepository: awilix.asClass(GuildBulletinBoardRepository).scoped(),
        characterBuddyRepository: awilix.asClass(CharacterBuddyRepository).scoped(),
        characterRealtimeStateRepository: awilix.asClass(CharacterRealtimeStateRepository).scoped(),
        noteRemoveRequestRepository: awilix.asClass(NoteRemoveRequestRepository).scoped(),
        familyMoveRequestRepository: awilix.asClass(FamilyMoveRequestRepository).scoped(),
        playerNpcRepository: awilix.asClass(PlayerNpcRepository).scoped(),
        rankingRepository: awilix.asClass(RankingRepository).scoped(),
        medalRankingRepository: awilix.asClass(MedalRankingRepository).scoped(),
        characterService: awilix.asClass(CharacterService).transient(),
        wzService: awilix.asClass(WzService).singleton(),
        buffService: awilix.asClass(BuffService).transient(),
        accountService: awilix.asClass(AccountService).transient(),
        characterOverviewService: awilix.asClass(CharacterOverviewService).transient(),
        sessionService: awilix.asClass(SessionService).transient(),
        rabbitmqService: awilix.asClass(RabbitMQService).singleton(),
        partyService: awilix.asClass(PartyService).transient(),
        guildService: awilix.asClass(GuildService).transient(),
        guildBulletinBoardService: awilix.asClass(GuildBulletinBoardService).transient(),
        allianceService: awilix.asClass(AllianceService).transient(),
        buddyService: awilix.asClass(BuddyService).transient(),
        serverTimeService: awilix.asClass(ServerTimeService).singleton(),
        operationLogService: awilix.asClass(OperationLogService).transient(),
        parcelService: awilix.asClass(ParcelService).transient(),
        shopService: awilix.asClass(ShopService).transient(),
        shopSearchRanking: awilix.asClass(ShopSearchRanking).singleton(),
        cashShopService: awilix.asClass(CashShopService).transient(),
        marriageService: awilix.asClass(MarriageService).transient(),
        playerNpcService: awilix.asClass(PlayerNpcService).transient(),
        rankingService: awilix.asClass(RankingService).transient(),
        medalService: awilix.asClass(MedalService).transient(),
        socialService: awilix.asClass(SocialService).transient(),
        familyService: awilix.asClass(FamilyService).transient(),
    });
    return container;
}
