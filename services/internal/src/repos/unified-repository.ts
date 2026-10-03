import type { InternalContext } from "../context/internal-context";
import { redisNameKey } from "../redis-name-key";

// Same as the login session TTL; RefreshSession extends it while the login session lives.
const CHARACTER_NAME_RESERVATION_TTL_SECONDS = 120;

export interface AccountIdentityRow {
    id: number;
    login_id: string;
    created_at: Date | string;
}

export interface CharacterNameRegistryRow {
    character_id: number;
    name: string;
    account_id: number;
    world_id: number;
    status?: string;
}

export interface GuildNameRegistryRow {
    guild_id: number;
    name: string;
    world_id: number;
    created_at?: Date | string;
}

export class UnifiedRepository {
    private readonly ctx: InternalContext;

    constructor(internalContext: InternalContext) {
        this.ctx = internalContext;
    }

    private pool() {
        return this.ctx.getPgUnifiedPool();
    }

    async findAccountByLoginId(loginId: string): Promise<AccountIdentityRow | null> {
        const { rows } = await this.pool().query(
            "SELECT id, login_id, created_at FROM account_identity WHERE login_id = $1",
            [loginId]
        );
        return rows[0] ?? null;
    }

    async insertAccountIdentity(loginId: string): Promise<AccountIdentityRow> {
        const { rows } = await this.pool().query(
            "INSERT INTO account_identity (login_id) VALUES ($1) RETURNING id, login_id, created_at",
            [loginId]
        );
        return rows[0];
    }

    async findCharacterNameEntry(name: string): Promise<CharacterNameRegistryRow | null> {
        const { rows } = await this.pool().query(
            "SELECT character_id, name, account_id, world_id, status FROM character_name_registry WHERE LOWER(name) = LOWER($1)",
            [name]
        );
        return rows[0] ?? null;
    }

    async nextCharacterId(): Promise<number> {
        const { rows } = await this.pool().query(
            "SELECT nextval(pg_get_serial_sequence('character_name_registry', 'character_id')) AS id"
        );
        return Number(rows[0].id);
    }

    async insertCharacterName(characterId: number, name: string, accountId: number, worldId: number): Promise<void> {
        await this.pool().query(
            `INSERT INTO character_name_registry (character_id, name, account_id, world_id)
             VALUES ($1, $2, $3, $4)`,
            [characterId, name, accountId, worldId]
        );
    }

    private nameReservationKey(name: string) {
        return redisNameKey(`reserve:${name.toLowerCase()}`);
    }

    private accountNameReservationKey(accountId: number) {
        return redisNameKey(`reserve_account:${accountId}`);
    }

    /** One name per account: reserving a new name releases the account's previous one. */
    async reserveCharacterName(name: string, accountId: number): Promise<boolean> {
        const { client } = this.ctx.getRedisUnifiedAccess();
        const script = `
local name_key = KEYS[1]
local account_key = KEYS[2]
local account_id = ARGV[1]
local name = ARGV[2]
local prefix = ARGV[3]
local ttl = tonumber(ARGV[4])
local owner = redis.call("GET", name_key)
if owner and owner ~= account_id then
  return 0
end
local prev = redis.call("GET", account_key)
if prev and prev ~= name and redis.call("GET", prefix .. prev) == account_id then
  redis.call("DEL", prefix .. prev)
end
redis.call("SET", name_key, account_id, "EX", ttl)
redis.call("SET", account_key, name, "EX", ttl)
return 1
`;
        const lower = name.toLowerCase();
        const ok = await client.eval(
            script,
            2,
            this.nameReservationKey(lower),
            this.accountNameReservationKey(accountId),
            String(accountId),
            lower,
            redisNameKey("reserve:"),
            String(CHARACTER_NAME_RESERVATION_TTL_SECONDS)
        );
        return Number(ok) === 1;
    }

    async refreshCharacterNameReservation(accountId: number): Promise<void> {
        const { client } = this.ctx.getRedisUnifiedAccess();
        const script = `
local account_key = KEYS[1]
local account_id = ARGV[1]
local prefix = ARGV[2]
local ttl = tonumber(ARGV[3])
local name = redis.call("GET", account_key)
if not name then
  return 0
end
redis.call("EXPIRE", account_key, ttl)
if redis.call("GET", prefix .. name) == account_id then
  redis.call("EXPIRE", prefix .. name, ttl)
end
return 1
`;
        await client.eval(
            script,
            1,
            this.accountNameReservationKey(accountId),
            String(accountId),
            redisNameKey("reserve:"),
            String(CHARACTER_NAME_RESERVATION_TTL_SECONDS)
        );
    }

    async releaseCharacterNameReservation(accountId: number): Promise<void> {
        const { client } = this.ctx.getRedisUnifiedAccess();
        const script = `
local account_key = KEYS[1]
local account_id = ARGV[1]
local prefix = ARGV[2]
local name = redis.call("GET", account_key)
if not name then
  return 0
end
redis.call("DEL", account_key)
if redis.call("GET", prefix .. name) == account_id then
  redis.call("DEL", prefix .. name)
end
return 1
`;
        await client.eval(script, 1, this.accountNameReservationKey(accountId), String(accountId), redisNameKey("reserve:"));
    }

    async deleteCharacterName(characterId: number) {
        await this.pool().query(
            "DELETE FROM character_name_registry WHERE character_id = $1",
            [characterId]
        );
    }

    async findGuildNameEntry(name: string): Promise<GuildNameRegistryRow | null> {
        const { rows } = await this.pool().query(
            "SELECT guild_id, name, world_id, created_at FROM guild_name_registry WHERE LOWER(name) = LOWER($1)",
            [name]
        );
        return rows[0] ?? null;
    }

    async reserveGuildName(name: string, worldId: number): Promise<GuildNameRegistryRow> {
        const { rows } = await this.pool().query(
            `INSERT INTO guild_name_registry (name, world_id)
             VALUES ($1, $2)
             RETURNING guild_id, name, world_id, created_at`,
            [name, worldId]
        );
        return rows[0];
    }

    async deleteGuildName(guildId: number) {
        await this.pool().query(
            "DELETE FROM guild_name_registry WHERE guild_id = $1",
            [guildId]
        );
    }
}
