import type { Pool, PoolClient } from "pg";
import type { InternalContext } from "../context/internal-context";
import type { RepositoryTxOptions } from "../types/repository-contracts";

export type HiredMerchantRow = {
    merchant_id: number;
    account_id: number;
    character_id: number;
    owner_name: string;
    channel_id: number | null;
    map_id: number | null;
    item_id: number;
    title: string;
    meso: number;
    items: unknown[];
    sold: unknown[];
    opened_at: Date;
    closed_at: Date | null;
};

export type HiredMerchantInsert = {
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

export type HiredMerchantUpdate = {
    merchantId: number;
    accountId: number;
    title: string;
    meso: number;
    items: unknown[];
    sold: unknown[];
};

const SELECT_COLS =
    "merchant_id, account_id, character_id, owner_name, channel_id, map_id, item_id, title, meso, items, sold, opened_at, closed_at";

export class HiredMerchantRepository {
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

    async findByAccount(worldId: number, accountId: number, options: RepositoryTxOptions = {}): Promise<HiredMerchantRow | null> {
        const client = options.txClient ?? this.pool(worldId, accountId);
        const res = await client.query(`SELECT ${SELECT_COLS} FROM hired_merchants WHERE account_id = $1`, [accountId]);
        return (res.rows[0] as HiredMerchantRow | undefined) ?? null;
    }

    async insert(worldId: number, merchant: HiredMerchantInsert, options: RepositoryTxOptions = {}): Promise<HiredMerchantRow> {
        const client = options.txClient ?? this.pool(worldId, merchant.accountId);
        const res = await client.query(
            `INSERT INTO hired_merchants
                (account_id, character_id, owner_name, channel_id, map_id, item_id, title, meso, items, sold, closed_at)
             VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
             RETURNING ${SELECT_COLS}`,
            [
                merchant.accountId,
                merchant.characterId,
                merchant.ownerName,
                merchant.channelId,
                merchant.mapId,
                merchant.itemId,
                merchant.title,
                merchant.meso,
                JSON.stringify(merchant.items),
                JSON.stringify(merchant.sold),
                merchant.closedAt,
            ]
        );
        return res.rows[0] as HiredMerchantRow;
    }

    async update(worldId: number, merchant: HiredMerchantUpdate, options: RepositoryTxOptions = {}): Promise<void> {
        const client = options.txClient ?? this.pool(worldId, merchant.accountId);
        await client.query(
            `UPDATE hired_merchants SET title = $3, meso = $4, items = $5, sold = $6
             WHERE merchant_id = $1 AND account_id = $2 AND closed_at IS NULL`,
            [
                merchant.merchantId,
                merchant.accountId,
                merchant.title,
                merchant.meso,
                JSON.stringify(merchant.items),
                JSON.stringify(merchant.sold),
            ]
        );
    }

    async close(worldId: number, merchantId: number, accountId: number, options: RepositoryTxOptions = {}): Promise<void> {
        const client = options.txClient ?? this.pool(worldId, accountId);
        await client.query(
            `UPDATE hired_merchants SET channel_id = NULL, map_id = NULL, closed_at = NOW()
             WHERE merchant_id = $1 AND account_id = $2 AND closed_at IS NULL`,
            [merchantId, accountId]
        );
    }

    async delete(worldId: number, merchantId: number, accountId: number, options: RepositoryTxOptions = {}): Promise<void> {
        const client = options.txClient ?? this.pool(worldId, accountId);
        await client.query("DELETE FROM hired_merchants WHERE merchant_id = $1 AND account_id = $2", [merchantId, accountId]);
    }

    async deleteEmpty(worldId: number, merchantId: number, accountId: number, options: RepositoryTxOptions = {}): Promise<void> {
        const client = options.txClient ?? this.pool(worldId, accountId);
        await client.query(
            "DELETE FROM hired_merchants WHERE merchant_id = $1 AND account_id = $2 AND meso = 0 AND jsonb_array_length(items) = 0",
            [merchantId, accountId]
        );
    }

    async closeChannel(worldId: number, channelId: number): Promise<number> {
        let count = 0;
        for (const pool of this.ctx.getPgDataPools(worldId)) {
            await pool.query(
                "DELETE FROM hired_merchants WHERE channel_id = $1 AND closed_at IS NULL AND meso = 0 AND jsonb_array_length(items) = 0",
                [channelId]
            );
            const res = await pool.query(
                `UPDATE hired_merchants SET channel_id = NULL, map_id = NULL, closed_at = NOW()
                 WHERE channel_id = $1 AND closed_at IS NULL`,
                [channelId]
            );
            count += res.rowCount ?? 0;
        }
        return count;
    }
}
