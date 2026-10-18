import type { QueryInterface } from "sequelize";

export interface MigrationModule {
    up: (queryInterface: QueryInterface, sequelizeModule: typeof import("sequelize")) => Promise<void>;
    down: (queryInterface: QueryInterface, sequelizeModule: typeof import("sequelize")) => Promise<void>;
}
