import type { Pool, PoolClient } from "pg";
import { randomUUID } from "node:crypto";
import type { InternalContext } from "../context/internal-context";
import type { RepositoryTxOptions } from "../types/repository-contracts";
import { RequestRepository } from "./request-repository";

export type CashItemRow = {
    serial: string;
    account_id: number;
    item: unknown;
};

export type CashGiftRow = {
    serial: string;
    item_id: number;
    sender_name: string;
    message: string;
};

export type CashRingRow = {
    serial: string;
    partner_serial: string;
    character_id: number;
    partner_id: number;
    partner_name: string;
    item_id: number;
};

export type CashCouponRow = {
    code: string;
    kind: number;
    value: number;
    used_at: Date | null;
    used_by: number | null;
    request_id: string | null;
};

export type CashBalance = {
    nxCash: number;
    maplePoint: number;
};

export enum CashTransferKind {
    Gift = 0,
    Ring = 1,
    TakeOut = 2,
    PutIn = 3,
    InventorySlot = 4,
    Meso = 5,
    Purchase = 6,
    StorageSlot = 7,
    CharacterSlot = 8,
    Credit = 9,
}

// The database the request row lives in: a world data shard, or the unified database for requests made with a
// balance change.
export type CashTransferRequest = {
    request_id: string;
    source_id: number;
    target_id: number;
    kind: CashTransferKind;
    detail: unknown;
    pool: Pool;
};

export class CashShopRepository {
    private readonly ctx: InternalContext;

    constructor(internalContext: InternalContext) {
        this.ctx = internalContext;
    }

    // Balances are account-wide and live in the unified database; every balance change holds this lock.
    async withBalanceLock<T>(accountId: number, fn: (client: PoolClient) => Promise<T>): Promise<T> {
        return this.ctx.withPgClient(this.ctx.getPgUnifiedPool(), async (client) => {
            await client.query("BEGIN");
            try {
                await client.query("SELECT pg_advisory_xact_lock($1, $2)", [0x6361736a, accountId]);
                const result = await fn(client);
                await client.query("COMMIT");
                return result;
            } catch (err) {
                await client.query("ROLLBACK").catch(() => {});
                throw err;
            }
        });
    }

    async withAccountLock<T>(
        worldId: number,
        accountId: number,
        fn: (client: PoolClient) => Promise<T>,
        options: RepositoryTxOptions = {}
    ): Promise<T> {
        return this.ctx.withPgDataTransaction(worldId, accountId, async (client) => {
            await client.query("SELECT pg_advisory_xact_lock($1, $2)", [0x63617368, accountId]);
            return fn(client);
        }, options);
    }

    async withCharacterLock<T>(worldId: number, characterId: number, fn: (client: PoolClient) => Promise<T>): Promise<T> {
        return this.ctx.withPgDataTransaction(worldId, characterId, async (client) => {
            await client.query("SELECT pg_advisory_xact_lock($1, $2)", [0x63617369, characterId]);
            return fn(client);
        });
    }

    async withAccount<T>(worldId: number, accountId: number, fn: (client: PoolClient) => Promise<T>, options: RepositoryTxOptions = {}): Promise<T> {
        return this.ctx.withPgDataTransaction(worldId, accountId, fn, options);
    }

    async queueTransfer(sourceId: number, targetId: number, kind: CashTransferKind, detail: unknown, client: PoolClient): Promise<void> {
        await client.query(
            "INSERT INTO cash_transfer_requests (request_id, source_id, target_id, kind, detail) VALUES ($1, $2, $3, $4, $5)",
            [randomUUID(), sourceId, targetId, kind, JSON.stringify(detail)]
        );
    }

    // Queued in the unified database with a balance change, for the shards of one world.
    async queueBalanceTransfer(worldId: number, sourceId: number, targetId: number, kind: CashTransferKind, detail: unknown, client: PoolClient): Promise<void> {
        await client.query(
            "INSERT INTO cash_transfer_requests (request_id, world_id, source_id, target_id, kind, detail) VALUES ($1, $2, $3, $4, $5, $6)",
            [randomUUID(), worldId, sourceId, targetId, kind, JSON.stringify(detail)]
        );
    }

