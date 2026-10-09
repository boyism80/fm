import type { MigrationModule } from "../../src/sequelize-migration";

const migration: MigrationModule = {
    async up(queryInterface) {
        await queryInterface.sequelize.query("ALTER TABLE shops ADD COLUMN IF NOT EXISTS store_bank_id BIGINT NULL;");
        await queryInterface.sequelize.query("ALTER TABLE shops ADD COLUMN IF NOT EXISTS claimed_at TIMESTAMPTZ NULL;");
        await queryInterface.sequelize.query(
            "CREATE UNIQUE INDEX IF NOT EXISTS shops_store_bank_idx ON shops (store_bank_id) WHERE store_bank_id IS NOT NULL;"
        );
    },

    async down(queryInterface) {
        await queryInterface.sequelize.query("DELETE FROM shops WHERE claimed_at IS NOT NULL;");
        await queryInterface.sequelize.query("DROP INDEX IF EXISTS shops_store_bank_idx;");
        await queryInterface.sequelize.query("ALTER TABLE shops DROP COLUMN IF EXISTS claimed_at;");
        await queryInterface.sequelize.query("ALTER TABLE shops DROP COLUMN IF EXISTS store_bank_id;");
    },
};

export default migration;
