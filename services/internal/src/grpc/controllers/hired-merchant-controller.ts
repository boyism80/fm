import type { HiredMerchantService } from "../../services/hired-merchant-service";
import type {
    ClaimStoreBankReply,
    ClaimStoreBankRequest,
    CloseChannelMerchantsReply,
    CloseChannelMerchantsRequest,
    FindHiredMerchantReply,
    FindHiredMerchantRequest,
    OpenHiredMerchantReply,
    OpenHiredMerchantRequest,
    SaveHiredMerchantReply,
    SaveHiredMerchantRequest,
} from "../../protobuf/generated/fminternal/internal_service";
import type { GrpcCall, GrpcCallback, GrpcErrorHandler } from "./types";
import { Controller, Method } from "../grpc-method-decorator";
import { makeSaveCharacterEntry } from "./character-controller";

@Controller("hiredMerchantController")
export class HiredMerchantGrpcController {
    private readonly hiredMerchantService: HiredMerchantService;
    private readonly grpcError: GrpcErrorHandler;

    constructor(hiredMerchantService: HiredMerchantService, grpcError: GrpcErrorHandler) {
        this.hiredMerchantService = hiredMerchantService;
        this.grpcError = grpcError;
    }

    @Method("findHiredMerchant")
    async findHiredMerchant(call: GrpcCall<FindHiredMerchantRequest>, callback: GrpcCallback<FindHiredMerchantReply>) {
        try {
            const req = call.request;
            callback(null, { merchant: await this.hiredMerchantService.findHiredMerchant(req.worldId, req.accountId) });
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("openHiredMerchant")
    async openHiredMerchant(call: GrpcCall<OpenHiredMerchantRequest>, callback: GrpcCallback<OpenHiredMerchantReply>) {
        try {
            const req = call.request;
            if (!req.merchant) {
                throw Object.assign(new Error("merchant is required"), { code: "INVALID_PAYLOAD" });
            }
            callback(null, await this.hiredMerchantService.openHiredMerchant(req.merchant));
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("saveHiredMerchant")
    async saveHiredMerchant(call: GrpcCall<SaveHiredMerchantRequest>, callback: GrpcCallback<SaveHiredMerchantReply>) {
        try {
            const req = call.request;
            if (!req.merchant) {
                throw Object.assign(new Error("merchant is required"), { code: "INVALID_PAYLOAD" });
            }
            await this.hiredMerchantService.saveHiredMerchant(
                req.merchant,
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
                merchant: await this.hiredMerchantService.claimStoreBank(req.worldId, req.accountId, req.characterId),
            });
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("closeChannelMerchants")
    async closeChannelMerchants(
        call: GrpcCall<CloseChannelMerchantsRequest>,
        callback: GrpcCallback<CloseChannelMerchantsReply>
    ) {
        try {
            const req = call.request;
            callback(null, { count: await this.hiredMerchantService.closeChannelMerchants(req.worldId, req.channelId) });
        } catch (err) {
            this.grpcError(err, callback);
        }
    }
}
