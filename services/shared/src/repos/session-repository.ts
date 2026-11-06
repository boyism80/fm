import type Redis from "ioredis";
import { redisSessionKey } from "../redis-session-key";
import { AccountSessionState, accountSessionStateFromRedisHash } from "../session-state";
import { SessionErrorCode } from "../protobuf/generated/fminternal/internal_service";
import type { CharacterPresentation, Debuff } from "../protobuf/generated/fminternal/internal_service";
import type { InternalContext } from "../context/internal-context";
import type { AccountSession } from "../session-types";
import type { CharacterRepository } from "./character-repository";
import { redisLockKey, type DistributedLock } from "../system/distributed-lock";

export type { AccountSession } from "../session-types";

type SessionHash = Record<string, string | undefined>;

const AS = AccountSessionState;

export class SessionRepository {
    private readonly ctx: InternalContext;
    private readonly characterRepo: CharacterRepository;
    private readonly distributedLock: DistributedLock;

    constructor(
        internalContext: InternalContext,
        characterRepository: CharacterRepository,
        distributedLock: DistributedLock,
    ) {
        this.ctx = internalContext;
        this.characterRepo = characterRepository;
        this.distributedLock = distributedLock;
    }

    private accountKey(accountId: number) {
        return redisSessionKey(`account:${accountId}`);
    }

    private accountHashToSession(hash: SessionHash): AccountSession | null {
        if (!hash || Object.keys(hash).length === 0) {
            return null;
        }
        const toNumberOrNull = (v: string | undefined) => (v === "" || v == null ? null : Number(v));
        return {
            version: Number(hash.version ?? 1),
            worldId: toNumberOrNull(hash.world_id),
            accountId: Number(hash.account_id ?? 0),
            state: accountSessionStateFromRedisHash(hash.state),
            loginServer: {
                id: hash.login_server_id || null,
                connected: (hash.login_server_connected || "0") === "1",
            },
            characterId: toNumberOrNull(hash.character_id),
            characterName: hash.character_name || null,
            channelId: toNumberOrNull(hash.channel_id),
            cashShopId: toNumberOrNull(hash.cash_shop_id),
            returnChannelId: toNumberOrNull(hash.return_channel_id),
            gameToGameTransfer: (hash.game_to_game_transfer || "0") === "1",
            debuffs: hash.debuffs ? (JSON.parse(hash.debuffs) as Debuff[]) : [],
            presentation: hash.presentation ? (JSON.parse(hash.presentation) as CharacterPresentation) : null,
            timestamps: {
                createdAt: hash.created_at || null,
                updatedAt: hash.updated_at || null,
                stateChangedAt: hash.state_changed_at || null,
            },
        };
    }

    // Every change reads the session and writes it back under the account's lock; a request that cannot take the lock fails.
    private async locked<T>(
        accountId: number,
        fn: (session: SessionHash, key: string, client: Redis) => Promise<T>,
    ): Promise<T> {
        const { client } = this.ctx.getRedisUnifiedAccess();
        await using _lock = await this.distributedLock.acquire(client, redisLockKey(`session:${accountId}`));
        const key = this.accountKey(accountId);
        return fn(await client.hgetall(key), key, client);
    }

    private exists(session: SessionHash) {
        return Object.keys(session).length > 0;
    }

    private isOwner(session: SessionHash, state: number, characterId: number, idField: string, id: number) {
        return (
            Number(session.state ?? 0) === state &&
            session.character_id === String(characterId) &&
            session[idField] === String(id)
        );
    }

    async beginLogin(accountId: number, loginServerId: string, now: string, ttlSeconds: number) {
        return this.locked(accountId, async (session, key, client): Promise<Array<number | string>> => {
            if (this.exists(session)) {
                return [0, SessionErrorCode.SESSION_ALREADY_LOGGED_IN];
            }
            await client
                .multi()
                .hset(key, {
                    version: "1",
                    world_id: "",
                    account_id: String(accountId),
                    state: String(AS.ACCOUNT_SESSION_STATE_LOGIN),
                    login_server_id: String(loginServerId),
                    login_server_connected: "1",
                    created_at: now,
                    updated_at: now,
                    state_changed_at: now,
                })
                .expire(key, ttlSeconds)
                .exec();
            return [1, SessionErrorCode.SESSION_NONE];
        });
    }

