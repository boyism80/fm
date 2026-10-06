import type { MigrationModule } from "../../src/sequelize-migration";

const migration: MigrationModule = {
    async up(queryInterface, Sequelize) {
        await queryInterface.addColumn("characters", "hp_ap_used", {
            type: Sequelize.SMALLINT,
            allowNull: false,
            defaultValue: 0,
        });
    },

    async down(queryInterface) {
        await queryInterface.removeColumn("characters", "hp_ap_used");
    },
};

export default migration;
