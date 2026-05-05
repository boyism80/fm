const { redisCacheKey } = require("../redis-cache-key");
import { ValueRepository } from "./value-repository";
import type { RepositoryQuery } from "../types/repository-contracts";
import type {
    CharacterRealtimeStateDeleteRow,
    CharacterRealtimeStateModel,
    CharacterRealtimeStateRow,
} from "../types/repository-models";

const SELECT_COLS = "world_id, character_id, party_id, guild_id, updated_at";

export type { CharacterRealtimeStateModel };

export class CharacterRealtimeStateRepository extends ValueRepository<CharacterRealtimeStateModel, CharacterRealtimeStateRow, number> {
    getKey(model: CharacterRealtimeStateModel) {
        return model.characterId;
    }

    getRedisKey(worldId: number, characterId: number) {
        return redisCacheKey(`w${worldId}:character-realtime-state:${characterId}`);
    }

    onSelect(characterId: number, worldId: number): RepositoryQuery {
        return {
            text: `SELECT ${SELECT_COLS} FROM character_realtime_state WHERE world_id = $1 AND character_id = $2`,
            values: [Number(worldId), Number(characterId)],
        };
    }

    onSelectMany(characterIds: number[], worldId: number): RepositoryQuery {
        return {
            text: `SELECT ${SELECT_COLS} FROM character_realtime_state WHERE world_id = $1 AND character_id = ANY($2::int[])`,
            values: [Number(worldId), characterIds.map(Number)],
        };
    }

    onUpsert(row: CharacterRealtimeStateRow): RepositoryQuery {
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

    onDelete(row: CharacterRealtimeStateDeleteRow): RepositoryQuery {
        return {
            text: "DELETE FROM character_realtime_state WHERE world_id = $1 AND character_id = $2",
            values: [Number(row.worldId), Number(row.characterId)],
        };
    }

    rowToModel(row: CharacterRealtimeStateRow): CharacterRealtimeStateModel {
        return {
            worldId: Number(row.world_id),
            characterId: Number(row.character_id),
            partyId: row.party_id == null ? null : Number(row.party_id),
            guildId: row.guild_id == null ? null : Number(row.guild_id),
            updatedAt: row.updated_at instanceof Date ? row.updated_at : row.updated_at ? new Date(row.updated_at) : undefined,
        };
    }

    modelToRow(model: CharacterRealtimeStateModel): CharacterRealtimeStateRow {
        return {
            world_id: model.worldId,
            character_id: model.characterId,
            party_id: model.partyId ?? null,
            guild_id: model.guildId ?? null,
        };
    }
}
