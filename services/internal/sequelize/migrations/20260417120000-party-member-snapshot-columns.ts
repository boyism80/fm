import type { MigrationModule } from "../../src/sequelize-migration";

const migration: MigrationModule = {
    async up(queryInterface, Sequelize) {
        await queryInterface.addColumn("party_members", "map_id", {
            type: Sequelize.INTEGER,
            allowNull: false,
            defaultValue: 0,
        });
        await queryInterface.addColumn("party_members", "channel_index", {
            type: Sequelize.INTEGER,
            allowNull: false,
            defaultValue: -2,
        });
        await queryInterface.addColumn("party_members", "door", {
            type: Sequelize.JSONB,
            allowNull: true,
        });
    },

    async down(queryInterface) {
        await queryInterface.removeColumn("party_members", "door");
        await queryInterface.removeColumn("party_members", "channel_index");
        await queryInterface.removeColumn("party_members", "map_id");
    },
};

export default migration;
