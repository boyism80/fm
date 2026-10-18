import type { MigrationModule } from "../../src/sequelize-migration";

const migration: MigrationModule = {
    async up(queryInterface) {
        await queryInterface.sequelize.query(`
            ALTER TABLE characters ADD COLUMN IF NOT EXISTS exp_reached_at TIMESTAMPTZ;
            UPDATE characters SET exp_reached_at = updated_at WHERE exp_reached_at IS NULL;
            ALTER TABLE characters ALTER COLUMN exp_reached_at SET DEFAULT NOW();
            ALTER TABLE characters ALTER COLUMN exp_reached_at SET NOT NULL;

            CREATE INDEX IF NOT EXISTS idx_characters_overall_rank
                ON characters (world_id, level DESC, exp DESC, exp_reached_at, id)
                WHERE NOT deleted AND level >= 10;
            CREATE INDEX IF NOT EXISTS idx_characters_fame_rank
                ON characters (world_id, population DESC, level DESC, exp DESC, id)
                WHERE NOT deleted;

            ALTER TABLE guilds ADD COLUMN IF NOT EXISTS weekly_gp INTEGER NOT NULL DEFAULT 0;
            ALTER TABLE guilds ADD COLUMN IF NOT EXISTS gp_week INTEGER NOT NULL DEFAULT 0;
            CREATE INDEX IF NOT EXISTS idx_guilds_weekly_rank
                ON guilds (world_id, gp_week, weekly_gp DESC, guild_id)
                WHERE disbanded_at IS NULL;

            CREATE TABLE IF NOT EXISTS ranking_snapshots (
                world_id INTEGER NOT NULL,
                kind TEXT NOT NULL,
                generation BIGINT NOT NULL,
                period TEXT NOT NULL,
                active BOOLEAN NOT NULL DEFAULT FALSE,
                computed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                PRIMARY KEY (world_id, kind, generation)
            );
            CREATE UNIQUE INDEX IF NOT EXISTS idx_ranking_snapshots_active
                ON ranking_snapshots (world_id, kind) WHERE active;

            CREATE TABLE IF NOT EXISTS rankings (
                world_id INTEGER NOT NULL,
                kind TEXT NOT NULL,
                generation BIGINT NOT NULL,
                rank INTEGER NOT NULL,
                subject_id BIGINT NOT NULL,
                name TEXT NOT NULL,
                class_id SMALLINT,
                level SMALLINT,
                value BIGINT NOT NULL,
                rank_diff INTEGER NOT NULL DEFAULT 0,
                class_rank INTEGER,
                class_rank_diff INTEGER,
                PRIMARY KEY (world_id, kind, generation, rank)
            );
            CREATE INDEX IF NOT EXISTS idx_rankings_subject
                ON rankings (world_id, kind, generation, subject_id);
        `);
    },

    async down(queryInterface) {
        await queryInterface.sequelize.query(`
            DROP TABLE IF EXISTS rankings;
            DROP TABLE IF EXISTS ranking_snapshots;
            DROP INDEX IF EXISTS idx_guilds_weekly_rank;
            ALTER TABLE guilds DROP COLUMN IF EXISTS gp_week;
            ALTER TABLE guilds DROP COLUMN IF EXISTS weekly_gp;
            DROP INDEX IF EXISTS idx_characters_fame_rank;
            DROP INDEX IF EXISTS idx_characters_overall_rank;
            ALTER TABLE characters DROP COLUMN IF EXISTS exp_reached_at;
        `);
    },
};

export default migration;
