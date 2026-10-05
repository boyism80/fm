import {
    GuildErrorCode,
    GuildMemberRank,
    type GuildBulletinBoardReplyEntry,
    type GuildBulletinBoardThreadDetail,
    type GuildBulletinBoardThreadEntry,
} from "../protobuf/generated/fminternal/internal_service";
import { AppConfiguration } from "../config/app-configuration";
import { InternalContext } from "../context/internal-context";
import { GuildRepository } from "../repos/guild-repository";
import { GuildMemberRepository } from "../repos/guild-member-repository";
import { GuildBulletinBoardRepository } from "../repos/guild-bulletin-board-repository";
import { CharacterRealtimeStateRepository } from "../repos/character-realtime-state-repository";
import type { GuildMemberModel } from "../repos/guild-member-repository";
import type { GuildBulletinBoardThreadModel } from "../types/repository-models";
import { DistributedLockService } from "./distributed-lock-service";
import { redisCacheKey } from "../redis-cache-key";

const messages = { GuildErrorCode };

const MAX_GUILD_BULLETIN_TITLE_LEN = 25;
const MAX_GUILD_BULLETIN_BODY_LEN = 600;
const MAX_GUILD_BULLETIN_REPLY_LEN = 25;
const GUILD_BULLETIN_THREADS_PER_PAGE = 10;
const GUILD_BULLETIN_ICON_CASH_MIN = 0x64;
const GUILD_BULLETIN_ICON_CASH_MAX = 0x6a;
const GUILD_BULLETIN_COOLDOWN_SEC = 60;

export type ListGuildBulletinBoardThreadsResult = {
    ok: boolean;
    code?: GuildErrorCode;
    threads?: GuildBulletinBoardThreadEntry[];
    listStart?: number;
    threadCount?: number;
    notice?: GuildBulletinBoardThreadEntry;
};

export type ShowGuildBulletinBoardThreadResult = {
    ok: boolean;
    code?: GuildErrorCode;
    thread?: GuildBulletinBoardThreadDetail;
};

export type CreateGuildBulletinBoardThreadResult = {
    ok: boolean;
    code?: GuildErrorCode;
    thread?: GuildBulletinBoardThreadDetail;
    threads?: GuildBulletinBoardThreadEntry[];
    listStart?: number;
    threadCount?: number;
    notice?: GuildBulletinBoardThreadEntry;
};

export type UpdateGuildBulletinBoardThreadResult = {
    ok: boolean;
    code?: GuildErrorCode;
    thread?: GuildBulletinBoardThreadDetail;
};

export type DeleteGuildBulletinBoardThreadResult = {
    ok: boolean;
    code?: GuildErrorCode;
};

export type CreateGuildBulletinBoardReplyResult = {
    ok: boolean;
    code?: GuildErrorCode;
    thread?: GuildBulletinBoardThreadDetail;
};

export type DeleteGuildBulletinBoardReplyResult = {
    ok: boolean;
    code?: GuildErrorCode;
    thread?: GuildBulletinBoardThreadDetail;
};

export type GuildMembershipResult =
    | { code: GuildErrorCode.GUILD_ERROR_NONE; guildId: number; member: GuildMemberModel }
    | { code: Exclude<GuildErrorCode, GuildErrorCode.GUILD_ERROR_NONE> };

export class GuildBulletinBoardService {
    private readonly ctx: InternalContext;
    private readonly app: AppConfiguration;
    private readonly guildRepo: GuildRepository;
    private readonly guildMemberRepo: GuildMemberRepository;
    private readonly guildBulletinBoardRepo: GuildBulletinBoardRepository;
    private readonly characterRealtimeStateRepo: CharacterRealtimeStateRepository;
    private readonly distributedLockService: DistributedLockService;

