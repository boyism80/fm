import type { MigrationModule } from "@fm/shared/sequelize-migration";

// The party row's leader_character_id is the only record of the leader; the member role copied it and nothing read it.
const migration: MigrationModule = {
    async up(queryInterface) {
        await queryInterface.removeColumn("party_members", "role");
    },

    async down(queryInterface, Sequelize) {
        await queryInterface.addColumn("party_members", "role", {
            type: Sequelize.SMALLINT,
            allowNull: false,
            defaultValue: 1,
        });
    },
};

export default migration;
