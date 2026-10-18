import type { MigrationModule } from "../../src/sequelize-migration";

const DEFAULT_GUILD_RANK_TITLES = ["Master", "Jr. Master", "Member", "Member", "Member"];

function sqlTextArrayLiteral(values: string[]): string {
    return `ARRAY[${values.map((t) => `'${t.replace(/'/g, "''")}'`).join(", ")}]::text[]`;
}

const DEFAULT_GUILD_RANK_TITLES_SQL = sqlTextArrayLiteral(DEFAULT_GUILD_RANK_TITLES);

const migration: MigrationModule = {
    async up(queryInterface) {
        await queryInterface.sequelize.query(`
            ALTER TABLE guilds
            ALTER COLUMN rank_titles DROP DEFAULT;
        `);

        await queryInterface.sequelize.query(`
            ALTER TABLE guilds
            ALTER COLUMN rank_titles TYPE TEXT[]
            USING (
                CASE
                    WHEN rank_titles IS NULL THEN ${DEFAULT_GUILD_RANK_TITLES_SQL}
                    WHEN jsonb_typeof(rank_titles::jsonb) = 'array' THEN ARRAY[
                        COALESCE(rank_titles::jsonb->>0, 'Master'),
                        COALESCE(rank_titles::jsonb->>1, 'Jr. Master'),
                        COALESCE(rank_titles::jsonb->>2, 'Member'),
                        COALESCE(rank_titles::jsonb->>3, 'Member'),
                        COALESCE(rank_titles::jsonb->>4, 'Member')
                    ]::text[]
                    ELSE ${DEFAULT_GUILD_RANK_TITLES_SQL}
                END
            );
        `);

        await queryInterface.sequelize.query(`
            ALTER TABLE guilds
            ALTER COLUMN rank_titles SET NOT NULL,
            ALTER COLUMN rank_titles SET DEFAULT ARRAY['Master','Jr. Master','Member','Member','Member']::text[];
        `);
    },

    async down(queryInterface) {
        await queryInterface.sequelize.query(`
            ALTER TABLE guilds
            ALTER COLUMN rank_titles DROP DEFAULT;
        `);

        await queryInterface.sequelize.query(`
            ALTER TABLE guilds
            ALTER COLUMN rank_titles TYPE JSONB
            USING to_jsonb(rank_titles);
        `);

        await queryInterface.sequelize.query(`
            ALTER TABLE guilds
            ALTER COLUMN rank_titles SET DEFAULT '["Master","Jr. Master","Member","Member","Member"]'::jsonb;
        `);
    },
};

export default migration;
