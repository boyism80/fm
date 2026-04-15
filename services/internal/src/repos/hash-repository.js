"use strict";

class HashRepository {
    constructor(internalContext) {
        this.ctx = internalContext;
    }

    getOwnerKey(_model) { throw new Error(`${this.constructor.name}.getOwnerKey not implemented`); }
    getItemKey(_model) { throw new Error(`${this.constructor.name}.getItemKey not implemented`); }
    getRedisHashKey(_worldId, _ownerKey) { throw new Error(`${this.constructor.name}.getRedisHashKey not implemented`); }
    onSelectByOwner(_ownerKey, _worldId) { throw new Error(`${this.constructor.name}.onSelectByOwner not implemented`); }
    onBulkUpsert(_rows) { throw new Error(`${this.constructor.name}.onBulkUpsert not implemented`); }
    onBulkDelete(_itemKeys, _ownerKey, _worldId) { throw new Error(`${this.constructor.name}.onBulkDelete not implemented`); }
    onDelete(row, worldId) {
        return this.onBulkDelete([String(this.getItemKey(row))], this.getOwnerKey(row), worldId);
    }
    rowToModel(_row) { throw new Error(`${this.constructor.name}.rowToModel not implemented`); }
    modelToRow(_model) { throw new Error(`${this.constructor.name}.modelToRow not implemented`); }

    getTtlSeconds() { return 300; }
    getShardHash(ownerKey) { return ownerKey; }

    _pool(worldId, ownerKey) {
        return this.ctx.getPgDataPool(worldId, this.getShardHash(ownerKey));
    }

    _redis(worldId, ownerKey) {
        return this.ctx.getRedisDataAccess(worldId, this.getShardHash(ownerKey)).client;
    }

    async getAll(worldId, ownerKey) {
        const hashKey = this.getRedisHashKey(worldId, ownerKey);
        const redis = this._redis(worldId, ownerKey);

        if (await redis.exists(hashKey)) {
            const fields = await redis.hgetall(hashKey);
            const result = new Map();
            if (fields) {
                for (const [field, json] of Object.entries(fields)) {
                    if (field === "_loaded") continue;
                    const row = JSON.parse(json);
                    if (!row.deleted) result.set(field, this.rowToModel(row));
                }
            }
            return result;
        }

        const pool = this._pool(worldId, ownerKey);
        const { text, values } = this.onSelectByOwner(ownerKey, worldId);
        const res = await pool.query(text, values);

        const pipeline = redis.pipeline();
        pipeline.hset(hashKey, "_loaded", "1");
        for (const row of res.rows) {
            if (!row.deleted) {
                pipeline.hset(hashKey, String(this.getItemKey(this.rowToModel(row))), JSON.stringify(row));
            }
        }
        pipeline.expire(hashKey, this.getTtlSeconds());
        await pipeline.exec();

        const result = new Map();
        for (const row of res.rows) {
            if (!row.deleted) {
                const model = this.rowToModel(row);
                result.set(String(this.getItemKey(model)), model);
            }
        }
        return result;
    }

    async setAll(worldId, models) {
        const dbGroups = new Map();
        for (const model of models) {
            const pool = this._pool(worldId, this.getOwnerKey(model));
            if (!dbGroups.has(pool)) dbGroups.set(pool, []);
            dbGroups.get(pool).push(model);
        }

        const saved = [];
        for (const [pool, groupModels] of dbGroups) {
            const rows = groupModels.map(m => this.modelToRow(m));
            const { text, values } = this.onBulkUpsert(rows);
            const res = await pool.query(text, values);

            const byOwner = new Map();
            for (const savedRow of res.rows) {
                const model = this.rowToModel(savedRow);
                const ownerKey = this.getOwnerKey(model);
                if (!byOwner.has(ownerKey)) byOwner.set(ownerKey, []);
                byOwner.get(ownerKey).push({ row: savedRow, model });
            }

            for (const [ownerKey, ownerItems] of byOwner) {
                const hashKey = this.getRedisHashKey(worldId, ownerKey);
                const redis = this._redis(worldId, ownerKey);
                if (await redis.exists(hashKey)) {
                    const pipeline = redis.pipeline();
                    for (const { row, model } of ownerItems) {
                        pipeline.hset(hashKey, String(this.getItemKey(model)), JSON.stringify(row));
                    }
                    await pipeline.exec();
                }
                saved.push(...ownerItems.map(i => i.model));
            }
        }

        return saved;
    }

    async set(worldId, model) {
        const result = await this.setAll(worldId, [model]);
        return result[0];
    }

    async delAll(worldId, ownerKey, itemKeys) {
        if (!itemKeys.length) return;
        const pool = this._pool(worldId, ownerKey);
        const { text, values } = this.onBulkDelete(itemKeys, ownerKey, worldId);
        await pool.query(text, values);

        const hashKey = this.getRedisHashKey(worldId, ownerKey);
        const redis = this._redis(worldId, ownerKey);
        if (await redis.exists(hashKey)) {
            await redis.hdel(hashKey, ...itemKeys.map(String));
        }
    }

    async del(worldId, ownerKey, itemKey) {
        return this.delAll(worldId, ownerKey, [itemKey]);
    }

    async delete(row) {
        const worldId = Number(row.worldId);
        return this.del(worldId, this.getOwnerKey(row), String(this.getItemKey(row)));
    }
}

module.exports = { HashRepository };
