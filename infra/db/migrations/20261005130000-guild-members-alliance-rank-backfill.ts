import type { MigrationModule } from "../../src/sequelize-migration";

const migration: MigrationModule = {
    async up(queryInterface) {
        await queryInterface.sequelize.query(`
            UPDATE guild_members m
            SET alliance_rank = 2, updated_at = NOW()
            FROM guilds g
            WHERE g.world_id = m.world_id
              AND g.guild_id = m.guild_id
              AND g.alliance_id IS NOT NULL
              AND g.disbanded_at IS NULL
              AND m.guild_rank = 1
              AND (m.alliance_rank IS NULL OR m.alliance_rank = 3);
        `);
        await queryInterface.sequelize.query(`
            UPDATE guild_members m
            SET alliance_rank = 5, updated_at = NOW()
            FROM guilds g
            WHERE g.world_id = m.world_id
              AND g.guild_id = m.guild_id
              AND g.alliance_id IS NOT NULL
              AND g.disbanded_at IS NULL
              AND m.alliance_rank IS NULL;
        `);
    },

    async down() {},
};

export default migration;
