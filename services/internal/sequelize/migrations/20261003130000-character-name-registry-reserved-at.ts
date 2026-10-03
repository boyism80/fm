import type { MigrationModule } from "../../src/sequelize-migration";

const migration: MigrationModule = {
    async up(queryInterface) {
        await queryInterface.sequelize.query(`
            ALTER TABLE character_name_registry
            ADD COLUMN IF NOT EXISTS reserved_at TIMESTAMPTZ NULL DEFAULT NULL;
        `);
    },

    async down(queryInterface) {
        await queryInterface.sequelize.query(`
            ALTER TABLE character_name_registry
            DROP COLUMN IF EXISTS reserved_at;
        `);
    },
};

export default migration;
