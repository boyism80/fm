import type { MigrationModule } from "../../src/sequelize-migration";

const migration: MigrationModule = {
    async up(queryInterface) {
        await queryInterface.sequelize.query(`
            CREATE TABLE IF NOT EXISTS player_npcs (
                world_id INTEGER NOT NULL,
                npc_id INTEGER NOT NULL,
                character_id INTEGER NOT NULL,
                name TEXT NOT NULL,
                gender SMALLINT NOT NULL,
                skin SMALLINT NOT NULL,
                face INTEGER NOT NULL,
                hair INTEGER NOT NULL,
                equips JSONB NOT NULL DEFAULT '[]'::jsonb,
                map_id INTEGER NOT NULL,
                x INTEGER NOT NULL,
                y INTEGER NOT NULL,
                created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                PRIMARY KEY (world_id, npc_id)
            );
        `);
    },

    async down(queryInterface) {
        await queryInterface.sequelize.query(`DROP TABLE IF EXISTS player_npcs;`);
    },
};

export default migration;
