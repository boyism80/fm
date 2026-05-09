import fs from "fs";
import path from "path";
import { fileURLToPath } from "url";
import yaml from "js-yaml";
import type {
    GameServerWorldConfig,
    InternalConfig,
    PgEndpoint,
    PgPoolConfig,
    RedisEndpoint,
    SequelizeDefineConfig,
    WorldPgConfig,
    WorldRedisConfig,
} from "./types/internal-config";

const DEFAULT_PG_POOL: PgPoolConfig = { min: 0, max: 10 };
const DEFAULT_REDIS = { password: "", db: 0, tls: false };
const DEFAULT_CHANNEL_NAME_PREFIX = "Channel";
const DEFAULT_RABBITMQ = { ip: "127.0.0.1", port: 5672, uid: "guest", pwd: "guest", vhost: "fm" };
const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

type RawPgEndpoint = Partial<PgEndpoint> & { pool?: Partial<PgPoolConfig> };
type RawRedisEndpoint = Partial<RedisEndpoint>;
type RawWorldPgConfig = { global: RawPgEndpoint; data: RawPgEndpoint[] };
type RawWorldRedisConfig = { global: RawRedisEndpoint; data: RawRedisEndpoint[] };
type RawGameServerChannel = Partial<{ channel_id: number; host: string; port: number; name: string }>;
type RawGameServerWorld = Partial<{ world_name: string; flag: number; event_message: string; channels: RawGameServerChannel[] }>;
type RawSequelize = Partial<Omit<InternalConfig["sequelize"], "define"> & { define: SequelizeDefineConfig }>;
type RawConfig = Partial<{
    app: Partial<InternalConfig["app"]>;
    grpc: Partial<InternalConfig["grpc"]>;
    postgresql: Partial<{ unified: RawPgEndpoint | null; worlds: Record<string, RawWorldPgConfig> }>;
    redis: Partial<{ unified: RawRedisEndpoint | null; worlds: Record<string, RawWorldRedisConfig> }>;
    game_servers: Partial<{ worlds: Record<string, RawGameServerWorld> }>;
    rabbitmq: Partial<InternalConfig["rabbitmq"]>;
    cache: Partial<InternalConfig["cache"]>;
    sequelize: RawSequelize;
    resources: Partial<InternalConfig["resources"]>;
}>;

function normalizePgEndpoint(endpoint?: RawPgEndpoint | null, defaults: Partial<PgEndpoint> = {}): PgEndpoint | null {
    if (!endpoint) {
        return null;
    }
    const poolRaw: Partial<PgPoolConfig> = endpoint.pool ?? {};
    return {
        host: endpoint.host ?? defaults.host ?? "127.0.0.1",
        port: endpoint.port ?? defaults.port ?? 5432,
        database: endpoint.database ?? defaults.database ?? "fm",
        username: endpoint.username ?? defaults.username ?? "fm",
        password: endpoint.password ?? defaults.password ?? "",
        ssl: endpoint.ssl ?? defaults.ssl ?? false,
        pool: {
            min: poolRaw.min ?? DEFAULT_PG_POOL.min,
            max: poolRaw.max ?? DEFAULT_PG_POOL.max,
        },
    };
}

function normalizeRedisEndpoint(endpoint?: RawRedisEndpoint | null, defaults: Partial<RedisEndpoint> = {}): RedisEndpoint | null {
    if (!endpoint) {
        return null;
    }
    return {
        host: endpoint.host ?? defaults.host ?? "127.0.0.1",
        port: endpoint.port ?? defaults.port ?? 6379,
        password: endpoint.password ?? defaults.password ?? DEFAULT_REDIS.password,
        db: endpoint.db ?? defaults.db ?? DEFAULT_REDIS.db,
        tls: endpoint.tls ?? defaults.tls ?? DEFAULT_REDIS.tls,
    };
}

function normalizePostgresqlWorldEntry(entry: RawWorldPgConfig | undefined, worldKey: string): WorldPgConfig {
    if (!entry) {
        throw new Error(`postgresql.worlds[${worldKey}] must be an object with global and data`);
    }
    const global = normalizePgEndpoint(entry.global);
    if (!global) {
        throw new Error(`postgresql.worlds[${worldKey}].global is required`);
    }
    const dataRaw = entry.data ?? [];
    if (dataRaw.length === 0) {
        throw new Error(`postgresql.worlds[${worldKey}].data must be a non-empty array of shard endpoints`);
    }
    const data = dataRaw.map((shard) => {
        const endpoint = normalizePgEndpoint(shard, global);
        if (!endpoint) {
            throw new Error(`postgresql.worlds[${worldKey}].data: invalid shard entry`);
        }
        return endpoint;
    });
    return { global, data };
}

