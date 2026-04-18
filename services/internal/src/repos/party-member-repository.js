"use strict";

const { redisCacheKey } = require("../redis-cache-key");
const { HashRepository } = require("./hash-repository");

const SELECT_COLS =
    "world_id, party_id, character_id, character_name, level, class_id, role, map_id, channel_index, door, joined_at, updated_at";

class PartyMemberRepository extends HashRepository {
    constructor(internalContext) {
        super(internalContext);
    }

    getGroupKey(model) {
        return model.partyId;
    }

    getItemKey(model) {
        return model.characterId;
    }

    getRedisHashKey(worldId, partyId) {
        return redisCacheKey(`w${worldId}:party-member:${partyId}`);
    }

    onSelect(partyId, worldId) {
        return {
            text: `SELECT ${SELECT_COLS} FROM party_members WHERE world_id = $1 AND party_id = $2`,
            values: [Number(worldId), Number(partyId)],
        };
    }

    onBulkUpsert(rows) {
        if (!rows.length) {
            return null;
        }
        const placeholders = rows
            .map(
                (_, i) =>
                    `($${i * 10 + 1},$${i * 10 + 2},$${i * 10 + 3},$${i * 10 + 4},$${i * 10 + 5},$${i * 10 + 6},$${i * 10 + 7},$${i * 10 + 8},$${i * 10 + 9},$${i * 10 + 10},NOW(),NOW())`
            )
            .join(",\n");
        return {
            text: `INSERT INTO party_members (world_id, party_id, character_id, character_name, level, class_id, role, map_id, channel_index, door, joined_at, updated_at)
                   VALUES
                   ${placeholders}
                   ON CONFLICT (world_id, party_id, character_id) DO UPDATE
                   SET character_name = EXCLUDED.character_name,
                       level = EXCLUDED.level,
                       class_id = EXCLUDED.class_id,
                       role = EXCLUDED.role,
                       map_id = EXCLUDED.map_id,
                       channel_index = EXCLUDED.channel_index,
                       door = EXCLUDED.door,
                       updated_at = NOW()
                   RETURNING ${SELECT_COLS}`,
            values: rows.flatMap((r) => [
                Number(r.world_id),
                Number(r.party_id),
                Number(r.character_id),
                String(r.character_name),
                Number(r.level),
                Number(r.class_id),
                r.role ?? "MEMBER",
                Number(r.map_id ?? 0),
                Number(r.channel_index ?? -2),
                r.door == null ? null : r.door,
            ]),
        };
    }

    onBulkDelete(itemKeys, partyId, worldId) {
        return {
            text: "DELETE FROM party_members WHERE world_id = $1 AND party_id = $2 AND character_id = ANY($3::int[])",
            values: [Number(worldId), Number(partyId), itemKeys.map(Number)],
        };
    }

    _doorFromRow(door) {
        if (door == null) {
            return null;
        }
        if (typeof door === "string") {
            try {
                return JSON.parse(door);
            } catch {
                return null;
            }
        }
        return door;
    }

    rowToModel(row) {
        return {
            worldId: Number(row.world_id),
            partyId: Number(row.party_id),
            characterId: Number(row.character_id),
            characterName: row.character_name,
            level: Number(row.level),
            classId: Number(row.class_id),
            role: row.role,
            mapId: Number(row.map_id ?? 0),
            channelIndex: Number(row.channel_index ?? -2),
            door: this._doorFromRow(row.door),
            joinedAt: row.joined_at instanceof Date ? row.joined_at : new Date(row.joined_at),
            updatedAt: row.updated_at instanceof Date ? row.updated_at : new Date(row.updated_at),
        };
    }

    modelToRow(model) {
        return {
            world_id: model.worldId,
            party_id: model.partyId,
            character_id: model.characterId,
            character_name: model.characterName,
            level: model.level,
            class_id: model.classId,
            role: model.role ?? "MEMBER",
            map_id: model.mapId ?? 0,
            channel_index: model.channelIndex ?? -2,
            door: model.door == null ? null : model.door,
        };
    }
}

module.exports = { PartyMemberRepository };
