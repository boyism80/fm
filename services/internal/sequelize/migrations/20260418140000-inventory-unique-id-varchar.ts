import type { MigrationModule } from "../../src/sequelize-migration";

const migration: MigrationModule = {
    async up(queryInterface) {
        await queryInterface.sequelize.query(`
            ALTER TABLE inventory
            ALTER COLUMN unique_id TYPE VARCHAR(40)
            USING (CASE WHEN unique_id IS NULL THEN NULL ELSE unique_id::text END);
        `);
    },

    async down(queryInterface) {
        await queryInterface.sequelize.query(`
            ALTER TABLE inventory
            ALTER COLUMN unique_id TYPE BIGINT            USING (
              CASE
                WHEN unique_id IS NULL OR trim(unique_id) = '' THEN NULL
                ELSE unique_id::bigint END
            );
        `);
    },
};

export default migration;
