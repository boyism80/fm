import { deserializeCharacterSession, serializeCharacterSession } from "../character-session-json";
import { redisSessionKey } from "../redis-session-key";
import {
    AccountSessionState,
    CharacterSessionState,
    accountSessionStateFromRedisHash,
    accountSessionStateToRedisHash,
} from "../session-state";
import { isRedisLogoutEvalArray, isRedisRefreshEvalArray } from "../types/redis-eval";
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

    private worldCharacterSessionsKey(worldId: number) {
        return redisSessionKey(`w${worldId}`);
    }

    private transitionCharacterPointerKey(accountId: number) {
        return redisSessionKey(`account:${accountId}:transition_character_id`);
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

    async beginTransitionAtomic(worldId: number, accountId: number, characterId: number, characterName: string, now: string, ttlSeconds: number) {
        const { client } = this.ctx.getRedisGlobalAccess(worldId);
        const accountKey = this.accountKey(accountId);
        const worldKey = this.worldCharacterSessionsKey(worldId);
        const pointerKey = this.transitionCharacterPointerKey(accountId);
        const transitionSessionJson = serializeCharacterSession({
            version: 1,
            worldId,
            accountId,
            characterId,
            characterName,
            state: CS.CHARACTER_SESSION_STATE_TRANSITION,
            gameServer: { id: null, worldId: null, channelId: null, connected: false },
            timestamps: { createdAt: now, updatedAt: now },
        });
        const script = `
local account_key = KEYS[1]
local world_key = KEYS[2]
local pointer_key = KEYS[3]
local world_id = tonumber(ARGV[1])
local account_id = tonumber(ARGV[2])
local character_id = tonumber(ARGV[3])
local character_name = ARGV[4]
local now = ARGV[5]
local ttl = tonumber(ARGV[6])
local character_session_json = ARGV[7]
local ERR_SESSION_NONE = 0
local ERR_SESSION_UNKNOWN = 1
local ERR_SESSION_NOT_FOUND = 3
if redis.call("EXISTS", account_key) == 0 then
  return {0, ERR_SESSION_NOT_FOUND}
end
redis.call("HSET", account_key,
  "world_id", tostring(world_id),
  "state", "${AS.ACCOUNT_SESSION_STATE_TRANSITION}",
  "updated_at", now,
  "state_changed_at", now
)
redis.call("HDEL", account_key, "character_id", "character_name")
redis.call("EXPIRE", account_key, ttl)
redis.call("HSET", world_key, character_name, character_session_json)
redis.call("EXPIRE", world_key, ttl)
redis.call("SET", pointer_key, tostring(character_id), "EX", ttl)
return {1, ERR_SESSION_NONE}
`;
        return (await client.eval(
            script,
            3,
            accountKey,
            worldKey,
            pointerKey,
            String(worldId),
            String(accountId),
            String(characterId),
            String(characterName),
            String(now),
            String(ttlSeconds),
            transitionSessionJson
        )) as Array<number | string>;
    }

    async attachGameSessionAtomic(
        worldId: number,
        accountId: number,
        characterId: number,
        characterName: string,
        channelId: number,
        now: string,
        ttlSeconds: number
    ) {
        const { client } = this.ctx.getRedisGlobalAccess(worldId);
        const accountKey = this.accountKey(accountId);
        const worldKey = this.worldCharacterSessionsKey(worldId);
        const pointerKey = this.transitionCharacterPointerKey(accountId);
        const existingRaw = await client.hget(worldKey, characterName);
        const prev = existingRaw ? deserializeCharacterSession(existingRaw) : null;
        const prevCreated = prev?.timestamps.createdAt ?? now;
        const gameServerId = `w${worldId}:ch${channelId}`;
        const onlineSessionJson = serializeCharacterSession({
            version: 1,
            worldId,
            accountId,
            characterId,
            characterName,
            state: CS.CHARACTER_SESSION_STATE_ONLINE,
            gameServer: {
                id: gameServerId,
                worldId,
                channelId,
                connected: true,
            },
            timestamps: { createdAt: prevCreated, updatedAt: now },
        });
        const newChKey = this.channelOnlineUsersKey(worldId, channelId);
        let oldChKey = newChKey;
        let incrF = "0";
        let decrF = "0";
        let touchTtlF = "0";
        if (
            prev?.state === CS.CHARACTER_SESSION_STATE_ONLINE &&
            prev.gameServer?.connected === true &&
            prev.gameServer.channelId === channelId
        ) {
            touchTtlF = "1";
        } else if (
            prev?.state === CS.CHARACTER_SESSION_STATE_ONLINE &&
            prev.gameServer?.connected === true &&
            prev.gameServer.channelId != null &&
            prev.gameServer.channelId !== channelId
        ) {
            incrF = "1";
            decrF = "1";
            oldChKey = this.channelOnlineUsersKey(worldId, prev.gameServer.channelId);
        } else {
            incrF = "1";
        }
        const channelUsersTtl = this.ctx.appConfiguration.getServerAliveTtlSeconds();
        const script = `
local account_key = KEYS[1]
local world_key = KEYS[2]
local pointer_key = KEYS[3]
local new_ch_key = KEYS[4]
local old_ch_key = KEYS[5]
local character_name = ARGV[1]
local now = ARGV[2]
local ttl = tonumber(ARGV[3])
local expected_account_id = tonumber(ARGV[4])
local expected_world_id = tonumber(ARGV[5])
local character_session_json = ARGV[6]
local incr_f = ARGV[7]
local decr_f = ARGV[8]
local touch_ttl_f = ARGV[9]
local users_ttl = tonumber(ARGV[10])
local ERR_SESSION_NONE = 0
local ERR_SESSION_UNKNOWN = 1
local ERR_SESSION_NOT_FOUND = 3
if redis.call("EXISTS", account_key) == 0 then
  return {0, ERR_SESSION_NOT_FOUND}
end
if redis.call("HEXISTS", world_key, character_name) ~= 1 then
  return {0, ERR_SESSION_UNKNOWN}
end
local state = tonumber(redis.call("HGET", account_key, "state") or "0")
if state ~= ${AS.ACCOUNT_SESSION_STATE_TRANSITION} and state ~= ${AS.ACCOUNT_SESSION_STATE_GAME} then
  return {0, ERR_SESSION_UNKNOWN}
end
local acc_acc = tonumber(redis.call("HGET", account_key, "account_id") or "0")
local world_id = tonumber(redis.call("HGET", account_key, "world_id") or "0")
if acc_acc ~= expected_account_id or world_id ~= expected_world_id then
  return {0, ERR_SESSION_UNKNOWN}
end

redis.call("HSET", account_key,
  "state", "${AS.ACCOUNT_SESSION_STATE_GAME}",
  "login_server_connected", "0",
  "updated_at", now,
  "state_changed_at", now
)
redis.call("EXPIRE", account_key, ttl)
redis.call("HSET", world_key, character_name, character_session_json)
redis.call("EXPIRE", world_key, ttl)
redis.call("DEL", pointer_key)
if incr_f == "1" then
  redis.call("INCR", new_ch_key)
end
if incr_f == "1" or touch_ttl_f == "1" then
  if redis.call("EXISTS", new_ch_key) == 1 then
    redis.call("EXPIRE", new_ch_key, users_ttl)
  end
end
if decr_f == "1" and old_ch_key ~= new_ch_key then
  local v = redis.call("GET", old_ch_key)
  if v and tonumber(v) > 0 then
    redis.call("DECR", old_ch_key)
  end
  if redis.call("EXISTS", old_ch_key) == 1 then
    redis.call("EXPIRE", old_ch_key, users_ttl)
  end
end
return {1, ERR_SESSION_NONE}
`;
        return (await client.eval(
            script,
            5,
            accountKey,
            worldKey,
            pointerKey,
            newChKey,
            oldChKey,
            String(characterName),
            String(now),
            String(ttlSeconds),
            String(accountId),
            String(worldId),
            onlineSessionJson,
            incrF,
            decrF,
            touchTtlF,
            String(channelUsersTtl)
        )) as Array<number | string>;
    }

    async refreshAtomic(
        worldId: number,
        accountId: number,
        loginTtl: number,
        transitionTtl: number,
        gameTtl: number,
        requestCharacterId?: number
    ) {
        const { client } = this.ctx.getRedisGlobalAccess(worldId);
        const accountKey = this.accountKey(accountId);
        const pointerKey = this.transitionCharacterPointerKey(accountId);
        const worldKey = this.worldCharacterSessionsKey(worldId);
        const script = `
local account_key = KEYS[1]
local pointer_key = KEYS[2]
local login_ttl = tonumber(ARGV[1])
local transition_ttl = tonumber(ARGV[2])
local game_ttl = tonumber(ARGV[3])
local request_character_id = ARGV[4] or ""
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
if state == ${AS.ACCOUNT_SESSION_STATE_TRANSITION} and (not request_character_id or request_character_id == "") then
  if redis.call("EXISTS", pointer_key) == 1 then
    redis.call("EXPIRE", pointer_key, ttl)
  end
end
return {1, ERR_SESSION_NONE, ttl, tostring(state)}
`;
        const reqCid =
            requestCharacterId != null && requestCharacterId > 0 ? String(requestCharacterId) : "";
        const raw = await client.eval(
            script,
            2,
            accountKey,
            pointerKey,
            String(loginTtl),
            String(transitionTtl),
            String(gameTtl),
            reqCid
        );
        if (!Array.isArray(raw) || !isRedisRefreshEvalArray(raw) || Number(raw[0]) !== 1) {
            return raw as Array<number | string>;
        }
        const ttl = Number(raw[2]);
        const state = Number(raw[3]);
        if (
            state === AS.ACCOUNT_SESSION_STATE_TRANSITION ||
            state === AS.ACCOUNT_SESSION_STATE_GAME
        ) {
            await client.expire(worldKey, ttl);
        }
        return raw as Array<number | string>;
    }

    async logoutAtomic(
        worldId: number,
        accountId: number,
        transitionTtl: number,
        requestCharacterName?: string,
        requestChannelId?: number
    ) {
        const { client } = this.ctx.getRedisGlobalAccess(worldId);
        const accountKey = this.accountKey(accountId);
        const pointerKey = this.transitionCharacterPointerKey(accountId);
        const worldKey = this.worldCharacterSessionsKey(worldId);
        const usersTtl = this.ctx.appConfiguration.getServerAliveTtlSeconds();
        const ch =
            requestChannelId != null && Number.isInteger(requestChannelId) && requestChannelId >= 0
                ? requestChannelId
                : -1;
        const channelUsersKey =
            ch >= 0 ? this.channelOnlineUsersKey(worldId, ch) : this.channelOnlineUsersNoopKey(worldId);
        const charName = requestCharacterName ?? "";
        const script = `
local account_key = KEYS[1]
local pointer_key = KEYS[2]
local world_key = KEYS[3]
local channel_users_key = KEYS[4]
local transition_ttl = tonumber(ARGV[1])
local char_name = ARGV[2]
local channel_id = tonumber(ARGV[3])
local users_ttl = tonumber(ARGV[4])
local ERR_SESSION_NONE = 0
local prev_state = 0
if redis.call("EXISTS", account_key) == 1 then
  prev_state = tonumber(redis.call("HGET", account_key, "state") or "0")
  if prev_state == ${AS.ACCOUNT_SESSION_STATE_TRANSITION} then
    redis.call("DEL", pointer_key)
    redis.call("EXPIRE", account_key, transition_ttl)
    return {1, ERR_SESSION_NONE, prev_state}
  end
end
redis.call("DEL", pointer_key)
redis.call("DEL", account_key)
if prev_state == ${AS.ACCOUNT_SESSION_STATE_GAME} and char_name ~= "" then
  redis.call("HDEL", world_key, char_name)
end
if prev_state == ${AS.ACCOUNT_SESSION_STATE_GAME} and channel_id >= 0 then
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
        const raw = await client.eval(
            script,
            4,
            accountKey,
            pointerKey,
            worldKey,
            channelUsersKey,
            String(transitionTtl),
            charName,
            String(ch),
            String(usersTtl)
        );
        return raw as Array<number | string>;
    }

    async getAccountSession(worldId: number, accountId: number): Promise<AccountSession | null> {
        const { client } = this.ctx.getRedisGlobalAccess(worldId);
        const raw = await client.hgetall(this.accountKey(accountId));
        return this.accountHashToSession(raw);
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

    async getCharacterSessionByName(worldId: number, characterName: string): Promise<CharacterSession | null> {
        const { client } = this.ctx.getRedisGlobalAccess(worldId);
        const raw = await client.hget(this.worldCharacterSessionsKey(worldId), characterName);
        if (!raw) {
            return null;
        }
        return deserializeCharacterSession(raw);
    }

    async setCharacterSession(worldId: number, _characterId: number, session: CharacterSession, ttlSeconds: number) {
        const { client } = this.ctx.getRedisGlobalAccess(worldId);
        const name = session.characterName ?? "";
        if (!name) {
            return;
        }
        const key = this.worldCharacterSessionsKey(worldId);
        const sessionJson = serializeCharacterSession(session);
        await client.hset(key, name, sessionJson);
        await client.expire(key, ttlSeconds);
    }

    async delCharacterSessionByName(worldId: number, characterName: string) {
        const { client } = this.ctx.getRedisGlobalAccess(worldId);
        await client.hdel(this.worldCharacterSessionsKey(worldId), characterName);
    }

    async refreshCharacterSessionByName(worldId: number, characterName: string, ttlSeconds: number) {
        const { client } = this.ctx.getRedisGlobalAccess(worldId);
        await client.expire(this.worldCharacterSessionsKey(worldId), ttlSeconds);
    }
}
