import bcrypt from "bcrypt";
import type { InternalContext } from "../context/internal-context";
import type { AccountRepository } from "../repos/account-repository";
import type { UnifiedRepository } from "../repos/unified-repository";

const BCRYPT_ROUNDS = 10;

export const LoginStatus = {
    SUCCESS: 0,
    WRONG_PASSWORD: 1,
    BANNED: 2,
    REGISTERED: 3,
};

export class AccountService {
    private readonly ctx: InternalContext;
    private readonly accountRepo: AccountRepository;
    private readonly unifiedRepo: UnifiedRepository;

    constructor(
        internalContext: InternalContext,
        accountRepository: AccountRepository,
        unifiedRepository: UnifiedRepository
    ) {
        this.ctx = internalContext;
        this.accountRepo = accountRepository;
        this.unifiedRepo = unifiedRepository;
    }

    private worldId() {
        return this.ctx.appConfiguration.app.world_id;
    }

    async loginAccount(loginId: string, password: string, macAddress: string | null, ipAddress: string | null, initialRole: number) {
        const worldId = this.worldId();
        const identity = await this.unifiedRepo.findAccountByLoginId(loginId);

        if (!identity) {
            const newIdentity = await this.unifiedRepo.insertAccountIdentity(loginId);
            const accountId = newIdentity.id;
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
            };
            const saved = await this.accountRepo.set(worldId, newAccount);
            return {
                status: LoginStatus.REGISTERED,
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
            return { status: LoginStatus.BANNED, accountId: account.accountId, role: account.role ?? 0 };
        }

        const match = await bcrypt.compare(password, account.passwordHash);
        if (!match) {
            return { status: LoginStatus.WRONG_PASSWORD, role: account.role ?? 0 };
        }
        if (macAddress || ipAddress) {
            const updated = { ...account, macAddress: macAddress ?? account.macAddress, lastLoginIp: ipAddress ?? account.lastLoginIp };
            await this.accountRepo.set(worldId, updated);
        }

        return {
            status: LoginStatus.SUCCESS,
            accountId: account.accountId,
            gender: account.gender,
            role: account.role,
            isChatBlocked: account.isChatBlocked,
            chatBlockedUntil: account.chatBlockedUntil,
            characterSlotCount: account.characterSlotCount,
        };
    }
}
