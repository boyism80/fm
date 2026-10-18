import type { MigrationModule } from "../../src/sequelize-migration";

const migration: MigrationModule = {
    async up(queryInterface) {
        await queryInterface.sequelize.query("ALTER TABLE inventory DROP CONSTRAINT IF EXISTS inventory_pkey;");
        await queryInterface.sequelize.query("ALTER TABLE inventory ADD COLUMN IF NOT EXISTS id BIGSERIAL NOT NULL;");
        await queryInterface.sequelize.query("ALTER TABLE inventory ADD CONSTRAINT inventory_pkey PRIMARY KEY (id);");
        await queryInterface.sequelize.query("ALTER TABLE inventory ALTER COLUMN unique_id DROP NOT NULL;");
    },

    async down(queryInterface) {
        await queryInterface.sequelize.query("ALTER TABLE inventory DROP CONSTRAINT IF EXISTS inventory_pkey;");
        await queryInterface.sequelize.query("ALTER TABLE inventory DROP COLUMN IF EXISTS id;");
        await queryInterface.sequelize.query("ALTER TABLE inventory ALTER COLUMN unique_id SET NOT NULL;");
        await queryInterface.sequelize.query("ALTER TABLE inventory ADD CONSTRAINT inventory_pkey PRIMARY KEY (unique_id);");
    },
};

export default migration;