    async beginTransition(
        worldId: number,
        accountId: number,
        characterId: number,
        characterName: string,
        clientIp: string,
        now: string,
        ttlSeconds: number,
        sourceChannelId: number | null,
        sourceCashShopId: number | null,
        debuffs: Debuff[],
        presentation: CharacterPresentation | null,
    ) {
        return this.locked(accountId, async (session, key, client): Promise<Array<number | string>> => {
            if (!this.exists(session)) {
                return [0, SessionErrorCode.SESSION_NOT_FOUND];
            }
            const state = Number(session.state ?? 0);
            const fromServer = sourceChannelId != null || sourceCashShopId != null;
            if (sourceCashShopId != null) {
                if (
                    !this.isOwner(
                        session,
                        AS.ACCOUNT_SESSION_STATE_CASH_SHOP,
                        characterId,
                        "cash_shop_id",
                        sourceCashShopId,
                    )
                ) {
                    return [0, SessionErrorCode.SESSION_NOT_OWNER];
                }
            } else if (sourceChannelId != null) {
                if (!this.isOwner(session, AS.ACCOUNT_SESSION_STATE_GAME, characterId, "channel_id", sourceChannelId)) {
                    return [0, SessionErrorCode.SESSION_NOT_OWNER];
                }
            } else if (
                state !== AS.ACCOUNT_SESSION_STATE_LOGIN &&
                (state !== AS.ACCOUNT_SESSION_STATE_TRANSITION || session.game_to_game_transfer === "1")
            ) {
                return [0, SessionErrorCode.SESSION_NOT_OWNER];
            }

            const write = client
                .multi()
                .hset(key, {
                    world_id: String(worldId),
                    state: String(AS.ACCOUNT_SESSION_STATE_TRANSITION),
                    character_id: String(characterId),
                    character_name: String(characterName),
                    client_ip: clientIp,
                    game_to_game_transfer: fromServer ? "1" : "0",
                    debuffs: fromServer && debuffs.length > 0 ? JSON.stringify(debuffs) : "",
                    presentation: fromServer && presentation != null ? JSON.stringify(presentation) : "",
                    updated_at: now,
                    state_changed_at: now,
                })
                .hdel(key, "channel_id", "cash_shop_id", "return_channel_id");
            if (sourceCashShopId == null && sourceChannelId != null) {
                write.hset(key, "return_channel_id", String(sourceChannelId));
            }
            await write.expire(key, ttlSeconds).exec();
            return [1, SessionErrorCode.SESSION_NONE];
        });
    }

    // Ends a transition that never reached its server: the session goes only if it is still that same transition,
    // so a client entering at the same moment either enters first or is refused.
    async expireTransition(accountId: number, stateChangedAt: string): Promise<boolean> {
        return this.locked(accountId, async (session, key, client) => {
            if (
                session.state !== String(AS.ACCOUNT_SESSION_STATE_TRANSITION) ||
                session.state_changed_at !== stateChangedAt
            ) {
                return false;
            }
            await client.del(key);
            return true;
        });
    }

    private entryError(
        session: SessionHash,
        accountId: number,
        worldId: number,
        characterId: number,
        clientIp: string,
    ): SessionErrorCode | null {
        if (!this.exists(session)) {
            return SessionErrorCode.SESSION_NOT_FOUND;
        }
        if (
            Number(session.state ?? 0) !== AS.ACCOUNT_SESSION_STATE_TRANSITION ||
            Number(session.account_id ?? 0) !== accountId ||
            Number(session.world_id ?? 0) !== worldId ||
            Number(session.character_id ?? 0) !== characterId
        ) {
            return SessionErrorCode.SESSION_UNKNOWN;
        }
        if (clientIp === "" || session.client_ip !== clientIp) {
            return SessionErrorCode.SESSION_CLIENT_IP_MISMATCH;
        }
        return null;
    }

    async enterCashShop(
        worldId: number,
        accountId: number,
        characterId: number,
        cashShopId: number,
        clientIp: string,
        now: string,
        ttlSeconds: number,
    ) {
        return this.locked(accountId, async (session, key, client): Promise<Array<number | string>> => {
            const error = this.entryError(session, accountId, worldId, characterId, clientIp);
            if (error != null) {
                return [0, error];
            }
            const returnChannelId = session.return_channel_id ?? "";
            if (returnChannelId === "") {
                return [0, SessionErrorCode.SESSION_UNKNOWN];
            }
            await client
                .multi()
                .hset(key, {
                    state: String(AS.ACCOUNT_SESSION_STATE_CASH_SHOP),
                    cash_shop_id: String(cashShopId),
                    updated_at: now,
                    state_changed_at: now,
                })
                .expire(key, ttlSeconds)
                .exec();
            return [1, SessionErrorCode.SESSION_NONE, returnChannelId];
        });
    }

