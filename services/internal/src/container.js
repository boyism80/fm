"use strict";

const awilix = require("awilix");
const { loadConfig } = require("./config");
const { AppConfiguration } = require("./config/app-configuration");
const { InternalContext } = require("./context/internal-context");
const { CharacterRepository } = require("./repos/character-repository");
const { InventoryRepository } = require("./repos/inventory-repository");
const { AccountRepository } = require("./repos/account-repository");
const { UnifiedRepository } = require("./repos/unified-repository");
const { CharacterOverviewRepository } = require("./repos/character-overview-repository");
const { SkillRepository } = require("./repos/skill-repository");
const { SessionRepository } = require("./repos/session-repository");
const { KeyLayoutRepository } = require("./repos/key-layout-repository");
const { PartyRepository } = require("./repos/party-repository");
const { PartyMemberRepository } = require("./repos/party-member-repository");
const { CharacterRealtimeStateRepository } = require("./repos/character-realtime-state-repository");
const { CharacterService } = require("./services/character-service");
const { SkillService } = require("./services/skill-service");
const { AccountService } = require("./services/account-service");
const { CharacterOverviewService } = require("./services/character-overview-service");
const { SessionService } = require("./services/session-service");
const { RabbitMQService } = require("./services/rabbitmq-service");
const { PartyService } = require("./services/party-service");

function createAppContainer() {
    const internalConfig = loadConfig();
    const container = awilix.createContainer({
        injectionMode: awilix.InjectionMode.CLASSIC,
    });
    container.register({
        internalConfig:              awilix.asValue(internalConfig),
        appConfiguration:            awilix.asClass(AppConfiguration).singleton(),
        internalContext:             awilix.asClass(InternalContext).singleton(),
        characterRepository:         awilix.asClass(CharacterRepository).singleton(),
        inventoryRepository:         awilix.asClass(InventoryRepository).singleton(),
        accountRepository:           awilix.asClass(AccountRepository).singleton(),
        unifiedRepository:           awilix.asClass(UnifiedRepository).singleton(),
        characterOverviewRepository: awilix.asClass(CharacterOverviewRepository).singleton(),
        skillRepository:             awilix.asClass(SkillRepository).singleton(),
        sessionRepository:           awilix.asClass(SessionRepository).singleton(),
        keyLayoutRepository:         awilix.asClass(KeyLayoutRepository).singleton(),
        partyRepository:             awilix.asClass(PartyRepository).singleton(),
        partyMemberRepository:       awilix.asClass(PartyMemberRepository).singleton(),
        characterRealtimeStateRepository: awilix.asClass(CharacterRealtimeStateRepository).singleton(),
        characterService:            awilix.asClass(CharacterService).singleton(),
        skillService:                awilix.asClass(SkillService).singleton(),
        accountService:              awilix.asClass(AccountService).singleton(),
        characterOverviewService:    awilix.asClass(CharacterOverviewService).singleton(),
        sessionService:              awilix.asClass(SessionService).singleton(),
        rabbitmqService:             awilix.asClass(RabbitMQService).singleton(),
        partyService:                awilix.asClass(PartyService).singleton(),
    });
    return container;
}

module.exports = { createAppContainer };
