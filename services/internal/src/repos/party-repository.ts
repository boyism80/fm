import { redisCacheKey } from "../redis-cache-key";
import { toPgInt } from "./pg-int";
import { ValueRepository } from "./value-repository";
import type { RepositoryQuery } from "../types/repository-contracts";
import type { PartyDeleteRow, PartyModel, PartyRow } from "../types/repository-models";

const SELECT_COLS = "world_id, party_id, leader_character_id, state, revision, disbanded_at, created_at, updated_at";

export type { PartyModel };

export class PartyRepository extends ValueRepository<PartyModel, PartyRow, number> {
    override getKey(model: PartyModel) {
        return model.partyId;
    }

    override getRedisKey(worldId: number, partyId: number) {
        return redisCacheKey(`w${worldId}:party:${partyId}`);
    }

    override onSelect(partyId: number, worldId: number): RepositoryQuery {
        return {
            text: `SELECT ${SELECT_COLS} FROM parties WHERE world_id = $1 AND party_id = $2`,
            values: [worldId, partyId],
        };
    }

    override onSelectMany(partyIds: number[], worldId: number): RepositoryQuery {
        return {
            text: `SELECT ${SELECT_COLS} FROM parties WHERE world_id = $1 AND party_id = ANY($2::bigint[])`,
            values: [worldId, partyIds],
        };
    }

    override onUpsert(row: PartyRow): RepositoryQuery {
        return {
            text: `INSERT INTO parties (world_id, party_id, leader_character_id, state, revision, disbanded_at, created_at, updated_at)
                   VALUES ($1,$2,$3,$4,$5,$6,NOW(),NOW())
                   ON CONFLICT (world_id, party_id) DO UPDATE
                   SET leader_character_id = EXCLUDED.leader_character_id,
                       state = EXCLUDED.state,
                       revision = EXCLUDED.revision,
                       disbanded_at = EXCLUDED.disbanded_at,
                       updated_at = NOW()
                   RETURNING ${SELECT_COLS}`,
            values: [
                row.world_id,
                row.party_id,
                row.leader_character_id,
                row.state,
                row.revision,
                row.disbanded_at ?? null,
            ],
        };
    }

    override onDelete(row: PartyDeleteRow): RepositoryQuery {
        return {
            text: "DELETE FROM parties WHERE world_id = $1 AND party_id = $2",
            values: [row.worldId, row.partyId],
        };
    }

    override normalizeRow(row: PartyRow): PartyRow {
        return {
            ...row,
            party_id: toPgInt(row.party_id),
            revision: toPgInt(row.revision),
        };
    }

    override rowToModel(row: PartyRow): PartyModel {
        return {
            worldId: row.world_id,
            partyId: toPgInt(row.party_id),
            leaderCharacterId: row.leader_character_id,
            state: row.state,
            revision: toPgInt(row.revision),
            disbandedAt: row.disbanded_at ? new Date(row.disbanded_at) : null,
            createdAt: row.created_at instanceof Date ? row.created_at : row.created_at ? new Date(row.created_at) : undefined,
            updatedAt: row.updated_at instanceof Date ? row.updated_at : row.updated_at ? new Date(row.updated_at) : undefined,
        };
    }

    override modelToRow(model: PartyModel): PartyRow {
        return {
            world_id: model.worldId,
            party_id: model.partyId,
            leader_character_id: model.leaderCharacterId,
            state: model.state ?? "ACTIVE",
            revision: model.revision ?? 1,
            disbanded_at: model.disbandedAt ?? null,
        };
    }
}
