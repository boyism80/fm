import type { MigrationModule } from "../../src/sequelize-migration";

const migration: MigrationModule = {
    async up(queryInterface) {
        await queryInterface.sequelize.query(`
            ALTER TABLE inventory
            ALTER COLUMN unique_id TYPE BIGINT
            USING (
              CASE
                WHEN unique_id IS NULL THEN NULL
                WHEN trim(unique_id::text) = '' THEN NULL
                WHEN trim(unique_id::text) ~ '^[0-9]+$' THEN trim(unique_id::text)::bigint
                ELSE NULL
              END
            );
        `);
    },

    async down(queryInterface) {
        await queryInterface.sequelize.query(`
            ALTER TABLE inventory
            ALTER COLUMN unique_id TYPE VARCHAR(40)
            USING (CASE WHEN unique_id IS NULL THEN NULL ELSE unique_id::text END);
        `);
    },
};

export default migration;
