import type { MigrationModule } from "../../src/sequelize-migration";

const migration: MigrationModule = {
    async up(queryInterface) {
        await queryInterface.sequelize.query(`
            ALTER TABLE character_buddies
            ADD COLUMN IF NOT EXISTS requested BOOLEAN NOT NULL DEFAULT FALSE;
        `);
    },

    async down(queryInterface) {
        await queryInterface.sequelize.query(`
            ALTER TABLE character_buddies
            DROP COLUMN IF EXISTS requested;
        `);
    },
};

export default migration;
