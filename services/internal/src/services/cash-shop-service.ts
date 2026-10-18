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
import type { AccountWorldRepository } from "@fm/shared/repos/account-world-repository";
import {
    CashTransferKind,
    type CashBalance,
    type CashBalanceRequestRepository,
    type CashRequestRepository,
    type CashShopRepository,
    type CashItemRow,
    type CashRingRow,
    type CashTransferRequest,
} from "@fm/shared/repos/cash-shop-repository";
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
const NO_REFUND: CashBalance = { nxCash: 0, maplePoint: 0 };

type BalanceTransfer = { targetId: number; kind: CashTransferKind; detail: unknown };

export class CashShopService {
    private readonly cashShopRepo: CashShopRepository;
    private readonly cashRequestRepo: CashRequestRepository;
    private readonly cashBalanceRequestRepo: CashBalanceRequestRepository;
    private readonly accountWorldRepo: AccountWorldRepository;
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
        cashBalanceRequestRepository: CashBalanceRequestRepository,
        accountWorldRepository: AccountWorldRepository,
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
        this.cashBalanceRequestRepo = cashBalanceRequestRepository;
        this.accountWorldRepo = accountWorldRepository;
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

    private static negate(delta: CashBalance): CashBalance {
        return { nxCash: -delta.nxCash, maplePoint: -delta.maplePoint };
    }

    // The balance is account-wide in the unified database, while what it pays for lives in a world's shards. The
    // change and the requests for what it bought commit together; settle places them in the shards once per request
    // id, or refunds them when they no longer fit. A request id makes a retried call return without charging again.
    private async charge(
        worldId: number,
        accountId: number,
        delta: CashBalance,
        requestId: string | null,
        transfers: BalanceTransfer[]
    ): Promise<{ result: CashShopResult; balance: CashBalance }> {
        return this.cashShopRepo.withBalanceLock(accountId, async (txClient) => {
            const balance = await this.cashShopRepo.findBalance(accountId, { txClient });
            if (requestId !== null && (await this.cashBalanceRequestRepo.find(worldId, accountId, requestId, { txClient })) != null) {
                return { result: CashShopResult.CASH_SHOP_RESULT_OK, balance };
            }
            if (this.affordable(balance, delta) === false) {
                return { result: CashShopResult.CASH_SHOP_RESULT_NOT_ENOUGH_CASH, balance };
            }
            if (requestId !== null) {
                await this.cashBalanceRequestRepo.insert(worldId, accountId, requestId, {}, txClient);
            }
            await this.cashShopRepo.addBalance(accountId, delta, { txClient });
            for (const transfer of transfers) {
                await this.cashShopRepo.queueBalanceTransfer(worldId, accountId, transfer.targetId, transfer.kind, transfer.detail, txClient);
            }
            return {
                result: CashShopResult.CASH_SHOP_RESULT_OK,
                balance: { nxCash: balance.nxCash + delta.nxCash, maplePoint: balance.maplePoint + delta.maplePoint },
            };
        });
    }

