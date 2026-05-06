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
import { SessionRepository } from "./repos/session-repository";
import { KeyLayoutRepository } from "./repos/key-layout-repository";
import { PartyRepository } from "./repos/party-repository";
import { PartyMemberRepository } from "./repos/party-member-repository";
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

export function createAppContainer() {
    const internalConfig = loadConfig();
    const container = awilix.createContainer({
        injectionMode: awilix.InjectionMode.CLASSIC,
    });
    container.register({
        internalConfig: awilix.asValue(internalConfig),
        appConfiguration: awilix.asClass(AppConfiguration).singleton(),
        internalContext: awilix.asClass(InternalContext).singleton(),
        characterRepository: awilix.asClass(CharacterRepository).singleton(),
        inventoryRepository: awilix.asClass(InventoryRepository).singleton(),
        accountRepository: awilix.asClass(AccountRepository).singleton(),
        unifiedRepository: awilix.asClass(UnifiedRepository).singleton(),
        characterOverviewRepository: awilix.asClass(CharacterOverviewRepository).singleton(),
        skillRepository: awilix.asClass(SkillRepository).singleton(),
        buffRepository: awilix.asClass(BuffRepository).singleton(),
        sessionRepository: awilix.asClass(SessionRepository).singleton(),
        keyLayoutRepository: awilix.asClass(KeyLayoutRepository).singleton(),
        partyRepository: awilix.asClass(PartyRepository).singleton(),
        partyMemberRepository: awilix.asClass(PartyMemberRepository).singleton(),
        characterRealtimeStateRepository: awilix.asClass(CharacterRealtimeStateRepository).singleton(),
        characterService: awilix.asClass(CharacterService).singleton(),
        wzService: awilix.asClass(WzService).singleton(),
        skillService: awilix.asClass(SkillService).singleton(),
        buffService: awilix.asClass(BuffService).singleton(),
        accountService: awilix.asClass(AccountService).singleton(),
        characterOverviewService: awilix.asClass(CharacterOverviewService).singleton(),
        sessionService: awilix.asClass(SessionService).singleton(),
        rabbitmqService: awilix.asClass(RabbitMQService).singleton(),
        partyService: awilix.asClass(PartyService).singleton(),
    });
    return container;
}
