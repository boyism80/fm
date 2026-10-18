import type { MigrationModule } from "../../src/sequelize-migration";

const migration: MigrationModule = {
    async up(queryInterface, Sequelize) {
        await queryInterface.createTable("guild_name_registry", {
            guild_id: {
                type: Sequelize.BIGINT,
                primaryKey: true,
                autoIncrement: true,
                allowNull: false,
            },
            name: {
                type: Sequelize.STRING(45),
                allowNull: false,
            },
            world_id: {
                type: Sequelize.SMALLINT,
                allowNull: false,
            },
            created_at: {
                type: Sequelize.DATE,
                allowNull: false,
                defaultValue: Sequelize.literal("NOW()"),
            },
        });
        await queryInterface.sequelize.query(
            "CREATE UNIQUE INDEX idx_guild_name_lower ON guild_name_registry (LOWER(name));"
        );
    },

    async down(queryInterface) {
        await queryInterface.sequelize.query("DROP INDEX IF EXISTS idx_guild_name_lower;");
        await queryInterface.dropTable("guild_name_registry");
    },
};

export default migration;
