import { redisCacheKey } from "../redis-cache-key";
import { toPgInt } from "./pg-int";
import { HashRepository } from "./hash-repository";
import type { RepositoryQuery } from "../types/repository-contracts";
import type { RecordModel, RecordRow } from "../types/repository-models";

const PER_ROW_PARAMS = 6;

export type { RecordModel };

function rowValues(row: RecordRow) {
    return [row.owner_id, row.record_key, row.value, row.text, row.period, row.updated_at_unix_ms];
}

export class CharacterRecordRepository extends HashRepository<RecordModel, RecordRow> {
    protected get table() {
        return "character_records";
    }

    protected get ownerColumn() {
        return "character_id";
    }

    private get selectCols() {
        return `${this.ownerColumn} AS owner_id, record_key, value, text, period, updated_at_unix_ms, updated_at`;
    }

    override getTtlSeconds() {
        return this.ctx.appConfiguration.getCharacterCacheTtlSeconds();
    }

    override getGroupKey(model: RecordModel) {
        return String(model.ownerId);
    }

    override getItemKey(model: RecordModel) {
        return model.key;
    }

    override getRedisHashKey(worldId: number, ownerId: string) {
        return redisCacheKey(`w${worldId}:${this.table}:${ownerId}`);
    }

    override onSelect(ownerId: string): RepositoryQuery {
        return {
            text: `SELECT ${this.selectCols} FROM ${this.table} WHERE ${this.ownerColumn} = $1`,
            values: [Number(ownerId)],
        };
    }

    override onBulkUpsert(rows: RecordRow[]): RepositoryQuery {
        if (!rows.length) {
            return { text: "", values: [] };
        }
        const placeholders = rows
            .map((_, i) => {
                const base = i * PER_ROW_PARAMS;
                return `($${base + 1},$${base + 2},$${base + 3},$${base + 4},$${base + 5},$${base + 6},NOW())`;
            })
            .join(",");
        return {
            text: `INSERT INTO ${this.table} (${this.ownerColumn}, record_key, value, text, period, updated_at_unix_ms, updated_at) VALUES
${placeholders}
ON CONFLICT (${this.ownerColumn}, record_key) DO UPDATE SET
  value = EXCLUDED.value,
  text = EXCLUDED.text,
  period = EXCLUDED.period,
  updated_at_unix_ms = EXCLUDED.updated_at_unix_ms,
  updated_at = NOW() RETURNING ${this.selectCols}`,
            values: rows.flatMap(rowValues),
        };
    }

    override onBulkDelete(itemKeys: string[], ownerId: string): RepositoryQuery {
        return {
            text: `DELETE FROM ${this.table} WHERE record_key = ANY($1::text[]) AND ${this.ownerColumn} = $2`,
            values: [itemKeys, Number(ownerId)],
        };
    }

    override normalizeRow(row: RecordRow): RecordRow {
        return {
            ...row,
            owner_id: toPgInt(row.owner_id),
            value: toPgInt(row.value),
            period: toPgInt(row.period),
            updated_at_unix_ms: toPgInt(row.updated_at_unix_ms),
        };
    }

    override rowToModel(row: RecordRow): RecordModel {
        return {
            ownerId: toPgInt(row.owner_id),
            key: row.record_key,
            value: toPgInt(row.value),
            text: row.text ?? "",
            period: toPgInt(row.period),
            updatedAtUnixMs: toPgInt(row.updated_at_unix_ms),
            updatedAt: row.updated_at instanceof Date ? row.updated_at : row.updated_at ? new Date(row.updated_at) : undefined,
        };
    }

    override modelToRow(model: RecordModel): RecordRow {
        return {
            owner_id: model.ownerId,
            record_key: model.key,
            value: model.value,
            text: model.text ?? "",
            period: model.period,
            updated_at_unix_ms: model.updatedAtUnixMs,
        };
    }
}
