import type { MigrationModule } from "@fm/shared/sequelize-migration";

const migration: MigrationModule = {
    async up(queryInterface) {
        await queryInterface.sequelize.query(`
            ALTER TABLE ranking_entry DROP CONSTRAINT IF EXISTS ranking_entry_round_id_filter_subject_id_key;
            ALTER TABLE ranking_entry DROP CONSTRAINT IF EXISTS ranking_entry_round_id_subject_id_filter_key;
            ALTER TABLE ranking_entry ADD CONSTRAINT ranking_entry_round_id_subject_id_filter_key UNIQUE (round_id, subject_id, filter);

            ALTER TABLE character_overview
                DROP COLUMN IF EXISTS rank,
                DROP COLUMN IF EXISTS rank_diff,
                DROP COLUMN IF EXISTS class_rank,
                DROP COLUMN IF EXISTS class_rank_diff;
        `);
    },

    async down(queryInterface) {
        await queryInterface.sequelize.query(`
            ALTER TABLE character_overview
                ADD COLUMN IF NOT EXISTS rank INTEGER NOT NULL DEFAULT 0,
                ADD COLUMN IF NOT EXISTS rank_diff INTEGER NOT NULL DEFAULT 0,
                ADD COLUMN IF NOT EXISTS class_rank INTEGER NOT NULL DEFAULT 0,
                ADD COLUMN IF NOT EXISTS class_rank_diff INTEGER NOT NULL DEFAULT 0;

            ALTER TABLE ranking_entry DROP CONSTRAINT IF EXISTS ranking_entry_round_id_subject_id_filter_key;
            ALTER TABLE ranking_entry ADD CONSTRAINT ranking_entry_round_id_filter_subject_id_key UNIQUE (round_id, filter, subject_id);
        `);
    },
};

export default migration;
