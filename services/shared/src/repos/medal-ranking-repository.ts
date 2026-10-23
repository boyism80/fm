import type { PoolClient } from "pg";
import type { InternalContext } from "../context/internal-context";

export const MEDAL_RANKING_SIZE = 5;

export type MedalRankingEntryRow = {
    characterId: number;
    name: string;
    score: number;
};

export type MedalRankingRow = {
    kind: string;
    month: number;
    entries: MedalRankingEntryRow[];
};

type Row = {
    kind: string;
    month: number;
    character_id: number;
    name: string;
    score: string;
};

// Ties keep whoever reached the score first.
const BOARD_ORDER = "score DESC, recorded_at, character_id";

export class MedalRankingRepository {
    private readonly ctx: InternalContext;

    constructor(internalContext: InternalContext) {
        this.ctx = internalContext;
    }

    async list(worldId: number): Promise<MedalRankingRow[]> {
        const { rows } = await this.ctx
            .getPgGlobalPool(worldId)
            .query(
                `SELECT kind, month, character_id, name, score FROM medal_rankings WHERE world_id = $1 ORDER BY kind, month, ${BOARD_ORDER}`,
                [worldId],
            );
        const rankings = new Map<string, MedalRankingRow>();
        for (const row of rows as Row[]) {
            const key = `${row.kind}|${row.month}`;
            const ranking = rankings.get(key) ?? { kind: row.kind, month: row.month, entries: [] };
            ranking.entries.push({ characterId: row.character_id, name: row.name, score: Number(row.score) });
            rankings.set(key, ranking);
        }
        return [...rankings.values()];
    }

    // A character keeps its best score of the month; the board keeps the top MEDAL_RANKING_SIZE.
    async submit(
        worldId: number,
        kind: string,
        month: number,
        characterId: number,
        name: string,
        score: number,
    ): Promise<{ previousHolderId: number; ranking: MedalRankingRow }> {
        return this.ctx.withPgGlobalTransaction(worldId, async (client) => {
            await client.query("SELECT pg_advisory_xact_lock(hashtext($1))", [`medal:${worldId}:${kind}`]);
            await client.query("DELETE FROM medal_rankings WHERE world_id = $1 AND kind = $2 AND month < $3", [
                worldId,
                kind,
                month,
            ]);
            const before = await this.board(client, worldId, kind, month);
            await client.query(
                `INSERT INTO medal_rankings (world_id, kind, month, character_id, name, score) VALUES ($1, $2, $3, $4, $5, $6)
                 ON CONFLICT (world_id, kind, month, character_id) DO UPDATE SET
                     name = EXCLUDED.name,
                     score = GREATEST(medal_rankings.score, EXCLUDED.score),
                     recorded_at = CASE WHEN EXCLUDED.score > medal_rankings.score THEN NOW() ELSE medal_rankings.recorded_at END`,
                [worldId, kind, month, characterId, name, score],
            );
            await client.query(
                `DELETE FROM medal_rankings WHERE world_id = $1 AND kind = $2 AND month = $3 AND character_id NOT IN (
                     SELECT character_id FROM medal_rankings WHERE world_id = $1 AND kind = $2 AND month = $3 ORDER BY ${BOARD_ORDER} LIMIT ${MEDAL_RANKING_SIZE})`,
                [worldId, kind, month],
            );
            const entries = await this.board(client, worldId, kind, month);
            return { previousHolderId: before[0]?.characterId ?? 0, ranking: { kind, month, entries } };
        });
    }

    private async board(
        client: PoolClient,
        worldId: number,
        kind: string,
        month: number,
    ): Promise<MedalRankingEntryRow[]> {
        const { rows } = await client.query(
            `SELECT kind, month, character_id, name, score FROM medal_rankings WHERE world_id = $1 AND kind = $2 AND month = $3 ORDER BY ${BOARD_ORDER}`,
            [worldId, kind, month],
        );
        return (rows as Row[]).map((row) => ({
            characterId: row.character_id,
            name: row.name,
            score: Number(row.score),
        }));
    }
}
