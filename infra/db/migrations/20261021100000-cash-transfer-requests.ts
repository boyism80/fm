import type { MigrationModule } from "@fm/shared/sequelize-migration";

const migration: MigrationModule = {
    async up(queryInterface) {
        await queryInterface.sequelize.query(`
            CREATE TABLE IF NOT EXISTS cash_transfer_requests (
                request_id UUID PRIMARY KEY,
                source_id INTEGER NOT NULL,
                target_id INTEGER NOT NULL,
                kind TEXT NOT NULL,
                detail JSONB NOT NULL,
                completed BOOLEAN NOT NULL DEFAULT FALSE,
                created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
            );
        `);
        await queryInterface.sequelize.query(
            "CREATE INDEX IF NOT EXISTS cash_transfer_requests_pending_idx ON cash_transfer_requests (created_at) WHERE completed = FALSE;"
        );
        await queryInterface.sequelize.query(`
            CREATE TABLE IF NOT EXISTS cash_transfer_receipts (
                request_id UUID PRIMARY KEY,
                created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
            );
        `);
    },

    async down(queryInterface) {
        await queryInterface.sequelize.query("DROP TABLE IF EXISTS cash_transfer_receipts;");
        await queryInterface.sequelize.query("DROP TABLE IF EXISTS cash_transfer_requests;");
    },
};

export default migration;
