import type { MigrationModule } from "../../src/sequelize-migration";

const migration: MigrationModule = {
    async up(queryInterface, Sequelize) {
        await queryInterface.createTable("operation_logs", {
            id: {
                type: Sequelize.BIGINT,
                autoIncrement: true,
                primaryKey: true,
            },
            channel_id: {
                type: Sequelize.INTEGER,
                allowNull: false,
            },
            character_id: {
                type: Sequelize.INTEGER,
                allowNull: false,
            },
            kind: {
                type: Sequelize.TEXT,
                allowNull: false,
            },
            meso: {
                type: Sequelize.BIGINT,
                allowNull: false,
                defaultValue: 0,
            },
            detail: {
                type: Sequelize.TEXT,
                allowNull: false,
                defaultValue: "",
            },
            created_at: {
                type: Sequelize.DATE,
                allowNull: false,
                defaultValue: Sequelize.literal("NOW()"),
            },
        });
        await queryInterface.addIndex("operation_logs", {
            fields: ["character_id"],
            name: "idx_operation_logs_character",
        });
    },

    async down(queryInterface) {
        await queryInterface.dropTable("operation_logs");
    },
};

export default migration;
