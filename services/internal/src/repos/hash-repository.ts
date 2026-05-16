import type { Pool } from "pg";
import { Repository } from "./repository";
import type { RepositoryQuery, RepositoryTxOptions } from "../types/repository-contracts";
import type { InternalContext } from "../context/internal-context";

export abstract class HashRepository<TModel = Record<string, unknown>, TRow = Record<string, unknown>> extends Repository<TModel, TRow, string> {
    private readonly localGroupCache: Map<string, Map<string, TRow>>;

    constructor(internalContext: InternalContext) {
        super(internalContext);
        this.localGroupCache = new Map();
    }

    protected logL1Group(method: string, result: "hit" | "miss", hashKey: string): void {
        console.log(`[L1][${this.constructor.name}][${method}] ${result} hash=${hashKey}`);
    }

    override async evictCache(worldId: number, key: string): Promise<void> {
        await this.evictGroupCache(worldId, key);
    }

    abstract getGroupKey(_model: TModel): string;
    abstract getItemKey(_model: TModel): string;
    abstract getRedisHashKey(_worldId: number, _groupKey: string): string;
    override getRedisKey(worldId: number, key: string): string {
        return this.getRedisHashKey(worldId, key);
    }
    abstract onSelect(_groupKey: string, _worldId: number): RepositoryQuery;
    abstract onBulkUpsert(_rows: TRow[]): RepositoryQuery;
    abstract onBulkDelete(_itemKeys: string[], _groupKey: string, _worldId: number): RepositoryQuery;
    override onDelete(row: unknown): RepositoryQuery {
        const model = row as TModel & { worldId?: number };
        const worldId = Number((row as { worldId: number }).worldId ?? model.worldId ?? 0);
        return this.onBulkDelete([this.getItemKey(model)], this.getGroupKey(model), worldId);
    }
    abstract rowToModel(_row: TRow): TModel;
    abstract modelToRow(_model: TModel): TRow;
    override getTtlSeconds(): number {
        return 300;
    }

    override getShardHash(groupKey: string): number {
        return Number(groupKey);
    }

    override async get(_worldId: number, _key: string, _options: RepositoryTxOptions = {}): Promise<never> {
        throw new Error(`${this.constructor.name}.get is not supported for hash repositories`);
    }

    override async getMany(_worldId: number, _keys: string[], _options: RepositoryTxOptions = {}): Promise<never> {
        throw new Error(`${this.constructor.name}.getMany is not supported for hash repositories`);
    }

    override async getAll(worldId: number, groupKey: string, options: RepositoryTxOptions = {}): Promise<Map<string, TModel>> {
        const hashKey = this.getRedisHashKey(worldId, groupKey);
        const redis = this.redis(worldId, groupKey);
        const useCache = !options.txClient;
        const localCached = this.localGroupCache.get(hashKey);
        if (localCached) {
            this.logL1Group("getAll", "hit", hashKey);
            const result = new Map<string, TModel>();
            for (const [itemKey, row] of localCached) {
                result.set(itemKey, this.rowToModel(this.normalizeRow(row)));
            }
            return result;
        }
        this.logL1Group("getAll", "miss", hashKey);

        if (useCache && await redis.exists(hashKey)) {
            const fields = await redis.hgetall(hashKey);
            const result = new Map<string, TModel>();
            const localRows = new Map<string, TRow>();
            if (fields) {
                for (const [field, json] of Object.entries(fields)) {
                    if (field === "_loaded") {
                        continue;
                    }
                    const parsed = JSON.parse(json) as TRow & { deleted?: boolean };
                    if (!parsed.deleted) {
                        const typedRow = this.normalizeRow(parsed);
                        localRows.set(field, typedRow);
                        result.set(field, this.rowToModel(typedRow));
                    }
                }
            }
            this.localGroupCache.set(hashKey, localRows);
            return result;
        }

        const pool = this.pool(worldId, groupKey) as Pool;
        const select = this.onSelect(groupKey, worldId);
        const res = await this.query(pool, select.text, select.values, options);
        const rows = res.rows as Array<TRow & { deleted?: boolean }>;

        if (useCache) {
            const pipeline = redis.pipeline();
            pipeline.hset(hashKey, "_loaded", "1");
            for (const raw of rows) {
                if (!raw.deleted) {
                    const row = this.normalizeRow(raw);
                    pipeline.hset(hashKey, this.getItemKey(this.rowToModel(row)), JSON.stringify(row));
                }
            }
            pipeline.expire(hashKey, this.getTtlSeconds());
            await pipeline.exec();
        }

        const result = new Map<string, TModel>();
        const localRows = new Map<string, TRow>();
        for (const raw of rows) {
            if (!raw.deleted) {
                const typedRow = this.normalizeRow(raw);
                const model = this.rowToModel(typedRow);
                result.set(this.getItemKey(model), model);
                localRows.set(this.getItemKey(model), typedRow);
            }
        }
        this.localGroupCache.set(hashKey, localRows);
        return result;
    }

    override async setAll(worldId: number, models: TModel[], options: RepositoryTxOptions = {}): Promise<TModel[]> {
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
            for (const raw of res.rows as TRow[]) {
                const savedRow = this.normalizeRow(raw);
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
                        pipeline.hset(hashKey, this.getItemKey(model), JSON.stringify(row));
                    }
                    await pipeline.exec();
                }
                const localCached = this.localGroupCache.get(hashKey);
                if (localCached) {
                    for (const { row, model } of groupItems) {
                        localCached.set(this.getItemKey(model), row);
                    }
                }
                saved.push(...groupItems.map((i) => i.model));
            }
        }

        return saved;
    }

    override async set(worldId: number, model: TModel, options: RepositoryTxOptions = {}): Promise<TModel> {
        const result = await this.setAll(worldId, [model], options);
        const first = result[0];
        if (!first) {
            throw new Error(`${this.constructor.name}.set returned no rows`);
        }
        return first;
    }

    override async delAll(worldId: number, groupKey: string, itemKeys: Array<string | number>, options: RepositoryTxOptions = {}): Promise<void> {
        if (!itemKeys.length) {
            return;
        }
        const pool = this.pool(worldId, groupKey) as Pool;
        const deleteQuery = this.onBulkDelete(itemKeys.map(String), groupKey, worldId);
        await this.query(pool, deleteQuery.text, deleteQuery.values, options);

        const hashKey = this.getRedisHashKey(worldId, groupKey);
        const redis = this.redis(worldId, groupKey);
        const keyStrings = itemKeys.map(String);
        if (!options.txClient && await redis.exists(hashKey)) {
            await redis.hdel(hashKey, ...keyStrings);
        }
        const localCached = this.localGroupCache.get(hashKey);
        if (localCached) {
            for (const itemKey of keyStrings) {
                localCached.delete(itemKey);
            }
        }
    }

    async del(worldId: number, groupKey: string, itemKey: string | number, options: RepositoryTxOptions = {}): Promise<void> {
        return this.delAll(worldId, groupKey, [itemKey], options);
    }

    override async delete(row: { worldId: number } & Record<string, unknown>, options: RepositoryTxOptions = {}): Promise<boolean> {
        const worldId = row.worldId;
        const model = row as TModel;
        await this.del(worldId, this.getGroupKey(model), this.getItemKey(model), options);
        return true;
    }

    async evictGroupCache(worldId: number, groupKey: string): Promise<void> {
        const hashKey = this.getRedisHashKey(worldId, groupKey);
        this.localGroupCache.delete(hashKey);
        const redis = this.redis(worldId, groupKey);
        await redis.del(hashKey).catch(() => {});
    }
}
