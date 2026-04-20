"use strict";

/**
 * Sequelize CLI config. Run migrations per shard, e.g.:
 *   set DATABASE_URL=postgres://fm:admin@127.0.0.1:5432/fm && npx sequelize-cli db:migrate
 */
module.exports = {
    development: {
        url: process.env.DATABASE_URL || "postgres://fm:admin@127.0.0.1:5432/fm",
        dialect: "postgres",
        dialectOptions:
            process.env.PGSSL === "1" ? { ssl: { require: true, rejectUnauthorized: false } } : {},
    },
    production: {
        url: process.env.DATABASE_URL,
        dialect: "postgres",
        dialectOptions:
            process.env.PGSSL === "1" ? { ssl: { require: true, rejectUnauthorized: false } } : {},
    },
};
