import type { MigrationModule } from "../../src/sequelize-migration";

const migration: MigrationModule = {
    async up(queryInterface) {
        await queryInterface.sequelize.query(`
            CREATE TABLE IF NOT EXISTS hired_merchants (
                merchant_id SERIAL PRIMARY KEY,
                account_id INTEGER NOT NULL UNIQUE,
                character_id INTEGER NOT NULL,
                owner_name VARCHAR(32) NOT NULL,
                channel_id INTEGER NULL,
                map_id INTEGER NULL,
                item_id INTEGER NOT NULL,
                title VARCHAR(64) NOT NULL DEFAULT '',
                meso INTEGER NOT NULL DEFAULT 0,
                items JSONB NOT NULL DEFAULT '[]'::jsonb,
                sold JSONB NOT NULL DEFAULT '[]'::jsonb,
                opened_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                closed_at TIMESTAMPTZ NULL
            );
        `);
        await queryInterface.sequelize.query(
            "CREATE INDEX IF NOT EXISTS hired_merchants_channel_id_idx ON hired_merchants (channel_id) WHERE closed_at IS NULL;"
        );
    },

    async down(queryInterface) {
        await queryInterface.sequelize.query("DROP TABLE IF EXISTS hired_merchants;");
    },
};

export default migration;
