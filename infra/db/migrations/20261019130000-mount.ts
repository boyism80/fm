import type { MigrationModule } from "../../src/sequelize-migration";

const migration: MigrationModule = {
    async up(queryInterface, Sequelize) {
        await queryInterface.createTable("mount", {
            character_id: {
                type: Sequelize.INTEGER,
                allowNull: false,
            },
            world_id: {
                type: Sequelize.INTEGER,
                allowNull: false,
            },
            level: {
                type: Sequelize.SMALLINT,
                allowNull: false,
                defaultValue: 1,
            },
            exp: {
                type: Sequelize.INTEGER,
                allowNull: false,
                defaultValue: 0,
            },
            fatigue: {
                type: Sequelize.SMALLINT,
                allowNull: false,
                defaultValue: 0,
            },
            updated_at: {
                type: Sequelize.DATE,
                allowNull: false,
                defaultValue: Sequelize.literal("NOW()"),
            },
        });
        await queryInterface.addConstraint("mount", {
            fields: ["character_id", "world_id"],
            type: "primary key",
            name: "pk_mount",
        });
        await queryInterface.sequelize.query(`
            INSERT INTO mount (character_id, world_id, level, exp, fatigue, updated_at)
            SELECT id, world_id, 1, 0, 0, NOW()
            FROM characters;
        `);
    },

    async down(queryInterface) {
        await queryInterface.dropTable("mount");
    },
};

export default migration;
