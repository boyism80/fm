import { redisCacheKey } from "../redis-cache-key";
import { HashRepository } from "./hash-repository";
import type { RepositoryQuery } from "../types/repository-contracts";
import type { SkillModel, SkillRow } from "../types/repository-models";

const SELECT_COLS = "character_id, skill_id, level, master_level, cooldown_end_unix_ms, updated_at";
const INSERT_COLS = "character_id, skill_id, level, master_level, cooldown_end_unix_ms, updated_at";
const ON_CONFLICT_SET = `
  level = EXCLUDED.level, master_level = EXCLUDED.master_level,
  cooldown_end_unix_ms = EXCLUDED.cooldown_end_unix_ms,
  updated_at = NOW()`;
const PER_ROW_PARAMS = 5;

export type { SkillModel };

function rowValues(row: SkillRow) {
    return [
        row.character_id, row.skill_id,
        row.level, row.master_level,
        row.cooldown_end_unix_ms ?? null,
    ];
}

export class SkillRepository extends HashRepository<SkillModel, SkillRow> {
    getTtlSeconds() {
        return this.ctx.appConfiguration.getCharacterCacheTtlSeconds();
    }

    getGroupKey(model: SkillModel) {
        return String(model.characterId);
    }

    getItemKey(model: SkillModel) {
        return String(model.skillId);
    }

    getRedisHashKey(worldId: number, characterId: string) {
        return redisCacheKey(`w${worldId}:skills:${characterId}`);
    }

    onSelect(characterId: string): RepositoryQuery {
        return {
            text: `SELECT ${SELECT_COLS} FROM character_skills WHERE character_id = $1`,
            values: [Number(characterId)],
        };
    }

    onBulkUpsert(rows: SkillRow[]): RepositoryQuery {
        if (!rows.length) {
            return { text: "", values: [] };
        }
        const placeholders = rows
            .map((_, ri) => `(${Array.from({ length: PER_ROW_PARAMS }, (_, ci) => `$${ri * PER_ROW_PARAMS + ci + 1}`).join(",")},NOW())`)
            .join(",\n");
        return {
            text: `INSERT INTO character_skills (${INSERT_COLS}) VALUES\n${placeholders}\nON CONFLICT (character_id, skill_id) DO UPDATE SET${ON_CONFLICT_SET} RETURNING ${SELECT_COLS}`,
            values: rows.flatMap(rowValues),
        };
    }

    onBulkDelete(itemKeys: string[], characterId: string): RepositoryQuery {
        return {
            text: "DELETE FROM character_skills WHERE skill_id = ANY($1::bigint[]) AND character_id = $2",
            values: [itemKeys.map(Number), Number(characterId)],
        };
    }

    rowToModel(row: SkillRow): SkillModel {
        return {
            characterId: Number(row.character_id),
            skillId: Number(row.skill_id),
            level: Number(row.level),
            masterLevel: Number(row.master_level),
            cooldownEndUnixMs: row.cooldown_end_unix_ms != null ? Number(row.cooldown_end_unix_ms) : null,
            updatedAt: row.updated_at instanceof Date ? row.updated_at : row.updated_at ? new Date(row.updated_at) : undefined,
        };
    }

    modelToRow(model: SkillModel): SkillRow {
        return {
            character_id: model.characterId,
            skill_id: model.skillId,
            level: model.level,
            master_level: model.masterLevel,
            cooldown_end_unix_ms: model.cooldownEndUnixMs ?? null,
        };
    }

    async replaceBySnapshot(worldId: number, characterId: number, models: SkillModel[]) {
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

            const incomingSkillIds = new Set(normalized.map((r) => String(r.skill_id)));
            const deleteSkillIds = (existingRows as SkillRow[]).map((r) => String(r.skill_id)).filter((id) => !incomingSkillIds.has(id));
            if (deleteSkillIds.length > 0) {
                const delQuery = this.onBulkDelete(deleteSkillIds, groupKey);
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
