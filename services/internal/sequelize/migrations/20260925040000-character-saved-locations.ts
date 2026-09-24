import type { MigrationModule } from "../../src/sequelize-migration";

const migration: MigrationModule = {
    async up(queryInterface, Sequelize) {
        await queryInterface.createTable("character_saved_locations", {
            character_id: {
                type: Sequelize.INTEGER,
                allowNull: false,
            },
            location_key: {
                type: Sequelize.TEXT,
                allowNull: false,
            },
            map_id: {
                type: Sequelize.INTEGER,
                allowNull: false,
            },
            updated_at: {
                type: Sequelize.DATE,
                allowNull: false,
                defaultValue: Sequelize.literal("NOW()"),
            },
        });
        await queryInterface.addConstraint("character_saved_locations", {
            fields: ["character_id", "location_key"],
            type: "primary key",
            name: "pk_character_saved_locations",
        });
        await queryInterface.addIndex("character_saved_locations", {
            fields: ["character_id"],
            name: "idx_character_saved_locations_owner",
        });
    },

    async down(queryInterface) {
        await queryInterface.dropTable("character_saved_locations");
    },
};

export default migration;
