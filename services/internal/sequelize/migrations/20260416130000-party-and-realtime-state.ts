import type { MigrationModule } from "../../src/sequelize-migration";

const migration: MigrationModule = {
    async up(queryInterface, Sequelize) {
        await queryInterface.createTable("parties", {
            world_id: {
                type: Sequelize.INTEGER,
                allowNull: false,
            },
            party_id: {
                type: Sequelize.BIGINT,
                allowNull: false,
            },
            leader_character_id: {
                type: Sequelize.INTEGER,
                allowNull: false,
            },
            state: {
                type: Sequelize.STRING(16),
                allowNull: false,
                defaultValue: "ACTIVE",
            },
            revision: {
                type: Sequelize.BIGINT,
                allowNull: false,
                defaultValue: 1,
            },
            disbanded_at: {
                type: Sequelize.DATE,
                allowNull: true,
            },
            created_at: {
                type: Sequelize.DATE,
                allowNull: false,
                defaultValue: Sequelize.literal("NOW()"),
            },
            updated_at: {
                type: Sequelize.DATE,
                allowNull: false,
                defaultValue: Sequelize.literal("NOW()"),
            },
        });
        await queryInterface.addConstraint("parties", {
            fields: ["world_id", "party_id"],
            type: "primary key",
            name: "pk_parties",
        });
        await queryInterface.addIndex("parties", {
            fields: ["world_id", "leader_character_id"],
            name: "idx_parties_world_leader",
        });

        await queryInterface.createTable("party_members", {
            world_id: { type: Sequelize.INTEGER, allowNull: false },
            party_id: { type: Sequelize.BIGINT, allowNull: false },
            character_id: { type: Sequelize.INTEGER, allowNull: false },
            character_name: { type: Sequelize.STRING(32), allowNull: false },
            level: { type: Sequelize.SMALLINT, allowNull: false, defaultValue: 1 },
            class_id: { type: Sequelize.SMALLINT, allowNull: false, defaultValue: 0 },
            role: { type: Sequelize.STRING(16), allowNull: false, defaultValue: "MEMBER" },
            joined_at: { type: Sequelize.DATE, allowNull: false, defaultValue: Sequelize.literal("NOW()") },
            updated_at: { type: Sequelize.DATE, allowNull: false, defaultValue: Sequelize.literal("NOW()") },
        });
        await queryInterface.addConstraint("party_members", {
            fields: ["world_id", "party_id", "character_id"],
            type: "primary key",
            name: "pk_party_members",
        });
        await queryInterface.addConstraint("party_members", {
            fields: ["world_id", "character_id"],
            type: "unique",
            name: "uk_party_members_world_character",
        });
        await queryInterface.addIndex("party_members", {
            fields: ["world_id", "party_id"],
            name: "idx_party_members_world_party",
        });

        await queryInterface.createTable("character_realtime_state", {
            world_id: { type: Sequelize.INTEGER, allowNull: false },
            character_id: { type: Sequelize.INTEGER, allowNull: false },
            party_id: { type: Sequelize.BIGINT, allowNull: true },
            guild_id: { type: Sequelize.INTEGER, allowNull: true },
            updated_at: { type: Sequelize.DATE, allowNull: false, defaultValue: Sequelize.literal("NOW()") },
        });
        await queryInterface.addConstraint("character_realtime_state", {
            fields: ["world_id", "character_id"],
            type: "primary key",
            name: "pk_character_realtime_state",
        });
        await queryInterface.addIndex("character_realtime_state", {
            fields: ["world_id", "party_id"],
            name: "idx_character_realtime_state_world_party",
        });
        await queryInterface.addIndex("character_realtime_state", {
            fields: ["world_id", "guild_id"],
            name: "idx_character_realtime_state_world_guild",
        });
    },

    async down(queryInterface) {
        await queryInterface.dropTable("character_realtime_state");
        await queryInterface.dropTable("party_members");
        await queryInterface.dropTable("parties");
    },
};

export default migration;
