import { redisCacheKey } from "../redis-cache-key";
import { ValueRepository } from "./value-repository";
import type { RepositoryQuery } from "../types/repository-contracts";
import type { MountDeleteModel, MountModel, MountRow } from "../types/repository-models";

const SELECT_COLS = "character_id, world_id, level, exp, fatigue, updated_at";
const INSERT_COLS = "character_id, world_id, level, exp, fatigue, updated_at";
const ON_CONFLICT_SET = "level = EXCLUDED.level, exp = EXCLUDED.exp, fatigue = EXCLUDED.fatigue, updated_at = NOW()";

export type { MountModel };

export class MountRepository extends ValueRepository<MountModel, MountRow, number> {
    override getKey(model: MountModel) {
        return model.characterId;
    }

    override getTtlSeconds() {
        return this.ctx.appConfiguration.getCharacterCacheTtlSeconds();
    }

    override getRedisKey(worldId: number, characterId: number) {
        return redisCacheKey(`w${worldId}:mount:${characterId}`);
    }

    override onSelect(characterId: number, worldId: number): RepositoryQuery {
        return {
            text: `SELECT ${SELECT_COLS} FROM mount WHERE character_id = $1 AND world_id = $2`,
            values: [characterId, worldId],
        };
    }

    override onSelectMany(characterIds: number[], worldId: number): RepositoryQuery {
        return {
            text: `SELECT ${SELECT_COLS} FROM mount WHERE character_id = ANY($1::int[]) AND world_id = $2`,
            values: [characterIds, worldId],
        };
    }

    override onUpsert(row: MountRow): RepositoryQuery {
        return {
            text: `INSERT INTO mount (${INSERT_COLS}) VALUES ($1,$2,$3,$4,$5,NOW()) ON CONFLICT (character_id, world_id) DO UPDATE SET ${ON_CONFLICT_SET} RETURNING ${SELECT_COLS}`,
            values: [row.character_id, row.world_id, row.level, row.exp, row.fatigue],
        };
    }

    override onBulkUpsert(rows: MountRow[]): RepositoryQuery | null {
        if (!rows.length) {
            return null;
        }
        const placeholders = rows.map((_, i) => `($${i * 5 + 1},$${i * 5 + 2},$${i * 5 + 3},$${i * 5 + 4},$${i * 5 + 5},NOW())`).join(",\n");
        return {
            text: `INSERT INTO mount (${INSERT_COLS}) VALUES\n${placeholders}\nON CONFLICT (character_id, world_id) DO UPDATE SET ${ON_CONFLICT_SET} RETURNING ${SELECT_COLS}`,
            values: rows.flatMap((r) => [r.character_id, r.world_id, r.level, r.exp, r.fatigue]),
        };
    }

    override onDelete(row: MountDeleteModel): RepositoryQuery {
        return {
            text: "DELETE FROM mount WHERE character_id = $1 AND world_id = $2",
            values: [row.characterId, row.worldId],
        };
    }

    override rowToModel(row: MountRow): MountModel {
        return {
            characterId: row.character_id,
            worldId: row.world_id,
            level: row.level,
            exp: row.exp,
            fatigue: row.fatigue,
            updatedAt: row.updated_at instanceof Date ? row.updated_at : row.updated_at ? new Date(row.updated_at) : undefined,
        };
    }

    override modelToRow(model: MountModel): MountRow {
        return {
            character_id: model.characterId,
            world_id: model.worldId,
            level: model.level,
            exp: model.exp,
            fatigue: model.fatigue,
        };
    }
}
