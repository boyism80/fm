"use strict";

const fs = require("fs");
const path = require("path");
const yaml = require("js-yaml");

const DEFAULT_PG_POOL = { min: 0, max: 10 };

const DEFAULT_REDIS = {
    password: "",
    db: 0,
    tls: false,
    key_prefix: "cache:",
};

const DEFAULT_CHANNEL_NAME_PREFIX = "Channel";
const DEFAULT_RABBITMQ = {
    ip: "127.0.0.1",
    port: 5672,
    uid: "guest",
    pwd: "guest",
    vhost: "fm",
};

function normalizePgEndpoint(ep, defaults = {}) {
    if (!ep || typeof ep !== "object") {
        return null;
    }
    const poolRaw = ep.pool ?? {};
    return {
        host: ep.host ?? defaults.host ?? "127.0.0.1",
        port: ep.port ?? defaults.port ?? 5432,
        database: ep.database ?? defaults.database ?? "fm",
        username: ep.username ?? defaults.username ?? "fm",
        password: ep.password ?? defaults.password ?? "",
        ssl: ep.ssl ?? defaults.ssl ?? false,
        pool: {
            min: poolRaw.min ?? DEFAULT_PG_POOL.min,
            max: poolRaw.max ?? DEFAULT_PG_POOL.max,
        },
    };
}

function normalizeRedisEndpoint(ep, defaults = {}) {
    if (!ep || typeof ep !== "object") {
        return null;
    }
    return {
        host: ep.host ?? defaults.host ?? "127.0.0.1",
        port: ep.port ?? defaults.port ?? 6379,
        password: ep.password ?? defaults.password ?? DEFAULT_REDIS.password,
        db: ep.db ?? defaults.db ?? DEFAULT_REDIS.db,
        tls: ep.tls ?? defaults.tls ?? DEFAULT_REDIS.tls,
        key_prefix: ep.key_prefix ?? defaults.key_prefix ?? DEFAULT_REDIS.key_prefix,
    };
}

/** Legacy: flat postgresql.host without worlds / unified wrapper. */
function isLegacyPostgresql(p) {
    return p && typeof p === "object" && typeof p.host === "string" && !p.worlds;
}

/** Legacy: flat redis.host without worlds. */
function isLegacyRedis(r) {
    return r && typeof r === "object" && typeof r.host === "string" && !r.worlds;
}

function migrateLegacyPostgresql(p, worldIdStr) {
    const base = normalizePgEndpoint(p);
    return {
        unified: null,
        worlds: {
            [worldIdStr]: {
                global: { ...base },
                data: [{ ...base }],
            },
        },
    };
}

function migrateLegacyRedis(r, worldIdStr) {
    const base = normalizeRedisEndpoint(r);
    return {
        unified: null,
        worlds: {
            [worldIdStr]: {
                global: { ...base },
                data: [{ ...base }],
            },
        },
    };
}

function normalizePostgresqlWorldEntry(entry, worldKey) {
    if (!entry || typeof entry !== "object") {
        throw new Error(`postgresql.worlds[${worldKey}] must be an object with global and data`);
    }
    const global = normalizePgEndpoint(entry.global);
    if (!global) {
        throw new Error(`postgresql.worlds[${worldKey}].global is required`);
    }
    const dataRaw = entry.data;
    if (!Array.isArray(dataRaw) || dataRaw.length === 0) {
        throw new Error(`postgresql.worlds[${worldKey}].data must be a non-empty array of shard endpoints`);
    }
    const data = dataRaw.map((shard) => {
        const s = normalizePgEndpoint(shard, global);
        if (!s) {
            throw new Error(`postgresql.worlds[${worldKey}].data: invalid shard entry`);
        }
        return s;
    });
    return { global, data };
}

