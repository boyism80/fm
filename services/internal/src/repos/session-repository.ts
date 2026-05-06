import { redisCacheKey } from "../redis-cache-key";
import type { InternalContext } from "../context/internal-context";

type SessionHash = Record<string, string>;

export interface AccountSession {
    version: number;
    worldId: number | null;
    accountId: number;
    state: string;
    character: { id: number | null; name: string | null };
    loginServer: { id: string | null; connected: boolean };
    gameServer: { id: string | null; worldId: number | null; channelId: number | null; connected: boolean };
    timestamps: { createdAt: string | null; updatedAt: string | null; stateChangedAt: string | null };
}

export interface CharacterSession {
    version: number;
    worldId: number;
    accountId: number;
    characterId: number;
    characterName: string | null;
    state: string;
    gameServer: { id: string | null; worldId: number | null; channelId: number | null; connected: boolean };
    timestamps: { createdAt: string | null; updatedAt: string | null };
}

export class SessionRepository {
    private readonly ctx: InternalContext;

    constructor(internalContext: InternalContext) {
        this.ctx = internalContext;
    }

    private accountKey(accountId: number) {
        return redisCacheKey(`session:account:${accountId}`);
    }

    private characterKey(worldId: number, characterId: number) {
        return redisCacheKey(`w${worldId}:session:character:${characterId}`);
    }

    private accountSessionToHash(session: AccountSession): SessionHash {
        return {
            version: String(session.version ?? 1),
            world_id: session.worldId == null ? "" : String(session.worldId),
            account_id: String(session.accountId ?? 0),
            state: session.state ?? "",
            character_id: session.character?.id == null ? "" : String(session.character.id),
            character_name: session.character?.name ?? "",
            login_server_id: session.loginServer?.id ?? "",
            login_server_connected: session.loginServer?.connected ? "1" : "0",
            game_server_id: session.gameServer?.id ?? "",
            game_server_world_id: session.gameServer?.worldId == null ? "" : String(session.gameServer.worldId),
            game_server_channel_id: session.gameServer?.channelId == null ? "" : String(session.gameServer.channelId),
            game_server_connected: session.gameServer?.connected ? "1" : "0",
            created_at: session.timestamps?.createdAt ?? "",
            updated_at: session.timestamps?.updatedAt ?? "",
            state_changed_at: session.timestamps?.stateChangedAt ?? "",
        };
    }

    private accountHashToSession(hash: SessionHash): AccountSession | null {
        if (!hash || Object.keys(hash).length === 0) {
            return null;
        }
        const toNumberOrNull = (v: unknown) => (v === "" || v == null ? null : Number(v));
        return {
            version: Number(hash.version ?? 1),
            worldId: toNumberOrNull(hash.world_id),
            accountId: Number(hash.account_id ?? 0),
            state: hash.state ?? "",
            character: {
                id: toNumberOrNull(hash.character_id),
                name: hash.character_name || null,
            },
            loginServer: {
                id: hash.login_server_id || null,
                connected: (hash.login_server_connected || "0") === "1",
            },
            gameServer: {
                id: hash.game_server_id || null,
                worldId: toNumberOrNull(hash.game_server_world_id),
                channelId: toNumberOrNull(hash.game_server_channel_id),
                connected: (hash.game_server_connected || "0") === "1",
            },
            timestamps: {
                createdAt: hash.created_at || null,
                updatedAt: hash.updated_at || null,
                stateChangedAt: hash.state_changed_at || null,
            },
        };
    }

    private characterSessionToHash(session: CharacterSession): SessionHash {
        return {
            version: String(session.version ?? 1),
            world_id: String(session.worldId ?? 0),
            account_id: String(session.accountId ?? 0),
            character_id: String(session.characterId ?? 0),
            character_name: session.characterName ?? "",
            state: session.state ?? "",
            game_server_id: session.gameServer?.id ?? "",
            game_server_world_id: session.gameServer?.worldId == null ? "" : String(session.gameServer.worldId),
            game_server_channel_id: session.gameServer?.channelId == null ? "" : String(session.gameServer.channelId),
            game_server_connected: session.gameServer?.connected ? "1" : "0",
            created_at: session.timestamps?.createdAt ?? "",
            updated_at: session.timestamps?.updatedAt ?? "",
        };
    }

