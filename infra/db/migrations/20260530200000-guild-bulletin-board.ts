import type { MigrationModule } from "../../src/sequelize-migration";

const migration: MigrationModule = {
    async up(queryInterface, Sequelize) {
        await queryInterface.createTable("guild_bulletin_board_threads", {
            world_id: {
                type: Sequelize.INTEGER,
                allowNull: false,
            },
            guild_id: {
                type: Sequelize.BIGINT,
                allowNull: false,
            },
            local_thread_id: {
                type: Sequelize.INTEGER,
                allowNull: false,
            },
            poster_character_id: {
                type: Sequelize.INTEGER,
                allowNull: false,
            },
            title: {
                type: Sequelize.STRING(25),
                allowNull: false,
            },
            body: {
                type: Sequelize.TEXT,
                allowNull: false,
            },
            icon: {
                type: Sequelize.SMALLINT,
                allowNull: false,
                defaultValue: 0,
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
        await queryInterface.addConstraint("guild_bulletin_board_threads", {
            fields: ["world_id", "guild_id", "local_thread_id"],
            type: "primary key",
            name: "pk_guild_bulletin_board_threads",
        });
        await queryInterface.addIndex("guild_bulletin_board_threads", {
            fields: ["world_id", "guild_id"],
            name: "idx_guild_bulletin_board_threads_world_guild",
        });
        await queryInterface.sequelize.query(`
            CREATE UNIQUE INDEX uk_guild_bulletin_board_notice
            ON guild_bulletin_board_threads (world_id, guild_id)
            WHERE local_thread_id = 0;
        `);
        await queryInterface.sequelize.query(`
            CREATE INDEX idx_guild_bulletin_board_threads_world_guild_local
            ON guild_bulletin_board_threads (world_id, guild_id, local_thread_id DESC);
        `);

        await queryInterface.createTable("guild_bulletin_board_replies", {
            world_id: {
                type: Sequelize.INTEGER,
                allowNull: false,
            },
            guild_id: {
                type: Sequelize.BIGINT,
                allowNull: false,
            },
            local_thread_id: {
                type: Sequelize.INTEGER,
                allowNull: false,
            },
            reply_id: {
                type: Sequelize.INTEGER,
                allowNull: false,
            },
            poster_character_id: {
                type: Sequelize.INTEGER,
                allowNull: false,
            },
            content: {
                type: Sequelize.STRING(25),
                allowNull: false,
            },
            created_at: {
                type: Sequelize.DATE,
                allowNull: false,
                defaultValue: Sequelize.literal("NOW()"),
            },
        });
        await queryInterface.addConstraint("guild_bulletin_board_replies", {
            fields: ["world_id", "guild_id", "local_thread_id", "reply_id"],
            type: "primary key",
            name: "pk_guild_bulletin_board_replies",
        });
        await queryInterface.addIndex("guild_bulletin_board_replies", {
            fields: ["world_id", "guild_id", "local_thread_id"],
            name: "idx_guild_bulletin_board_replies_thread",
        });
    },

    async down(queryInterface) {
        await queryInterface.dropTable("guild_bulletin_board_replies");
        await queryInterface.dropTable("guild_bulletin_board_threads");
    },
};

export default migration;
