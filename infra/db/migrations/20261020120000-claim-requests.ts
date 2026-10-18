import type { MigrationModule } from "@fm/shared/sequelize-migration";

const migration: MigrationModule = {
    async up(queryInterface) {
        await queryInterface.sequelize.query(`
            CREATE TABLE IF NOT EXISTS parcel_claim_requests (
                request_id BIGINT PRIMARY KEY,
                receiver_id INTEGER NOT NULL,
                parcel JSONB NOT NULL,
                created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
            );
        `);
        await queryInterface.sequelize.query(`
            CREATE TABLE IF NOT EXISTS note_remove_requests (
                request_id BIGINT PRIMARY KEY,
                character_id INTEGER NOT NULL,
                population INTEGER NOT NULL,
                created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
            );
        `);
        await queryInterface.sequelize.query(`
            CREATE TABLE IF NOT EXISTS wedding_gift_send_requests (
                request_id BIGINT PRIMARY KEY,
                receiver_id INTEGER NOT NULL,
                aborted BOOLEAN NOT NULL DEFAULT FALSE,
                created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
            );
        `);
        await queryInterface.sequelize.query(`
            CREATE TABLE IF NOT EXISTS wedding_gift_claim_requests (
                request_id BIGINT PRIMARY KEY,
                receiver_id INTEGER NOT NULL,
                gift JSONB NOT NULL,
                created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
            );
        `);
    },

    async down(queryInterface) {
        await queryInterface.sequelize.query("DROP TABLE IF EXISTS wedding_gift_claim_requests;");
        await queryInterface.sequelize.query("DROP TABLE IF EXISTS wedding_gift_send_requests;");
        await queryInterface.sequelize.query("DROP TABLE IF EXISTS note_remove_requests;");
        await queryInterface.sequelize.query("DROP TABLE IF EXISTS parcel_claim_requests;");
    },
};

export default migration;
