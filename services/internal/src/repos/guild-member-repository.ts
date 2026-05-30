import { redisCacheKey } from "../redis-cache-key";
import { toPgInt } from "./pg-int";
import { HashRepository } from "./hash-repository";
import type { RepositoryQuery, RepositoryTxOptions } from "../types/repository-contracts";
import type { GuildMemberModel, GuildMemberRow } from "../types/repository-models";
import { GuildMemberRank } from "../protobuf/generated/fminternal/internal_service";

const DEFAULT_GUILD_MEMBER_RANK = GuildMemberRank.GUILD_MEMBER_RANK_MEMBER;

const SELECT_COLS = "world_id, guild_id, character_id, character_name, level, class_id, guild_rank, joined_at, updated_at";

export type { GuildMemberModel };

export class GuildMemberRepository extends HashRepository<GuildMemberModel, GuildMemberRow> {
    override getGroupKey(model: GuildMemberModel) {
        return String(model.guildId);
    }

    override getItemKey(model: GuildMemberModel) {
        return String(model.characterId);
    }

    override getRedisHashKey(worldId: number, guildId: string) {
        return redisCacheKey(`w${worldId}:guild-member:${guildId}`);
    }

    override onSelect(guildId: string, worldId: number): RepositoryQuery {
        return {
            text: `SELECT ${SELECT_COLS} FROM guild_members WHERE world_id = $1 AND guild_id = $2 ORDER BY guild_rank ASC, joined_at ASC, character_id ASC`,
            values: [worldId, Number(guildId)],
        };
    }

    override onBulkUpsert(rows: GuildMemberRow[]): RepositoryQuery {
        if (!rows.length) {
            return { text: "", values: [] };
        }
        const placeholders = rows
            .map((_, i) => `($${i * 7 + 1},$${i * 7 + 2},$${i * 7 + 3},$${i * 7 + 4},$${i * 7 + 5},$${i * 7 + 6},$${i * 7 + 7},NOW(),NOW())`)
            .join(",\n");
        return {
            text: `INSERT INTO guild_members (world_id, guild_id, character_id, character_name, level, class_id, guild_rank, joined_at, updated_at)
                   VALUES
                   ${placeholders}
                   ON CONFLICT (world_id, guild_id, character_id) DO UPDATE
                   SET character_name = EXCLUDED.character_name,
                       level = EXCLUDED.level,
                       class_id = EXCLUDED.class_id,
                       guild_rank = EXCLUDED.guild_rank,
                       updated_at = NOW()
                   RETURNING ${SELECT_COLS}`,
            values: rows.flatMap((r) => [
                r.world_id,
                r.guild_id,
                r.character_id,
                r.character_name,
                r.level,
                r.class_id,
                r.guild_rank ?? DEFAULT_GUILD_MEMBER_RANK,
            ]),
        };
    }

    override onBulkDelete(itemKeys: string[], guildId: string, worldId: number): RepositoryQuery {
        return {
            text: "DELETE FROM guild_members WHERE world_id = $1 AND guild_id = $2 AND character_id = ANY($3::int[])",
            values: [worldId, Number(guildId), itemKeys.map(Number)],
        };
    }

    async deleteAllForGuild(worldId: number, guildId: number, options: RepositoryTxOptions = {}) {
        const groupKey = String(guildId);
        const pool = this.pool(worldId, groupKey);
        await this.query(
            pool,
            "DELETE FROM guild_members WHERE world_id = $1 AND guild_id = $2",
            [worldId, guildId],
            options
        );
        if (!options.txClient) {
            await this.evictGroupCache(worldId, groupKey);
        }
    }

    override normalizeRow(row: GuildMemberRow): GuildMemberRow {
        return {
            ...row,
            guild_id: toPgInt(row.guild_id),
        };
    }

    override rowToModel(row: GuildMemberRow): GuildMemberModel {
        return {
            worldId: row.world_id,
            guildId: toPgInt(row.guild_id),
            characterId: row.character_id,
            characterName: row.character_name,
            level: row.level,
            classId: row.class_id,
            guildRank: (row.guild_rank ?? DEFAULT_GUILD_MEMBER_RANK) as GuildMemberRank,
            joinedAt: row.joined_at instanceof Date ? row.joined_at : new Date(row.joined_at),
            updatedAt: row.updated_at instanceof Date ? row.updated_at : new Date(row.updated_at),
        };
    }

    override modelToRow(model: GuildMemberModel): GuildMemberRow {
        return {
            world_id: model.worldId,
            guild_id: model.guildId,
            character_id: model.characterId,
            character_name: model.characterName,
            level: model.level,
            class_id: model.classId,
            guild_rank: model.guildRank ?? DEFAULT_GUILD_MEMBER_RANK,
            joined_at: model.joinedAt ?? new Date(),
            updated_at: model.updatedAt ?? new Date(),
        };
    }
}
