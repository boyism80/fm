import type { MigrationModule } from "@fm/shared/sequelize-migration";

const migration: MigrationModule = {
    async up(q) {
        await q.sequelize.query(`
            CREATE TABLE families (
                leader_id INTEGER PRIMARY KEY,
                precept VARCHAR(200) NOT NULL DEFAULT ''
            );
            ALTER TABLE social_family ADD COLUMN leader_id INTEGER;

            WITH RECURSIVE tree AS (
                SELECT character_id, character_id AS root FROM social_family WHERE senior_id = 0
                UNION ALL
                SELECT f.character_id, tree.root FROM social_family f JOIN tree ON f.senior_id = tree.character_id
            )
            UPDATE social_family f SET leader_id = tree.root
            FROM tree
            WHERE f.character_id = tree.character_id
              AND (f.senior_id <> 0 OR EXISTS (SELECT 1 FROM social_family j WHERE j.senior_id = f.character_id));

            INSERT INTO families (leader_id, precept)
            SELECT character_id, notice FROM social_family WHERE leader_id = character_id;

            ALTER TABLE social_family ADD CONSTRAINT social_family_leader_fk FOREIGN KEY (leader_id) REFERENCES families (leader_id);
            CREATE INDEX social_family_leader ON social_family (leader_id) WHERE leader_id IS NOT NULL;
            ALTER TABLE social_family DROP COLUMN notice;
            DROP TABLE social_family_invites, social_family_summons;
        `);
    },

    async down(q) {
        await q.sequelize.query(`
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
            ALTER TABLE social_family ADD COLUMN notice VARCHAR(200) NOT NULL DEFAULT '';

            UPDATE social_family f SET notice = families.precept
            FROM families
            WHERE f.character_id = families.leader_id;

            DROP INDEX social_family_leader;
            ALTER TABLE social_family DROP COLUMN leader_id;
            DROP TABLE families;
        `);
    },
};

export default migration;
