import type { MigrationModule } from "../../src/sequelize-migration";

const migration: MigrationModule = {
    async up(queryInterface, Sequelize) {
        await queryInterface.sequelize.query(`
            CREATE SEQUENCE party_id_seq AS BIGINT START WITH 1 INCREMENT BY 1 MINVALUE 1 NO MAXVALUE CACHE 1;
        `);
        const rows = await queryInterface.sequelize.query(`SELECT COALESCE(MAX(party_id), 0)::bigint AS m FROM parties`, {
            type: Sequelize.QueryTypes.SELECT,
        });
        const m = (rows[0] as { m?: string | number } | undefined)?.m != null ? String((rows[0] as { m: string | number }).m) : "0";
        const isCalled = Number(m) > 0;
        await queryInterface.sequelize.query("SELECT setval('party_id_seq', $1::bigint, $2)", {
            bind: [m, isCalled],
        });
    },

    async down(queryInterface) {
        await queryInterface.sequelize.query("DROP SEQUENCE IF EXISTS party_id_seq");
    },
};

export default migration;
