import type { MigrationModule } from "../../src/sequelize-migration";

const migration: MigrationModule = {
    async up(queryInterface) {
        await queryInterface.sequelize.query(
            "ALTER TABLE characters ADD COLUMN IF NOT EXISTS teleport_stones INTEGER[] NOT NULL DEFAULT '{}';"
        );
        await queryInterface.sequelize.query(
            "ALTER TABLE characters ADD COLUMN IF NOT EXISTS vip_teleport_stones INTEGER[] NOT NULL DEFAULT '{}';"
        );
    },

    async down(queryInterface) {
        await queryInterface.sequelize.query("ALTER TABLE characters DROP COLUMN IF EXISTS vip_teleport_stones;");
        await queryInterface.sequelize.query("ALTER TABLE characters DROP COLUMN IF EXISTS teleport_stones;");
    },
};

export default migration;
