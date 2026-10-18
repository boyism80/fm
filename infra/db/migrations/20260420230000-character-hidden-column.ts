import type { MigrationModule } from "@fm/shared/sequelize-migration";

const migration: MigrationModule = {
    async up(queryInterface, Sequelize) {
        const columns = await queryInterface.describeTable("characters");
        if (columns.hidden) {
            return;
        }
        await queryInterface.addColumn("characters", "hidden", {
            type: Sequelize.BOOLEAN,
            allowNull: false,
            defaultValue: false,
        });
    },

    async down(queryInterface) {
        await queryInterface.removeColumn("characters", "hidden");
    },
};

export default migration;