    constructor(
        internalContext: InternalContext,
        appConfiguration: AppConfiguration,
        guildRepository: GuildRepository,
        guildMemberRepository: GuildMemberRepository,
        guildBulletinBoardRepository: GuildBulletinBoardRepository,
        characterRealtimeStateRepository: CharacterRealtimeStateRepository,
        distributedLockService: DistributedLockService
    ) {
        this.ctx = internalContext;
        this.app = appConfiguration;
        this.guildRepo = guildRepository;
        this.guildMemberRepo = guildMemberRepository;
        this.guildBulletinBoardRepo = guildBulletinBoardRepository;
        this.characterRealtimeStateRepo = characterRealtimeStateRepository;
        this.distributedLockService = distributedLockService;
    }

    private assertWorld(worldId: number) {
        const wid = String(worldId);
        if (!this.app.postgresql.worlds[wid]) {
            const err = new Error(`Unknown world_id: ${worldId}`) as Error & { code?: string };
            err.code = "UNKNOWN_WORLD";
            throw err;
        }
    }

    private assertCharacterId(characterId: number) {
        const n = characterId;
        if (!Number.isInteger(n) || n <= 0 || n > 0xffffffff) {
            const err = new Error("character_id must be a positive uint32") as Error & { code?: string };
            err.code = "INVALID_CHARACTER_ID";
            throw err;
        }
    }

    private canChangeMemberRank(rank: GuildMemberRank | undefined) {
        return rank === GuildMemberRank.GUILD_MEMBER_RANK_MASTER
            || rank === GuildMemberRank.GUILD_MEMBER_RANK_JUNIOR;
    }

    private truncateBulletinField(value: string, maxLen: number): string {
        if (value.length <= maxLen) {
            return value;
        }
        return value.substring(0, maxLen);
    }

    private isValidBulletinIcon(icon: number): boolean {
        if (icon >= 0 && icon <= 2) {
            return true;
        }
        if (icon >= GUILD_BULLETIN_ICON_CASH_MIN && icon <= GUILD_BULLETIN_ICON_CASH_MAX) {
            return true;
        }
        return false;
    }

    private canModifyBulletinContent(
        posterCharacterId: number,
        characterId: number,
        guildRank: GuildMemberRank | undefined
    ): boolean {
        if (posterCharacterId === characterId) {
            return true;
        }
        return this.canChangeMemberRank(guildRank);
    }

    private bulletinThreadToEntry(thread: GuildBulletinBoardThreadModel): GuildBulletinBoardThreadEntry {
        return {
            localThreadId: thread.localThreadId,
            posterCharacterId: thread.posterCharacterId,
            title: thread.title,
            timestampUnixMs: thread.createdAt.getTime(),
            icon: thread.icon,
            replyCount: thread.replyCount ?? 0,
        };
    }

    private async loadGuildBulletinBoardThreadListPage(
        worldId: number,
        guildId: number,
        listStart: number
    ): Promise<{
        threads: GuildBulletinBoardThreadEntry[];
        threadCount: number;
        notice?: GuildBulletinBoardThreadEntry;
        listStart: number;
    }> {
        const offset = Math.max(0, listStart);
        const [threadCount, pageRows, noticeRow] = await Promise.all([
            this.guildBulletinBoardRepo.countThreads(worldId, guildId),
            this.guildBulletinBoardRepo.listThreads(
                worldId,
                guildId,
                offset,
                GUILD_BULLETIN_THREADS_PER_PAGE
            ),
            this.guildBulletinBoardRepo.getNoticeThread(worldId, guildId),
        ]);
        const threads = pageRows.map((t) => this.bulletinThreadToEntry(t));
        const notice = noticeRow ? this.bulletinThreadToEntry(noticeRow) : undefined;
        return {
            threads,
            threadCount,
            notice,
            listStart: offset,
        };
    }

