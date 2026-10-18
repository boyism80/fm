import type { MigrationModule } from "@fm/shared/sequelize-migration";

const migration: MigrationModule = {
    async up(queryInterface, Sequelize) {
        await queryInterface.addColumn("character_quests", "start_time_unix_ms", {
            type: Sequelize.BIGINT,
            allowNull: true,
        });
    },

    async down(queryInterface) {
        await queryInterface.removeColumn("character_quests", "start_time_unix_ms");
    },
};

export default migration;
