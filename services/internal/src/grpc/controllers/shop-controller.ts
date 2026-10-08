import type { ShopService } from "../../services/shop-service";
import type {
    ClaimStoreBankReply,
    ClaimStoreBankRequest,
    CloseChannelShopsReply,
    CloseChannelShopsRequest,
    FindEntrustedShopReply,
    FindEntrustedShopRequest,
    FindPopularShopSearchesReply,
    FindPopularShopSearchesRequest,
    OpenShopReply,
    OpenShopRequest,
    SaveShopReply,
    SaveShopRequest,
    SearchShopsReply,
    SearchShopsRequest,
} from "../../protobuf/generated/fminternal/internal_service";
import type { GrpcCall, GrpcCallback, GrpcErrorHandler } from "./types";
import { Controller, Method } from "../grpc-method-decorator";
import { makeSaveCharacterEntry } from "./character-controller";

@Controller("shopController")
export class ShopGrpcController {
    private readonly shopService: ShopService;
    private readonly grpcError: GrpcErrorHandler;

    constructor(shopService: ShopService, grpcError: GrpcErrorHandler) {
        this.shopService = shopService;
        this.grpcError = grpcError;
    }

    @Method("findEntrustedShop")
    async findEntrustedShop(call: GrpcCall<FindEntrustedShopRequest>, callback: GrpcCallback<FindEntrustedShopReply>) {
        try {
            const req = call.request;
            callback(null, await this.shopService.findEntrustedShop(req.worldId, req.accountId));
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("openShop")
    async openShop(call: GrpcCall<OpenShopRequest>, callback: GrpcCallback<OpenShopReply>) {
        try {
            const req = call.request;
            if (!req.shop) {
                throw Object.assign(new Error("shop is required"), { code: "INVALID_PAYLOAD" });
            }
            callback(null, await this.shopService.openShop(req.shop));
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("saveShop")
    async saveShop(call: GrpcCall<SaveShopRequest>, callback: GrpcCallback<SaveShopReply>) {
        try {
            const req = call.request;
            if (!req.shop) {
                throw Object.assign(new Error("shop is required"), { code: "INVALID_PAYLOAD" });
            }
            await this.shopService.saveShop(
                req.shop,
                (req.characters ?? []).map((entry) => makeSaveCharacterEntry(entry)),
                req.close === true
            );
            callback(null, {});
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("claimStoreBank")
    async claimStoreBank(call: GrpcCall<ClaimStoreBankRequest>, callback: GrpcCallback<ClaimStoreBankReply>) {
        try {
            const req = call.request;
            callback(null, { shops: await this.shopService.claimStoreBank(req.worldId, req.accountId, req.characterId) });
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("closeChannelShops")
    async closeChannelShops(call: GrpcCall<CloseChannelShopsRequest>, callback: GrpcCallback<CloseChannelShopsReply>) {
        try {
            const req = call.request;
            callback(null, { count: await this.shopService.closeChannelShops(req.worldId, req.channelId) });
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("searchShops")
    async searchShops(call: GrpcCall<SearchShopsRequest>, callback: GrpcCallback<SearchShopsReply>) {
        try {
            const req = call.request;
            callback(null, { entries: await this.shopService.searchShops(req.worldId, req.itemId, req.highFirst === true) });
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("findPopularShopSearches")
    async findPopularShopSearches(
        call: GrpcCall<FindPopularShopSearchesRequest>,
        callback: GrpcCallback<FindPopularShopSearchesReply>
    ) {
        try {
            callback(null, { itemIds: this.shopService.findPopularShopSearches(call.request.worldId) });
        } catch (err) {
            this.grpcError(err, callback);
        }
    }
}
