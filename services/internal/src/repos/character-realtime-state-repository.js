"use strict";

const { ValueRepository } = require("./value-repository");

const SELECT_COLS = "world_id, character_id, party_id, guild_id, updated_at";

class CharacterRealtimeStateRepository extends ValueRepository {
    constructor(internalContext) {
        super(internalContext);
    }

    getKey(model) {
        return model.characterId;
    }

    getRedisKey(worldId, characterId) {
        const { keyPrefix } = this.ctx.getRedisDataAccess(worldId, characterId);
        return `${keyPrefix}fm:w${worldId}:character-realtime-state:${characterId}`;
    }

    onSelect(characterId, worldId) {
        return {
            text: `SELECT ${SELECT_COLS} FROM character_realtime_state WHERE world_id = $1 AND character_id = $2`,
            values: [Number(worldId), Number(characterId)],
        };
    }

    onSelectMany(characterIds, worldId) {
        return {
            text: `SELECT ${SELECT_COLS} FROM character_realtime_state WHERE world_id = $1 AND character_id = ANY($2::int[])`,
            values: [Number(worldId), characterIds.map(Number)],
        };
    }

    onUpsert(row) {
        return {
            text: `INSERT INTO character_realtime_state (world_id, character_id, party_id, guild_id, updated_at)
                   VALUES ($1,$2,$3,$4,NOW())
                   ON CONFLICT (world_id, character_id) DO UPDATE
                   SET party_id = EXCLUDED.party_id,
                       guild_id = EXCLUDED.guild_id,
                       updated_at = NOW()
                   RETURNING ${SELECT_COLS}`,
            values: [
                Number(row.world_id),
                Number(row.character_id),
                row.party_id == null ? null : Number(row.party_id),
                row.guild_id == null ? null : Number(row.guild_id),
            ],
        };
    }

    onDelete(row) {
        return {
            text: "DELETE FROM character_realtime_state WHERE world_id = $1 AND character_id = $2",
            values: [Number(row.worldId), Number(row.characterId)],
        };
    }

    rowToModel(row) {
        return {
            worldId: Number(row.world_id),
            characterId: Number(row.character_id),
            partyId: row.party_id == null ? null : Number(row.party_id),
            guildId: row.guild_id == null ? null : Number(row.guild_id),
            updatedAt: row.updated_at instanceof Date ? row.updated_at : new Date(row.updated_at),
        };
    }

    modelToRow(model) {
        return {
            world_id: model.worldId,
            character_id: model.characterId,
            party_id: model.partyId ?? null,
            guild_id: model.guildId ?? null,
        };
    }
}

module.exports = { CharacterRealtimeStateRepository };
