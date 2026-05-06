import { redisCacheKey } from "../redis-cache-key";
import { ValueRepository } from "./value-repository";
import type { RepositoryQuery } from "../types/repository-contracts";
import type { KeyLayoutDeleteModel, KeyLayoutModel, KeyLayoutRow } from "../types/repository-models";

const SELECT_COLS = "character_id, world_id, key_layout_json, updated_at";
const INSERT_COLS = "character_id, world_id, key_layout_json, updated_at";
const ON_CONFLICT_SET = "key_layout_json = EXCLUDED.key_layout_json, updated_at = NOW()";

export type { KeyLayoutModel };

export class KeyLayoutRepository extends ValueRepository<KeyLayoutModel, KeyLayoutRow, number> {
    override getKey(model: KeyLayoutModel) {
        return model.characterId;
    }

    override getTtlSeconds() {
        return this.ctx.appConfiguration.getCharacterCacheTtlSeconds();
    }

    override getRedisKey(worldId: number, characterId: number) {
        return redisCacheKey(`w${worldId}:keylayout:${characterId}`);
    }

    override onSelect(characterId: number, worldId: number): RepositoryQuery {
        return {
            text: `SELECT ${SELECT_COLS} FROM key_layout WHERE character_id = $1 AND world_id = $2`,
            values: [characterId, worldId],
        };
    }

    override onSelectMany(characterIds: number[], worldId: number): RepositoryQuery {
        return {
            text: `SELECT ${SELECT_COLS} FROM key_layout WHERE character_id = ANY($1::int[]) AND world_id = $2`,
            values: [characterIds, worldId],
        };
    }

    override onUpsert(row: KeyLayoutRow): RepositoryQuery {
        return {
            text: `INSERT INTO key_layout (${INSERT_COLS}) VALUES ($1,$2,$3,NOW()) ON CONFLICT (character_id, world_id) DO UPDATE SET ${ON_CONFLICT_SET} RETURNING ${SELECT_COLS}`,
            values: [row.character_id, row.world_id, row.key_layout_json ?? "{}"],
        };
    }

    override onBulkUpsert(rows: KeyLayoutRow[]): RepositoryQuery | null {
        if (!rows.length) {
            return null;
        }
        const placeholders = rows.map((_, i) => `($${i * 3 + 1},$${i * 3 + 2},$${i * 3 + 3},NOW())`).join(",\n");
        return {
            text: `INSERT INTO key_layout (${INSERT_COLS}) VALUES\n${placeholders}\nON CONFLICT (character_id, world_id) DO UPDATE SET ${ON_CONFLICT_SET} RETURNING ${SELECT_COLS}`,
            values: rows.flatMap((r) => [r.character_id, r.world_id, r.key_layout_json ?? "{}"]),
        };
    }

    override onDelete(row: KeyLayoutDeleteModel): RepositoryQuery {
        return {
            text: "DELETE FROM key_layout WHERE character_id = $1 AND world_id = $2",
            values: [row.characterId, row.worldId],
        };
    }

    override rowToModel(row: KeyLayoutRow): KeyLayoutModel {
        return {
            characterId: row.character_id,
            worldId: row.world_id,
            keyLayoutJson: typeof row.key_layout_json === "string" ? row.key_layout_json : JSON.stringify(row.key_layout_json ?? {}),
            updatedAt: row.updated_at instanceof Date ? row.updated_at : row.updated_at ? new Date(row.updated_at) : undefined,
        };
    }

    override modelToRow(model: KeyLayoutModel): KeyLayoutRow {
        return {
            character_id: model.characterId,
            world_id: model.worldId,
            key_layout_json: model.keyLayoutJson ?? "{}",
        };
    }
}
