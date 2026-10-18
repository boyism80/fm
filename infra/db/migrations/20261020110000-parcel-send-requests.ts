import type { MigrationModule } from "../../src/sequelize-migration";

const migration: MigrationModule = {
    async up(queryInterface) {
        await queryInterface.sequelize.query(`
            CREATE TABLE IF NOT EXISTS parcel_send_requests (
                request_id BIGINT PRIMARY KEY,
                receiver_id INTEGER NOT NULL,
                aborted BOOLEAN NOT NULL DEFAULT FALSE,
                created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
            );
        `);
    },

    async down(queryInterface) {
        await queryInterface.sequelize.query("DROP TABLE IF EXISTS parcel_send_requests;");
    },
};

export default migration;
