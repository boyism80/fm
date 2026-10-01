import { Repository } from "./repository";
import type { RepositoryTxOptions } from "./repository";
import type { InternalContext } from "../context/internal-context";
import type { DistributedLockGuard } from "../system/distributed-lock";

export abstract class ValueRepository<TModel = Record<string, unknown>, TRow = Record<string, unknown>, TKey = unknown> extends Repository<TModel, TRow, TKey> {
    private readonly localCache: Map<string, TRow | null>;

    constructor(internalContext: InternalContext) {
        super(internalContext);
        this.localCache = new Map();
    }

    override async invalidateCache(worldId: number, key: TKey): Promise<void> {
        this.localCache.delete(this.getRedisKey(worldId, key));
        await super.invalidateCache(worldId, key);
    }

    override async get(worldId: number, key: TKey, options: RepositoryTxOptions = {}): Promise<TModel | null> {
        const redis = this.redis(worldId, key);
        const redisKey = this.getRedisKey(worldId, key);
        const useCache = !options.txClient;
        if (this.localCache.has(redisKey)) {
            const cachedRow = this.localCache.get(redisKey);
            return cachedRow == null ? null : this.rowToModel(this.normalizeRow(cachedRow));
        }

        if (useCache) {
            const cached = await redis.get(redisKey);
            if (cached) {
                try {
                    const parsed = JSON.parse(cached) as TRow & { deleted?: boolean };
                    if (parsed.deleted) {
                        this.localCache.set(redisKey, null);
                        return null;
                    }
                    const row = this.normalizeRow(parsed);
                    this.localCache.set(redisKey, row);
                    return this.rowToModel(row);
                } catch {
                    await redis.del(redisKey).catch(() => {});
                }
            }
        }

        const guard = useCache ? await this.tryLockCache(worldId, key) : null;
        try {
            const pool = this.pool(worldId, key);
            const select = this.onSelect(key, worldId);
            const res = await this.query(pool, select.text, select.values, options);
            if (!res.rows.length) {
                this.localCache.set(redisKey, null);
                return null;
            }
            const raw = res.rows[0] as TRow & { deleted?: boolean };
            if (raw.deleted) {
                this.localCache.set(redisKey, null);
                return null;
            }
            const row = this.normalizeRow(raw);

            if (guard) {
                await redis.set(redisKey, JSON.stringify(row), "EX", this.getTtlSeconds()).catch(() => {});
            }
            this.localCache.set(redisKey, row);
            return this.rowToModel(row);
        } finally {
            await guard?.release();
        }
    }

    override async set(worldId: number, model: TModel, options: RepositoryTxOptions = {}): Promise<TModel> {
        const key = this.getKey(model);
        const row = this.modelToRow(model);
        const pool = this.pool(worldId, key);
        const upsert = this.onUpsert(row);
        const res = await this.query(pool, upsert.text, upsert.values, options);
        const savedRow = this.normalizeRow(res.rows[0] as TRow);
        if (!options.txClient) {
            await this.invalidateCache(worldId, key);
        }
        this.localCache.set(this.getRedisKey(worldId, key), savedRow);
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
                    if (cachedRow != null) {
                        results.set(key, this.rowToModel(this.normalizeRow(cachedRow)));
                    }
                    continue;
                }
                if (useCache && cached[i]) {
                    try {
                        const parsed = JSON.parse(cached[i] as string) as TRow & { deleted?: boolean };
                        if (!parsed.deleted) {
                            const row = this.normalizeRow(parsed);
                            this.localCache.set(redisKey, row);
                            results.set(key, this.rowToModel(row));
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
            const guards = new Map<string, DistributedLockGuard>();
            try {
                if (useCache) {
                    for (const key of missedKeys) {
                        const guard = await this.tryLockCache(worldId, key);
                        if (guard) {
                            guards.set(this.getRedisKey(worldId, key), guard);
                        }
                    }
                }

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

                for (const raw of rows) {
                    if (!raw.deleted) {
                        const row = this.normalizeRow(raw);
                        const model = this.rowToModel(row);
                        const key = this.getKey(model);
                        const redisKey = this.getRedisKey(worldId, key);
                        if (guards.has(redisKey)) {
                            await this.redis(worldId, key).set(redisKey, JSON.stringify(row), "EX", this.getTtlSeconds()).catch(() => {});
                        }
                        this.localCache.set(redisKey, row);
                        results.set(key, model);
                    }
                }
            } finally {
                for (const guard of guards.values()) {
                    await guard.release();
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

            for (const raw of savedRows) {
                const savedRow = this.normalizeRow(raw);
                const model = this.rowToModel(savedRow);
                const key = this.getKey(model);
                if (!options.txClient) {
                    await this.invalidateCache(worldId, key);
                }
                this.localCache.set(this.getRedisKey(worldId, key), savedRow);
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
        if (deleted && !options.txClient) {
            await this.invalidateCache(worldId, key);
        }
        if (deleted) {
            this.localCache.set(this.getRedisKey(worldId, key), null);
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
