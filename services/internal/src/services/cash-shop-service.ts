import { randomBytes } from "node:crypto";
import type { PoolClient } from "pg";
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
} from "@fm/shared/protobuf/generated/fminternal/internal_service";
import type { AccountRepository } from "@fm/shared/repos/account-repository";
import { CashTransferKind, type CashBalance, type CashRequestRepository, type CashShopRepository, type CashItemRow, type CashRingRow, type CashTransferRequest } from "@fm/shared/repos/cash-shop-repository";
import type { CharacterRepository } from "@fm/shared/repos/character-repository";
import type { StorageRepository } from "@fm/shared/repos/storage-repository";
import type { InventoryRepository, InventoryModel } from "@fm/shared/repos/inventory-repository";
import type { SessionRepository } from "@fm/shared/repos/session-repository";
import type { UnifiedRepository } from "@fm/shared/repos/unified-repository";
import type { RepositoryTxOptions } from "@fm/shared/types/repository-contracts";
import type { RabbitMQService } from "./rabbitmq-service";
import type { SocialService } from "./social-service";
import type { DistributedLockService } from "./distributed-lock-service";

const LOCKER_CAPACITY = 100;

class LockerFull extends Error {
    constructor() {
        super("locker full");
    }
}

enum TransferOutcome {
    Placed = 0,
    Full = 1,
    Done = 2,
    Skipped = 3,
}
const DEFAULT_STORAGE_SLOTS = 4;
const COUPON_ALPHABET = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789";
const COUPON_LENGTH = 16;
const MAX_MESO = 2147483647;

export class CashShopService {
    private readonly cashShopRepo: CashShopRepository;
    private readonly cashRequestRepo: CashRequestRepository;
    private readonly accountRepo: AccountRepository;
    private readonly storageRepo: StorageRepository;
    private readonly inventoryRepo: InventoryRepository;
    private readonly characterRepo: CharacterRepository;
    private readonly unifiedRepo: UnifiedRepository;
    private readonly sessionRepo: SessionRepository;
    private readonly rabbitmqService: RabbitMQService;
    private readonly socialService: SocialService;
    private readonly distributedLockService: DistributedLockService;

