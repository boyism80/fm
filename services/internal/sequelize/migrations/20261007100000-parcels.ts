import type { MigrationModule } from "../../src/sequelize-migration";

const migration: MigrationModule = {
    async up(queryInterface) {
        await queryInterface.sequelize.query(`
            CREATE TABLE IF NOT EXISTS parcels (
                parcel_id SERIAL PRIMARY KEY,
                receiver_id INTEGER NOT NULL,
                sender_name VARCHAR(32) NOT NULL,
                meso INTEGER NOT NULL DEFAULT 0,
                quick BOOLEAN NOT NULL DEFAULT FALSE,
                message VARCHAR(200) NOT NULL DEFAULT '',
                item JSONB NULL,
                notified BOOLEAN NOT NULL DEFAULT FALSE,
                sent_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
            );
        `);
        await queryInterface.sequelize.query(
            "CREATE INDEX IF NOT EXISTS parcels_receiver_id_idx ON parcels (receiver_id);"
        );
    },

    async down(queryInterface) {
        await queryInterface.sequelize.query("DROP TABLE IF EXISTS parcels;");
    },
};

export default migration;
