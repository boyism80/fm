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
    abstract evictCache(worldId: number, key: TKey): Promise<void>;
    abstract get(worldId: number, key: TKey, options?: RepositoryTxOptions): Promise<TModel | null>;
    abstract set(worldId: number, model: TModel, options?: RepositoryTxOptions): Promise<TModel>;
    abstract getMany(worldId: number, keys: TKey[], options?: RepositoryTxOptions): Promise<Map<TKey, TModel>>;
    abstract getAll(worldId: number, key: TKey, options?: RepositoryTxOptions): Promise<unknown>;
    abstract setAll(worldId: number, models: TModel[], options?: RepositoryTxOptions): Promise<TModel[]>;
    abstract delAll(worldId: number, key: TKey, itemKeys: Array<string | number>, options?: RepositoryTxOptions): Promise<void>;
    abstract delete(row: { worldId: number } & Record<string, unknown>, options?: RepositoryTxOptions): Promise<boolean>;
}
