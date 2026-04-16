"use strict";

module.exports = {
    async up(queryInterface, Sequelize) {
        await queryInterface.createTable("key_layout", {
            character_id: {
                type: Sequelize.INTEGER,
                allowNull: false,
            },
            world_id: {
                type: Sequelize.INTEGER,
                allowNull: false,
            },
            key_layout_json: {
                type: Sequelize.JSONB,
                allowNull: false,
                defaultValue: {},
            },
            updated_at: {
                type: Sequelize.DATE,
                allowNull: false,
                defaultValue: Sequelize.literal("NOW()"),
            },
        });
        await queryInterface.addConstraint("key_layout", {
            fields: ["character_id", "world_id"],
            type: "primary key",
            name: "pk_key_layout",
        });
        await queryInterface.sequelize.query(`
            INSERT INTO key_layout (character_id, world_id, key_layout_json, updated_at)
            SELECT id, world_id, '{}'::jsonb, NOW()
            FROM characters;
        `);
    },

    async down(queryInterface) {
        await queryInterface.dropTable("key_layout");
    },
};
