"use strict";

const { HashRepository } = require("./hash-repository");

const SELECT_COLS = `character_id, account_id, world_id, name, gender, skin_color, face, hair,
  level, class_id, map_id, spawn_point, rank, rank_diff, class_rank, class_rank_diff,
  base_looks, overlays, deleted, updated_at`;

const INSERT_COLS = `character_id, account_id, world_id, name, gender, skin_color, face, hair,
  level, class_id, map_id, spawn_point, rank, rank_diff, class_rank, class_rank_diff,
  base_looks, overlays, deleted, updated_at`;

const ON_CONFLICT_SET = `
  name = EXCLUDED.name, gender = EXCLUDED.gender, skin_color = EXCLUDED.skin_color,
  face = EXCLUDED.face, hair = EXCLUDED.hair, level = EXCLUDED.level,
  class_id = EXCLUDED.class_id, map_id = EXCLUDED.map_id, spawn_point = EXCLUDED.spawn_point,
  rank = EXCLUDED.rank, rank_diff = EXCLUDED.rank_diff,
  class_rank = EXCLUDED.class_rank, class_rank_diff = EXCLUDED.class_rank_diff,
  base_looks = EXCLUDED.base_looks, overlays = EXCLUDED.overlays,
  deleted = FALSE, updated_at = NOW()`;

const PER_ROW_PARAMS = 18;

function rowValues(row) {
    return [
        row.character_id, row.account_id,   row.world_id,      row.name,
        row.gender,       row.skin_color,    row.face,          row.hair,
        row.level,        row.class_id,      row.map_id,        row.spawn_point,
        row.rank,         row.rank_diff,     row.class_rank,    row.class_rank_diff,
        JSON.stringify(row.base_looks ?? {}),
        JSON.stringify(row.overlays ?? {}),
    ];
}

class CharacterOverviewRepository extends HashRepository {
    constructor(internalContext) {
        super(internalContext);
    }

    getShardHash(ownerKey) {
        return ownerKey;
    }

    getTtlSeconds() {
        return this.ctx.appConfiguration.getCharacterCacheTtlSeconds();
    }

    getGroupKey(model) {
        return model.accountId;
    }

    getItemKey(model) {
        return String(model.characterId);
    }

    getRedisHashKey(worldId, accountId) {
        const { keyPrefix } = this.ctx.getRedisDataAccess(worldId, accountId);
        return `${keyPrefix}fm:w${worldId}:overview:${accountId}`;
    }

    onSelect(accountId, worldId) {
        return {
            text: `SELECT ${SELECT_COLS} FROM character_overview WHERE account_id = $1 AND world_id = $2 AND NOT deleted`,
            values: [Number(accountId), Number(worldId)],
        };
    }

    onBulkUpsert(rows) {
        if (!rows.length) return { text: "", values: [] };
        const placeholders = rows
            .map((_, ri) =>
                `(${Array.from({ length: PER_ROW_PARAMS }, (_, ci) => `$${ri * PER_ROW_PARAMS + ci + 1}`).join(",")},FALSE,NOW())`
            )
            .join(",\n");
        return {
            text: `INSERT INTO character_overview (${INSERT_COLS}) VALUES\n${placeholders}\nON CONFLICT (character_id) DO UPDATE SET${ON_CONFLICT_SET} RETURNING ${SELECT_COLS}`,
            values: rows.flatMap(rowValues),
        };
    }

    onBulkDelete(itemKeys, accountId, _worldId) {
        return {
            text: `UPDATE character_overview SET deleted = TRUE, updated_at = NOW() WHERE character_id = ANY($1::integer[]) AND account_id = $2`,
            values: [itemKeys.map(Number), Number(accountId)],
        };
    }

    rowToModel(row) {
        return {
            characterId:   Number(row.character_id),
            accountId:     Number(row.account_id),
            worldId:       Number(row.world_id),
            name:          row.name,
            gender:        Number(row.gender),
            skinColor:     Number(row.skin_color),
            face:          Number(row.face),
            hair:          Number(row.hair),
            level:         Number(row.level),
            classId:       Number(row.class_id),
            mapId:         Number(row.map_id),
            spawnPoint:    Number(row.spawn_point),
            rank:          Number(row.rank),
            rankDiff:      Number(row.rank_diff),
            classRank:     Number(row.class_rank),
            classRankDiff: Number(row.class_rank_diff),
            baseLooks:     typeof row.base_looks === "string" ? JSON.parse(row.base_looks) : (row.base_looks ?? {}),
            overlays:      typeof row.overlays === "string" ? JSON.parse(row.overlays) : (row.overlays ?? {}),
        };
    }

    modelToRow(model) {
        return {
            character_id:   model.characterId,
            account_id:     model.accountId,
            world_id:       model.worldId,
            name:           model.name,
            gender:         model.gender,
            skin_color:     model.skinColor,
            face:           model.face,
            hair:           model.hair,
            level:          model.level,
            class_id:       model.classId,
            map_id:         model.mapId,
            spawn_point:    model.spawnPoint,
            rank:           model.rank ?? 0,
            rank_diff:      model.rankDiff ?? 0,
            class_rank:     model.classRank ?? 0,
            class_rank_diff: model.classRankDiff ?? 0,
            base_looks:     model.baseLooks ?? {},
            overlays:       model.overlays ?? {},
        };
    }

}

module.exports = { CharacterOverviewRepository };
