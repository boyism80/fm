import type { MigrationModule } from "../../src/sequelize-migration";

const migration: MigrationModule = {
    async up(queryInterface, Sequelize) {
        await queryInterface.createTable("db_failures", {
            id: { type: Sequelize.BIGINT, autoIncrement: true, primaryKey: true },
            world_id: { type: Sequelize.INTEGER, allowNull: false },
            database_role: { type: Sequelize.TEXT, allowNull: false },
            operation: { type: Sequelize.TEXT, allowNull: false },
            subject: { type: Sequelize.TEXT, allowNull: false, defaultValue: "" },
            statement: { type: Sequelize.TEXT, allowNull: false, defaultValue: "" },
            error_code: { type: Sequelize.TEXT, allowNull: false, defaultValue: "" },
            error_message: { type: Sequelize.TEXT, allowNull: false },
            created_at: { type: Sequelize.DATE, allowNull: false, defaultValue: Sequelize.literal("NOW()") },
        });
        await queryInterface.addIndex("db_failures", {
            fields: ["created_at"],
            name: "idx_db_failures_created_at",
        });
    },

    async down(queryInterface) {
        await queryInterface.dropTable("db_failures");
    },
};

export default migration;