function normalizePostgresql(rawPg: RawConfig["postgresql"], worldId: number): InternalConfig["postgresql"] {
    const pgObject = rawPg ?? {};
    const unified = pgObject.unified ? normalizePgEndpoint(pgObject.unified) : null;
    const worldsIn = pgObject.worlds ?? {};
    const worlds: Record<string, WorldPgConfig> = {};
    for (const key of Object.keys(worldsIn)) {
        worlds[key] = normalizePostgresqlWorldEntry(worldsIn[key], key);
    }
    const worldKey = String(worldId);
    if (!worlds[worldKey]) {
        throw new Error(
            `postgresql.worlds["${worldKey}"] is missing (app.world_id=${worldId}). Define per-world global + data[] shards.`
        );
    }
    return { unified, worlds };
}

function normalizeRedisWorldEntry(entry: RawWorldRedisConfig | undefined, worldKey: string): WorldRedisConfig {
    if (!entry) {
        throw new Error(`redis.worlds[${worldKey}] must be an object with global and data`);
    }
    const global = normalizeRedisEndpoint(entry.global);
    if (!global) {
        throw new Error(`redis.worlds[${worldKey}].global is required`);
    }
    const dataRaw = entry.data ?? [];
    if (dataRaw.length === 0) {
        throw new Error(`redis.worlds[${worldKey}].data must be a non-empty array of shard endpoints`);
    }
    const data = dataRaw.map((shard) => {
        const endpoint = normalizeRedisEndpoint(shard, global);
        if (!endpoint) {
            throw new Error(`redis.worlds[${worldKey}].data: invalid shard entry`);
        }
        return endpoint;
    });
    return { global, data };
}

function normalizeRedis(rawRedis: RawConfig["redis"], worldId: number): InternalConfig["redis"] {
    const redisObject = rawRedis ?? {};
    const unified = redisObject.unified ? normalizeRedisEndpoint(redisObject.unified) : null;
    const worldsIn = redisObject.worlds ?? {};
    const worlds: Record<string, WorldRedisConfig> = {};
    for (const key of Object.keys(worldsIn)) {
        worlds[key] = normalizeRedisWorldEntry(worldsIn[key], key);
    }
    const worldKey = String(worldId);
    if (!worlds[worldKey]) {
        throw new Error(
            `redis.worlds["${worldKey}"] is missing (app.world_id=${worldId}). Define per-world global + data[] shards.`
        );
    }
    return { unified, worlds };
}

function normalizeGameServers(rawGameServers: RawConfig["game_servers"]): InternalConfig["game_servers"] {
    const gameServers = rawGameServers ?? {};
    const worldsIn = gameServers.worlds ?? {};
    const worlds: Record<string, GameServerWorldConfig> = {};
    for (const worldKey of Object.keys(worldsIn)) {
        const world = worldsIn[worldKey] ?? {};
        const channelsRaw = world.channels ?? [];
        const channels = channelsRaw.map((channel, idx) => {
            const ch = channel ?? {};
            if (!ch.host) {
                throw new Error(`game_servers.worlds[${worldKey}].channels[${idx}].host is required`);
            }
            const port = ch.port ?? 0;
            if (!Number.isFinite(port) || port <= 0) {
                throw new Error(`game_servers.worlds[${worldKey}].channels[${idx}].port must be > 0`);
            }
            const channelId = ch.channel_id ?? idx;
            return {
                channel_id: channelId,
                host: ch.host,
                port,
                name: ch.name || `${DEFAULT_CHANNEL_NAME_PREFIX} ${channelId + 1}`,
            };
        });
        worlds[worldKey] = {
            world_name: world.world_name || `World-${worldKey}`,
            flag: world.flag ?? 0,
            event_message: world.event_message ?? "",
            channels,
        };
    }
    return { worlds };
}

