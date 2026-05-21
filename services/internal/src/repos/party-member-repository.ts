import { redisCacheKey } from "../redis-cache-key";
import { toPgInt } from "./pg-int";
import { HashRepository } from "./hash-repository";
import type { RepositoryQuery } from "../types/repository-contracts";
import type { PartyMemberModel, PartyMemberRow } from "../types/repository-models";
import { PartyMemberRole } from "../protobuf/generated/fminternal/internal_service";

const DEFAULT_PARTY_MEMBER_ROLE = PartyMemberRole.PARTY_MEMBER_ROLE_MEMBER;

const SELECT_COLS = "world_id, party_id, character_id, character_name, level, class_id, role, map_id, door, joined_at, updated_at";

export type { PartyMemberModel };

export class PartyMemberRepository extends HashRepository<PartyMemberModel, PartyMemberRow> {
    override getGroupKey(model: PartyMemberModel) {
        return String(model.partyId);
    }

    override getItemKey(model: PartyMemberModel) {
        return String(model.characterId);
    }

    override getRedisHashKey(worldId: number, partyId: string) {
        return redisCacheKey(`w${worldId}:party-member:${partyId}`);
    }

    override onSelect(partyId: string, worldId: number): RepositoryQuery {
        return {
            text: `SELECT ${SELECT_COLS} FROM party_members WHERE world_id = $1 AND party_id = $2 ORDER BY joined_at ASC, character_id ASC`,
            values: [worldId, Number(partyId)],
        };
    }

    override onBulkUpsert(rows: PartyMemberRow[]): RepositoryQuery {
        if (!rows.length) {
            return { text: "", values: [] };
        }
        const placeholders = rows
            .map((_, i) => `($${i * 9 + 1},$${i * 9 + 2},$${i * 9 + 3},$${i * 9 + 4},$${i * 9 + 5},$${i * 9 + 6},$${i * 9 + 7},$${i * 9 + 8},$${i * 9 + 9},NOW(),NOW())`)
            .join(",\n");
        return {
            text: `INSERT INTO party_members (world_id, party_id, character_id, character_name, level, class_id, role, map_id, door, joined_at, updated_at)
                   VALUES
                   ${placeholders}
                   ON CONFLICT (world_id, party_id, character_id) DO UPDATE
                   SET character_name = EXCLUDED.character_name,
                       level = EXCLUDED.level,
                       class_id = EXCLUDED.class_id,
                       role = EXCLUDED.role,
                       map_id = EXCLUDED.map_id,
                       door = EXCLUDED.door,
                       updated_at = NOW()
                   RETURNING ${SELECT_COLS}`,
            values: rows.flatMap((r) => [
                r.world_id,
                r.party_id,
                r.character_id,
                r.character_name,
                r.level,
                r.class_id,
                r.role ?? DEFAULT_PARTY_MEMBER_ROLE,
                r.map_id ?? 0,
                r.door == null ? null : r.door,
            ]),
        };
    }

    override onBulkDelete(itemKeys: string[], partyId: string, worldId: number): RepositoryQuery {
        return {
            text: "DELETE FROM party_members WHERE world_id = $1 AND party_id = $2 AND character_id = ANY($3::int[])",
            values: [worldId, Number(partyId), itemKeys.map(Number)],
        };
    }

    doorFromRow(door: PartyMemberRow["door"]): PartyMemberModel["door"] {
        if (door == null || typeof door !== "object") {
            return null;
        }
        const town = Number(door.town);
        const target = Number(door.target);
        const x = Number(door.x);
        const y = Number(door.y);
        if (Number.isFinite(town) && Number.isFinite(target) && Number.isFinite(x) && Number.isFinite(y)) {
            return { town, target, x, y };
        }
        return null;
    }

    override normalizeRow(row: PartyMemberRow): PartyMemberRow {
        return {
            ...row,
            party_id: toPgInt(row.party_id),
        };
    }

    override rowToModel(row: PartyMemberRow): PartyMemberModel {
        return {
            worldId: row.world_id,
            partyId: toPgInt(row.party_id),
            characterId: row.character_id,
            characterName: row.character_name,
            level: row.level,
            classId: row.class_id,
            role: row.role ?? DEFAULT_PARTY_MEMBER_ROLE,
            mapId: row.map_id ?? 0,
            door: this.doorFromRow(row.door),
            joinedAt: row.joined_at instanceof Date ? row.joined_at : new Date(row.joined_at),
            updatedAt: row.updated_at instanceof Date ? row.updated_at : new Date(row.updated_at),
        };
    }

    override modelToRow(model: PartyMemberModel): PartyMemberRow {
        return {
            world_id: model.worldId,
            party_id: model.partyId,
            character_id: model.characterId,
            character_name: model.characterName,
            level: model.level,
            class_id: model.classId,
            role: model.role ?? DEFAULT_PARTY_MEMBER_ROLE,
            map_id: model.mapId ?? 0,
            door: model.door,
            joined_at: model.joinedAt ?? new Date(),
            updated_at: model.updatedAt ?? new Date(),
        };
    }
}