    // A request without a refund was queued on a shard (the buyer's half of a ring, already paid by the gift); it is
    // closed there and the item is lost.
    private async refundBalance(request: CashTransferRequest, refund: CashBalance): Promise<void> {
        if (refund.nxCash === 0 && refund.maplePoint === 0) {
            console.error(`[cash-shop] transfer id=${request.request_id} did not fit and has no refund`);
            await this.cashShopRepo.finishTransfer(request);
            return;
        }
        await this.cashShopRepo.withBalanceLock(request.source_id, async (txClient) => {
            if (await this.cashShopRepo.claimTransfer(request.request_id, txClient)) {
                await this.cashShopRepo.addBalance(request.source_id, refund, { txClient });
            }
        });
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
            case CashTransferKind.Purchase:
                await this.settlePurchase(worldId, request);
                return;
            case CashTransferKind.StorageSlot:
                await this.settleStorageSlot(worldId, request);
                return;
            case CashTransferKind.CharacterSlot:
                await this.settleCharacterSlot(worldId, request);
                return;
            case CashTransferKind.Credit:
                await this.settleCredit(request);
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
                if (await this.cashShopRepo.isTransferCompleted(request)) {
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

    private async finishSettled(request: CashTransferRequest, result: TransferOutcome): Promise<boolean> {
        if (result === TransferOutcome.Done || result === TransferOutcome.Full) {
            return false;
        }
        await this.cashShopRepo.finishTransfer(request);
        return result === TransferOutcome.Placed;
    }

    private async settleGift(worldId: number, request: CashTransferRequest): Promise<void> {
        const detail = request.detail as {
            accountId: number;
            items: unknown[];
            message: string;
            recipientId: number;
            senderName: string;
            buyerItem?: { accountId: number; item: unknown };
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
            // The buyer's half of a ring is queued only once the partner's half is placed, so a refund never leaves
            // the buyer holding it.
            if (detail.buyerItem != null) {
                await this.cashShopRepo.queueTransfer(detail.accountId, detail.buyerItem.accountId, CashTransferKind.Purchase, {
                    accountId: detail.buyerItem.accountId,
                    items: [detail.buyerItem.item],
                    refund: NO_REFUND,
                }, txClient);
            }
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
        if ((await this.finishSettled(request, result)) === false) {
            return;
        }
        await this.notifyGift(worldId, detail.recipientId, detail.senderName);
    }

    private async settleRing(worldId: number, request: CashTransferRequest): Promise<void> {
        const result = await this.receive(worldId, request, async (txClient) => {
            await this.cashShopRepo.insertRing(worldId, request.detail as CashRingRow, { txClient });
            return TransferOutcome.Placed;
        });
        await this.finishSettled(request, result);
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
        if ((await this.finishSettled(request, result)) === false) {
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
        await this.finishSettled(request, result);
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
        await this.finishSettled(request, result);
    }

    private async settlePurchase(worldId: number, request: CashTransferRequest): Promise<void> {
        const detail = request.detail as { accountId: number; items: unknown[]; refund: CashBalance };
        const result = await this.receive(worldId, request, async (txClient) => {
            const items = await this.cashShopRepo.findItems(worldId, detail.accountId, { txClient });
            if (items.length + detail.items.length > LOCKER_CAPACITY) {
                return TransferOutcome.Full;
            }
            for (const item of detail.items) {
                const cashItem = CashItem.fromJSON(item);
                await this.cashShopRepo.insertItem(worldId, detail.accountId, String(cashItem.item?.uniqueId ?? 0), CashItem.toJSON(cashItem), { txClient });
            }
            return TransferOutcome.Placed;
        });
        if (result === TransferOutcome.Full) {
            await this.refundBalance(request, detail.refund);
            return;
        }
        await this.finishSettled(request, result);
    }

    private async settleStorageSlot(worldId: number, request: CashTransferRequest): Promise<void> {
        const detail = request.detail as { amount: number; max: number };
        const result = await this.receive(worldId, request, async (txClient) => {
            const storage = await this.storageRepo.get(worldId, request.target_id, { txClient });
            await this.storageRepo.set(
                worldId,
                {
                    accountId: request.target_id,
                    worldId,
                    slots: Math.min((storage?.slots ?? DEFAULT_STORAGE_SLOTS) + detail.amount, detail.max),
                    meso: storage?.meso ?? 0,
                },
                { txClient }
            );
            return TransferOutcome.Placed;
        });
        await this.finishSettled(request, result);
        await this.storageRepo.invalidateCache(worldId, request.target_id);
    }

    private async settleCharacterSlot(worldId: number, request: CashTransferRequest): Promise<void> {
        const detail = request.detail as { amount: number; max: number };
        const result = await this.receive(worldId, request, async (txClient) => {
            const slots = await this.accountWorldRepo.findCharacterSlotCount(worldId, request.target_id, { txClient });
            await this.accountWorldRepo.setCharacterSlotCount(worldId, request.target_id, Math.min(slots + detail.amount, detail.max), { txClient });
            return TransferOutcome.Placed;
        });
        await this.finishSettled(request, result);
    }

    // A credit lands in the unified balance, recorded under the request id in the same transaction.
    private async settleCredit(request: CashTransferRequest): Promise<void> {
        const detail = request.detail as { delta: CashBalance };
        await this.cashShopRepo.withBalanceLock(request.target_id, async (txClient) => {
            if (await this.cashBalanceRequestRepo.reserve(0, request.target_id, request.request_id, {}, txClient)) {
                await this.cashShopRepo.addBalance(request.target_id, detail.delta, { txClient });
            }
        });
        await this.cashShopRepo.finishTransfer(request);
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
        await this.finishSettled(request, result);
    }

    private async refundGift(worldId: number, request: CashTransferRequest): Promise<void> {
        const detail = request.detail as {
            price?: number;
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
        await this.cashShopRepo.withBalanceLock(request.source_id, async (txClient) => {
            const claimed = await this.cashShopRepo.claimTransfer(request.request_id, txClient);
            if (!claimed) {
                return;
            }
            await this.cashShopRepo.addBalance(request.source_id, { nxCash: price, maplePoint: 0 }, { txClient });
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
        const [balance, items, wishlist, characterSlotCount, storage, gifts] = await Promise.all([
            this.cashShopRepo.findBalance(accountId),
            this.cashShopRepo.findItems(worldId, accountId),
            this.cashShopRepo.findWishlist(worldId, characterId),
            this.accountWorldRepo.findCharacterSlotCount(worldId, accountId),
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
            characterSlotCount,
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
        await using _lock = await this.distributedLockService.acquireWorldDataLock(req.worldId, `cash_account:${req.accountId}`);
        if ((await this.reservedLockerCount(req.worldId, req.accountId)) + req.items.length > LOCKER_CAPACITY) {
            const balance = await this.cashShopRepo.findBalance(req.accountId);
            return { result: CashShopResult.CASH_SHOP_RESULT_LOCKER_FULL, nxCash: balance.nxCash, maplePoint: balance.maplePoint };
        }

        const delta = this.debit(req.currency, req.price);
        const charged = await this.charge(req.worldId, req.accountId, delta, null, [{
            targetId: req.accountId,
            kind: CashTransferKind.Purchase,
            detail: { accountId: req.accountId, items: req.items.map((item) => CashItem.toJSON(item)), refund: CashShopService.negate(delta) },
        }]);
        if (charged.result === CashShopResult.CASH_SHOP_RESULT_OK) {
            await this.settleTransfers(req.worldId);
        }
        return { result: charged.result, nxCash: charged.balance.nxCash, maplePoint: charged.balance.maplePoint };
    }

    async takeOut(worldId: number, accountId: number, characterId: number, item: InventoryModel): Promise<CashShopResult> {
        const result = await this.cashShopRepo.withAccountLock(worldId, accountId, async (txClient) => {
            const row = await this.cashShopRepo.deleteItem(worldId, accountId, String(item.uniqueId), { txClient });
            if (row == null) {
                return CashShopResult.CASH_SHOP_RESULT_NOT_FOUND;
            }

            await this.cashShopRepo.deleteGift(worldId, accountId, String(item.uniqueId), { txClient });
            await this.cashShopRepo.queueTransfer(accountId, characterId, CashTransferKind.TakeOut, item, txClient);
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
            await this.cashShopRepo.queueTransfer(characterId, accountId, CashTransferKind.PutIn, {
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

    async addCash(accountId: number, nxCash: number, maplePoint: number) {
        return this.cashShopRepo.withBalanceLock(accountId, async (txClient) => {
            const balance = await this.cashShopRepo.findBalance(accountId, { txClient });
            const delta = {
                nxCash: Math.max(nxCash, -balance.nxCash),
                maplePoint: Math.max(maplePoint, -balance.maplePoint),
            };
            await this.cashShopRepo.addBalance(accountId, delta, { txClient });
            return {
                nxCash: balance.nxCash + delta.nxCash,
                maplePoint: balance.maplePoint + delta.maplePoint,
            };
        });
    }

    async expandSlot(req: ExpandCashSlotRequest): Promise<ExpandCashSlotReply> {
        try {
            const delta = this.debit(req.currency, req.price);
            const detail = { amount: req.amount, max: req.max, refund: CashShopService.negate(delta) };
            let slots = 0;
            let transfer: BalanceTransfer;
            switch (req.kind) {
                case CashSlotKind.CASH_SLOT_KIND_INVENTORY: {
                    const character = await this.characterRepo.get(req.worldId, req.characterId);
                    const index = req.inventoryType - 1;
                    if (character == null || character.slotLimits == null || index < 0 || index >= character.slotLimits.length) {
                        return this.slotFailure(req.accountId, CashShopResult.CASH_SHOP_RESULT_NOT_FOUND);
                    }
                    slots = (character.slotLimits[index] ?? 0) + req.amount;
                    transfer = { targetId: req.characterId, kind: CashTransferKind.InventorySlot, detail: { ...detail, index } };
                    break;
                }
                case CashSlotKind.CASH_SLOT_KIND_STORAGE: {
                    const storage = await this.storageRepo.get(req.worldId, req.accountId);
                    slots = (storage?.slots ?? DEFAULT_STORAGE_SLOTS) + req.amount;
                    transfer = { targetId: req.accountId, kind: CashTransferKind.StorageSlot, detail };
                    break;
                }
                case CashSlotKind.CASH_SLOT_KIND_CHARACTER: {
                    slots = (await this.accountWorldRepo.findCharacterSlotCount(req.worldId, req.accountId)) + req.amount;
                    transfer = { targetId: req.accountId, kind: CashTransferKind.CharacterSlot, detail };
                    break;
                }
                default:
                    return this.slotFailure(req.accountId, CashShopResult.CASH_SHOP_RESULT_NOT_FOUND);
            }
            if (slots > req.max) {
                return this.slotFailure(req.accountId, CashShopResult.CASH_SHOP_RESULT_SLOT_LIMIT);
            }

            const charged = await this.charge(req.worldId, req.accountId, delta, null, [transfer]);
            if (charged.result !== CashShopResult.CASH_SHOP_RESULT_OK) {
                return { result: charged.result, slots: 0, nxCash: charged.balance.nxCash, maplePoint: charged.balance.maplePoint };
            }
            await this.settleTransfers(req.worldId);
            return { result: charged.result, slots, nxCash: charged.balance.nxCash, maplePoint: charged.balance.maplePoint };
        } finally {
            if (req.kind === CashSlotKind.CASH_SLOT_KIND_INVENTORY) {
                await this.characterRepo.invalidateCache(req.worldId, req.characterId);
            }
        }
    }

    private async slotFailure(accountId: number, result: CashShopResult): Promise<ExpandCashSlotReply> {
        const balance = await this.cashShopRepo.findBalance(accountId);
        return { result, slots: 0, nxCash: balance.nxCash, maplePoint: balance.maplePoint };
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

    private async reservedLockerCount(worldId: number, accountId: number): Promise<number> {
        const items = await this.cashShopRepo.findItems(worldId, accountId);
        return items.length + (await this.cashShopRepo.pendingLockerItemCount(worldId, accountId));
    }

    // Both accounts' lockers in this world are held while the gift is checked and charged.
    async gift(req: GiftCashItemRequest): Promise<GiftCashItemReply> {
        const fail = async (result: CashShopResult) => {
            const balance = await this.cashShopRepo.findBalance(req.accountId);
            return { result, nxCash: balance.nxCash, maplePoint: balance.maplePoint };
        };
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

        await using _locks = await this.distributedLockService.acquireWorldDataLocks(req.worldId, [
            `cash_account:${req.accountId}`,
            `cash_account:${entry.account_id}`,
        ]);
        if ((await this.cashBalanceRequestRepo.find(req.worldId, req.accountId, req.requestId)) != null) {
            return fail(CashShopResult.CASH_SHOP_RESULT_OK);
        }
        if ((await this.reservedLockerCount(req.worldId, entry.account_id)) + req.items.length > LOCKER_CAPACITY) {
            return fail(CashShopResult.CASH_SHOP_RESULT_RECIPIENT_LOCKER_FULL);
        }

        const charged = await this.charge(req.worldId, req.accountId, this.debit(CashCurrency.CASH_CURRENCY_NX_CASH, req.price), req.requestId, [{
            targetId: entry.account_id,
            kind: CashTransferKind.Gift,
            detail: {
                accountId: entry.account_id,
                items: req.items.map((item) => CashItem.toJSON(item)),
                message: req.message,
                recipientId: entry.character_id,
                senderName: req.items[0]?.buyerName ?? "",
                price: req.price,
            },
        }]);
        if (charged.result === CashShopResult.CASH_SHOP_RESULT_OK) {
            await this.settleTransfers(req.worldId);
        }
        return { result: charged.result, nxCash: charged.balance.nxCash, maplePoint: charged.balance.maplePoint };
    }

    async buyRing(req: BuyCashRingRequest): Promise<BuyCashRingReply> {
        const item = req.item!;
        const partnerItem = req.partnerItem!;
        const fail = async (result: CashShopResult) => {
            const balance = await this.cashShopRepo.findBalance(req.accountId);
            return { result, nxCash: balance.nxCash, maplePoint: balance.maplePoint };
        };
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

        await using _locks = await this.distributedLockService.acquireWorldDataLocks(req.worldId, [
            `cash_account:${req.accountId}`,
            `cash_account:${entry.account_id}`,
        ]);
        if ((await this.cashBalanceRequestRepo.find(req.worldId, req.accountId, req.requestId)) != null) {
            return fail(CashShopResult.CASH_SHOP_RESULT_OK);
        }
        const balance = await this.cashShopRepo.findBalance(req.accountId);
        if (balance.nxCash < req.price) {
            return fail(CashShopResult.CASH_SHOP_RESULT_NOT_ENOUGH_CASH);
        }
        if ((await this.reservedLockerCount(req.worldId, req.accountId)) >= LOCKER_CAPACITY) {
            return fail(CashShopResult.CASH_SHOP_RESULT_LOCKER_FULL);
        }
        if ((await this.reservedLockerCount(req.worldId, entry.account_id)) + 1 > LOCKER_CAPACITY) {
            return fail(CashShopResult.CASH_SHOP_RESULT_RECIPIENT_LOCKER_FULL);
        }

        const serial = String(item.item?.uniqueId ?? 0);
        const partnerSerial = String(partnerItem.item?.uniqueId ?? 0);
        const itemId = item.item?.itemId ?? 0;
        const charged = await this.charge(req.worldId, req.accountId, this.debit(CashCurrency.CASH_CURRENCY_NX_CASH, req.price), req.requestId, [
            {
                targetId: entry.account_id,
                kind: CashTransferKind.Gift,
                detail: {
                    accountId: entry.account_id,
                    items: [CashItem.toJSON(partnerItem)],
                    message: req.message,
                    recipientId: entry.character_id,
                    senderName: req.characterName,
                    price: req.price,
                    buyerItem: { accountId: req.accountId, item: CashItem.toJSON(item) },
                    rings: [
                        { characterId: req.characterId, serial },
                        { characterId: entry.character_id, serial: partnerSerial },
                    ],
                },
            },
            {
                targetId: req.characterId,
                kind: CashTransferKind.Ring,
                detail: { serial, partner_serial: partnerSerial, character_id: req.characterId, partner_id: entry.character_id, partner_name: partner.name, item_id: itemId },
            },
            {
                targetId: entry.character_id,
                kind: CashTransferKind.Ring,
                detail: { serial: partnerSerial, partner_serial: serial, character_id: entry.character_id, partner_id: req.characterId, partner_name: req.characterName, item_id: itemId },
            },
        ]);
        if (charged.result === CashShopResult.CASH_SHOP_RESULT_OK) {
            // The second pass places the buyer's half, queued when the partner's half was placed.
            await this.settleTransfers(req.worldId);
            await this.settleTransfers(req.worldId);
        }
        return { result: charged.result, nxCash: charged.balance.nxCash, maplePoint: charged.balance.maplePoint };
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

    // The item leaves the locker with a credit request in the same transaction; settle adds the maple points.
    async payBack(req: PayBackCashItemRequest): Promise<PayBackCashItemReply> {
        const removed = await this.cashShopRepo.withAccountLock(req.worldId, req.accountId, async (txClient) => {
            const row = await this.cashShopRepo.deleteItem(req.worldId, req.accountId, String(req.serial), { txClient });
            if (row == null) {
                return false;
            }
            await this.cashShopRepo.queueTransfer(req.accountId, req.accountId, CashTransferKind.Credit, { delta: { nxCash: 0, maplePoint: req.maplePoint } }, txClient);
            return true;
        });
        if (removed) {
            await this.settleTransfers(req.worldId);
        }
        const balance = await this.cashShopRepo.findBalance(req.accountId);
        return {
            result: removed ? CashShopResult.CASH_SHOP_RESULT_OK : CashShopResult.CASH_SHOP_RESULT_NOT_FOUND,
            nxCash: balance.nxCash,
            maplePoint: balance.maplePoint,
        };
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
        if (coupon.used_at != null && (coupon.used_by !== accountId || coupon.request_id == null || (await this.couponGranted(worldId, accountId, coupon.kind, coupon.request_id)))) {
            return { result: CashShopResult.CASH_SHOP_RESULT_COUPON_USED, kind: coupon.kind, value: coupon.value };
        }
        return { result: CashShopResult.CASH_SHOP_RESULT_OK, kind: coupon.kind, value: coupon.value };
    }

    // A balance coupon is granted in the unified database and an item or meso coupon in the account's shard, each
    // recorded under the coupon's request id in the same transaction.
    private async couponGranted(worldId: number, accountId: number, kind: CashCouponKind, requestId: string): Promise<boolean> {
        if (kind === CashCouponKind.CASH_COUPON_KIND_NX_CASH || kind === CashCouponKind.CASH_COUPON_KIND_MAPLE_POINT) {
            return (await this.cashBalanceRequestRepo.find(worldId, accountId, requestId)) != null;
        }
        return (await this.cashRequestRepo.find(worldId, accountId, requestId)) != null;
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

        let meso = character.meso ?? 0;
        let result: CashShopResult;
        switch (coupon.kind) {
            case CashCouponKind.CASH_COUPON_KIND_NX_CASH:
            case CashCouponKind.CASH_COUPON_KIND_MAPLE_POINT: {
                const delta = coupon.kind === CashCouponKind.CASH_COUPON_KIND_NX_CASH ? { nxCash: coupon.value, maplePoint: 0 } : { nxCash: 0, maplePoint: coupon.value };
                result = await this.cashShopRepo.withBalanceLock(req.accountId, async (txClient) => {
                    if ((await this.cashBalanceRequestRepo.reserve(req.worldId, req.accountId, requestId, {}, txClient)) === false) {
                        return CashShopResult.CASH_SHOP_RESULT_COUPON_USED;
                    }
                    await this.cashShopRepo.addBalance(req.accountId, delta, { txClient });
                    return CashShopResult.CASH_SHOP_RESULT_OK;
                });
                break;
            }
            case CashCouponKind.CASH_COUPON_KIND_ITEM:
            case CashCouponKind.CASH_COUPON_KIND_MESO:
                result = await this.cashShopRepo.withAccountLock(req.worldId, req.accountId, async (txClient) => {
                    if (coupon.kind === CashCouponKind.CASH_COUPON_KIND_ITEM) {
                        const items = await this.cashShopRepo.findItems(req.worldId, req.accountId, { txClient });
                        if (items.length >= LOCKER_CAPACITY) {
                            return CashShopResult.CASH_SHOP_RESULT_LOCKER_FULL;
                        }
                    }
                    if ((await this.cashRequestRepo.reserve(req.worldId, req.accountId, requestId, {}, txClient)) === false) {
                        return CashShopResult.CASH_SHOP_RESULT_COUPON_USED;
                    }
                    if (coupon.kind === CashCouponKind.CASH_COUPON_KIND_ITEM) {
                        if (req.item?.item == null) {
                            throw new Error(`coupon ${req.code}: item is required`);
                        }
                        await this.cashShopRepo.insertItem(req.worldId, req.accountId, String(req.item.item.uniqueId ?? 0), CashItem.toJSON(req.item), { txClient });
                    } else {
                        meso = Math.min(meso + coupon.value, MAX_MESO);
                        await this.cashShopRepo.queueTransfer(req.accountId, req.characterId, CashTransferKind.Meso, { meso: coupon.value }, txClient);
                    }
                    return CashShopResult.CASH_SHOP_RESULT_OK;
                });
                break;
            default:
                throw new Error(`coupon ${req.code}: unknown kind ${coupon.kind}`);
        }
        if (result !== CashShopResult.CASH_SHOP_RESULT_OK) {
            return result === CashShopResult.CASH_SHOP_RESULT_COUPON_USED ? used : { result, nxCash: 0, maplePoint: 0, meso: 0 };
        }
        if (coupon.kind === CashCouponKind.CASH_COUPON_KIND_MESO) {
            await this.settleTransfers(req.worldId);
        }
        const balance = await this.cashShopRepo.findBalance(req.accountId);
        return { result, nxCash: balance.nxCash, maplePoint: balance.maplePoint, meso };
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
