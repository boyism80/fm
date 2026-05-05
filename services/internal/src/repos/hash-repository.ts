import type { Pool } from "pg";
import { Repository } from "./repository";
import type { RepositoryQuery, RepositoryTxOptions } from "../types/repository-contracts";

export class HashRepository<TModel = Record<string, unknown>, TRow = Record<string, unknown>> extends Repository<TModel, TRow, string> {
    getGroupKey(_model: TModel): string { throw new Error(`${this.constructor.name}.getGroupKey not implemented`); }
    getItemKey(_model: TModel): string | number { throw new Error(`${this.constructor.name}.getItemKey not implemented`); }
    getRedisHashKey(_worldId: number, _groupKey: string): string { throw new Error(`${this.constructor.name}.getRedisHashKey not implemented`); }
    getRedisKey(worldId: number, key: string): string { return this.getRedisHashKey(worldId, key); }
    onSelect(_groupKey: string, _worldId: number): RepositoryQuery { throw new Error(`${this.constructor.name}.onSelect not implemented`); }
    onBulkUpsert(_rows: TRow[]): RepositoryQuery { throw new Error(`${this.constructor.name}.onBulkUpsert not implemented`); }
    onBulkDelete(_itemKeys: string[], _groupKey: string, _worldId: number): RepositoryQuery { throw new Error(`${this.constructor.name}.onBulkDelete not implemented`); }
    onDelete(row: unknown): RepositoryQuery {
        const model = row as TModel & { worldId?: number };
        const worldId = Number((row as { worldId: number }).worldId ?? model.worldId ?? 0);
        return this.onBulkDelete([String(this.getItemKey(model))], this.getGroupKey(model), worldId);
    }
    rowToModel(_row: TRow): TModel { throw new Error(`${this.constructor.name}.rowToModel not implemented`); }
    modelToRow(_model: TModel): TRow { throw new Error(`${this.constructor.name}.modelToRow not implemented`); }
    getTtlSeconds(): number { return 300; }
    getShardHash(groupKey: string): number { return Number(groupKey); }

    async get(): Promise<never> {
        throw new Error(`${this.constructor.name}.get is not supported for hash repositories`);
    }

    async getMany(): Promise<never> {
        throw new Error(`${this.constructor.name}.getMany is not supported for hash repositories`);
    }

    async getAll(worldId: number, groupKey: string, options: RepositoryTxOptions = {}): Promise<Map<string, TModel>> {
        const hashKey = this.getRedisHashKey(worldId, groupKey);
        const redis = this.redis(worldId, groupKey);
        const useCache = !options.txClient;

        if (useCache && await redis.exists(hashKey)) {
            const fields = await redis.hgetall(hashKey);
            const result = new Map<string, TModel>();
            if (fields) {
                for (const [field, json] of Object.entries(fields)) {
                    if (field === "_loaded") {
                        continue;
                    }
                    const row = JSON.parse(json) as TRow & { deleted?: boolean };
                    if (!row.deleted) {
                        result.set(field, this.rowToModel(row));
                    }
                }
            }
            return result;
        }

        const pool = this.pool(worldId, groupKey) as Pool;
        const select = this.onSelect(groupKey, worldId);
        const res = await this.query(pool, select.text, select.values, options);
        const rows = res.rows as Array<TRow & { deleted?: boolean }>;

        if (useCache) {
            const pipeline = redis.pipeline();
            pipeline.hset(hashKey, "_loaded", "1");
            for (const row of rows) {
                if (!row.deleted) {
                    pipeline.hset(hashKey, String(this.getItemKey(this.rowToModel(row))), JSON.stringify(row));
                }
            }
            pipeline.expire(hashKey, this.getTtlSeconds());
            await pipeline.exec();
        }

        const result = new Map<string, TModel>();
        for (const row of rows) {
            if (!row.deleted) {
                const model = this.rowToModel(row);
                result.set(String(this.getItemKey(model)), model);
            }
        }
        return result;
    }

    async setAll(worldId: number, models: TModel[], options: RepositoryTxOptions = {}): Promise<TModel[]> {
        const dbGroups = new Map<Pool, TModel[]>();
        for (const model of models) {
            const pool = this.pool(worldId, this.getGroupKey(model));
            if (!dbGroups.has(pool)) {
                dbGroups.set(pool, []);
            }
            dbGroups.get(pool)?.push(model);
        }
        if (options.txClient && dbGroups.size > 1) {
            throw new Error(`${this.constructor.name}.setAll with txClient supports only single data shard`);
        }

        const saved: TModel[] = [];
        for (const [pool, groupModels] of dbGroups) {
            const rows = groupModels.map((m) => this.modelToRow(m));
            const bulkUpsert = this.onBulkUpsert(rows);
            const res = await this.query(pool, bulkUpsert.text, bulkUpsert.values, options);

            const byGroup = new Map<string, Array<{ row: TRow; model: TModel }>>();
            for (const savedRow of res.rows as TRow[]) {
                const model = this.rowToModel(savedRow);
                const key = this.getGroupKey(model);
                if (!byGroup.has(key)) {
                    byGroup.set(key, []);
                }
                byGroup.get(key)?.push({ row: savedRow, model });
            }

            for (const [groupKey, groupItems] of byGroup) {
                const hashKey = this.getRedisHashKey(worldId, groupKey);
                const redis = this.redis(worldId, groupKey);
                if (!options.txClient && await redis.exists(hashKey)) {
                    const pipeline = redis.pipeline();
                    for (const { row, model } of groupItems) {
                        pipeline.hset(hashKey, String(this.getItemKey(model)), JSON.stringify(row));
                    }
                    await pipeline.exec();
                }
                saved.push(...groupItems.map((i) => i.model));
            }
        }

        return saved;
    }

    async set(worldId: number, model: TModel, options: RepositoryTxOptions = {}): Promise<TModel> {
        const result = await this.setAll(worldId, [model], options);
        const first = result[0];
        if (!first) {
            throw new Error(`${this.constructor.name}.set returned no rows`);
        }
        return first;
    }

    async delAll(worldId: number, groupKey: string, itemKeys: Array<string | number>, options: RepositoryTxOptions = {}): Promise<void> {
        if (!itemKeys.length) {
            return;
        }
        const pool = this.pool(worldId, groupKey) as Pool;
        const deleteQuery = this.onBulkDelete(itemKeys.map(String), groupKey, worldId);
        await this.query(pool, deleteQuery.text, deleteQuery.values, options);

        const hashKey = this.getRedisHashKey(worldId, groupKey);
        const redis = this.redis(worldId, groupKey);
        if (!options.txClient && await redis.exists(hashKey)) {
            await redis.hdel(hashKey, ...itemKeys.map(String));
        }
    }

    async del(worldId: number, groupKey: string, itemKey: string | number, options: RepositoryTxOptions = {}): Promise<void> {
        return this.delAll(worldId, groupKey, [itemKey], options);
    }

    async delete(row: { worldId: number } & Record<string, unknown>, options: RepositoryTxOptions = {}): Promise<boolean> {
        const worldId = Number(row.worldId);
        await this.del(worldId, this.getGroupKey(row as unknown as TModel), String(this.getItemKey(row as unknown as TModel)), options);
        return true;
    }

    async evictGroupCache(worldId: number, groupKey: string): Promise<void> {
        const hashKey = this.getRedisHashKey(worldId, groupKey);
        const redis = this.redis(worldId, groupKey);
        await redis.del(hashKey).catch(() => {});
    }
}
