"use strict";

const { redisCacheKey } = require("../redis-cache-key");
const { ValueRepository } = require("./value-repository");

const SELECT_COLS = "world_id, party_id, leader_character_id, state, revision, disbanded_at, created_at, updated_at";

class PartyRepository extends ValueRepository {
    constructor(internalContext) {
        super(internalContext);
    }

    getKey(model) {
        return model.partyId;
    }

    getRedisKey(worldId, partyId) {
        return redisCacheKey(`w${worldId}:party:${partyId}`);
    }

    onSelect(partyId, worldId) {
        return {
            text: `SELECT ${SELECT_COLS} FROM parties WHERE world_id = $1 AND party_id = $2`,
            values: [Number(worldId), Number(partyId)],
        };
    }

    onSelectMany(partyIds, worldId) {
        return {
            text: `SELECT ${SELECT_COLS} FROM parties WHERE world_id = $1 AND party_id = ANY($2::bigint[])`,
            values: [Number(worldId), partyIds.map(Number)],
        };
    }

    onUpsert(row) {
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

    onDelete(row) {
        return {
            text: "DELETE FROM parties WHERE world_id = $1 AND party_id = $2",
            values: [Number(row.worldId), Number(row.partyId)],
        };
    }

    rowToModel(row) {
        return {
            worldId: Number(row.world_id),
            partyId: Number(row.party_id),
            leaderCharacterId: Number(row.leader_character_id),
            state: row.state,
            revision: Number(row.revision),
            disbandedAt: row.disbanded_at ? new Date(row.disbanded_at) : null,
            createdAt: row.created_at instanceof Date ? row.created_at : new Date(row.created_at),
            updatedAt: row.updated_at instanceof Date ? row.updated_at : new Date(row.updated_at),
        };
    }

    modelToRow(model) {
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

module.exports = { PartyRepository };
