import type { MigrationModule } from "@fm/shared/sequelize-migration";

const migration: MigrationModule = {
    async up(q, Sequelize) {
        await q.addColumn("characters", "family_reputation", {
            type: Sequelize.JSONB,
            allowNull: false,
            defaultValue: {},
        });
        await q.sequelize.query(`
            CREATE TABLE social_family_credit (
                character_id INTEGER NOT NULL,
                receiver_id INTEGER NOT NULL,
                reputation BIGINT NOT NULL,
                PRIMARY KEY (character_id, receiver_id)
            );
            CREATE TABLE social_family_moves (
                operation_id TEXT PRIMARY KEY,
                character_id INTEGER NOT NULL,
                payer_id INTEGER NOT NULL,
                feature INTEGER NOT NULL,
                status TEXT NOT NULL DEFAULT 'pending',
                expires_at TIMESTAMPTZ NOT NULL DEFAULT NOW()+INTERVAL '30 seconds'
            );
            CREATE INDEX social_family_moves_pending ON social_family_moves (expires_at) WHERE status='pending';
        `);
    },

    async down(q) {
        await q.sequelize.query("DROP TABLE social_family_moves, social_family_credit");
        await q.removeColumn("characters", "family_reputation");
    },
};

export default migration;
