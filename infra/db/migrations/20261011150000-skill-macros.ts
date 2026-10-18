import type { MigrationModule } from "../../src/sequelize-migration";

const migration: MigrationModule = {
    async up(queryInterface, Sequelize) {
        await queryInterface.createTable("character_skill_macros", {
            character_id: { type: Sequelize.INTEGER, allowNull: false },
            position: { type: Sequelize.SMALLINT, allowNull: false },
            name: { type: Sequelize.STRING(30), allowNull: false, defaultValue: "" },
            shout: { type: Sequelize.BOOLEAN, allowNull: false, defaultValue: false },
            skill1: { type: Sequelize.BIGINT, allowNull: false, defaultValue: 0 },
            skill2: { type: Sequelize.BIGINT, allowNull: false, defaultValue: 0 },
            skill3: { type: Sequelize.BIGINT, allowNull: false, defaultValue: 0 },
            updated_at: { type: Sequelize.DATE, allowNull: false, defaultValue: Sequelize.literal("NOW()") },
        });
        await queryInterface.addConstraint("character_skill_macros", {
            fields: ["character_id", "position"],
            type: "primary key",
            name: "pk_character_skill_macros",
        });
        await queryInterface.addConstraint("character_skill_macros", {
            fields: ["position"],
            type: "check",
            where: { position: { [Sequelize.Op.between]: [0, 4] } },
            name: "ck_character_skill_macros_position",
        });
    },

    async down(queryInterface) {
        await queryInterface.dropTable("character_skill_macros");
    },
};

export default migration;
