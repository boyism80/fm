import type { MigrationModule } from "../../src/sequelize-migration";

const migration: MigrationModule = {
    async up(queryInterface) {
        await queryInterface.sequelize.query(`
            CREATE TABLE IF NOT EXISTS storages (
                account_id INTEGER NOT NULL,
                world_id INTEGER NOT NULL,
                slots SMALLINT NOT NULL DEFAULT 4,
                meso INTEGER NOT NULL DEFAULT 0,
                created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                PRIMARY KEY (account_id, world_id)
            );
        `);
        await queryInterface.sequelize.query(`
            CREATE TABLE IF NOT EXISTS storage_items (
                unique_id BIGINT NULL,
                owner_id INTEGER NOT NULL,
                inventory_type SMALLINT NOT NULL,
                item_id INTEGER NOT NULL,
                slot SMALLINT NOT NULL,
                count SMALLINT NOT NULL DEFAULT 1,
                expiration TIMESTAMPTZ NULL,
                enhance_chance SMALLINT NULL,
                enhance_count SMALLINT NOT NULL DEFAULT 0,
                flag SMALLINT NULL,
                skill_bonus SMALLINT NULL,
                owner_name VARCHAR(32) NULL,
                equip_bonus_stats TEXT NULL,
                created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                PRIMARY KEY (owner_id, inventory_type, slot)
            );
        `);
    },

    async down(queryInterface) {
        await queryInterface.sequelize.query("DROP TABLE IF EXISTS storage_items;");
        await queryInterface.sequelize.query("DROP TABLE IF EXISTS storages;");
    },
};

export default migration;
