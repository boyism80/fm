import type {
    ShopService,
} from "../../services/shop-service";
import type {
    AbortOpenShopReply,
    AbortOpenShopRequest,
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
} from "@fm/shared/protobuf/generated/fminternal/internal_service";
import type {
    GrpcCall,
    GrpcCallback,
} from "./types";
import { Controller, Method } from "../grpc-method-decorator";
import { makeSaveCharacterEntry } from "./character-controller";

@Controller("shopController")
export class ShopGrpcController {
    private readonly shopService: ShopService;

    constructor(shopService: ShopService) {
        this.shopService = shopService;
    }

    @Method("findEntrustedShop")
    async findEntrustedShop(call: GrpcCall<FindEntrustedShopRequest>, callback: GrpcCallback<FindEntrustedShopReply>) {
        const req = call.request;
        callback(null, await this.shopService.findEntrustedShop(req.worldId, req.accountId));
    }

    @Method("openShop")
    async openShop(call: GrpcCall<OpenShopRequest>, callback: GrpcCallback<OpenShopReply>) {
        const req = call.request;
        if (!req.shop) {
            throw Object.assign(new Error("shop is required"), { code: "INVALID_PAYLOAD" });
        }
        callback(null, await this.shopService.openShop(req.shop, req.requestId));
    }

    @Method("abortOpenShop")
    async abortOpenShop(call: GrpcCall<AbortOpenShopRequest>, callback: GrpcCallback<AbortOpenShopReply>) {
        const req = call.request;
        callback(null, {
            aborted: await this.shopService.abortOpenShop(req.worldId, req.accountId, req.requestId),
        });
    }

    @Method("saveShop")
    async saveShop(call: GrpcCall<SaveShopRequest>, callback: GrpcCallback<SaveShopReply>) {
        const req = call.request;
        if (!req.shop) {
            throw Object.assign(new Error("shop is required"), { code: "INVALID_PAYLOAD" });
        }
        await this.shopService.saveShop(
            req.shop,
            (req.characters ?? []).map((entry) => makeSaveCharacterEntry(entry)),
            req.close === true,
            req.storeBank ?? [],
        );
        callback(null, {});
    }

    @Method("closeChannelShops")
    async closeChannelShops(call: GrpcCall<CloseChannelShopsRequest>, callback: GrpcCallback<CloseChannelShopsReply>) {
        const req = call.request;
        callback(null, { count: await this.shopService.closeChannelShops(req.worldId, req.channelId) });
    }

    @Method("searchShops")
    async searchShops(call: GrpcCall<SearchShopsRequest>, callback: GrpcCallback<SearchShopsReply>) {
        const req = call.request;
        callback(null, {
            entries: await this.shopService.searchShops(req.worldId, req.itemId, req.highFirst === true),
        });
    }

    @Method("findPopularShopSearches")
    async findPopularShopSearches(
        call: GrpcCall<FindPopularShopSearchesRequest>,
        callback: GrpcCallback<FindPopularShopSearchesReply>,
    ) {
        callback(null, { itemIds: this.shopService.findPopularShopSearches(call.request.worldId) });
    }
}
