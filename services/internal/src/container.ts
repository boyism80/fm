import * as awilix from "awilix";
import { loadConfig } from "./config";
import { AppConfiguration } from "./config/app-configuration";
import { InternalContext } from "./context/internal-context";
import { CharacterRepository } from "./repos/character-repository";
import { InventoryRepository } from "./repos/inventory-repository";
import { AccountRepository } from "./repos/account-repository";
import { UnifiedRepository } from "./repos/unified-repository";
import { CharacterOverviewRepository } from "./repos/character-overview-repository";
import { SkillRepository } from "./repos/skill-repository";
import { BuffRepository } from "./repos/buff-repository";
import { QuestRepository } from "./repos/quest-repository";
import { SavedLocationRepository } from "./repos/saved-location-repository";
import { CharacterRecordRepository } from "./repos/character-record-repository";
import { AccountRecordRepository } from "./repos/account-record-repository";
import { MonsterBookRepository } from "./repos/monster-book-repository";
import { SessionRepository } from "./repos/session-repository";
import { KeyLayoutRepository } from "./repos/key-layout-repository";
import { SkillMacroRepository } from "./repos/skill-macro-repository";
import { MountRepository } from "./repos/mount-repository";
import { UnsettledRequestRepository } from "./repos/unsettled-request-repository";
import { StorageRepository } from "./repos/storage-repository";
import { StorageItemRepository } from "./repos/storage-item-repository";
import { ParcelClaimRequestRepository, ParcelRepository, ParcelReturnRequestRepository, ParcelSendRequestRepository } from "./repos/parcel-repository";
import { ShopOpenRequestRepository, ShopRepository } from "./repos/shop-repository";
import { CashRequestRepository, CashShopRepository } from "./repos/cash-shop-repository";
import {
    MarriageRepository,
    MarriageRequestRepository,
    WeddingGiftClaimRequestRepository,
    WeddingGiftSendRequestRepository,
} from "./repos/marriage-repository";
import { PartyRepository } from "./repos/party-repository";
import { PartyMemberRepository } from "./repos/party-member-repository";
import { AllianceChangeRequestRepository, AllianceCreateRequestRepository, AllianceRepository } from "./repos/alliance-repository";
import { GuildChangeRequestRepository, GuildCreateRequestRepository, GuildRepository } from "./repos/guild-repository";
import { GuildMemberRepository } from "./repos/guild-member-repository";
import { GuildBulletinBoardRepository } from "./repos/guild-bulletin-board-repository";
import { CharacterBuddyRepository } from "./repos/character-buddy-repository";
import { CharacterRealtimeStateRepository } from "./repos/character-realtime-state-repository";
import { NoteRemoveRequestRepository } from "./repos/note-request-repository";
import { FamilyMoveRequestRepository } from "./repos/family-request-repository";
import { PlayerNpcRepository } from "./repos/player-npc-repository";
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
import { DistributedLock } from "./system/distributed-lock";
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

export function createAppContainer() {
    const internalConfig = loadConfig();
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
        socialService: awilix.asClass(SocialService).transient(),
        familyService: awilix.asClass(FamilyService).transient(),
    });
    return container;
}
