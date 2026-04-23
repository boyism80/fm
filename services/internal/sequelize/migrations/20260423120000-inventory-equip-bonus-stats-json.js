"use strict";


module.exports = {
    async up(queryInterface, Sequelize) {
        await queryInterface.addColumn("inventory", "equip_bonus_stats", {
            type: Sequelize.TEXT,
            allowNull: true,
        });
    },

    async down(queryInterface) {
        await queryInterface.removeColumn("inventory", "equip_bonus_stats");
    },
};
