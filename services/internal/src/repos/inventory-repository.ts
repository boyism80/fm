import { redisCacheKey } from "../redis-cache-key";
import { toPgIntOrNull } from "./pg-int";
import { HashRepository } from "./hash-repository";
import type { RepositoryQuery } from "../types/repository-contracts";
import type { InventoryModel, InventoryRow } from "../types/repository-models";

const SELECT_COLS = `unique_id, owner_id, inventory_type, item_id, slot, count, expiration,
  enhance_chance, enhance_count, flag, skill_bonus, owner_name, equip_bonus_stats, created_at, updated_at`;

const INSERT_COLS = `unique_id, owner_id, inventory_type, item_id, slot, count, expiration,
  enhance_chance, enhance_count, flag, skill_bonus, owner_name, equip_bonus_stats, updated_at`;

const PER_ROW_PARAMS = 13;

export type { InventoryModel };

function rowValues(row: InventoryRow) {
    let uid = row.unique_id;
    if (uid === undefined) {
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

export class InventoryRepository extends HashRepository<InventoryModel, InventoryRow> {
    override getTtlSeconds() {
        return this.ctx.appConfiguration.getItemCacheTtlSeconds();
    }

    override getGroupKey(model: InventoryModel) {
        return String(model.ownerId);
    }

    override getItemKey(model: InventoryModel) {
        return `${model.inventoryType}:${model.slot}`;
    }

    override getRedisHashKey(worldId: number, ownerId: string) {
        return redisCacheKey(`w${worldId}:inventory:${ownerId}`);
    }

    override onSelect(ownerId: string): RepositoryQuery {
        return {
            text: `SELECT ${SELECT_COLS} FROM inventory WHERE owner_id = $1`,
            values: [Number(ownerId)],
        };
    }

    override onBulkUpsert(rows: InventoryRow[]): RepositoryQuery {
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

    override onBulkDelete(itemKeys: string[], ownerId: string): RepositoryQuery {
        if (!itemKeys.length) {
            return { text: "", values: [] };
        }
        const tuples = itemKeys.map((k) => {
            const [t, s] = k.split(":");
            return [Number(t), Number(s)];
        });
        const values = [Number(ownerId), ...tuples.flat()];
        const conds = tuples.map((_, i) => `(inventory_type = $${2 + i * 2} AND slot = $${3 + i * 2})`).join(" OR ");
        return {
            text: `DELETE FROM inventory WHERE owner_id = $1 AND (${conds})`,
            values,
        };
    }

    override normalizeRow(row: InventoryRow): InventoryRow {
        return {
            ...row,
            unique_id: toPgIntOrNull(row.unique_id),
        };
    }

    override rowToModel(row: InventoryRow): InventoryModel {
        return {
            uniqueId: toPgIntOrNull(row.unique_id),
            ownerId: row.owner_id,
            inventoryType: row.inventory_type,
            itemId: row.item_id,
            slot: row.slot,
            count: row.count,
            expiration: row.expiration ? new Date(row.expiration) : null,
            enhanceChance: row.enhance_chance,
            enhanceCount: row.enhance_count,
            flag: row.flag,
            skillBonus: row.skill_bonus,
            ownerName: row.owner_name ?? null,
            equipBonusStats: typeof row.equip_bonus_stats === "string" ? JSON.parse(row.equip_bonus_stats) : (row.equip_bonus_stats ?? {}),
            updatedAt: row.updated_at instanceof Date ? row.updated_at : row.updated_at ? new Date(row.updated_at) : undefined,
        };
    }

    override modelToRow(model: InventoryModel): InventoryRow {
        let uniqueId = model.uniqueId;
        if (uniqueId === 0) {
            uniqueId = null;
        }
        return {
            unique_id: uniqueId ?? null,
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
            equip_bonus_stats: model.equipBonusStats ?? {},
        };
    }

    async replaceBySnapshot(worldId: number, ownerId: number, models: InventoryModel[]) {
        const groupKey = String(ownerId);
        const pool = this.pool(worldId, groupKey);
        const normalized = models.map((m) => {
            const row = this.modelToRow(m);
            row.owner_id = ownerId;
            return row;
        });

        await pool.query("BEGIN");
        try {
            await pool.query("DELETE FROM inventory WHERE owner_id = $1", [ownerId]);
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

        const redis = this.redis(worldId, groupKey);
        await redis.del(this.getRedisHashKey(worldId, String(ownerId))).catch(() => {});
    }
}
