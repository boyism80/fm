import type { MigrationModule } from "@fm/shared/sequelize-migration";

const migration: MigrationModule = {
    async up(queryInterface) {
        await queryInterface.sequelize.query("ALTER TABLE cash_transfer_requests ADD COLUMN IF NOT EXISTS world_id INTEGER;");
    },

    async down(queryInterface) {
        await queryInterface.sequelize.query("ALTER TABLE cash_transfer_requests DROP COLUMN IF EXISTS world_id;");
    },
};

export default migration;
