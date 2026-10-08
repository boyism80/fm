import type { MigrationModule } from "../../src/sequelize-migration";

const TABLES = ["character_records", "account_records"];

const migration: MigrationModule = {
    async up(queryInterface, Sequelize) {
        for (const table of TABLES) {
            await queryInterface.removeColumn(table, "period");
            await queryInterface.addColumn(table, "expires_at_unix_ms", {
                type: Sequelize.BIGINT,
                allowNull: false,
                defaultValue: 0,
            });
        }
    },

    async down(queryInterface, Sequelize) {
        for (const table of TABLES) {
            await queryInterface.removeColumn(table, "expires_at_unix_ms");
            await queryInterface.addColumn(table, "period", {
                type: Sequelize.SMALLINT,
                allowNull: false,
                defaultValue: 0,
            });
        }
    },
};

export default migration;
