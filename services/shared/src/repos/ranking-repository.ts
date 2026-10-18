import type { Pool } from "pg";
import type { InternalContext } from "../context/internal-context";
import type { RankingKind } from "../ranking-filters";

export type RankingKeys = [number, number, number, number];

export type RankingBounds = {
    lo: RankingKeys | null;
    hi: RankingKeys | null;
};

export type RankingInput = {
    keys: RankingKeys;
    subjectId: number;
    name: string;
    worldId: number;
    classId: number | null;
    jobGroup: number | null;
    jobId: number | null;
    level: number | null;
    value: number;
};

export type RankingCombination = {
    worldId: number;
    jobGroup: number | null;
    jobId: number | null;
    count: number;
};

export type RankingRow = {
    filter: string;
    rankNo: number;
    input: RankingInput;
};

export type PublishedRound = {
    roundId: number;
    periodToken: string;
    table: string;
};

type InputRow = {
    k1: string;
    k2: string;
    k3: string;
    k4: string;
    subject_id: string;
    name: string;
    world_id: number;
    class_id: number | null;
    job_group: number | null;
    job_id: number | null;
    level: number | null;
    value: string;
};

const DELETE_BATCH = 10000;
const KEY_COLUMNS = "(k1, k2, k3, k4)";
const ROUND_TABLE = /^ranking_entry_(overall|fame|guild)_(\d+)_(\w+)$/;

// A round is written into its own table and published by attaching that table to ranking_entry, so readers see
// only finished rounds. Inputs are a per-shard copy of the ranked rows taken when the round starts.
export class RankingRepository {
    private readonly ctx: InternalContext;

    constructor(internalContext: InternalContext) {
        this.ctx = internalContext;
    }

    private unified(): Pool {
        return this.ctx.getPgUnifiedPool();
    }

    static periodToken(period: string): string {
        return period.replace(/[^0-9A-Za-z]/g, "_");
    }

    static roundTable(kind: RankingKind, roundId: number, period: string): string {
        return `ranking_entry_${kind}_${roundId}_${RankingRepository.periodToken(period)}`;
    }

    async nextRoundId(): Promise<number> {
        const { rows } = await this.unified().query("SELECT nextval('ranking_round_seq') AS id");
        return Number((rows[0] as { id: string }).id);
    }

    async createRoundTable(table: string, roundId: number): Promise<void> {
        await this.unified().query(`DROP TABLE IF EXISTS ${table}`);
        await this.unified().query(`CREATE TABLE ${table} (LIKE ranking_entry INCLUDING ALL)`);
        await this.unified().query(`ALTER TABLE ${table} ADD CONSTRAINT ${table}_round CHECK (round_id = ${Number(roundId)})`);
    }

    async insertRows(table: string, roundId: number, rows: RankingRow[]): Promise<void> {
        await this.unified().query(
            `INSERT INTO ${table} (round_id, filter, rank_no, subject_id, name, world_id, class_id, level, value)
             SELECT $1, e.filter, e.rank_no, e.subject_id, e.name, e.world_id, e.class_id, e.level, e.value
             FROM unnest($2::text[], $3::bigint[], $4::bigint[], $5::text[], $6::smallint[], $7::smallint[], $8::smallint[], $9::bigint[])
                 AS e(filter, rank_no, subject_id, name, world_id, class_id, level, value)
             ON CONFLICT DO NOTHING`,
            [
                roundId,
                rows.map((r) => r.filter),
                rows.map((r) => r.rankNo),
                rows.map((r) => r.input.subjectId),
                rows.map((r) => r.input.name),
                rows.map((r) => r.input.worldId),
                rows.map((r) => r.input.classId),
                rows.map((r) => r.input.level),
                rows.map((r) => r.input.value),
            ]
        );
    }

    async filterTotals(table: string): Promise<Map<string, { count: number; last: number }>> {
        const { rows } = await this.unified().query(`SELECT filter, COUNT(*)::bigint AS n, MAX(rank_no) AS last FROM ${table} GROUP BY filter`);
        const totals = new Map<string, { count: number; last: number }>();
        for (const row of rows as Array<{ filter: string; n: string; last: string }>) {
            totals.set(row.filter, { count: Number(row.n), last: Number(row.last) });
        }
        return totals;
    }

    async attach(table: string, roundId: number): Promise<void> {
        await this.unified().query(`ALTER TABLE ranking_entry ATTACH PARTITION ${table} FOR VALUES IN (${Number(roundId)})`);
    }

    async publishedRounds(kind: RankingKind): Promise<PublishedRound[]> {
        const { rows } = await this.unified().query(
            `SELECT c.relname FROM pg_inherits i JOIN pg_class c ON c.oid = i.inhrelid
             WHERE i.inhparent = 'ranking_entry'::regclass`
        );
        return RankingRepository.rounds(rows as Array<{ relname: string }>, kind).sort((a, b) => b.roundId - a.roundId);
    }

