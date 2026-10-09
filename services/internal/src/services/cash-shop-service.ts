import { randomBytes } from "node:crypto";
import {
    CashCouponKind,
    CashCurrency,
    CashItem,
    CashShopResult,
    CashSlotKind,
    type BuyCashItemReply,
    type BuyCashItemRequest,
    type BuyCashQuestItemReply,
    type BuyCashRingReply,
    type BuyCashRingRequest,
    type CashRing,
    type ExpandCashSlotReply,
    type ExpandCashSlotRequest,
    type FindCashCouponReply,
    type GiftCashItemReply,
    type GiftCashItemRequest,
    type PayBackCashItemReply,
    type PayBackCashItemRequest,
    type RedeemCashCouponReply,
    type RedeemCashCouponRequest,
} from "../protobuf/generated/fminternal/internal_service";
import type { AccountRepository } from "../repos/account-repository";
import type { CashBalance, CashShopRepository, CashItemRow } from "../repos/cash-shop-repository";
import type { CharacterRepository } from "../repos/character-repository";
import type { StorageRepository } from "../repos/storage-repository";
import type { InventoryRepository, InventoryModel } from "../repos/inventory-repository";
import type { SessionRepository } from "../repos/session-repository";
import type { UnifiedRepository } from "../repos/unified-repository";
import type { RabbitMQService } from "./rabbitmq-service";

const LOCKER_CAPACITY = 100;
const DEFAULT_STORAGE_SLOTS = 4;
const COUPON_ALPHABET = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789";
const COUPON_LENGTH = 16;
const AMQ_DIRECT_EXCHANGE = "amq.direct";

export class CashShopService {
    private readonly cashShopRepo: CashShopRepository;
    private readonly accountRepo: AccountRepository;
    private readonly storageRepo: StorageRepository;
    private readonly inventoryRepo: InventoryRepository;
    private readonly characterRepo: CharacterRepository;
    private readonly unifiedRepo: UnifiedRepository;
    private readonly sessionRepo: SessionRepository;
    private readonly rabbitmqService: RabbitMQService;

    constructor(
        cashShopRepository: CashShopRepository,
        accountRepository: AccountRepository,
        storageRepository: StorageRepository,
        inventoryRepository: InventoryRepository,
        characterRepository: CharacterRepository,
        unifiedRepository: UnifiedRepository,
        sessionRepository: SessionRepository,
        rabbitmqService: RabbitMQService
    ) {
        this.cashShopRepo = cashShopRepository;
        this.accountRepo = accountRepository;
        this.storageRepo = storageRepository;
        this.inventoryRepo = inventoryRepository;
        this.characterRepo = characterRepository;
        this.unifiedRepo = unifiedRepository;
        this.sessionRepo = sessionRepository;
        this.rabbitmqService = rabbitmqService;
    }

    private toProto(row: CashItemRow): CashItem {
        return CashItem.fromJSON(row.item);
    }

    private debit(currency: CashCurrency, price: number): CashBalance {
        if (currency === CashCurrency.CASH_CURRENCY_MAPLE_POINT) {
            return { nxCash: 0, maplePoint: -price };
        }
        return { nxCash: -price, maplePoint: 0 };
    }

    private affordable(balance: CashBalance, delta: CashBalance): boolean {
        return balance.nxCash + delta.nxCash >= 0 && balance.maplePoint + delta.maplePoint >= 0;
    }

    async loadAccount(worldId: number, accountId: number, characterId: number) {
        const [balance, items, wishlist, account, storage, gifts] = await Promise.all([
            this.cashShopRepo.findBalance(worldId, accountId),
            this.cashShopRepo.findItems(worldId, accountId),
            this.cashShopRepo.findWishlist(worldId, characterId),
            this.accountRepo.get(worldId, accountId),
            this.storageRepo.get(worldId, accountId),
            this.cashShopRepo.takeGifts(worldId, accountId),
        ]);
        const now = Date.now();
        const locker = items.map((row) => this.toProto(row));
        const expired = locker
            .filter((item) => (item.item?.expirationUnixMs ?? 0) > 0 && (item.item?.expirationUnixMs ?? 0) <= now)
            .map((item) => item.item?.uniqueId ?? "0");
        await this.cashShopRepo.deleteItems(worldId, accountId, expired);
        return {
            nxCash: balance.nxCash,
            maplePoint: balance.maplePoint,
            locker,
            wishlist,
            characterSlotCount: account?.characterSlotCount ?? 0,
            storageSlotCount: storage?.slots ?? DEFAULT_STORAGE_SLOTS,
            gifts: gifts.map((row) => ({
                serial: row.serial,
                itemId: row.item_id,
                senderName: row.sender_name,
                message: row.message,
            })),
            expired,
        };
    }

