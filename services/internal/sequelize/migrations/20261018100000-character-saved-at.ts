import type { MigrationModule } from "../../src/sequelize-migration";

const migration: MigrationModule = {
    async up(queryInterface) {
        await queryInterface.sequelize.query("ALTER TABLE characters ADD COLUMN IF NOT EXISTS saved_at TIMESTAMPTZ NULL;");
    },

    async down(queryInterface) {
        await queryInterface.sequelize.query("ALTER TABLE characters DROP COLUMN IF EXISTS saved_at;");
    },
};

export default migration;
