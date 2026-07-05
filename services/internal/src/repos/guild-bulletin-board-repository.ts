import type { Pool, PoolClient } from "pg";
import type { InternalContext } from "../context/internal-context";
import type { RepositoryTxOptions } from "../types/repository-contracts";
import type {
    GuildBulletinBoardReplyModel,
    GuildBulletinBoardReplyRow,
    GuildBulletinBoardThreadModel,
    GuildBulletinBoardThreadRow,
    GuildBulletinBoardThreadWithRepliesModel,
} from "../types/repository-models";
import { toPgInt } from "./pg-int";

const NOTICE_LOCAL_THREAD_ID = 0;

export type CreateGuildBulletinBoardThreadInput = {
    guildId: number;
    posterCharacterId: number;
    title: string;
    body: string;
    icon: number;
    isNotice: boolean;
    createdAt?: Date;
};

export type CreateGuildBulletinBoardReplyInput = {
    guildId: number;
    localThreadId: number;
    posterCharacterId: number;
    content: string;
    createdAt?: Date;
};

type GuildBulletinBoardThreadWithRepliesRow = GuildBulletinBoardThreadRow & {
    reply_id: number | null;
    reply_poster_character_id: number | null;
    reply_content: string | null;
    reply_created_at: Date | string | null;
};

export class GuildBulletinBoardRepository {
    private readonly ctx: InternalContext;

    constructor(internalContext: InternalContext) {
        this.ctx = internalContext;
    }

    getShardHash(guildId: number): number {
        return guildId;
    }

    pool(worldId: number, guildId: number): Pool {
        return this.ctx.getPgDataPool(worldId, this.getShardHash(guildId));
    }

    private async runInDataTransaction<T>(
        worldId: number,
        guildId: number,
        options: RepositoryTxOptions,
        fn: (client: PoolClient) => Promise<T>
    ): Promise<T> {
        if (options.txClient) {
            return fn(options.txClient);
        }
        return this.ctx.withPgDataTransaction(worldId, guildId, fn);
    }

    async countThreads(
        worldId: number,
        guildId: number,
        options: RepositoryTxOptions = {}
    ): Promise<number> {
        const client = options.txClient ?? this.pool(worldId, guildId);
        const res = await client.query(
            `SELECT COUNT(*)::int AS cnt
             FROM guild_bulletin_board_threads
             WHERE world_id = $1 AND guild_id = $2`,
            [worldId, guildId]
        );
        const row = res.rows[0] as { cnt: number } | undefined;
        if (!row) {
            return 0;
        }
        return Number(row.cnt);
    }

    async getNoticeThread(
        worldId: number,
        guildId: number,
        options: RepositoryTxOptions = {}
    ): Promise<GuildBulletinBoardThreadModel | null> {
        const client = options.txClient ?? this.pool(worldId, guildId);
        const res = await client.query(
            `SELECT t.world_id, t.guild_id, t.local_thread_id, t.poster_character_id, t.title, t.body, t.icon,
                    t.created_at, t.updated_at,
                    COALESCE(r.reply_count, 0) AS reply_count
             FROM guild_bulletin_board_threads t
             LEFT JOIN (
                 SELECT world_id, guild_id, local_thread_id, COUNT(*)::int AS reply_count
                 FROM guild_bulletin_board_replies
                 WHERE world_id = $1 AND guild_id = $2
                 GROUP BY world_id, guild_id, local_thread_id
             ) r ON r.world_id = t.world_id
                AND r.guild_id = t.guild_id
                AND r.local_thread_id = t.local_thread_id
             WHERE t.world_id = $1 AND t.guild_id = $2 AND t.local_thread_id = $3`,
            [worldId, guildId, NOTICE_LOCAL_THREAD_ID]
        );
        const row = res.rows[0] as GuildBulletinBoardThreadRow | undefined;
        if (!row) {
            return null;
        }
        return this.rowToThreadModel(row);
    }