    async buy(req: BuyCashItemRequest): Promise<BuyCashItemReply> {
        return this.cashShopRepo.withAccountLock(req.worldId, req.accountId, async (txClient) => {
            const balance = await this.cashShopRepo.findBalance(req.worldId, req.accountId, { txClient });
            const fail = (result: CashShopResult) => ({ result, nxCash: balance.nxCash, maplePoint: balance.maplePoint });
            const delta = this.debit(req.currency, req.price);
            if (this.affordable(balance, delta) === false) {
                return fail(CashShopResult.CASH_SHOP_RESULT_NOT_ENOUGH_CASH);
            }

            const items = await this.cashShopRepo.findItems(req.worldId, req.accountId, { txClient });
            if (items.length + req.items.length > LOCKER_CAPACITY) {
                return fail(CashShopResult.CASH_SHOP_RESULT_LOCKER_FULL);
            }

            await this.cashShopRepo.addBalance(req.worldId, req.accountId, delta, { txClient });
            for (const item of req.items) {
                await this.cashShopRepo.insertItem(req.worldId, req.accountId, String(item.item?.uniqueId ?? 0), CashItem.toJSON(item), { txClient });
            }
            return {
                result: CashShopResult.CASH_SHOP_RESULT_OK,
                nxCash: balance.nxCash + delta.nxCash,
                maplePoint: balance.maplePoint + delta.maplePoint,
            };
        });
    }

    async takeOut(worldId: number, accountId: number, characterId: number, item: InventoryModel): Promise<CashShopResult> {
        return this.cashShopRepo.withAccountLock(worldId, accountId, async (txClient) => {
            const row = await this.cashShopRepo.deleteItem(worldId, accountId, String(item.uniqueId), { txClient });
            if (row == null) {
                return CashShopResult.CASH_SHOP_RESULT_NOT_FOUND;
            }

            const inventory = await this.inventoryRepo.getAll(worldId, String(characterId));
            await this.inventoryRepo.replaceBySnapshot(worldId, String(characterId), [...inventory.values(), item]);
            return CashShopResult.CASH_SHOP_RESULT_OK;
        });
    }

    async putIn(worldId: number, accountId: number, characterId: number, item: CashItem): Promise<CashShopResult> {
        const serial = String(item.item?.uniqueId ?? 0);
        return this.cashShopRepo.withAccountLock(worldId, accountId, async (txClient) => {
            const items = await this.cashShopRepo.findItems(worldId, accountId, { txClient });
            if (items.length >= LOCKER_CAPACITY) {
                return CashShopResult.CASH_SHOP_RESULT_LOCKER_FULL;
            }

            const inventory = [...(await this.inventoryRepo.getAll(worldId, String(characterId))).values()];
            const kept = inventory.filter((model) => String(model.uniqueId) !== serial || model.slot < 0);
            if (kept.length === inventory.length) {
                return CashShopResult.CASH_SHOP_RESULT_NOT_FOUND;
            }

            await this.inventoryRepo.replaceBySnapshot(worldId, String(characterId), kept);
            await this.cashShopRepo.insertItem(worldId, accountId, serial, CashItem.toJSON(item), { txClient });
            return CashShopResult.CASH_SHOP_RESULT_OK;
        });
    }

    async addCash(worldId: number, accountId: number, nxCash: number, maplePoint: number) {
        return this.cashShopRepo.withAccountLock(worldId, accountId, async (txClient) => {
            const balance = await this.cashShopRepo.findBalance(worldId, accountId, { txClient });
            const delta = {
                nxCash: Math.max(nxCash, -balance.nxCash),
                maplePoint: Math.max(maplePoint, -balance.maplePoint),
            };
            await this.cashShopRepo.addBalance(worldId, accountId, delta, { txClient });
            return {
                nxCash: balance.nxCash + delta.nxCash,
                maplePoint: balance.maplePoint + delta.maplePoint,
            };
        });
    }

