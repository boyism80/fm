import type { InternalContext } from "../context/internal-context";

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
            "SELECT character_id, name, account_id, world_id, status FROM character_name_registry WHERE LOWER(name) = LOWER($1) AND reserved_at IS NULL",
            [name]
        );
        return rows[0] ?? null;
    }

    /**
     * Atomically reserves the name in Postgres.
     * - New name: inserts a pending row and returns the new character_id.
     * - Same account re-reservation: refreshes reserved_at and returns the existing character_id.
     * - Expired reservation by another account: deletes the stale row and inserts fresh (new serial id),
     *   so no orphan data from a prior failed creation can conflict with the new character_id.
     * - Confirmed or actively reserved by another account: returns null.
     * One pending reservation per account: any other pending row for this account is deleted first.
     */
    async reserveCharacterName(name: string, accountId: number, worldId: number): Promise<number | null> {
        return this.ctx.withPgClient(this.pool(), async (client) => {
            try {
                await client.query("BEGIN");

                // Release any prior pending reservation for this account under a different name.
                await client.query(
                    `DELETE FROM character_name_registry WHERE account_id = $1 AND reserved_at IS NOT NULL AND LOWER(name) != LOWER($2)`,
                    [accountId, name]
                );

                const { rows } = await client.query(
                    `SELECT character_id, account_id, reserved_at,
                      reserved_at > NOW() - make_interval(secs => $2::numeric) AS reservation_active
                     FROM character_name_registry WHERE LOWER(name) = LOWER($1) FOR UPDATE`,
                    [name, CHARACTER_NAME_RESERVATION_TTL_SECONDS]
                );

                if (rows.length === 0) {
                    const { rows: ins } = await client.query(
                        `INSERT INTO character_name_registry (name, account_id, world_id, reserved_at) VALUES ($1, $2, $3, NOW()) RETURNING character_id`,
                        [name, accountId, worldId]
                    );
                    await client.query("COMMIT");
                    return Number(ins[0].character_id);
                }

                const { character_id, account_id: existingAccountId, reserved_at } = rows[0];

                if (reserved_at === null) {
                    await client.query("ROLLBACK");
                    return null;
                }

                if (Number(existingAccountId) === accountId) {
                    await client.query(
                        `UPDATE character_name_registry SET reserved_at = NOW() WHERE character_id = $1`,
                        [character_id]
                    );
                    await client.query("COMMIT");
                    return Number(character_id);
                }

                if (rows[0].reservation_active) {
                    await client.query("ROLLBACK");
                    return null;
                }

                // Expired reservation by a different account: delete and re-insert to get a
                // fresh serial character_id. Reusing the old id risks conflicting with orphan
                // data left by a failed creation that never ran its cleanup (e.g. process kill).
                await client.query(
                    `DELETE FROM character_name_registry WHERE character_id = $1`,
                    [character_id]
                );
                const { rows: ins } = await client.query(
                    `INSERT INTO character_name_registry (name, account_id, world_id, reserved_at) VALUES ($1, $2, $3, NOW()) RETURNING character_id`,
                    [name, accountId, worldId]
                );
                await client.query("COMMIT");
                return Number(ins[0].character_id);
            } catch (err) {
                await client.query("ROLLBACK").catch(() => {});
                throw err;
            }
        });
    }

    async confirmCharacterName(name: string): Promise<void> {
        await this.pool().query(
            `UPDATE character_name_registry SET reserved_at = NULL WHERE LOWER(name) = LOWER($1) AND reserved_at IS NOT NULL`,
            [name]
        );
    }

    async refreshCharacterNameReservation(accountId: number): Promise<void> {
        await this.pool().query(
            `UPDATE character_name_registry SET reserved_at = NOW() WHERE account_id = $1 AND reserved_at IS NOT NULL`,
            [accountId]
        );
    }

    async releaseCharacterNameReservation(accountId: number): Promise<void> {
        await this.pool().query(
            `DELETE FROM character_name_registry WHERE account_id = $1 AND reserved_at IS NOT NULL`,
            [accountId]
        );
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
