import { Pool, type PoolClient } from "pg";
import { AsyncLocalStorage } from "node:async_hooks";
import Redis, { type RedisOptions } from "ioredis";
import { getPostgresqlWorld, getRedisWorld } from "../config";
import type { AppConfiguration } from "../config/app-configuration";
import type { InternalConfig, PgEndpoint, RedisEndpoint } from "../types/internal-config";
import type { RepositoryTxOptions } from "../types/repository-contracts";

type TxFn<T> = (client: PoolClient) => Promise<T>;
type DatabaseOperation = { name: string; subject: string };
type PoolLocation = { worldId: number; role: string };

const PG_CONNECTION_TIMEOUT_MS = 10000;

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
    private readonly pgFailureLog: Record<string, Pool>;
    private readonly pgData: Record<string, Pool[]>;
    private readonly redisUnified: Redis | null;
    private readonly redisGlobal: Record<string, Redis[]>;
    private readonly redisData: Record<string, Redis[]>;
    private readonly clientPools = new WeakMap<PoolClient, Pool>();
    private readonly poolLocations = new WeakMap<Pool, PoolLocation>();
    private readonly afterCommit = new WeakMap<PoolClient, Map<string, () => Promise<void>>>();
    private readonly databaseOperation = new AsyncLocalStorage<DatabaseOperation>();

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
        this.pgUnified = this.config.postgresql.unified ? this.createPool(this.config.postgresql.unified, -1, "unified") : null;
        this.pgGlobal = {};
        this.pgFailureLog = {};
        this.pgData = {};
        this.redisUnified = this.config.redis.unified ? InternalContext.createRedis(this.config.redis.unified) : null;
        this.redisGlobal = {};
        this.redisData = {};

        for (const worldKey of Object.keys(this.config.postgresql.worlds)) {
            const worldId = Number(worldKey);
            const worldPg = getPostgresqlWorld(this.config, worldId);
            this.pgGlobal[worldKey] = this.createPool(worldPg.global, worldId, "global");
            this.pgFailureLog[worldKey] = InternalContext.createFailureLogPool(worldPg.global);
            this.pgData[worldKey] = worldPg.data.map((ep, index) => this.createPool(ep, worldId, `data-${index}`));
        }

        for (const worldKey of Object.keys(this.config.redis.worlds)) {
            const worldId = Number(worldKey);
            const worldRedis = getRedisWorld(this.config, worldId);
            this.redisGlobal[worldKey] = [InternalContext.createRedis(worldRedis.global)];
            this.redisData[worldKey] = worldRedis.data.map((ep) => InternalContext.createRedis(ep));
        }
    }

    private createPool(endpoint: PgEndpoint, worldId: number, role: string): Pool {
        const pool = new Pool({
            host: endpoint.host,
            port: endpoint.port,
            database: endpoint.database,
            user: endpoint.username,
            password: endpoint.password,
            ssl: endpoint.ssl ? { rejectUnauthorized: false } : false,
            max: endpoint.pool.max,
            min: endpoint.pool.min,
            connectionTimeoutMillis: PG_CONNECTION_TIMEOUT_MS,
        });
        pool.on("error", (error) => {
            console.error(`[pg] idle client error host=${endpoint.host} database=${endpoint.database}`, error);
        });
        this.poolLocations.set(pool, { worldId, role });
        const query = pool.query.bind(pool) as (statement: string, values?: unknown[]) => Promise<unknown>;
        pool.query = ((statement: string, values?: unknown[]) =>
            this.observeQuery(pool, statement, () => query(statement, values))) as Pool["query"];
        return pool;
    }

    private static createFailureLogPool(endpoint: PgEndpoint): Pool {
        const pool = new Pool({
            host: endpoint.host,
            port: endpoint.port,
            database: endpoint.database,
            user: endpoint.username,
            password: endpoint.password,
            ssl: endpoint.ssl ? { rejectUnauthorized: false } : false,
            max: 1,
            connectionTimeoutMillis: PG_CONNECTION_TIMEOUT_MS,
        });
        pool.on("error", (error) => {
            console.error("[pg] failure log idle client error", error);
        });
        return pool;
    }

    private async observeQuery<T>(pool: Pool, statement: string, query: () => Promise<T>): Promise<T> {
        try {
            return await query();
        } catch (error) {
            await this.recordDatabaseFailure(pool, statement, error);
            throw error;
        }
    }

    private async recordDatabaseFailure(pool: Pool, statement: string, error: unknown): Promise<void> {
        const location = this.poolLocations.get(pool);
        if (!location) {
            return;
        }
        const operation = this.databaseOperation.getStore();
        const failure = error as { code?: string; message?: string };
        const worldId = location.worldId >= 0 ? location.worldId : this.config.app.world_id;
        await this.insertFailure(worldId, location.role, operation?.name ?? "query", operation?.subject ?? "", statement, failure.code ?? "", String(failure.message ?? error));
    }

    private async insertFailure(worldId: number, role: string, operation: string, subject: string, statement: string, code: string, message: string): Promise<void> {
        const global = this.pgFailureLog[String(worldId)];
        if (!global) {
            console.error(`[pg] failure log has no global pool world=${worldId} operation=${operation}`);
            return;
        }
        try {
            await global.query(
                "INSERT INTO db_failures (world_id, database_role, operation, subject, statement, error_code, error_message) VALUES ($1, $2, $3, $4, $5, $6, $7)",
                [worldId, role, operation, subject.slice(0, 4000), statement.slice(0, 4000), code, message.slice(0, 4000)],
            );
        } catch (logError) {
            console.error("[pg] failure log write failed", logError, { worldId, role, operation, subject, message });
        }
    }

    withDatabaseOperation<T>(name: string, subject: string, work: () => Promise<T>): Promise<T> {
        return this.databaseOperation.run({ name, subject }, work);
    }

    deferCacheInvalidation(client: PoolClient, key: string, invalidate: () => Promise<void>): void {
        const pending = this.afterCommit.get(client);
        if (!pending) {
            throw new Error("Cache invalidation requires an active transaction");
        }
        pending.set(key, invalidate);
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

    getPgDataPools(worldId: number): Pool[] {
        const worldKey = String(worldId);
        return this.pgData[worldKey] ?? [];
    }

    async withPgClient<T>(pool: Pool, fn: TxFn<T>): Promise<T> {
        let client: PoolClient;
        try {
            client = await pool.connect();
        } catch (error) {
            await this.recordDatabaseFailure(pool, "CONNECT", error);
            throw error;
        }
        this.clientPools.set(client, pool);
        const originalQuery = client.query;
        const query = client.query.bind(client) as (statement: string, values?: unknown[]) => Promise<unknown>;
        client.query = ((statement: string, values?: unknown[]) =>
            this.observeQuery(pool, statement, () => query(statement, values))) as PoolClient["query"];
        let lost: Error | undefined;
        const onError = (error: Error) => {
            lost = error;
            console.error("[pg] checked-out client error", error);
        };
        client.on("error", onError);
        try {
            return await fn(client);
        } finally {
            client.removeListener("error", onError);
            client.query = originalQuery;
            client.release(lost);
        }
    }

    private async withPgTransaction<T>(pool: Pool, fn: TxFn<T>): Promise<T> {
        return this.withPgClient(pool, async (client) => {
            this.afterCommit.set(client, new Map());
            let result: T;
            try {
                await client.query("BEGIN");
                result = await fn(client);
            } catch (error) {
                this.afterCommit.delete(client);
                await client.query("ROLLBACK").catch(() => {});
                throw error;
            }
            const pending = this.afterCommit.get(client);
            this.afterCommit.delete(client);
            try {
                await client.query("COMMIT");
            } catch (error) {
                await client.query("ROLLBACK").catch(() => {});
                for (const invalidate of pending?.values() ?? []) {
                    await invalidate();
                }
                throw error;
            }
            for (const invalidate of pending?.values() ?? []) {
                await invalidate();
            }
            return result;
        });
    }

    // A transaction client runs a query only on its own pool; another shard's client would hit the wrong database.
    pgClient(pool: Pool, options: RepositoryTxOptions = {}): Pool | PoolClient {
        if (options.txClient != null && this.clientPools.get(options.txClient) === pool) {
            return options.txClient;
        }
        return pool;
    }

    // A transaction already open on the same shard is joined, so the caller never waits for a second connection
    // from a pool it is holding.
    async withPgDataTransaction<T>(worldId: number, hash: number, fn: TxFn<T>, options: RepositoryTxOptions = {}): Promise<T> {
        const pool = this.getPgDataPool(worldId, hash);
        if (options.txClient != null && this.clientPools.get(options.txClient) === pool) {
            return fn(options.txClient);
        }
        return this.withPgTransaction(pool, fn);
    }

    async withPgGlobalTransaction<T>(worldId: number, fn: TxFn<T>): Promise<T> {
        const worldKey = String(worldId);
        const pool = this.pgGlobal[worldKey];
        if (!pool) {
            throw new Error(`Unknown world_id for global transaction: ${worldId}`);
        }
        return this.withPgTransaction(pool, fn);
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

    getRedisUnifiedAccess(): { client: Redis } {
        if (!this.redisUnified) {
            throw new Error("Unified Redis client is not configured (redis.unified missing)");
        }
        return { client: this.redisUnified };
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
            const failureLogPool = this.pgFailureLog[worldId];
            if (failureLogPool) {
                pools.push(failureLogPool);
            }
            for (const pool of worldDataPools) {
                pools.push(pool);
            }
        }
        await Promise.all(pools.map((pool) => pool.end().catch(() => {})));

        const clients: Redis[] = [];
        if (this.redisUnified) {
            clients.push(this.redisUnified);
        }
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
