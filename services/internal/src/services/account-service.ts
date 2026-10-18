import bcrypt from "bcrypt";
import { LoginAccountReply_Status } from "@fm/shared/protobuf/generated/fminternal/internal_service";
import type { InternalContext } from "@fm/shared/context/internal-context";
import type { AccountRepository } from "@fm/shared/repos/account-repository";
import type { AccountWorldRepository } from "@fm/shared/repos/account-world-repository";
import type { AccountModel } from "@fm/shared/types/repository-models";
import type { PgEndpoint } from "@fm/shared/types/internal-config";
import type { UnifiedRepository } from "@fm/shared/repos/unified-repository";
import type { DistributedLockService } from "./distributed-lock-service";

const BCRYPT_ROUNDS = 10;
const SECOND_PASSWORD_MAX_FAILURES = 5;
const SECOND_PASSWORD_LOCK_MS = 15 * 60 * 1000;

export type SecondPasswordResult =
    | "success"
    | "invalid"
    | "already_set"
    | "not_set"
    | "locked";

export class AccountService {
    private readonly ctx: InternalContext;
    private readonly accountRepo: AccountRepository;
    private readonly accountWorldRepo: AccountWorldRepository;
    private readonly unifiedRepo: UnifiedRepository;
    private readonly distributedLockService: DistributedLockService;

    constructor(
        internalContext: InternalContext,
        accountRepository: AccountRepository,
        accountWorldRepository: AccountWorldRepository,
        unifiedRepository: UnifiedRepository,
        distributedLockService: DistributedLockService
    ) {
        this.ctx = internalContext;
        this.accountRepo = accountRepository;
        this.accountWorldRepo = accountWorldRepository;
        this.unifiedRepo = unifiedRepository;
        this.distributedLockService = distributedLockService;
    }

    async loginAccount(loginId: string, password: string, macAddress: string | null, ipAddress: string | null, initialRole: number) {
        const identity = await this.unifiedRepo.findAccountByLoginId(loginId);

        if (!identity || identity.registered_at === null) {
            const accountId = identity ? identity.id : (await this.unifiedRepo.insertAccountIdentity(loginId)).id;
            const passwordHash = await bcrypt.hash(password, BCRYPT_ROUNDS);
            const role = Number.isInteger(initialRole) ? initialRole : 0;
            let saved: AccountModel;
            try {
                saved = await this.accountRepo.create({
                    accountId,
                    loginId,
                    passwordHash,
                    gender: 0,
                    role,
                    isBanned: false,
                    banReason: null,
                    isChatBlocked: false,
                    chatBlockedUntil: null,
                    lastLoginIp: ipAddress ?? null,
                    macAddress: macAddress ?? null,
                    secondPasswordHash: null,
                    secondPasswordFailures: 0,
                    secondPasswordLockedUntil: null,
                });
                await this.unifiedRepo.confirmAccountIdentity(accountId);
            } catch (err) {
                await this.unifiedRepo.deletePendingAccountIdentity(accountId).catch(() => {});
                throw err;
            }
            return {
                status: LoginAccountReply_Status.REGISTERED,
                accountId: saved.accountId,
                gender: saved.gender,
                role: saved.role,
                isChatBlocked: saved.isChatBlocked,
                chatBlockedUntil: saved.chatBlockedUntil,
            };
        }

        const account = await this.accountRepo.get(identity.id);
        if (!account) {
            throw new Error(`Account identity found but account data missing for id=${identity.id}`);
        }
        if (account.isBanned) {
            return { status: LoginAccountReply_Status.BANNED, accountId: account.accountId, role: account.role ?? 0 };
        }

        const match = await bcrypt.compare(password, account.passwordHash);
        if (!match) {
            return { status: LoginAccountReply_Status.WRONG_PASSWORD, role: account.role ?? 0 };
        }
        if (macAddress || ipAddress) {
            await this.accountRepo.setLogin(account.accountId, ipAddress ?? account.lastLoginIp ?? null, macAddress ?? account.macAddress ?? null);
        }

        return {
            status: LoginAccountReply_Status.SUCCESS,
            accountId: account.accountId,
            gender: account.gender,
            role: account.role,
            isChatBlocked: account.isChatBlocked,
            chatBlockedUntil: account.chatBlockedUntil,
        };
    }

