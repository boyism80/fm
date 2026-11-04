import type { MigrationModule } from "@fm/shared/sequelize-migration";

const migration: MigrationModule = {
    async up(queryInterface) {
        await queryInterface.sequelize.query("UPDATE characters SET hp = 50 WHERE hp < 1");
    },

    async down() {},
};

export default migration;
