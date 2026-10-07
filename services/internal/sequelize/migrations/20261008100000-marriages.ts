import type { MigrationModule } from "../../src/sequelize-migration";

const migration: MigrationModule = {
    async up(queryInterface, Sequelize) {
        await queryInterface.sequelize.query(`
            CREATE TABLE IF NOT EXISTS marriages (
                marriage_id SERIAL PRIMARY KEY,
                groom_id INTEGER NOT NULL,
                bride_id INTEGER NOT NULL,
                groom_name VARCHAR(32) NOT NULL,
                bride_name VARCHAR(32) NOT NULL,
                groom_item_id INTEGER NOT NULL,
                bride_item_id INTEGER NOT NULL,
                status SMALLINT NOT NULL,
                ticket_item_id INTEGER NOT NULL DEFAULT 0,
                groom_wished BOOLEAN NOT NULL DEFAULT FALSE,
                bride_wished BOOLEAN NOT NULL DEFAULT FALSE,
                groom_wishes TEXT[] NOT NULL DEFAULT '{}',
                bride_wishes TEXT[] NOT NULL DEFAULT '{}',
                guests INTEGER[] NOT NULL DEFAULT '{}',
                divorce_requester_id INTEGER NOT NULL DEFAULT 0,
                divorce_requested_at TIMESTAMPTZ NULL,
                ended BOOLEAN NOT NULL DEFAULT FALSE,
                created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
            );
        `);
        await queryInterface.sequelize.query(
            "CREATE UNIQUE INDEX IF NOT EXISTS marriages_groom_active_idx ON marriages (groom_id) WHERE ended = FALSE;"
        );
        await queryInterface.sequelize.query(
            "CREATE UNIQUE INDEX IF NOT EXISTS marriages_bride_active_idx ON marriages (bride_id) WHERE ended = FALSE;"
        );
        await queryInterface.sequelize.query(`
            CREATE TABLE IF NOT EXISTS wedding_gifts (
                gift_id SERIAL PRIMARY KEY,
                receiver_id INTEGER NOT NULL,
                sender_name VARCHAR(32) NOT NULL,
                item JSONB NOT NULL,
                created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
            );
        `);
        await queryInterface.sequelize.query(
            "CREATE INDEX IF NOT EXISTS wedding_gifts_receiver_id_idx ON wedding_gifts (receiver_id);"
        );
        await queryInterface.addColumn("inventory", "marriage_id", {
            type: Sequelize.INTEGER,
            allowNull: false,
            defaultValue: 0,
        });
        await queryInterface.addColumn("storage_items", "marriage_id", {
            type: Sequelize.INTEGER,
            allowNull: false,
            defaultValue: 0,
        });
    },

    async down(queryInterface) {
        await queryInterface.removeColumn("storage_items", "marriage_id");
        await queryInterface.removeColumn("inventory", "marriage_id");
        await queryInterface.sequelize.query("DROP TABLE IF EXISTS wedding_gifts;");
        await queryInterface.sequelize.query("DROP TABLE IF EXISTS marriages;");
    },
};

export default migration;