    async setAccountRole(accountId: number, role: number) {
        const account = await this.accountRepo.get(accountId);
        if (!account) {
            return false;
        }
        await this.accountRepo.setRole(accountId, role);
        return true;
    }

    async setSecondPassword(accountId: number, password: string): Promise<SecondPasswordResult> {
        if (password.length < 4 || password.length > 16) {
            return "invalid";
        }

        await using _lock = await this.distributedLockService.acquireUnifiedLock(`second_password:a${accountId}`);
        const account = await this.accountRepo.get(accountId);
        if (!account) {
            return "invalid";
        }
        if (account.secondPasswordHash) {
            return "already_set";
        }

        const secondPasswordHash = await bcrypt.hash(password, BCRYPT_ROUNDS);
        await this.accountRepo.setSecondPassword(accountId, secondPasswordHash, 0, null);
        return "success";
    }

    async removeSecondPassword(accountId: number, password: string): Promise<SecondPasswordResult> {
        await using _lock = await this.distributedLockService.acquireUnifiedLock(`second_password:a${accountId}`);
        const account = await this.accountRepo.get(accountId);
        if (!account) {
            return "invalid";
        }
        if (!account.secondPasswordHash) {
            return "not_set";
        }

        const verified = await this.checkSecondPassword(account, password);
        if (verified !== "success") {
            return verified;
        }

        await this.accountRepo.setSecondPassword(accountId, null, 0, null);
        return "success";
    }

    async authorizeCharacterSelection(accountId: number, password?: string): Promise<SecondPasswordResult> {
        await using _lock = await this.distributedLockService.acquireUnifiedLock(`second_password:a${accountId}`);
        const account = await this.accountRepo.get(accountId);
        if (!account) {
            return "invalid";
        }
        if (!account.secondPasswordHash) {
            return password == null ? "success" : "not_set";
        }
        if (password == null) {
            return "invalid";
        }
        return this.checkSecondPassword(account, password);
    }

    // Accounts used to live in the data shards of the configured world. Each world's shards hand their slot counts
    // to account_worlds, and shards other than the unified database hand their account rows to it.
    async moveAccountsToUnified() {
        const config = this.ctx.appConfiguration.postgresql;
        const unified = config.unified ? AccountService.endpointKey(config.unified) : "";
        for (const [worldKey, world] of Object.entries(config.worlds)) {
            const worldId = Number(worldKey);
            const pools = this.ctx.getPgDataPools(worldId);
            const seen = new Set<string>();
            for (let i = 0; i < world.data.length; i += 1) {
                const key = AccountService.endpointKey(world.data[i]!);
                const pool = pools[i]!;
                if (seen.has(key)) {
                    continue;
                }
                seen.add(key);
                await this.accountWorldRepo.takeSlotsFromAccounts(worldId, pool);
                if (key !== unified) {
                    const moved = await this.accountRepo.moveFrom(pool);
                    if (moved > 0) {
                        console.log(`[account] moved ${moved} accounts from world ${worldId} shard ${key} to unified`);
                    }
                }
            }
        }
    }

    private static endpointKey(endpoint: PgEndpoint): string {
        return `${endpoint.host}:${endpoint.port}/${endpoint.database}`;
    }

    private async checkSecondPassword(account: AccountModel, password: string): Promise<SecondPasswordResult> {
        const now = Date.now();
        const lockedUntil = account.secondPasswordLockedUntil ? new Date(account.secondPasswordLockedUntil).getTime() : 0;
        if (lockedUntil > now) {
            return "locked";
        }

        const match = await bcrypt.compare(password, account.secondPasswordHash!);
        if (match) {
            if (account.secondPasswordFailures !== 0 || account.secondPasswordLockedUntil) {
                await this.accountRepo.setSecondPasswordFailures(account.accountId, 0, null);
            }
            return "success";
        }

        const failures = lockedUntil > 0 && lockedUntil <= now ? 1 : account.secondPasswordFailures + 1;
        await this.accountRepo.setSecondPasswordFailures(
            account.accountId,
            failures,
            failures >= SECOND_PASSWORD_MAX_FAILURES ? new Date(now + SECOND_PASSWORD_LOCK_MS).toISOString() : null
        );
        return failures >= SECOND_PASSWORD_MAX_FAILURES ? "locked" : "invalid";
    }
}
