import { redisCacheKey } from "../redis-cache-key";
import { toPgInt, toPgIntOrNull } from "./pg-int";
import { ValueRepository } from "./value-repository";
import type { RepositoryQuery } from "../types/repository-contracts";
import type { GuildDeleteRow, GuildModel, GuildRow } from "../types/repository-models";
import {
    DEFAULT_GUILD_LOGO,
    DEFAULT_GUILD_RANK_TITLES,
    type GuildLogo,
    type GuildRankTitles,
} from "../types/guild-json";

const SELECT_COLS = `world_id, guild_id, name, leader_character_id, gp, capacity, notice, logo, rank_titles,
  alliance_id, revision, disbanded_at, created_at, updated_at`;

export type { GuildModel };

export class GuildRepository extends ValueRepository<GuildModel, GuildRow, number> {
    override getKey(model: GuildModel) {
        return model.guildId;
    }

    override getRedisKey(worldId: number, guildId: number) {
        return redisCacheKey(`w${worldId}:guild:${guildId}`);
    }

    logoFromRow(logo: GuildRow["logo"]): GuildLogo {
        if (logo == null || typeof logo !== "object") {
            return { ...DEFAULT_GUILD_LOGO };
        }
        const logoId = Number(logo.logo);
        const logoColor = Number(logo.logoColor);
        const logoBG = Number(logo.logoBG);
        const logoBGColor = Number(logo.logoBGColor);
        if (
            Number.isFinite(logoId) &&
            Number.isFinite(logoColor) &&
            Number.isFinite(logoBG) &&
            Number.isFinite(logoBGColor)
        ) {
            return {
                logo: logoId,
                logoColor,
                logoBG,
                logoBGColor,
            };
        }
        return { ...DEFAULT_GUILD_LOGO };
    }

    rankTitlesFromRow(rankTitles: GuildRow["rank_titles"]): GuildRankTitles {
        if (!Array.isArray(rankTitles) || rankTitles.length !== 5) {
            return [...DEFAULT_GUILD_RANK_TITLES];
        }
        return rankTitles.map((entry) => String(entry ?? "")) as GuildRankTitles;
    }

    override onSelect(guildId: number, worldId: number): RepositoryQuery {
        return {
            text: `SELECT ${SELECT_COLS} FROM guilds WHERE world_id = $1 AND guild_id = $2 AND disbanded_at IS NULL`,
            values: [worldId, guildId],
        };
    }

    override onSelectMany(guildIds: number[], worldId: number): RepositoryQuery {
        return {
            text: `SELECT ${SELECT_COLS} FROM guilds WHERE world_id = $1 AND guild_id = ANY($2::bigint[]) AND disbanded_at IS NULL`,
            values: [worldId, guildIds],
        };
    }

    override onUpsert(row: GuildRow): RepositoryQuery {
        return {
            text: `INSERT INTO guilds (world_id, guild_id, name, leader_character_id, gp, capacity, notice, logo, rank_titles, alliance_id, revision, disbanded_at, created_at, updated_at)
                   VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9::text[],$10,$11,$12,NOW(),NOW())
                   ON CONFLICT (world_id, guild_id) DO UPDATE
                   SET name = EXCLUDED.name,
                       leader_character_id = EXCLUDED.leader_character_id,
                       gp = EXCLUDED.gp,
                       capacity = EXCLUDED.capacity,
                       notice = EXCLUDED.notice,
                       logo = EXCLUDED.logo,
                       rank_titles = EXCLUDED.rank_titles,
                       alliance_id = EXCLUDED.alliance_id,
                       revision = EXCLUDED.revision,
                       disbanded_at = EXCLUDED.disbanded_at,
                       updated_at = NOW()
                   RETURNING ${SELECT_COLS}`,
            values: [
                row.world_id,
                row.guild_id,
                row.name,
                row.leader_character_id,
                row.gp,
                row.capacity,
                row.notice,
                row.logo,
                row.rank_titles,
                row.alliance_id ?? null,
                row.revision,
                row.disbanded_at ?? null,
            ],
        };
    }

    override onDelete(row: GuildDeleteRow): RepositoryQuery {
        return {
            text: "DELETE FROM guilds WHERE world_id = $1 AND guild_id = $2",
            values: [row.worldId, row.guildId],
        };
    }

    override normalizeRow(row: GuildRow): GuildRow {
        return {
            ...row,
            guild_id: toPgInt(row.guild_id),
            alliance_id: toPgIntOrNull(row.alliance_id),
            revision: toPgInt(row.revision),
        };
    }

    override rowToModel(row: GuildRow): GuildModel {
        return {
            worldId: row.world_id,
            guildId: toPgInt(row.guild_id),
            name: row.name,
            leaderCharacterId: row.leader_character_id,
            gp: row.gp,
            capacity: row.capacity,
            notice: row.notice,
            logo: this.logoFromRow(row.logo),
            rankTitles: this.rankTitlesFromRow(row.rank_titles),
            allianceId: toPgIntOrNull(row.alliance_id),
            revision: toPgInt(row.revision),
            disbandedAt: row.disbanded_at ? new Date(row.disbanded_at) : null,
            createdAt: row.created_at instanceof Date ? row.created_at : row.created_at ? new Date(row.created_at) : undefined,
            updatedAt: row.updated_at instanceof Date ? row.updated_at : row.updated_at ? new Date(row.updated_at) : undefined,
        };
    }

    override modelToRow(model: GuildModel): GuildRow {
        return {
            world_id: model.worldId,
            guild_id: model.guildId,
            name: model.name,
            leader_character_id: model.leaderCharacterId,
            gp: model.gp,
            capacity: model.capacity,
            notice: model.notice,
            logo: model.logo,
            rank_titles: model.rankTitles,
            alliance_id: model.allianceId ?? null,
            revision: model.revision ?? 1,
            disbanded_at: model.disbandedAt ?? null,
        };
    }
}
