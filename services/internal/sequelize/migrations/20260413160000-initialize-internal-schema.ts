import type { MigrationModule } from "../../src/sequelize-migration";

const migration: MigrationModule = {
    async up(queryInterface, Sequelize) {
        await queryInterface.sequelize.query("DROP TABLE IF EXISTS character_overview CASCADE;");
        await queryInterface.sequelize.query("DROP TABLE IF EXISTS character_skills CASCADE;");
        await queryInterface.sequelize.query("DROP TABLE IF EXISTS inventory CASCADE;");
        await queryInterface.sequelize.query("DROP TABLE IF EXISTS characters CASCADE;");
        await queryInterface.sequelize.query("DROP TABLE IF EXISTS accounts CASCADE;");
        await queryInterface.sequelize.query("DROP TABLE IF EXISTS character_name_registry CASCADE;");
        await queryInterface.sequelize.query("DROP TABLE IF EXISTS account_identity CASCADE;");

        await queryInterface.createTable("account_identity", {
            id: {
                type: Sequelize.INTEGER,
                primaryKey: true,
                autoIncrement: true,
            },
            login_id: {
                type: Sequelize.STRING(64),
                allowNull: false,
                unique: true,
            },
            created_at: {
                type: Sequelize.DATE,
                allowNull: false,
                defaultValue: Sequelize.literal("NOW()"),
            },
        });

        await queryInterface.createTable("character_name_registry", {
            character_id: {
                type: Sequelize.INTEGER,
                primaryKey: true,
                autoIncrement: true,
            },
            name: {
                type: Sequelize.STRING(32),
                allowNull: false,
                unique: true,
            },
            account_id: {
                type: Sequelize.INTEGER,
                allowNull: false,
            },
            world_id: {
                type: Sequelize.SMALLINT,
                allowNull: false,
            },
            status: {
                type: Sequelize.SMALLINT,
                allowNull: false,
                defaultValue: 0,
            },
            created_at: {
                type: Sequelize.DATE,
                allowNull: false,
                defaultValue: Sequelize.literal("NOW()"),
            },
        });
        await queryInterface.sequelize.query(
            "CREATE UNIQUE INDEX idx_character_name_lower ON character_name_registry (LOWER(name));"
        );

        await queryInterface.createTable("accounts", {
            id: {
                type: Sequelize.INTEGER,
                primaryKey: true,
                allowNull: false,
            },
            login_id: {
                type: Sequelize.STRING(64),
                allowNull: false,
            },
            password_hash: {
                type: Sequelize.STRING(128),
                allowNull: false,
            },
            gender: {
                type: Sequelize.SMALLINT,
                allowNull: false,
                defaultValue: 0,
            },
            role: {
                type: Sequelize.SMALLINT,
                allowNull: false,
                defaultValue: 0,
            },
            is_banned: {
                type: Sequelize.BOOLEAN,
                allowNull: false,
                defaultValue: false,
            },
            ban_reason: {
                type: Sequelize.TEXT,
                allowNull: true,
            },
            is_chat_blocked: {
                type: Sequelize.BOOLEAN,
                allowNull: false,
                defaultValue: false,
            },
            chat_blocked_until: {
                type: Sequelize.DATE,
                allowNull: true,
            },
            character_slot_count: {
                type: Sequelize.SMALLINT,
                allowNull: false,
                defaultValue: 6,
            },
            last_login_ip: {
                type: Sequelize.STRING(64),
                allowNull: true,
            },
            mac_address: {
                type: Sequelize.STRING(32),
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

        await queryInterface.createTable("characters", {
            id: {
                type: Sequelize.INTEGER,
                primaryKey: true,
                allowNull: false,
                comment: "character id (uint32)",
            },
            account_id: {
                type: Sequelize.INTEGER,
                allowNull: false,
            },
            world_id: {
                type: Sequelize.INTEGER,
                allowNull: false,
            },
            name: {
                type: Sequelize.STRING(32),
                allowNull: false,
            },
            gender: { type: Sequelize.SMALLINT, allowNull: false, defaultValue: 0 },
            skin_color: { type: Sequelize.SMALLINT, allowNull: false, defaultValue: 0 },
            face: { type: Sequelize.INTEGER, allowNull: false, defaultValue: 0 },
            hair: { type: Sequelize.INTEGER, allowNull: false, defaultValue: 0 },
            level: { type: Sequelize.SMALLINT, allowNull: false, defaultValue: 1 },
            class_id: { type: Sequelize.SMALLINT, allowNull: false, defaultValue: 0 },
            role: { type: Sequelize.SMALLINT, allowNull: false, defaultValue: 0 },
            str: { type: Sequelize.SMALLINT, allowNull: false, defaultValue: 0 },
            dex: { type: Sequelize.SMALLINT, allowNull: false, defaultValue: 0 },
            int_stat: { type: Sequelize.SMALLINT, allowNull: false, defaultValue: 0 },
            luk: { type: Sequelize.SMALLINT, allowNull: false, defaultValue: 0 },
            hp: { type: Sequelize.SMALLINT, allowNull: false, defaultValue: 0 },
            max_hp: { type: Sequelize.SMALLINT, allowNull: false, defaultValue: 0 },
            mp: { type: Sequelize.SMALLINT, allowNull: false, defaultValue: 0 },
            max_mp: { type: Sequelize.SMALLINT, allowNull: false, defaultValue: 0 },
            ability_point: { type: Sequelize.SMALLINT, allowNull: false, defaultValue: 0 },
            exp: { type: Sequelize.INTEGER, allowNull: false, defaultValue: 0 },
            map_id: { type: Sequelize.INTEGER, allowNull: false, defaultValue: 0 },
            spawn_point: { type: Sequelize.SMALLINT, allowNull: false, defaultValue: 0 },
            pos_x: { type: Sequelize.SMALLINT, allowNull: false, defaultValue: 0 },
            pos_y: { type: Sequelize.SMALLINT, allowNull: false, defaultValue: 0 },
            stance: { type: Sequelize.SMALLINT, allowNull: false, defaultValue: 0 },
            meso: { type: Sequelize.INTEGER, allowNull: false, defaultValue: 0 },
            skill_point: { type: Sequelize.SMALLINT, allowNull: false, defaultValue: 0 },
            hidden: { type: Sequelize.BOOLEAN, allowNull: false, defaultValue: false },
            deleted: { type: Sequelize.BOOLEAN, allowNull: false, defaultValue: false },
            created_at: { type: Sequelize.DATE, allowNull: false, defaultValue: Sequelize.literal("NOW()") },
            updated_at: { type: Sequelize.DATE, allowNull: false, defaultValue: Sequelize.literal("NOW()") },
        });
        await queryInterface.sequelize.query(
            "CREATE INDEX idx_characters_world_alive ON characters (world_id) WHERE deleted = false;"
        );
        await queryInterface.sequelize.query(
            "CREATE INDEX idx_characters_account_world ON characters (account_id, world_id) WHERE deleted = false;"
        );

        await queryInterface.createTable("inventory", {
            unique_id: {
                type: Sequelize.BIGINT,
                primaryKey: true,
                allowNull: false,
            },
            owner_id: {
                type: Sequelize.INTEGER,
                allowNull: false,
            },
            item_id: {
                type: Sequelize.INTEGER,
                allowNull: false,
            },
            slot: {
                type: Sequelize.SMALLINT,
                allowNull: false,
            },
            count: {
                type: Sequelize.SMALLINT,
                allowNull: false,
                defaultValue: 1,
            },
            expiration: {
                type: Sequelize.DATE,
                allowNull: true,
            },
            enhance_chance: {
                type: Sequelize.SMALLINT,
                allowNull: true,
            },
            enhance_count: {
                type: Sequelize.SMALLINT,
                allowNull: true,
                defaultValue: 0,
            },
            flag: {
                type: Sequelize.SMALLINT,
                allowNull: true,
            },
            skill_bonus: {
                type: Sequelize.SMALLINT,
                allowNull: true,
            },
            owner_name: {
                type: Sequelize.STRING(32),
                allowNull: true,
            },
            deleted: {
                type: Sequelize.BOOLEAN,
                allowNull: false,
                defaultValue: false,
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
        await queryInterface.sequelize.query(
            "CREATE INDEX idx_inventory_owner_alive ON inventory (owner_id) WHERE deleted = false;"
        );

        await queryInterface.createTable("character_overview", {
            character_id: {
                type: Sequelize.INTEGER,
                primaryKey: true,
                allowNull: false,
            },
            account_id: {
                type: Sequelize.INTEGER,
                allowNull: false,
            },
            world_id: {
                type: Sequelize.SMALLINT,
                allowNull: false,
            },
            name: {
                type: Sequelize.STRING(32),
                allowNull: false,
            },
            gender: { type: Sequelize.SMALLINT, allowNull: false, defaultValue: 0 },
            skin_color: { type: Sequelize.SMALLINT, allowNull: false, defaultValue: 0 },
            face: { type: Sequelize.INTEGER, allowNull: false, defaultValue: 0 },
            hair: { type: Sequelize.INTEGER, allowNull: false, defaultValue: 0 },
            level: { type: Sequelize.SMALLINT, allowNull: false, defaultValue: 1 },
            class_id: { type: Sequelize.SMALLINT, allowNull: false, defaultValue: 0 },
            map_id: { type: Sequelize.INTEGER, allowNull: false, defaultValue: 0 },
            spawn_point: { type: Sequelize.SMALLINT, allowNull: false, defaultValue: 0 },
            rank: { type: Sequelize.INTEGER, allowNull: false, defaultValue: 0 },
            rank_diff: { type: Sequelize.INTEGER, allowNull: false, defaultValue: 0 },
            class_rank: { type: Sequelize.INTEGER, allowNull: false, defaultValue: 0 },
            class_rank_diff: { type: Sequelize.INTEGER, allowNull: false, defaultValue: 0 },
            base_looks: {
                type: Sequelize.JSONB,
                allowNull: false,
                defaultValue: {},
            },
            overlays: {
                type: Sequelize.JSONB,
                allowNull: false,
                defaultValue: {},
            },
            deleted: { type: Sequelize.BOOLEAN, allowNull: false, defaultValue: false },
            updated_at: {
                type: Sequelize.DATE,
                allowNull: false,
                defaultValue: Sequelize.literal("NOW()"),
            },
        });
        await queryInterface.sequelize.query(
            "CREATE INDEX idx_character_overview_account ON character_overview (account_id, world_id) WHERE deleted = false;"
        );

        await queryInterface.createTable("character_skills", {
            character_id: {
                type: Sequelize.INTEGER,
                allowNull: false,
            },
            skill_id: {
                type: Sequelize.BIGINT,
                allowNull: false,
            },
            level: {
                type: Sequelize.INTEGER,
                allowNull: false,
                defaultValue: 0,
            },
            master_level: {
                type: Sequelize.INTEGER,
                allowNull: false,
                defaultValue: 0,
            },
            cooldown_end_unix_ms: {
                type: Sequelize.BIGINT,
                allowNull: true,
            },
            updated_at: {
                type: Sequelize.DATE,
                allowNull: false,
                defaultValue: Sequelize.literal("NOW()"),
            },
        });
        await queryInterface.addConstraint("character_skills", {
            fields: ["character_id", "skill_id"],
            type: "primary key",
            name: "pk_character_skills",
        });
        await queryInterface.addIndex("character_skills", {
            fields: ["character_id"],
            name: "idx_character_skills_owner",
        });
    },

    async down(queryInterface) {
        await queryInterface.dropTable("character_overview");
        await queryInterface.dropTable("character_skills");
        await queryInterface.dropTable("inventory");
        await queryInterface.dropTable("characters");
        await queryInterface.dropTable("accounts");
        await queryInterface.dropTable("character_name_registry");
        await queryInterface.dropTable("account_identity");
    },
};

export default migration;
