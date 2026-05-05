const { redisCacheKey } = require("../redis-cache-key");
import { ValueRepository } from "./value-repository";
import type { RepositoryQuery } from "../types/repository-contracts";
import type { PartyDeleteRow, PartyModel, PartyRow } from "../types/repository-models";

const SELECT_COLS = "world_id, party_id, leader_character_id, state, revision, disbanded_at, created_at, updated_at";

export type { PartyModel };

export class PartyRepository extends ValueRepository<PartyModel, PartyRow, number> {
    getKey(model: PartyModel) {
        return model.partyId;
    }

    getRedisKey(worldId: number, partyId: number) {
        return redisCacheKey(`w${worldId}:party:${partyId}`);
    }

    onSelect(partyId: number, worldId: number): RepositoryQuery {
        return {
            text: `SELECT ${SELECT_COLS} FROM parties WHERE world_id = $1 AND party_id = $2`,
            values: [Number(worldId), Number(partyId)],
        };
    }

    onSelectMany(partyIds: number[], worldId: number): RepositoryQuery {
        return {
            text: `SELECT ${SELECT_COLS} FROM parties WHERE world_id = $1 AND party_id = ANY($2::bigint[])`,
            values: [Number(worldId), partyIds.map(Number)],
        };
    }

    onUpsert(row: PartyRow): RepositoryQuery {
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
                Number(row.world_id),
                Number(row.party_id),
                Number(row.leader_character_id),
                row.state,
                Number(row.revision),
                row.disbanded_at ?? null,
            ],
        };
    }

    onDelete(row: PartyDeleteRow): RepositoryQuery {
        return {
            text: "DELETE FROM parties WHERE world_id = $1 AND party_id = $2",
            values: [Number(row.worldId), Number(row.partyId)],
        };
    }

    rowToModel(row: PartyRow): PartyModel {
        return {
            worldId: Number(row.world_id),
            partyId: Number(row.party_id),
            leaderCharacterId: Number(row.leader_character_id),
            state: row.state,
            revision: Number(row.revision),
            disbandedAt: row.disbanded_at ? new Date(row.disbanded_at) : null,
            createdAt: row.created_at instanceof Date ? row.created_at : row.created_at ? new Date(row.created_at) : undefined,
            updatedAt: row.updated_at instanceof Date ? row.updated_at : row.updated_at ? new Date(row.updated_at) : undefined,
        };
    }

    modelToRow(model: PartyModel): PartyRow {
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
