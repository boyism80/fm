import type { MigrationModule } from "@fm/shared/sequelize-migration";

// A reply belongs to its thread: deleting the thread deletes its replies, and a reply cannot outlive a thread deleted
// while it was being written, whichever request commits first.
const migration: MigrationModule = {
    async up(queryInterface) {
        await queryInterface.sequelize.query(`
            DELETE FROM guild_bulletin_board_replies r
            WHERE NOT EXISTS (
                SELECT 1 FROM guild_bulletin_board_threads t
                WHERE t.world_id = r.world_id AND t.guild_id = r.guild_id AND t.local_thread_id = r.local_thread_id
            );
            ALTER TABLE guild_bulletin_board_replies
                ADD CONSTRAINT fk_guild_bulletin_board_replies_thread
                FOREIGN KEY (world_id, guild_id, local_thread_id)
                REFERENCES guild_bulletin_board_threads (world_id, guild_id, local_thread_id)
                ON DELETE CASCADE;
        `);
    },

    async down(queryInterface) {
        await queryInterface.sequelize.query(`
            ALTER TABLE guild_bulletin_board_replies DROP CONSTRAINT fk_guild_bulletin_board_replies_thread;
        `);
    },
};

export default migration;