    private characterHashToSession(hash: SessionHash): CharacterSession | null {
        if (!hash || Object.keys(hash).length === 0) {
            return null;
        }
        const toNumberOrNull = (v: unknown) => (v === "" || v == null ? null : Number(v));
        return {
            version: Number(hash.version ?? 1),
            worldId: Number(hash.world_id ?? 0),
            accountId: Number(hash.account_id ?? 0),
            characterId: Number(hash.character_id ?? 0),
            characterName: hash.character_name || null,
            state: hash.state || "",
            gameServer: {
                id: hash.game_server_id || null,
                worldId: toNumberOrNull(hash.game_server_world_id),
                channelId: toNumberOrNull(hash.game_server_channel_id),
                connected: (hash.game_server_connected || "0") === "1",
            },
            timestamps: {
                createdAt: hash.created_at || null,
                updatedAt: hash.updated_at || null,
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
  "state", "LOGIN",
  "character_id", "",
  "character_name", "",
  "login_server_id", login_server_id,
  "login_server_connected", "1",
  "game_server_id", "",
  "game_server_world_id", "",
  "game_server_channel_id", "",
  "game_server_connected", "0",
  "created_at", now,
  "updated_at", now,
  "state_changed_at", now
)
redis.call("EXPIRE", key, ttl)
return {1, ERR_SESSION_NONE}
`;
        return client.eval(script, 1, accountKey, String(accountId), String(loginServerId), String(now), String(ttlSeconds));
    }

    async beginTransitionAtomic(worldId: number, accountId: number, characterId: number, characterName: string, now: string, ttlSeconds: number) {
        const { client } = this.ctx.getRedisGlobalAccess(worldId);
        const accountKey = this.accountKey(accountId);
        const script = `
local key = KEYS[1]
local world_id = tonumber(ARGV[1])
local character_id = tonumber(ARGV[2])
local character_name = ARGV[3]
local now = ARGV[4]
local ttl = tonumber(ARGV[5])
local ERR_SESSION_NONE = 0
local ERR_SESSION_UNKNOWN = 1
local ERR_SESSION_NOT_FOUND = 3
if redis.call("EXISTS", key) == 0 then
  return {0, ERR_SESSION_NOT_FOUND}
end
local state = redis.call("HGET", key, "state")
if state ~= "LOGIN" and state ~= "TRANSITION" then
  return {0, ERR_SESSION_UNKNOWN}
end
redis.call("HSET", key,
  "world_id", tostring(world_id),
  "state", "TRANSITION",
  "character_id", tostring(character_id),
  "character_name", character_name,
  "updated_at", now,
  "state_changed_at", now
)
redis.call("EXPIRE", key, ttl)
return {1, ERR_SESSION_NONE}
`;
        return client.eval(script, 1, accountKey, String(worldId), String(characterId), String(characterName), String(now), String(ttlSeconds));
    }

    async attachGameSessionAtomic(worldId: number, accountId: number, channelId: number, now: string, ttlSeconds: number) {
        const { client } = this.ctx.getRedisGlobalAccess(worldId);
        const accountKey = this.accountKey(accountId);
        const characterKey = this.characterKey(worldId, 0);
        const script = `
local account_key = KEYS[1]
local character_key_prefix = KEYS[2]
local channel_id = tonumber(ARGV[1])
local now = ARGV[2]
local ttl = tonumber(ARGV[3])
local ERR_SESSION_NONE = 0
local ERR_SESSION_UNKNOWN = 1
local ERR_SESSION_NOT_FOUND = 3
if redis.call("EXISTS", account_key) == 0 then
  return {0, ERR_SESSION_NOT_FOUND}
end
local state = redis.call("HGET", account_key, "state")
if state ~= "TRANSITION" and state ~= "GAME" then
  return {0, ERR_SESSION_UNKNOWN}
end
local character_id = tonumber(redis.call("HGET", account_key, "character_id"))
local character_name = redis.call("HGET", account_key, "character_name") or ""
if not character_id or character_id == 0 then
  return {0, ERR_SESSION_UNKNOWN}
end
local world_id = tonumber(redis.call("HGET", account_key, "world_id") or "0")
local account_id = tonumber(redis.call("HGET", account_key, "account_id") or "0")
local game_server_id = "w" .. tostring(world_id) .. ":ch" .. tostring(channel_id)
local character_key = string.gsub(character_key_prefix, ":0$", ":" .. tostring(character_id))

redis.call("HSET", account_key,
  "state", "GAME",
  "game_server_id", game_server_id,
  "game_server_world_id", tostring(world_id),
  "game_server_channel_id", tostring(channel_id),
  "game_server_connected", "1",
  "login_server_connected", "0",
  "updated_at", now,
  "state_changed_at", now
)
redis.call("EXPIRE", account_key, ttl)

redis.call("HSET", character_key,
  "version", "1",
  "world_id", tostring(world_id),
  "account_id", tostring(account_id),
  "character_id", tostring(character_id),
  "character_name", character_name,
  "state", "ONLINE",
  "game_server_id", game_server_id,
  "game_server_world_id", tostring(world_id),
  "game_server_channel_id", tostring(channel_id),
  "game_server_connected", "1",
  "created_at", now,
  "updated_at", now
)
redis.call("EXPIRE", character_key, ttl)
return {1, ERR_SESSION_NONE}
`;
        return client.eval(script, 2, accountKey, characterKey, String(channelId), String(now), String(ttlSeconds));
    }

    async refreshAtomic(worldId: number, accountId: number, loginTtl: number, transitionTtl: number, gameTtl: number) {
        const { client } = this.ctx.getRedisGlobalAccess(worldId);
        const accountKey = this.accountKey(accountId);
        const script = `
local account_key = KEYS[1]
local character_key_prefix = KEYS[2]
local login_ttl = tonumber(ARGV[1])
local transition_ttl = tonumber(ARGV[2])
local game_ttl = tonumber(ARGV[3])
local ERR_SESSION_NONE = 0
local ERR_SESSION_NOT_FOUND = 3
if redis.call("EXISTS", account_key) == 0 then
  return {0, ERR_SESSION_NOT_FOUND}
end
local state = redis.call("HGET", account_key, "state")
local ttl = login_ttl
if state == "TRANSITION" then
  ttl = transition_ttl
elseif state == "GAME" then
  ttl = game_ttl
end
redis.call("EXPIRE", account_key, ttl)
if state == "GAME" then
  local character_id = redis.call("HGET", account_key, "character_id")
  if character_id and character_id ~= "" then
    local character_key = string.gsub(character_key_prefix, ":0$", ":" .. tostring(character_id))
    redis.call("EXPIRE", character_key, ttl)
  end
end
return {1, ERR_SESSION_NONE}
`;
        return client.eval(script, 2, accountKey, this.characterKey(worldId, 0), String(loginTtl), String(transitionTtl), String(gameTtl));
    }

    async logoutAtomic(worldId: number, accountId: number, transitionTtl: number) {
        const { client } = this.ctx.getRedisGlobalAccess(worldId);
        const accountKey = this.accountKey(accountId);
        const script = `
local account_key = KEYS[1]
local character_key_prefix = KEYS[2]
local transition_ttl = tonumber(ARGV[1])
local ERR_SESSION_NONE = 0
if redis.call("EXISTS", account_key) == 1 then
  local state = redis.call("HGET", account_key, "state")
  if state == "TRANSITION" then
    redis.call("EXPIRE", account_key, transition_ttl)
    return {1, ERR_SESSION_NONE}
  end
  local character_id = redis.call("HGET", account_key, "character_id")
  if character_id and character_id ~= "" then
    local character_key = string.gsub(character_key_prefix, ":0$", ":" .. tostring(character_id))
    redis.call("DEL", character_key)
  end
end
redis.call("DEL", account_key)
return {1, ERR_SESSION_NONE}
`;
        return client.eval(script, 2, accountKey, this.characterKey(worldId, 0), String(transitionTtl));
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

    async getCharacterSession(worldId: number, characterId: number): Promise<CharacterSession | null> {
        const { client } = this.ctx.getRedisGlobalAccess(worldId);
        const raw = await client.hgetall(this.characterKey(worldId, characterId));
        return this.characterHashToSession(raw);
    }

    async setCharacterSession(worldId: number, characterId: number, session: CharacterSession, ttlSeconds: number) {
        const { client } = this.ctx.getRedisGlobalAccess(worldId);
        const multi = client.multi();
        multi.hset(this.characterKey(worldId, characterId), this.characterSessionToHash(session));
        multi.expire(this.characterKey(worldId, characterId), ttlSeconds);
        await multi.exec();
    }

    async delCharacterSession(worldId: number, characterId: number) {
        const { client } = this.ctx.getRedisGlobalAccess(worldId);
        await client.del(this.characterKey(worldId, characterId));
    }

    async refreshCharacterSession(worldId: number, characterId: number, ttlSeconds: number) {
        const { client } = this.ctx.getRedisGlobalAccess(worldId);
        await client.expire(this.characterKey(worldId, characterId), ttlSeconds);
    }
}
