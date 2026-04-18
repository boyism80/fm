"use strict";

class Repository {
    constructor(internalContext) {
        this.ctx = internalContext;
    }

    getKey(_model) { throw new Error(`${this.constructor.name}.getKey not implemented`); }
    onSelect(_key, _worldId) { throw new Error(`${this.constructor.name}.onSelect not implemented`); }
    onUpsert(_row) { throw new Error(`${this.constructor.name}.onUpsert not implemented`); }
    onDelete(_row) { throw new Error(`${this.constructor.name}.onDelete not implemented`); }
    getRedisKey(_worldId, _key) { throw new Error(`${this.constructor.name}.getRedisKey not implemented`); }
    rowToModel(_row) { throw new Error(`${this.constructor.name}.rowToModel not implemented`); }
    modelToRow(_model) { throw new Error(`${this.constructor.name}.modelToRow not implemented`); }

    getTtlSeconds() { return 300; }
    getShardHash(key) { return key; }
    onSelectMany(_keys, _worldId) { return null; }
    onBulkUpsert(_rows) { return null; }
    getDeleteKey(row) { return this.getKey(row); }

    _pool(worldId, key) {
        return this.ctx.getPgDataPool(worldId, this.getShardHash(key));
    }

    _query(pool, text, values, options = undefined) {
        const txClient = options?.txClient;
        if (txClient) {
            return txClient.query(text, values);
        }
        return pool.query(text, values);
    }

    _redis(worldId, key) {
        return this.ctx.getRedisDataAccess(worldId, this.getShardHash(key)).client;
    }

    _groupByPgShard(worldId, keys) {
        const groups = new Map();
        for (const key of keys) {
            const pool = this._pool(worldId, key);
            if (!groups.has(pool)) groups.set(pool, []);
            groups.get(pool).push(key);
        }
        return groups;
    }

    _groupModelsByPgShard(worldId, models) {
        const groups = new Map();
        for (const model of models) {
            const pool = this._pool(worldId, this.getKey(model));
            if (!groups.has(pool)) groups.set(pool, []);
            groups.get(pool).push(model);
        }
        return groups;
    }

    _groupByRedisShard(worldId, keys) {
        const groups = new Map();
        for (const key of keys) {
            const { client } = this.ctx.getRedisDataAccess(worldId, this.getShardHash(key));
            if (!groups.has(client)) groups.set(client, { client, keys: [] });
            groups.get(client).keys.push(key);
        }
        return [...groups.values()];
    }

    async evictCache(worldId, key) {
        const redis = this._redis(worldId, key);
        await redis.del(this.getRedisKey(worldId, key)).catch(() => {});
    }

    async get(worldId, key, options = undefined) {
        const redis = this._redis(worldId, key);
        const redisKey = this.getRedisKey(worldId, key);
        const useCache = !options?.txClient;

        if (useCache) {
            const cached = await redis.get(redisKey);
            if (cached) {
                try {
                    const row = JSON.parse(cached);
                    if (row.deleted) return null;
                    return this.rowToModel(row);
                } catch {
                    await redis.del(redisKey).catch(() => {});
                }
            }
        }

        const pool = this._pool(worldId, key);
        const { text, values } = this.onSelect(key, worldId);
        const res = await this._query(pool, text, values, options);
        if (!res.rows.length) return null;
        const row = res.rows[0];
        if (row.deleted) return null;

        if (useCache) {
            await redis.set(redisKey, JSON.stringify(row), "EX", this.getTtlSeconds()).catch(() => {});
        }
        return this.rowToModel(row);
    }

    async set(worldId, model, options = undefined) {
        const key = this.getKey(model);
        const row = this.modelToRow(model);
        const pool = this._pool(worldId, key);
        const { text, values } = this.onUpsert(row);
        const res = await this._query(pool, text, values, options);
        const savedRow = res.rows[0];
        if (!options?.txClient) {
            const redis = this._redis(worldId, key);
            const redisKey = this.getRedisKey(worldId, key);
            await redis.set(redisKey, JSON.stringify(savedRow), "EX", this.getTtlSeconds());
        }
        return this.rowToModel(savedRow);
    }

