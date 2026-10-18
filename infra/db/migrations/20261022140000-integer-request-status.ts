import type { MigrationModule } from "../../src/sequelize-migration";

const migration: MigrationModule = {
    async up(queryInterface) {
        await queryInterface.sequelize.query(`
            ALTER TABLE requests ALTER COLUMN status DROP DEFAULT;
            ALTER TABLE requests
                ALTER COLUMN status TYPE INTEGER
                USING (
                    CASE status
                        WHEN 'applied' THEN 0
                        WHEN 'aborted' THEN 1
                        WHEN 'completed' THEN 2
                        ELSE 0
                    END
                );
            ALTER TABLE requests ALTER COLUMN status SET DEFAULT 0;

            ALTER TABLE cash_transfer_requests
                ALTER COLUMN kind TYPE INTEGER
                USING (
                    CASE kind
                        WHEN 'gift' THEN 0
                        WHEN 'ring' THEN 1
                        WHEN 'take_out' THEN 2
                        WHEN 'put_in' THEN 3
                        WHEN 'inventory_slot' THEN 4
                        WHEN 'meso' THEN 5
                        ELSE 0
                    END
                );
        `);
    },

    async down(queryInterface) {
        await queryInterface.sequelize.query(`
            ALTER TABLE cash_transfer_requests
                ALTER COLUMN kind TYPE TEXT
                USING (
                    CASE kind
                        WHEN 0 THEN 'gift'
                        WHEN 1 THEN 'ring'
                        WHEN 2 THEN 'take_out'
                        WHEN 3 THEN 'put_in'
                        WHEN 4 THEN 'inventory_slot'
                        WHEN 5 THEN 'meso'
                        ELSE 'gift'
                    END
                );

            ALTER TABLE requests ALTER COLUMN status DROP DEFAULT;
            ALTER TABLE requests
                ALTER COLUMN status TYPE TEXT
                USING (
                    CASE status
                        WHEN 0 THEN 'applied'
                        WHEN 1 THEN 'aborted'
                        WHEN 2 THEN 'completed'
                        ELSE 'applied'
                    END
                );
            ALTER TABLE requests ALTER COLUMN status SET DEFAULT 'applied';
        `);
    },
};

export default migration;
