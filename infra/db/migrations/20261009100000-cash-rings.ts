import type { MigrationModule } from "../../src/sequelize-migration";

const migration: MigrationModule = {
    async up(queryInterface) {
        await queryInterface.sequelize.query(`
            CREATE TABLE IF NOT EXISTS cash_rings (
                serial BIGINT PRIMARY KEY,
                partner_serial BIGINT NOT NULL,
                character_id INTEGER NOT NULL,
                partner_id INTEGER NOT NULL,
                partner_name VARCHAR(32) NOT NULL,
                item_id INTEGER NOT NULL,
                created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
            );
        `);
        await queryInterface.sequelize.query(
            "CREATE INDEX IF NOT EXISTS cash_rings_character_id_idx ON cash_rings (character_id);"
        );
    },

    async down(queryInterface) {
        await queryInterface.sequelize.query("DROP TABLE IF EXISTS cash_rings;");
    },
};

export default migration;
