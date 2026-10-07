import type { CashShopService } from "../../services/cash-shop-service";
import type {
    AddCashReply,
    AddCashRequest,
    BuyCashItemReply,
    BuyCashItemRequest,
    PutInCashItemReply,
    PutInCashItemRequest,
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
            if (!call.request.item?.item) {
                throw Object.assign(new Error("item is required"), { code: "INVALID_PAYLOAD" });
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
}
