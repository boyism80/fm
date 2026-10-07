import {
    CashCurrency,
    CashItem,
    CashShopResult,
    type BuyCashItemReply,
    type BuyCashItemRequest,
} from "../protobuf/generated/fminternal/internal_service";
import type { AccountRepository } from "../repos/account-repository";
import type { CashShopRepository, CashItemRow } from "../repos/cash-shop-repository";
import type { StorageRepository } from "../repos/storage-repository";
import type { InventoryRepository, InventoryModel } from "../repos/inventory-repository";

const LOCKER_CAPACITY = 100;

export class CashShopService {
    private readonly cashShopRepo: CashShopRepository;
    private readonly accountRepo: AccountRepository;
    private readonly storageRepo: StorageRepository;
    private readonly inventoryRepo: InventoryRepository;

    constructor(
        cashShopRepository: CashShopRepository,
        accountRepository: AccountRepository,
        storageRepository: StorageRepository,
        inventoryRepository: InventoryRepository
    ) {
        this.cashShopRepo = cashShopRepository;
        this.accountRepo = accountRepository;
        this.storageRepo = storageRepository;
        this.inventoryRepo = inventoryRepository;
    }

    private toProto(row: CashItemRow): CashItem {
        return CashItem.fromJSON(row.item);
    }

    async loadAccount(worldId: number, accountId: number, characterId: number) {
        const [balance, items, wishlist, account, storage] = await Promise.all([
            this.cashShopRepo.findBalance(worldId, accountId),
            this.cashShopRepo.findItems(worldId, accountId),
            this.cashShopRepo.findWishlist(worldId, characterId),
            this.accountRepo.get(worldId, accountId),
            this.storageRepo.get(worldId, accountId),
        ]);
        return {
            nxCash: balance.nxCash,
            maplePoint: balance.maplePoint,
            locker: items.map((row) => this.toProto(row)),
            wishlist,
            characterSlotCount: account?.characterSlotCount ?? 0,
            storageSlotCount: storage?.slots ?? 4,
        };
    }

    async buy(req: BuyCashItemRequest): Promise<BuyCashItemReply> {
        const serial = String(req.item?.item?.uniqueId ?? 0);
        return this.cashShopRepo.withAccountLock(req.worldId, req.accountId, async (txClient) => {
            const balance = await this.cashShopRepo.findBalance(req.worldId, req.accountId, { txClient });
            const fail = (result: CashShopResult) => ({ result, nxCash: balance.nxCash, maplePoint: balance.maplePoint });
            const delta = { nxCash: 0, maplePoint: 0 };
            if (req.currency === CashCurrency.CASH_CURRENCY_MAPLE_POINT) {
                delta.maplePoint = -req.price;
            } else {
                delta.nxCash = -req.price;
            }
            if (balance.nxCash + delta.nxCash < 0 || balance.maplePoint + delta.maplePoint < 0) {
                return fail(CashShopResult.CASH_SHOP_RESULT_NOT_ENOUGH_CASH);
            }

            const items = await this.cashShopRepo.findItems(req.worldId, req.accountId, { txClient });
            if (items.length >= LOCKER_CAPACITY) {
                return fail(CashShopResult.CASH_SHOP_RESULT_LOCKER_FULL);
            }

            await this.cashShopRepo.addBalance(req.worldId, req.accountId, delta, { txClient });
            await this.cashShopRepo.insertItem(req.worldId, req.accountId, serial, CashItem.toJSON(req.item!), { txClient });
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

    async setWishlist(worldId: number, characterId: number, commoditySns: number[]): Promise<void> {
        await this.cashShopRepo.setWishlist(worldId, characterId, commoditySns.slice(0, 10));
    }
}