    async enterGame(
        worldId: number,
        accountId: number,
        characterId: number,
        channelId: number,
        clientIp: string,
        now: string,
        ttlSeconds: number,
    ) {
        return this.locked(accountId, async (session, key, client): Promise<Array<number | string>> => {
            const error = this.entryError(session, accountId, worldId, characterId, clientIp);
            if (error != null) {
                return [0, error];
            }
            await client
                .multi()
                .hset(key, {
                    state: String(AS.ACCOUNT_SESSION_STATE_GAME),
                    login_server_connected: "0",
                    channel_id: String(channelId),
                    game_to_game_transfer: "0",
                    updated_at: now,
                    state_changed_at: now,
                })
                .hdel(key, "debuffs", "presentation")
                .expire(key, ttlSeconds)
                .exec();
            return [1, SessionErrorCode.SESSION_NONE];
        });
    }

    private ownerMatches(
        session: SessionHash,
        owner: { characterId: number; channelId?: number; cashShopId?: number },
    ): boolean {
        if (owner.cashShopId != null) {
            return this.isOwner(
                session,
                AS.ACCOUNT_SESSION_STATE_CASH_SHOP,
                owner.characterId,
                "cash_shop_id",
                owner.cashShopId,
            );
        }
        return this.isOwner(
            session,
            AS.ACCOUNT_SESSION_STATE_GAME,
            owner.characterId,
            "channel_id",
            owner.channelId ?? -1,
        );
    }

    async refresh(
        accountId: number,
        loginTtl: number,
        transitionTtl: number,
        gameTtl: number,
        owner: { characterId: number; channelId?: number; cashShopId?: number } | null,
    ) {
        return this.locked(accountId, async (session, key, client): Promise<Array<number | string>> => {
            if (!this.exists(session)) {
                return [0, SessionErrorCode.SESSION_NOT_FOUND, 0, "0"];
            }
            const state = Number(session.state ?? 0);
            if (owner != null && !this.ownerMatches(session, owner)) {
                return [0, SessionErrorCode.SESSION_NOT_OWNER, 0, String(state)];
            }
            let ttl = loginTtl;
            if (state === AS.ACCOUNT_SESSION_STATE_TRANSITION) {
                ttl = transitionTtl;
            } else if (state === AS.ACCOUNT_SESSION_STATE_GAME || state === AS.ACCOUNT_SESSION_STATE_CASH_SHOP) {
                ttl = gameTtl;
            }
            await client.expire(key, ttl);
            return [1, SessionErrorCode.SESSION_NONE, ttl, String(state)];
        });
    }

    async logout(
        accountId: number,
        keep: boolean,
        requestChannelId?: number,
        ownerCharacterId?: number,
        requestCashShopId?: number,
    ) {
        return this.locked(accountId, async (session, key, client): Promise<Array<number | string>> => {
            const state = Number(session.state ?? 0);
            if (keep) {
                return [1, SessionErrorCode.SESSION_NONE, state];
            }
            if (
                ownerCharacterId != null &&
                !this.ownerMatches(session, {
                    characterId: ownerCharacterId,
                    channelId: requestChannelId,
                    cashShopId: requestCashShopId,
                })
            ) {
                return [0, SessionErrorCode.SESSION_NOT_OWNER, state];
            }
            await client.del(key);
            return [1, SessionErrorCode.SESSION_NONE, state];
        });
    }

    async getAccountSession(accountId: number): Promise<AccountSession | null> {
        const { client } = this.ctx.getRedisUnifiedAccess();
        const raw = await client.hgetall(this.accountKey(accountId));
        return this.accountHashToSession(raw);
    }

    async getAccountSessionTtl(accountId: number): Promise<number> {
        const { client } = this.ctx.getRedisUnifiedAccess();
        return client.ttl(this.accountKey(accountId));
    }

    private async findCharacterSession(worldId: number, characterId: number): Promise<AccountSession | null> {
        const row = await this.characterRepo.get(worldId, characterId);
        if (row == null) {
            return null;
        }
        const account = await this.getAccountSession(row.accountId);
        if (account == null || account.worldId !== worldId || account.characterId !== characterId) {
            return null;
        }
        return account;
    }

    async findPartyChannel(worldId: number, characterId: number): Promise<number | null> {
        const account = await this.findCharacterSession(worldId, characterId);
        if (account?.state === AS.ACCOUNT_SESSION_STATE_CASH_SHOP) {
            return -1;
        }
        return this.findChannel(worldId, characterId);
    }

    async findChannel(worldId: number, characterId: number): Promise<number | null> {
        const account = await this.findCharacterSession(worldId, characterId);
        if (account == null || account.state !== AS.ACCOUNT_SESSION_STATE_GAME) {
            return null;
        }
        if (account.channelId == null || Number.isInteger(account.channelId) === false || account.channelId < 0) {
            return null;
        }
        return account.channelId;
    }
}
