import type { MigrationModule } from "@fm/shared/sequelize-migration";

const migration: MigrationModule = {
    async up(queryInterface) {
        await queryInterface.sequelize.query(`
            CREATE TABLE IF NOT EXISTS account_worlds (
                account_id INTEGER NOT NULL,
                world_id INTEGER NOT NULL,
                character_slot_count SMALLINT NOT NULL,
                PRIMARY KEY (account_id, world_id)
            );
        `);
    },

    async down(queryInterface) {
        await queryInterface.sequelize.query("DROP TABLE IF EXISTS account_worlds;");
    },
};

export default migration;
