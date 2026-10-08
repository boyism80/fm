import type { Pool, PoolClient } from "pg";
import type { InternalContext } from "../context/internal-context";
import type { RepositoryTxOptions } from "../types/repository-contracts";

export type ShopRow = {
    shop_id: number;
    kind: number;
    account_id: number;
    character_id: number;
    owner_name: string;
    channel_id: number | null;
    map_id: number | null;
    sn: number | null;
    item_id: number;
    title: string;
    meso: number;
    items: unknown[];
    sold: unknown[];
    opened_at: Date;
    closed_at: Date | null;
};

export type ShopInsert = {
    kind: number;
    accountId: number;
    characterId: number;
    ownerName: string;
    channelId: number | null;
    mapId: number | null;
    itemId: number;
    title: string;
    meso: number;
    items: unknown[];
    sold: unknown[];
    closedAt: Date | null;
};

export type ShopUpdate = {
    shopId: number;
    accountId: number;
    sn: number | null;
    title: string;
    meso: number;
    items: unknown[];
    sold: unknown[];
};

const SELECT_COLS =
    "shop_id, kind, account_id, character_id, owner_name, channel_id, map_id, sn, item_id, title, meso, items, sold, opened_at, closed_at";

export class ShopRepository {
    private readonly ctx: InternalContext;

    constructor(internalContext: InternalContext) {
        this.ctx = internalContext;
    }

    pool(worldId: number, accountId: number): Pool {
        return this.ctx.getPgDataPool(worldId, accountId);
    }

    async withAccountLock<T>(worldId: number, accountId: number, fn: (client: PoolClient) => Promise<T>): Promise<T> {
        return this.ctx.withPgDataTransaction(worldId, accountId, async (client) => {
            await client.query("SELECT pg_advisory_xact_lock($1, $2)", [0x6d657263, accountId]);
            return fn(client);
        });
    }

    async findByAccount(worldId: number, accountId: number, options: RepositoryTxOptions = {}): Promise<ShopRow[]> {
        const client = options.txClient ?? this.pool(worldId, accountId);
        const res = await client.query(
            `SELECT ${SELECT_COLS} FROM shops WHERE account_id = $1 ORDER BY closed_at NULLS FIRST, shop_id`,
            [accountId]
        );
        return res.rows as ShopRow[];
    }

    async findOpen(worldId: number): Promise<ShopRow[]> {
        const rows = new Map<string, ShopRow>();
        for (const pool of this.ctx.getPgDataPools(worldId)) {
            const res = await pool.query(`SELECT ${SELECT_COLS} FROM shops WHERE closed_at IS NULL AND sn IS NOT NULL`);
            for (const row of res.rows as ShopRow[]) {
                rows.set(`${row.account_id}:${row.shop_id}`, row);
            }
        }
        return [...rows.values()];
    }

    async insert(worldId: number, shop: ShopInsert, options: RepositoryTxOptions = {}): Promise<ShopRow> {
        const client = options.txClient ?? this.pool(worldId, shop.accountId);
        const res = await client.query(
            `INSERT INTO shops
                (kind, account_id, character_id, owner_name, channel_id, map_id, item_id, title, meso, items, sold, closed_at)
             VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
             RETURNING ${SELECT_COLS}`,
            [
                shop.kind,
                shop.accountId,
                shop.characterId,
                shop.ownerName,
                shop.channelId,
                shop.mapId,
                shop.itemId,
                shop.title,
                shop.meso,
                JSON.stringify(shop.items),
                JSON.stringify(shop.sold),
                shop.closedAt,
            ]
        );
        return res.rows[0] as ShopRow;
    }

    async update(worldId: number, shop: ShopUpdate, options: RepositoryTxOptions = {}): Promise<void> {
        const client = options.txClient ?? this.pool(worldId, shop.accountId);
        await client.query(
            `UPDATE shops SET sn = $3, title = $4, meso = $5, items = $6, sold = $7
             WHERE shop_id = $1 AND account_id = $2 AND closed_at IS NULL`,
            [
                shop.shopId,
                shop.accountId,
                shop.sn,
                shop.title,
                shop.meso,
                JSON.stringify(shop.items),
                JSON.stringify(shop.sold),
            ]
        );
    }

    async close(worldId: number, shopId: number, accountId: number, options: RepositoryTxOptions = {}): Promise<void> {
        const client = options.txClient ?? this.pool(worldId, accountId);
        await client.query(
            `UPDATE shops SET channel_id = NULL, map_id = NULL, sn = NULL, closed_at = NOW()
             WHERE shop_id = $1 AND account_id = $2 AND closed_at IS NULL`,
            [shopId, accountId]
        );
    }

    async deleteClosed(worldId: number, accountId: number, characterId: number, options: RepositoryTxOptions = {}): Promise<ShopRow[]> {
        const client = options.txClient ?? this.pool(worldId, accountId);
        const res = await client.query(
            `DELETE FROM shops WHERE account_id = $1 AND character_id = $2 AND closed_at IS NOT NULL
             RETURNING ${SELECT_COLS}`,
            [accountId, characterId]
        );
        return res.rows as ShopRow[];
    }

    async deleteEmpty(worldId: number, shopId: number, accountId: number, options: RepositoryTxOptions = {}): Promise<void> {
        const client = options.txClient ?? this.pool(worldId, accountId);
        await client.query(
            "DELETE FROM shops WHERE shop_id = $1 AND account_id = $2 AND meso = 0 AND jsonb_array_length(items) = 0",
            [shopId, accountId]
        );
    }

    async closeChannel(worldId: number, channelId: number): Promise<number> {
        let count = 0;
        for (const pool of this.ctx.getPgDataPools(worldId)) {
            await pool.query(
                "DELETE FROM shops WHERE channel_id = $1 AND closed_at IS NULL AND meso = 0 AND jsonb_array_length(items) = 0",
                [channelId]
            );
            const res = await pool.query(
                `UPDATE shops SET channel_id = NULL, map_id = NULL, sn = NULL, closed_at = NOW()
                 WHERE channel_id = $1 AND closed_at IS NULL`,
                [channelId]
            );
            count += res.rowCount ?? 0;
        }
        return count;
    }
}
