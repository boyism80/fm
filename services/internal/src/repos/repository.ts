import type { Pool, QueryResult } from "pg";
import type Redis from "ioredis";
import type { InternalContext } from "../context/internal-context";
import type { RepositoryQuery, RepositoryTxOptions } from "../types/repository-contracts";

export type { RepositoryQuery, RepositoryTxOptions };

export abstract class Repository<TModel = Record<string, unknown>, TRow = Record<string, unknown>, TKey = unknown> {
    protected readonly ctx: InternalContext;

    constructor(internalContext: InternalContext) {
        this.ctx = internalContext;
    }

    abstract getKey(_model: TModel): TKey;
    abstract onSelect(_key: TKey, _worldId: number): RepositoryQuery;
    abstract onUpsert(_row: TRow): RepositoryQuery;
    abstract onDelete(_row: unknown): RepositoryQuery;
    abstract getRedisKey(_worldId: number, _key: TKey): string;
    abstract rowToModel(_row: TRow): TModel;
    abstract modelToRow(_model: TModel): TRow;
    getTtlSeconds(): number {
        return 300;
    }

    getShardHash(key: TKey): number {
        return Number(key);
    }

    onSelectMany(_keys: TKey[], _worldId: number): RepositoryQuery | null {
        return null;
    }

    onBulkUpsert(_rows: TRow[]): RepositoryQuery | null {
        return null;
    }

    getDeleteKey(row: unknown): TKey {
        return this.getKey(row as TModel);
    }

    protected pool(worldId: number, key: TKey): Pool {
        return this.ctx.getPgDataPool(worldId, this.getShardHash(key));
    }

    protected query(pool: Pool, text: string, values: unknown[], options: RepositoryTxOptions = {}): Promise<QueryResult> {
        const txClient = options.txClient;
        if (txClient) {
            return txClient.query(text, values);
        }
        return pool.query(text, values);
    }

    protected redis(worldId: number, key: TKey): Redis {
        return this.ctx.getRedisDataAccess(worldId, this.getShardHash(key)).client;
    }

    protected groupByPgShard(worldId: number, keys: TKey[]): Map<Pool, TKey[]> {
        const groups = new Map<Pool, TKey[]>();
        for (const key of keys) {
            const pool = this.pool(worldId, key);
            if (!groups.has(pool)) {
                groups.set(pool, []);
            }
            groups.get(pool)?.push(key);
        }
        return groups;
    }

    protected groupModelsByPgShard(worldId: number, models: TModel[]): Map<Pool, TModel[]> {
        const groups = new Map<Pool, TModel[]>();
        for (const model of models) {
            const pool = this.pool(worldId, this.getKey(model));
            if (!groups.has(pool)) {
                groups.set(pool, []);
            }
            groups.get(pool)?.push(model);
        }
        return groups;
    }

    protected groupByRedisShard(worldId: number, keys: TKey[]): Array<{ client: Redis; keys: TKey[] }> {
        const groups = new Map<Redis, { client: Redis; keys: TKey[] }>();
        for (const key of keys) {
            const { client } = this.ctx.getRedisDataAccess(worldId, this.getShardHash(key));
            if (!groups.has(client)) {
                groups.set(client, { client, keys: [] });
            }
            groups.get(client)?.keys.push(key);
        }
        return [...groups.values()];
    }

    async evictCache(worldId: number, key: TKey): Promise<void> {
        const redis = this.redis(worldId, key);
        await redis.del(this.getRedisKey(worldId, key)).catch(() => {});
    }

    async get(worldId: number, key: TKey, options: RepositoryTxOptions = {}): Promise<TModel | null> {
        const redis = this.redis(worldId, key);
        const redisKey = this.getRedisKey(worldId, key);
        const useCache = !options.txClient;

        if (useCache) {
            const cached = await redis.get(redisKey);
            if (cached) {
                try {
                    const row = JSON.parse(cached) as TRow & { deleted?: boolean };
                    if (row.deleted) {
                        return null;
                    }
                    return this.rowToModel(row);
                } catch {
                    await redis.del(redisKey).catch(() => {});
                }
            }
        }

        const pool = this.pool(worldId, key);
        const select = this.onSelect(key, worldId);
        const res = await this.query(pool, select.text, select.values, options);
        if (!res.rows.length) {
            return null;
        }
        const row = res.rows[0] as TRow & { deleted?: boolean };
        if (row.deleted) {
            return null;
        }

        if (useCache) {
            await redis.set(redisKey, JSON.stringify(row), "EX", this.getTtlSeconds()).catch(() => {});
        }
        return this.rowToModel(row);
    }

    async set(worldId: number, model: TModel, options: RepositoryTxOptions = {}): Promise<TModel> {
        const key = this.getKey(model);
        const row = this.modelToRow(model);
        const pool = this.pool(worldId, key);
        const upsert = this.onUpsert(row);
        const res = await this.query(pool, upsert.text, upsert.values, options);
        const savedRow = res.rows[0] as TRow;
        if (!options.txClient) {
            const redis = this.redis(worldId, key);
            const redisKey = this.getRedisKey(worldId, key);
            await redis.set(redisKey, JSON.stringify(savedRow), "EX", this.getTtlSeconds());
        }
        return this.rowToModel(savedRow);
    }

