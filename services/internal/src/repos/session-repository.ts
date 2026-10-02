import { redisSessionKey } from "../redis-session-key";
import {
    AccountSessionState,
    CharacterSessionState,
    accountSessionStateFromRedisHash,
    accountSessionStateToRedisHash,
} from "../session-state";
import type { InternalContext } from "../context/internal-context";
import type { AccountSession, CharacterSession } from "../session-types";

export type { AccountSession, CharacterSession } from "../session-types";

type SessionHash = Record<string, string>;

const AS = AccountSessionState;
const CS = CharacterSessionState;

export class SessionRepository {
    private readonly ctx: InternalContext;

    constructor(internalContext: InternalContext) {
        this.ctx = internalContext;
    }

    private accountKey(accountId: number) {
        return redisSessionKey(`account:${accountId}`);
    }

    private channelOnlineUsersKey(worldId: number, channelId: number) {
        return redisSessionKey(`w${worldId}:ch${channelId}:online_users`);
    }

    private channelOnlineUsersNoopKey(worldId: number) {
        return redisSessionKey(`w${worldId}:ch_:online_users_noop`);
    }

    async touchChannelOnlineUsersTtl(worldId: number, channelId: number) {
        if (!Number.isInteger(channelId) || channelId < 0) {
            return;
        }
        const { client } = this.ctx.getRedisGlobalAccess(worldId);
        const ttl = this.ctx.appConfiguration.getServerAliveTtlSeconds();
        await client.expire(this.channelOnlineUsersKey(worldId, channelId), ttl);
    }

    async getChannelOnlineUserCount(worldId: number, channelId: number): Promise<number> {
        if (!Number.isInteger(channelId) || channelId < 0) {
            return 0;
        }
        const { client } = this.ctx.getRedisGlobalAccess(worldId);
        const raw = await client.get(this.channelOnlineUsersKey(worldId, channelId));
        if (raw == null || raw === "") {
            return 0;
        }
        const n = Number(raw);
        if (!Number.isFinite(n) || n < 0) {
            return 0;
        }
        return Math.floor(n);
    }

