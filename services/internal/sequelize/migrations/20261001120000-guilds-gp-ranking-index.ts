import type { MigrationModule } from "../../src/sequelize-migration";

const migration: MigrationModule = {
    async up(queryInterface) {
        await queryInterface.sequelize.query(`
            CREATE INDEX idx_guilds_world_gp_ranking
            ON guilds (world_id, gp DESC, guild_id ASC)
            WHERE disbanded_at IS NULL;
        `);
    },

    async down(queryInterface) {
        await queryInterface.removeIndex("guilds", "idx_guilds_world_gp_ranking");
    },
};

export default migration;