    async expandSlot(req: ExpandCashSlotRequest): Promise<ExpandCashSlotReply> {
        return this.cashShopRepo.withAccountLock(req.worldId, req.accountId, async (txClient) => {
            const balance = await this.cashShopRepo.findBalance(req.worldId, req.accountId, { txClient });
            const fail = (result: CashShopResult) => ({ result, slots: 0, nxCash: balance.nxCash, maplePoint: balance.maplePoint });
            const delta = this.debit(req.currency, req.price);
            if (this.affordable(balance, delta) === false) {
                return fail(CashShopResult.CASH_SHOP_RESULT_NOT_ENOUGH_CASH);
            }

            let slots = 0;
            switch (req.kind) {
                case CashSlotKind.CASH_SLOT_KIND_INVENTORY: {
                    const character = await this.characterRepo.get(req.worldId, req.characterId);
                    const index = req.inventoryType - 1;
                    if (character == null || character.slotLimits == null || index < 0 || index >= character.slotLimits.length) {
                        return fail(CashShopResult.CASH_SHOP_RESULT_NOT_FOUND);
                    }
                    slots = (character.slotLimits[index] ?? 0) + req.amount;
                    if (slots > req.max) {
                        return fail(CashShopResult.CASH_SHOP_RESULT_SLOT_LIMIT);
                    }
                    const slotLimits = [...character.slotLimits];
                    slotLimits[index] = slots;
                    await this.characterRepo.set(req.worldId, { ...character, slotLimits });
                    break;
                }
                case CashSlotKind.CASH_SLOT_KIND_STORAGE: {
                    const storage = await this.storageRepo.get(req.worldId, req.accountId);
                    slots = (storage?.slots ?? DEFAULT_STORAGE_SLOTS) + req.amount;
                    if (slots > req.max) {
                        return fail(CashShopResult.CASH_SHOP_RESULT_SLOT_LIMIT);
                    }
                    await this.storageRepo.set(req.worldId, {
                        accountId: req.accountId,
                        worldId: req.worldId,
                        slots,
                        meso: storage?.meso ?? 0,
                    });
                    break;
                }
                case CashSlotKind.CASH_SLOT_KIND_CHARACTER: {
                    const account = await this.accountRepo.get(req.worldId, req.accountId);
                    if (account == null) {
                        return fail(CashShopResult.CASH_SHOP_RESULT_NOT_FOUND);
                    }
                    slots = account.characterSlotCount + req.amount;
                    if (slots > req.max) {
                        return fail(CashShopResult.CASH_SHOP_RESULT_SLOT_LIMIT);
                    }
                    await this.accountRepo.set(req.worldId, { ...account, characterSlotCount: slots });
                    break;
                }
                default:
                    return fail(CashShopResult.CASH_SHOP_RESULT_NOT_FOUND);
            }

            await this.cashShopRepo.addBalance(req.worldId, req.accountId, delta, { txClient });
            return {
                result: CashShopResult.CASH_SHOP_RESULT_OK,
                slots,
                nxCash: balance.nxCash + delta.nxCash,
                maplePoint: balance.maplePoint + delta.maplePoint,
            };
        });
    }

    private async deliverGift(worldId: number, recipientAccountId: number, items: CashItem[], message: string): Promise<void> {
        await this.cashShopRepo.withAccount(worldId, recipientAccountId, async (txClient) => {
            for (const item of items) {
                const serial = String(item.item?.uniqueId ?? 0);
                await this.cashShopRepo.insertItem(worldId, recipientAccountId, serial, CashItem.toJSON(item), { txClient });
                await this.cashShopRepo.insertGift(
                    worldId,
                    recipientAccountId,
                    { serial, item_id: item.item?.itemId ?? 0, sender_name: item.buyerName, message },
                    { txClient }
                );
            }
        });
    }

    private async notifyGift(worldId: number, recipientId: number, senderName: string): Promise<void> {
        const channelId = await this.sessionRepo.findChannel(worldId, recipientId);
        if (channelId == null || channelId < 0) {
            return;
        }
        await this.rabbitmqService.assertDirectExchange(AMQ_DIRECT_EXCHANGE);
        await this.rabbitmqService.publish(AMQ_DIRECT_EXCHANGE, `fm.${worldId}.${channelId}.parcel`, "cash_gift_arrived", {
            occurred_at: new Date().toISOString(),
            character_id: recipientId,
            sender_name: senderName,
        });
    }

