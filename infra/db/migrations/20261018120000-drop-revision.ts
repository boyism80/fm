import type { MigrationModule } from "../../src/sequelize-migration";

const TABLES = ["parties", "guilds", "alliances"];

const migration: MigrationModule = {
    async up(queryInterface) {
        for (const table of TABLES) {
            await queryInterface.sequelize.query(`ALTER TABLE ${table} DROP COLUMN IF EXISTS revision;`);
        }
    },

    async down(queryInterface) {
        for (const table of TABLES) {
            await queryInterface.sequelize.query(`ALTER TABLE ${table} ADD COLUMN IF NOT EXISTS revision BIGINT NOT NULL DEFAULT 1;`);
        }
    },
};

export default migration;
