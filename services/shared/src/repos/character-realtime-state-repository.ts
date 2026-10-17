import { redisCacheKey } from "../redis-cache-key";
import { toPgIntOrNull } from "./pg-int";
import { ValueRepository } from "./value-repository";
import type { RepositoryQuery, RepositoryTxOptions } from "../types/repository-contracts";
import { DEFAULT_BUDDY_CAPACITY } from "../buddy-defaults";
import type {
    CharacterRealtimeStateDeleteRow,
    CharacterRealtimeStateModel,
    CharacterRealtimeStateRow,
} from "../types/repository-models";
const SELECT_COLS = "world_id, character_id, party_id, guild_id, buddy_capacity, updated_at";

export type { CharacterRealtimeStateModel };

export class CharacterRealtimeStateRepository extends ValueRepository<CharacterRealtimeStateModel, CharacterRealtimeStateRow, number> {
    override getKey(model: CharacterRealtimeStateModel) {
        return model.characterId;
    }

    override getRedisKey(worldId: number, characterId: number) {
        return redisCacheKey(`w${worldId}:character-realtime-state:${characterId}`);
    }

    override onSelect(characterId: number, worldId: number): RepositoryQuery {
        return {
            text: `SELECT ${SELECT_COLS} FROM character_realtime_state WHERE world_id = $1 AND character_id = $2`,
            values: [worldId, characterId],
        };
    }

    override onSelectMany(characterIds: number[], worldId: number): RepositoryQuery {
        return {
            text: `SELECT ${SELECT_COLS} FROM character_realtime_state WHERE world_id = $1 AND character_id = ANY($2::int[])`,
            values: [worldId, characterIds],
        };
    }

    override onUpsert(row: CharacterRealtimeStateRow): RepositoryQuery {
        return {
            text: `INSERT INTO character_realtime_state (world_id, character_id, party_id, guild_id, buddy_capacity, updated_at)
                   VALUES ($1,$2,$3,$4,COALESCE($5, ${DEFAULT_BUDDY_CAPACITY}),NOW())
                   ON CONFLICT (world_id, character_id) DO UPDATE
                   SET party_id = EXCLUDED.party_id,
                       guild_id = EXCLUDED.guild_id,
                       buddy_capacity = COALESCE(EXCLUDED.buddy_capacity, character_realtime_state.buddy_capacity),
                       updated_at = NOW()
                   RETURNING ${SELECT_COLS}`,
            values: [
                row.world_id,
                row.character_id,
                row.party_id,
                row.guild_id,
                row.buddy_capacity ?? null,
            ],
        };
    }

    async getBuddyCapacity(worldId: number, characterId: number, options: RepositoryTxOptions = {}): Promise<number> {
        const state = await this.get(worldId, characterId, options);
        if (state?.buddyCapacity == null) {
            return DEFAULT_BUDDY_CAPACITY;
        }
        return state.buddyCapacity;
    }

    override onDelete(row: CharacterRealtimeStateDeleteRow): RepositoryQuery {
        return {
            text: "DELETE FROM character_realtime_state WHERE world_id = $1 AND character_id = $2",
            values: [row.worldId, row.characterId],
        };
    }

    override normalizeRow(row: CharacterRealtimeStateRow): CharacterRealtimeStateRow {
        return {
            ...row,
            party_id: toPgIntOrNull(row.party_id),
            guild_id: toPgIntOrNull(row.guild_id),
        };
    }

    override rowToModel(row: CharacterRealtimeStateRow): CharacterRealtimeStateModel {
        return {
            worldId: row.world_id,
            characterId: row.character_id,
            partyId: toPgIntOrNull(row.party_id),
            guildId: toPgIntOrNull(row.guild_id),
            buddyCapacity: row.buddy_capacity == null ? undefined : Number(row.buddy_capacity),
            updatedAt: row.updated_at instanceof Date ? row.updated_at : row.updated_at ? new Date(row.updated_at) : undefined,
        };
    }

    override modelToRow(model: CharacterRealtimeStateModel): CharacterRealtimeStateRow {
        return {
            world_id: model.worldId,
            character_id: model.characterId,
            party_id: model.partyId ?? null,
            guild_id: model.guildId ?? null,
            buddy_capacity: model.buddyCapacity ?? null,
        };
    }
}