    async pendingTransfers(worldId: number): Promise<CashTransferRequest[]> {
        const requests: CashTransferRequest[] = [];
        const select = "SELECT request_id, source_id, target_id, kind, detail FROM cash_transfer_requests WHERE completed = FALSE";
        for (const pool of this.ctx.getPgDataPools(worldId)) {
            const res = await pool.query(`${select} AND world_id IS NULL ORDER BY created_at LIMIT 100`);
            requests.push(...(res.rows as Omit<CashTransferRequest, "pool">[]).map((row) => ({ ...row, pool })));
        }
        const unified = this.ctx.getPgUnifiedPool();
        const res = await unified.query(`${select} AND world_id = $1 ORDER BY created_at LIMIT 100`, [worldId]);
        requests.push(...(res.rows as Omit<CashTransferRequest, "pool">[]).map((row) => ({ ...row, pool: unified })));
        return requests;
    }

    async receiveTransfer(worldId: number, request: CashTransferRequest, fn: (client: PoolClient) => Promise<void>): Promise<void> {
        await this.ctx.withPgDataTransaction(worldId, request.target_id, async (client) => {
            if (
                request.kind === CashTransferKind.Gift ||
                request.kind === CashTransferKind.PutIn ||
                request.kind === CashTransferKind.Purchase ||
                request.kind === CashTransferKind.StorageSlot ||
                request.kind === CashTransferKind.CharacterSlot
            ) {
                await client.query("SELECT pg_advisory_xact_lock($1, $2)", [0x63617368, request.target_id]);
            } else if (request.kind === CashTransferKind.InventorySlot || request.kind === CashTransferKind.Meso) {
                await client.query("SELECT pg_advisory_xact_lock($1, $2)", [0x63617369, request.target_id]);
            }
            await fn(client);
        });
    }

    // Items on their way into the account's locker of this world, from any queue.
    async pendingLockerItemCount(worldId: number, accountId: number): Promise<number> {
        const pending = `SELECT COALESCE(SUM(jsonb_array_length(detail->'items')), 0) AS count
                         FROM cash_transfer_requests
                         WHERE completed = FALSE AND kind = ANY($2::int[]) AND detail->>'accountId' = $1::text
                           AND jsonb_typeof(detail->'items') = 'array'`;
        const kinds = [CashTransferKind.Gift, CashTransferKind.Purchase];
        let count = 0;
        for (const pool of this.ctx.getPgDataPools(worldId)) {
            const res = await pool.query(`${pending} AND world_id IS NULL`, [accountId, kinds]);
            count += Number(res.rows[0]?.count ?? 0);
        }
        const res = await this.ctx.getPgUnifiedPool().query(`${pending} AND world_id = $3`, [accountId, kinds, worldId]);
        return count + Number(res.rows[0]?.count ?? 0);
    }

    async finishTransfer(request: CashTransferRequest): Promise<void> {
        await request.pool.query("UPDATE cash_transfer_requests SET completed = TRUE WHERE request_id = $1", [request.request_id]);
    }

    async isTransferCompleted(request: CashTransferRequest): Promise<boolean> {
        const res = await request.pool.query("SELECT completed FROM cash_transfer_requests WHERE request_id = $1", [request.request_id]);
        const row = res.rows[0] as { completed: boolean } | undefined;
        return row?.completed === true;
    }

    async claimTransfer(requestId: string, client: PoolClient): Promise<boolean> {
        const res = await client.query(
            "UPDATE cash_transfer_requests SET completed = TRUE WHERE request_id = $1 AND completed = FALSE RETURNING request_id",
            [requestId]
        );
        return (res.rowCount ?? 0) > 0;
    }

    async completeRingTransfers(sourceId: number, serials: string[], client: PoolClient): Promise<void> {
        if (serials.length === 0) {
            return;
        }
        await client.query(
            `UPDATE cash_transfer_requests
             SET completed = TRUE
             WHERE completed = FALSE AND source_id = $1 AND kind = $3 AND detail->>'serial' = ANY($2::text[])`,
            [sourceId, serials, CashTransferKind.Ring]
        );
    }