    async gift(req: GiftCashItemRequest): Promise<GiftCashItemReply> {
        const reply = await this.cashShopRepo.withAccountLock(req.worldId, req.accountId, async (txClient) => {
            const balance = await this.cashShopRepo.findBalance(req.worldId, req.accountId, { txClient });
            const fail = (result: CashShopResult) => ({ result, nxCash: balance.nxCash, maplePoint: balance.maplePoint, recipientId: 0 });
            const delta = this.debit(CashCurrency.CASH_CURRENCY_NX_CASH, req.price);
            if (this.affordable(balance, delta) === false) {
                return fail(CashShopResult.CASH_SHOP_RESULT_NOT_ENOUGH_CASH);
            }

            const entry = await this.unifiedRepo.findCharacterNameEntry(req.recipientName);
            if (entry == null || entry.world_id !== req.worldId) {
                return fail(CashShopResult.CASH_SHOP_RESULT_NOT_FOUND);
            }
            if (entry.account_id === req.accountId) {
                return fail(CashShopResult.CASH_SHOP_RESULT_SAME_ACCOUNT);
            }
            const recipient = await this.characterRepo.get(req.worldId, entry.character_id);
            if (recipient == null) {
                return fail(CashShopResult.CASH_SHOP_RESULT_NOT_FOUND);
            }
            if (req.gender !== 2 && recipient.gender !== req.gender) {
                return fail(CashShopResult.CASH_SHOP_RESULT_GENDER);
            }
            const lockerItems = await this.cashShopRepo.findItems(req.worldId, entry.account_id);
            if (lockerItems.length + req.items.length > LOCKER_CAPACITY) {
                return fail(CashShopResult.CASH_SHOP_RESULT_RECIPIENT_LOCKER_FULL);
            }

            await this.deliverGift(req.worldId, entry.account_id, req.items, req.message);
            await this.cashShopRepo.addBalance(req.worldId, req.accountId, delta, { txClient });
            return {
                result: CashShopResult.CASH_SHOP_RESULT_OK,
                nxCash: balance.nxCash + delta.nxCash,
                maplePoint: balance.maplePoint + delta.maplePoint,
                recipientId: entry.character_id,
            };
        });
        if (reply.result !== CashShopResult.CASH_SHOP_RESULT_OK) {
            return { result: reply.result, nxCash: reply.nxCash, maplePoint: reply.maplePoint };
        }

        await this.notifyGift(req.worldId, reply.recipientId, req.items[0]?.buyerName ?? "");
        return { result: reply.result, nxCash: reply.nxCash, maplePoint: reply.maplePoint };
    }