function normalizePostgresql(rawPg, worldIdStr) {
    let pg = rawPg;
    if (isLegacyPostgresql(pg)) {
        pg = migrateLegacyPostgresql(pg, worldIdStr);
    }
    if (!pg || typeof pg !== "object") {
        pg = {};
    }
    const unified = pg.unified ? normalizePgEndpoint(pg.unified) : null;
    const worldsIn = pg.worlds && typeof pg.worlds === "object" ? pg.worlds : {};
    const worlds = {};
    for (const key of Object.keys(worldsIn)) {
        worlds[key] = normalizePostgresqlWorldEntry(worldsIn[key], key);
    }
    if (!worlds[worldIdStr]) {
        throw new Error(
            `postgresql.worlds["${worldIdStr}"] is missing (app.world_id=${worldIdStr}). ` +
                "Define per-world global + data[] shards, or use the legacy flat postgresql block for a single node."
        );
    }
    return { unified, worlds };
}

function normalizeRedisWorldEntry(entry, worldKey) {
    if (!entry || typeof entry !== "object") {
        throw new Error(`redis.worlds[${worldKey}] must be an object with global and data`);
    }
    const global = normalizeRedisEndpoint(entry.global);
    if (!global) {
        throw new Error(`redis.worlds[${worldKey}].global is required`);
    }
    const dataRaw = entry.data;
    if (!Array.isArray(dataRaw) || dataRaw.length === 0) {
        throw new Error(`redis.worlds[${worldKey}].data must be a non-empty array of shard endpoints`);
    }
    const data = dataRaw.map((shard) => {
        const s = normalizeRedisEndpoint(shard, global);
        if (!s) {
            throw new Error(`redis.worlds[${worldKey}].data: invalid shard entry`);
        }
        return s;
    });
    return { global, data };
}

function normalizeRedis(rawRedis, worldIdStr) {
    let r = rawRedis;
    if (isLegacyRedis(r)) {
        r = migrateLegacyRedis(r, worldIdStr);
    }
    if (!r || typeof r !== "object") {
        r = {};
    }
    const unified = r.unified ? normalizeRedisEndpoint(r.unified) : null;
    const worldsIn = r.worlds && typeof r.worlds === "object" ? r.worlds : {};
    const worlds = {};
    for (const key of Object.keys(worldsIn)) {
        worlds[key] = normalizeRedisWorldEntry(worldsIn[key], key);
    }
    if (!worlds[worldIdStr]) {
        throw new Error(
            `redis.worlds["${worldIdStr}"] is missing (app.world_id=${worldIdStr}). ` +
                "Define per-world global + data[] shards, or use the legacy flat redis block."
        );
    }
    return { unified, worlds };
}

function normalizeGameServers(rawGameServers) {
    const gameServers = rawGameServers && typeof rawGameServers === "object" ? rawGameServers : {};
    const worldsIn = gameServers.worlds && typeof gameServers.worlds === "object" ? gameServers.worlds : {};
    const worlds = {};
    for (const worldKey of Object.keys(worldsIn)) {
        const world = worldsIn[worldKey] || {};
        const channelsRaw = Array.isArray(world.channels) ? world.channels : [];
        const channels = channelsRaw.map((ch, idx) => {
            if (!ch || typeof ch !== "object") {
                throw new Error(`game_servers.worlds[${worldKey}].channels[${idx}] must be an object`);
            }
            if (typeof ch.host !== "string" || !ch.host) {
                throw new Error(`game_servers.worlds[${worldKey}].channels[${idx}].host is required`);
            }
            const port = Number(ch.port);
            if (!Number.isFinite(port) || port <= 0) {
                throw new Error(`game_servers.worlds[${worldKey}].channels[${idx}].port must be > 0`);
            }
            const channelId = Number.isFinite(Number(ch.channel_id)) ? Number(ch.channel_id) : idx;
            return {
                channel_id: channelId,
                host: ch.host,
                port,
                name: ch.name || `${DEFAULT_CHANNEL_NAME_PREFIX} ${channelId + 1}`,
            };
        });
        worlds[worldKey] = {
            world_name: world.world_name || `World-${worldKey}`,
            flag: Number.isFinite(Number(world.flag)) ? Number(world.flag) : 0,
            event_message: world.event_message || "",
            channels,
        };
    }
    return { worlds };
}

