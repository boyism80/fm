"use strict";

const { redisCacheKey } = require("../redis-cache-key");
const { HashRepository } = require("./hash-repository");

const SELECT_COLS = `character_id, skill_id, level, master_level, cooldown_end_unix_ms, updated_at`;

const INSERT_COLS = `character_id, skill_id, level, master_level, cooldown_end_unix_ms, updated_at`;

const ON_CONFLICT_SET = `
  level = EXCLUDED.level, master_level = EXCLUDED.master_level,
  cooldown_end_unix_ms = EXCLUDED.cooldown_end_unix_ms,
  updated_at = NOW()`;

const PER_ROW_PARAMS = 5;

function rowValues(row) {
    return [
        row.character_id, row.skill_id,
        row.level,        row.master_level,
        row.cooldown_end_unix_ms ?? null,
    ];
}

class SkillRepository extends HashRepository {
    constructor(internalContext) {
        super(internalContext);
    }

    getTtlSeconds() {
        return this.ctx.appConfiguration.getCharacterCacheTtlSeconds();
    }

    getGroupKey(model) {
        return model.characterId;
    }

    getItemKey(model) {
        return String(model.skillId);
    }

    getRedisHashKey(worldId, characterId) {
        return redisCacheKey(`w${worldId}:skills:${characterId}`);
    }

    onSelect(characterId, _worldId) {
        return {
            text:   `SELECT ${SELECT_COLS} FROM character_skills WHERE character_id = $1`,
            values: [Number(characterId)],
        };
    }

    onBulkUpsert(rows) {
        if (!rows.length) return { text: "", values: [] };
        const placeholders = rows
            .map((_, ri) =>
                `(${Array.from({ length: PER_ROW_PARAMS }, (_, ci) => `$${ri * PER_ROW_PARAMS + ci + 1}`).join(",")},NOW())`
            )
            .join(",\n");
        return {
            text:   `INSERT INTO character_skills (${INSERT_COLS}) VALUES\n${placeholders}\nON CONFLICT (character_id, skill_id) DO UPDATE SET${ON_CONFLICT_SET} RETURNING ${SELECT_COLS}`,
            values: rows.flatMap(rowValues),
        };
    }

    onBulkDelete(itemKeys, characterId, _worldId) {
        return {
            text:   `DELETE FROM character_skills WHERE skill_id = ANY($1::bigint[]) AND character_id = $2`,
            values: [itemKeys.map(Number), Number(characterId)],
        };
    }

    rowToModel(row) {
        return {
            characterId:        Number(row.character_id),
            skillId:            Number(row.skill_id),
            level:              Number(row.level),
            masterLevel:        Number(row.master_level),
            cooldownEndUnixMs:  row.cooldown_end_unix_ms != null ? Number(row.cooldown_end_unix_ms) : null,
            updatedAt:          row.updated_at instanceof Date ? row.updated_at : new Date(row.updated_at),
        };
    }

    modelToRow(model) {
        return {
            character_id:         model.characterId,
            skill_id:             model.skillId,
            level:                model.level,
            master_level:         model.masterLevel,
            cooldown_end_unix_ms: model.cooldownEndUnixMs ?? null,
        };
    }

    async replaceBySnapshot(worldId, characterId, models) {
        const pool = this._pool(worldId, characterId);
        const { text: selectText, values: selectValues } = this.onSelect(characterId, worldId);
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
            const deleteSkillIds = existingRows
                .map((r) => String(r.skill_id))
                .filter((id) => !incomingSkillIds.has(id));
            if (deleteSkillIds.length > 0) {
                const delQuery = this.onBulkDelete(deleteSkillIds, characterId, worldId);
                await pool.query(delQuery.text, delQuery.values);
            }

            await pool.query("COMMIT");
        } catch (err) {
            await pool.query("ROLLBACK");
            throw err;
        }

        const redis = this._redis(worldId, characterId);
        await redis.del(this.getRedisHashKey(worldId, characterId)).catch(() => {});
    }

}

module.exports = { SkillRepository };