    async buyRing(req: BuyCashRingRequest): Promise<BuyCashRingReply> {
        const item = req.item!;
        const partnerItem = req.partnerItem!;
        const reply = await this.cashShopRepo.withAccountLock(req.worldId, req.accountId, async (txClient) => {
            const balance = await this.cashShopRepo.findBalance(req.worldId, req.accountId, { txClient });
            const fail = (result: CashShopResult) => ({ result, nxCash: balance.nxCash, maplePoint: balance.maplePoint, partnerId: 0 });
            const delta = this.debit(CashCurrency.CASH_CURRENCY_NX_CASH, req.price);
            if (this.affordable(balance, delta) === false) {
                return fail(CashShopResult.CASH_SHOP_RESULT_NOT_ENOUGH_CASH);
            }

            const entry = await this.unifiedRepo.findCharacterNameEntry(req.partnerName);
            if (entry == null || entry.world_id !== req.worldId || entry.character_id === req.characterId) {
                return fail(CashShopResult.CASH_SHOP_RESULT_NOT_FOUND);
            }
            if (entry.account_id === req.accountId) {
                return fail(CashShopResult.CASH_SHOP_RESULT_SAME_ACCOUNT);
            }
            const partner = await this.characterRepo.get(req.worldId, entry.character_id);
            if (partner == null) {
                return fail(CashShopResult.CASH_SHOP_RESULT_NOT_FOUND);
            }
            if (req.couple && partner.gender === req.gender) {
                return fail(CashShopResult.CASH_SHOP_RESULT_GENDER);
            }
            const lockerItems = await this.cashShopRepo.findItems(req.worldId, req.accountId, { txClient });
            if (lockerItems.length >= LOCKER_CAPACITY) {
                return fail(CashShopResult.CASH_SHOP_RESULT_LOCKER_FULL);
            }
            const partnerLockerItems = await this.cashShopRepo.findItems(req.worldId, entry.account_id);
            if (partnerLockerItems.length >= LOCKER_CAPACITY) {
                return fail(CashShopResult.CASH_SHOP_RESULT_RECIPIENT_LOCKER_FULL);
            }

            const serial = String(item.item?.uniqueId ?? 0);
            const partnerSerial = String(partnerItem.item?.uniqueId ?? 0);
            const itemId = item.item?.itemId ?? 0;
            await this.deliverGift(req.worldId, entry.account_id, [partnerItem], req.message);
            await this.cashShopRepo.insertItem(req.worldId, req.accountId, serial, CashItem.toJSON(item), { txClient });
            await this.cashShopRepo.insertRing(req.worldId, {
                serial,
                partner_serial: partnerSerial,
                character_id: req.characterId,
                partner_id: entry.character_id,
                partner_name: partner.name,
                item_id: itemId,
            });
            await this.cashShopRepo.insertRing(req.worldId, {
                serial: partnerSerial,
                partner_serial: serial,
                character_id: entry.character_id,
                partner_id: req.characterId,
                partner_name: req.characterName,
                item_id: itemId,
            });
            await this.cashShopRepo.addBalance(req.worldId, req.accountId, delta, { txClient });
            return {
                result: CashShopResult.CASH_SHOP_RESULT_OK,
                nxCash: balance.nxCash + delta.nxCash,
                maplePoint: balance.maplePoint + delta.maplePoint,
                partnerId: entry.character_id,
            };
        });
        if (reply.result !== CashShopResult.CASH_SHOP_RESULT_OK) {
            return { result: reply.result, nxCash: reply.nxCash, maplePoint: reply.maplePoint };
        }

        await this.notifyGift(req.worldId, reply.partnerId, req.characterName);
        return { result: reply.result, nxCash: reply.nxCash, maplePoint: reply.maplePoint };
    }

    async getRings(worldId: number, characterId: number): Promise<CashRing[]> {
        const rows = await this.cashShopRepo.findRings(worldId, characterId);
        return rows.map((row) => ({
            serial: row.serial,
            partnerSerial: row.partner_serial,
            partnerId: row.partner_id,
            partnerName: row.partner_name,
            itemId: row.item_id,
        }));
    }

    async payBack(req: PayBackCashItemRequest): Promise<PayBackCashItemReply> {
        return this.cashShopRepo.withAccountLock(req.worldId, req.accountId, async (txClient) => {
            const balance = await this.cashShopRepo.findBalance(req.worldId, req.accountId, { txClient });
            const row = await this.cashShopRepo.deleteItem(req.worldId, req.accountId, String(req.serial), { txClient });
            if (row == null) {
                return { result: CashShopResult.CASH_SHOP_RESULT_NOT_FOUND, nxCash: balance.nxCash, maplePoint: balance.maplePoint };
            }

            await this.cashShopRepo.addBalance(req.worldId, req.accountId, { nxCash: 0, maplePoint: req.maplePoint }, { txClient });
            return {
                result: CashShopResult.CASH_SHOP_RESULT_OK,
                nxCash: balance.nxCash,
                maplePoint: balance.maplePoint + req.maplePoint,
            };
        });
    }

    async buyQuestItem(worldId: number, characterId: number, price: number, item: InventoryModel): Promise<BuyCashQuestItemReply> {
        const character = await this.characterRepo.get(worldId, characterId);
        if (character == null) {
            return { result: CashShopResult.CASH_SHOP_RESULT_NOT_FOUND, meso: 0 };
        }
        const meso = character.meso ?? 0;
        if (meso < price) {
            return { result: CashShopResult.CASH_SHOP_RESULT_NOT_ENOUGH_MESO, meso };
        }

        const inventory = await this.inventoryRepo.getAll(worldId, String(characterId));
        await this.inventoryRepo.replaceBySnapshot(worldId, String(characterId), [...inventory.values(), item]);
        await this.characterRepo.set(worldId, { ...character, meso: meso - price });
        return { result: CashShopResult.CASH_SHOP_RESULT_OK, meso: meso - price };
    }