function withDefaults(raw) {
    const d = raw && typeof raw === "object" ? raw : {};
    const app = {
        log_level: d.app?.log_level ?? "info",
        world_id: d.app?.world_id ?? 0,
    };
    const worldIdStr = String(app.world_id);

    return {
        app,
        grpc: {
            host: d.grpc?.host ?? "0.0.0.0",
            port: d.grpc?.port ?? 50051,
        },
        postgresql: normalizePostgresql(d.postgresql, worldIdStr),
        redis: normalizeRedis(d.redis, worldIdStr),
        game_servers: normalizeGameServers(d.game_servers),
        rabbitmq: {
            ip: d.rabbitmq?.ip ?? DEFAULT_RABBITMQ.ip,
            port: d.rabbitmq?.port ?? DEFAULT_RABBITMQ.port,
            uid: d.rabbitmq?.uid ?? DEFAULT_RABBITMQ.uid,
            pwd: d.rabbitmq?.pwd ?? DEFAULT_RABBITMQ.pwd,
            vhost: d.rabbitmq?.vhost ?? DEFAULT_RABBITMQ.vhost,
        },
        cache: {
            write_strategy: d.cache?.write_strategy ?? "write-through",
            character_ttl_seconds: d.cache?.character_ttl_seconds ?? 300,
            item_ttl_seconds: d.cache?.item_ttl_seconds ?? 300,
        },
        sequelize: {
            dialect: d.sequelize?.dialect ?? "postgres",
            migration_storage: d.sequelize?.migration_storage ?? "sequelize",
            auto_migrate_on_startup: d.sequelize?.auto_migrate_on_startup ?? false,
            define: d.sequelize?.define ?? { underscored: true },
        },
    };
}

function getPostgresqlWorld(cfg, worldId) {
    const w = String(worldId ?? cfg.app.world_id);
    const block = cfg.postgresql.worlds[w];
    if (!block) {
        throw new Error(`No postgresql.worlds["${w}"]`);
    }
    return block;
}

function getRedisWorld(cfg, worldId) {
    const w = String(worldId ?? cfg.app.world_id);
    const block = cfg.redis.worlds[w];
    if (!block) {
        throw new Error(`No redis.worlds["${w}"]`);
    }
    return block;
}

function getPostgresGlobal(cfg, worldId) {
    return getPostgresqlWorld(cfg, worldId).global;
}

function getRedisGlobal(cfg, worldId) {
    return getRedisWorld(cfg, worldId).global;
}

/**
 * Pick a PostgreSQL data shard from a stable hash (e.g. account id). Matches fb-style routing.
 * @param {number|string|null|undefined} hash
 */
function pickPostgresDataShard(worldPg, hash) {
    const { data } = worldPg;
    if (data.length === 1) {
        return data[0];
    }
    const h = Number(hash);
    const idx = Number.isFinite(h) ? Math.abs(Math.trunc(h)) % data.length : 0;
    return data[idx];
}

/**
 * @param {number|string|null|undefined} hash
 */
function pickRedisDataShard(worldRedis, hash) {
    const { data } = worldRedis;
    if (data.length === 1) {
        return data[0];
    }
    const h = Number(hash);
    const idx = Number.isFinite(h) ? Math.abs(Math.trunc(h)) % data.length : 0;
    return data[idx];
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

module.exports = {
    loadConfig,
    withDefaults,
    resolveConfigPath,
    getPostgresqlWorld,
    getRedisWorld,
    getPostgresGlobal,
    getRedisGlobal,
    pickPostgresDataShard,
    pickRedisDataShard,
};
