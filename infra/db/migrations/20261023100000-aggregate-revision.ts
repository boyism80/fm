import type { MigrationModule } from "@fm/shared/sequelize-migration";

const migration: MigrationModule = {
    async up(queryInterface) {
        await queryInterface.sequelize.query(`
            ALTER TABLE guilds ADD COLUMN IF NOT EXISTS revision BIGINT NOT NULL DEFAULT 0;
            ALTER TABLE parties ADD COLUMN IF NOT EXISTS revision BIGINT NOT NULL DEFAULT 0;
            ALTER TABLE alliances ADD COLUMN IF NOT EXISTS revision BIGINT NOT NULL DEFAULT 0;
        `);
    },

    async down(queryInterface) {
        await queryInterface.sequelize.query(`
            ALTER TABLE alliances DROP COLUMN IF EXISTS revision;
            ALTER TABLE parties DROP COLUMN IF EXISTS revision;
            ALTER TABLE guilds DROP COLUMN IF EXISTS revision;
        `);
    },
};

export default migration;
