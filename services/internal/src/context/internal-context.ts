import { Pool, type PoolClient } from "pg";
import Redis, { type RedisOptions } from "ioredis";
import { getPostgresqlWorld, getRedisWorld } from "../config";
import type { AppConfiguration } from "../config/app-configuration";
import type { InternalConfig, PgEndpoint, RedisEndpoint } from "../types/internal-config";

type TxFn<T> = (client: PoolClient) => Promise<T>;

function shardIndex(hash: number, shardCount: number): number {
    if (shardCount <= 1) {
        return 0;
    }
    const n = hash;
    if (!Number.isFinite(n)) {
        return 0;
    }
    return Math.abs(Math.trunc(n)) % shardCount;
}

export class InternalContext {
    readonly appConfiguration: AppConfiguration;
    private readonly config: Pick<InternalConfig, "app" | "grpc" | "postgresql" | "redis" | "cache" | "sequelize">;
    private readonly pgUnified: Pool | null;
    private readonly pgGlobal: Record<string, Pool>;
    private readonly pgData: Record<string, Pool[]>;
    private readonly redisGlobal: Record<string, Redis[]>;
    private readonly redisData: Record<string, Redis[]>;

    constructor(appConfiguration: AppConfiguration) {
        this.appConfiguration = appConfiguration;
        const cfg = appConfiguration.raw;
        this.config = {
            app: cfg.app,
            grpc: cfg.grpc,
            postgresql: cfg.postgresql,
            redis: cfg.redis,
            cache: cfg.cache,
            sequelize: cfg.sequelize,
        };
        this.pgUnified = this.config.postgresql.unified ? InternalContext.createPool(this.config.postgresql.unified) : null;
        this.pgGlobal = {};
        this.pgData = {};
        this.redisGlobal = {};
        this.redisData = {};

        for (const worldKey of Object.keys(this.config.postgresql.worlds)) {
            const worldId = Number(worldKey);
            const worldPg = getPostgresqlWorld(this.config, worldId);
            this.pgGlobal[worldKey] = InternalContext.createPool(worldPg.global);
            this.pgData[worldKey] = worldPg.data.map((ep) => InternalContext.createPool(ep));
        }

        for (const worldKey of Object.keys(this.config.redis.worlds)) {
            const worldId = Number(worldKey);
            const worldRedis = getRedisWorld(this.config, worldId);
            this.redisGlobal[worldKey] = [InternalContext.createRedis(worldRedis.global)];
            this.redisData[worldKey] = worldRedis.data.map((ep) => InternalContext.createRedis(ep));
        }
    }

    private static createPool(endpoint: PgEndpoint): Pool {
        return new Pool({
            host: endpoint.host,
            port: endpoint.port,
            database: endpoint.database,
            user: endpoint.username,
            password: endpoint.password,
            ssl: endpoint.ssl ? { rejectUnauthorized: false } : false,
            max: endpoint.pool.max,
            min: endpoint.pool.min,
        });
    }

    private static createRedis(endpoint: RedisEndpoint): Redis {
        const options: RedisOptions = {
            host: endpoint.host,
            port: endpoint.port,
            password: endpoint.password || undefined,
            db: endpoint.db,
            maxRetriesPerRequest: 2,
        };
        if (endpoint.tls) {
            options.tls = {};
        }
        return new Redis(options);
    }

    getPgDataPool(worldId: number, hash: number): Pool {
        const worldKey = String(worldId);
        const worldPg = getPostgresqlWorld(this.config, worldId);
        const idx = shardIndex(hash, worldPg.data.length);
        const pool = this.pgData[worldKey]?.[idx];
        if (!pool) {
            throw new Error(`No pg data shard for world ${worldKey} index ${idx}`);
        }
        return pool;
    }

    async withPgDataTransaction<T>(worldId: number, hash: number, fn: TxFn<T>): Promise<T> {
        const pool = this.getPgDataPool(worldId, hash);
        const client = await pool.connect();
        try {
            await client.query("BEGIN");
            const result = await fn(client);
            await client.query("COMMIT");
            return result;
        } catch (error) {
            try {
                await client.query("ROLLBACK");
            } catch {
            }
            throw error;
        } finally {
            client.release();
        }
    }

    async withPgGlobalTransaction<T>(worldId: number, fn: TxFn<T>): Promise<T> {
        const worldKey = String(worldId);
        const pool = this.pgGlobal[worldKey];
        if (!pool) {
            throw new Error(`Unknown world_id for global transaction: ${worldId}`);
        }
        const client = await pool.connect();
        try {
            await client.query("BEGIN");
            const result = await fn(client);
            await client.query("COMMIT");
            return result;
        } catch (error) {
            try {
                await client.query("ROLLBACK");
            } catch {
            }
            throw error;
        } finally {
            client.release();
        }
    }

    getPgUnifiedPool(): Pool {
        if (!this.pgUnified) {
            throw new Error("Unified PostgreSQL pool is not configured (postgresql.unified missing)");
        }
        return this.pgUnified;
    }

    getRedisDataAccess(worldId: number, hash: number): { client: Redis } {
        const worldKey = String(worldId);
        const worldRedis = getRedisWorld(this.config, worldId);
        const idx = shardIndex(hash, worldRedis.data.length);
        const client = this.redisData[worldKey]?.[idx];
        if (!client) {
            throw new Error(`No redis data shard for world ${worldKey} index ${idx}`);
        }
        return { client };
    }

    getRedisGlobalAccess(worldId: number): { client: Redis } {
        const worldKey = String(worldId);
        const client = this.redisGlobal[worldKey]?.[0];
        if (!client) {
            throw new Error(`Unknown world_id for redis global access: ${worldId}`);
        }
        return { client };
    }

    async close(): Promise<void> {
        const pools: Pool[] = [];
        if (this.pgUnified) {
            pools.push(this.pgUnified);
        }
        for (const worldId of Object.keys(this.pgGlobal)) {
            const worldGlobalPool = this.pgGlobal[worldId];
            if (worldGlobalPool) {
                pools.push(worldGlobalPool);
            }
            const worldDataPools = this.pgData[worldId] ?? [];
            for (const pool of worldDataPools) {
                pools.push(pool);
            }
        }
        await Promise.all(pools.map((pool) => pool.end().catch(() => {})));

        const clients: Redis[] = [];
        for (const worldId of Object.keys(this.redisData)) {
            const worldGlobalClients = this.redisGlobal[worldId] ?? [];
            for (const client of worldGlobalClients) {
                clients.push(client);
            }
            const worldDataClients = this.redisData[worldId] ?? [];
            for (const client of worldDataClients) {
                clients.push(client);
            }
        }
        for (const client of clients) {
            try {
                client.disconnect();
            } catch {
            }
        }
    }
}
