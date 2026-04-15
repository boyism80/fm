"use strict";

class BaseRepository {
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
            const { client, keyPrefix } = this.ctx.getRedisDataAccess(worldId, this.getShardHash(key));
            if (!groups.has(client)) groups.set(client, { client, keyPrefix, keys: [] });
            groups.get(client).keys.push(key);
        }
        return [...groups.values()];
    }

    async get(worldId, key) {
        const redis = this._redis(worldId, key);
        const redisKey = this.getRedisKey(worldId, key);

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

        const pool = this._pool(worldId, key);
        const { text, values } = this.onSelect(key, worldId);
        const res = await pool.query(text, values);
        if (!res.rows.length) return null;
        const row = res.rows[0];
        if (row.deleted) return null;

        await redis.set(redisKey, JSON.stringify(row), "EX", this.getTtlSeconds()).catch(() => {});
        return this.rowToModel(row);
    }

    async set(worldId, model) {
        const key = this.getKey(model);
        const row = this.modelToRow(model);
        const pool = this._pool(worldId, key);
        const { text, values } = this.onUpsert(row);
        const res = await pool.query(text, values);
        const savedRow = res.rows[0];
        const redis = this._redis(worldId, key);
        const redisKey = this.getRedisKey(worldId, key);
        await redis.set(redisKey, JSON.stringify(savedRow), "EX", this.getTtlSeconds());
        return this.rowToModel(savedRow);
    }

    async getMany(worldId, keys) {
        const results = new Map();
        const redisGroups = this._groupByRedisShard(worldId, keys);
        const dbMissByPool = new Map();

        for (const { client, keys: groupKeys } of redisGroups) {
            const redisKeys = groupKeys.map(k => this.getRedisKey(worldId, k));
            const cached = await client.mget(...redisKeys);

            for (let i = 0; i < groupKeys.length; i++) {
                const key = groupKeys[i];
                if (cached[i]) {
                    try {
                        const row = JSON.parse(cached[i]);
                        if (!row.deleted) results.set(key, this.rowToModel(row));
                    } catch {
                        // corrupted cache entry → fall through to DB
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
                const res = await pool.query(bulkQuery.text, bulkQuery.values);
                rows = res.rows;
            } else {
                rows = [];
                for (const key of missedKeys) {
                    const { text, values } = this.onSelect(key, worldId);
                    const res = await pool.query(text, values);
                    rows.push(...res.rows);
                }
            }

            for (const row of rows) {
                if (!row.deleted) {
                    const model = this.rowToModel(row);
                    const key = this.getKey(model);
                    const redis = this._redis(worldId, key);
                    const redisKey = this.getRedisKey(worldId, key);
                    await redis.set(redisKey, JSON.stringify(row), "EX", this.getTtlSeconds()).catch(() => {});
                    results.set(key, model);
                }
            }
        }

        return results;
    }

    async setAll(worldId, models) {
        const saved = [];
        const pgGroups = this._groupModelsByPgShard(worldId, models);

        for (const [pool, groupModels] of pgGroups) {
            const rows = groupModels.map(m => this.modelToRow(m));
            const bulkQuery = this.onBulkUpsert(rows);
            let savedRows;

            if (bulkQuery) {
                const res = await pool.query(bulkQuery.text, bulkQuery.values);
                savedRows = res.rows;
            } else {
                savedRows = [];
                for (const row of rows) {
                    const { text, values } = this.onUpsert(row);
                    const res = await pool.query(text, values);
                    savedRows.push(...res.rows);
                }
            }

            for (const savedRow of savedRows) {
                const model = this.rowToModel(savedRow);
                const key = this.getKey(model);
                const redis = this._redis(worldId, key);
                const redisKey = this.getRedisKey(worldId, key);
                await redis.set(redisKey, JSON.stringify(savedRow), "EX", this.getTtlSeconds());
                saved.push(model);
            }
        }

        return saved;
    }

    async delete(row) {
        const worldId = row.worldId;
        const key = this.getDeleteKey(row);
        const pool = this._pool(worldId, key);
        const { text, values } = this.onDelete(row);
        const res = await pool.query(text, values);
        const deleted = Number(res.rowCount || 0) > 0;
        if (deleted) {
            const redis = this._redis(worldId, key);
            await redis.del(this.getRedisKey(worldId, key)).catch(() => {});
        }
        return deleted;
    }
}

module.exports = { BaseRepository };
