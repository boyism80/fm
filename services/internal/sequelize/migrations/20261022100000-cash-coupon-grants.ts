import type { MigrationModule } from "../../src/sequelize-migration";

const migration: MigrationModule = {
    async up(queryInterface) {
        await queryInterface.sequelize.query(`
            CREATE TABLE IF NOT EXISTS cash_coupon_grants (
                code TEXT PRIMARY KEY,
                account_id INTEGER NOT NULL,
                created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
            );
        `);
    },

    async down(queryInterface) {
        await queryInterface.sequelize.query("DROP TABLE IF EXISTS cash_coupon_grants;");
    },
};

export default migration;
