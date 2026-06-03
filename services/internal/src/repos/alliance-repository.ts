import { redisCacheKey } from "../redis-cache-key";
import { toPgInt } from "./pg-int";
import { ValueRepository } from "./value-repository";
import type { RepositoryQuery } from "../types/repository-contracts";
import type { AllianceModel, AllianceRow } from "../types/repository-models";
import {
    DEFAULT_ALLIANCE_CAPACITY,
    DEFAULT_ALLIANCE_RANK_TITLES,
    type AllianceRankTitles,
} from "../types/alliance-json";

const SELECT_COLS = `world_id, alliance_id, name, leader_character_id, guild_ids, rank_titles,
  capacity, notice, revision, disbanded_at, created_at, updated_at`;

export type { AllianceModel };

export class AllianceRepository extends ValueRepository<AllianceModel, AllianceRow, number> {
    override getKey(model: AllianceModel) {
        return model.allianceId;
    }

    override getRedisKey(worldId: number, allianceId: number) {
        return redisCacheKey(`w${worldId}:alliance:${allianceId}`);
    }

    guildIdsFromRow(guildIds: AllianceRow["guild_ids"]): number[] {
        if (!Array.isArray(guildIds)) {
            return [];
        }
        const out: number[] = [];
        for (const entry of guildIds) {
            if (entry == null) {
                continue;
            }
            const n = Number(entry);
            if (Number.isFinite(n) && n > 0) {
                out.push(n);
            }
        }
        return out;
    }

    rankTitlesFromRow(rankTitles: AllianceRow["rank_titles"]): AllianceRankTitles {
        if (!Array.isArray(rankTitles) || rankTitles.length !== 5) {
            return [...DEFAULT_ALLIANCE_RANK_TITLES];
        }
        return rankTitles.map((entry) => String(entry ?? "")) as AllianceRankTitles;
    }

    override onSelect(allianceId: number, worldId: number): RepositoryQuery {
        return {
            text: `SELECT ${SELECT_COLS} FROM alliances WHERE world_id = $1 AND alliance_id = $2 AND disbanded_at IS NULL`,
            values: [worldId, allianceId],
        };
    }

    override onUpsert(row: AllianceRow): RepositoryQuery {
        return {
            text: `INSERT INTO alliances (world_id, alliance_id, name, leader_character_id, guild_ids, rank_titles, capacity, notice, revision, disbanded_at, created_at, updated_at)
                   VALUES ($1, COALESCE($2, nextval('alliance_id_seq')), $3, $4, $5::bigint[], $6::text[], $7, $8, $9, $10, NOW(), NOW())
                   ON CONFLICT (world_id, alliance_id) DO UPDATE
                   SET name = EXCLUDED.name,
                       leader_character_id = EXCLUDED.leader_character_id,
                       guild_ids = EXCLUDED.guild_ids,
                       rank_titles = EXCLUDED.rank_titles,
                       capacity = EXCLUDED.capacity,
                       notice = EXCLUDED.notice,
                       revision = EXCLUDED.revision,
                       disbanded_at = EXCLUDED.disbanded_at,
                       updated_at = NOW()
                   RETURNING ${SELECT_COLS}`,
            values: [
                row.world_id,
                row.alliance_id > 0 ? row.alliance_id : null,
                row.name,
                row.leader_character_id,
                row.guild_ids,
                row.rank_titles,
                row.capacity,
                row.notice,
                row.revision,
                row.disbanded_at ?? null,
            ],
        };
    }

    override onDelete(row: { worldId: number; allianceId: number }): RepositoryQuery {
        return {
            text: "DELETE FROM alliances WHERE world_id = $1 AND alliance_id = $2",
            values: [row.worldId, row.allianceId],
        };
    }

    override normalizeRow(row: AllianceRow): AllianceRow {
        return {
            ...row,
            alliance_id: toPgInt(row.alliance_id),
            revision: toPgInt(row.revision),
        };
    }

    override rowToModel(row: AllianceRow): AllianceModel {
        return {
            worldId: row.world_id,
            allianceId: toPgInt(row.alliance_id),
            name: row.name,
            leaderCharacterId: row.leader_character_id,
            guildIds: this.guildIdsFromRow(row.guild_ids),
            rankTitles: this.rankTitlesFromRow(row.rank_titles),
            capacity: row.capacity,
            notice: row.notice,
            revision: toPgInt(row.revision),
            disbandedAt: row.disbanded_at ? new Date(row.disbanded_at) : null,
            createdAt: row.created_at instanceof Date ? row.created_at : row.created_at ? new Date(row.created_at) : undefined,
            updatedAt: row.updated_at instanceof Date ? row.updated_at : row.updated_at ? new Date(row.updated_at) : undefined,
        };
    }

    override modelToRow(model: AllianceModel): AllianceRow {
        return {
            world_id: model.worldId,
            alliance_id: model.allianceId,
            name: model.name,
            leader_character_id: model.leaderCharacterId,
            guild_ids: model.guildIds,
            rank_titles: model.rankTitles,
            capacity: model.capacity ?? DEFAULT_ALLIANCE_CAPACITY,
            notice: model.notice ?? "",
            revision: model.revision ?? 1,
            disbanded_at: model.disbandedAt ?? null,
        };
    }
}
