import type { MigrationModule } from "../../src/sequelize-migration";

const migration: MigrationModule = {
    async up(queryInterface) {
        await queryInterface.sequelize.query(`
            ALTER TABLE accounts
                ADD COLUMN IF NOT EXISTS nx_cash INTEGER NOT NULL DEFAULT 0,
                ADD COLUMN IF NOT EXISTS maple_point INTEGER NOT NULL DEFAULT 0;
        `);
        await queryInterface.sequelize.query(`
            CREATE TABLE IF NOT EXISTS cash_items (
                serial BIGINT PRIMARY KEY,
                account_id INTEGER NOT NULL,
                item JSONB NOT NULL,
                created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
            );
        `);
        await queryInterface.sequelize.query(
            "CREATE INDEX IF NOT EXISTS cash_items_account_id_idx ON cash_items (account_id);"
        );
        await queryInterface.sequelize.query(`
            CREATE TABLE IF NOT EXISTS cash_wishlists (
                character_id INTEGER PRIMARY KEY,
                commodity_sns INTEGER[] NOT NULL DEFAULT '{}'
            );
        `);
    },

    async down(queryInterface) {
        await queryInterface.sequelize.query("DROP TABLE IF EXISTS cash_wishlists;");
        await queryInterface.sequelize.query("DROP TABLE IF EXISTS cash_items;");
        await queryInterface.sequelize.query("ALTER TABLE accounts DROP COLUMN IF EXISTS nx_cash, DROP COLUMN IF EXISTS maple_point;");
    },
};

export default migration;
