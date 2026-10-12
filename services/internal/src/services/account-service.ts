import bcrypt from "bcrypt";
import { LoginAccountReply_Status } from "../protobuf/generated/fminternal/internal_service";
import type { InternalContext } from "../context/internal-context";
import type { AccountRepository } from "../repos/account-repository";
import type { AccountModel } from "../types/repository-models";
import type { UnifiedRepository } from "../repos/unified-repository";
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
    private readonly unifiedRepo: UnifiedRepository;
    private readonly distributedLockService: DistributedLockService;

    constructor(
        internalContext: InternalContext,
        accountRepository: AccountRepository,
        unifiedRepository: UnifiedRepository,
        distributedLockService: DistributedLockService
    ) {
        this.ctx = internalContext;
        this.accountRepo = accountRepository;
        this.unifiedRepo = unifiedRepository;
        this.distributedLockService = distributedLockService;
    }

    private worldId() {
        return this.ctx.appConfiguration.app.world_id;
    }

    async loginAccount(loginId: string, password: string, macAddress: string | null, ipAddress: string | null, initialRole: number) {
        const worldId = this.worldId();
        const identity = await this.unifiedRepo.findAccountByLoginId(loginId);

        if (!identity || identity.registered_at === null) {
            const accountId = identity ? identity.id : (await this.unifiedRepo.insertAccountIdentity(loginId)).id;
            const passwordHash = await bcrypt.hash(password, BCRYPT_ROUNDS);
            const role = Number.isInteger(initialRole) ? initialRole : 0;
            const newAccount = {
                accountId,
                loginId,
                passwordHash,
                gender: 0,
                role,
                isBanned: false,
                banReason: null,
                isChatBlocked: false,
                chatBlockedUntil: null,
                characterSlotCount: 6,
                lastLoginIp: ipAddress ?? null,
                macAddress: macAddress ?? null,
                secondPasswordHash: null,
                secondPasswordFailures: 0,
                secondPasswordLockedUntil: null,
            };
            let saved: Awaited<ReturnType<AccountRepository["set"]>>;
            try {
                saved = await this.accountRepo.set(worldId, newAccount);
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
                characterSlotCount: saved.characterSlotCount,
            };
        }

        const account = await this.accountRepo.get(worldId, identity.id);
        if (!account) {
            throw new Error(`Account identity found but world data missing for id=${identity.id}`);
        }
        if (account.isBanned) {
            return { status: LoginAccountReply_Status.BANNED, accountId: account.accountId, role: account.role ?? 0 };
        }

        const match = await bcrypt.compare(password, account.passwordHash);
        if (!match) {
            return { status: LoginAccountReply_Status.WRONG_PASSWORD, role: account.role ?? 0 };
        }
        if (macAddress || ipAddress) {
            await this.accountRepo.setLogin(worldId, account.accountId, ipAddress ?? account.lastLoginIp ?? null, macAddress ?? account.macAddress ?? null);
        }

        return {
            status: LoginAccountReply_Status.SUCCESS,
            accountId: account.accountId,
            gender: account.gender,
            role: account.role,
            isChatBlocked: account.isChatBlocked,
            chatBlockedUntil: account.chatBlockedUntil,
            characterSlotCount: account.characterSlotCount,
        };
    }

    async setAccountRole(accountId: number, role: number) {
        const account = await this.accountRepo.get(this.worldId(), accountId);
        if (!account) {
            return false;
        }
        await this.accountRepo.setRole(this.worldId(), accountId, role);
        return true;
    }

    async setSecondPassword(worldId: number, accountId: number, password: string): Promise<SecondPasswordResult> {
        if (password.length < 4 || password.length > 16) {
            return "invalid";
        }

        await using _lock = await this.distributedLockService.acquireWorldDataLock(worldId, `second_password:a${accountId}`);
        const account = await this.accountRepo.get(worldId, accountId);
        if (!account) {
            return "invalid";
        }
        if (account.secondPasswordHash) {
            return "already_set";
        }

        const secondPasswordHash = await bcrypt.hash(password, BCRYPT_ROUNDS);
        await this.accountRepo.setSecondPassword(worldId, accountId, secondPasswordHash, 0, null);
        return "success";
    }

    async removeSecondPassword(worldId: number, accountId: number, password: string): Promise<SecondPasswordResult> {
        await using _lock = await this.distributedLockService.acquireWorldDataLock(worldId, `second_password:a${accountId}`);
        const account = await this.accountRepo.get(worldId, accountId);
        if (!account) {
            return "invalid";
        }
        if (!account.secondPasswordHash) {
            return "not_set";
        }

        const verified = await this.checkSecondPassword(worldId, account, password);
        if (verified !== "success") {
            return verified;
        }

        await this.accountRepo.setSecondPassword(worldId, accountId, null, 0, null);
        return "success";
    }

    async authorizeCharacterSelection(worldId: number, accountId: number, password?: string): Promise<SecondPasswordResult> {
        await using _lock = await this.distributedLockService.acquireWorldDataLock(worldId, `second_password:a${accountId}`);
        const account = await this.accountRepo.get(worldId, accountId);
        if (!account) {
            return "invalid";
        }
        if (!account.secondPasswordHash) {
            return password == null ? "success" : "not_set";
        }
        if (password == null) {
            return "invalid";
        }
        return this.checkSecondPassword(worldId, account, password);
    }

    private async checkSecondPassword(worldId: number, account: AccountModel, password: string): Promise<SecondPasswordResult> {
        const now = Date.now();
        const lockedUntil = account.secondPasswordLockedUntil ? new Date(account.secondPasswordLockedUntil).getTime() : 0;
        if (lockedUntil > now) {
            return "locked";
        }

        const match = await bcrypt.compare(password, account.secondPasswordHash!);
        if (match) {
            if (account.secondPasswordFailures !== 0 || account.secondPasswordLockedUntil) {
                await this.accountRepo.setSecondPasswordFailures(worldId, account.accountId, 0, null);
            }
            return "success";
        }

        const failures = lockedUntil > 0 && lockedUntil <= now ? 1 : account.secondPasswordFailures + 1;
        await this.accountRepo.setSecondPasswordFailures(
            worldId,
            account.accountId,
            failures,
            failures >= SECOND_PASSWORD_MAX_FAILURES ? new Date(now + SECOND_PASSWORD_LOCK_MS).toISOString() : null
        );
        return failures >= SECOND_PASSWORD_MAX_FAILURES ? "locked" : "invalid";
    }
}
