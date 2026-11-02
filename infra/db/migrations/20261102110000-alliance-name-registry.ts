import type { MigrationModule } from "@fm/shared/sequelize-migration";

const migration: MigrationModule = {
    async up(queryInterface, Sequelize) {
        await queryInterface.createTable("alliance_name_registry", {
            world_id: {
                type: Sequelize.SMALLINT,
                primaryKey: true,
                allowNull: false,
            },
            alliance_id: {
                type: Sequelize.BIGINT,
                primaryKey: true,
                allowNull: false,
            },
            name: {
                type: Sequelize.STRING(45),
                allowNull: false,
            },
            created_at: {
                type: Sequelize.DATE,
                allowNull: false,
                defaultValue: Sequelize.literal("NOW()"),
            },
        });
        await queryInterface.sequelize.query(
            "CREATE UNIQUE INDEX idx_alliance_name_lower ON alliance_name_registry (LOWER(name));"
        );
    },

    async down(queryInterface) {
        await queryInterface.sequelize.query("DROP INDEX IF EXISTS idx_alliance_name_lower;");
        await queryInterface.dropTable("alliance_name_registry");
    },
};

export default migration;
