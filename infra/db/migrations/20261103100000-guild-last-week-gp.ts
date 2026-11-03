import type { MigrationModule } from "@fm/shared/sequelize-migration";

const migration: MigrationModule = {
    async up(queryInterface) {
        await queryInterface.sequelize.query(`
            ALTER TABLE guilds ADD COLUMN IF NOT EXISTS last_week_gp INTEGER NOT NULL DEFAULT 0;
            ALTER TABLE guilds ADD COLUMN IF NOT EXISTS last_gp_week INTEGER NOT NULL DEFAULT 0;
        `);
    },

    async down(queryInterface) {
        await queryInterface.sequelize.query(`
            ALTER TABLE guilds DROP COLUMN IF EXISTS last_gp_week;
            ALTER TABLE guilds DROP COLUMN IF EXISTS last_week_gp;
        `);
    },
};

export default migration;
