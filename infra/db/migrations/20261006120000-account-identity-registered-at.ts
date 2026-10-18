import type { MigrationModule } from "../../src/sequelize-migration";

const migration: MigrationModule = {
    async up(queryInterface) {
        await queryInterface.sequelize.query(`
            ALTER TABLE account_identity
            ADD COLUMN IF NOT EXISTS registered_at TIMESTAMPTZ NULL DEFAULT NULL;
        `);
        await queryInterface.sequelize.query(`
            UPDATE account_identity SET registered_at = created_at WHERE registered_at IS NULL;
        `);
    },

    async down(queryInterface) {
        await queryInterface.sequelize.query(`
            ALTER TABLE account_identity
            DROP COLUMN IF EXISTS registered_at;
        `);
    },
};

export default migration;
