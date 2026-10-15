import type { MigrationModule } from "../../src/sequelize-migration";

const migration: MigrationModule = {
    async up(queryInterface) {
        await queryInterface.sequelize.query(
            "ALTER TABLE cash_transfer_requests ADD COLUMN IF NOT EXISTS source_serial BIGINT;"
        );
        await queryInterface.sequelize.query(
            "CREATE UNIQUE INDEX IF NOT EXISTS cash_transfer_requests_source_serial_idx ON cash_transfer_requests (source_serial);"
        );
    },

    async down(queryInterface) {
        await queryInterface.sequelize.query(
            "DROP INDEX IF EXISTS cash_transfer_requests_source_serial_idx;"
        );
        await queryInterface.sequelize.query(
            "ALTER TABLE cash_transfer_requests DROP COLUMN IF EXISTS source_serial;"
        );
    },
};

export default migration;
