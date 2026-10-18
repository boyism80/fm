import type { MigrationModule } from "../../src/sequelize-migration";

const migration: MigrationModule = {
    async up(queryInterface, Sequelize) {
        await queryInterface.addColumn("accounts", "second_password_hash", {
            type: Sequelize.STRING(128),
            allowNull: true,
        });
        await queryInterface.addColumn("accounts", "second_password_failures", {
            type: Sequelize.SMALLINT,
            allowNull: false,
            defaultValue: 0,
        });
        await queryInterface.addColumn("accounts", "second_password_locked_until", {
            type: Sequelize.DATE,
            allowNull: true,
        });
    },

    async down(queryInterface) {
        await queryInterface.removeColumn("accounts", "second_password_locked_until");
        await queryInterface.removeColumn("accounts", "second_password_failures");
        await queryInterface.removeColumn("accounts", "second_password_hash");
    },
};

export default migration;
