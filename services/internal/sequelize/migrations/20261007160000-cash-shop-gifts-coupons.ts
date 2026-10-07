import type { MigrationModule } from "../../src/sequelize-migration";

const migration: MigrationModule = {
    async up(queryInterface) {
        await queryInterface.sequelize.query(
            "ALTER TABLE characters ADD COLUMN IF NOT EXISTS slot_limits INTEGER[] NOT NULL DEFAULT '{32,32,32,32,32}';"
        );
        await queryInterface.sequelize.query(`
            CREATE TABLE IF NOT EXISTS cash_gifts (
                serial BIGINT PRIMARY KEY,
                account_id INTEGER NOT NULL,
                item_id INTEGER NOT NULL,
                sender_name TEXT NOT NULL,
                message TEXT NOT NULL,
                created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
            );
        `);
        await queryInterface.sequelize.query(
            "CREATE INDEX IF NOT EXISTS cash_gifts_account_id_idx ON cash_gifts (account_id);"
        );
        await queryInterface.sequelize.query(`
            CREATE TABLE IF NOT EXISTS cash_coupons (
                code TEXT PRIMARY KEY,
                kind SMALLINT NOT NULL,
                value INTEGER NOT NULL,
                used_by INTEGER,
                used_at TIMESTAMPTZ,
                created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
            );
        `);
    },

    async down(queryInterface) {
        await queryInterface.sequelize.query("DROP TABLE IF EXISTS cash_coupons;");
        await queryInterface.sequelize.query("DROP TABLE IF EXISTS cash_gifts;");
        await queryInterface.sequelize.query("ALTER TABLE characters DROP COLUMN IF EXISTS slot_limits;");
    },
};

export default migration;
