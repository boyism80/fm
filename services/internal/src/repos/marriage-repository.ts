import type { PoolClient } from "pg";
import type { InternalContext } from "../context/internal-context";
import type { RepositoryTxOptions } from "../types/repository-contracts";

export type MarriageRow = {
    marriage_id: number;
    groom_id: number;
    bride_id: number;
    groom_name: string;
    bride_name: string;
    groom_item_id: number;
    bride_item_id: number;
    status: number;
    ticket_item_id: number;
    groom_wished: boolean;
    bride_wished: boolean;
    groom_wishes: string[];
    bride_wishes: string[];
    guests: number[];
    divorce_requester_id: number;
    divorce_requested_at: Date | null;
};

export type MarriageInsert = {
    groomId: number;
    brideId: number;
    groomName: string;
    brideName: string;
    groomItemId: number;
    brideItemId: number;
    status: number;
};

export type WeddingGiftRow = {
    gift_id: number;
    receiver_id: number;
    sender_name: string;
    item: unknown;
};

const SELECT_COLS = `marriage_id, groom_id, bride_id, groom_name, bride_name, groom_item_id, bride_item_id, status,
  ticket_item_id, groom_wished, bride_wished, groom_wishes, bride_wishes, guests, divorce_requester_id, divorce_requested_at`;

const GIFT_COLS = "gift_id, receiver_id, sender_name, item";

export class MarriageRepository {
    private readonly ctx: InternalContext;

    constructor(internalContext: InternalContext) {
        this.ctx = internalContext;
    }

    async withTransaction<T>(worldId: number, fn: (client: PoolClient) => Promise<T>): Promise<T> {
        return this.ctx.withPgGlobalTransaction(worldId, fn);
    }

    async findByCharacter(client: PoolClient, characterId: number): Promise<MarriageRow | null> {
        const res = await client.query(
            `SELECT ${SELECT_COLS} FROM marriages WHERE ended = FALSE AND (groom_id = $1 OR bride_id = $1) FOR UPDATE`,
            [characterId]
        );
        return (res.rows[0] as MarriageRow | undefined) ?? null;
    }

    async find(client: PoolClient, marriageId: number): Promise<MarriageRow | null> {
        const res = await client.query(`SELECT ${SELECT_COLS} FROM marriages WHERE ended = FALSE AND marriage_id = $1 FOR UPDATE`, [
            marriageId,
        ]);
        return (res.rows[0] as MarriageRow | undefined) ?? null;
    }

    async insert(client: PoolClient, marriage: MarriageInsert): Promise<MarriageRow> {
        const res = await client.query(
            `INSERT INTO marriages (groom_id, bride_id, groom_name, bride_name, groom_item_id, bride_item_id, status)
             VALUES ($1, $2, $3, $4, $5, $6, $7)
             RETURNING ${SELECT_COLS}`,
            [
                marriage.groomId,
                marriage.brideId,
                marriage.groomName,
                marriage.brideName,
                marriage.groomItemId,
                marriage.brideItemId,
                marriage.status,
            ]
        );
        return res.rows[0] as MarriageRow;
    }

    async update(client: PoolClient, row: MarriageRow): Promise<void> {
        await client.query(
            `UPDATE marriages SET groom_item_id = $2, bride_item_id = $3, status = $4, ticket_item_id = $5,
               groom_wished = $6, bride_wished = $7, groom_wishes = $8, bride_wishes = $9, guests = $10,
               divorce_requester_id = $11, divorce_requested_at = $12
             WHERE marriage_id = $1`,
            [
                row.marriage_id,
                row.groom_item_id,
                row.bride_item_id,
                row.status,
                row.ticket_item_id,
                row.groom_wished,
                row.bride_wished,
                row.groom_wishes,
                row.bride_wishes,
                row.guests,
                row.divorce_requester_id,
                row.divorce_requested_at,
            ]
        );
    }