    constructor(
        cashShopRepository: CashShopRepository,
        cashRequestRepository: CashRequestRepository,
        accountRepository: AccountRepository,
        storageRepository: StorageRepository,
        inventoryRepository: InventoryRepository,
        characterRepository: CharacterRepository,
        unifiedRepository: UnifiedRepository,
        sessionRepository: SessionRepository,
        rabbitmqService: RabbitMQService,
        socialService: SocialService,
        distributedLockService: DistributedLockService
    ) {
        this.cashShopRepo = cashShopRepository;
        this.cashRequestRepo = cashRequestRepository;
        this.accountRepo = accountRepository;
        this.storageRepo = storageRepository;
        this.inventoryRepo = inventoryRepository;
        this.characterRepo = characterRepository;
        this.unifiedRepo = unifiedRepository;
        this.sessionRepo = sessionRepository;
        this.rabbitmqService = rabbitmqService;
        this.socialService = socialService;
        this.distributedLockService = distributedLockService;
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

    async settleTransfers(worldId: number): Promise<void> {
        for (const request of await this.cashShopRepo.pendingTransfers(worldId)) {
            try {
                await this.settle(worldId, request);
            } catch (err) {
                console.error(`[cash-shop] transfer pending id=${request.request_id} kind=${request.kind}`, err);
            }
        }
    }

    private async settle(worldId: number, request: CashTransferRequest): Promise<void> {
        switch (request.kind) {
            case CashTransferKind.Gift:
                await this.settleGift(worldId, request);
                return;
            case CashTransferKind.Ring:
                await this.settleRing(worldId, request);
                return;
            case CashTransferKind.TakeOut:
                await this.settleTakeOut(worldId, request);
                return;
            case CashTransferKind.PutIn:
                await this.settlePutIn(worldId, request);
                return;
            case CashTransferKind.InventorySlot:
                await this.settleSlot(worldId, request);
                return;
            case CashTransferKind.Meso:
                await this.settleMeso(worldId, request);
                return;
            default:
                throw new Error(`unknown cash transfer kind ${request.kind}`);
        }
    }

    private async receive(
        worldId: number,
        request: CashTransferRequest,
        place: (txClient: PoolClient) => Promise<TransferOutcome.Placed | TransferOutcome.Full>
    ): Promise<TransferOutcome> {
        let result: TransferOutcome = TransferOutcome.Skipped;
        try {
            await this.cashShopRepo.receiveTransfer(worldId, request, async (txClient) => {
                if (await this.cashShopRepo.isTransferCompleted(worldId, request)) {
                    result = TransferOutcome.Done;
                    return;
                }
                result = await place(txClient);
                if (result === TransferOutcome.Full) {
                    throw new LockerFull();
                }
            });
        } catch (err) {
            if (err instanceof LockerFull) {
                return TransferOutcome.Full;
            }
            throw err;
        }
        return result;
    }

    private async finishSettled(worldId: number, request: CashTransferRequest, result: TransferOutcome): Promise<boolean> {
        if (result === TransferOutcome.Done || result === TransferOutcome.Full) {
            return false;
        }
        await this.cashShopRepo.finishTransfer(worldId, request);
        return result === TransferOutcome.Placed;
    }

    private async settleGift(worldId: number, request: CashTransferRequest): Promise<void> {
        const detail = request.detail as {
            accountId: number;
            items: unknown[];
            message: string;
            recipientId: number;
            senderName: string;
        };
        const result = await this.receive(worldId, request, async (txClient) => {
            const items = await this.cashShopRepo.findItems(worldId, detail.accountId, { txClient });
            if (items.length + detail.items.length > LOCKER_CAPACITY) {
                return TransferOutcome.Full;
            }
            await this.grantGift(
                worldId,
                detail.accountId,
                detail.items.map((item) => CashItem.fromJSON(item)),
                detail.message,
                { txClient }
            );
            return TransferOutcome.Placed;
        });
        if (result === TransferOutcome.Full) {
            try {
                await this.refundGift(worldId, request);
            } catch (refundError) {
                console.error(`[cash-shop] transfer refund id=${request.request_id}`, refundError);
            }
            return;
        }
        if ((await this.finishSettled(worldId, request, result)) === false) {
            return;
        }
        await this.notifyGift(worldId, detail.recipientId, detail.senderName);
    }

    private async settleRing(worldId: number, request: CashTransferRequest): Promise<void> {
        const result = await this.receive(worldId, request, async (txClient) => {
            await this.cashShopRepo.insertRing(worldId, request.detail as CashRingRow, { txClient });
            return TransferOutcome.Placed;
        });
        await this.finishSettled(worldId, request, result);
    }

    private async settleTakeOut(worldId: number, request: CashTransferRequest): Promise<void> {
        const result = await this.receive(worldId, request, async (txClient) => {
            const item = request.detail as InventoryModel;
            if (item.expiration != null) {
                item.expiration = new Date(item.expiration);
            }
            const inventory = await this.inventoryRepo.getAll(worldId, String(request.target_id), { txClient });
            await this.inventoryRepo.replaceBySnapshot(worldId, String(request.target_id), [...inventory.values(), item], { txClient });
            return TransferOutcome.Placed;
        });
        if ((await this.finishSettled(worldId, request, result)) === false) {
            return;
        }
        await this.inventoryRepo.invalidateCache(worldId, String(request.target_id));
    }

    private async settlePutIn(worldId: number, request: CashTransferRequest): Promise<void> {
        const detail = request.detail as { accountId: number; serial: string; item: unknown };
        const result = await this.receive(worldId, request, async (txClient) => {
            const items = await this.cashShopRepo.findItems(worldId, detail.accountId, { txClient });
            if (items.length >= LOCKER_CAPACITY) {
                return TransferOutcome.Full;
            }
            await this.cashShopRepo.insertItem(worldId, detail.accountId, detail.serial, detail.item, { txClient });
            return TransferOutcome.Placed;
        });
        if (result === TransferOutcome.Full) {
            try {
                await this.refundPutIn(worldId, request);
            } catch (refundError) {
                console.error(`[cash-shop] transfer refund id=${request.request_id}`, refundError);
            }
            return;
        }
        await this.finishSettled(worldId, request, result);
    }

    private async settleSlot(worldId: number, request: CashTransferRequest): Promise<void> {
        const detail = request.detail as { index: number; amount: number; max: number };
        const result = await this.receive(worldId, request, async (txClient) => {
            const character = await this.characterRepo.get(worldId, request.target_id, { txClient });
            if (character == null || character.slotLimits == null) {
                throw new Error(`character not found: ${request.target_id}`);
            }
            const slotLimits = [...character.slotLimits];
            slotLimits[detail.index] = Math.min((slotLimits[detail.index] ?? 0) + detail.amount, detail.max);
            await this.characterRepo.set(worldId, { ...character, slotLimits }, { txClient });
            return TransferOutcome.Placed;
        });
        await this.finishSettled(worldId, request, result);
    }

    private async settleMeso(worldId: number, request: CashTransferRequest): Promise<void> {
        const detail = request.detail as { meso: number };
        const result = await this.receive(worldId, request, async (txClient) => {
            const character = await this.characterRepo.get(worldId, request.target_id, { txClient });
            if (character == null) {
                throw new Error(`character not found: ${request.target_id}`);
            }
            await this.characterRepo.set(worldId, { ...character, meso: Math.min((character.meso ?? 0) + detail.meso, MAX_MESO) }, { txClient });
            return TransferOutcome.Placed;
        });
        await this.finishSettled(worldId, request, result);
    }

    private async refundGift(worldId: number, request: CashTransferRequest): Promise<void> {
        const detail = request.detail as {
            price?: number;
            buyerSerial?: string;
            rings?: { characterId: number; serial: string }[];
        };
        if (typeof detail.price !== "number") {
            throw new Error(`gift transfer ${request.request_id} has no price`);
        }
        const price = detail.price;
        const rings = detail.rings ?? [];
        for (const ring of rings) {
            await this.cashShopRepo.deleteRing(worldId, ring.characterId, ring.serial);
        }
        await this.cashShopRepo.withAccountLock(worldId, request.source_id, async (txClient) => {
            const claimed = await this.cashShopRepo.claimTransfer(request.request_id, txClient);
            if (!claimed) {
                return;
            }
            await this.cashShopRepo.addBalance(worldId, request.source_id, { nxCash: price, maplePoint: 0 }, { txClient });
            if (detail.buyerSerial != null && detail.buyerSerial !== "") {
                await this.cashShopRepo.deleteItem(worldId, request.source_id, detail.buyerSerial, { txClient });
            }
            await this.cashShopRepo.completeRingTransfers(
                request.source_id,
                rings.map((ring) => ring.serial),
                txClient
            );
        });
        for (const ring of rings) {
            await this.cashShopRepo.deleteRing(worldId, ring.characterId, ring.serial);
        }
    }

    private async refundPutIn(worldId: number, request: CashTransferRequest): Promise<void> {
        const detail = request.detail as { serial: string; inventory?: InventoryModel };
        if (detail.inventory == null) {
            throw new Error(`put_in transfer ${request.request_id} has no item`);
        }
        const item = detail.inventory;
        if (item.expiration != null) {
            item.expiration = new Date(item.expiration);
        }
        await this.cashShopRepo.withCharacterLock(worldId, request.source_id, async (txClient) => {
            const claimed = await this.cashShopRepo.claimTransfer(request.request_id, txClient);
            if (!claimed) {
                return;
            }
            const inventory = [...(await this.inventoryRepo.getAll(worldId, String(request.source_id), { txClient })).values()];
            if (inventory.some((model) => String(model.uniqueId) === detail.serial)) {
                return;
            }
            const used = new Set(inventory.filter((model) => model.inventoryType === item.inventoryType).map((model) => model.slot));
            if (used.has(item.slot)) {
                let slot = 1;
                while (used.has(slot)) {
                    slot++;
                }
                item.slot = slot;
            }
            await this.inventoryRepo.replaceBySnapshot(worldId, String(request.source_id), [...inventory, item], { txClient });
        });
        await this.inventoryRepo.invalidateCache(worldId, String(request.source_id));
    }

    async loadAccount(worldId: number, accountId: number, characterId: number) {
        const [balance, items, wishlist, account, storage, gifts] = await Promise.all([
            this.cashShopRepo.findBalance(worldId, accountId),
            this.cashShopRepo.findItems(worldId, accountId),
            this.cashShopRepo.findWishlist(worldId, characterId),
            this.accountRepo.get(worldId, accountId),
            this.storageRepo.get(worldId, accountId),
            this.cashShopRepo.findGifts(worldId, accountId),
        ]);
        const now = Date.now();
        const locker = items.map((row) => this.toProto(row));
        const expired = locker
            .filter((item) => (item.item?.expirationUnixMs ?? 0) > 0 && (item.item?.expirationUnixMs ?? 0) <= now)
            .map((item) => item.item?.uniqueId ?? "0");
        await this.cashShopRepo.deleteItems(worldId, accountId, expired);
        const visibleGifts = gifts.filter((gift) => !expired.includes(gift.serial));
        for (const gift of visibleGifts) {
            await this.socialService.gift(worldId, characterId, gift.sender_name, gift.serial);
        }
        return {
            nxCash: balance.nxCash,
            maplePoint: balance.maplePoint,
            locker,
            wishlist,
            characterSlotCount: account?.characterSlotCount ?? 0,
            storageSlotCount: storage?.slots ?? DEFAULT_STORAGE_SLOTS,
            gifts: visibleGifts.map((row) => ({
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
        const result = await this.cashShopRepo.withAccountLock(worldId, accountId, async (txClient) => {
            const row = await this.cashShopRepo.deleteItem(worldId, accountId, String(item.uniqueId), { txClient });
            if (row == null) {
                return CashShopResult.CASH_SHOP_RESULT_NOT_FOUND;
            }

            await this.cashShopRepo.deleteGift(worldId, accountId, String(item.uniqueId), { txClient });
            await this.cashShopRepo.queueTransfer(worldId, accountId, characterId, CashTransferKind.TakeOut, item, txClient);
            return CashShopResult.CASH_SHOP_RESULT_OK;
        });
        if (result === CashShopResult.CASH_SHOP_RESULT_OK) {
            await this.settleTransfers(worldId);
        }
        return result;
    }

    async putIn(worldId: number, accountId: number, characterId: number, item: CashItem): Promise<CashShopResult> {
        const serial = String(item.item?.uniqueId ?? 0);
        const result = await this.cashShopRepo.withCharacterLock(worldId, characterId, async (txClient) => {
            const items = await this.cashShopRepo.findItems(worldId, accountId);
            if (items.length >= LOCKER_CAPACITY) {
                return CashShopResult.CASH_SHOP_RESULT_LOCKER_FULL;
            }

            const inventory = [...(await this.inventoryRepo.getAll(worldId, String(characterId), { txClient })).values()];
            const kept = inventory.filter((model) => String(model.uniqueId) !== serial || model.slot < 0);
            if (kept.length === inventory.length) {
                return CashShopResult.CASH_SHOP_RESULT_NOT_FOUND;
            }

            await this.inventoryRepo.replaceBySnapshot(worldId, String(characterId), kept, { txClient });
            const removed = inventory.find((model) => String(model.uniqueId) === serial && model.slot >= 0);
            await this.cashShopRepo.queueTransfer(worldId, characterId, accountId, CashTransferKind.PutIn, {
                accountId,
                serial,
                item: CashItem.toJSON(item),
                inventory: removed,
            }, txClient);
            return CashShopResult.CASH_SHOP_RESULT_OK;
        });
        await this.inventoryRepo.invalidateCache(worldId, String(characterId));
        if (result === CashShopResult.CASH_SHOP_RESULT_OK) {
            await this.settleTransfers(worldId);
        }
        return result;
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
        try {
            const reply = await this.cashShopRepo.withAccountLock(req.worldId, req.accountId, async (txClient) => {
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
                        await this.cashShopRepo.queueTransfer(
                            req.worldId,
                            req.accountId,
                            req.characterId,
                            CashTransferKind.InventorySlot,
                            { index, amount: req.amount, max: req.max },
                            txClient
                        );
                        break;
                    }
                    case CashSlotKind.CASH_SLOT_KIND_STORAGE: {
                        const storage = await this.storageRepo.get(req.worldId, req.accountId, { txClient });
                        slots = (storage?.slots ?? DEFAULT_STORAGE_SLOTS) + req.amount;
                        if (slots > req.max) {
                            return fail(CashShopResult.CASH_SHOP_RESULT_SLOT_LIMIT);
                        }
                        await this.storageRepo.set(
                            req.worldId,
                            {
                                accountId: req.accountId,
                                worldId: req.worldId,
                                slots,
                                meso: storage?.meso ?? 0,
                            },
                            { txClient }
                        );
                        break;
                    }
                    case CashSlotKind.CASH_SLOT_KIND_CHARACTER: {
                        const account = await this.accountRepo.get(req.worldId, req.accountId, { txClient });
                        if (account == null) {
                            return fail(CashShopResult.CASH_SHOP_RESULT_NOT_FOUND);
                        }
                        slots = account.characterSlotCount + req.amount;
                        if (slots > req.max) {
                            return fail(CashShopResult.CASH_SHOP_RESULT_SLOT_LIMIT);
                        }
                        await this.accountRepo.setCharacterSlotCount(req.worldId, req.accountId, slots, { txClient });
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
            if (reply.result === CashShopResult.CASH_SHOP_RESULT_OK && req.kind === CashSlotKind.CASH_SLOT_KIND_INVENTORY) {
                await this.settleTransfers(req.worldId);
            }
            return reply;
        } finally {
            switch (req.kind) {
                case CashSlotKind.CASH_SLOT_KIND_INVENTORY:
                    await this.characterRepo.invalidateCache(req.worldId, req.characterId);
                    break;
                case CashSlotKind.CASH_SLOT_KIND_STORAGE:
                    await this.storageRepo.invalidateCache(req.worldId, req.accountId);
                    break;
                case CashSlotKind.CASH_SLOT_KIND_CHARACTER:
                    await this.accountRepo.invalidateCache(req.worldId, req.accountId);
                    break;
            }
        }
    }

    private async grantGift(
        worldId: number,
        recipientAccountId: number,
        items: CashItem[],
        message: string,
        options: RepositoryTxOptions
    ): Promise<void> {
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
        }, options);
    }

    // The gift is already committed; a failed notice must not report the purchase as failed.
    private async notifyGift(worldId: number, recipientId: number, senderName: string): Promise<void> {
        try {
            const channelId = await this.sessionRepo.findChannel(worldId, recipientId);
            if (channelId == null || channelId < 0) {
                return;
            }
            await this.rabbitmqService.publish(`fm.${worldId}.${channelId}.parcel`, "cash_gift_arrived", {
                character_id: recipientId,
                sender_name: senderName,
            });
        } catch (err) {
            console.error("[cash-shop] gift notice failed", err);
        }
    }

    private async withAccounts<T>(
        worldId: number,
        senderId: number,
        recipientId: number,
        fn: (senderTx: PoolClient) => Promise<T>
    ): Promise<T> {
        await using _locks = await this.distributedLockService.acquireWorldDataLocks(worldId, [
            `cash_account:${senderId}`,
            `cash_account:${recipientId}`,
        ]);
        return this.cashShopRepo.withAccountLock(worldId, senderId, fn);
    }

    private async reservedLockerCount(worldId: number, accountId: number, txClient: PoolClient): Promise<number> {
        const items = await this.cashShopRepo.findItems(worldId, accountId, { txClient });
        return items.length + (await this.cashShopRepo.pendingGiftItemCount(worldId, accountId));
    }

    async gift(req: GiftCashItemRequest): Promise<GiftCashItemReply> {
        const opened = await this.cashShopRepo.findBalance(req.worldId, req.accountId);
        const fail = (result: CashShopResult, balance: CashBalance = opened) => ({
            result,
            nxCash: balance.nxCash,
            maplePoint: balance.maplePoint,
            recipientId: 0,
        });
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

        const reply = await this.withAccounts(req.worldId, req.accountId, entry.account_id, async (txClient) => {
            const balance = await this.cashShopRepo.findBalance(req.worldId, req.accountId, { txClient });
            if ((await this.cashRequestRepo.find(req.worldId, req.accountId, req.requestId, { txClient })) != null) {
                return { result: CashShopResult.CASH_SHOP_RESULT_OK, nxCash: balance.nxCash, maplePoint: balance.maplePoint, recipientId: 0 };
            }
            const used = await this.reservedLockerCount(req.worldId, entry.account_id, txClient);
            if (used + req.items.length > LOCKER_CAPACITY) {
                return fail(CashShopResult.CASH_SHOP_RESULT_RECIPIENT_LOCKER_FULL, balance);
            }
            const delta = this.debit(CashCurrency.CASH_CURRENCY_NX_CASH, req.price);
            if (this.affordable(balance, delta) === false) {
                return fail(CashShopResult.CASH_SHOP_RESULT_NOT_ENOUGH_CASH, balance);
            }

            await this.cashRequestRepo.insert(req.worldId, req.accountId, req.requestId, {}, txClient);
            await this.cashShopRepo.addBalance(req.worldId, req.accountId, delta, { txClient });
            await this.cashShopRepo.queueTransfer(req.worldId, req.accountId, entry.account_id, CashTransferKind.Gift, {
                accountId: entry.account_id,
                items: req.items.map((item) => CashItem.toJSON(item)),
                message: req.message,
                recipientId: entry.character_id,
                senderName: req.items[0]?.buyerName ?? "",
                price: req.price,
            }, txClient);
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

        await this.settleTransfers(req.worldId);
        return { result: reply.result, nxCash: reply.nxCash, maplePoint: reply.maplePoint };
    }

    async buyRing(req: BuyCashRingRequest): Promise<BuyCashRingReply> {
        const item = req.item!;
        const partnerItem = req.partnerItem!;
        const opened = await this.cashShopRepo.findBalance(req.worldId, req.accountId);
        const early = (result: CashShopResult) => ({ result, nxCash: opened.nxCash, maplePoint: opened.maplePoint, partnerId: 0 });
        const entry = await this.unifiedRepo.findCharacterNameEntry(req.partnerName);
        if (entry == null || entry.world_id !== req.worldId || entry.character_id === req.characterId) {
            return early(CashShopResult.CASH_SHOP_RESULT_NOT_FOUND);
        }
        if (entry.account_id === req.accountId) {
            return early(CashShopResult.CASH_SHOP_RESULT_SAME_ACCOUNT);
        }
        const partner = await this.characterRepo.get(req.worldId, entry.character_id);
        if (partner == null) {
            return early(CashShopResult.CASH_SHOP_RESULT_NOT_FOUND);
        }
        if (req.couple && partner.gender === req.gender) {
            return early(CashShopResult.CASH_SHOP_RESULT_GENDER);
        }

        const reply = await this.withAccounts(req.worldId, req.accountId, entry.account_id, async (txClient) => {
            const balance = await this.cashShopRepo.findBalance(req.worldId, req.accountId, { txClient });
            const fail = (result: CashShopResult) => ({ result, nxCash: balance.nxCash, maplePoint: balance.maplePoint, partnerId: 0 });
            if ((await this.cashRequestRepo.find(req.worldId, req.accountId, req.requestId, { txClient })) != null) {
                return { result: CashShopResult.CASH_SHOP_RESULT_OK, nxCash: balance.nxCash, maplePoint: balance.maplePoint, partnerId: 0 };
            }
            const delta = this.debit(CashCurrency.CASH_CURRENCY_NX_CASH, req.price);
            if (this.affordable(balance, delta) === false) {
                return fail(CashShopResult.CASH_SHOP_RESULT_NOT_ENOUGH_CASH);
            }
            if ((await this.reservedLockerCount(req.worldId, req.accountId, txClient)) >= LOCKER_CAPACITY) {
                return fail(CashShopResult.CASH_SHOP_RESULT_LOCKER_FULL);
            }
            if ((await this.reservedLockerCount(req.worldId, entry.account_id, txClient)) + 1 > LOCKER_CAPACITY) {
                return fail(CashShopResult.CASH_SHOP_RESULT_RECIPIENT_LOCKER_FULL);
            }

            const serial = String(item.item?.uniqueId ?? 0);
            const partnerSerial = String(partnerItem.item?.uniqueId ?? 0);
            const itemId = item.item?.itemId ?? 0;
            await this.cashRequestRepo.insert(req.worldId, req.accountId, req.requestId, {}, txClient);
            await this.cashShopRepo.insertItem(req.worldId, req.accountId, serial, CashItem.toJSON(item), { txClient });
            await this.cashShopRepo.queueTransfer(req.worldId, req.accountId, entry.account_id, CashTransferKind.Gift, {
                accountId: entry.account_id,
                items: [CashItem.toJSON(partnerItem)],
                message: req.message,
                recipientId: entry.character_id,
                senderName: req.characterName,
                price: req.price,
                buyerSerial: serial,
                rings: [
                    { characterId: req.characterId, serial },
                    { characterId: entry.character_id, serial: partnerSerial },
                ],
            }, txClient);
            await this.cashShopRepo.queueTransfer(req.worldId, req.accountId, req.characterId, CashTransferKind.Ring, {
                serial,
                partner_serial: partnerSerial,
                character_id: req.characterId,
                partner_id: entry.character_id,
                partner_name: partner.name,
                item_id: itemId,
            }, txClient);
            await this.cashShopRepo.queueTransfer(req.worldId, req.accountId, entry.character_id, CashTransferKind.Ring, {
                serial: partnerSerial,
                partner_serial: serial,
                character_id: entry.character_id,
                partner_id: req.characterId,
                partner_name: req.characterName,
                item_id: itemId,
            }, txClient);
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

        await this.settleTransfers(req.worldId);
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
        const reply = await this.cashShopRepo.withCharacterLock(worldId, characterId, async (txClient) => {
            const character = await this.characterRepo.get(worldId, characterId, { txClient });
            if (character == null) {
                return { result: CashShopResult.CASH_SHOP_RESULT_NOT_FOUND, meso: 0 };
            }
            const meso = character.meso ?? 0;
            if (meso < price) {
                return { result: CashShopResult.CASH_SHOP_RESULT_NOT_ENOUGH_MESO, meso };
            }

            const inventory = await this.inventoryRepo.getAll(worldId, String(characterId), { txClient });
            await this.inventoryRepo.replaceBySnapshot(worldId, String(characterId), [...inventory.values(), item], { txClient });
            await this.characterRepo.set(worldId, { ...character, meso: meso - price }, { txClient });
            return { result: CashShopResult.CASH_SHOP_RESULT_OK, meso: meso - price };
        });
        await this.inventoryRepo.invalidateCache(worldId, String(characterId));
        return reply;
    }

    async findCoupon(worldId: number, code: string, accountId: number): Promise<FindCashCouponReply> {
        const coupon = await this.cashShopRepo.findCoupon(worldId, code);
        if (coupon == null) {
            return { result: CashShopResult.CASH_SHOP_RESULT_NOT_FOUND, kind: CashCouponKind.CASH_COUPON_KIND_NX_CASH, value: 0 };
        }
        if (
            coupon.used_at != null &&
            (coupon.used_by !== accountId || coupon.request_id == null || (await this.cashRequestRepo.find(worldId, accountId, coupon.request_id)) != null)
        ) {
            return { result: CashShopResult.CASH_SHOP_RESULT_COUPON_USED, kind: coupon.kind, value: coupon.value };
        }
        return { result: CashShopResult.CASH_SHOP_RESULT_OK, kind: coupon.kind, value: coupon.value };
    }

    // The claim commits before the grant and keeps the first request id, so a retry by this account grants under that id at most once.
    async redeemCoupon(req: RedeemCashCouponRequest): Promise<RedeemCashCouponReply> {
        const used = { result: CashShopResult.CASH_SHOP_RESULT_COUPON_USED, nxCash: 0, maplePoint: 0, meso: 0 };
        const coupon = await this.cashShopRepo.claimCoupon(req.worldId, req.code, req.accountId, req.requestId);
        if (coupon == null || coupon.request_id == null) {
            return used;
        }
        const requestId = coupon.request_id;
        const character = await this.characterRepo.get(req.worldId, req.characterId);
        if (character == null) {
            throw new Error(`character not found: ${req.characterId}`);
        }

        const reply = await this.cashShopRepo.withAccountLock(req.worldId, req.accountId, async (txClient) => {
            const current = await this.cashShopRepo.findBalance(req.worldId, req.accountId, { txClient });
            if (coupon.kind === CashCouponKind.CASH_COUPON_KIND_ITEM) {
                const items = await this.cashShopRepo.findItems(req.worldId, req.accountId, { txClient });
                if (items.length >= LOCKER_CAPACITY) {
                    return { result: CashShopResult.CASH_SHOP_RESULT_LOCKER_FULL, nxCash: 0, maplePoint: 0, meso: 0 };
                }
            }
            if ((await this.cashRequestRepo.reserve(req.worldId, req.accountId, requestId, {}, txClient)) === false) {
                return used;
            }

            const delta = { nxCash: 0, maplePoint: 0 };
            let meso = character.meso ?? 0;
            switch (coupon.kind) {
                case CashCouponKind.CASH_COUPON_KIND_NX_CASH:
                    delta.nxCash = coupon.value;
                    break;
                case CashCouponKind.CASH_COUPON_KIND_MAPLE_POINT:
                    delta.maplePoint = coupon.value;
                    break;
                case CashCouponKind.CASH_COUPON_KIND_ITEM:
                    if (req.item?.item == null) {
                        throw new Error(`coupon ${req.code}: item is required`);
                    }
                    await this.cashShopRepo.insertItem(req.worldId, req.accountId, String(req.item.item.uniqueId ?? 0), CashItem.toJSON(req.item), { txClient });
                    break;
                case CashCouponKind.CASH_COUPON_KIND_MESO:
                    meso = Math.min(meso + coupon.value, MAX_MESO);
                    await this.cashShopRepo.queueTransfer(req.worldId, req.accountId, req.characterId, CashTransferKind.Meso, { meso: coupon.value }, txClient);
                    break;
            }
            await this.cashShopRepo.addBalance(req.worldId, req.accountId, delta, { txClient });
            return {
                result: CashShopResult.CASH_SHOP_RESULT_OK,
                nxCash: current.nxCash + delta.nxCash,
                maplePoint: current.maplePoint + delta.maplePoint,
                meso,
            };
        });
        if (reply.result === CashShopResult.CASH_SHOP_RESULT_OK && coupon.kind === CashCouponKind.CASH_COUPON_KIND_MESO) {
            await this.settleTransfers(req.worldId);
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
