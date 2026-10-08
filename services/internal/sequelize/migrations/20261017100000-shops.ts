import type { MigrationModule } from "../../src/sequelize-migration";

const migration: MigrationModule = {
    async up(queryInterface) {
        await queryInterface.sequelize.query("ALTER TABLE entrusted_shops RENAME TO shops;");
        await queryInterface.sequelize.query("ALTER INDEX entrusted_shops_channel_id_idx RENAME TO shops_channel_id_idx;");
        await queryInterface.sequelize.query("ALTER TABLE shops ADD COLUMN kind SMALLINT NOT NULL DEFAULT 5;");
        await queryInterface.sequelize.query("ALTER TABLE shops ADD COLUMN sn INTEGER NULL;");
        await queryInterface.sequelize.query("ALTER TABLE shops DROP CONSTRAINT IF EXISTS hired_merchants_account_id_key;");
        await queryInterface.sequelize.query(
            "CREATE UNIQUE INDEX shops_entrusted_account_idx ON shops (account_id) WHERE kind = 5 AND closed_at IS NULL;"
        );
        await queryInterface.sequelize.query(
            "CREATE UNIQUE INDEX shops_personal_character_idx ON shops (character_id) WHERE kind = 4 AND closed_at IS NULL;"
        );
        await queryInterface.sequelize.query("CREATE INDEX shops_account_id_idx ON shops (account_id);");
    },

    async down(queryInterface) {
        await queryInterface.sequelize.query("DELETE FROM shops WHERE kind <> 5;");
        await queryInterface.sequelize.query("DROP INDEX IF EXISTS shops_account_id_idx;");
        await queryInterface.sequelize.query("DROP INDEX IF EXISTS shops_personal_character_idx;");
        await queryInterface.sequelize.query("DROP INDEX IF EXISTS shops_entrusted_account_idx;");
        await queryInterface.sequelize.query("ALTER TABLE shops ADD CONSTRAINT hired_merchants_account_id_key UNIQUE (account_id);");
        await queryInterface.sequelize.query("ALTER TABLE shops DROP COLUMN sn;");
        await queryInterface.sequelize.query("ALTER TABLE shops DROP COLUMN kind;");
        await queryInterface.sequelize.query("ALTER INDEX shops_channel_id_idx RENAME TO entrusted_shops_channel_id_idx;");
        await queryInterface.sequelize.query("ALTER TABLE shops RENAME TO entrusted_shops;");
    },
};

export default migration;