    private accountSessionToHash(session: AccountSession): SessionHash {
        return {
            version: String(session.version ?? 1),
            world_id: session.worldId == null ? "" : String(session.worldId),
            account_id: String(session.accountId ?? 0),
            state: accountSessionStateToRedisHash(session.state),
            login_server_id: session.loginServer?.id ?? "",
            login_server_connected: session.loginServer?.connected ? "1" : "0",
            character_id: session.characterId == null ? "" : String(session.characterId),
            character_name: session.characterName ?? "",
            channel_id: session.gameServer?.channelId == null ? "" : String(session.gameServer.channelId),
            game_server_connected: session.gameServer?.connected ? "1" : "0",
            game_to_game_transfer: session.gameToGameTransfer ? "1" : "0",
            created_at: session.timestamps?.createdAt ?? "",
            updated_at: session.timestamps?.updatedAt ?? "",
            state_changed_at: session.timestamps?.stateChangedAt ?? "",
        };
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
            gameServer: {
                channelId: toNumberOrNull(hash.channel_id),
                connected: (hash.game_server_connected || "0") === "1",
            },
            gameToGameTransfer: (hash.game_to_game_transfer || "0") === "1",
            timestamps: {
                createdAt: hash.created_at || null,
                updatedAt: hash.updated_at || null,
                stateChangedAt: hash.state_changed_at || null,
            },
        };
    }

    async beginLoginAtomic(worldId: number, accountId: number, loginServerId: string, now: string, ttlSeconds: number) {
        const { client } = this.ctx.getRedisGlobalAccess(worldId);
        const accountKey = this.accountKey(accountId);
        const script = `
local key = KEYS[1]
local account_id = tonumber(ARGV[1])
local login_server_id = ARGV[2]
local now = ARGV[3]
local ttl = tonumber(ARGV[4])
local ERR_SESSION_NONE = 0
local ERR_SESSION_ALREADY_LOGGED_IN = 2
if redis.call("EXISTS", key) == 1 then
  return {0, ERR_SESSION_ALREADY_LOGGED_IN}
end
redis.call("HSET", key,
  "version", "1",
  "world_id", "",
  "account_id", tostring(account_id),
  "state", "${AS.ACCOUNT_SESSION_STATE_LOGIN}",
  "login_server_id", login_server_id,
  "login_server_connected", "1",
  "created_at", now,
  "updated_at", now,
  "state_changed_at", now
)
redis.call("EXPIRE", key, ttl)
return {1, ERR_SESSION_NONE}
`;
        return (await client.eval(script, 1, accountKey, String(accountId), String(loginServerId), String(now), String(ttlSeconds))) as Array<
            number | string
        >;
    }

    async beginTransitionAtomic(
        worldId: number,
        accountId: number,
        characterId: number,
        characterName: string,
        now: string,
        ttlSeconds: number,
        gameToGameTransfer: boolean,
        prevChannelId: number | null
    ) {
        const { client } = this.ctx.getRedisGlobalAccess(worldId);
        const accountKey = this.accountKey(accountId);
        const prevCh = prevChannelId != null && Number.isInteger(prevChannelId) && prevChannelId >= 0 ? prevChannelId : -1;
        const prevChannelUsersKey =
            prevCh >= 0 ? this.channelOnlineUsersKey(worldId, prevCh) : this.channelOnlineUsersNoopKey(worldId);
        const usersTtl = this.ctx.appConfiguration.getServerAliveTtlSeconds();
        const script = `
local account_key = KEYS[1]
local prev_channel_users_key = KEYS[2]
local world_id = ARGV[1]
local character_id = ARGV[2]
local character_name = ARGV[3]
local now = ARGV[4]
local ttl = tonumber(ARGV[5])
local game_to_game_transfer = ARGV[6]
local prev_channel_id = ARGV[7]
local users_ttl = tonumber(ARGV[8])
local ERR_SESSION_NONE = 0
local ERR_SESSION_NOT_FOUND = 3
if redis.call("EXISTS", account_key) == 0 then
  return {0, ERR_SESSION_NOT_FOUND}
end
local prev_state = tonumber(redis.call("HGET", account_key, "state") or "0")
local connected = redis.call("HGET", account_key, "game_server_connected") or "0"
local channel_id = redis.call("HGET", account_key, "channel_id") or ""
if prev_state == ${AS.ACCOUNT_SESSION_STATE_GAME} and connected == "1" and prev_channel_id ~= "-1" and channel_id == prev_channel_id then
  local v = redis.call("GET", prev_channel_users_key)
  if v and tonumber(v) > 0 then
    redis.call("DECR", prev_channel_users_key)
  end
  if redis.call("EXISTS", prev_channel_users_key) == 1 then
    redis.call("EXPIRE", prev_channel_users_key, users_ttl)
  end
end
redis.call("HSET", account_key,
  "world_id", world_id,
  "state", "${AS.ACCOUNT_SESSION_STATE_TRANSITION}",
  "character_id", character_id,
  "character_name", character_name,
  "game_server_connected", "0",
  "game_to_game_transfer", game_to_game_transfer,
  "updated_at", now,
  "state_changed_at", now
)
redis.call("HDEL", account_key, "channel_id")
redis.call("EXPIRE", account_key, ttl)
return {1, ERR_SESSION_NONE}
`;
        return (await client.eval(
            script,
            2,
            accountKey,
            prevChannelUsersKey,
            String(worldId),
            String(characterId),
            String(characterName),
            String(now),
            String(ttlSeconds),
            gameToGameTransfer ? "1" : "0",
            String(prevCh),
            String(usersTtl)
        )) as Array<number | string>;
    }

    async attachGameSessionAtomic(
        worldId: number,
        accountId: number,
        characterId: number,
        channelId: number,
        now: string,
        ttlSeconds: number
    ) {
        const { client } = this.ctx.getRedisGlobalAccess(worldId);
        const accountKey = this.accountKey(accountId);
        const channelUsersKey = this.channelOnlineUsersKey(worldId, channelId);
        const channelUsersTtl = this.ctx.appConfiguration.getServerAliveTtlSeconds();
        const script = `
local account_key = KEYS[1]
local channel_users_key = KEYS[2]
local expected_account_id = tonumber(ARGV[1])
local expected_world_id = tonumber(ARGV[2])
local expected_character_id = tonumber(ARGV[3])
local channel_id = ARGV[4]
local now = ARGV[5]
local ttl = tonumber(ARGV[6])
local users_ttl = tonumber(ARGV[7])
local ERR_SESSION_NONE = 0
local ERR_SESSION_UNKNOWN = 1
local ERR_SESSION_NOT_FOUND = 3
if redis.call("EXISTS", account_key) == 0 then
  return {0, ERR_SESSION_NOT_FOUND}
end
local state = tonumber(redis.call("HGET", account_key, "state") or "0")
if state ~= ${AS.ACCOUNT_SESSION_STATE_TRANSITION} then
  return {0, ERR_SESSION_UNKNOWN}
end
local account_id = tonumber(redis.call("HGET", account_key, "account_id") or "0")
local world_id = tonumber(redis.call("HGET", account_key, "world_id") or "0")
local character_id = tonumber(redis.call("HGET", account_key, "character_id") or "0")
if account_id ~= expected_account_id or world_id ~= expected_world_id or character_id ~= expected_character_id then
  return {0, ERR_SESSION_UNKNOWN}
end

redis.call("HSET", account_key,
  "state", "${AS.ACCOUNT_SESSION_STATE_GAME}",
  "login_server_connected", "0",
  "channel_id", channel_id,
  "game_server_connected", "1",
  "game_to_game_transfer", "0",
  "updated_at", now,
  "state_changed_at", now
)
redis.call("EXPIRE", account_key, ttl)
redis.call("INCR", channel_users_key)
redis.call("EXPIRE", channel_users_key, users_ttl)
return {1, ERR_SESSION_NONE}
`;
        return (await client.eval(
            script,
            2,
            accountKey,
            channelUsersKey,
            String(accountId),
            String(worldId),
            String(characterId),
            String(channelId),
            String(now),
            String(ttlSeconds),
            String(channelUsersTtl)
        )) as Array<number | string>;
    }

    async refreshAtomic(worldId: number, accountId: number, loginTtl: number, transitionTtl: number, gameTtl: number) {
        const { client } = this.ctx.getRedisGlobalAccess(worldId);
        const accountKey = this.accountKey(accountId);
        const script = `
local account_key = KEYS[1]
local login_ttl = tonumber(ARGV[1])
local transition_ttl = tonumber(ARGV[2])
local game_ttl = tonumber(ARGV[3])
local ERR_SESSION_NONE = 0
local ERR_SESSION_NOT_FOUND = 3
if redis.call("EXISTS", account_key) == 0 then
  return {0, ERR_SESSION_NOT_FOUND, 0, "0"}
end
local state = tonumber(redis.call("HGET", account_key, "state") or "0")
local ttl = login_ttl
if state == ${AS.ACCOUNT_SESSION_STATE_TRANSITION} then
  ttl = transition_ttl
elseif state == ${AS.ACCOUNT_SESSION_STATE_GAME} then
  ttl = game_ttl
end
redis.call("EXPIRE", account_key, ttl)
return {1, ERR_SESSION_NONE, ttl, tostring(state)}
`;
        return (await client.eval(script, 1, accountKey, String(loginTtl), String(transitionTtl), String(gameTtl))) as Array<
            number | string
        >;
    }

    async logoutAtomic(worldId: number, accountId: number, keep: boolean, requestChannelId?: number) {
        const { client } = this.ctx.getRedisGlobalAccess(worldId);
        const accountKey = this.accountKey(accountId);
        const usersTtl = this.ctx.appConfiguration.getServerAliveTtlSeconds();
        const ch =
            requestChannelId != null && Number.isInteger(requestChannelId) && requestChannelId >= 0
                ? requestChannelId
                : -1;
        const channelUsersKey =
            ch >= 0 ? this.channelOnlineUsersKey(worldId, ch) : this.channelOnlineUsersNoopKey(worldId);
        const script = `
local account_key = KEYS[1]
local channel_users_key = KEYS[2]
local keep = ARGV[1]
local channel_id = tonumber(ARGV[2])
local users_ttl = tonumber(ARGV[3])
local ERR_SESSION_NONE = 0
local prev_state = 0
local connected = "0"
if redis.call("EXISTS", account_key) == 1 then
  prev_state = tonumber(redis.call("HGET", account_key, "state") or "0")
  connected = redis.call("HGET", account_key, "game_server_connected") or "0"
end
if keep == "1" then
  return {1, ERR_SESSION_NONE, prev_state}
end
redis.call("DEL", account_key)
if prev_state == ${AS.ACCOUNT_SESSION_STATE_GAME} and connected == "1" and channel_id >= 0 then
  local v = redis.call("GET", channel_users_key)
  if v and tonumber(v) > 0 then
    redis.call("DECR", channel_users_key)
  end
  if redis.call("EXISTS", channel_users_key) == 1 then
    redis.call("EXPIRE", channel_users_key, users_ttl)
  end
end
return {1, ERR_SESSION_NONE, prev_state}
`;
        return (await client.eval(script, 2, accountKey, channelUsersKey, keep ? "1" : "0", String(ch), String(usersTtl))) as Array<
            number | string
        >;
    }

    async getAccountSession(worldId: number, accountId: number): Promise<AccountSession | null> {
        const { client } = this.ctx.getRedisGlobalAccess(worldId);
        const raw = await client.hgetall(this.accountKey(accountId));
        return this.accountHashToSession(raw);
    }

    async getAccountSessionTtl(worldId: number, accountId: number): Promise<number> {
        const { client } = this.ctx.getRedisGlobalAccess(worldId);
        return client.ttl(this.accountKey(accountId));
    }

    async setAccountSession(worldId: number, accountId: number, session: AccountSession, ttlSeconds: number) {
        const { client } = this.ctx.getRedisGlobalAccess(worldId);
        await client.hset(this.accountKey(accountId), this.accountSessionToHash(session));
        await client.expire(this.accountKey(accountId), ttlSeconds);
    }

    async delAccountSession(worldId: number, accountId: number) {
        const { client } = this.ctx.getRedisGlobalAccess(worldId);
        await client.del(this.accountKey(accountId));
    }

    async refreshAccountSession(worldId: number, accountId: number, ttlSeconds: number) {
        const { client } = this.ctx.getRedisGlobalAccess(worldId);
        await client.expire(this.accountKey(accountId), ttlSeconds);
    }

    async getCharacterSession(worldId: number, accountId: number, characterId: number): Promise<CharacterSession | null> {
        const account = await this.getAccountSession(worldId, accountId);
        if (account == null || account.worldId !== worldId || account.characterId !== characterId) {
            return null;
        }
        let state: CharacterSessionState;
        switch (account.state) {
            case AS.ACCOUNT_SESSION_STATE_TRANSITION:
                state = CS.CHARACTER_SESSION_STATE_TRANSITION;
                break;
            case AS.ACCOUNT_SESSION_STATE_GAME:
                state = CS.CHARACTER_SESSION_STATE_ONLINE;
                break;
            default:
                return null;
        }
        return {
            worldId,
            accountId,
            characterId,
            characterName: account.characterName,
            state,
            gameServer: account.gameServer,
        };
    }
}
