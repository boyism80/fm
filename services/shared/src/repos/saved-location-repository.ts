import { redisCacheKey } from "../redis-cache-key";
import { toPgInt } from "./pg-int";
import { HashRepository } from "./hash-repository";
import type { RepositoryQuery } from "../types/repository-contracts";
import type { SavedLocationModel, SavedLocationRow } from "../types/repository-models";

const SELECT_COLS = "character_id, location_key, map_id, updated_at";
const INSERT_COLS = "character_id, location_key, map_id, updated_at";
const ON_CONFLICT_SET = `
  map_id = EXCLUDED.map_id,
  updated_at = NOW()`;
const PER_ROW_PARAMS = 3;

export type { SavedLocationModel };

function rowValues(row: SavedLocationRow) {
    return [row.character_id, row.location_key, row.map_id];
}

export class SavedLocationRepository extends HashRepository<SavedLocationModel, SavedLocationRow> {
    override getTtlSeconds() {
        return this.ctx.appConfiguration.getCharacterCacheTtlSeconds();
    }

    override getGroupKey(model: SavedLocationModel) {
        return String(model.characterId);
    }

    override getItemKey(model: SavedLocationModel) {
        return model.name;
    }

    override getRedisHashKey(worldId: number, characterId: string) {
        return redisCacheKey(`w${worldId}:saved_locations:${characterId}`);
    }

    override onSelect(characterId: string): RepositoryQuery {
        return {
            text: `SELECT ${SELECT_COLS} FROM character_saved_locations WHERE character_id = $1`,
            values: [Number(characterId)],
        };
    }

    override onBulkUpsert(rows: SavedLocationRow[]): RepositoryQuery {
        if (!rows.length) {
            return { text: "", values: [] };
        }
        const placeholders = rows
            .map((_, i) => {
                const base = i * PER_ROW_PARAMS;
                return `($${base + 1},$${base + 2},$${base + 3},NOW())`;
            })
            .join(",");
        return {
            text: `INSERT INTO character_saved_locations (${INSERT_COLS}) VALUES\n${placeholders}\nON CONFLICT (character_id, location_key) DO UPDATE SET${ON_CONFLICT_SET} RETURNING ${SELECT_COLS}`,
            values: rows.flatMap(rowValues),
        };
    }

    override onBulkDelete(itemKeys: string[], characterId: string): RepositoryQuery {
        return {
            text: "DELETE FROM character_saved_locations WHERE location_key = ANY($1::text[]) AND character_id = $2",
            values: [itemKeys, Number(characterId)],
        };
    }

    override normalizeRow(row: SavedLocationRow): SavedLocationRow {
        return {
            ...row,
            map_id: toPgInt(row.map_id),
        };
    }

    override rowToModel(row: SavedLocationRow): SavedLocationModel {
        return {
            characterId: row.character_id,
            name: row.location_key,
            mapId: toPgInt(row.map_id),
            updatedAt: row.updated_at instanceof Date ? row.updated_at : row.updated_at ? new Date(row.updated_at) : undefined,
        };
    }

    override modelToRow(model: SavedLocationModel): SavedLocationRow {
        return {
            character_id: model.characterId,
            location_key: model.name,
            map_id: model.mapId,
        };
    }
}
