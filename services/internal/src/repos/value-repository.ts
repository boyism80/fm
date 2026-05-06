import { Repository } from "./repository";
import type { RepositoryTxOptions } from "./repository";
import type { InternalContext } from "../context/internal-context";

export abstract class ValueRepository<TModel = Record<string, unknown>, TRow = Record<string, unknown>, TKey = unknown> extends Repository<TModel, TRow, TKey> {
    private readonly localCache: Map<string, TRow | null>;

    constructor(internalContext: InternalContext) {
        super(internalContext);
        this.localCache = new Map();
    }

    protected logL1(method: string, result: "hit" | "miss", key: string): void {
        console.log(`[L1][${this.constructor.name}][${method}] ${result} key=${key}`);
    }

    override async evictCache(worldId: number, key: TKey): Promise<void> {
        const redisKey = this.getRedisKey(worldId, key);
        this.localCache.delete(redisKey);
        const redis = this.redis(worldId, key);
        await redis.del(redisKey).catch(() => {});
    }

    override async get(worldId: number, key: TKey, options: RepositoryTxOptions = {}): Promise<TModel | null> {
        const redis = this.redis(worldId, key);
        const redisKey = this.getRedisKey(worldId, key);
        const useCache = !options.txClient;
        if (this.localCache.has(redisKey)) {
            const cachedRow = this.localCache.get(redisKey);
            this.logL1("get", "hit", redisKey);
            return cachedRow == null ? null : this.rowToModel(cachedRow);
        }
        this.logL1("get", "miss", redisKey);

        if (useCache) {
            const cached = await redis.get(redisKey);
            if (cached) {
                try {
                    const row = JSON.parse(cached) as TRow & { deleted?: boolean };
                    if (row.deleted) {
                        this.localCache.set(redisKey, null);
                        return null;
                    }
                    this.localCache.set(redisKey, row as TRow);
                    return this.rowToModel(row as TRow);
                } catch {
                    await redis.del(redisKey).catch(() => {});
                }
            }
        }

        const pool = this.pool(worldId, key);
        const select = this.onSelect(key, worldId);
        const res = await this.query(pool, select.text, select.values, options);
        if (!res.rows.length) {
            this.localCache.set(redisKey, null);
            return null;
        }
        const row = res.rows[0] as TRow & { deleted?: boolean };
        if (row.deleted) {
            this.localCache.set(redisKey, null);
            return null;
        }

        if (useCache) {
            await redis.set(redisKey, JSON.stringify(row), "EX", this.getTtlSeconds()).catch(() => {});
        }
        this.localCache.set(redisKey, row as TRow);
        return this.rowToModel(row as TRow);
    }

    override async set(worldId: number, model: TModel, options: RepositoryTxOptions = {}): Promise<TModel> {
        const key = this.getKey(model);
        const row = this.modelToRow(model);
        const pool = this.pool(worldId, key);
        const upsert = this.onUpsert(row);
        const res = await this.query(pool, upsert.text, upsert.values, options);
        const savedRow = res.rows[0] as TRow;
        this.localCache.set(this.getRedisKey(worldId, key), savedRow);
        if (!options.txClient) {
            const redis = this.redis(worldId, key);
            const redisKey = this.getRedisKey(worldId, key);
            await redis.set(redisKey, JSON.stringify(savedRow), "EX", this.getTtlSeconds());
        }
        return this.rowToModel(savedRow);
    }

    override async getMany(worldId: number, keys: TKey[], options: RepositoryTxOptions = {}): Promise<Map<TKey, TModel>> {
        const results = new Map<TKey, TModel>();
        const useCache = !options.txClient;
        if (options.txClient) {
            const pgGroups = this.groupByPgShard(worldId, keys);
            if (pgGroups.size > 1) {
                throw new Error(`${this.constructor.name}.getMany with txClient supports only single data shard`);
            }
        }
        const redisGroups = this.groupByRedisShard(worldId, keys);
        const dbMissByPool = new Map<ReturnType<typeof this.pool>, TKey[]>();

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
                const redisKey = this.getRedisKey(worldId, key);
                if (this.localCache.has(redisKey)) {
                    const cachedRow = this.localCache.get(redisKey);
                    this.logL1("getMany", "hit", redisKey);
                    if (cachedRow != null) {
                        results.set(key, this.rowToModel(cachedRow));
                    }
                    continue;
                }
                this.logL1("getMany", "miss", redisKey);
                if (useCache && cached[i]) {
                    try {
                        const row = JSON.parse(cached[i] as string) as TRow & { deleted?: boolean };
                        if (!row.deleted) {
                            this.localCache.set(redisKey, row as TRow);
                            results.set(key, this.rowToModel(row as TRow));
                        } else {
                            this.localCache.set(redisKey, null);
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
                    this.localCache.set(this.getRedisKey(worldId, key), row as TRow);
                    results.set(key, model);
                }
            }
        }
        return results;
    }

    override async setAll(worldId: number, models: TModel[], options: RepositoryTxOptions = {}): Promise<TModel[]> {
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
                this.localCache.set(this.getRedisKey(worldId, key), savedRow);
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

    override async delete(row: { worldId: number } & Record<string, unknown>, options: RepositoryTxOptions = {}): Promise<boolean> {
        const worldId = row.worldId;
        const key = this.getDeleteKey(row);
        const pool = this.pool(worldId, key);
        const del = this.onDelete(row);
        const res = await this.query(pool, del.text, del.values, options);
        const deleted = (res.rowCount ?? 0) > 0;
        const redisKey = this.getRedisKey(worldId, key);
        if (deleted && !options.txClient) {
            const redis = this.redis(worldId, key);
            await redis.del(redisKey).catch(() => {});
        }
        if (deleted) {
            this.localCache.set(redisKey, null);
        }
        return deleted;
    }

    override async getAll(_worldId: number, _key: TKey, _options: RepositoryTxOptions = {}): Promise<never> {
        throw new Error(`${this.constructor.name}.getAll is not supported for value repositories`);
    }

    override async delAll(
        _worldId: number,
        _key: TKey,
        _itemKeys: Array<string | number>,
        _options: RepositoryTxOptions = {}
    ): Promise<never> {
        throw new Error(`${this.constructor.name}.delAll is not supported for value repositories`);
    }
}
