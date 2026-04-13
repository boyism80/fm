"use strict";

const fs = require("fs");
const path = require("path");
const yaml = require("js-yaml");

function withDefaults(raw) {
    const d = raw && typeof raw === "object" ? raw : {};
    return {
        app: {
            log_level: d.app?.log_level ?? "info",
            world_id: d.app?.world_id ?? 1,
        },
        grpc: {
            host: d.grpc?.host ?? "0.0.0.0",
            port: d.grpc?.port ?? 50051,
        },
        postgresql: {
            host: d.postgresql?.host ?? "127.0.0.1",
            port: d.postgresql?.port ?? 5432,
            database: d.postgresql?.database ?? "fm",
            username: d.postgresql?.username ?? "fm",
            password: d.postgresql?.password ?? "",
            ssl: d.postgresql?.ssl ?? false,
            pool: {
                min: d.postgresql?.pool?.min ?? 0,
                max: d.postgresql?.pool?.max ?? 10,
            },
        },
        redis: {
            host: d.redis?.host ?? "127.0.0.1",
            port: d.redis?.port ?? 6379,
            password: d.redis?.password ?? "",
            db: d.redis?.db ?? 0,
            tls: d.redis?.tls ?? false,
            key_prefix: d.redis?.key_prefix ?? "cache:",
        },
        cache: {
            write_strategy: d.cache?.write_strategy ?? "write-through",
        },
        sequelize: {
            dialect: d.sequelize?.dialect ?? "postgres",
            migration_storage: d.sequelize?.migration_storage ?? "sequelize",
            define: d.sequelize?.define ?? { underscored: true },
        },
    };
}

function resolveConfigPath() {
    const envPath = process.env.FM_INTERNAL_CONFIG;
    if (envPath) {
        return envPath;
    }
    return path.join(__dirname, "..", "config.yaml");
}

function loadConfig() {
    const configPath = resolveConfigPath();
    if (!fs.existsSync(configPath)) {
        throw new Error(
            `Internal config not found: ${configPath}. Set FM_INTERNAL_CONFIG or add services/internal/config.yaml`
        );
    }
    const raw = yaml.load(fs.readFileSync(configPath, "utf8"));
    const config = withDefaults(raw);
    return { configPath, ...config };
}

module.exports = { loadConfig, withDefaults, resolveConfigPath };
