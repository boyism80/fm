import type { MigrationModule } from "@fm/shared/sequelize-migration";

const migration: MigrationModule = {
    async up(queryInterface) {
        await queryInterface.sequelize.query(`
            ALTER TABLE character_overview
                ADD COLUMN IF NOT EXISTS cash_weapon INTEGER NOT NULL DEFAULT 0,
                ADD COLUMN IF NOT EXISTS str SMALLINT NOT NULL DEFAULT 0,
                ADD COLUMN IF NOT EXISTS dex SMALLINT NOT NULL DEFAULT 0,
                ADD COLUMN IF NOT EXISTS int_stat SMALLINT NOT NULL DEFAULT 0,
                ADD COLUMN IF NOT EXISTS luk SMALLINT NOT NULL DEFAULT 0,
                ADD COLUMN IF NOT EXISTS hp SMALLINT NOT NULL DEFAULT 0,
                ADD COLUMN IF NOT EXISTS max_hp SMALLINT NOT NULL DEFAULT 0,
                ADD COLUMN IF NOT EXISTS mp SMALLINT NOT NULL DEFAULT 0,
                ADD COLUMN IF NOT EXISTS max_mp SMALLINT NOT NULL DEFAULT 0,
                ADD COLUMN IF NOT EXISTS ability_point SMALLINT NOT NULL DEFAULT 0,
                ADD COLUMN IF NOT EXISTS skill_point SMALLINT NOT NULL DEFAULT 0,
                ADD COLUMN IF NOT EXISTS exp INTEGER NOT NULL DEFAULT 0,
                ADD COLUMN IF NOT EXISTS population SMALLINT NOT NULL DEFAULT 0;
        `);
    },

    async down(queryInterface) {
        await queryInterface.sequelize.query(`
            ALTER TABLE character_overview
                DROP COLUMN IF EXISTS cash_weapon,
                DROP COLUMN IF EXISTS str,
                DROP COLUMN IF EXISTS dex,
                DROP COLUMN IF EXISTS int_stat,
                DROP COLUMN IF EXISTS luk,
                DROP COLUMN IF EXISTS hp,
                DROP COLUMN IF EXISTS max_hp,
                DROP COLUMN IF EXISTS mp,
                DROP COLUMN IF EXISTS max_mp,
                DROP COLUMN IF EXISTS ability_point,
                DROP COLUMN IF EXISTS skill_point,
                DROP COLUMN IF EXISTS exp,
                DROP COLUMN IF EXISTS population;
        `);
    },
};

export default migration;
