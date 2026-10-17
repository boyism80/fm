import { redisCacheKey } from "../redis-cache-key";
import { toPgInt } from "./pg-int";
import { HashRepository } from "./hash-repository";
import type { RepositoryQuery } from "../types/repository-contracts";
import type { MonsterBookCardModel, MonsterBookCardRow } from "../types/repository-models";

const SELECT_COLS = "character_id, card_id, count, updated_at";
const INSERT_COLS = "character_id, card_id, count, updated_at";
const ON_CONFLICT_SET = `
  count = EXCLUDED.count,
  updated_at = NOW()`;
const PER_ROW_PARAMS = 3;

export type { MonsterBookCardModel };

function rowValues(row: MonsterBookCardRow) {
    return [row.character_id, row.card_id, row.count];
}

export class MonsterBookRepository extends HashRepository<MonsterBookCardModel, MonsterBookCardRow> {
    override getTtlSeconds() {
        return this.ctx.appConfiguration.getCharacterCacheTtlSeconds();
    }

    override getGroupKey(model: MonsterBookCardModel) {
        return String(model.characterId);
    }

    override getItemKey(model: MonsterBookCardModel) {
        return String(model.cardId);
    }

    override getRedisHashKey(worldId: number, characterId: string) {
        return redisCacheKey(`w${worldId}:monster_book:${characterId}`);
    }

    override onSelect(characterId: string): RepositoryQuery {
        return {
            text: `SELECT ${SELECT_COLS} FROM character_monster_book WHERE character_id = $1`,
            values: [Number(characterId)],
        };
    }

    override onBulkUpsert(rows: MonsterBookCardRow[]): RepositoryQuery {
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
            text: `INSERT INTO character_monster_book (${INSERT_COLS}) VALUES\n${placeholders}\nON CONFLICT (character_id, card_id) DO UPDATE SET${ON_CONFLICT_SET} RETURNING ${SELECT_COLS}`,
            values: rows.flatMap(rowValues),
        };
    }

    override onBulkDelete(itemKeys: string[], characterId: string): RepositoryQuery {
        return {
            text: "DELETE FROM character_monster_book WHERE card_id = ANY($1::int[]) AND character_id = $2",
            values: [itemKeys.map(Number), Number(characterId)],
        };
    }

    override normalizeRow(row: MonsterBookCardRow): MonsterBookCardRow {
        return {
            ...row,
            card_id: toPgInt(row.card_id),
            count: toPgInt(row.count),
        };
    }

    override rowToModel(row: MonsterBookCardRow): MonsterBookCardModel {
        return {
            characterId: row.character_id,
            cardId: toPgInt(row.card_id),
            count: toPgInt(row.count),
        };
    }

    override modelToRow(model: MonsterBookCardModel): MonsterBookCardRow {
        return {
            character_id: model.characterId,
            card_id: model.cardId,
            count: model.count,
        };
    }
}