    async getMany(worldId, keys, options = undefined) {
        const results = new Map();
        const useCache = !options?.txClient;
        if (options?.txClient) {
            const pgGroups = this._groupByPgShard(worldId, keys);
            if (pgGroups.size > 1) {
                throw new Error(`${this.constructor.name}.getMany with txClient supports only single data shard`);
            }
        }
        const redisGroups = this._groupByRedisShard(worldId, keys);
        const dbMissByPool = new Map();

        for (const { client, keys: groupKeys } of redisGroups) {
            let cached = [];
            if (useCache) {
                const redisKeys = groupKeys.map(k => this.getRedisKey(worldId, k));
                cached = await client.mget(...redisKeys);
            }

            for (let i = 0; i < groupKeys.length; i++) {
                const key = groupKeys[i];
                if (useCache && cached[i]) {
                    try {
                        const row = JSON.parse(cached[i]);
                        if (!row.deleted) results.set(key, this.rowToModel(row));
                    } catch {
                        // corrupted cache entry -> fall through to DB
                    }
                }
                if (!results.has(key)) {
                    const pool = this._pool(worldId, key);
                    if (!dbMissByPool.has(pool)) dbMissByPool.set(pool, []);
                    dbMissByPool.get(pool).push(key);
                }
            }
        }

        for (const [pool, missedKeys] of dbMissByPool) {
            const bulkQuery = this.onSelectMany(missedKeys, worldId);
            let rows;

            if (bulkQuery) {
                const res = await this._query(pool, bulkQuery.text, bulkQuery.values, options);
                rows = res.rows;
            } else {
                rows = [];
                for (const key of missedKeys) {
                    const { text, values } = this.onSelect(key, worldId);
                    const res = await this._query(pool, text, values, options);
                    rows.push(...res.rows);
                }
            }

            for (const row of rows) {
                if (!row.deleted) {
                    const model = this.rowToModel(row);
                    const key = this.getKey(model);
                    if (useCache) {
                        const redis = this._redis(worldId, key);
                        const redisKey = this.getRedisKey(worldId, key);
                        await redis.set(redisKey, JSON.stringify(row), "EX", this.getTtlSeconds()).catch(() => {});
                    }
                    results.set(key, model);
                }
            }
        }

        return results;
    }

    async setAll(worldId, models, options = undefined) {
        const saved = [];
        const pgGroups = this._groupModelsByPgShard(worldId, models);
        if (options?.txClient && pgGroups.size > 1) {
            throw new Error(`${this.constructor.name}.setAll with txClient supports only single data shard`);
        }

        for (const [pool, groupModels] of pgGroups) {
            const rows = groupModels.map(m => this.modelToRow(m));
            const bulkQuery = this.onBulkUpsert(rows);
            let savedRows;

            if (bulkQuery) {
                const res = await this._query(pool, bulkQuery.text, bulkQuery.values, options);
                savedRows = res.rows;
            } else {
                savedRows = [];
                for (const row of rows) {
                    const { text, values } = this.onUpsert(row);
                    const res = await this._query(pool, text, values, options);
                    savedRows.push(...res.rows);
                }
            }

            for (const savedRow of savedRows) {
                const model = this.rowToModel(savedRow);
                const key = this.getKey(model);
                if (!options?.txClient) {
                    const redis = this._redis(worldId, key);
                    const redisKey = this.getRedisKey(worldId, key);
                    await redis.set(redisKey, JSON.stringify(savedRow), "EX", this.getTtlSeconds());
                }
                saved.push(model);
            }
        }

        return saved;
    }

    async delete(row, options = undefined) {
        const worldId = row.worldId;
        const key = this.getDeleteKey(row);
        const pool = this._pool(worldId, key);
        const { text, values } = this.onDelete(row);
        const res = await this._query(pool, text, values, options);
        const deleted = Number(res.rowCount || 0) > 0;
        if (deleted && !options?.txClient) {
            const redis = this._redis(worldId, key);
            await redis.del(this.getRedisKey(worldId, key)).catch(() => {});
        }
        return deleted;
    }
}

module.exports = { Repository };
