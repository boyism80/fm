import type { Pool } from "pg";
import type { InternalContext } from "../context/internal-context";
import type { RepositoryTxOptions } from "../types/repository-contracts";

export const DEFAULT_CHARACTER_SLOT_COUNT = 6;

// Per-world account data, on the account's data shard of that world. A missing row means the defaults.
export class AccountWorldRepository {
    private readonly ctx: InternalContext;

    constructor(internalContext: InternalContext) {
        this.ctx = internalContext;
    }

    async findCharacterSlotCount(
        worldId: number,
        accountId: number,
        options: RepositoryTxOptions = {},
    ): Promise<number> {
        const pool = this.ctx.getPgDataPool(worldId, accountId);
        const { rows } = await this.ctx
            .pgClient(pool, options)
            .query("SELECT character_slot_count FROM account_worlds WHERE account_id = $1 AND world_id = $2", [
                accountId,
                worldId,
            ]);
        return (
            (rows[0] as { character_slot_count: number } | undefined)?.character_slot_count ??
            DEFAULT_CHARACTER_SLOT_COUNT
        );
    }

    async setCharacterSlotCount(
        worldId: number,
        accountId: number,
        count: number,
        options: RepositoryTxOptions = {},
    ): Promise<void> {
        const pool = this.ctx.getPgDataPool(worldId, accountId);
        await this.ctx.pgClient(pool, options).query(
            `INSERT INTO account_worlds (account_id, world_id, character_slot_count) VALUES ($1, $2, $3)
             ON CONFLICT (account_id, world_id) DO UPDATE SET character_slot_count = EXCLUDED.character_slot_count`,
            [accountId, worldId, count],
        );
    }

    // Slot counts bought before slots became per world were stored on the shard's account rows. They are copied
    // once and reset on the account rows, so a later run finds nothing to copy.
    async takeSlotsFromAccounts(worldId: number, shard: Pool): Promise<void> {
        await this.ctx.withPgClient(shard, async (client) => {
            await client.query("BEGIN");
            try {
                await client.query(
                    `INSERT INTO account_worlds (account_id, world_id, character_slot_count)
                     SELECT id, $1, character_slot_count FROM accounts WHERE character_slot_count <> $2
                     ON CONFLICT (account_id, world_id) DO NOTHING`,
                    [worldId, DEFAULT_CHARACTER_SLOT_COUNT],
                );
                await client.query("UPDATE accounts SET character_slot_count = $1 WHERE character_slot_count <> $1", [
                    DEFAULT_CHARACTER_SLOT_COUNT,
                ]);
                await client.query("COMMIT");
            } catch (err) {
                await client.query("ROLLBACK").catch(() => {});
                throw err;
            }
        });
    }
}
