const { redisCacheKey } = require("../redis-cache-key");
import { ValueRepository } from "./value-repository";
import type { RepositoryQuery } from "../types/repository-contracts";
import type { AccountDeleteModel, AccountModel, AccountRow } from "../types/repository-models";

export type { AccountModel };

export class AccountRepository extends ValueRepository<AccountModel, AccountRow, number> {
    getKey(model: AccountModel) {
        return model.accountId;
    }

    getTtlSeconds() {
        return 300;
    }

    getRedisKey(worldId: number, accountId: number) {
        return redisCacheKey(`w${worldId}:account:${accountId}`);
    }

    onSelect(accountId: number): RepositoryQuery {
        return {
            text: "SELECT * FROM accounts WHERE id = $1",
            values: [accountId],
        };
    }

    onUpsert(row: AccountRow): RepositoryQuery {
        return {
            text: `INSERT INTO accounts
                     (id, login_id, password_hash, gender, role, is_banned, ban_reason,
                      is_chat_blocked, chat_blocked_until, character_slot_count, last_login_ip,
                      mac_address, created_at, updated_at)
                   VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,NOW(),NOW())
                   ON CONFLICT (id) DO UPDATE SET
                     password_hash        = EXCLUDED.password_hash,
                     gender               = EXCLUDED.gender,
                     role                 = EXCLUDED.role,
                     is_banned            = EXCLUDED.is_banned,
                     ban_reason           = EXCLUDED.ban_reason,
                     is_chat_blocked      = EXCLUDED.is_chat_blocked,
                     chat_blocked_until   = EXCLUDED.chat_blocked_until,
                     character_slot_count = EXCLUDED.character_slot_count,
                     last_login_ip        = EXCLUDED.last_login_ip,
                     mac_address          = EXCLUDED.mac_address,
                     updated_at           = NOW()
                   RETURNING *`,
            values: [
                row.id,
                row.login_id,
                row.password_hash,
                row.gender,
                row.role,
                row.is_banned,
                row.ban_reason ?? null,
                row.is_chat_blocked,
                row.chat_blocked_until ?? null,
                row.character_slot_count,
                row.last_login_ip ?? null,
                row.mac_address ?? null,
            ],
        };
    }

    onDelete(row: AccountDeleteModel): RepositoryQuery {
        return {
            text: "DELETE FROM accounts WHERE id = $1",
            values: [Number(row.accountId)],
        };
    }

    rowToModel(row: AccountRow): AccountModel {
        return {
            accountId: Number(row.id),
            loginId: row.login_id,
            passwordHash: row.password_hash,
            gender: Number(row.gender),
            role: Number(row.role),
            isBanned: Boolean(row.is_banned),
            banReason: row.ban_reason ?? null,
            isChatBlocked: Boolean(row.is_chat_blocked),
            chatBlockedUntil: row.chat_blocked_until ?? null,
            characterSlotCount: Number(row.character_slot_count),
            lastLoginIp: row.last_login_ip ?? null,
            macAddress: row.mac_address ?? null,
        };
    }

    modelToRow(model: AccountModel): AccountRow {
        return {
            id: model.accountId,
            login_id: model.loginId,
            password_hash: model.passwordHash,
            gender: model.gender,
            role: model.role ?? 0,
            is_banned: model.isBanned,
            ban_reason: model.banReason ?? null,
            is_chat_blocked: model.isChatBlocked,
            chat_blocked_until: model.chatBlockedUntil ?? null,
            character_slot_count: model.characterSlotCount,
            last_login_ip: model.lastLoginIp ?? null,
            mac_address: model.macAddress ?? null,
        };
    }
}
