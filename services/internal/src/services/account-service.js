"use strict";

const bcrypt = require("bcrypt");

const BCRYPT_ROUNDS = 10;

const LoginStatus = {
    SUCCESS: 0,
    WRONG_PASSWORD: 1,
    BANNED: 2,
    REGISTERED: 3,
};

class AccountService {
    constructor(internalContext, accountRepository, unifiedRepository) {
        this.ctx = internalContext;
        this.accountRepo = accountRepository;
        this.unifiedRepo = unifiedRepository;
    }

    _worldId() {
        return this.ctx.appConfiguration.app.world_id;
    }

    async loginAccount(loginId, password, macAddress, ipAddress, initialRole) {
        const worldId = this._worldId();
        const identity = await this.unifiedRepo.findAccountByLoginId(loginId);

        if (!identity) {
            const newIdentity = await this.unifiedRepo.insertAccountIdentity(loginId);
            const accountId = newIdentity.id;
            const passwordHash = await bcrypt.hash(password, BCRYPT_ROUNDS);

            const role = Number.isInteger(Number(initialRole)) ? Number(initialRole) : 0;
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
            const saved = await this.accountRepo.save(worldId, newAccount);

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

        const account = await this.accountRepo.getById(worldId, identity.id);
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
            await this.accountRepo.save(worldId, updated);
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

module.exports = { AccountService, LoginStatus };
