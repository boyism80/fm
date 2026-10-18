import type { MigrationModule } from "../../src/sequelize-migration";

const TABLES = [
    { table: "character_records", owner: "character_id" },
    { table: "account_records", owner: "account_id" },
];

const migration: MigrationModule = {
    async up(queryInterface, Sequelize) {
        for (const { table, owner } of TABLES) {
            await queryInterface.createTable(table, {
                [owner]: {
                    type: Sequelize.INTEGER,
                    allowNull: false,
                },
                record_key: {
                    type: Sequelize.TEXT,
                    allowNull: false,
                },
                value: {
                    type: Sequelize.BIGINT,
                    allowNull: false,
                    defaultValue: 0,
                },
                text: {
                    type: Sequelize.TEXT,
                    allowNull: false,
                    defaultValue: "",
                },
                period: {
                    type: Sequelize.SMALLINT,
                    allowNull: false,
                    defaultValue: 0,
                },
                updated_at_unix_ms: {
                    type: Sequelize.BIGINT,
                    allowNull: false,
                    defaultValue: 0,
                },
                updated_at: {
                    type: Sequelize.DATE,
                    allowNull: false,
                    defaultValue: Sequelize.literal("NOW()"),
                },
            });
            await queryInterface.addConstraint(table, {
                fields: [owner, "record_key"],
                type: "primary key",
                name: `pk_${table}`,
            });
            await queryInterface.addIndex(table, {
                fields: [owner],
                name: `idx_${table}_owner`,
            });
        }
    },

    async down(queryInterface) {
        for (const { table } of TABLES) {
            await queryInterface.dropTable(table);
        }
    },
};

export default migration;