export function withDefaults(raw: RawConfig): Omit<InternalConfig, "configPath"> {
    const doc = raw;
    const appObject = doc.app ?? {};
    const app = {
        log_level: appObject.log_level ?? "info",
        world_id: appObject.world_id ?? 0,
        server_alive_ttl_seconds: appObject.server_alive_ttl_seconds ?? 45,
    };
    const worldId = app.world_id;
    const grpcObject = doc.grpc ?? {};
    const rabbitmqObject = doc.rabbitmq ?? {};
    const cacheObject = doc.cache ?? {};
    const sequelizeObject = doc.sequelize ?? {};
    const resourcesObject = doc.resources ?? {};
    return {
        app,
        grpc: {
            host: grpcObject.host ?? "0.0.0.0",
            port: grpcObject.port ?? 50051,
        },
        postgresql: normalizePostgresql(doc.postgresql, worldId),
        redis: normalizeRedis(doc.redis, worldId),
        game_servers: normalizeGameServers(doc.game_servers),
        rabbitmq: {
            ip: rabbitmqObject.ip ?? DEFAULT_RABBITMQ.ip,
            port: rabbitmqObject.port ?? DEFAULT_RABBITMQ.port,
            uid: rabbitmqObject.uid ?? DEFAULT_RABBITMQ.uid,
            pwd: rabbitmqObject.pwd ?? DEFAULT_RABBITMQ.pwd,
            vhost: rabbitmqObject.vhost ?? DEFAULT_RABBITMQ.vhost,
        },
        cache: {
            write_strategy: cacheObject.write_strategy ?? "write-through",
            character_ttl_seconds: cacheObject.character_ttl_seconds ?? 300,
            item_ttl_seconds: cacheObject.item_ttl_seconds ?? 300,
        },
        sequelize: {
            dialect: sequelizeObject.dialect ?? "postgres",
            migration_storage: sequelizeObject.migration_storage ?? "sequelize",
            auto_migrate_on_startup: sequelizeObject.auto_migrate_on_startup ?? false,
            define: sequelizeObject.define ?? {},
        },
        resources: {
            wz_root: resourcesObject.wz_root ?? path.resolve(__dirname, "..", "..", "..", "resources", "wz"),
        },
    };
}

export function getPostgresqlWorld(cfg: Pick<InternalConfig, "app" | "postgresql">, worldId?: number): WorldPgConfig {
    const worldKey = String(worldId ?? cfg.app.world_id);
    const block = cfg.postgresql.worlds[worldKey];
    if (!block) {
        throw new Error(`No postgresql.worlds["${worldKey}"]`);
    }
    return block;
}

export function getRedisWorld(cfg: Pick<InternalConfig, "app" | "redis">, worldId?: number): WorldRedisConfig {
    const worldKey = String(worldId ?? cfg.app.world_id);
    const block = cfg.redis.worlds[worldKey];
    if (!block) {
        throw new Error(`No redis.worlds["${worldKey}"]`);
    }
    return block;
}

export function getPostgresGlobal(cfg: Pick<InternalConfig, "app" | "postgresql">, worldId?: number): PgEndpoint {
    return getPostgresqlWorld(cfg, worldId).global;
}

export function getRedisGlobal(cfg: Pick<InternalConfig, "app" | "redis">, worldId?: number): RedisEndpoint {
    return getRedisWorld(cfg, worldId).global;
}

export function pickPostgresDataShard(worldPg: WorldPgConfig, hash: number): PgEndpoint {
    if (worldPg.data.length === 1) {
        const single = worldPg.data[0];
        if (!single) {
            throw new Error("postgres data shard is empty");
        }
        return single;
    }
    const n = hash;
    const idx = Number.isFinite(n) ? Math.abs(Math.trunc(n)) % worldPg.data.length : 0;
    const shard = worldPg.data[idx];
    if (!shard) {
        throw new Error(`postgres data shard index out of range: ${idx}`);
    }
    return shard;
}

export function pickRedisDataShard(worldRedis: WorldRedisConfig, hash: number): RedisEndpoint {
    if (worldRedis.data.length === 1) {
        const single = worldRedis.data[0];
        if (!single) {
            throw new Error("redis data shard is empty");
        }
        return single;
    }
    const n = hash;
    const idx = Number.isFinite(n) ? Math.abs(Math.trunc(n)) % worldRedis.data.length : 0;
    const shard = worldRedis.data[idx];
    if (!shard) {
        throw new Error(`redis data shard index out of range: ${idx}`);
    }
    return shard;
}

export function resolveConfigPath(): string {
    if (process.env.FM_INTERNAL_CONFIG) {
        return process.env.FM_INTERNAL_CONFIG;
    }
    return path.join(__dirname, "..", "config.yaml");
}

export function loadConfig(): InternalConfig {
    const configPath = resolveConfigPath();
    if (!fs.existsSync(configPath)) {
        throw new Error(`Internal config not found: ${configPath}. Set FM_INTERNAL_CONFIG or add services/internal/config.yaml`);
    }
    const raw = (yaml.load(fs.readFileSync(configPath, "utf8")) ?? {}) as RawConfig;
    return { configPath, ...withDefaults(raw) };
}
