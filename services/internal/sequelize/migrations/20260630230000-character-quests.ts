import type { MigrationModule } from "../../src/sequelize-migration";

const migration: MigrationModule = {
    async up(queryInterface, Sequelize) {
        await queryInterface.createTable("character_quests", {
            character_id: {
                type: Sequelize.INTEGER,
                allowNull: false,
            },
            quest_id: {
                type: Sequelize.INTEGER,
                allowNull: false,
            },
            status: {
                type: Sequelize.SMALLINT,
                allowNull: false,
            },
            mob_kills: {
                type: Sequelize.JSONB,
                allowNull: false,
                defaultValue: {},
            },
            status_record: {
                type: Sequelize.TEXT,
                allowNull: false,
                defaultValue: "",
            },
            unknown2: {
                type: Sequelize.JSONB,
                allowNull: false,
                defaultValue: {},
            },
            completion_time_unix_ms: {
                type: Sequelize.BIGINT,
                allowNull: true,
            },
            forfeited: {
                type: Sequelize.INTEGER,
                allowNull: false,
                defaultValue: 0,
            },
            updated_at: {
                type: Sequelize.DATE,
                allowNull: false,
                defaultValue: Sequelize.literal("NOW()"),
            },
        });
        await queryInterface.addConstraint("character_quests", {
            fields: ["character_id", "quest_id"],
            type: "primary key",
            name: "pk_character_quests",
        });
        await queryInterface.addIndex("character_quests", {
            fields: ["character_id"],
            name: "idx_character_quests_owner",
        });
    },

    async down(queryInterface) {
        await queryInterface.dropTable("character_quests");
    },
};

export default migration;
