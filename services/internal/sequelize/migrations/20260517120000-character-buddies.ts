import { DEFAULT_BUDDY_CAPACITY } from "../../src/buddy-defaults";
import type { MigrationModule } from "../../src/sequelize-migration";

const migration: MigrationModule = {
    async up(queryInterface, Sequelize) {
        await queryInterface.addColumn("character_realtime_state", "buddy_capacity", {
            type: Sequelize.SMALLINT,
            allowNull: false,
            defaultValue: DEFAULT_BUDDY_CAPACITY,
        });

        await queryInterface.createTable("character_buddies", {
            character_id: {
                type: Sequelize.INTEGER,
                allowNull: false,
            },
            buddy_character_id: {
                type: Sequelize.INTEGER,
                allowNull: false,
            },
            group_name: {
                type: Sequelize.STRING(16),
                allowNull: false,
                defaultValue: "",
            },
            pending: {
                type: Sequelize.BOOLEAN,
                allowNull: false,
                defaultValue: true,
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
        await queryInterface.addConstraint("character_buddies", {
            fields: ["character_id", "buddy_character_id"],
            type: "primary key",
            name: "pk_character_buddies",
        });
        await queryInterface.addIndex("character_buddies", {
            fields: ["character_id"],
            name: "idx_character_buddies_owner",
        });
        await queryInterface.addIndex("character_buddies", {
            fields: ["buddy_character_id"],
            name: "idx_character_buddies_buddy",
        });
        await queryInterface.addIndex("character_buddies", {
            fields: ["character_id", "pending"],
            name: "idx_character_buddies_owner_pending",
        });
    },

    async down(queryInterface) {
        await queryInterface.dropTable("character_buddies");
        await queryInterface.removeColumn("character_realtime_state", "buddy_capacity");
    },
};

export default migration;