    async getMany(worldId: number, keys: TKey[], options: RepositoryTxOptions = {}): Promise<Map<TKey, TModel>> {
        const results = new Map<TKey, TModel>();
        const useCache = !options.txClient;
        if (options.txClient) {
            const pgGroups = this.groupByPgShard(worldId, keys);
            if (pgGroups.size > 1) {
                throw new Error(`${this.constructor.name}.getMany with txClient supports only single data shard`);
            }
        }
        const redisGroups = this.groupByRedisShard(worldId, keys);
        const dbMissByPool = new Map<Pool, TKey[]>();

        for (const { client, keys: groupKeys } of redisGroups) {
            let cached: Array<string | null> = [];
            if (useCache) {
                const redisKeys = groupKeys.map((k) => this.getRedisKey(worldId, k));
                cached = await client.mget(...redisKeys);
            }

            for (let i = 0; i < groupKeys.length; i += 1) {
                const key = groupKeys[i];
                if (key === undefined) {
                    continue;
                }
                if (useCache && cached[i]) {
                    try {
                        const row = JSON.parse(cached[i] as string) as TRow & { deleted?: boolean };
                        if (!row.deleted) {
                            results.set(key, this.rowToModel(row));
                        }
                    } catch {
                    }
                }
                if (!results.has(key)) {
                    const pool = this.pool(worldId, key);
                    if (!dbMissByPool.has(pool)) {
                        dbMissByPool.set(pool, []);
                    }
                    dbMissByPool.get(pool)?.push(key);
                }
            }
        }

        for (const [pool, missedKeys] of dbMissByPool) {
            const bulkQuery = this.onSelectMany(missedKeys, worldId);
            let rows: Array<TRow & { deleted?: boolean }>;
            if (bulkQuery) {
                const res = await this.query(pool, bulkQuery.text, bulkQuery.values, options);
                rows = res.rows as Array<TRow & { deleted?: boolean }>;
            } else {
                rows = [];
                for (const key of missedKeys) {
                    const select = this.onSelect(key, worldId);
                    const res = await this.query(pool, select.text, select.values, options);
                    rows.push(...(res.rows as Array<TRow & { deleted?: boolean }>));
                }
            }

            for (const row of rows) {
                if (!row.deleted) {
                    const model = this.rowToModel(row);
                    const key = this.getKey(model);
                    if (useCache) {
                        const redis = this.redis(worldId, key);
                        const redisKey = this.getRedisKey(worldId, key);
                        await redis.set(redisKey, JSON.stringify(row), "EX", this.getTtlSeconds()).catch(() => {});
                    }
                    results.set(key, model);
                }
            }
        }
        return results;
    }

    async setAll(worldId: number, models: TModel[], options: RepositoryTxOptions = {}): Promise<TModel[]> {
        const saved: TModel[] = [];
        const pgGroups = this.groupModelsByPgShard(worldId, models);
        if (options.txClient && pgGroups.size > 1) {
            throw new Error(`${this.constructor.name}.setAll with txClient supports only single data shard`);
        }

        for (const [pool, groupModels] of pgGroups) {
            const rows = groupModels.map((m) => this.modelToRow(m));
            const bulkQuery = this.onBulkUpsert(rows);
            let savedRows: TRow[];
            if (bulkQuery) {
                const res = await this.query(pool, bulkQuery.text, bulkQuery.values, options);
                savedRows = res.rows as TRow[];
            } else {
                savedRows = [];
                for (const row of rows) {
                    const upsert = this.onUpsert(row);
                    const res = await this.query(pool, upsert.text, upsert.values, options);
                    savedRows.push(...(res.rows as TRow[]));
                }
            }

            for (const savedRow of savedRows) {
                const model = this.rowToModel(savedRow);
                const key = this.getKey(model);
                if (!options.txClient) {
                    const redis = this.redis(worldId, key);
                    const redisKey = this.getRedisKey(worldId, key);
                    await redis.set(redisKey, JSON.stringify(savedRow), "EX", this.getTtlSeconds());
                }
                saved.push(model);
            }
        }
        return saved;
    }

    async delete(row: { worldId: number } & Record<string, unknown>, options: RepositoryTxOptions = {}): Promise<boolean> {
        const worldId = row.worldId;
        const key = this.getDeleteKey(row);
        const pool = this.pool(worldId, key);
        const del = this.onDelete(row);
        const res = await this.query(pool, del.text, del.values, options);
        const deleted = (res.rowCount ?? 0) > 0;
        if (deleted && !options.txClient) {
            const redis = this.redis(worldId, key);
            await redis.del(this.getRedisKey(worldId, key)).catch(() => {});
        }
        return deleted;
    }
}
