"use strict";

class SessionRepository {
    constructor(internalContext) {
        this.ctx = internalContext;
    }

    _accountKey(worldId, accountId) {
        const { keyPrefix } = this.ctx.getRedisGlobalAccess(worldId);
        return `${keyPrefix}fm:w${worldId}:session:account:${accountId}`;
    }

    _characterKey(worldId, characterId) {
        const { keyPrefix } = this.ctx.getRedisGlobalAccess(worldId);
        return `${keyPrefix}fm:w${worldId}:session:character:${characterId}`;
    }

    _characterToAccountKey(worldId, characterId) {
        const { keyPrefix } = this.ctx.getRedisGlobalAccess(worldId);
        return `${keyPrefix}fm:w${worldId}:session:index:character_to_account:${characterId}`;
    }

    _accountSessionToHash(session) {
        return {
            version: String(session.version ?? 1),
            world_id: String(session.worldId ?? 0),
            account_id: String(session.accountId ?? 0),
            state: String(session.state ?? ""),
            character_id: session.character?.id == null ? "" : String(session.character.id),
            character_name: session.character?.name == null ? "" : String(session.character.name),
            login_server_id: session.loginServer?.id == null ? "" : String(session.loginServer.id),
            login_server_connected: session.loginServer?.connected ? "1" : "0",
            game_server_id: session.gameServer?.id == null ? "" : String(session.gameServer.id),
            game_server_world_id: session.gameServer?.worldId == null ? "" : String(session.gameServer.worldId),
            game_server_channel_id: session.gameServer?.channelId == null ? "" : String(session.gameServer.channelId),
            game_server_connected: session.gameServer?.connected ? "1" : "0",
            created_at: session.timestamps?.createdAt == null ? "" : String(session.timestamps.createdAt),
            updated_at: session.timestamps?.updatedAt == null ? "" : String(session.timestamps.updatedAt),
            state_changed_at: session.timestamps?.stateChangedAt == null ? "" : String(session.timestamps.stateChangedAt),
        };
    }

    _accountHashToSession(hash) {
        if (!hash || Object.keys(hash).length === 0) {
            return null;
        }
        const toNumberOrNull = (v) => (v === "" || v == null ? null : Number(v));
        return {
            version: Number(hash.version ?? 1),
            worldId: Number(hash.world_id ?? 0),
            accountId: Number(hash.account_id ?? 0),
            state: hash.state ?? "",
            character: {
                id: toNumberOrNull(hash.character_id),
                name: hash.character_name || null,
            },
            loginServer: {
                id: hash.login_server_id || null,
                connected: String(hash.login_server_connected || "0") === "1",
            },
            gameServer: {
                id: hash.game_server_id || null,
                worldId: toNumberOrNull(hash.game_server_world_id),
                channelId: toNumberOrNull(hash.game_server_channel_id),
                connected: String(hash.game_server_connected || "0") === "1",
            },
            timestamps: {
                createdAt: hash.created_at || null,
                updatedAt: hash.updated_at || null,
                stateChangedAt: hash.state_changed_at || null,
            },
        };
    }

    _characterSessionToHash(session) {
        return {
            version: String(session.version ?? 1),
            world_id: String(session.worldId ?? 0),
            account_id: String(session.accountId ?? 0),
            character_id: String(session.characterId ?? 0),
            character_name: session.characterName == null ? "" : String(session.characterName),
            state: session.state == null ? "" : String(session.state),
            game_server_id: session.gameServer?.id == null ? "" : String(session.gameServer.id),
            game_server_world_id: session.gameServer?.worldId == null ? "" : String(session.gameServer.worldId),
            game_server_channel_id: session.gameServer?.channelId == null ? "" : String(session.gameServer.channelId),
            game_server_connected: session.gameServer?.connected ? "1" : "0",
            created_at: session.timestamps?.createdAt == null ? "" : String(session.timestamps.createdAt),
            updated_at: session.timestamps?.updatedAt == null ? "" : String(session.timestamps.updatedAt),
        };
    }