    async listThreads(
        worldId: number,
        guildId: number,
        offset: number,
        limit: number,
        options: RepositoryTxOptions = {}
    ): Promise<GuildBulletinBoardThreadModel[]> {
        const client = options.txClient ?? this.pool(worldId, guildId);
        const res = await client.query(
            `SELECT t.world_id, t.guild_id, t.local_thread_id, t.poster_character_id, t.title, t.body, t.icon,
                    t.created_at, t.updated_at,
                    COALESCE(r.reply_count, 0) AS reply_count
             FROM guild_bulletin_board_threads t
             LEFT JOIN (
                 SELECT world_id, guild_id, local_thread_id, COUNT(*)::int AS reply_count
                 FROM guild_bulletin_board_replies
                 WHERE world_id = $1 AND guild_id = $2
                 GROUP BY world_id, guild_id, local_thread_id
             ) r ON r.world_id = t.world_id
                AND r.guild_id = t.guild_id
                AND r.local_thread_id = t.local_thread_id
             WHERE t.world_id = $1 AND t.guild_id = $2
             ORDER BY t.local_thread_id DESC
             LIMIT $3 OFFSET $4`,
            [worldId, guildId, limit, offset]
        );
        return (res.rows as GuildBulletinBoardThreadRow[]).map((row) => this.rowToThreadModel(row));
    }

    async getThread(
        worldId: number,
        guildId: number,
        localThreadId: number,
        options: RepositoryTxOptions = {}
    ): Promise<GuildBulletinBoardThreadModel | null> {
        const client = options.txClient ?? this.pool(worldId, guildId);
        const res = await client.query(
            `SELECT t.world_id, t.guild_id, t.local_thread_id, t.poster_character_id, t.title, t.body, t.icon,
                    t.created_at, t.updated_at
             FROM guild_bulletin_board_threads t
             WHERE t.world_id = $1 AND t.guild_id = $2 AND t.local_thread_id = $3`,
            [worldId, guildId, localThreadId]
        );
        const row = res.rows[0] as GuildBulletinBoardThreadRow | undefined;
        if (!row) {
            return null;
        }
        return this.rowToThreadModel(row);
    }

    async getThreadWithReplies(
        worldId: number,
        guildId: number,
        localThreadId: number,
        options: RepositoryTxOptions = {}
    ): Promise<GuildBulletinBoardThreadWithRepliesModel | null> {
        const client = options.txClient ?? this.pool(worldId, guildId);
        const res = await client.query(
            `SELECT t.world_id, t.guild_id, t.local_thread_id, t.poster_character_id, t.title, t.body, t.icon,
                    t.created_at, t.updated_at,
                    r.reply_id,
                    r.poster_character_id AS reply_poster_character_id,
                    r.content AS reply_content,
                    r.created_at AS reply_created_at
             FROM guild_bulletin_board_threads t
             LEFT JOIN guild_bulletin_board_replies r
                ON r.world_id = t.world_id
               AND r.guild_id = t.guild_id
               AND r.local_thread_id = t.local_thread_id
             WHERE t.world_id = $1 AND t.guild_id = $2 AND t.local_thread_id = $3
             ORDER BY r.reply_id ASC NULLS FIRST`,
            [worldId, guildId, localThreadId]
        );
        const rows = res.rows as GuildBulletinBoardThreadWithRepliesRow[];
        const first = rows[0];
        if (!first) {
            return null;
        }
        const thread = this.rowToThreadModel(first);
        const replies: GuildBulletinBoardReplyModel[] = [];
        for (const row of rows) {
            if (row.reply_id == null) {
                continue;
            }
            replies.push(
                this.rowToReplyModel({
                    world_id: row.world_id,
                    guild_id: row.guild_id,
                    local_thread_id: row.local_thread_id,
                    reply_id: row.reply_id,
                    poster_character_id: row.reply_poster_character_id,
                    content: row.reply_content,
                    created_at: row.reply_created_at,
                } as GuildBulletinBoardReplyRow)
            );
        }
        return { thread, replies };
    }

    async listReplies(
        worldId: number,
        guildId: number,
        localThreadId: number,
        options: RepositoryTxOptions = {}
    ): Promise<GuildBulletinBoardReplyModel[]> {
        const client = options.txClient ?? this.pool(worldId, guildId);
        const res = await client.query(
            `SELECT world_id, guild_id, local_thread_id, reply_id, poster_character_id, content, created_at
             FROM guild_bulletin_board_replies
             WHERE world_id = $1 AND guild_id = $2 AND local_thread_id = $3
             ORDER BY reply_id ASC`,
            [worldId, guildId, localThreadId]
        );
        return (res.rows as GuildBulletinBoardReplyRow[]).map((row) => this.rowToReplyModel(row));
    }

