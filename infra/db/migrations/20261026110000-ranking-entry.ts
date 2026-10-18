import type { MigrationModule } from "@fm/shared/sequelize-migration";

const migration: MigrationModule = {
    async up(queryInterface) {
        await queryInterface.sequelize.query(`
            DROP TABLE IF EXISTS rankings;
            DROP TABLE IF EXISTS ranking_snapshots;
            DROP TABLE IF EXISTS ranking_ranges;
            DROP TABLE IF EXISTS ranking_entries;
            DROP TABLE IF EXISTS ranking_inputs;

            CREATE TABLE IF NOT EXISTS ranking_entry (
                round_id BIGINT NOT NULL,
                filter TEXT NOT NULL,
                rank_no BIGINT NOT NULL,
                subject_id BIGINT NOT NULL,
                name TEXT NOT NULL,
                world_id SMALLINT NOT NULL,
                class_id SMALLINT,
                level SMALLINT,
                value BIGINT NOT NULL,
                PRIMARY KEY (round_id, filter, rank_no),
                UNIQUE (round_id, filter, subject_id)
            ) PARTITION BY LIST (round_id);
            CREATE INDEX IF NOT EXISTS idx_ranking_entry_name ON ranking_entry (round_id, name);
            CREATE SEQUENCE IF NOT EXISTS ranking_round_seq;

            CREATE TABLE IF NOT EXISTS ranking_inputs (
                kind TEXT NOT NULL,
                round_id BIGINT NOT NULL,
                k1 BIGINT NOT NULL,
                k2 BIGINT NOT NULL,
                k3 BIGINT NOT NULL,
                k4 BIGINT NOT NULL,
                subject_id BIGINT NOT NULL,
                name TEXT NOT NULL,
                world_id SMALLINT NOT NULL,
                class_id SMALLINT,
                job_group SMALLINT,
                job_id SMALLINT,
                level SMALLINT,
                value BIGINT NOT NULL,
                PRIMARY KEY (round_id, k1, k2, k3, k4)
            );
        `);
    },

    async down(queryInterface) {
        await queryInterface.sequelize.query(`
            DROP TABLE IF EXISTS ranking_inputs;
            DROP SEQUENCE IF EXISTS ranking_round_seq;
            DROP TABLE IF EXISTS ranking_entry;
        `);
    },
};

export default migration;
