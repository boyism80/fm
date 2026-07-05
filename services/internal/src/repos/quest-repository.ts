import { redisCacheKey } from "../redis-cache-key";
import { toPgInt, toPgIntOrNull } from "./pg-int";
import { HashRepository } from "./hash-repository";
import type { RepositoryQuery } from "../types/repository-contracts";
import type { QuestModel, QuestRow } from "../types/repository-models";

const SELECT_COLS =
    "character_id, quest_id, status, mob_kills, status_record, unknown2, completion_time_unix_ms, forfeited, updated_at";
const INSERT_COLS =
    "character_id, quest_id, status, mob_kills, status_record, unknown2, completion_time_unix_ms, forfeited, updated_at";
const ON_CONFLICT_SET = `
  status = EXCLUDED.status,
  mob_kills = EXCLUDED.mob_kills,
  status_record = EXCLUDED.status_record,
  unknown2 = EXCLUDED.unknown2,
  completion_time_unix_ms = EXCLUDED.completion_time_unix_ms,
  forfeited = EXCLUDED.forfeited,
  updated_at = NOW()`;
const PER_ROW_PARAMS = 8;

export type { QuestModel };

function parseNumberMap(raw: QuestRow["mob_kills"]): Record<string, number> {
    if (typeof raw === "string") {
        try {
            const parsed = JSON.parse(raw) as Record<string, number>;
            return parsed && typeof parsed === "object" ? parsed : {};
        } catch {
            return {};
        }
    }
    return raw ?? {};
}

function parseStringMap(raw: QuestRow["unknown2"]): Record<string, string> {
    if (typeof raw === "string") {
        try {
            const parsed = JSON.parse(raw) as Record<string, string>;
            return parsed && typeof parsed === "object" ? parsed : {};
        } catch {
            return {};
        }
    }
    return raw ?? {};
}

function rowValues(row: QuestRow) {
    return [
        row.character_id,
        row.quest_id,
        row.status,
        JSON.stringify(parseNumberMap(row.mob_kills)),
        row.status_record,
        JSON.stringify(parseStringMap(row.unknown2)),
        row.completion_time_unix_ms ?? null,
        row.forfeited,
    ];
}

export class QuestRepository extends HashRepository<QuestModel, QuestRow> {
    override getTtlSeconds() {
        return this.ctx.appConfiguration.getCharacterCacheTtlSeconds();
    }

    override getGroupKey(model: QuestModel) {
        return String(model.characterId);
    }

    override getItemKey(model: QuestModel) {
        return String(model.questId);
    }

    override getRedisHashKey(worldId: number, characterId: string) {
        return redisCacheKey(`w${worldId}:quests:${characterId}`);
    }

    override onSelect(characterId: string): RepositoryQuery {
        return {
            text: `SELECT ${SELECT_COLS} FROM character_quests WHERE character_id = $1`,
            values: [Number(characterId)],
        };
    }

    override onBulkUpsert(rows: QuestRow[]): RepositoryQuery {
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
            text: `INSERT INTO character_quests (${INSERT_COLS}) VALUES\n${placeholders}\nON CONFLICT (character_id, quest_id) DO UPDATE SET${ON_CONFLICT_SET} RETURNING ${SELECT_COLS}`,
            values: rows.flatMap(rowValues),
        };
    }

    override onBulkDelete(itemKeys: string[], characterId: string): RepositoryQuery {
        return {
            text: "DELETE FROM character_quests WHERE quest_id = ANY($1::int[]) AND character_id = $2",
            values: [itemKeys.map((k) => Number(k)), Number(characterId)],
        };
    }

    override normalizeRow(row: QuestRow): QuestRow {
        return {
            ...row,
            quest_id: toPgInt(row.quest_id),
            status: toPgInt(row.status),
            completion_time_unix_ms: toPgIntOrNull(row.completion_time_unix_ms),
            forfeited: toPgInt(row.forfeited),
        };
    }

    override rowToModel(row: QuestRow): QuestModel {
        return {
            characterId: row.character_id,
            questId: toPgInt(row.quest_id),
            status: toPgInt(row.status),
            mobKills: parseNumberMap(row.mob_kills),
            statusRecord: row.status_record ?? "",
            unknown2: parseStringMap(row.unknown2),
            completionTimeUnixMs: toPgInt(row.completion_time_unix_ms),
            forfeited: toPgInt(row.forfeited),
            updatedAt: row.updated_at instanceof Date ? row.updated_at : row.updated_at ? new Date(row.updated_at) : undefined,
        };
    }

    override modelToRow(model: QuestModel): QuestRow {
        return {
            character_id: model.characterId,
            quest_id: model.questId,
            status: model.status,
            mob_kills: model.mobKills ?? {},
            status_record: model.statusRecord ?? "",
            unknown2: model.unknown2 ?? {},
            completion_time_unix_ms: model.completionTimeUnixMs || null,
            forfeited: model.forfeited ?? 0,
        };
    }

    async replaceBySnapshot(worldId: number, characterId: number, models: QuestModel[]) {
        const groupKey = String(characterId);
        const pool = this.pool(worldId, groupKey);
        const { text: selectText, values: selectValues } = this.onSelect(groupKey);
        const existingRes = await pool.query(selectText, selectValues);
        const existingRows = existingRes.rows ?? [];

        const normalized = models.map((m) => {
            const row = this.modelToRow(m);
            row.character_id = characterId;
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

            const incoming = new Set(normalized.map((r) => String(r.quest_id)));
            const deleteIds = (existingRows as QuestRow[])
                .map((r) => String(r.quest_id))
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
