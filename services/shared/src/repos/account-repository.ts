import type { Pool } from "pg";
import type { InternalContext } from "../context/internal-context";
import type { RepositoryQueryValue } from "../types/repository-contracts";
import type { AccountModel, AccountRow } from "../types/repository-models";

export type { AccountModel };

const MOVE_BATCH = 1000;

// Account-wide data lives in the unified database. Character slots are per world (AccountWorldRepository).
export class AccountRepository {
    private readonly ctx: InternalContext;

    constructor(internalContext: InternalContext) {
        this.ctx = internalContext;
    }

    private pool() {
        return this.ctx.getPgUnifiedPool();
    }

    async get(accountId: number): Promise<AccountModel | null> {
        const { rows } = await this.pool().query("SELECT * FROM accounts WHERE id = $1", [accountId]);
        const row = rows[0] as AccountRow | undefined;
        if (!row) {
            return null;
        }
        return AccountRepository.rowToModel(row);
    }

    async create(model: AccountModel): Promise<AccountModel> {
        const { rows } = await this.pool().query(
            `INSERT INTO accounts
               (id, login_id, password_hash, gender, role, is_banned, ban_reason, is_chat_blocked, chat_blocked_until,
                last_login_ip, mac_address, second_password_hash, second_password_failures, second_password_locked_until,
                created_at, updated_at)
             VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, NOW(), NOW())
             RETURNING *`,
            [
                model.accountId,
                model.loginId,
                model.passwordHash,
                model.gender,
                model.role ?? 0,
                model.isBanned,
                model.banReason ?? null,
                model.isChatBlocked,
                model.chatBlockedUntil ?? null,
                model.lastLoginIp ?? null,
                model.macAddress ?? null,
                model.secondPasswordHash ?? null,
                model.secondPasswordFailures,
                model.secondPasswordLockedUntil ?? null,
            ]
        );
        return AccountRepository.rowToModel(rows[0] as AccountRow);
    }

    // An existing account is written column by column, never as a whole row. The second password state is
    // written under the `second_password:a{id}` lock and the other columns are not, so writing back a
    // whole stale row would undo failures and lockouts (or a login) recorded in the meantime.
    async setSecondPassword(accountId: number, hash: string | null, failures: number, lockedUntil: string | null): Promise<void> {
        await this.update(
            accountId,
            "second_password_hash = $2, second_password_failures = $3, second_password_locked_until = $4",
            [hash, failures, lockedUntil]
        );
    }

    async setSecondPasswordFailures(accountId: number, failures: number, lockedUntil: string | null): Promise<void> {
        await this.update(accountId, "second_password_failures = $2, second_password_locked_until = $3", [failures, lockedUntil]);
    }

    async setLogin(accountId: number, lastLoginIp: string | null, macAddress: string | null): Promise<void> {
        await this.update(accountId, "last_login_ip = $2, mac_address = $3", [lastLoginIp, macAddress]);
    }

    async setRole(accountId: number, role: number): Promise<void> {
        await this.update(accountId, "role = $2", [role]);
    }

    private async update(accountId: number, columns: string, values: RepositoryQueryValue[]): Promise<void> {
        await this.pool().query(`UPDATE accounts SET ${columns}, updated_at = NOW() WHERE id = $1`, [accountId, ...values]);
    }

    // Accounts used to live in a world data shard. Rows still there are copied to the unified database and
    // removed from the shard, so the move happens once.
    async moveFrom(shard: Pool): Promise<number> {
        let moved = 0;
        while (true) {
            const { rows } = await shard.query(`SELECT * FROM accounts ORDER BY id LIMIT ${MOVE_BATCH}`);
            if (rows.length === 0) {
                return moved;
            }
            await this.pool().query(
                "INSERT INTO accounts SELECT * FROM jsonb_populate_recordset(NULL::accounts, $1) ON CONFLICT (id) DO NOTHING",
                [JSON.stringify(rows)]
            );
            await shard.query("DELETE FROM accounts WHERE id = ANY($1)", [(rows as AccountRow[]).map((row) => row.id)]);
            moved += rows.length;
        }
    }

    private static rowToModel(row: AccountRow): AccountModel {
        return {
            accountId: row.id,
            loginId: row.login_id,
            passwordHash: row.password_hash,
            gender: row.gender,
            role: row.role,
            isBanned: row.is_banned,
            banReason: row.ban_reason ?? null,
            isChatBlocked: row.is_chat_blocked,
            chatBlockedUntil: row.chat_blocked_until ?? null,
            lastLoginIp: row.last_login_ip ?? null,
            macAddress: row.mac_address ?? null,
            secondPasswordHash: row.second_password_hash ?? null,
            secondPasswordFailures: row.second_password_failures,
            secondPasswordLockedUntil: row.second_password_locked_until ?? null,
        };
    }
}
