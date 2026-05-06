import { redisCacheKey } from "../redis-cache-key";
import { HashRepository } from "./hash-repository";
import type { RepositoryQuery } from "../types/repository-contracts";
import type { BuffModel, BuffRow, BuffFlagValueRow } from "../types/repository-models";

const SELECT_COLS =
    "character_id, buff_source_id, kind, remaining_duration_ms, skill_level, causer_id, flag_values, updated_at";
const INSERT_COLS =
    "character_id, buff_source_id, kind, remaining_duration_ms, skill_level, causer_id, flag_values, updated_at";
const ON_CONFLICT_SET = `
  kind = EXCLUDED.kind,
  remaining_duration_ms = EXCLUDED.remaining_duration_ms,
  skill_level = EXCLUDED.skill_level,
  causer_id = EXCLUDED.causer_id,
  flag_values = EXCLUDED.flag_values,
  updated_at = NOW()`;
const PER_ROW_PARAMS = 7;

export type { BuffModel };

function parseFlagValues(raw: BuffRow["flag_values"]): BuffFlagValueRow[] {
    if (Array.isArray(raw)) {
        return raw;
    }
    if (typeof raw === "string") {
        try {
            const parsed = JSON.parse(raw) as unknown;
            return Array.isArray(parsed) ? (parsed as BuffFlagValueRow[]) : [];
        } catch {
            return [];
        }
    }
    return [];
}

function rowValues(row: BuffRow) {
    return [
        row.character_id,
        row.buff_source_id,
        row.kind,
        row.remaining_duration_ms ?? null,
        row.skill_level ?? null,
        row.causer_id ?? null,
        JSON.stringify(parseFlagValues(row.flag_values)),
    ];
}

export class BuffRepository extends HashRepository<BuffModel, BuffRow> {
    getTtlSeconds() {
        return this.ctx.appConfiguration.getCharacterCacheTtlSeconds();
    }

    getGroupKey(model: BuffModel) {
        return String(model.characterId);
    }

    getItemKey(model: BuffModel) {
        return String(model.buffSourceId);
    }

    getRedisHashKey(worldId: number, characterId: string) {
        return redisCacheKey(`w${worldId}:buffs:${characterId}`);
    }

    onSelect(characterId: string): RepositoryQuery {
        return {
            text: `SELECT ${SELECT_COLS} FROM character_buffs WHERE character_id = $1`,
            values: [Number(characterId)],
        };
    }

    onBulkUpsert(rows: BuffRow[]): RepositoryQuery {
        if (!rows.length) {
            return { text: "", values: [] };
        }
        const placeholders = rows
            .map(
                (_, ri) =>
                    `(${Array.from({ length: PER_ROW_PARAMS }, (_, ci) => `$${ri * PER_ROW_PARAMS + ci + 1}`).join(",")},NOW())`
            )
            .join(",\n");
        return {
            text: `INSERT INTO character_buffs (${INSERT_COLS}) VALUES\n${placeholders}\nON CONFLICT (character_id, buff_source_id) DO UPDATE SET${ON_CONFLICT_SET} RETURNING ${SELECT_COLS}`,
            values: rows.flatMap(rowValues),
        };
    }

    onBulkDelete(itemKeys: string[], characterId: string): RepositoryQuery {
        return {
            text: "DELETE FROM character_buffs WHERE buff_source_id = ANY($1::int[]) AND character_id = $2",
            values: [itemKeys.map((k) => Number(k)), Number(characterId)],
        };
    }

    rowToModel(row: BuffRow): BuffModel {
        const flags = parseFlagValues(row.flag_values);
        return {
            characterId: Number(row.character_id),
            buffSourceId: Number(row.buff_source_id),
            kind: Number(row.kind),
            flagValues: flags.map((f) => ({
                mask: Number(f.mask),
                position: Number(f.position),
                value: Number(f.value),
            })),
            remainingDurationMs:
                row.remaining_duration_ms != null ? Number(row.remaining_duration_ms) : null,
            skillLevel: row.skill_level != null ? Number(row.skill_level) : null,
            causerId: row.causer_id != null ? Number(row.causer_id) : null,
            updatedAt: row.updated_at instanceof Date ? row.updated_at : row.updated_at ? new Date(row.updated_at) : undefined,
        };
    }

    modelToRow(model: BuffModel): BuffRow {
        return {
            character_id: model.characterId,
            buff_source_id: model.buffSourceId,
            kind: model.kind,
            remaining_duration_ms: model.remainingDurationMs,
            skill_level: model.skillLevel,
            causer_id: model.causerId,
            flag_values: model.flagValues.map((f) => ({
                mask: f.mask,
                position: f.position,
                value: f.value,
            })),
        };
    }

    async replaceBySnapshot(worldId: number, characterId: number, models: BuffModel[]) {
        const groupKey = String(characterId);
        const pool = this.pool(worldId, groupKey);
        const { text: selectText, values: selectValues } = this.onSelect(groupKey);
        const existingRes = await pool.query(selectText, selectValues);
        const existingRows = existingRes.rows ?? [];

        const normalized = models.map((m) => {
            const row = this.modelToRow(m);
            row.character_id = Number(characterId);
            return row;
        });

        await pool.query("BEGIN");
        try {
            if (normalized.length > 0) {
                const upsert = this.onBulkUpsert(normalized);
                if (upsert.text) {
                    await pool.query(upsert.text, upsert.values);
                }
            }

            const incoming = new Set(normalized.map((r) => String(r.buff_source_id)));
            const deleteIds = (existingRows as BuffRow[])
                .map((r) => String(r.buff_source_id))
                .filter((id) => !incoming.has(id));
            if (deleteIds.length > 0) {
                const delQuery = this.onBulkDelete(deleteIds, groupKey);
                await pool.query(delQuery.text, delQuery.values);
            }

            await pool.query("COMMIT");
        } catch (err) {
            await pool.query("ROLLBACK");
            throw err;
        }

        const redis = this.redis(worldId, groupKey);
        await redis.del(this.getRedisHashKey(worldId, String(characterId))).catch(() => {});
    }
}
