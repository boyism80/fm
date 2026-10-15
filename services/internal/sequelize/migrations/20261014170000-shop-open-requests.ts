import type { MigrationModule } from "../../src/sequelize-migration";

const migration: MigrationModule = {
    async up(queryInterface, Sequelize) {
        await queryInterface.createTable("shop_open_requests", {
            request_id: { type: Sequelize.BIGINT, primaryKey: true },
            account_id: { type: Sequelize.INTEGER, allowNull: false },
            shop_id: { type: Sequelize.BIGINT, allowNull: true },
            aborted: { type: Sequelize.BOOLEAN, allowNull: false, defaultValue: false },
            created_at: { type: Sequelize.DATE, allowNull: false, defaultValue: Sequelize.literal("NOW()") },
        });
        await queryInterface.addIndex("shop_open_requests", {
            fields: ["account_id"],
            name: "idx_shop_open_requests_account",
        });
    },

    async down(queryInterface) {
        await queryInterface.dropTable("shop_open_requests");
    },
};

export default migration;
