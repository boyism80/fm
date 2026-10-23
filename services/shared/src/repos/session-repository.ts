import { redisSessionKey } from "../redis-session-key";
import { AccountSessionState, accountSessionStateFromRedisHash } from "../session-state";
import { SessionErrorCode } from "../protobuf/generated/fminternal/internal_service";
import type { Debuff } from "../protobuf/generated/fminternal/internal_service";
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
            timestamps: {
                createdAt: hash.created_at || null,
                updatedAt: hash.updated_at || null,
                stateChangedAt: hash.state_changed_at || null,
            },
        };
    }

    async beginLogin(accountId: number, loginServerId: string, now: string, ttlSeconds: number) {
        const { client } = this.ctx.getRedisUnifiedAccess();
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
        return (await client.eval(
            script,
            1,
            accountKey,
            String(accountId),
            String(loginServerId),
            String(now),
            String(ttlSeconds),
        )) as Array<number | string>;
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
    ) {
        const { client } = this.ctx.getRedisUnifiedAccess();
        const accountKey = this.accountKey(accountId);
        const sourceCh =
            sourceChannelId != null && Number.isInteger(sourceChannelId) && sourceChannelId >= 0 ? sourceChannelId : -1;
        const sourceCs =
            sourceCashShopId != null && Number.isInteger(sourceCashShopId) && sourceCashShopId >= 0
                ? sourceCashShopId
                : -1;
        const script = `
local account_key = KEYS[1]
local world_id = ARGV[1]
local character_id = ARGV[2]
local character_name = ARGV[3]
local now = ARGV[4]
local ttl = tonumber(ARGV[5])
local source_channel_id = ARGV[6]
local client_ip = ARGV[7]
local debuffs = ARGV[8]
local source_cash_shop_id = ARGV[9]
if redis.call("EXISTS", account_key) == 0 then
  return {0, ${SessionErrorCode.SESSION_NOT_FOUND}}
end
local prev_state = tonumber(redis.call("HGET", account_key, "state") or "0")
local game_to_game_transfer = "0"
local return_channel_id = ""
if source_cash_shop_id ~= "-1" then
  local session_character_id = redis.call("HGET", account_key, "character_id") or ""
  local cash_shop_id = redis.call("HGET", account_key, "cash_shop_id") or ""
  if prev_state ~= ${AS.ACCOUNT_SESSION_STATE_CASH_SHOP} or session_character_id ~= character_id or cash_shop_id ~= source_cash_shop_id then
    return {0, ${SessionErrorCode.SESSION_NOT_OWNER}}
  end
  game_to_game_transfer = "1"
elseif source_channel_id == "-1" then
  local from_game = redis.call("HGET", account_key, "game_to_game_transfer") or "0"
  if prev_state ~= ${AS.ACCOUNT_SESSION_STATE_LOGIN} and (prev_state ~= ${AS.ACCOUNT_SESSION_STATE_TRANSITION} or from_game == "1") then
    return {0, ${SessionErrorCode.SESSION_NOT_OWNER}}
  end
else
  local session_character_id = redis.call("HGET", account_key, "character_id") or ""
  local channel_id = redis.call("HGET", account_key, "channel_id") or ""
  if prev_state ~= ${AS.ACCOUNT_SESSION_STATE_GAME} or session_character_id ~= character_id or channel_id ~= source_channel_id then
    return {0, ${SessionErrorCode.SESSION_NOT_OWNER}}
  end
  game_to_game_transfer = "1"
  return_channel_id = source_channel_id
end
redis.call("HSET", account_key,
  "world_id", world_id,
  "state", "${AS.ACCOUNT_SESSION_STATE_TRANSITION}",
  "character_id", character_id,
  "character_name", character_name,
  "client_ip", client_ip,
  "game_to_game_transfer", game_to_game_transfer,
  "debuffs", debuffs,
  "updated_at", now,
  "state_changed_at", now
)
redis.call("HDEL", account_key, "channel_id", "cash_shop_id", "return_channel_id")
if return_channel_id ~= "" then
  redis.call("HSET", account_key, "return_channel_id", return_channel_id)
end
redis.call("EXPIRE", account_key, ttl)
return {1, ${SessionErrorCode.SESSION_NONE}}
`;
        return (await client.eval(
            script,
            1,
            accountKey,
            String(worldId),
            String(characterId),
            String(characterName),
            String(now),
            String(ttlSeconds),
            String(sourceCh),
            clientIp,
            (sourceCh >= 0 || sourceCs >= 0) && debuffs.length > 0 ? JSON.stringify(debuffs) : "",
            String(sourceCs),
        )) as Array<number | string>;
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
        const { client } = this.ctx.getRedisUnifiedAccess();
        const accountKey = this.accountKey(accountId);
        const script = `
local account_key = KEYS[1]
local expected_account_id = tonumber(ARGV[1])
local expected_world_id = tonumber(ARGV[2])
local expected_character_id = tonumber(ARGV[3])
local cash_shop_id = ARGV[4]
local now = ARGV[5]
local ttl = tonumber(ARGV[6])
local client_ip = ARGV[7]
if redis.call("EXISTS", account_key) == 0 then
  return {0, ${SessionErrorCode.SESSION_NOT_FOUND}}
end
local state = tonumber(redis.call("HGET", account_key, "state") or "0")
local return_channel_id = redis.call("HGET", account_key, "return_channel_id") or ""
if state ~= ${AS.ACCOUNT_SESSION_STATE_TRANSITION} or return_channel_id == "" then
  return {0, ${SessionErrorCode.SESSION_UNKNOWN}}
end
local account_id = tonumber(redis.call("HGET", account_key, "account_id") or "0")
local world_id = tonumber(redis.call("HGET", account_key, "world_id") or "0")
local character_id = tonumber(redis.call("HGET", account_key, "character_id") or "0")
if account_id ~= expected_account_id or world_id ~= expected_world_id or character_id ~= expected_character_id then
  return {0, ${SessionErrorCode.SESSION_UNKNOWN}}
end
local transition_ip = redis.call("HGET", account_key, "client_ip") or ""
if client_ip == "" or transition_ip ~= client_ip then
  return {0, ${SessionErrorCode.SESSION_CLIENT_IP_MISMATCH}}
end

redis.call("HSET", account_key,
  "state", "${AS.ACCOUNT_SESSION_STATE_CASH_SHOP}",
  "cash_shop_id", cash_shop_id,
  "updated_at", now,
  "state_changed_at", now
)
redis.call("EXPIRE", account_key, ttl)
return {1, ${SessionErrorCode.SESSION_NONE}, return_channel_id}
`;
        return (await client.eval(
            script,
            1,
            accountKey,
            String(accountId),
            String(worldId),
            String(characterId),
            String(cashShopId),
            String(now),
            String(ttlSeconds),
            clientIp,
        )) as Array<number | string>;
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
        const { client } = this.ctx.getRedisUnifiedAccess();
        const accountKey = this.accountKey(accountId);
        const script = `
local account_key = KEYS[1]
local expected_account_id = tonumber(ARGV[1])
local expected_world_id = tonumber(ARGV[2])
local expected_character_id = tonumber(ARGV[3])
local channel_id = ARGV[4]
local now = ARGV[5]
local ttl = tonumber(ARGV[6])
local client_ip = ARGV[7]
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
local transition_ip = redis.call("HGET", account_key, "client_ip") or ""
if client_ip == "" or transition_ip ~= client_ip then
  return {0, ${SessionErrorCode.SESSION_CLIENT_IP_MISMATCH}}
end

redis.call("HSET", account_key,
  "state", "${AS.ACCOUNT_SESSION_STATE_GAME}",
  "login_server_connected", "0",
  "channel_id", channel_id,
  "game_to_game_transfer", "0",
  "updated_at", now,
  "state_changed_at", now
)
redis.call("HDEL", account_key, "debuffs")
redis.call("EXPIRE", account_key, ttl)
return {1, ${SessionErrorCode.SESSION_NONE}}
`;
        return (await client.eval(
            script,
            1,
            accountKey,
            String(accountId),
            String(worldId),
            String(characterId),
            String(channelId),
            String(now),
            String(ttlSeconds),
            clientIp,
        )) as Array<number | string>;
    }

    async refresh(
        accountId: number,
        loginTtl: number,
        transitionTtl: number,
        gameTtl: number,
        owner: { characterId: number; channelId?: number; cashShopId?: number } | null,
    ) {
        const { client } = this.ctx.getRedisUnifiedAccess();
        const accountKey = this.accountKey(accountId);
        const script = `
local account_key = KEYS[1]
local login_ttl = tonumber(ARGV[1])
local transition_ttl = tonumber(ARGV[2])
local game_ttl = tonumber(ARGV[3])
local owner_character_id = ARGV[4]
local owner_channel_id = ARGV[5]
local owner_cash_shop_id = ARGV[6]
if redis.call("EXISTS", account_key) == 0 then
  return {0, ${SessionErrorCode.SESSION_NOT_FOUND}, 0, "0"}
end
local state = tonumber(redis.call("HGET", account_key, "state") or "0")
if owner_character_id ~= "" then
  local character_id = redis.call("HGET", account_key, "character_id") or ""
  if owner_cash_shop_id ~= "" then
    local cash_shop_id = redis.call("HGET", account_key, "cash_shop_id") or ""
    if state ~= ${AS.ACCOUNT_SESSION_STATE_CASH_SHOP} or character_id ~= owner_character_id or cash_shop_id ~= owner_cash_shop_id then
      return {0, ${SessionErrorCode.SESSION_NOT_OWNER}, 0, tostring(state)}
    end
  else
    local channel_id = redis.call("HGET", account_key, "channel_id") or ""
    if state ~= ${AS.ACCOUNT_SESSION_STATE_GAME} or character_id ~= owner_character_id or channel_id ~= owner_channel_id then
      return {0, ${SessionErrorCode.SESSION_NOT_OWNER}, 0, tostring(state)}
    end
  end
end
local ttl = login_ttl
if state == ${AS.ACCOUNT_SESSION_STATE_TRANSITION} then
  ttl = transition_ttl
elseif state == ${AS.ACCOUNT_SESSION_STATE_GAME} or state == ${AS.ACCOUNT_SESSION_STATE_CASH_SHOP} then
  ttl = game_ttl
end
redis.call("EXPIRE", account_key, ttl)
return {1, ${SessionErrorCode.SESSION_NONE}, ttl, tostring(state)}
`;
        return (await client.eval(
            script,
            1,
            accountKey,
            String(loginTtl),
            String(transitionTtl),
            String(gameTtl),
            owner ? String(owner.characterId) : "",
            owner?.channelId != null ? String(owner.channelId) : "",
            owner?.cashShopId != null ? String(owner.cashShopId) : "",
        )) as Array<number | string>;
    }

    async logout(
        accountId: number,
        keep: boolean,
        requestChannelId?: number,
        ownerCharacterId?: number,
        requestCashShopId?: number,
    ) {
        const { client } = this.ctx.getRedisUnifiedAccess();
        const accountKey = this.accountKey(accountId);
        const ch =
            requestChannelId != null && Number.isInteger(requestChannelId) && requestChannelId >= 0
                ? requestChannelId
                : -1;
        const cs =
            requestCashShopId != null && Number.isInteger(requestCashShopId) && requestCashShopId >= 0
                ? requestCashShopId
                : -1;
        const script = `
local account_key = KEYS[1]
local keep = ARGV[1]
local owner_character_id = ARGV[3]
local request_cash_shop_id = ARGV[4]
local prev_state = 0
if redis.call("EXISTS", account_key) == 1 then
  prev_state = tonumber(redis.call("HGET", account_key, "state") or "0")
end
if keep == "1" then
  return {1, ${SessionErrorCode.SESSION_NONE}, prev_state}
end
if owner_character_id ~= "" then
  local character_id = redis.call("HGET", account_key, "character_id") or ""
  if request_cash_shop_id ~= "-1" then
    local session_cash_shop_id = redis.call("HGET", account_key, "cash_shop_id") or ""
    if prev_state ~= ${AS.ACCOUNT_SESSION_STATE_CASH_SHOP} or character_id ~= owner_character_id or session_cash_shop_id ~= request_cash_shop_id then
      return {0, ${SessionErrorCode.SESSION_NOT_OWNER}, prev_state}
    end
  else
    local session_channel_id = redis.call("HGET", account_key, "channel_id") or ""
    if prev_state ~= ${AS.ACCOUNT_SESSION_STATE_GAME} or character_id ~= owner_character_id or session_channel_id ~= ARGV[2] then
      return {0, ${SessionErrorCode.SESSION_NOT_OWNER}, prev_state}
    end
  end
end
redis.call("DEL", account_key)
return {1, ${SessionErrorCode.SESSION_NONE}, prev_state}
`;
        return (await client.eval(
            script,
            1,
            accountKey,
            keep ? "1" : "0",
            String(ch),
            ownerCharacterId != null ? String(ownerCharacterId) : "",
            String(cs),
        )) as Array<number | string>;
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
