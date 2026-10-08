import type { EntrustedShopService } from "../../services/entrusted-shop-service";
import type {
    ClaimStoreBankReply,
    ClaimStoreBankRequest,
    CloseChannelEntrustedShopsReply,
    CloseChannelEntrustedShopsRequest,
    FindEntrustedShopReply,
    FindEntrustedShopRequest,
    OpenEntrustedShopReply,
    OpenEntrustedShopRequest,
    SaveEntrustedShopReply,
    SaveEntrustedShopRequest,
} from "../../protobuf/generated/fminternal/internal_service";
import type { GrpcCall, GrpcCallback, GrpcErrorHandler } from "./types";
import { Controller, Method } from "../grpc-method-decorator";
import { makeSaveCharacterEntry } from "./character-controller";

@Controller("entrustedShopController")
export class EntrustedShopGrpcController {
    private readonly entrustedShopService: EntrustedShopService;
    private readonly grpcError: GrpcErrorHandler;

    constructor(entrustedShopService: EntrustedShopService, grpcError: GrpcErrorHandler) {
        this.entrustedShopService = entrustedShopService;
        this.grpcError = grpcError;
    }

    @Method("findEntrustedShop")
    async findEntrustedShop(call: GrpcCall<FindEntrustedShopRequest>, callback: GrpcCallback<FindEntrustedShopReply>) {
        try {
            const req = call.request;
            callback(null, { shop: await this.entrustedShopService.findEntrustedShop(req.worldId, req.accountId) });
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("openEntrustedShop")
    async openEntrustedShop(call: GrpcCall<OpenEntrustedShopRequest>, callback: GrpcCallback<OpenEntrustedShopReply>) {
        try {
            const req = call.request;
            if (!req.shop) {
                throw Object.assign(new Error("shop is required"), { code: "INVALID_PAYLOAD" });
            }
            callback(null, await this.entrustedShopService.openEntrustedShop(req.shop));
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("saveEntrustedShop")
    async saveEntrustedShop(call: GrpcCall<SaveEntrustedShopRequest>, callback: GrpcCallback<SaveEntrustedShopReply>) {
        try {
            const req = call.request;
            if (!req.shop) {
                throw Object.assign(new Error("shop is required"), { code: "INVALID_PAYLOAD" });
            }
            await this.entrustedShopService.saveEntrustedShop(
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
            callback(null, {
                shop: await this.entrustedShopService.claimStoreBank(req.worldId, req.accountId, req.characterId),
            });
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("closeChannelEntrustedShops")
    async closeChannelEntrustedShops(
        call: GrpcCall<CloseChannelEntrustedShopsRequest>,
        callback: GrpcCallback<CloseChannelEntrustedShopsReply>
    ) {
        try {
            const req = call.request;
            callback(null, { count: await this.entrustedShopService.closeChannelEntrustedShops(req.worldId, req.channelId) });
        } catch (err) {
            this.grpcError(err, callback);
        }
    }
}
