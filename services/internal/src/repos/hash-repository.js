"use strict";

const { Repository } = require("./repository");

class HashRepository extends Repository {

    getGroupKey(_model) { throw new Error(`${this.constructor.name}.getGroupKey not implemented`); }
    getItemKey(_model) { throw new Error(`${this.constructor.name}.getItemKey not implemented`); }
    getRedisHashKey(_worldId, _groupKey) { throw new Error(`${this.constructor.name}.getRedisHashKey not implemented`); }
    getRedisKey(worldId, key) { return this.getRedisHashKey(worldId, key); }
    onSelect(_groupKey, _worldId) { throw new Error(`${this.constructor.name}.onSelect not implemented`); }
    onBulkUpsert(_rows) { throw new Error(`${this.constructor.name}.onBulkUpsert not implemented`); }
    onBulkDelete(_itemKeys, _groupKey, _worldId) { throw new Error(`${this.constructor.name}.onBulkDelete not implemented`); }
    onDelete(row, worldId) {
        return this.onBulkDelete([String(this.getItemKey(row))], this.getGroupKey(row), worldId);
    }
    rowToModel(_row) { throw new Error(`${this.constructor.name}.rowToModel not implemented`); }
    modelToRow(_model) { throw new Error(`${this.constructor.name}.modelToRow not implemented`); }

    getTtlSeconds() { return 300; }
    getShardHash(groupKey) { return groupKey; }

    async get() {
        throw new Error(`${this.constructor.name}.get is not supported for hash repositories`);
    }

    async getMany() {
        throw new Error(`${this.constructor.name}.getMany is not supported for hash repositories`);
    }

    async getAll(worldId, groupKey) {
        const hashKey = this.getRedisHashKey(worldId, groupKey);
        const redis = this._redis(worldId, groupKey);

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

        const pool = this._pool(worldId, groupKey);
        const { text, values } = this.onSelect(groupKey, worldId);
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
            const pool = this._pool(worldId, this.getGroupKey(model));
            if (!dbGroups.has(pool)) dbGroups.set(pool, []);
            dbGroups.get(pool).push(model);
        }

        const saved = [];
        for (const [pool, groupModels] of dbGroups) {
            const rows = groupModels.map(m => this.modelToRow(m));
            const { text, values } = this.onBulkUpsert(rows);
            const res = await pool.query(text, values);

            const byGroup = new Map();
            for (const savedRow of res.rows) {
                const model = this.rowToModel(savedRow);
                const groupKey = this.getGroupKey(model);
                if (!byGroup.has(groupKey)) byGroup.set(groupKey, []);
                byGroup.get(groupKey).push({ row: savedRow, model });
            }

            for (const [groupKey, groupItems] of byGroup) {
                const hashKey = this.getRedisHashKey(worldId, groupKey);
                const redis = this._redis(worldId, groupKey);
                if (await redis.exists(hashKey)) {
                    const pipeline = redis.pipeline();
                    for (const { row, model } of groupItems) {
                        pipeline.hset(hashKey, String(this.getItemKey(model)), JSON.stringify(row));
                    }
                    await pipeline.exec();
                }
                saved.push(...groupItems.map(i => i.model));
            }
        }

        return saved;
    }

    async set(worldId, model) {
        const result = await this.setAll(worldId, [model]);
        return result[0];
    }

    async delAll(worldId, groupKey, itemKeys) {
        if (!itemKeys.length) return;
        const pool = this._pool(worldId, groupKey);
        const { text, values } = this.onBulkDelete(itemKeys, groupKey, worldId);
        await pool.query(text, values);

        const hashKey = this.getRedisHashKey(worldId, groupKey);
        const redis = this._redis(worldId, groupKey);
        if (await redis.exists(hashKey)) {
            await redis.hdel(hashKey, ...itemKeys.map(String));
        }
    }

    async del(worldId, groupKey, itemKey) {
        return this.delAll(worldId, groupKey, [itemKey]);
    }

    async delete(row) {
        const worldId = Number(row.worldId);
        return this.del(worldId, this.getGroupKey(row), String(this.getItemKey(row)));
    }
}

module.exports = { HashRepository };
