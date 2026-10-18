import type { Pool } from "pg";
import type { RankingBounds, RankingInput, RankingRepository } from "@fm/shared/repos/ranking-repository";

const FETCH_SIZE = 2000;

// Reads one shard's part of a range in key order, a batch at a time. The inputs are frozen, so each batch is a
// short query that resumes after the last key; no transaction stays open while the range is merged.
export class RangeReader {
    private readonly rankings: RankingRepository;
    private readonly pool: Pool;
    private readonly roundId: number;
    private readonly bounds: RankingBounds;
    private rows: RankingInput[] = [];
    private index = 0;
    private exhausted = false;

    constructor(rankings: RankingRepository, pool: Pool, roundId: number, bounds: RankingBounds) {
        this.rankings = rankings;
        this.pool = pool;
        this.roundId = roundId;
        this.bounds = bounds;
    }

    async peek(): Promise<RankingInput | undefined> {
        if (this.index >= this.rows.length && !this.exhausted) {
            const after = this.rows.length > 0 ? this.rows[this.rows.length - 1]!.keys : null;
            this.rows = await this.rankings.readRange(this.pool, this.roundId, this.bounds, after, FETCH_SIZE);
            this.index = 0;
            this.exhausted = this.rows.length < FETCH_SIZE;
        }
        return this.rows[this.index];
    }

    take(): RankingInput {
        const row = this.rows[this.index]!;
        this.index += 1;
        return row;
    }
}
