import type { MigrationModule } from "@fm/shared/sequelize-migration";

const migration: MigrationModule = {
    async up(q, Sequelize) {
        await q.addColumn("characters", "family_progress", {
            type: Sequelize.BIGINT,
            allowNull: false,
            defaultValue: 0,
        });
        await q.addColumn("characters", "population_history", {
            type: Sequelize.JSONB,
            allowNull: false,
            defaultValue: { last_given_at: null, given: {} },
        });
        await q.sequelize.query(`
            CREATE TABLE social_notes (
                id SERIAL PRIMARY KEY,
                receiver_id INTEGER NOT NULL,
                sender_id INTEGER NOT NULL,
                sender_name VARCHAR(13) NOT NULL,
                message VARCHAR(200) NOT NULL,
                gift SMALLINT NOT NULL DEFAULT 0,
                created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
            );
            CREATE INDEX social_notes_receiver ON social_notes (receiver_id, id);

            CREATE TABLE social_thanks (
                serial TEXT PRIMARY KEY,
                character_id INTEGER NOT NULL,
                sender_id INTEGER NOT NULL,
                consumed BOOLEAN NOT NULL DEFAULT FALSE
            );
            CREATE INDEX social_thanks_receiver ON social_thanks (character_id, sender_id) WHERE NOT consumed;

            CREATE TABLE social_family (
                character_id INTEGER PRIMARY KEY,
                senior_id INTEGER NOT NULL DEFAULT 0,
                name VARCHAR(13) NOT NULL,
                level INTEGER NOT NULL,
                class_id INTEGER NOT NULL,
                reputation INTEGER NOT NULL DEFAULT 0,
                total_reputation INTEGER NOT NULL DEFAULT 0,
                notice VARCHAR(200) NOT NULL DEFAULT '',
                progress BIGINT NOT NULL DEFAULT 0,
                buff_until TIMESTAMPTZ,
                buff_feature INTEGER NOT NULL DEFAULT 0
            );
            CREATE INDEX social_family_senior ON social_family (senior_id);

            CREATE TABLE social_family_invites (
                receiver_id INTEGER PRIMARY KEY,
                sender_id INTEGER NOT NULL,
                channel_id INTEGER NOT NULL,
                expires_at TIMESTAMPTZ NOT NULL
            );

            CREATE TABLE social_family_summons (
                receiver_id INTEGER PRIMARY KEY,
                sender_id INTEGER NOT NULL,
                channel_id INTEGER NOT NULL,
                expires_at TIMESTAMPTZ NOT NULL
            );

            CREATE TABLE social_family_uses (
                character_id INTEGER NOT NULL,
                feature INTEGER NOT NULL,
                day DATE NOT NULL,
                operation_id TEXT NOT NULL,
                PRIMARY KEY (character_id, feature, day)
            );
        `);
    },

    async down(q) {
        await q.sequelize.query(`
            DROP TABLE
                social_family_uses,
                social_family_summons,
                social_family_invites,
                social_family,
                social_thanks,
                social_notes;
        `);
        await q.removeColumn("characters", "population_history");
        await q.removeColumn("characters", "family_progress");
    },
};

export default migration;