    async deleteRing(worldId: number, characterId: number, serial: string): Promise<void> {
        await this.ctx.getPgDataPool(worldId, characterId).query(
            "DELETE FROM cash_rings WHERE character_id = $1 AND serial = $2",
            [characterId, serial]
        );
    }

    async findBalance(accountId: number, options: RepositoryTxOptions = {}): Promise<CashBalance> {
        const client = this.ctx.pgClient(this.ctx.getPgUnifiedPool(), options);
        const res = await client.query("SELECT nx_cash, maple_point FROM accounts WHERE id = $1", [accountId]);
        const row = res.rows[0] as { nx_cash: number; maple_point: number } | undefined;
        return { nxCash: row?.nx_cash ?? 0, maplePoint: row?.maple_point ?? 0 };
    }

    async addBalance(accountId: number, delta: CashBalance, options: RepositoryTxOptions = {}): Promise<void> {
        const client = this.ctx.pgClient(this.ctx.getPgUnifiedPool(), options);
        await client.query(
            "UPDATE accounts SET nx_cash = nx_cash + $2, maple_point = maple_point + $3 WHERE id = $1",
            [accountId, delta.nxCash, delta.maplePoint]
        );
    }

    async findItems(worldId: number, accountId: number, options: RepositoryTxOptions = {}): Promise<CashItemRow[]> {
        const client = this.ctx.pgClient(this.ctx.getPgDataPool(worldId, accountId), options);
        const res = await client.query(
            "SELECT serial, account_id, item FROM cash_items WHERE account_id = $1 ORDER BY created_at, serial",
            [accountId]
        );
        return res.rows as CashItemRow[];
    }

    async insertItem(worldId: number, accountId: number, serial: string, item: unknown, options: RepositoryTxOptions = {}): Promise<void> {
        const client = this.ctx.pgClient(this.ctx.getPgDataPool(worldId, accountId), options);
        await client.query(
            "INSERT INTO cash_items (serial, account_id, item) VALUES ($1, $2, $3)",
            [serial, accountId, JSON.stringify(item)]
        );
    }

    // An unopened gift is the item itself, so it leaves with the item on every path out of the locker.
    async deleteItem(worldId: number, accountId: number, serial: string, options: RepositoryTxOptions = {}): Promise<CashItemRow | null> {
        const client = this.ctx.pgClient(this.ctx.getPgDataPool(worldId, accountId), options);
        const res = await client.query(
            "DELETE FROM cash_items WHERE account_id = $1 AND serial = $2 RETURNING serial, account_id, item",
            [accountId, serial]
        );
        const row = (res.rows[0] as CashItemRow | undefined) ?? null;
        if (row == null) {
            return null;
        }
        await client.query("DELETE FROM cash_gifts WHERE account_id = $1 AND serial = $2", [accountId, serial]);
        return row;
    }

    async insertGift(worldId: number, accountId: number, gift: CashGiftRow, options: RepositoryTxOptions = {}): Promise<void> {
        const client = this.ctx.pgClient(this.ctx.getPgDataPool(worldId, accountId), options);
        await client.query(
            "INSERT INTO cash_gifts (serial, account_id, item_id, sender_name, message) VALUES ($1, $2, $3, $4, $5)",
            [gift.serial, accountId, gift.item_id, gift.sender_name, gift.message]
        );
    }

    async findGifts(worldId: number, accountId: number): Promise<CashGiftRow[]> {
        const res = await this.ctx.getPgDataPool(worldId, accountId).query(
            "SELECT serial, item_id, sender_name, message FROM cash_gifts WHERE account_id = $1 ORDER BY serial",
            [accountId]
        );
        return res.rows as CashGiftRow[];
    }

