import type { PoolClient } from "pg";
import type { InternalContext } from "../context/internal-context";
import type { RepositoryTxOptions } from "../types/repository-contracts";

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

export type CashCouponRow = {
    code: string;
    kind: number;
    value: number;
    used_at: Date | null;
};

export type CashBalance = {
    nxCash: number;
    maplePoint: number;
};

export class CashShopRepository {
    private readonly ctx: InternalContext;

    constructor(internalContext: InternalContext) {
        this.ctx = internalContext;
    }

    async withAccountLock<T>(worldId: number, accountId: number, fn: (client: PoolClient) => Promise<T>): Promise<T> {
        return this.ctx.withPgDataTransaction(worldId, accountId, async (client) => {
            await client.query("SELECT pg_advisory_xact_lock($1, $2)", [0x63617368, accountId]);
            return fn(client);
        });
    }

    async withAccount<T>(worldId: number, accountId: number, fn: (client: PoolClient) => Promise<T>): Promise<T> {
        return this.ctx.withPgDataTransaction(worldId, accountId, fn);
    }

    async findBalance(worldId: number, accountId: number, options: RepositoryTxOptions = {}): Promise<CashBalance> {
        const client = options.txClient ?? this.ctx.getPgDataPool(worldId, accountId);
        const res = await client.query("SELECT nx_cash, maple_point FROM accounts WHERE id = $1", [accountId]);
        const row = res.rows[0] as { nx_cash: number; maple_point: number } | undefined;
        return { nxCash: row?.nx_cash ?? 0, maplePoint: row?.maple_point ?? 0 };
    }

    async addBalance(worldId: number, accountId: number, delta: CashBalance, options: RepositoryTxOptions = {}): Promise<void> {
        const client = options.txClient ?? this.ctx.getPgDataPool(worldId, accountId);
        await client.query(
            "UPDATE accounts SET nx_cash = nx_cash + $2, maple_point = maple_point + $3 WHERE id = $1",
            [accountId, delta.nxCash, delta.maplePoint]
        );
    }

    async findItems(worldId: number, accountId: number, options: RepositoryTxOptions = {}): Promise<CashItemRow[]> {
        const client = options.txClient ?? this.ctx.getPgDataPool(worldId, accountId);
        const res = await client.query(
            "SELECT serial, account_id, item FROM cash_items WHERE account_id = $1 ORDER BY created_at, serial",
            [accountId]
        );
        return res.rows as CashItemRow[];
    }

    async insertItem(worldId: number, accountId: number, serial: string, item: unknown, options: RepositoryTxOptions = {}): Promise<void> {
        const client = options.txClient ?? this.ctx.getPgDataPool(worldId, accountId);
        await client.query(
            "INSERT INTO cash_items (serial, account_id, item) VALUES ($1, $2, $3)",
            [serial, accountId, JSON.stringify(item)]
        );
    }

    async deleteItem(worldId: number, accountId: number, serial: string, options: RepositoryTxOptions = {}): Promise<CashItemRow | null> {
        const client = options.txClient ?? this.ctx.getPgDataPool(worldId, accountId);
        const res = await client.query(
            "DELETE FROM cash_items WHERE account_id = $1 AND serial = $2 RETURNING serial, account_id, item",
            [accountId, serial]
        );
        return (res.rows[0] as CashItemRow | undefined) ?? null;
    }

    async insertGift(worldId: number, accountId: number, gift: CashGiftRow, options: RepositoryTxOptions = {}): Promise<void> {
        const client = options.txClient ?? this.ctx.getPgDataPool(worldId, accountId);
        await client.query(
            "INSERT INTO cash_gifts (serial, account_id, item_id, sender_name, message) VALUES ($1, $2, $3, $4, $5)",
            [gift.serial, accountId, gift.item_id, gift.sender_name, gift.message]
        );
    }

    async takeGifts(worldId: number, accountId: number): Promise<CashGiftRow[]> {
        const res = await this.ctx.getPgDataPool(worldId, accountId).query(
            "DELETE FROM cash_gifts WHERE account_id = $1 RETURNING serial, item_id, sender_name, message",
            [accountId]
        );
        return res.rows as CashGiftRow[];
    }

    async deleteItems(worldId: number, accountId: number, serials: string[]): Promise<void> {
        await this.ctx.getPgDataPool(worldId, accountId).query(
            "DELETE FROM cash_items WHERE account_id = $1 AND serial = ANY($2::bigint[])",
            [accountId, serials]
        );
    }

    async findCoupon(worldId: number, code: string): Promise<CashCouponRow | null> {
        return this.ctx.withPgGlobalTransaction(worldId, async (client) => {
            const res = await client.query("SELECT code, kind, value, used_at FROM cash_coupons WHERE code = $1", [code]);
            return (res.rows[0] as CashCouponRow | undefined) ?? null;
        });
    }

    async withCouponClaim<T>(
        worldId: number,
        code: string,
        accountId: number,
        fn: (coupon: CashCouponRow) => Promise<T>
    ): Promise<T | null> {
        return this.ctx.withPgGlobalTransaction(worldId, async (client) => {
            const res = await client.query(
                "UPDATE cash_coupons SET used_by = $2, used_at = NOW() WHERE code = $1 AND used_at IS NULL RETURNING code, kind, value, used_at",
                [code, accountId]
            );
            const coupon = res.rows[0] as CashCouponRow | undefined;
            if (coupon == null) {
                return null;
            }
            return fn(coupon);
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
