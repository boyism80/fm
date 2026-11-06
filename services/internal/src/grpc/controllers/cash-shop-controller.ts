import type { CashShopService } from "../../services/cash-shop-service";
import type {
    AddCashReply,
    AddCashRequest,
    BuyCashItemReply,
    BuyCashItemRequest,
    BuyCashQuestItemReply,
    BuyCashQuestItemRequest,
    CreateCashCouponsReply,
    CreateCashCouponsRequest,
    ExpandCashSlotReply,
    ExpandCashSlotRequest,
    FindCashCouponReply,
    FindCashCouponRequest,
    BuyCashRingReply,
    BuyCashRingRequest,
    GiftCashItemReply,
    GiftCashItemRequest,
    PayBackCashItemReply,
    PayBackCashItemRequest,
    PutInCashItemReply,
    PutInCashItemRequest,
    RedeemCashCouponReply,
    RedeemCashCouponRequest,
    ReloadCashShopReply,
    ReloadCashShopRequest,
    SetCashWishlistReply,
    SetCashWishlistRequest,
    TakeOutCashItemReply,
    TakeOutCashItemRequest,
} from "@fm/shared/protobuf/generated/fminternal/internal_service";
import type { GrpcCall, GrpcCallback } from "./types";
import { Controller, Method } from "../grpc-method-decorator";
import { grpcMapper } from "../mappers";
import { INVENTORY_MODEL, INVENTORY_PROTO } from "../inventory-proto";
import type { Inventory } from "@fm/shared/protobuf/generated/fminternal/internal_service";
import type { InventoryModel } from "@fm/shared/repos/inventory-repository";

@Controller("cashShopController")
export class CashShopGrpcController {
    private readonly cashShopService: CashShopService;

    constructor(cashShopService: CashShopService) {
        this.cashShopService = cashShopService;
    }

    @Method("buyCashItem")
    async buyCashItem(call: GrpcCall<BuyCashItemRequest>, callback: GrpcCallback<BuyCashItemReply>) {
        if (call.request.items.length === 0 || call.request.items.some((item) => !item.item?.uniqueId)) {
            throw Object.assign(new Error("items with unique_id are required"), { code: "INVALID_PAYLOAD" });
        }
        callback(null, await this.cashShopService.buy(call.request));
    }

    @Method("takeOutCashItem")
    async takeOutCashItem(call: GrpcCall<TakeOutCashItemRequest>, callback: GrpcCallback<TakeOutCashItemReply>) {
        const req = call.request;
        if (!req.item?.uniqueId) {
            throw Object.assign(new Error("item with unique_id is required"), { code: "INVALID_PAYLOAD" });
        }
        const item = grpcMapper.map<Inventory, InventoryModel>(req.item, INVENTORY_PROTO, INVENTORY_MODEL);
        const { result, slot } = await this.cashShopService.takeOut(req.worldId, req.accountId, req.characterId, item);
        callback(null, { result, slot });
    }

    @Method("putInCashItem")
    async putInCashItem(call: GrpcCall<PutInCashItemRequest>, callback: GrpcCallback<PutInCashItemReply>) {
        const req = call.request;
        if (!req.item?.item?.uniqueId) {
            throw Object.assign(new Error("item with unique_id is required"), { code: "INVALID_PAYLOAD" });
        }
        const result = await this.cashShopService.putIn(req.worldId, req.accountId, req.characterId, req.item);
        callback(null, { result });
    }

    @Method("setCashWishlist")
    async setCashWishlist(call: GrpcCall<SetCashWishlistRequest>, callback: GrpcCallback<SetCashWishlistReply>) {
        await this.cashShopService.setWishlist(
            call.request.worldId,
            call.request.characterId,
            call.request.commoditySns,
        );
        callback(null, { ok: true });
    }

    @Method("addCash")
    async addCash(call: GrpcCall<AddCashRequest>, callback: GrpcCallback<AddCashReply>) {
        const req = call.request;
        callback(null, await this.cashShopService.addCash(req.accountId, req.nxCash, req.maplePoint));
    }

