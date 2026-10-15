import type { MigrationModule } from "../../src/sequelize-migration";

const migration: MigrationModule = {
    async up(queryInterface) {
        await queryInterface.sequelize.query(`
            DROP TABLE IF EXISTS cash_coupon_grants;
            ALTER TABLE cash_coupons ADD COLUMN IF NOT EXISTS request_id TEXT;
            ALTER TABLE cash_transfer_receipts ALTER COLUMN request_id TYPE TEXT;
        `);
    },

    async down(queryInterface) {
        await queryInterface.sequelize.query(`
            ALTER TABLE cash_transfer_receipts ALTER COLUMN request_id TYPE UUID USING request_id::uuid;
            ALTER TABLE cash_coupons DROP COLUMN IF EXISTS request_id;
        `);
    },
};

export default migration;
