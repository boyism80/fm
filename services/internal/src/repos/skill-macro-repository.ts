import { redisCacheKey } from "../redis-cache-key";
import { toPgInt } from "./pg-int";
import { HashRepository } from "./hash-repository";
import type { RepositoryQuery } from "../types/repository-contracts";
import type { SkillMacroModel, SkillMacroRow } from "../types/repository-models";

const SELECT_COLS = "character_id, position, name, shout, skill1, skill2, skill3, updated_at";
const INSERT_COLS = SELECT_COLS;
const PER_ROW_PARAMS = 7;

export type { SkillMacroModel };

function rowValues(row: SkillMacroRow) {
    return [row.character_id, row.position, row.name, row.shout, row.skill1, row.skill2, row.skill3];
}

export class SkillMacroRepository extends HashRepository<SkillMacroModel, SkillMacroRow> {
    override getTtlSeconds() {
        return this.ctx.appConfiguration.getCharacterCacheTtlSeconds();
    }

    override getGroupKey(model: SkillMacroModel) {
        return String(model.characterId);
    }

    override getItemKey(model: SkillMacroModel) {
        return String(model.position);
    }

    override getRedisHashKey(worldId: number, characterId: string) {
        return redisCacheKey(`w${worldId}:skill_macros:${characterId}`);
    }

    override onSelect(characterId: string): RepositoryQuery {
        return { text: `SELECT ${SELECT_COLS} FROM character_skill_macros WHERE character_id = $1 ORDER BY position`, values: [Number(characterId)] };
    }

    override onBulkUpsert(rows: SkillMacroRow[]): RepositoryQuery {
        if (!rows.length) {
            return { text: "", values: [] };
        }
        const placeholders = rows.map((_, i) => {
            const base = i * PER_ROW_PARAMS;
            return `($${base + 1},$${base + 2},$${base + 3},$${base + 4},$${base + 5},$${base + 6},$${base + 7},NOW())`;
        }).join(",");
        return {
            text: `INSERT INTO character_skill_macros (${INSERT_COLS}) VALUES\n${placeholders}\nON CONFLICT (character_id, position) DO UPDATE SET name = EXCLUDED.name, shout = EXCLUDED.shout, skill1 = EXCLUDED.skill1, skill2 = EXCLUDED.skill2, skill3 = EXCLUDED.skill3, updated_at = NOW() RETURNING ${SELECT_COLS}`,
            values: rows.flatMap(rowValues),
        };
    }

    override onBulkDelete(itemKeys: string[], characterId: string): RepositoryQuery {
        return { text: "DELETE FROM character_skill_macros WHERE position = ANY($1::smallint[]) AND character_id = $2", values: [itemKeys.map(Number), Number(characterId)] };
    }

    override rowToModel(row: SkillMacroRow): SkillMacroModel {
        return {
            characterId: toPgInt(row.character_id),
            position: toPgInt(row.position),
            name: row.name,
            shout: row.shout,
            skill1: toPgInt(row.skill1),
            skill2: toPgInt(row.skill2),
            skill3: toPgInt(row.skill3),
        };
    }

    override modelToRow(model: SkillMacroModel): SkillMacroRow {
        return {
            character_id: model.characterId,
            position: model.position,
            name: model.name,
            shout: model.shout,
            skill1: model.skill1,
            skill2: model.skill2,
            skill3: model.skill3,
        };
    }
}