    async end(client: PoolClient, marriageId: number): Promise<void> {
        await client.query("UPDATE marriages SET ended = TRUE WHERE marriage_id = $1", [marriageId]);
    }

    async findRequest(client: PoolClient, requestId: string): Promise<{ aborted: boolean } | null> {
        const res = await client.query("SELECT aborted FROM marriage_requests WHERE request_id = $1", [requestId]);
        return (res.rows[0] as { aborted: boolean } | undefined) ?? null;
    }

    // A plain insert: a concurrent abort makes it fail, so the change rolls back with it.
    async insertRequest(client: PoolClient, requestId: string): Promise<void> {
        await client.query("INSERT INTO marriage_requests (request_id) VALUES ($1)", [requestId]);
    }

    async abortRequest(client: PoolClient, requestId: string): Promise<void> {
        await client.query("INSERT INTO marriage_requests (request_id, aborted) VALUES ($1, TRUE) ON CONFLICT DO NOTHING", [requestId]);
    }

    async findGifts(worldId: number, receiverId: number): Promise<WeddingGiftRow[]> {
        const res = await this.ctx
            .getPgDataPool(worldId, receiverId)
            .query(`SELECT ${GIFT_COLS} FROM wedding_gifts WHERE receiver_id = $1 ORDER BY gift_id`, [receiverId]);
        return res.rows as WeddingGiftRow[];
    }

    async insertGift(worldId: number, receiverId: number, senderName: string, item: unknown, options: RepositoryTxOptions = {}): Promise<void> {
        const client = options.txClient ?? this.ctx.getPgDataPool(worldId, receiverId);
        await client.query("INSERT INTO wedding_gifts (receiver_id, sender_name, item) VALUES ($1, $2, $3)", [
            receiverId,
            senderName,
            JSON.stringify(item),
        ]);
    }

    async withReceiverLock<T>(worldId: number, receiverId: number, fn: (client: PoolClient) => Promise<T>): Promise<T> {
        return this.ctx.withPgDataTransaction(worldId, receiverId, async (client) => {
            await client.query("SELECT pg_advisory_xact_lock($1, $2)", [0x67696674, receiverId]);
            return fn(client);
        });
    }

    async deleteGift(client: PoolClient, receiverId: number, giftId: number): Promise<WeddingGiftRow | null> {
        const res = await client.query(`DELETE FROM wedding_gifts WHERE receiver_id = $1 AND gift_id = $2 RETURNING ${GIFT_COLS}`, [receiverId, giftId]);
        return (res.rows[0] as WeddingGiftRow | undefined) ?? null;
    }

    async findGiftSendRequest(client: PoolClient, requestId: string): Promise<{ aborted: boolean } | null> {
        const res = await client.query("SELECT aborted FROM wedding_gift_send_requests WHERE request_id = $1", [requestId]);
        return (res.rows[0] as { aborted: boolean } | undefined) ?? null;
    }

    async insertGiftSendRequest(client: PoolClient, receiverId: number, requestId: string, aborted: boolean): Promise<void> {
        await client.query(
            "INSERT INTO wedding_gift_send_requests (request_id, receiver_id, aborted) VALUES ($1, $2, $3)",
            [requestId, receiverId, aborted]
        );
    }

    async findGiftClaimRequest(client: PoolClient, requestId: string): Promise<unknown | null> {
        const res = await client.query("SELECT gift FROM wedding_gift_claim_requests WHERE request_id = $1", [requestId]);
        return (res.rows[0] as { gift: unknown } | undefined)?.gift ?? null;
    }

    async insertGiftClaimRequest(client: PoolClient, receiverId: number, requestId: string, gift: unknown): Promise<void> {
        await client.query(
            "INSERT INTO wedding_gift_claim_requests (request_id, receiver_id, gift) VALUES ($1, $2, $3)",
            [requestId, receiverId, JSON.stringify(gift)]
        );
    }
}
