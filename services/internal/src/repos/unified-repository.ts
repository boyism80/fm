import type { InternalContext } from "../context/internal-context";

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

    async reserveCharacterName(name: string, accountId: number, worldId: number): Promise<CharacterNameRegistryRow> {
        const { rows } = await this.pool().query(
            `INSERT INTO character_name_registry (name, account_id, world_id)
             VALUES ($1, $2, $3)
             RETURNING character_id, name, account_id, world_id`,
            [name, accountId, worldId]
        );
        return rows[0];
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
