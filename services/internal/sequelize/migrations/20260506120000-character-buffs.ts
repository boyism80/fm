import type { MigrationModule } from "../../src/sequelize-migration";

const migration: MigrationModule = {
    async up(queryInterface, Sequelize) {
        await queryInterface.createTable("character_buffs", {
            character_id: {
                type: Sequelize.INTEGER,
                allowNull: false,
            },
            buff_source_id: {
                type: Sequelize.INTEGER,
                allowNull: false,
            },
            kind: {
                type: Sequelize.SMALLINT,
                allowNull: false,
            },
            remaining_duration_ms: {
                type: Sequelize.BIGINT,
                allowNull: true,
            },
            skill_level: {
                type: Sequelize.SMALLINT,
                allowNull: true,
            },
            causer_id: {
                type: Sequelize.INTEGER,
                allowNull: true,
            },
            flag_values: {
                type: Sequelize.JSONB,
                allowNull: false,
                defaultValue: [],
            },
            updated_at: {
                type: Sequelize.DATE,
                allowNull: false,
                defaultValue: Sequelize.literal("NOW()"),
            },
        });
        await queryInterface.addConstraint("character_buffs", {
            fields: ["character_id", "buff_source_id"],
            type: "primary key",
            name: "pk_character_buffs",
        });
        await queryInterface.addIndex("character_buffs", {
            fields: ["character_id"],
            name: "idx_character_buffs_owner",
        });
    },

    async down(queryInterface) {
        await queryInterface.dropTable("character_buffs");
    },
};

export default migration;
