import { redisCacheKey } from "../redis-cache-key";
import { HashRepository } from "./hash-repository";
import type { RepositoryQuery, RepositoryTxOptions } from "../types/repository-contracts";
import type { CharacterBuddyModel, CharacterBuddyRow } from "../types/repository-models";

const SELECT_COLS = "character_id, buddy_character_id, group_name, pending, created_at, updated_at";

export type { CharacterBuddyModel };

export class CharacterBuddyRepository extends HashRepository<CharacterBuddyModel, CharacterBuddyRow> {
    override getGroupKey(model: CharacterBuddyModel) {
        return String(model.characterId);
    }

    override getItemKey(model: CharacterBuddyModel) {
        return String(model.buddyCharacterId);
    }

    override getRedisHashKey(worldId: number, characterId: string) {
        return redisCacheKey(`w${worldId}:character-buddy:${characterId}`);
    }

    override onSelect(characterId: string, _worldId: number): RepositoryQuery {
        return {
            text: `SELECT ${SELECT_COLS} FROM character_buddies WHERE character_id = $1 ORDER BY buddy_character_id ASC`,
            values: [Number(characterId)],
        };
    }

    override onBulkUpsert(rows: CharacterBuddyRow[]): RepositoryQuery {
        if (!rows.length) {
            return { text: "", values: [] };
        }
        const placeholders = rows
            .map((_, i) => `($${i * 4 + 1},$${i * 4 + 2},$${i * 4 + 3},$${i * 4 + 4},NOW(),NOW())`)
            .join(",\n");
        return {
            text: `INSERT INTO character_buddies (character_id, buddy_character_id, group_name, pending, created_at, updated_at)
                   VALUES
                   ${placeholders}
                   ON CONFLICT (character_id, buddy_character_id) DO UPDATE
                   SET group_name = EXCLUDED.group_name,
                       pending = EXCLUDED.pending,
                       updated_at = NOW()
                   RETURNING ${SELECT_COLS}`,
            values: rows.flatMap((r) => [r.character_id, r.buddy_character_id, r.group_name, r.pending]),
        };
    }

    override onBulkDelete(itemKeys: string[], characterId: string, _worldId: number): RepositoryQuery {
        return {
            text: "DELETE FROM character_buddies WHERE character_id = $1 AND buddy_character_id = ANY($2::int[])",
            values: [Number(characterId), itemKeys.map(Number)],
        };
    }

    override rowToModel(row: CharacterBuddyRow): CharacterBuddyModel {
        return {
            characterId: row.character_id,
            buddyCharacterId: row.buddy_character_id,
            groupName: row.group_name ?? "",
            pending: row.pending === true,
            createdAt: row.created_at instanceof Date ? row.created_at : row.created_at ? new Date(row.created_at) : undefined,
            updatedAt: row.updated_at instanceof Date ? row.updated_at : row.updated_at ? new Date(row.updated_at) : undefined,
        };
    }

    override modelToRow(model: CharacterBuddyModel): CharacterBuddyRow {
        return {
            character_id: model.characterId,
            buddy_character_id: model.buddyCharacterId,
            group_name: model.groupName,
            pending: model.pending,
        };
    }

    async listAcceptedBuddyCharacterIds(worldId: number, characterId: number, options: RepositoryTxOptions = {}): Promise<number[]> {
        const all = await this.getAll(worldId, String(characterId), options);
        const ids: number[] = [];
        for (const model of all.values()) {
            if (!model.pending) {
                ids.push(model.buddyCharacterId);
            }
        }
        ids.sort((a, b) => a - b);
        return ids;
    }

    async countAccepted(worldId: number, characterId: number, options: RepositoryTxOptions = {}): Promise<number> {
        const pool = this.pool(worldId, String(characterId));
        const res = await this.query(
            pool,
            "SELECT COUNT(*)::int AS cnt FROM character_buddies WHERE character_id = $1 AND NOT pending",
            [characterId],
            options
        );
        return Number(res.rows[0]?.cnt ?? 0);
    }

    async countAll(worldId: number, characterId: number, options: RepositoryTxOptions = {}): Promise<number> {
        const pool = this.pool(worldId, String(characterId));
        const res = await this.query(
            pool,
            "SELECT COUNT(*)::int AS cnt FROM character_buddies WHERE character_id = $1",
            [characterId],
            options
        );
        return Number(res.rows[0]?.cnt ?? 0);
    }
}
