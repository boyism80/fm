import type { MigrationModule } from "../../src/sequelize-migration";

async function syncSequenceFromMax(
    queryInterface: import("sequelize").QueryInterface,
    Sequelize: typeof import("sequelize"),
    sequenceName: string,
    maxSql: string,
    countSql: string
) {
    await queryInterface.sequelize.query(
        `ALTER SEQUENCE ${sequenceName} MINVALUE 1`
    );
    const maxRows = await queryInterface.sequelize.query(
        `SELECT GREATEST(COALESCE((${maxSql}), 0), 1)::bigint AS m`,
        { type: Sequelize.QueryTypes.SELECT }
    );
    const m = String((maxRows[0] as { m?: string | number })?.m ?? 1);
    const countRows = await queryInterface.sequelize.query(countSql, {
        type: Sequelize.QueryTypes.SELECT,
    });
    const count = Number((countRows[0] as { c?: string | number })?.c ?? 0);
    const isCalled = count > 0;
    await queryInterface.sequelize.query(`SELECT setval('${sequenceName}', $1::bigint, $2)`, {
        bind: [m, isCalled],
    });
}

const migration: MigrationModule = {
    async up(queryInterface, Sequelize) {
        await queryInterface.sequelize.query(`
            DO $$
            DECLARE r RECORD;
            DECLARE new_id BIGINT;
            BEGIN
                FOR r IN SELECT world_id, party_id FROM parties WHERE party_id < 1 LOOP
                    new_id := nextval('party_id_seq');
                    UPDATE party_members
                       SET party_id = new_id
                     WHERE world_id = r.world_id AND party_id = r.party_id;
                    UPDATE character_realtime_state
                       SET party_id = new_id
                     WHERE world_id = r.world_id AND party_id = r.party_id;
                    UPDATE parties
                       SET party_id = new_id
                     WHERE world_id = r.world_id AND party_id = r.party_id;
                END LOOP;
            END $$;
        `);

        await queryInterface.sequelize.query(`
            DO $$
            DECLARE r RECORD;
            DECLARE new_id BIGINT;
            BEGIN
                FOR r IN SELECT guild_id FROM guild_name_registry WHERE guild_id < 1 LOOP
                    new_id := nextval(pg_get_serial_sequence('guild_name_registry', 'guild_id'));
                    UPDATE guild_bulletin_board_replies
                       SET guild_id = new_id
                     WHERE guild_id = r.guild_id;
                    UPDATE guild_bulletin_board_threads
                       SET guild_id = new_id
                     WHERE guild_id = r.guild_id;
                    UPDATE guild_members SET guild_id = new_id WHERE guild_id = r.guild_id;
                    UPDATE character_realtime_state SET guild_id = new_id WHERE guild_id = r.guild_id;
                    UPDATE guilds SET guild_id = new_id WHERE guild_id = r.guild_id;
                    UPDATE guild_name_registry SET guild_id = new_id WHERE guild_id = r.guild_id;
                END LOOP;
            END $$;
        `);

        await queryInterface.sequelize.query(`
            DO $$
            DECLARE r RECORD;
            DECLARE new_id BIGINT;
            BEGIN
                FOR r IN SELECT world_id, alliance_id FROM alliances WHERE alliance_id < 1 LOOP
                    new_id := nextval('alliance_id_seq');
                    UPDATE guilds
                       SET alliance_id = new_id
                     WHERE world_id = r.world_id AND alliance_id = r.alliance_id;
                    UPDATE alliances
                       SET alliance_id = new_id
                     WHERE world_id = r.world_id AND alliance_id = r.alliance_id;
                END LOOP;
            END $$;
        `);

        await syncSequenceFromMax(
            queryInterface,
            Sequelize,
            "party_id_seq",
            "SELECT MAX(party_id) FROM parties",
            "SELECT COUNT(*) FROM parties"
        );

        await syncSequenceFromMax(
            queryInterface,
            Sequelize,
            "alliance_id_seq",
            "SELECT MAX(alliance_id) FROM alliances",
            "SELECT COUNT(*) FROM alliances"
        );

        const guildSeqRows = await queryInterface.sequelize.query(
            `SELECT pg_get_serial_sequence('guild_name_registry', 'guild_id') AS seq`,
            { type: Sequelize.QueryTypes.SELECT }
        );
        const guildSeq = (guildSeqRows[0] as { seq?: string | null })?.seq;
        if (guildSeq) {
            await queryInterface.sequelize.query(`ALTER SEQUENCE ${guildSeq} MINVALUE 1`);
            const maxRows = await queryInterface.sequelize.query(
                `SELECT GREATEST(
                    COALESCE((SELECT MAX(guild_id) FROM guild_name_registry), 0),
                    COALESCE((SELECT MAX(guild_id) FROM guilds), 0),
                    1
                )::bigint AS m`,
                { type: Sequelize.QueryTypes.SELECT }
            );
            const m = String((maxRows[0] as { m?: string | number })?.m ?? 1);
            const countRows = await queryInterface.sequelize.query(
                "SELECT COUNT(*) FROM guild_name_registry",
                { type: Sequelize.QueryTypes.SELECT }
            );
            const count = Number((countRows[0] as { c?: string | number })?.c ?? 0);
            await queryInterface.sequelize.query(`SELECT setval($1::regclass, $2::bigint, $3)`, {
                bind: [guildSeq, m, count > 0],
            });
        }

        await queryInterface.sequelize.query(`
            ALTER TABLE parties
                ADD CONSTRAINT chk_parties_party_id_min_one CHECK (party_id >= 1);
            ALTER TABLE guild_name_registry
                ADD CONSTRAINT chk_guild_name_registry_guild_id_min_one CHECK (guild_id >= 1);
            ALTER TABLE guilds
                ADD CONSTRAINT chk_guilds_guild_id_min_one CHECK (guild_id >= 1);
            ALTER TABLE alliances
                ADD CONSTRAINT chk_alliances_alliance_id_min_one CHECK (alliance_id >= 1);
            ALTER TABLE guilds
                ADD CONSTRAINT chk_guilds_alliance_id_min_one
                CHECK (alliance_id IS NULL OR alliance_id >= 1);
        `);
    },

    async down(queryInterface) {
        await queryInterface.sequelize.query(`
            ALTER TABLE guilds DROP CONSTRAINT IF EXISTS chk_guilds_alliance_id_min_one;
            ALTER TABLE alliances DROP CONSTRAINT IF EXISTS chk_alliances_alliance_id_min_one;
            ALTER TABLE guilds DROP CONSTRAINT IF EXISTS chk_guilds_guild_id_min_one;
            ALTER TABLE guild_name_registry DROP CONSTRAINT IF EXISTS chk_guild_name_registry_guild_id_min_one;
            ALTER TABLE parties DROP CONSTRAINT IF EXISTS chk_parties_party_id_min_one;
        `);
    },
};

export default migration;
