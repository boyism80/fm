import type { MigrationModule } from "../../src/sequelize-migration";

const migration: MigrationModule = {
    async up(queryInterface, Sequelize) {
        await queryInterface.addColumn("inventory", "pet", {
            type: Sequelize.JSONB,
            allowNull: true,
        });
        await queryInterface.addColumn("storage_items", "pet", {
            type: Sequelize.JSONB,
            allowNull: true,
        });
        await queryInterface.addColumn("characters", "pet_hp_item", {
            type: Sequelize.INTEGER,
            allowNull: false,
            defaultValue: 0,
        });
        await queryInterface.addColumn("characters", "pet_mp_item", {
            type: Sequelize.INTEGER,
            allowNull: false,
            defaultValue: 0,
        });
        await queryInterface.addColumn("characters", "summoned_pet", {
            type: Sequelize.BIGINT,
            allowNull: false,
            defaultValue: 0,
        });
    },

    async down(queryInterface) {
        await queryInterface.removeColumn("characters", "summoned_pet");
        await queryInterface.removeColumn("characters", "pet_mp_item");
        await queryInterface.removeColumn("characters", "pet_hp_item");
        await queryInterface.removeColumn("storage_items", "pet");
        await queryInterface.removeColumn("inventory", "pet");
    },
};

export default migration;
