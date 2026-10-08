import type { MigrationModule } from "../../src/sequelize-migration";

const TABLES = ["character_records", "account_records"];

const migration: MigrationModule = {
    async up(queryInterface, Sequelize) {
        for (const table of TABLES) {
            await queryInterface.addColumn(table, "expires_at", {
                type: Sequelize.DATE,
                allowNull: true,
            });
            await queryInterface.addColumn(table, "recorded_at", {
                type: Sequelize.DATE,
                allowNull: false,
                defaultValue: Sequelize.literal("NOW()"),
            });
            await queryInterface.sequelize.query(`
                UPDATE ${table} SET
                    expires_at = CASE WHEN expires_at_unix_ms = 0 THEN NULL ELSE to_timestamp(expires_at_unix_ms / 1000.0) END,
                    recorded_at = CASE WHEN updated_at_unix_ms = 0 THEN updated_at ELSE to_timestamp(updated_at_unix_ms / 1000.0) END
            `);
            await queryInterface.removeColumn(table, "expires_at_unix_ms");
            await queryInterface.removeColumn(table, "updated_at_unix_ms");
        }
    },

    async down(queryInterface, Sequelize) {
        for (const table of TABLES) {
            await queryInterface.addColumn(table, "expires_at_unix_ms", {
                type: Sequelize.BIGINT,
                allowNull: false,
                defaultValue: 0,
            });
            await queryInterface.addColumn(table, "updated_at_unix_ms", {
                type: Sequelize.BIGINT,
                allowNull: false,
                defaultValue: 0,
            });
            await queryInterface.sequelize.query(`
                UPDATE ${table} SET
                    expires_at_unix_ms = COALESCE(FLOOR(EXTRACT(EPOCH FROM expires_at) * 1000), 0),
                    updated_at_unix_ms = FLOOR(EXTRACT(EPOCH FROM recorded_at) * 1000)
            `);
            await queryInterface.removeColumn(table, "expires_at");
            await queryInterface.removeColumn(table, "recorded_at");
        }
    },
};

export default migration;