    _characterHashToSession(hash) {
        if (!hash || Object.keys(hash).length === 0) {
            return null;
        }
        const toNumberOrNull = (v) => (v === "" || v == null ? null : Number(v));
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
                connected: String(hash.game_server_connected || "0") === "1",
            },
            timestamps: {
                createdAt: hash.created_at || null,
                updatedAt: hash.updated_at || null,
            },
        };
    }

    async beginLoginAtomic(worldId, accountId, loginServerId, now, ttlSeconds) {
        const { client } = this.ctx.getRedisGlobalAccess(worldId);
        const accountKey = this._accountKey(worldId, accountId);
        const script = `
local key = KEYS[1]
local world_id = tonumber(ARGV[1])
local account_id = tonumber(ARGV[2])
local login_server_id = ARGV[3]
local now = ARGV[4]
local ttl = tonumber(ARGV[5])
if redis.call("EXISTS", key) == 1 then
  return {0, "ALREADY_LOGGED_IN"}
end
redis.call("HSET", key,
  "version", "1",
  "world_id", tostring(world_id),
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
return {1, ""}
`;
        return client.eval(
            script,
            1,
            accountKey,
            String(worldId),
            String(accountId),
            String(loginServerId),
            String(now),
            String(ttlSeconds)
        );
    }

    async beginTransitionAtomic(worldId, accountId, characterId, characterName, now, ttlSeconds) {
        const { client } = this.ctx.getRedisGlobalAccess(worldId);
        const accountKey = this._accountKey(worldId, accountId);
        const script = `
local key = KEYS[1]
local character_id = tonumber(ARGV[1])
local character_name = ARGV[2]
local now = ARGV[3]
local ttl = tonumber(ARGV[4])
if redis.call("EXISTS", key) == 0 then
  return {0, "SESSION_NOT_FOUND"}
end
local state = redis.call("HGET", key, "state")
if state ~= "LOGIN" and state ~= "TRANSITION" then
  return {0, "INVALID_STATE"}
end
redis.call("HSET", key,
  "state", "TRANSITION",
  "character_id", tostring(character_id),
  "character_name", character_name,
  "updated_at", now,
  "state_changed_at", now
)
redis.call("EXPIRE", key, ttl)
return {1, ""}
`;
        return client.eval(
            script,
            1,
            accountKey,
            String(characterId),
            String(characterName),
            String(now),
            String(ttlSeconds)
        );
    }

    async attachGameSessionAtomic(worldId, accountId, channelId, now, ttlSeconds) {
        const { client } = this.ctx.getRedisGlobalAccess(worldId);
        const accountKey = this._accountKey(worldId, accountId);
        const characterKey = this._characterKey(worldId, 0);
        const indexKey = this._characterToAccountKey(worldId, 0);
        const script = `
local account_key = KEYS[1]
local character_key_prefix = KEYS[2]
local index_key_prefix = KEYS[3]
local channel_id = tonumber(ARGV[1])
local now = ARGV[2]
local ttl = tonumber(ARGV[3])
if redis.call("EXISTS", account_key) == 0 then
  return {0, "SESSION_NOT_FOUND"}
end
local state = redis.call("HGET", account_key, "state")
if state ~= "TRANSITION" and state ~= "GAME" then
  return {0, "INVALID_STATE"}
end
local character_id = tonumber(redis.call("HGET", account_key, "character_id"))
local character_name = redis.call("HGET", account_key, "character_name") or ""
if not character_id or character_id == 0 then
  return {0, "CHARACTER_NOT_SELECTED"}
end
local world_id = tonumber(redis.call("HGET", account_key, "world_id") or "0")
local account_id = tonumber(redis.call("HGET", account_key, "account_id") or "0")
local game_server_id = "w" .. tostring(world_id) .. ":ch" .. tostring(channel_id)
local character_key = string.gsub(character_key_prefix, ":0$", ":" .. tostring(character_id))
local index_key = string.gsub(index_key_prefix, ":0$", ":" .. tostring(character_id))

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
redis.call("SET", index_key, tostring(account_id), "EX", ttl)
return {1, ""}
`;
        return client.eval(
            script,
            3,
            accountKey,
            characterKey,
            indexKey,
            String(channelId),
            String(now),
            String(ttlSeconds)
        );
    }

    async refreshAtomic(worldId, accountId, loginTtl, transitionTtl, gameTtl) {
        const { client } = this.ctx.getRedisGlobalAccess(worldId);
        const accountKey = this._accountKey(worldId, accountId);
        const script = `
local account_key = KEYS[1]
local character_key_prefix = KEYS[2]
local index_key_prefix = KEYS[3]
local login_ttl = tonumber(ARGV[1])
local transition_ttl = tonumber(ARGV[2])
local game_ttl = tonumber(ARGV[3])
if redis.call("EXISTS", account_key) == 0 then
  return {0, "SESSION_NOT_FOUND"}
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
    local index_key = string.gsub(index_key_prefix, ":0$", ":" .. tostring(character_id))
    redis.call("EXPIRE", character_key, ttl)
    redis.call("EXPIRE", index_key, ttl)
  end
end
return {1, ""}
`;
        return client.eval(
            script,
            3,
            accountKey,
            this._characterKey(worldId, 0),
            this._characterToAccountKey(worldId, 0),
            String(loginTtl),
            String(transitionTtl),
            String(gameTtl)
        );
    }

    async logoutAtomic(worldId, accountId, transitionTtl) {
        const { client } = this.ctx.getRedisGlobalAccess(worldId);
        const accountKey = this._accountKey(worldId, accountId);
        const script = `
local account_key = KEYS[1]
local character_key_prefix = KEYS[2]
local index_key_prefix = KEYS[3]
local transition_ttl = tonumber(ARGV[1])
if redis.call("EXISTS", account_key) == 1 then
  local state = redis.call("HGET", account_key, "state")
  if state == "TRANSITION" then
    redis.call("EXPIRE", account_key, transition_ttl)
    return {1, "DEFERRED_TRANSITION"}
  end
  local character_id = redis.call("HGET", account_key, "character_id")
  if character_id and character_id ~= "" then
    local character_key = string.gsub(character_key_prefix, ":0$", ":" .. tostring(character_id))
    local index_key = string.gsub(index_key_prefix, ":0$", ":" .. tostring(character_id))
    redis.call("DEL", character_key, index_key)
  end
end
redis.call("DEL", account_key)
return {1, ""}
`;
        return client.eval(
            script,
            3,
            accountKey,
            this._characterKey(worldId, 0),
            this._characterToAccountKey(worldId, 0),
            String(transitionTtl)
        );
    }

    async getAccountSession(worldId, accountId) {
        const { client } = this.ctx.getRedisGlobalAccess(worldId);
        const raw = await client.hgetall(this._accountKey(worldId, accountId));
        return this._accountHashToSession(raw);
    }

    async setAccountSession(worldId, accountId, session, ttlSeconds) {
        const { client } = this.ctx.getRedisGlobalAccess(worldId);
        await client.hset(this._accountKey(worldId, accountId), this._accountSessionToHash(session));
        await client.expire(this._accountKey(worldId, accountId), ttlSeconds);
    }

    async delAccountSession(worldId, accountId) {
        const { client } = this.ctx.getRedisGlobalAccess(worldId);
        await client.del(this._accountKey(worldId, accountId));
    }

    async refreshAccountSession(worldId, accountId, ttlSeconds) {
        const { client } = this.ctx.getRedisGlobalAccess(worldId);
        await client.expire(this._accountKey(worldId, accountId), ttlSeconds);
    }

    async getCharacterSession(worldId, characterId) {
        const { client } = this.ctx.getRedisGlobalAccess(worldId);
        const raw = await client.hgetall(this._characterKey(worldId, characterId));
        return this._characterHashToSession(raw);
    }

    async setCharacterSession(worldId, characterId, session, ttlSeconds) {
        const { client } = this.ctx.getRedisGlobalAccess(worldId);
        const multi = client.multi();
        multi.hset(this._characterKey(worldId, characterId), this._characterSessionToHash(session));
        multi.expire(this._characterKey(worldId, characterId), ttlSeconds);
        multi.set(this._characterToAccountKey(worldId, characterId), String(session.accountId), "EX", ttlSeconds);
        await multi.exec();
    }

    async delCharacterSession(worldId, characterId) {
        const { client } = this.ctx.getRedisGlobalAccess(worldId);
        await client.del(
            this._characterKey(worldId, characterId),
            this._characterToAccountKey(worldId, characterId)
        );
    }

    async refreshCharacterSession(worldId, characterId, ttlSeconds) {
        const { client } = this.ctx.getRedisGlobalAccess(worldId);
        const multi = client.multi();
        multi.expire(this._characterKey(worldId, characterId), ttlSeconds);
        multi.expire(this._characterToAccountKey(worldId, characterId), ttlSeconds);
        await multi.exec();
    }
}

module.exports = { SessionRepository };

