import type { MigrationModule } from "../../src/sequelize-migration";

const DEFAULT_ALLIANCE_RANK_TITLES = ["Master", "Jr. Master", "Member", "Member", "Member"];

const migration: MigrationModule = {
    async up(queryInterface, Sequelize) {
        await queryInterface.sequelize.query(`
            CREATE SEQUENCE alliance_id_seq AS BIGINT START WITH 1 INCREMENT BY 1 MINVALUE 1 NO MAXVALUE CACHE 1;
        `);

        await queryInterface.createTable("alliances", {
            world_id: {
                type: Sequelize.INTEGER,
                allowNull: false,
            },
            alliance_id: {
                type: Sequelize.BIGINT,
                allowNull: false,
                defaultValue: Sequelize.literal("nextval('alliance_id_seq')"),
            },
            name: {
                type: Sequelize.STRING(45),
                allowNull: false,
            },
            leader_character_id: {
                type: Sequelize.INTEGER,
                allowNull: false,
            },
            guild_ids: {
                type: Sequelize.ARRAY(Sequelize.BIGINT),
                allowNull: false,
                defaultValue: [],
            },
            rank_titles: {
                type: Sequelize.ARRAY(Sequelize.TEXT),
                allowNull: false,
                defaultValue: DEFAULT_ALLIANCE_RANK_TITLES,
            },
            capacity: {
                type: Sequelize.INTEGER,
                allowNull: false,
                defaultValue: 2,
            },
            notice: {
                type: Sequelize.STRING(100),
                allowNull: false,
                defaultValue: "",
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
        await queryInterface.addConstraint("alliances", {
            fields: ["world_id", "alliance_id"],
            type: "primary key",
            name: "pk_alliances",
        });
        await queryInterface.addIndex("alliances", {
            fields: ["world_id", "leader_character_id"],
            name: "idx_alliances_world_leader",
        });
        await queryInterface.sequelize.query(`
            CREATE UNIQUE INDEX uk_alliances_world_name
            ON alliances (world_id, name)
            WHERE disbanded_at IS NULL;
        `);

        await queryInterface.addColumn("guilds", "alliance_id", {
            type: Sequelize.BIGINT,
            allowNull: true,
        });
        await queryInterface.addIndex("guilds", {
            fields: ["world_id", "alliance_id"],
            name: "idx_guilds_world_alliance",
        });

        await queryInterface.addColumn("guild_members", "alliance_rank", {
            type: Sequelize.SMALLINT,
            allowNull: true,
        });
    },

    async down(queryInterface) {
        await queryInterface.removeColumn("guild_members", "alliance_rank");
        await queryInterface.removeIndex("guilds", "idx_guilds_world_alliance");
        await queryInterface.sequelize.query("DROP INDEX IF EXISTS uk_alliances_world_name");
        await queryInterface.removeColumn("guilds", "alliance_id");
        await queryInterface.dropTable("alliances");
        await queryInterface.sequelize.query("DROP SEQUENCE IF EXISTS alliance_id_seq");
    },
};

export default migration;