    async unattachedTables(kind: RankingKind): Promise<PublishedRound[]> {
        const { rows } = await this.unified().query(
            `SELECT c.relname FROM pg_class c
             WHERE c.relkind = 'r' AND c.relname LIKE 'ranking\\_entry\\_%'
               AND NOT EXISTS (SELECT 1 FROM pg_inherits i WHERE i.inhrelid = c.oid)`
        );
        return RankingRepository.rounds(rows as Array<{ relname: string }>, kind);
    }

    async detachAndDrop(table: string): Promise<void> {
        await this.unified().query(`ALTER TABLE ranking_entry DETACH PARTITION ${table}`);
        await this.unified().query(`DROP TABLE ${table}`);
    }

    async dropTable(table: string): Promise<void> {
        await this.unified().query(`DROP TABLE IF EXISTS ${table}`);
    }

    async clearInputs(shard: Pool, kind: RankingKind): Promise<void> {
        await this.deleteInBatches(shard, "ranking_inputs", "kind = $1", [kind]);
    }

    async freezeInputs(shard: Pool, select: string): Promise<number> {
        const res = await shard.query(
            `INSERT INTO ranking_inputs (kind, round_id, k1, k2, k3, k4, subject_id, name, world_id, class_id, job_group, job_id, level, value) ${select}`
        );
        return res.rowCount ?? 0;
    }

    async sampleKeys(shard: Pool, roundId: number, fraction: number): Promise<RankingKeys[]> {
        const { rows } = await shard.query("SELECT k1, k2, k3, k4 FROM ranking_inputs WHERE round_id = $1 AND random() < $2", [roundId, fraction]);
        return (rows as InputRow[]).map((row) => [Number(row.k1), Number(row.k2), Number(row.k3), Number(row.k4)]);
    }

    async countRange(shard: Pool, roundId: number, bounds: RankingBounds): Promise<RankingCombination[]> {
        const values: unknown[] = [roundId];
        const where = RankingRepository.boundsSql(bounds, values);
        const { rows } = await shard.query(
            `SELECT world_id, job_group, job_id, COUNT(*)::bigint AS n FROM ranking_inputs
             WHERE round_id = $1${where} GROUP BY world_id, job_group, job_id`,
            values
        );
        return (rows as Array<{ world_id: number; job_group: number | null; job_id: number | null; n: string }>).map((row) => ({
            worldId: row.world_id,
            jobGroup: row.job_group,
            jobId: row.job_id,
            count: Number(row.n),
        }));
    }

    async readRange(shard: Pool, roundId: number, bounds: RankingBounds, after: RankingKeys | null, limit: number): Promise<RankingInput[]> {
        const values: unknown[] = [roundId];
        let where = RankingRepository.boundsSql(bounds, values);
        if (after !== null) {
            values.push(...after);
            where += ` AND ${KEY_COLUMNS} > ($${values.length - 3}, $${values.length - 2}, $${values.length - 1}, $${values.length})`;
        }
        values.push(limit);
        const { rows } = await shard.query(
            `SELECT k1, k2, k3, k4, subject_id, name, world_id, class_id, job_group, job_id, level, value FROM ranking_inputs
             WHERE round_id = $1${where}
             ORDER BY k1, k2, k3, k4
             LIMIT $${values.length}`,
            values
        );
        return (rows as InputRow[]).map((row) => ({
            keys: [Number(row.k1), Number(row.k2), Number(row.k3), Number(row.k4)],
            subjectId: Number(row.subject_id),
            name: row.name,
            worldId: row.world_id,
            classId: row.class_id,
            jobGroup: row.job_group,
            jobId: row.job_id,
            level: row.level,
            value: Number(row.value),
        }));
    }

    // Deleting a shard's rows in batches keeps every statement short, however large the shard grows.
    private async deleteInBatches(pool: Pool, table: string, where: string, values: unknown[]): Promise<void> {
        while (true) {
            const res = await pool.query(
                `DELETE FROM ${table} WHERE ctid IN (SELECT ctid FROM ${table} WHERE ${where} LIMIT ${DELETE_BATCH})`,
                values
            );
            if ((res.rowCount ?? 0) < DELETE_BATCH) {
                return;
            }
        }
    }

    private static rounds(rows: Array<{ relname: string }>, kind: RankingKind): PublishedRound[] {
        const rounds: PublishedRound[] = [];
        for (const { relname } of rows) {
            const match = ROUND_TABLE.exec(relname);
            if (match === null || match[1] !== kind) {
                continue;
            }
            rounds.push({ roundId: Number(match[2]), periodToken: match[3]!, table: relname });
        }
        return rounds;
    }

    private static boundsSql(bounds: RankingBounds, values: unknown[]): string {
        let sql = "";
        if (bounds.lo !== null) {
            values.push(...bounds.lo);
            sql += ` AND ${KEY_COLUMNS} >= ($${values.length - 3}, $${values.length - 2}, $${values.length - 1}, $${values.length})`;
        }
        if (bounds.hi !== null) {
            values.push(...bounds.hi);
            sql += ` AND ${KEY_COLUMNS} < ($${values.length - 3}, $${values.length - 2}, $${values.length - 1}, $${values.length})`;
        }
        return sql;
    }
}
