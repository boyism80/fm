import type { QueryInterface } from "sequelize";

export async function up(queryInterface: QueryInterface) {
    await queryInterface.removeColumn("party_members", "channel_index");
}

export async function down(queryInterface: QueryInterface) {
    await queryInterface.addColumn("party_members", "channel_index", {
        type: "INTEGER",
        allowNull: false,
        defaultValue: -2,
    });
}
