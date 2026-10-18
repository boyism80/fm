import type { MigrationModule } from "../../src/sequelize-migration";

const migration: MigrationModule = {
    async up(queryInterface, Sequelize) {
        await queryInterface.createTable("character_monster_book", {
            character_id: {
                type: Sequelize.INTEGER,
                allowNull: false,
            },
            card_id: {
                type: Sequelize.INTEGER,
                allowNull: false,
            },
            count: {
                type: Sequelize.SMALLINT,
                allowNull: false,
            },
            updated_at: {
                type: Sequelize.DATE,
                allowNull: false,
                defaultValue: Sequelize.literal("NOW()"),
            },
        });
        await queryInterface.addConstraint("character_monster_book", {
            fields: ["character_id", "card_id"],
            type: "primary key",
            name: "pk_character_monster_book",
        });
        await queryInterface.addIndex("character_monster_book", {
            fields: ["character_id"],
            name: "idx_character_monster_book_owner",
        });
        await queryInterface.addColumn("characters", "monster_book_cover", {
            type: Sequelize.INTEGER,
            allowNull: false,
            defaultValue: 0,
        });
    },

    async down(queryInterface) {
        await queryInterface.removeColumn("characters", "monster_book_cover");
        await queryInterface.dropTable("character_monster_book");
    },
};

export default migration;