    @Method("expandCashSlot")
    async expandCashSlot(call: GrpcCall<ExpandCashSlotRequest>, callback: GrpcCallback<ExpandCashSlotReply>) {
        callback(null, await this.cashShopService.expandSlot(call.request));
    }

    @Method("giftCashItem")
    async giftCashItem(call: GrpcCall<GiftCashItemRequest>, callback: GrpcCallback<GiftCashItemReply>) {
        if (call.request.items.length === 0 || call.request.items.some((item) => !item.item?.uniqueId)) {
            throw Object.assign(new Error("items with unique_id are required"), { code: "INVALID_PAYLOAD" });
        }
        callback(null, await this.cashShopService.gift(call.request));
    }

    @Method("buyCashRing")
    async buyCashRing(call: GrpcCall<BuyCashRingRequest>, callback: GrpcCallback<BuyCashRingReply>) {
        if (!call.request.item?.item?.uniqueId || !call.request.partnerItem?.item?.uniqueId) {
            throw Object.assign(new Error("item and partner_item with unique_id are required"), {
                code: "INVALID_PAYLOAD",
            });
        }
        callback(null, await this.cashShopService.buyRing(call.request));
    }

    @Method("payBackCashItem")
    async payBackCashItem(call: GrpcCall<PayBackCashItemRequest>, callback: GrpcCallback<PayBackCashItemReply>) {
        callback(null, await this.cashShopService.payBack(call.request));
    }

    @Method("buyCashQuestItem")
    async buyCashQuestItem(call: GrpcCall<BuyCashQuestItemRequest>, callback: GrpcCallback<BuyCashQuestItemReply>) {
        const req = call.request;
        if (!req.item?.uniqueId) {
            throw Object.assign(new Error("item with unique_id is required"), { code: "INVALID_PAYLOAD" });
        }
        const item = grpcMapper.map<Inventory, InventoryModel>(req.item, INVENTORY_PROTO, INVENTORY_MODEL);
        callback(null, await this.cashShopService.buyQuestItem(req.worldId, req.characterId, req.price, item));
    }

    @Method("reloadCashShop")
    async reloadCashShop(call: GrpcCall<ReloadCashShopRequest>, callback: GrpcCallback<ReloadCashShopReply>) {
        const req = call.request;
        const state = await this.cashShopService.reload(req.worldId, req.accountId, req.characterId);
        callback(null, {
            ...state,
            inventory: state.inventory.map((item) =>
                grpcMapper.map<InventoryModel, Inventory>(item, INVENTORY_MODEL, INVENTORY_PROTO),
            ),
        });
    }

    @Method("findCashCoupon")
    async findCashCoupon(call: GrpcCall<FindCashCouponRequest>, callback: GrpcCallback<FindCashCouponReply>) {
        callback(
            null,
            await this.cashShopService.findCoupon(
                call.request.worldId,
                call.request.code.toUpperCase(),
                call.request.accountId,
            ),
        );
    }

    @Method("redeemCashCoupon")
    async redeemCashCoupon(call: GrpcCall<RedeemCashCouponRequest>, callback: GrpcCallback<RedeemCashCouponReply>) {
        callback(
            null,
            await this.cashShopService.redeemCoupon({ ...call.request, code: call.request.code.toUpperCase() }),
        );
    }

    @Method("createCashCoupons")
    async createCashCoupons(call: GrpcCall<CreateCashCouponsRequest>, callback: GrpcCallback<CreateCashCouponsReply>) {
        const req = call.request;
        if (req.count < 1 || req.count > 100) {
            throw Object.assign(new Error("count must be 1..100"), { code: "INVALID_PAYLOAD" });
        }
        callback(null, {
            codes: await this.cashShopService.createCoupons(req.worldId, req.kind, req.value, req.count),
        });
    }
}
