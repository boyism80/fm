"use strict";

const { redisCacheKey } = require("../redis-cache-key");
const { ValueRepository } = require("./value-repository");

const SELECT_COLS = "character_id, world_id, key_layout_json, updated_at";
const INSERT_COLS = "character_id, world_id, key_layout_json, updated_at";
const ON_CONFLICT_SET = "key_layout_json = EXCLUDED.key_layout_json, updated_at = NOW()";

class KeyLayoutRepository extends ValueRepository {
    constructor(internalContext) {
        super(internalContext);
    }

    getKey(model) {
        return model.characterId;
    }

    getTtlSeconds() {
        return this.ctx.appConfiguration.getCharacterCacheTtlSeconds();
    }

    getRedisKey(worldId, characterId) {
        return redisCacheKey(`w${worldId}:keylayout:${characterId}`);
    }

    onSelect(characterId, worldId) {
        return {
            text: `SELECT ${SELECT_COLS} FROM key_layout WHERE character_id = $1 AND world_id = $2`,
            values: [Number(characterId), Number(worldId)],
        };
    }

    onSelectMany(characterIds, worldId) {
        return {
            text: `SELECT ${SELECT_COLS} FROM key_layout WHERE character_id = ANY($1::int[]) AND world_id = $2`,
            values: [characterIds.map(Number), Number(worldId)],
        };
    }

    onUpsert(row) {
        return {
            text: `INSERT INTO key_layout (${INSERT_COLS}) VALUES ($1,$2,$3,NOW()) ON CONFLICT (character_id, world_id) DO UPDATE SET ${ON_CONFLICT_SET} RETURNING ${SELECT_COLS}`,
            values: [row.character_id, row.world_id, row.key_layout_json ?? "{}"],
        };
    }

    onBulkUpsert(rows) {
        if (!rows.length) return null;
        const placeholders = rows
            .map((_, i) => `($${i * 3 + 1},$${i * 3 + 2},$${i * 3 + 3},NOW())`)
            .join(",\n");
        return {
            text: `INSERT INTO key_layout (${INSERT_COLS}) VALUES\n${placeholders}\nON CONFLICT (character_id, world_id) DO UPDATE SET ${ON_CONFLICT_SET} RETURNING ${SELECT_COLS}`,
            values: rows.flatMap((r) => [r.character_id, r.world_id, r.key_layout_json ?? "{}"]),
        };
    }

    onDelete(row) {
        return {
            text: "DELETE FROM key_layout WHERE character_id = $1 AND world_id = $2",
            values: [Number(row.characterId), Number(row.worldId)],
        };
    }

    rowToModel(row) {
        return {
            characterId: Number(row.character_id),
            worldId: Number(row.world_id),
            keyLayoutJson: typeof row.key_layout_json === "string"
                ? row.key_layout_json
                : JSON.stringify(row.key_layout_json ?? {}),
            updatedAt: row.updated_at instanceof Date ? row.updated_at : new Date(row.updated_at),
        };
    }

    modelToRow(model) {
        return {
            character_id: model.characterId,
            world_id: model.worldId,
            key_layout_json: model.keyLayoutJson ?? "{}",
        };
    }
}

module.exports = { KeyLayoutRepository };
