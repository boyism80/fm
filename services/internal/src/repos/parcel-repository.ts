import type { Pool, PoolClient } from "pg";
import type { InternalContext } from "../context/internal-context";
import type { RepositoryTxOptions } from "../types/repository-contracts";

export type ParcelRow = {
    parcel_id: number;
    receiver_id: number;
    sender_name: string;
    meso: number;
    quick: boolean;
    message: string;
    item: unknown | null;
    notified: boolean;
    sent_at: Date;
};

export type ParcelInsert = {
    receiverId: number;
    senderName: string;
    meso: number;
    quick: boolean;
    message: string;
    item: unknown | null;
};

const SELECT_COLS = "parcel_id, receiver_id, sender_name, meso, quick, message, item, notified, sent_at";

export class ParcelRepository {
    private readonly ctx: InternalContext;

    constructor(internalContext: InternalContext) {
        this.ctx = internalContext;
    }

    pool(worldId: number, receiverId: number): Pool {
        return this.ctx.getPgDataPool(worldId, receiverId);
    }

    async withReceiverLock<T>(worldId: number, receiverId: number, fn: (client: PoolClient) => Promise<T>): Promise<T> {
        return this.ctx.withPgDataTransaction(worldId, receiverId, async (client) => {
            await client.query("SELECT pg_advisory_xact_lock($1, $2)", [0x70617263, receiverId]);
            return fn(client);
        });
    }

    async findByReceiver(worldId: number, receiverId: number, options: RepositoryTxOptions = {}): Promise<ParcelRow[]> {
        const client = options.txClient ?? this.pool(worldId, receiverId);
        const res = await client.query(
            `SELECT ${SELECT_COLS} FROM parcels WHERE receiver_id = $1 ORDER BY parcel_id`,
            [receiverId]
        );
        return res.rows as ParcelRow[];
    }

    async insert(worldId: number, parcel: ParcelInsert, options: RepositoryTxOptions = {}): Promise<ParcelRow> {
        const client = options.txClient ?? this.pool(worldId, parcel.receiverId);
        const res = await client.query(
            `INSERT INTO parcels (receiver_id, sender_name, meso, quick, message, item)
             VALUES ($1, $2, $3, $4, $5, $6)
             RETURNING ${SELECT_COLS}`,
            [
                parcel.receiverId,
                parcel.senderName,
                parcel.meso,
                parcel.quick,
                parcel.message,
                parcel.item === null ? null : JSON.stringify(parcel.item),
            ]
        );
        return res.rows[0] as ParcelRow;
    }

    async delete(worldId: number, receiverId: number, parcelId: number, options: RepositoryTxOptions = {}): Promise<ParcelRow | null> {
        const client = options.txClient ?? this.pool(worldId, receiverId);
        const res = await client.query(
            `DELETE FROM parcels WHERE receiver_id = $1 AND parcel_id = $2 RETURNING ${SELECT_COLS}`,
            [receiverId, parcelId]
        );
        return (res.rows[0] as ParcelRow | undefined) ?? null;
    }

    async deleteMany(worldId: number, receiverId: number, parcelIds: number[], options: RepositoryTxOptions = {}): Promise<void> {
        if (parcelIds.length === 0) {
            return;
        }
        const client = options.txClient ?? this.pool(worldId, receiverId);
        await client.query("DELETE FROM parcels WHERE receiver_id = $1 AND parcel_id = ANY($2::int[])", [receiverId, parcelIds]);
    }

    async markNotified(worldId: number, receiverId: number, parcelIds: number[], options: RepositoryTxOptions = {}): Promise<void> {
        if (parcelIds.length === 0) {
            return;
        }
        const client = options.txClient ?? this.pool(worldId, receiverId);
        await client.query(
            "UPDATE parcels SET notified = TRUE WHERE receiver_id = $1 AND parcel_id = ANY($2::int[])",
            [receiverId, parcelIds]
        );
    }
}
