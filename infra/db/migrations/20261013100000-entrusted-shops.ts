import type { MigrationModule } from "../../src/sequelize-migration";

const migration: MigrationModule = {
    async up(queryInterface) {
        await queryInterface.sequelize.query("ALTER TABLE hired_merchants RENAME TO entrusted_shops;");
        await queryInterface.sequelize.query("ALTER TABLE entrusted_shops RENAME COLUMN merchant_id TO shop_id;");
        await queryInterface.sequelize.query("ALTER INDEX hired_merchants_channel_id_idx RENAME TO entrusted_shops_channel_id_idx;");
    },

    async down(queryInterface) {
        await queryInterface.sequelize.query("ALTER INDEX entrusted_shops_channel_id_idx RENAME TO hired_merchants_channel_id_idx;");
        await queryInterface.sequelize.query("ALTER TABLE entrusted_shops RENAME COLUMN shop_id TO merchant_id;");
        await queryInterface.sequelize.query("ALTER TABLE entrusted_shops RENAME TO hired_merchants;");
    },
};

export default migration;
