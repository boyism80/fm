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
    SetCashWishlistReply,
    SetCashWishlistRequest,
    TakeOutCashItemReply,
    TakeOutCashItemRequest,
} from "../../protobuf/generated/fminternal/internal_service";
import type { GrpcCall, GrpcCallback, GrpcErrorHandler } from "./types";
import { Controller, Method } from "../grpc-method-decorator";
import { grpcMapper } from "../mappers";
import { INVENTORY_MODEL, INVENTORY_PROTO } from "../inventory-proto";
import type { Inventory } from "../../protobuf/generated/fminternal/internal_service";
import type { InventoryModel } from "../../repos/inventory-repository";

@Controller("cashShopController")
export class CashShopGrpcController {
    private readonly cashShopService: CashShopService;
    private readonly grpcError: GrpcErrorHandler;

    constructor(cashShopService: CashShopService, grpcError: GrpcErrorHandler) {
        this.cashShopService = cashShopService;
        this.grpcError = grpcError;
    }

    @Method("buyCashItem")
    async buyCashItem(call: GrpcCall<BuyCashItemRequest>, callback: GrpcCallback<BuyCashItemReply>) {
        try {
            if (call.request.items.length === 0 || call.request.items.some((item) => !item.item?.uniqueId)) {
                throw Object.assign(new Error("items with unique_id are required"), { code: "INVALID_PAYLOAD" });
            }
            callback(null, await this.cashShopService.buy(call.request));
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("takeOutCashItem")
    async takeOutCashItem(call: GrpcCall<TakeOutCashItemRequest>, callback: GrpcCallback<TakeOutCashItemReply>) {
        try {
            const req = call.request;
            if (!req.item?.uniqueId) {
                throw Object.assign(new Error("item with unique_id is required"), { code: "INVALID_PAYLOAD" });
            }
            const item = grpcMapper.map<Inventory, InventoryModel>(req.item, INVENTORY_PROTO, INVENTORY_MODEL);
            const result = await this.cashShopService.takeOut(req.worldId, req.accountId, req.characterId, item);
            callback(null, { result });
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("putInCashItem")
    async putInCashItem(call: GrpcCall<PutInCashItemRequest>, callback: GrpcCallback<PutInCashItemReply>) {
        try {
            const req = call.request;
            if (!req.item?.item?.uniqueId) {
                throw Object.assign(new Error("item with unique_id is required"), { code: "INVALID_PAYLOAD" });
            }
            const result = await this.cashShopService.putIn(req.worldId, req.accountId, req.characterId, req.item);
            callback(null, { result });
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("setCashWishlist")
    async setCashWishlist(call: GrpcCall<SetCashWishlistRequest>, callback: GrpcCallback<SetCashWishlistReply>) {
        try {
            await this.cashShopService.setWishlist(call.request.worldId, call.request.characterId, call.request.commoditySns);
            callback(null, { ok: true });
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("addCash")
    async addCash(call: GrpcCall<AddCashRequest>, callback: GrpcCallback<AddCashReply>) {
        try {
            const req = call.request;
            callback(null, await this.cashShopService.addCash(req.worldId, req.accountId, req.nxCash, req.maplePoint));
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("expandCashSlot")
    async expandCashSlot(call: GrpcCall<ExpandCashSlotRequest>, callback: GrpcCallback<ExpandCashSlotReply>) {
        try {
            callback(null, await this.cashShopService.expandSlot(call.request));
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("giftCashItem")
    async giftCashItem(call: GrpcCall<GiftCashItemRequest>, callback: GrpcCallback<GiftCashItemReply>) {
        try {
            if (call.request.items.length === 0 || call.request.items.some((item) => !item.item?.uniqueId)) {
                throw Object.assign(new Error("items with unique_id are required"), { code: "INVALID_PAYLOAD" });
            }
            callback(null, await this.cashShopService.gift(call.request));
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("buyCashRing")
    async buyCashRing(call: GrpcCall<BuyCashRingRequest>, callback: GrpcCallback<BuyCashRingReply>) {
        try {
            if (!call.request.item?.item?.uniqueId || !call.request.partnerItem?.item?.uniqueId) {
                throw Object.assign(new Error("item and partner_item with unique_id are required"), { code: "INVALID_PAYLOAD" });
            }
            callback(null, await this.cashShopService.buyRing(call.request));
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("payBackCashItem")
    async payBackCashItem(call: GrpcCall<PayBackCashItemRequest>, callback: GrpcCallback<PayBackCashItemReply>) {
        try {
            callback(null, await this.cashShopService.payBack(call.request));
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("buyCashQuestItem")
    async buyCashQuestItem(call: GrpcCall<BuyCashQuestItemRequest>, callback: GrpcCallback<BuyCashQuestItemReply>) {
        try {
            const req = call.request;
            if (!req.item?.uniqueId) {
                throw Object.assign(new Error("item with unique_id is required"), { code: "INVALID_PAYLOAD" });
            }
            const item = grpcMapper.map<Inventory, InventoryModel>(req.item, INVENTORY_PROTO, INVENTORY_MODEL);
            callback(null, await this.cashShopService.buyQuestItem(req.worldId, req.characterId, req.price, item));
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("findCashCoupon")
    async findCashCoupon(call: GrpcCall<FindCashCouponRequest>, callback: GrpcCallback<FindCashCouponReply>) {
        try {
            callback(null, await this.cashShopService.findCoupon(call.request.worldId, call.request.code.toUpperCase()));
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("redeemCashCoupon")
    async redeemCashCoupon(call: GrpcCall<RedeemCashCouponRequest>, callback: GrpcCallback<RedeemCashCouponReply>) {
        try {
            callback(null, await this.cashShopService.redeemCoupon({ ...call.request, code: call.request.code.toUpperCase() }));
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("createCashCoupons")
    async createCashCoupons(call: GrpcCall<CreateCashCouponsRequest>, callback: GrpcCallback<CreateCashCouponsReply>) {
        try {
            const req = call.request;
            if (req.count < 1 || req.count > 100) {
                throw Object.assign(new Error("count must be 1..100"), { code: "INVALID_PAYLOAD" });
            }
            callback(null, { codes: await this.cashShopService.createCoupons(req.worldId, req.kind, req.value, req.count) });
        } catch (err) {
            this.grpcError(err, callback);
        }
    }
}
