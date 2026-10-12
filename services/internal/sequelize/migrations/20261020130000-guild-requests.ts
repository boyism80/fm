import type { MigrationModule } from "../../src/sequelize-migration";

const migration: MigrationModule = {
    async up(queryInterface) {
        await queryInterface.sequelize.query(`
            CREATE TABLE IF NOT EXISTS guild_requests (
                request_id BIGINT PRIMARY KEY,
                owner_id INTEGER NOT NULL,
                result_id INTEGER,
                created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
            );
        `);
    },

    async down(queryInterface) {
        await queryInterface.sequelize.query("DROP TABLE IF EXISTS guild_requests;");
    },
};

export default migration;
