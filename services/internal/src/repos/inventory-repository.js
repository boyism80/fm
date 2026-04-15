"use strict";

const { HashRepository } = require("./hash-repository");

const SELECT_COLS = `unique_id, owner_id, item_id, slot, count, expiration,
  enchant_chance, flag, skill_bonus, owner_name, deleted, created_at, updated_at`;

const INSERT_COLS = `unique_id, owner_id, item_id, slot, count, expiration,
  enchant_chance, flag, skill_bonus, owner_name, deleted, updated_at`;

const ON_CONFLICT_SET = `
  owner_id = EXCLUDED.owner_id, item_id = EXCLUDED.item_id, slot = EXCLUDED.slot,
  count = EXCLUDED.count, expiration = EXCLUDED.expiration,
  enchant_chance = EXCLUDED.enchant_chance, flag = EXCLUDED.flag,
  skill_bonus = EXCLUDED.skill_bonus, owner_name = EXCLUDED.owner_name,
  deleted = FALSE, updated_at = NOW()`;

const PER_ROW_PARAMS = 11;

function rowValues(row) {
    return [
        row.unique_id, row.owner_id,  row.item_id,       row.slot,
        row.count,     row.expiration, row.enchant_chance, row.flag,
        row.skill_bonus, row.owner_name, false,
    ];
}

class InventoryRepository extends HashRepository {
    constructor(internalContext) {
        super(internalContext);
    }

    getTtlSeconds() {
        return this.ctx.appConfiguration.getItemCacheTtlSeconds();
    }

    getOwnerKey(model) {
        return model.ownerId;
    }

    getItemKey(model) {
        return String(model.uniqueId);
    }

    getRedisHashKey(worldId, ownerId) {
        const { keyPrefix } = this.ctx.getRedisDataAccess(worldId, ownerId);
        return `${keyPrefix}fm:w${worldId}:inventory:${ownerId}`;
    }

    onSelectByOwner(ownerId, _worldId) {
        return {
            text: `SELECT ${SELECT_COLS} FROM inventory WHERE owner_id = $1 AND NOT deleted`,
            values: [Number(ownerId)],
        };
    }

    onBulkUpsert(rows) {
        if (!rows.length) return { text: "", values: [] };
        const placeholders = rows
            .map((_, ri) =>
                `(${Array.from({ length: PER_ROW_PARAMS }, (_, ci) => `$${ri * PER_ROW_PARAMS + ci + 1}`).join(",")},NOW())`
            )
            .join(",\n");
        return {
            text: `INSERT INTO inventory (${INSERT_COLS}) VALUES\n${placeholders}\nON CONFLICT (unique_id) DO UPDATE SET${ON_CONFLICT_SET} RETURNING ${SELECT_COLS}`,
            values: rows.flatMap(rowValues),
        };
    }

    onBulkDelete(itemKeys, ownerId, _worldId) {
        return {
            text: `UPDATE inventory SET deleted = TRUE, updated_at = NOW() WHERE unique_id = ANY($1::bigint[]) AND owner_id = $2`,
            values: [itemKeys.map(String), Number(ownerId)],
        };
    }

    rowToModel(row) {
        return {
            uniqueId:      String(row.unique_id),
            ownerId:       Number(row.owner_id),
            itemId:        Number(row.item_id),
            slot:          Number(row.slot),
            count:         Number(row.count),
            expiration:    row.expiration ? new Date(row.expiration) : null,
            enchantChance: row.enchant_chance != null ? Number(row.enchant_chance) : null,
            flag:          row.flag != null ? Number(row.flag) : null,
            skillBonus:    row.skill_bonus != null ? Number(row.skill_bonus) : null,
            ownerName:     row.owner_name ?? null,
            updatedAt:     row.updated_at instanceof Date ? row.updated_at : new Date(row.updated_at),
        };
    }

    modelToRow(model) {
        return {
            unique_id:      String(model.uniqueId),
            owner_id:       model.ownerId,
            item_id:        model.itemId,
            slot:           model.slot,
            count:          model.count,
            expiration:     model.expiration ?? null,
            enchant_chance: model.enchantChance ?? null,
            flag:           model.flag ?? null,
            skill_bonus:    model.skillBonus ?? null,
            owner_name:     model.ownerName ?? null,
        };
    }

}

module.exports = { InventoryRepository };