    private async loadBulletinThreadDetail(
        worldId: number,
        guildId: number,
        localThreadId: number
    ): Promise<GuildBulletinBoardThreadDetail | null> {
        const loaded = await this.guildBulletinBoardRepo.getThreadWithReplies(worldId, guildId, localThreadId);
        if (!loaded) {
            return null;
        }
        const replies: GuildBulletinBoardReplyEntry[] = loaded.replies.map((r) => {
            return {
                replyId: r.replyId,
                posterCharacterId: r.posterCharacterId,
                timestampUnixMs: r.createdAt.getTime(),
                content: r.content,
            };
        });
        return {
            localThreadId: loaded.thread.localThreadId,
            posterCharacterId: loaded.thread.posterCharacterId,
            timestampUnixMs: loaded.thread.createdAt.getTime(),
            title: loaded.thread.title,
            body: loaded.thread.body,
            icon: loaded.thread.icon,
            replies,
        };
    }

    private guildBulletinBoardCooldownKey(worldId: number, characterId: number) {
        return redisCacheKey(`w${worldId}:guild-bbs-cooldown:${characterId}`);
    }

    private async setBulletinBoardCooldown(worldId: number, characterId: number) {
        const { client } = this.ctx.getRedisGlobalAccess(worldId);
        const key = this.guildBulletinBoardCooldownKey(worldId, characterId);
        const acquired = await client.set(key, "1", "EX", GUILD_BULLETIN_COOLDOWN_SEC, "NX");
        return acquired === "OK";
    }

    private async unsetBulletinBoardCooldown(worldId: number, characterId: number) {
        const { client } = this.ctx.getRedisGlobalAccess(worldId);
        const key = this.guildBulletinBoardCooldownKey(worldId, characterId);
        await client.del(key);
    }

    private async assertGuildMembership(
        worldId: number,
        characterId: number
    ): Promise<GuildMembershipResult> {
        const state = await this.characterRealtimeStateRepo.get(worldId, characterId);
        if (state?.guildId == null) {
            return { code: messages.GuildErrorCode.GUILD_ERROR_NOT_IN_GUILD };
        }
        const guildId = state.guildId;
        if (!Number.isInteger(guildId) || guildId < 1) {
            return { code: messages.GuildErrorCode.GUILD_ERROR_NOT_IN_GUILD };
        }
        const guild = await this.guildRepo.get(worldId, guildId);
        if (!guild) {
            return { code: messages.GuildErrorCode.GUILD_ERROR_GUILD_NOT_FOUND };
        }
        const member = await this.guildMemberRepo.getItem(worldId, String(guildId), characterId);
        if (!member) {
            return { code: messages.GuildErrorCode.GUILD_ERROR_NOT_IN_GUILD };
        }
        return {
            code: messages.GuildErrorCode.GUILD_ERROR_NONE,
            guildId,
            member,
        };
    }

    async listGuildBulletinBoardThreads(
        worldId: number,
        characterId: number,
        page: number
    ): Promise<ListGuildBulletinBoardThreadsResult> {
        this.assertWorld(worldId);
        this.assertCharacterId(characterId);

        const membership = await this.assertGuildMembership(worldId, characterId);
        if (membership.code !== messages.GuildErrorCode.GUILD_ERROR_NONE) {
            return { ok: false, code: membership.code };
        }

        const listStart = Math.max(0, page) * GUILD_BULLETIN_THREADS_PER_PAGE;
        const listPage = await this.loadGuildBulletinBoardThreadListPage(
            worldId,
            membership.guildId,
            listStart
        );
        return {
            ok: true,
            threads: listPage.threads,
            listStart: listPage.listStart,
            threadCount: listPage.threadCount,
            notice: listPage.notice,
        };
    }