    async getReply(
        worldId: number,
        guildId: number,
        localThreadId: number,
        replyId: number,
        options: RepositoryTxOptions = {}
    ): Promise<GuildBulletinBoardReplyModel | null> {
        const client = options.txClient ?? this.pool(worldId, guildId);
        const res = await client.query(
            `SELECT world_id, guild_id, local_thread_id, reply_id, poster_character_id, content, created_at
             FROM guild_bulletin_board_replies
             WHERE world_id = $1 AND guild_id = $2 AND local_thread_id = $3 AND reply_id = $4`,
            [worldId, guildId, localThreadId, replyId]
        );
        const row = res.rows[0] as GuildBulletinBoardReplyRow | undefined;
        if (!row) {
            return null;
        }
        return this.rowToReplyModel(row);
    }

    async createThread(
        worldId: number,
        input: CreateGuildBulletinBoardThreadInput,
        options: RepositoryTxOptions = {}
    ): Promise<GuildBulletinBoardThreadModel> {
        const createdAt = input.createdAt ?? new Date();
        return this.runInDataTransaction(worldId, input.guildId, options, async (client) => {
            const guildRes = await client.query(
                `SELECT guild_id
                 FROM guilds
                 WHERE world_id = $1 AND guild_id = $2 AND disbanded_at IS NULL
                 FOR UPDATE`,
                [worldId, input.guildId]
            );
            if ((guildRes.rowCount ?? 0) === 0) {
                throw new Error(`guild not found for bulletin thread: world=${worldId} guild=${input.guildId}`);
            }

            const nextRes = await client.query(
                `SELECT
                    CASE
                        WHEN $3::boolean THEN $4::int
                        ELSE GREATEST(1, COALESCE(MAX(local_thread_id), 0) + 1)
                    END AS next_local_thread_id
                 FROM guild_bulletin_board_threads
                 WHERE world_id = $1 AND guild_id = $2`,
                [worldId, input.guildId, input.isNotice, NOTICE_LOCAL_THREAD_ID]
            );

            const localThreadId = Number(
                (nextRes.rows[0] as { next_local_thread_id: number | string }).next_local_thread_id
            );

            const res = await client.query(
                `INSERT INTO guild_bulletin_board_threads
                    (world_id, guild_id, local_thread_id, poster_character_id, title, body, icon, created_at, updated_at)
                 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $8)
                 RETURNING world_id, guild_id, local_thread_id, poster_character_id, title, body, icon, created_at, updated_at`,
                [
                    worldId,
                    input.guildId,
                    localThreadId,
                    input.posterCharacterId,
                    input.title,
                    input.body,
                    input.icon,
                    createdAt,
                ]
            );
            return this.rowToThreadModel(res.rows[0] as GuildBulletinBoardThreadRow);
        });
    }

    async updateThread(
        worldId: number,
        model: GuildBulletinBoardThreadModel,
        options: RepositoryTxOptions = {}
    ): Promise<GuildBulletinBoardThreadModel | null> {
        const client = options.txClient ?? this.pool(worldId, model.guildId);
        const res = await client.query(
            `UPDATE guild_bulletin_board_threads
             SET title = $4, body = $5, icon = $6, updated_at = $7
             WHERE world_id = $1 AND guild_id = $2 AND local_thread_id = $3
             RETURNING world_id, guild_id, local_thread_id, poster_character_id, title, body, icon, created_at, updated_at`,
            [
                worldId,
                model.guildId,
                model.localThreadId,
                model.title,
                model.body,
                model.icon,
                model.updatedAt,
            ]
        );
        const row = res.rows[0] as GuildBulletinBoardThreadRow | undefined;
        if (!row) {
            return null;
        }
        return this.rowToThreadModel(row);
    }

    async deleteThread(
        worldId: number,
        guildId: number,
        localThreadId: number,
        options: RepositoryTxOptions = {}
    ): Promise<boolean> {
        return this.runInDataTransaction(worldId, guildId, options, async (client) => {
            await client.query(
                `DELETE FROM guild_bulletin_board_replies
                 WHERE world_id = $1 AND guild_id = $2 AND local_thread_id = $3`,
                [worldId, guildId, localThreadId]
            );
            const res = await client.query(
                `DELETE FROM guild_bulletin_board_threads
                 WHERE world_id = $1 AND guild_id = $2 AND local_thread_id = $3`,
                [worldId, guildId, localThreadId]
            );
            return (res.rowCount ?? 0) > 0;
        });
    }

