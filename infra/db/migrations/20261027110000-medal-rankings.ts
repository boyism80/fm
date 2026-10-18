import type { MigrationModule } from "@fm/shared/sequelize-migration";

const migration: MigrationModule = {
    async up(queryInterface) {
        await queryInterface.sequelize.query(`
            CREATE TABLE IF NOT EXISTS medal_rankings (
                world_id INTEGER NOT NULL,
                kind TEXT NOT NULL,
                month INTEGER NOT NULL,
                character_id INTEGER NOT NULL,
                name TEXT NOT NULL,
                score BIGINT NOT NULL,
                recorded_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                PRIMARY KEY (world_id, kind, month, character_id)
            );
        `);
    },

    async down(queryInterface) {
        await queryInterface.sequelize.query(`DROP TABLE IF EXISTS medal_rankings;`);
    },
};

export default migration;