    async deleteItems(worldId: number, accountId: number, serials: string[]): Promise<void> {
        if (serials.length === 0) {
            return;
        }
        await this.withAccount(worldId, accountId, async (client) => {
            await client.query("DELETE FROM cash_items WHERE account_id = $1 AND serial = ANY($2::bigint[])", [accountId, serials]);
            await client.query("DELETE FROM cash_gifts WHERE account_id = $1 AND serial = ANY($2::bigint[])", [accountId, serials]);
        });
    }

    async findCoupon(worldId: number, code: string): Promise<CashCouponRow | null> {
        return this.ctx.withPgGlobalTransaction(worldId, async (client) => {
            const res = await client.query("SELECT code, kind, value, used_at, used_by, request_id FROM cash_coupons WHERE code = $1", [code]);
            return (res.rows[0] as CashCouponRow | undefined) ?? null;
        });
    }

    async claimCoupon(worldId: number, code: string, accountId: number, requestId: string): Promise<CashCouponRow | null> {
        return this.ctx.withPgGlobalTransaction(worldId, async (client) => {
            const res = await client.query(
                `UPDATE cash_coupons SET used_by = $2, used_at = COALESCE(used_at, NOW()), request_id = COALESCE(request_id, $3)
                 WHERE code = $1 AND (used_at IS NULL OR (used_by = $2 AND request_id IS NOT NULL))
                 RETURNING code, kind, value, used_at, used_by, request_id`,
                [code, accountId, requestId]
            );
            return (res.rows[0] as CashCouponRow | undefined) ?? null;
        });
    }

    async insertCoupons(worldId: number, kind: number, value: number, codes: string[]): Promise<string[]> {
        return this.ctx.withPgGlobalTransaction(worldId, async (client) => {
            const res = await client.query(
                "INSERT INTO cash_coupons (code, kind, value) SELECT unnest($1::text[]), $2, $3 ON CONFLICT (code) DO NOTHING RETURNING code",
                [codes, kind, value]
            );
            return (res.rows as Array<{ code: string }>).map((row) => row.code);
        });
    }

    async insertRing(worldId: number, ring: CashRingRow, options: RepositoryTxOptions = {}): Promise<void> {
        await this.ctx.pgClient(this.ctx.getPgDataPool(worldId, ring.character_id), options).query(
            `INSERT INTO cash_rings (serial, partner_serial, character_id, partner_id, partner_name, item_id)
             VALUES ($1, $2, $3, $4, $5, $6)`,
            [ring.serial, ring.partner_serial, ring.character_id, ring.partner_id, ring.partner_name, ring.item_id]
        );
    }

    async findRings(worldId: number, characterId: number): Promise<CashRingRow[]> {
        const res = await this.ctx.getPgDataPool(worldId, characterId).query(
            `SELECT serial, partner_serial, character_id, partner_id, partner_name, item_id
             FROM cash_rings WHERE character_id = $1 ORDER BY serial`,
            [characterId]
        );
        return res.rows as CashRingRow[];
    }

    async findWishlist(worldId: number, characterId: number): Promise<number[]> {
        const res = await this.ctx.getPgDataPool(worldId, characterId).query(
            "SELECT commodity_sns FROM cash_wishlists WHERE character_id = $1",
            [characterId]
        );
        const row = res.rows[0] as { commodity_sns: number[] } | undefined;
        return row?.commodity_sns ?? [];
    }

    async setWishlist(worldId: number, characterId: number, commoditySns: number[]): Promise<void> {
        await this.ctx.getPgDataPool(worldId, characterId).query(
            `INSERT INTO cash_wishlists (character_id, commodity_sns) VALUES ($1, $2)
             ON CONFLICT (character_id) DO UPDATE SET commodity_sns = EXCLUDED.commodity_sns`,
            [characterId, commoditySns]
        );
    }
}

// A transfer records whether it was placed or found the locker full, so every later settle follows the first decision.
export class CashRequestRepository extends RequestRepository<{ full?: boolean }> {
    protected readonly kind = "cash";
}

// Requests whose effect is a balance change, recorded next to the balance in the unified database.
export class CashBalanceRequestRepository extends RequestRepository {
    protected readonly kind = "cash_balance";

    protected override pool(): Pool {
        return this.ctx.getPgUnifiedPool();
    }
}
