import { redisCacheKey } from "../redis-cache-key";
import { ValueRepository } from "./value-repository";
import type { RepositoryQuery } from "../types/repository-contracts";
import type { CharacterDeleteModel, CharacterModel, CharacterRow } from "../types/repository-models";

const SELECT_COLS = `id, account_id, world_id, name, gender, skin_color, face, hair, level, class_id, role,
  str, dex, int_stat, luk, hp, max_hp, mp, max_mp, ability_point, exp,
  map_id, spawn_point, pos_x, pos_y, stance, meso, skill_point, hidden, deleted, created_at, updated_at`;

const ON_CONFLICT_SET = `
  account_id = EXCLUDED.account_id, world_id = EXCLUDED.world_id, name = EXCLUDED.name, gender = EXCLUDED.gender,
  skin_color = EXCLUDED.skin_color, face = EXCLUDED.face, hair = EXCLUDED.hair,
  level = EXCLUDED.level, class_id = EXCLUDED.class_id, role = EXCLUDED.role, str = EXCLUDED.str,
  dex = EXCLUDED.dex, int_stat = EXCLUDED.int_stat, luk = EXCLUDED.luk,
  hp = EXCLUDED.hp, max_hp = EXCLUDED.max_hp, mp = EXCLUDED.mp, max_mp = EXCLUDED.max_mp,
  ability_point = EXCLUDED.ability_point, exp = EXCLUDED.exp, map_id = EXCLUDED.map_id,
  spawn_point = EXCLUDED.spawn_point, pos_x = EXCLUDED.pos_x, pos_y = EXCLUDED.pos_y,
  stance = EXCLUDED.stance, meso = EXCLUDED.meso, skill_point = EXCLUDED.skill_point, hidden = EXCLUDED.hidden,
  deleted = FALSE, updated_at = NOW()`;

const INSERT_COLS = `id, account_id, world_id, name, gender, skin_color, face, hair, level, class_id, role,
  str, dex, int_stat, luk, hp, max_hp, mp, max_mp, ability_point, exp,
  map_id, spawn_point, pos_x, pos_y, stance, meso, skill_point, hidden, deleted, updated_at`;

const PER_ROW_PARAMS = 29;

export type { CharacterModel };

function rowValues(row: CharacterRow) {
    return [
        row.id, row.account_id, row.world_id, row.name,
        row.gender, row.skin_color, row.face, row.hair,
        row.level, row.class_id, row.role, row.str,
        row.dex, row.int_stat, row.luk, row.hp,
        row.max_hp, row.mp, row.max_mp, row.ability_point,
        row.exp, row.map_id, row.spawn_point, row.pos_x,
        row.pos_y, row.stance, row.meso, row.skill_point,
        row.hidden,
    ];
}

export class CharacterRepository extends ValueRepository<CharacterModel, CharacterRow, number> {
    getKey(model: CharacterModel) {
        return model.characterId;
    }

    getTtlSeconds() {
        return this.ctx.appConfiguration.getCharacterCacheTtlSeconds();
    }

    getRedisKey(worldId: number, characterId: number) {
        return redisCacheKey(`w${worldId}:character:${characterId}`);
    }

    onSelect(characterId: number, worldId: number): RepositoryQuery {
        return {
            text: `SELECT ${SELECT_COLS} FROM characters WHERE id = $1 AND world_id = $2 AND NOT deleted`,
            values: [Number(characterId), Number(worldId)],
        };
    }

    onSelectMany(characterIds: number[], worldId: number): RepositoryQuery {
        return {
            text: `SELECT ${SELECT_COLS} FROM characters WHERE id = ANY($1::int[]) AND world_id = $2 AND NOT deleted`,
            values: [characterIds.map(Number), Number(worldId)],
        };
    }

    onUpsert(row: CharacterRow): RepositoryQuery {
        return {
            text: `INSERT INTO characters (${INSERT_COLS}) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26,$27,$28,$29,FALSE,NOW()) ON CONFLICT (id) DO UPDATE SET${ON_CONFLICT_SET} RETURNING ${SELECT_COLS}`,
            values: rowValues(row),
        };
    }

    onBulkUpsert(rows: CharacterRow[]): RepositoryQuery | null {
        if (!rows.length) {
            return null;
        }
        const placeholders = rows
            .map((_, ri) => `(${Array.from({ length: PER_ROW_PARAMS }, (_, ci) => `$${ri * PER_ROW_PARAMS + ci + 1}`).join(",")},FALSE,NOW())`)
            .join(",\n");
        return {
            text: `INSERT INTO characters (${INSERT_COLS}) VALUES\n${placeholders}\nON CONFLICT (id) DO UPDATE SET${ON_CONFLICT_SET} RETURNING ${SELECT_COLS}`,
            values: rows.flatMap(rowValues),
        };
    }

    onDelete(row: CharacterDeleteModel): RepositoryQuery {
        return {
            text: "UPDATE characters SET deleted = true, updated_at = NOW() WHERE id = $1 AND account_id = $2 AND world_id = $3 AND deleted = false",
            values: [Number(row.characterId), Number(row.accountId), Number(row.worldId)],
        };
    }

    rowToModel(row: CharacterRow): CharacterModel {
        return {
            characterId: Number(row.id),
            accountId: Number(row.account_id),
            worldId: Number(row.world_id),
            name: row.name,
            gender: Number(row.gender),
            skinColor: Number(row.skin_color),
            face: Number(row.face),
            hair: Number(row.hair),
            level: Number(row.level),
            classId: Number(row.class_id),
            role: Number(row.role),
            str: Number(row.str),
            dex: Number(row.dex),
            intStat: Number(row.int_stat),
            luk: Number(row.luk),
            hp: Number(row.hp),
            maxHp: Number(row.max_hp),
            mp: Number(row.mp),
            maxMp: Number(row.max_mp),
            abilityPoint: Number(row.ability_point),
            exp: Number(row.exp),
            mapId: Number(row.map_id),
            spawnPoint: Number(row.spawn_point),
            positionX: Number(row.pos_x),
            positionY: Number(row.pos_y),
            stance: Number(row.stance),
            meso: Number(row.meso),
            skillPoint: Number(row.skill_point),
            hidden: row.hidden,
            updatedAt: row.updated_at instanceof Date ? row.updated_at : row.updated_at ? new Date(row.updated_at) : undefined,
        };
    }

    modelToRow(model: CharacterModel): CharacterRow {
        return {
            id: model.characterId,
            account_id: model.accountId,
            world_id: model.worldId,
            name: model.name,
            gender: model.gender,
            skin_color: model.skinColor,
            face: model.face,
            hair: model.hair,
            level: model.level,
            class_id: model.classId,
            role: model.role ?? 0,
            str: model.str ?? 0,
            dex: model.dex ?? 0,
            int_stat: model.intStat ?? 0,
            luk: model.luk ?? 0,
            hp: model.hp ?? 0,
            max_hp: model.maxHp ?? 0,
            mp: model.mp ?? 0,
            max_mp: model.maxMp ?? 0,
            ability_point: model.abilityPoint ?? 0,
            exp: model.exp ?? 0,
            map_id: model.mapId,
            spawn_point: model.spawnPoint,
            pos_x: model.positionX ?? 0,
            pos_y: model.positionY ?? 0,
            stance: model.stance ?? 0,
            meso: model.meso ?? 0,
            skill_point: model.skillPoint ?? 0,
            hidden: model.hidden ?? false,
        };
    }
}
