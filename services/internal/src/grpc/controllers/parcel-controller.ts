import type { ParcelService } from "../../services/parcel-service";
import type {
    CheckParcelArrivalsReply,
    CheckParcelArrivalsRequest,
    ClaimParcelReply,
    ClaimParcelRequest,
    DeleteParcelReply,
    DeleteParcelRequest,
    LoadParcelsReply,
    LoadParcelsRequest,
    SendParcelReply,
    SendParcelRequest,
} from "../../protobuf/generated/fminternal/internal_service";
import type { GrpcCall, GrpcCallback, GrpcErrorHandler } from "./types";
import { Controller, Method } from "../grpc-method-decorator";
import { makeSaveCharacterEntry } from "./character-controller";

@Controller("parcelController")
export class ParcelGrpcController {
    private readonly parcelService: ParcelService;
    private readonly grpcError: GrpcErrorHandler;

    constructor(parcelService: ParcelService, grpcError: GrpcErrorHandler) {
        this.parcelService = parcelService;
        this.grpcError = grpcError;
    }

    @Method("loadParcels")
    async loadParcels(call: GrpcCall<LoadParcelsRequest>, callback: GrpcCallback<LoadParcelsReply>) {
        try {
            callback(null, await this.parcelService.loadParcels(call.request.worldId, call.request.characterId));
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("sendParcel")
    async sendParcel(call: GrpcCall<SendParcelRequest>, callback: GrpcCallback<SendParcelReply>) {
        try {
            const req = call.request;
            if (!req.parcel) {
                throw Object.assign(new Error("parcel is required"), { code: "INVALID_PAYLOAD" });
            }
            const result = await this.parcelService.sendParcel({
                worldId: req.worldId,
                recipientName: req.recipientName ?? "",
                parcel: req.parcel,
                senderAccountId: req.senderAccountId,
                sender: req.sender ? makeSaveCharacterEntry(req.sender) : undefined,
                oneOfAKind: req.oneOfAKind === true,
            });
            callback(null, { result });
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("claimParcel")
    async claimParcel(call: GrpcCall<ClaimParcelRequest>, callback: GrpcCallback<ClaimParcelReply>) {
        try {
            const req = call.request;
            callback(null, await this.parcelService.claimParcel(req.worldId, req.characterId, req.parcelId));
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("deleteParcel")
    async deleteParcel(call: GrpcCall<DeleteParcelRequest>, callback: GrpcCallback<DeleteParcelReply>) {
        try {
            const req = call.request;
            const result = await this.parcelService.deleteParcel(req.worldId, req.characterId, req.parcelId);
            callback(null, { result });
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("checkParcelArrivals")
    async checkParcelArrivals(
        call: GrpcCall<CheckParcelArrivalsRequest>,
        callback: GrpcCallback<CheckParcelArrivalsReply>
    ) {
        try {
            const req = call.request;
            callback(null, await this.parcelService.checkParcelArrivals(req.worldId, req.characterId));
        } catch (err) {
            this.grpcError(err, callback);
        }
    }
}
