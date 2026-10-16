import type { Pool, PoolClient } from "pg";
import type { InternalContext } from "../context/internal-context";
import type { RepositoryTxOptions } from "../types/repository-contracts";

export enum RequestStatus {
    Applied = 0,
    Aborted = 1,
    Completed = 2,
}

export type RequestRow<TOutcome> = { requestId: string; ownerId: number } & (
    | { status: RequestStatus.Applied | RequestStatus.Completed; outcome: TOutcome }
    | { status: RequestStatus.Aborted; outcome: Partial<TOutcome> }
);

type RequestRecord = {
    request_id: string;
    owner_id: number;
    status: RequestStatus;
    outcome: Record<string, unknown>;
};

const SELECT_COLS = "request_id, owner_id, status, outcome";

// A request record lives in the database that holds the request's effect and is written in that transaction.
export abstract class RequestRepository<TOutcome extends object = Record<string, never>> {
    protected readonly ctx: InternalContext;
    protected abstract readonly kind: string;

    constructor(internalContext: InternalContext) {
        this.ctx = internalContext;
    }

    // A repository whose records live in the global DB overrides this.
    protected pool(worldId: number, ownerId: number): Pool {
        return this.ctx.getPgDataPool(worldId, ownerId);
    }

    private tx(worldId: number, ownerId: number, txClient: PoolClient): PoolClient {
        if (this.ctx.pgClient(this.pool(worldId, ownerId), { txClient }) !== txClient) {
            throw new Error(`${this.kind} request: the transaction is not on the request's database`);
        }
        return txClient;
    }

    private toRow(record: RequestRecord): RequestRow<TOutcome> {
        return { requestId: record.request_id, ownerId: record.owner_id, status: record.status, outcome: record.outcome } as RequestRow<TOutcome>;
    }

    async find(worldId: number, ownerId: number, requestId: string, options: RepositoryTxOptions = {}): Promise<RequestRow<TOutcome> | null> {
        const res = await this.ctx
            .pgClient(this.pool(worldId, ownerId), options)
            .query(`SELECT ${SELECT_COLS} FROM requests WHERE request_id = $1 AND kind = $2 AND owner_id = $3`, [requestId, this.kind, ownerId]);
        const record = res.rows[0] as RequestRecord | undefined;
        return record == null ? null : this.toRow(record);
    }

    // A plain insert: a concurrent abort makes it fail, so the effect rolls back with it.
    async insert(worldId: number, ownerId: number, requestId: string, outcome: TOutcome, txClient: PoolClient): Promise<void> {
        await this.tx(worldId, ownerId, txClient).query("INSERT INTO requests (request_id, kind, owner_id, outcome) VALUES ($1, $2, $3, $4)", [
            requestId,
            this.kind,
            ownerId,
            JSON.stringify(outcome),
        ]);
    }

    async reserve(worldId: number, ownerId: number, requestId: string, outcome: TOutcome, txClient: PoolClient): Promise<boolean> {
        const res = await this.tx(worldId, ownerId, txClient).query(
            "INSERT INTO requests (request_id, kind, owner_id, outcome) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING RETURNING request_id",
            [requestId, this.kind, ownerId, JSON.stringify(outcome)]
        );
        return (res.rowCount ?? 0) > 0;
    }

    // The tombstone goes in first, so a request that arrives later is rejected. Null when the id belongs to another owner.
    async abort(worldId: number, ownerId: number, requestId: string, txClient: PoolClient): Promise<RequestRow<TOutcome> | null> {
        await this.tx(worldId, ownerId, txClient).query(
            "INSERT INTO requests (request_id, kind, owner_id, status) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING",
            [requestId, this.kind, ownerId, RequestStatus.Aborted]
        );
        return this.find(worldId, ownerId, requestId, { txClient });
    }

    async transition(
        worldId: number,
        ownerId: number,
        requestId: string,
        from: RequestStatus.Applied | RequestStatus.Completed,
        to: RequestStatus,
        txClient: PoolClient
    ): Promise<{ requestId: string; ownerId: number; status: RequestStatus; outcome: TOutcome } | null> {
        const res = await this.tx(worldId, ownerId, txClient).query(
            `UPDATE requests SET status = $5 WHERE request_id = $1 AND kind = $2 AND owner_id = $3 AND status = $4 RETURNING ${SELECT_COLS}`,
            [requestId, this.kind, ownerId, from, to]
        );
        const record = res.rows[0] as RequestRecord | undefined;
        if (record == null) {
            return null;
        }
        return { requestId: record.request_id, ownerId: record.owner_id, status: record.status, outcome: record.outcome as TOutcome };
    }
}
