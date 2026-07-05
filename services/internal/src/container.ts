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
import { SessionRepository } from "./repos/session-repository";
import { KeyLayoutRepository } from "./repos/key-layout-repository";
import { PartyRepository } from "./repos/party-repository";
import { PartyMemberRepository } from "./repos/party-member-repository";
import { AllianceRepository } from "./repos/alliance-repository";
import { GuildRepository } from "./repos/guild-repository";
import { GuildMemberRepository } from "./repos/guild-member-repository";
import { GuildBulletinBoardRepository } from "./repos/guild-bulletin-board-repository";
import { CharacterBuddyRepository } from "./repos/character-buddy-repository";
import { CharacterRealtimeStateRepository } from "./repos/character-realtime-state-repository";
import { CharacterService } from "./services/character-service";
import { WzService } from "./services/wz-service";
import { SkillService } from "./services/skill-service";
import { BuffService } from "./services/buff-service";
import { AccountService } from "./services/account-service";
import { CharacterOverviewService } from "./services/character-overview-service";
import { SessionService } from "./services/session-service";
import { RabbitMQService } from "./services/rabbitmq-service";
import { PartyService } from "./services/party-service";
import { GuildService } from "./services/guild-service";
import { BuddyService } from "./services/buddy-service";
import { DistributedLock } from "./system/distributed-lock";
import { DistributedLockService } from "./services/distributed-lock-service";

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
        sessionRepository: awilix.asClass(SessionRepository).scoped(),
        keyLayoutRepository: awilix.asClass(KeyLayoutRepository).scoped(),
        partyRepository: awilix.asClass(PartyRepository).scoped(),
        partyMemberRepository: awilix.asClass(PartyMemberRepository).scoped(),
        allianceRepository: awilix.asClass(AllianceRepository).scoped(),
        guildRepository: awilix.asClass(GuildRepository).scoped(),
        guildMemberRepository: awilix.asClass(GuildMemberRepository).scoped(),
        guildBulletinBoardRepository: awilix.asClass(GuildBulletinBoardRepository).scoped(),
        characterBuddyRepository: awilix.asClass(CharacterBuddyRepository).scoped(),
        characterRealtimeStateRepository: awilix.asClass(CharacterRealtimeStateRepository).scoped(),
        characterService: awilix.asClass(CharacterService).transient(),
        wzService: awilix.asClass(WzService).singleton(),
        skillService: awilix.asClass(SkillService).transient(),
        buffService: awilix.asClass(BuffService).transient(),
        accountService: awilix.asClass(AccountService).transient(),
        characterOverviewService: awilix.asClass(CharacterOverviewService).transient(),
        sessionService: awilix.asClass(SessionService).transient(),
        rabbitmqService: awilix.asClass(RabbitMQService).singleton(),
        partyService: awilix.asClass(PartyService).transient(),
        guildService: awilix.asClass(GuildService).transient(),
        buddyService: awilix.asClass(BuddyService).transient(),
    });
    return container;
}
