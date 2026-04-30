"use strict";

const { redisCacheKey } = require("../redis-cache-key");
const { HashRepository } = require("./hash-repository");

const SELECT_COLS = `unique_id, owner_id, inventory_type, item_id, slot, count, expiration,
  enhance_chance, enhance_count, flag, skill_bonus, owner_name, equip_bonus_stats, created_at, updated_at`;

const INSERT_COLS = `unique_id, owner_id, inventory_type, item_id, slot, count, expiration,
  enhance_chance, enhance_count, flag, skill_bonus, owner_name, equip_bonus_stats, updated_at`;

const PER_ROW_PARAMS = 13;

function rowValues(row) {
    let uid = row.unique_id;
    if (uid === "" || uid === undefined) {
        uid = null;
    }
    return [
        uid,
        row.owner_id,
        row.inventory_type,
        row.item_id,
        row.slot,
        row.count,
        row.expiration,
        row.enhance_chance,
        row.enhance_count,
        row.flag,
        row.skill_bonus,
        row.owner_name,
        row.equip_bonus_stats,
    ];
}

class InventoryRepository extends HashRepository {
    constructor(internalContext) {
        super(internalContext);
    }

    getTtlSeconds() {
        return this.ctx.appConfiguration.getItemCacheTtlSeconds();
    }

    getGroupKey(model) {
        return model.ownerId;
    }

    getItemKey(model) {
        return `${model.inventoryType}:${model.slot}`;
    }

    getRedisHashKey(worldId, ownerId) {
        return redisCacheKey(`w${worldId}:inventory:${ownerId}`);
    }

    onSelect(ownerId, _worldId) {
        return {
            text: `SELECT ${SELECT_COLS} FROM inventory WHERE owner_id = $1`,
            values: [Number(ownerId)],
        };
    }

    onBulkUpsert(rows) {
        if (!rows.length) {
            return { text: "", values: [] };
        }
        const placeholders = rows
            .map((_, ri) => {
                const base = ri * PER_ROW_PARAMS + 1;
                return `(${Array.from({ length: PER_ROW_PARAMS }, (_, ci) => `$${base + ci}`).join(",")},NOW())`;
            })
            .join(",\n");
        return {
            text: `INSERT INTO inventory (${INSERT_COLS}) VALUES\n${placeholders}\nRETURNING ${SELECT_COLS}`,
            values: rows.flatMap(rowValues),
        };
    }

    onBulkDelete(itemKeys, ownerId, _worldId) {
        if (!itemKeys.length) {
            return { text: "", values: [] };
        }
        const tuples = itemKeys.map((k) => {
            const [t, s] = String(k).split(":");
            return [Number(t), Number(s)];
        });
        const values = [Number(ownerId), ...tuples.flat()];
        const conds = tuples
            .map((_, i) => `(inventory_type = $${2 + i * 2} AND slot = $${3 + i * 2})`)
            .join(" OR ");
        return {
            text: `DELETE FROM inventory WHERE owner_id = $1 AND (${conds})`,
            values,
        };
    }

    rowToModel(row) {
        return {
            uniqueId:
                row.unique_id != null && row.unique_id !== ""
                    ? Number(row.unique_id)
                    : null,
            ownerId: Number(row.owner_id),
            inventoryType: Number(row.inventory_type),
            itemId: Number(row.item_id),
            slot: Number(row.slot),
            count: Number(row.count),
            expiration: row.expiration ? new Date(row.expiration) : null,
            enhanceChance:
                row.enhance_chance != null ? Number(row.enhance_chance) : null,
            enhanceCount:
                row.enhance_count != null ? Number(row.enhance_count) : null,
            flag: row.flag != null ? Number(row.flag) : null,
            skillBonus:
                row.skill_bonus != null ? Number(row.skill_bonus) : null,
            ownerName: row.owner_name ?? null,
            equipBonusStats: JSON.parse(row.equip_bonus_stats),
            updatedAt:
                row.updated_at instanceof Date
                    ? row.updated_at
                    : new Date(row.updated_at),
        };
    }

    modelToRow(model) {
        let uniqueId = model.uniqueId;
        if (uniqueId === "" || uniqueId === 0 || uniqueId === "0") {
            uniqueId = null;
        }
        return {
            unique_id: uniqueId != null ? uniqueId : null,
            owner_id: model.ownerId,
            inventory_type: model.inventoryType,
            item_id: model.itemId,
            slot: model.slot,
            count: model.count,
            expiration: model.expiration ?? null,
            enhance_chance: model.enhanceChance ?? null,
            enhance_count: model.enhanceCount ?? 0,
            flag: model.flag ?? null,
            skill_bonus: model.skillBonus ?? null,
            owner_name: model.ownerName ?? null,
            equip_bonus_stats: JSON.stringify(model.equipBonusStats),
        };
    }

    async replaceBySnapshot(worldId, ownerId, models) {
        const pool = this._pool(worldId, ownerId);
        const normalized = models.map((m) => {
            const row = this.modelToRow(m);
            row.owner_id = Number(ownerId);
            return row;
        });

        await pool.query("BEGIN");
        try {
            await pool.query(`DELETE FROM inventory WHERE owner_id = $1`, [
                Number(ownerId),
            ]);
            if (normalized.length > 0) {
                const insert = this.onBulkUpsert(normalized);
                if (insert.text) {
                    await pool.query(insert.text, insert.values);
                }
            }
            await pool.query("COMMIT");
        } catch (err) {
            await pool.query("ROLLBACK");
            throw err;
        }

        const redis = this._redis(worldId, ownerId);
        await redis.del(this.getRedisHashKey(worldId, ownerId)).catch(() => {});
    }
}

module.exports = { InventoryRepository };