    async showGuildBulletinBoardThread(
        worldId: number,
        characterId: number,
        localThreadId: number
    ): Promise<ShowGuildBulletinBoardThreadResult> {
        this.assertWorld(worldId);
        this.assertCharacterId(characterId);

        const membership = await this.assertGuildMembership(worldId, characterId);
        if (membership.code !== messages.GuildErrorCode.GUILD_ERROR_NONE) {
            return { ok: false, code: membership.code };
        }

        const thread = await this.loadBulletinThreadDetail(worldId, membership.guildId, localThreadId);
        if (!thread) {
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_BULLETIN_THREAD_NOT_FOUND };
        }
        return { ok: true, thread };
    }

    async createGuildBulletinBoardThread(
        worldId: number,
        characterId: number,
        notice: boolean,
        title: string,
        body: string,
        icon: number
    ): Promise<CreateGuildBulletinBoardThreadResult> {
        this.assertWorld(worldId);
        this.assertCharacterId(characterId);

        const membership = await this.assertGuildMembership(worldId, characterId);
        if (membership.code !== messages.GuildErrorCode.GUILD_ERROR_NONE) {
            return { ok: false, code: membership.code };
        }

        if (!this.isValidBulletinIcon(icon)) {
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_INVALID_BULLETIN_CONTENT };
        }
        const normalizedTitle = this.truncateBulletinField(title, MAX_GUILD_BULLETIN_TITLE_LEN);
        const normalizedBody = this.truncateBulletinField(body, MAX_GUILD_BULLETIN_BODY_LEN);

        await using _guildLock = await this.distributedLockService.acquireWorldDataLock(
            worldId,
            `guild:${membership.guildId}`,
        );

        const isCooldown = await this.setBulletinBoardCooldown(worldId, characterId);
        if (!isCooldown) {
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_BULLETIN_COOLDOWN };
        }

        try {
            const created = await this.guildBulletinBoardRepo.createThread(worldId, {
                guildId: membership.guildId,
                posterCharacterId: characterId,
                title: normalizedTitle,
                body: normalizedBody,
                icon,
                isNotice: notice,
            });
            const thread = await this.loadBulletinThreadDetail(worldId, membership.guildId, created.localThreadId);
            if (!thread) {
                return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_UNKNOWN };
            }
            const page = await this.loadGuildBulletinBoardThreadListPage(worldId, membership.guildId, 0);
            return {
                ok: true,
                thread,
                threads: page.threads,
                listStart: page.listStart,
                threadCount: page.threadCount,
                notice: page.notice,
            };
        } catch {
            await this.unsetBulletinBoardCooldown(worldId, characterId);
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_UNKNOWN };
        }
    }

    async updateGuildBulletinBoardThread(
        worldId: number,
        characterId: number,
        localThreadId: number,
        title: string,
        body: string,
        icon: number
    ): Promise<UpdateGuildBulletinBoardThreadResult> {
        this.assertWorld(worldId);
        this.assertCharacterId(characterId);

        const membership = await this.assertGuildMembership(worldId, characterId);
        if (membership.code !== messages.GuildErrorCode.GUILD_ERROR_NONE) {
            return { ok: false, code: membership.code };
        }

        if (!this.isValidBulletinIcon(icon)) {
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_INVALID_BULLETIN_CONTENT };
        }
        const normalizedTitle = this.truncateBulletinField(title, MAX_GUILD_BULLETIN_TITLE_LEN);
        const normalizedBody = this.truncateBulletinField(body, MAX_GUILD_BULLETIN_BODY_LEN);

        const existing = await this.guildBulletinBoardRepo.getThread(worldId, membership.guildId, localThreadId);
        if (!existing) {
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_BULLETIN_THREAD_NOT_FOUND };
        }
        if (!this.canModifyBulletinContent(existing.posterCharacterId, characterId, membership.member.guildRank)) {
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_NOT_AUTHORIZED };
        }
        existing.title = normalizedTitle;
        existing.body = normalizedBody;
        existing.icon = icon;
        existing.updatedAt = new Date();
        const updated = await this.guildBulletinBoardRepo.updateThread(worldId, existing);
        if (!updated) {
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_BULLETIN_THREAD_NOT_FOUND };
        }
        const thread = await this.loadBulletinThreadDetail(worldId, membership.guildId, localThreadId);
        if (!thread) {
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_BULLETIN_THREAD_NOT_FOUND };
        }
        return { ok: true, thread };
    }

    async deleteGuildBulletinBoardThread(
        worldId: number,
        characterId: number,
        localThreadId: number
    ): Promise<DeleteGuildBulletinBoardThreadResult> {
        this.assertWorld(worldId);
        this.assertCharacterId(characterId);

        const membership = await this.assertGuildMembership(worldId, characterId);
        if (membership.code !== messages.GuildErrorCode.GUILD_ERROR_NONE) {
            return { ok: false, code: membership.code };
        }

        const existing = await this.guildBulletinBoardRepo.getThread(worldId, membership.guildId, localThreadId);
        if (!existing) {
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_BULLETIN_THREAD_NOT_FOUND };
        }
        if (!this.canModifyBulletinContent(existing.posterCharacterId, characterId, membership.member.guildRank)) {
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_NOT_AUTHORIZED };
        }
        await this.guildBulletinBoardRepo.deleteThread(worldId, membership.guildId, localThreadId);
        return { ok: true };
    }

    async createGuildBulletinBoardReply(
        worldId: number,
        characterId: number,
        localThreadId: number,
        content: string
    ): Promise<CreateGuildBulletinBoardReplyResult> {
        this.assertWorld(worldId);
        this.assertCharacterId(characterId);

        const membership = await this.assertGuildMembership(worldId, characterId);
        if (membership.code !== messages.GuildErrorCode.GUILD_ERROR_NONE) {
            return { ok: false, code: membership.code };
        }

        const normalizedContent = this.truncateBulletinField(content, MAX_GUILD_BULLETIN_REPLY_LEN);

        await using _guildLock = await this.distributedLockService.acquireWorldDataLock(
            worldId,
            `guild:${membership.guildId}`,
        );

        const isCooldown = await this.setBulletinBoardCooldown(worldId, characterId);
        if (!isCooldown) {
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_BULLETIN_COOLDOWN };
        }

        try {
            const created = await this.guildBulletinBoardRepo.createReply(worldId, {
                guildId: membership.guildId,
                localThreadId,
                posterCharacterId: characterId,
                content: normalizedContent,
            });
            if (!created) {
                await this.unsetBulletinBoardCooldown(worldId, characterId);
                return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_BULLETIN_THREAD_NOT_FOUND };
            }
            const thread = await this.loadBulletinThreadDetail(worldId, membership.guildId, localThreadId);
            if (!thread) {
                return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_BULLETIN_THREAD_NOT_FOUND };
            }
            return { ok: true, thread };
        } catch {
            await this.unsetBulletinBoardCooldown(worldId, characterId);
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_UNKNOWN };
        }
    }

    async deleteGuildBulletinBoardReply(
        worldId: number,
        characterId: number,
        localThreadId: number,
        replyId: number
    ): Promise<DeleteGuildBulletinBoardReplyResult> {
        this.assertWorld(worldId);
        this.assertCharacterId(characterId);

        const membership = await this.assertGuildMembership(worldId, characterId);
        if (membership.code !== messages.GuildErrorCode.GUILD_ERROR_NONE) {
            return { ok: false, code: membership.code };
        }

        const target = await this.guildBulletinBoardRepo.getReply(
            worldId,
            membership.guildId,
            localThreadId,
            replyId
        );
        if (!target) {
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_BULLETIN_REPLY_NOT_FOUND };
        }
        if (!this.canModifyBulletinContent(target.posterCharacterId, characterId, membership.member.guildRank)) {
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_NOT_AUTHORIZED };
        }
        const deleted = await this.guildBulletinBoardRepo.deleteReply(
            worldId,
            membership.guildId,
            localThreadId,
            replyId
        );
        if (!deleted) {
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_BULLETIN_REPLY_NOT_FOUND };
        }
        const thread = await this.loadBulletinThreadDetail(worldId, membership.guildId, localThreadId);
        if (!thread) {
            return { ok: false, code: messages.GuildErrorCode.GUILD_ERROR_BULLETIN_THREAD_NOT_FOUND };
        }
        return { ok: true, thread };
    }
}
