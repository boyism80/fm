"use strict";

const fs = require("fs");
const path = require("path");
const { Sequelize } = require("sequelize");

function _endpointToUrl(ep) {
    const user = encodeURIComponent(ep.username ?? "");
    const pass = encodeURIComponent(ep.password ?? "");
    const auth = pass ? `${user}:${pass}` : user;
    return `postgres://${auth}@${ep.host}:${ep.port}/${ep.database}`;
}

function _collectPgEndpoints(cfg) {
    const out = [];
    if (cfg.postgresql.unified) {
        out.push({ name: "unified", endpoint: cfg.postgresql.unified });
    }
    for (const [worldId, world] of Object.entries(cfg.postgresql.worlds ?? {})) {
        if (world.global) {
            out.push({ name: `world-${worldId}-global`, endpoint: world.global });
        }
        for (let i = 0; i < (world.data ?? []).length; i++) {
            out.push({ name: `world-${worldId}-data-${i}`, endpoint: world.data[i] });
        }
    }
    return out;
}

async function _ensureSequelizeMeta(queryInterface) {
    const qi = queryInterface;
    const [rows] = await qi.sequelize.query(
        "SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'SequelizeMeta') AS exists"
    );
    if (!rows[0]?.exists) {
        await qi.createTable("SequelizeMeta", {
            name: {
                type: Sequelize.STRING,
                allowNull: false,
                primaryKey: true,
            },
        });
    }
}

async function _runMigrationsOnEndpoint(baseDir, target) {
    const migrationsDir = path.join(baseDir, "sequelize", "migrations");
    const files = fs.readdirSync(migrationsDir)
        .filter((f) => f.endsWith(".js"))
        .sort();
    const sequelize = new Sequelize(_endpointToUrl(target.endpoint), {
        dialect: "postgres",
        logging: false,
        dialectOptions: target.endpoint.ssl ? { ssl: { require: true, rejectUnauthorized: false } } : {},
    });
    try {
        const qi = sequelize.getQueryInterface();
        await _ensureSequelizeMeta(qi);
        const [appliedRows] = await sequelize.query("SELECT name FROM \"SequelizeMeta\" ORDER BY name");
        const applied = new Set(appliedRows.map((r) => r.name));
        for (const file of files) {
            if (applied.has(file)) continue;
            const migrationPath = path.join(migrationsDir, file);
            delete require.cache[require.resolve(migrationPath)];
            const migration = require(migrationPath);
            await migration.up(qi, Sequelize);
            await sequelize.query("INSERT INTO \"SequelizeMeta\" (name) VALUES ($1)", { bind: [file] });
            console.log(`[auto-migrate] ${target.name}: applied ${file}`);
        }
    } finally {
        await sequelize.close();
    }
}

async function autoMigrateAllIfEnabled(internalConfig) {
    const enabled = Boolean(internalConfig.sequelize?.auto_migrate_on_startup);
    if (!enabled) return;
    const baseDir = path.join(__dirname, "..");
    const targets = _collectPgEndpoints(internalConfig);
    for (const target of targets) {
        await _runMigrationsOnEndpoint(baseDir, target);
    }
}

module.exports = { autoMigrateAllIfEnabled };
