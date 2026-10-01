import type { Pool } from "pg";
import { Repository } from "./repository";
import type { RepositoryQuery, RepositoryTxOptions } from "../types/repository-contracts";
import type { InternalContext } from "../context/internal-context";

// KEYS[1] hash. ARGV[1] ttl, then field/value pairs. A reader that loaded pre-commit rows must not replace a hash a writer already refreshed.
const HASH_FILL_IF_ABSENT_SCRIPT = `
if redis.call("EXISTS", KEYS[1]) == 1 then
    return 0
end
redis.call("HSET", KEYS[1], "_loaded", "1")
for i = 2, #ARGV, 2 do
    redis.call("HSET", KEYS[1], ARGV[i], ARGV[i + 1])
end
redis.call("EXPIRE", KEYS[1], ARGV[1])
return 1
`;

// KEYS[1] hash. ARGV field/value pairs. Writing into an expired key would leave a partial hash without a TTL.
const HASH_SET_IF_PRESENT_SCRIPT = `
if redis.call("EXISTS", KEYS[1]) == 0 then
    return 0
end
for i = 1, #ARGV, 2 do
    redis.call("HSET", KEYS[1], ARGV[i], ARGV[i + 1])
end
return 1
`;

export abstract class HashRepository<TModel = Record<string, unknown>, TRow = Record<string, unknown>> extends Repository<TModel, TRow, string> {
    private readonly localGroupCache: Map<string, Map<string, TRow>>;

    constructor(internalContext: InternalContext) {
        super(internalContext);
        this.localGroupCache = new Map();
    }

    override async refreshCache(worldId: number, key: string): Promise<void> {
        await this.refreshGroupCache(worldId, key);
    }

    override getKey(model: TModel): string {
        return this.getGroupKey(model);
    }

    override onUpsert(row: TRow): RepositoryQuery {
        return this.onBulkUpsert([row]);
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
            const result = new Map<string, TModel>();
            for (const [itemKey, row] of localCached) {
                result.set(itemKey, this.rowToModel(this.normalizeRow(row)));
            }
            return result;
        }

        const fields = useCache ? await redis.hgetall(hashKey) : null;
        if (fields?._loaded) {
            const result = new Map<string, TModel>();
            const localRows = new Map<string, TRow>();
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
            this.localGroupCache.set(hashKey, localRows);
            return result;
        }

        const pool = this.pool(worldId, groupKey) as Pool;
        const select = this.onSelect(groupKey, worldId);
        const res = await this.query(pool, select.text, select.values, options);
        const rows = res.rows as Array<TRow & { deleted?: boolean }>;

        if (useCache) {
            await redis.eval(HASH_FILL_IF_ABSENT_SCRIPT, 1, hashKey, this.getTtlSeconds(), ...this.rowsToHashFields(rows));
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

    async getItem(
        worldId: number,
        groupKey: string,
        itemKey: string | number,
        options: RepositoryTxOptions = {}
    ): Promise<TModel | null> {
        const all = await this.getAll(worldId, groupKey, options);
        return all.get(String(itemKey)) ?? null;
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
                if (!options.txClient) {
                    const fields: string[] = [];
                    for (const { row, model } of groupItems) {
                        fields.push(this.getItemKey(model), JSON.stringify(row));
                    }
                    await redis.eval(HASH_SET_IF_PRESENT_SCRIPT, 1, hashKey, ...fields);
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

    private rowsToHashFields(rows: Array<TRow & { deleted?: boolean }>): string[] {
        const fields: string[] = [];
        for (const raw of rows) {
            if (!raw.deleted) {
                const row = this.normalizeRow(raw);
                fields.push(this.getItemKey(this.rowToModel(row)), JSON.stringify(row));
            }
        }
        return fields;
    }

    // Writers call this after commit while holding the entity lock. The hash is rebuilt in one MULTI, so the key
    // stays present and a reader that loaded pre-commit rows cannot fill it back.
    async refreshGroupCache(worldId: number, groupKey: string): Promise<void> {
        const hashKey = this.getRedisHashKey(worldId, groupKey);
        this.localGroupCache.delete(hashKey);
        const redis = this.redis(worldId, groupKey);
        try {
            const select = this.onSelect(groupKey, worldId);
            const res = await this.pool(worldId, groupKey).query(select.text, select.values);
            const fields = this.rowsToHashFields(res.rows as Array<TRow & { deleted?: boolean }>);
            const multi = redis.multi().del(hashKey).hset(hashKey, "_loaded", "1");
            if (fields.length > 0) {
                multi.hset(hashKey, ...fields);
            }
            await multi.expire(hashKey, this.getTtlSeconds()).exec();
        } catch {
            await redis.del(hashKey).catch(() => {});
        }
    }
}
