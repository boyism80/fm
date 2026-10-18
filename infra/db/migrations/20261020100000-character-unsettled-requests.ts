import type { MigrationModule } from "@fm/shared/sequelize-migration";

const migration: MigrationModule = {
    async up(queryInterface, Sequelize) {
        await queryInterface.createTable("character_unsettled_requests", {
            character_id: { type: Sequelize.INTEGER, allowNull: false },
            id: { type: Sequelize.BIGINT, allowNull: false },
            type: { type: Sequelize.SMALLINT, allowNull: false },
            target: { type: Sequelize.BIGINT, allowNull: false, defaultValue: 0 },
            detail: { type: Sequelize.TEXT, allowNull: false, defaultValue: "{}" },
            updated_at: { type: Sequelize.DATE, allowNull: false, defaultValue: Sequelize.literal("NOW()") },
        });
        await queryInterface.addConstraint("character_unsettled_requests", {
            fields: ["character_id", "id"],
            type: "primary key",
            name: "pk_character_unsettled_requests",
        });
    },

    async down(queryInterface) {
        await queryInterface.dropTable("character_unsettled_requests");
    },
};

export default migration;
