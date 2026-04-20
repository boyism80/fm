"use strict";

/**
 * Inventory rows are keyed by (owner_id, inventory_type, slot).
 * `id` and soft-delete (`deleted`) are removed; snapshots replace rows via DELETE + INSERT.
 *
 * @param {import("sequelize").QueryInterface} queryInterface
 */
module.exports = {
    async up(queryInterface) {
        const { sequelize } = queryInterface;
        await sequelize.query(`DELETE FROM inventory WHERE deleted = TRUE`);
        await sequelize.query(`
            ALTER TABLE inventory ADD COLUMN IF NOT EXISTS inventory_type SMALLINT;
        `);
        await sequelize.query(`
            UPDATE inventory SET inventory_type = 1 WHERE slot < 0;
        `);
        await sequelize.query(`
            UPDATE inventory SET inventory_type = CASE WHEN (item_id / 10000) >= 100 AND (item_id / 10000) < 200 THEN 1
              WHEN (item_id / 10000) >= 200 AND (item_id / 10000) < 300 THEN 2
              WHEN (item_id / 10000) >= 300 AND (item_id / 10000) < 400 THEN 3
              WHEN (item_id / 10000) >= 400 AND (item_id / 10000) < 500 THEN 4
              WHEN (item_id / 10000) >= 500 AND (item_id / 10000) < 600 THEN 5
              ELSE 4
            END
            WHERE slot >= 0 AND inventory_type IS NULL;
        `);
        await sequelize.query(`
            ALTER TABLE inventory ALTER COLUMN inventory_type SET NOT NULL;
        `);
        await sequelize.query(`
            DELETE FROM inventory a
            USING inventory b
            WHERE a.owner_id = b.owner_id
              AND a.inventory_type = b.inventory_type
              AND a.slot = b.slot
              AND a.id < b.id;
        `);
        await sequelize.query(`DROP INDEX IF EXISTS idx_inventory_owner_alive`);
        await sequelize.query(`ALTER TABLE inventory DROP CONSTRAINT IF EXISTS inventory_pkey`);
        await sequelize.query(`ALTER TABLE inventory DROP COLUMN IF EXISTS id`);
        await sequelize.query(`ALTER TABLE inventory DROP COLUMN IF EXISTS deleted`);
        await sequelize.query(`
            ALTER TABLE inventory
            ADD CONSTRAINT inventory_pkey PRIMARY KEY (owner_id, inventory_type, slot);
        `);
        await sequelize.query(`
            CREATE INDEX IF NOT EXISTS idx_inventory_owner ON inventory (owner_id);
        `);
    },

    async down(queryInterface) {
        const { sequelize } = queryInterface;
        await sequelize.query(`DROP INDEX IF EXISTS idx_inventory_owner`);
        await sequelize.query(`ALTER TABLE inventory DROP CONSTRAINT IF EXISTS inventory_pkey`);
        await sequelize.query(`
            ALTER TABLE inventory ADD COLUMN IF NOT EXISTS deleted BOOLEAN NOT NULL DEFAULT FALSE;
        `);
        await sequelize.query(`
            ALTER TABLE inventory ADD COLUMN IF NOT EXISTS id BIGSERIAL NOT NULL;
        `);
        await sequelize.query(`
            ALTER TABLE inventory ADD CONSTRAINT inventory_pkey PRIMARY KEY (id);
        `);
        await sequelize.query(`
            CREATE INDEX IF NOT EXISTS idx_inventory_owner_alive ON inventory (owner_id) WHERE deleted = false;
        `);
        await sequelize.query(`ALTER TABLE inventory DROP COLUMN IF EXISTS inventory_type`);
    },
};
