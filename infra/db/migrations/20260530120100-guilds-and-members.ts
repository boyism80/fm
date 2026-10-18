import type { MigrationModule } from "../../src/sequelize-migration";

const DEFAULT_GUILD_LOGO = JSON.stringify({
    logo: 0,
    logoColor: 0,
    logoBG: 0,
    logoBGColor: 0,
});

const DEFAULT_GUILD_RANK_TITLES = JSON.stringify([
    "Master",
    "Jr. Master",
    "Member",
    "Member",
    "Member",
]);

const migration: MigrationModule = {
    async up(queryInterface, Sequelize) {
        await queryInterface.createTable("guilds", {
            world_id: {
                type: Sequelize.INTEGER,
                allowNull: false,
            },
            guild_id: {
                type: Sequelize.BIGINT,
                allowNull: false,
            },
            name: {
                type: Sequelize.STRING(45),
                allowNull: false,
            },
            leader_character_id: {
                type: Sequelize.INTEGER,
                allowNull: false,
            },
            gp: {
                type: Sequelize.INTEGER,
                allowNull: false,
                defaultValue: 0,
            },
            capacity: {
                type: Sequelize.INTEGER,
                allowNull: false,
                defaultValue: 10,
            },
            notice: {
                type: Sequelize.STRING(100),
                allowNull: false,
                defaultValue: "",
            },
            logo: {
                type: Sequelize.JSONB,
                allowNull: false,
                defaultValue: Sequelize.literal(`'${DEFAULT_GUILD_LOGO}'::jsonb`),
            },
            rank_titles: {
                type: Sequelize.JSONB,
                allowNull: false,
                defaultValue: Sequelize.literal(`'${DEFAULT_GUILD_RANK_TITLES}'::jsonb`),
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
        await queryInterface.addConstraint("guilds", {
            fields: ["world_id", "guild_id"],
            type: "primary key",
            name: "pk_guilds",
        });
        await queryInterface.addIndex("guilds", {
            fields: ["world_id", "leader_character_id"],
            name: "idx_guilds_world_leader",
        });

        await queryInterface.createTable("guild_members", {
            world_id: { type: Sequelize.INTEGER, allowNull: false },
            guild_id: { type: Sequelize.BIGINT, allowNull: false },
            character_id: { type: Sequelize.INTEGER, allowNull: false },
            character_name: { type: Sequelize.STRING(32), allowNull: false },
            level: { type: Sequelize.SMALLINT, allowNull: false, defaultValue: 1 },
            class_id: { type: Sequelize.SMALLINT, allowNull: false, defaultValue: 0 },
            guild_rank: { type: Sequelize.SMALLINT, allowNull: false, defaultValue: 5 },
            joined_at: { type: Sequelize.DATE, allowNull: false, defaultValue: Sequelize.literal("NOW()") },
            updated_at: { type: Sequelize.DATE, allowNull: false, defaultValue: Sequelize.literal("NOW()") },
        });
        await queryInterface.addConstraint("guild_members", {
            fields: ["world_id", "guild_id", "character_id"],
            type: "primary key",
            name: "pk_guild_members",
        });
        await queryInterface.addConstraint("guild_members", {
            fields: ["world_id", "character_id"],
            type: "unique",
            name: "uk_guild_members_world_character",
        });
        await queryInterface.addIndex("guild_members", {
            fields: ["world_id", "guild_id"],
            name: "idx_guild_members_world_guild",
        });
    },

    async down(queryInterface) {
        await queryInterface.dropTable("guild_members");
        await queryInterface.dropTable("guilds");
    },
};

export default migration;
