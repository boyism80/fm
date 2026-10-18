import type { MigrationModule } from "../../src/sequelize-migration";

const migration: MigrationModule = {
    async up(queryInterface) {
        await queryInterface.renameColumn("character_quests", "unknown2", "record_ex");
    },

    async down(queryInterface) {
        await queryInterface.renameColumn("character_quests", "record_ex", "unknown2");
    },
};

export default migration;
