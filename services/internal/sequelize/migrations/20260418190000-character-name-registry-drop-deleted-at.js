"use strict";

/**
 * Name registry rows are removed on character delete; `deleted_at` is unused.
 * Replace partial unique index with unconditional LOWER(name) uniqueness.
 *
 * @param {import("sequelize").QueryInterface} queryInterface
 */
module.exports = {
    async up(queryInterface) {
        const { sequelize } = queryInterface;
        await sequelize.query(`
DO $bd$
BEGIN
  IF EXISTS (
    SELECT 1 FROM information_schema.columns c
    WHERE c.table_schema = 'public'
      AND c.table_name = 'character_name_registry'
      AND c.column_name = 'deleted_at'
  ) THEN
    DELETE FROM character_name_registry WHERE deleted_at IS NOT NULL;
  END IF;
END $bd$;
`);
        await sequelize.query(`DROP INDEX IF EXISTS idx_character_name_lower`);
        await sequelize.query(
            `ALTER TABLE character_name_registry DROP COLUMN IF EXISTS deleted_at`
        );
        await sequelize.query(
            `CREATE UNIQUE INDEX idx_character_name_lower ON character_name_registry (LOWER(name))`
        );
    },

    async down(queryInterface) {
        const { sequelize } = queryInterface;
        await sequelize.query(`DROP INDEX IF EXISTS idx_character_name_lower`);
        await sequelize.query(`
            ALTER TABLE character_name_registry
            ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP
        `);
        await sequelize.query(`
            CREATE UNIQUE INDEX idx_character_name_lower ON character_name_registry (LOWER(name))
            WHERE deleted_at IS NULL
        `);
    },
};
