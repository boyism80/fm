import { redisCacheKey } from "../redis-cache-key";
import { toPgInt } from "./pg-int";
import { HashRepository } from "./hash-repository";
import type { RepositoryQuery } from "../types/repository-contracts";
import type { UnsettledRequestModel, UnsettledRequestRow } from "../types/repository-models";

const SELECT_COLS = "character_id, id, type, target, detail, updated_at";
const PER_ROW_PARAMS = 5;

export type { UnsettledRequestModel };

function rowValues(row: UnsettledRequestRow) {
    return [row.character_id, row.id, row.type, row.target, row.detail];
}

export class UnsettledRequestRepository extends HashRepository<UnsettledRequestModel, UnsettledRequestRow> {
    override getTtlSeconds() {
        return this.ctx.appConfiguration.getCharacterCacheTtlSeconds();
    }

    override getGroupKey(model: UnsettledRequestModel) {
        return String(model.characterId);
    }

    override getItemKey(model: UnsettledRequestModel) {
        return model.id;
    }

    override getRedisHashKey(worldId: number, characterId: string) {
        return redisCacheKey(`w${worldId}:unsettled_requests:${characterId}`);
    }

    override onSelect(characterId: string): RepositoryQuery {
        return { text: `SELECT ${SELECT_COLS} FROM character_unsettled_requests WHERE character_id = $1 ORDER BY id`, values: [Number(characterId)] };
    }

    override onBulkUpsert(rows: UnsettledRequestRow[]): RepositoryQuery {
        if (!rows.length) {
            return { text: "", values: [] };
        }
        const placeholders = rows.map((_, i) => {
            const base = i * PER_ROW_PARAMS;
            return `($${base + 1},$${base + 2},$${base + 3},$${base + 4},$${base + 5},NOW())`;
        }).join(",");
        return {
            text: `INSERT INTO character_unsettled_requests (${SELECT_COLS}) VALUES\n${placeholders}\nON CONFLICT (character_id, id) DO UPDATE SET type = EXCLUDED.type, target = EXCLUDED.target, detail = EXCLUDED.detail, updated_at = NOW() RETURNING ${SELECT_COLS}`,
            values: rows.flatMap(rowValues),
        };
    }

    override onBulkDelete(itemKeys: string[], characterId: string): RepositoryQuery {
        return { text: "DELETE FROM character_unsettled_requests WHERE id = ANY($1::bigint[]) AND character_id = $2", values: [itemKeys, Number(characterId)] };
    }

    override rowToModel(row: UnsettledRequestRow): UnsettledRequestModel {
        return {
            characterId: toPgInt(row.character_id),
            id: String(row.id),
            type: toPgInt(row.type),
            target: String(row.target),
            detail: row.detail,
        };
    }

    override modelToRow(model: UnsettledRequestModel): UnsettledRequestRow {
        return {
            character_id: model.characterId,
            id: model.id,
            type: model.type,
            target: model.target,
            detail: model.detail,
        };
    }
}
