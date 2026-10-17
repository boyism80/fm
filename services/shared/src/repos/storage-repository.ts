import { redisCacheKey } from "../redis-cache-key";
import { ValueRepository } from "./value-repository";
import type { RepositoryQuery } from "../types/repository-contracts";
import type { StorageDeleteModel, StorageModel, StorageRow } from "../types/repository-models";

const SELECT_COLS = "account_id, world_id, slots, meso, updated_at";
const INSERT_COLS = "account_id, world_id, slots, meso, updated_at";
const ON_CONFLICT_SET = "slots = EXCLUDED.slots, meso = EXCLUDED.meso, updated_at = NOW()";

export type { StorageModel };

export class StorageRepository extends ValueRepository<StorageModel, StorageRow, number> {
    override getKey(model: StorageModel) {
        return model.accountId;
    }

    override getTtlSeconds() {
        return this.ctx.appConfiguration.getItemCacheTtlSeconds();
    }

    override getRedisKey(worldId: number, accountId: number) {
        return redisCacheKey(`w${worldId}:storage:${accountId}`);
    }

    override onSelect(accountId: number, worldId: number): RepositoryQuery {
        return {
            text: `SELECT ${SELECT_COLS} FROM storages WHERE account_id = $1 AND world_id = $2`,
            values: [accountId, worldId],
        };
    }

    override onUpsert(row: StorageRow): RepositoryQuery {
        return {
            text: `INSERT INTO storages (${INSERT_COLS}) VALUES ($1,$2,$3,$4,NOW()) ON CONFLICT (account_id, world_id) DO UPDATE SET ${ON_CONFLICT_SET} RETURNING ${SELECT_COLS}`,
            values: [row.account_id, row.world_id, row.slots, row.meso],
        };
    }

    override onDelete(row: StorageDeleteModel): RepositoryQuery {
        return {
            text: "DELETE FROM storages WHERE account_id = $1 AND world_id = $2",
            values: [row.accountId, row.worldId],
        };
    }

    override rowToModel(row: StorageRow): StorageModel {
        return {
            accountId: row.account_id,
            worldId: row.world_id,
            slots: row.slots,
            meso: row.meso,
            updatedAt: row.updated_at instanceof Date ? row.updated_at : row.updated_at ? new Date(row.updated_at) : undefined,
        };
    }

    override modelToRow(model: StorageModel): StorageRow {
        return {
            account_id: model.accountId,
            world_id: model.worldId,
            slots: model.slots,
            meso: model.meso,
        };
    }
}
