import { redisSessionKey } from "../redis-session-key";
import { AccountSessionState, accountSessionStateFromRedisHash } from "../session-state";
import { SessionErrorCode } from "../protobuf/generated/fminternal/internal_service";
import type { InternalContext } from "../context/internal-context";
import type { AccountSession } from "../session-types";
import type { CharacterRepository } from "./character-repository";

export type { AccountSession } from "../session-types";

type SessionHash = Record<string, string>;

const AS = AccountSessionState;

export class SessionRepository {
    private readonly ctx: InternalContext;
    private readonly characterRepo: CharacterRepository;

    constructor(internalContext: InternalContext, characterRepository: CharacterRepository) {
        this.ctx = internalContext;
        this.characterRepo = characterRepository;
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
            gameToGameTransfer: (hash.game_to_game_transfer || "0") === "1",
            timestamps: {
                createdAt: hash.created_at || null,
                updatedAt: hash.updated_at || null,
                stateChangedAt: hash.state_changed_at || null,
            },
        };
    }

    async beginLogin(worldId: number, accountId: number, loginServerId: string, now: string, ttlSeconds: number) {
        const { client } = this.ctx.getRedisGlobalAccess(worldId);
        const accountKey = this.accountKey(accountId);
        const script = `
local key = KEYS[1]
local account_id = tonumber(ARGV[1])
local login_server_id = ARGV[2]
local now = ARGV[3]
local ttl = tonumber(ARGV[4])
if redis.call("EXISTS", key) == 1 then
  return {0, ${SessionErrorCode.SESSION_ALREADY_LOGGED_IN}}
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
return {1, ${SessionErrorCode.SESSION_NONE}}
`;
        return (await client.eval(script, 1, accountKey, String(accountId), String(loginServerId), String(now), String(ttlSeconds))) as Array<
            number | string
        >;
    }

    async beginTransition(
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
if redis.call("EXISTS", account_key) == 0 then
  return {0, ${SessionErrorCode.SESSION_NOT_FOUND}}
end
local prev_state = tonumber(redis.call("HGET", account_key, "state") or "0")
local channel_id = redis.call("HGET", account_key, "channel_id") or ""
if prev_state == ${AS.ACCOUNT_SESSION_STATE_GAME} and prev_channel_id ~= "-1" and channel_id == prev_channel_id then
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
  "game_to_game_transfer", game_to_game_transfer,
  "updated_at", now,
  "state_changed_at", now
)
redis.call("HDEL", account_key, "channel_id")
redis.call("EXPIRE", account_key, ttl)
return {1, ${SessionErrorCode.SESSION_NONE}}
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

    async enterGame(
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
if redis.call("EXISTS", account_key) == 0 then
  return {0, ${SessionErrorCode.SESSION_NOT_FOUND}}
end
local state = tonumber(redis.call("HGET", account_key, "state") or "0")
if state ~= ${AS.ACCOUNT_SESSION_STATE_TRANSITION} then
  return {0, ${SessionErrorCode.SESSION_UNKNOWN}}
end
local account_id = tonumber(redis.call("HGET", account_key, "account_id") or "0")
local world_id = tonumber(redis.call("HGET", account_key, "world_id") or "0")
local character_id = tonumber(redis.call("HGET", account_key, "character_id") or "0")
if account_id ~= expected_account_id or world_id ~= expected_world_id or character_id ~= expected_character_id then
  return {0, ${SessionErrorCode.SESSION_UNKNOWN}}
end

redis.call("HSET", account_key,
  "state", "${AS.ACCOUNT_SESSION_STATE_GAME}",
  "login_server_connected", "0",
  "channel_id", channel_id,
  "game_to_game_transfer", "0",
  "updated_at", now,
  "state_changed_at", now
)
redis.call("EXPIRE", account_key, ttl)
redis.call("INCR", channel_users_key)
redis.call("EXPIRE", channel_users_key, users_ttl)
return {1, ${SessionErrorCode.SESSION_NONE}}
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

    async refresh(worldId: number, accountId: number, loginTtl: number, transitionTtl: number, gameTtl: number) {
        const { client } = this.ctx.getRedisGlobalAccess(worldId);
        const accountKey = this.accountKey(accountId);
        const script = `
local account_key = KEYS[1]
local login_ttl = tonumber(ARGV[1])
local transition_ttl = tonumber(ARGV[2])
local game_ttl = tonumber(ARGV[3])
if redis.call("EXISTS", account_key) == 0 then
  return {0, ${SessionErrorCode.SESSION_NOT_FOUND}, 0, "0"}
end
local state = tonumber(redis.call("HGET", account_key, "state") or "0")
local ttl = login_ttl
if state == ${AS.ACCOUNT_SESSION_STATE_TRANSITION} then
  ttl = transition_ttl
elseif state == ${AS.ACCOUNT_SESSION_STATE_GAME} then
  ttl = game_ttl
end
redis.call("EXPIRE", account_key, ttl)
return {1, ${SessionErrorCode.SESSION_NONE}, ttl, tostring(state)}
`;
        return (await client.eval(script, 1, accountKey, String(loginTtl), String(transitionTtl), String(gameTtl))) as Array<
            number | string
        >;
    }

    async logout(worldId: number, accountId: number, keep: boolean, requestChannelId?: number) {
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
local prev_state = 0
if redis.call("EXISTS", account_key) == 1 then
  prev_state = tonumber(redis.call("HGET", account_key, "state") or "0")
end
if keep == "1" then
  return {1, ${SessionErrorCode.SESSION_NONE}, prev_state}
end
redis.call("DEL", account_key)
if prev_state == ${AS.ACCOUNT_SESSION_STATE_GAME} and channel_id >= 0 then
  local v = redis.call("GET", channel_users_key)
  if v and tonumber(v) > 0 then
    redis.call("DECR", channel_users_key)
  end
  if redis.call("EXISTS", channel_users_key) == 1 then
    redis.call("EXPIRE", channel_users_key, users_ttl)
  end
end
return {1, ${SessionErrorCode.SESSION_NONE}, prev_state}
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

    async findChannel(worldId: number, characterId: number): Promise<number | null> {
        const row = await this.characterRepo.get(worldId, characterId);
        if (row == null) {
            return null;
        }
        const account = await this.getAccountSession(worldId, row.accountId);
        if (account == null || account.state !== AS.ACCOUNT_SESSION_STATE_GAME) {
            return null;
        }
        if (account.worldId !== worldId || account.characterId !== characterId) {
            return null;
        }
        if (account.channelId == null || Number.isInteger(account.channelId) === false || account.channelId < 0) {
            return null;
        }
        return account.channelId;
    }
}