    async findCoupon(worldId: number, code: string): Promise<FindCashCouponReply> {
        const coupon = await this.cashShopRepo.findCoupon(worldId, code);
        if (coupon == null) {
            return { result: CashShopResult.CASH_SHOP_RESULT_NOT_FOUND, kind: CashCouponKind.CASH_COUPON_KIND_NX_CASH, value: 0 };
        }
        if (coupon.used_at != null) {
            return { result: CashShopResult.CASH_SHOP_RESULT_COUPON_USED, kind: coupon.kind, value: coupon.value };
        }
        return { result: CashShopResult.CASH_SHOP_RESULT_OK, kind: coupon.kind, value: coupon.value };
    }

    async redeemCoupon(req: RedeemCashCouponRequest): Promise<RedeemCashCouponReply> {
        const reply = await this.cashShopRepo.withCouponClaim(req.worldId, req.code, req.accountId, async (coupon) => {
            const character = await this.characterRepo.get(req.worldId, req.characterId);
            if (character == null) {
                throw new Error(`character not found: ${req.characterId}`);
            }
            let meso = character.meso ?? 0;
            const balance = await this.cashShopRepo.withAccountLock(req.worldId, req.accountId, async (txClient) => {
                const current = await this.cashShopRepo.findBalance(req.worldId, req.accountId, { txClient });
                const delta = { nxCash: 0, maplePoint: 0 };
                switch (coupon.kind) {
                    case CashCouponKind.CASH_COUPON_KIND_NX_CASH:
                        delta.nxCash = coupon.value;
                        break;
                    case CashCouponKind.CASH_COUPON_KIND_MAPLE_POINT:
                        delta.maplePoint = coupon.value;
                        break;
                    case CashCouponKind.CASH_COUPON_KIND_ITEM: {
                        if (req.item?.item == null) {
                            throw new Error(`coupon ${req.code}: item is required`);
                        }
                        const items = await this.cashShopRepo.findItems(req.worldId, req.accountId, { txClient });
                        if (items.length >= LOCKER_CAPACITY) {
                            throw Object.assign(new Error("locker full"), { result: CashShopResult.CASH_SHOP_RESULT_LOCKER_FULL });
                        }
                        await this.cashShopRepo.insertItem(req.worldId, req.accountId, String(req.item.item.uniqueId ?? 0), CashItem.toJSON(req.item), { txClient });
                        break;
                    }
                    case CashCouponKind.CASH_COUPON_KIND_MESO:
                        meso = Math.min(meso + coupon.value, 2147483647);
                        await this.characterRepo.set(req.worldId, { ...character, meso });
                        break;
                }
                await this.cashShopRepo.addBalance(req.worldId, req.accountId, delta, { txClient });
                return { nxCash: current.nxCash + delta.nxCash, maplePoint: current.maplePoint + delta.maplePoint };
            });
            return { result: CashShopResult.CASH_SHOP_RESULT_OK, ...balance, meso };
        }).catch((err: { result?: CashShopResult }) => {
            if (err.result == null) {
                throw err;
            }
            return { result: err.result, nxCash: 0, maplePoint: 0, meso: 0 };
        });
        if (reply == null) {
            return { result: CashShopResult.CASH_SHOP_RESULT_COUPON_USED, nxCash: 0, maplePoint: 0, meso: 0 };
        }
        return reply;
    }

    async createCoupons(worldId: number, kind: CashCouponKind, value: number, count: number): Promise<string[]> {
        const codes = Array.from({ length: count }, () => {
            const bytes = randomBytes(COUPON_LENGTH);
            return Array.from(bytes, (b) => COUPON_ALPHABET[b % COUPON_ALPHABET.length]).join("");
        });
        return this.cashShopRepo.insertCoupons(worldId, kind, value, codes);
    }

    async getWishlist(worldId: number, characterId: number): Promise<number[]> {
        return this.cashShopRepo.findWishlist(worldId, characterId);
    }

    async setWishlist(worldId: number, characterId: number, commoditySns: number[]): Promise<void> {
        await this.cashShopRepo.setWishlist(worldId, characterId, commoditySns.slice(0, 10));
    }
}
