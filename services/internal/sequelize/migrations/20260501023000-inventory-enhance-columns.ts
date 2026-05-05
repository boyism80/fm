import type { MigrationModule } from "../../src/sequelize-migration";

const migration: MigrationModule = {
    async up(queryInterface, Sequelize) {
        await queryInterface.renameColumn("inventory", "enchant_chance", "enhance_chance");
        await queryInterface.addColumn("inventory", "enhance_count", {
            type: Sequelize.SMALLINT,
            allowNull: true,
            defaultValue: 0,
        });
    },

    async down(queryInterface) {
        await queryInterface.removeColumn("inventory", "enhance_count");
        await queryInterface.renameColumn("inventory", "enhance_chance", "enchant_chance");
    },
};

export default migration;