    async createReply(
        worldId: number,
        input: CreateGuildBulletinBoardReplyInput,
        options: RepositoryTxOptions = {}
    ): Promise<GuildBulletinBoardReplyModel | null> {
        const createdAt = input.createdAt ?? new Date();
        return this.runInDataTransaction(worldId, input.guildId, options, async (client) => {
            const threadRes = await client.query(
                `SELECT local_thread_id
                 FROM guild_bulletin_board_threads
                 WHERE world_id = $1 AND guild_id = $2 AND local_thread_id = $3
                 FOR UPDATE`,
                [worldId, input.guildId, input.localThreadId]
            );
            if ((threadRes.rowCount ?? 0) === 0) {
                return null;
            }

            const nextRes = await client.query(
                `SELECT COALESCE(MAX(reply_id), -1) + 1 AS next_reply_id
                 FROM guild_bulletin_board_replies
                 WHERE world_id = $1 AND guild_id = $2 AND local_thread_id = $3
                 FOR UPDATE`,
                [worldId, input.guildId, input.localThreadId]
            );
            const replyId = Number(
                (nextRes.rows[0] as { next_reply_id: number | string }).next_reply_id
            );

            const res = await client.query(
                `INSERT INTO guild_bulletin_board_replies
                    (world_id, guild_id, local_thread_id, reply_id, poster_character_id, content, created_at)
                 VALUES ($1, $2, $3, $4, $5, $6, $7)
                 RETURNING world_id, guild_id, local_thread_id, reply_id, poster_character_id, content, created_at`,
                [
                    worldId,
                    input.guildId,
                    input.localThreadId,
                    replyId,
                    input.posterCharacterId,
                    input.content,
                    createdAt,
                ]
            );
            return this.rowToReplyModel(res.rows[0] as GuildBulletinBoardReplyRow);
        });
    }

    async deleteReply(
        worldId: number,
        guildId: number,
        localThreadId: number,
        replyId: number,
        options: RepositoryTxOptions = {}
    ): Promise<boolean> {
        const client = options.txClient ?? this.pool(worldId, guildId);
        const res = await client.query(
            `DELETE FROM guild_bulletin_board_replies
             WHERE world_id = $1 AND guild_id = $2 AND local_thread_id = $3 AND reply_id = $4`,
            [worldId, guildId, localThreadId, replyId]
        );
        return (res.rowCount ?? 0) > 0;
    }

    async deleteAllForGuild(worldId: number, guildId: number, options: RepositoryTxOptions = {}): Promise<void> {
        await this.runInDataTransaction(worldId, guildId, options, async (client) => {
            await client.query(
                `DELETE FROM guild_bulletin_board_replies WHERE world_id = $1 AND guild_id = $2`,
                [worldId, guildId]
            );
            await client.query(
                `DELETE FROM guild_bulletin_board_threads WHERE world_id = $1 AND guild_id = $2`,
                [worldId, guildId]
            );
        });
    }

    rowToThreadModel(row: GuildBulletinBoardThreadRow): GuildBulletinBoardThreadModel {
        return {
            worldId: row.world_id,
            guildId: toPgInt(row.guild_id),
            localThreadId: row.local_thread_id,
            posterCharacterId: row.poster_character_id,
            title: row.title,
            body: row.body,
            icon: row.icon,
            createdAt: row.created_at instanceof Date ? row.created_at : new Date(row.created_at),
            updatedAt: row.updated_at instanceof Date ? row.updated_at : new Date(row.updated_at),
            replyCount: row.reply_count == null ? undefined : Number(row.reply_count),
        };
    }

    rowToReplyModel(row: GuildBulletinBoardReplyRow): GuildBulletinBoardReplyModel {
        return {
            worldId: row.world_id,
            guildId: toPgInt(row.guild_id),
            localThreadId: row.local_thread_id,
            replyId: row.reply_id,
            posterCharacterId: row.poster_character_id,
            content: row.content,
            createdAt: row.created_at instanceof Date ? row.created_at : new Date(row.created_at),
        };
    }
}
