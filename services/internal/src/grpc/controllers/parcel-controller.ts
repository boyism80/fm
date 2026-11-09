import type { ParcelService } from "../../services/parcel-service";
import type {
    AbortParcelSendReply,
    AbortParcelSendRequest,
    CheckParcelArrivalsReply,
    CheckParcelArrivalsRequest,
    ClaimParcelReply,
    ClaimParcelRequest,
    DeleteParcelReply,
    DeleteParcelRequest,
    FindParcelRecipientReply,
    FindParcelRecipientRequest,
    LoadParcelsReply,
    LoadParcelsRequest,
    ReturnParcelReply,
    ReturnParcelRequest,
    SendParcelReply,
    SendParcelRequest,
} from "@fm/shared/protobuf/generated/fminternal/internal_service";
import type { GrpcCall, GrpcCallback } from "./types";
import { Controller, Method } from "../grpc-method-decorator";

@Controller("parcelController")
export class ParcelGrpcController {
    private readonly parcelService: ParcelService;

    constructor(parcelService: ParcelService) {
        this.parcelService = parcelService;
    }

    @Method("loadParcels")
    async loadParcels(call: GrpcCall<LoadParcelsRequest>, callback: GrpcCallback<LoadParcelsReply>) {
        callback(null, await this.parcelService.loadParcels(call.request.worldId, call.request.characterId));
    }

    @Method("sendParcel")
    async sendParcel(call: GrpcCall<SendParcelRequest>, callback: GrpcCallback<SendParcelReply>) {
        const req = call.request;
        if (!req.parcel) {
            throw Object.assign(new Error("parcel is required"), { code: "INVALID_PAYLOAD" });
        }
        const result = await this.parcelService.sendParcel({
            worldId: req.worldId,
            recipientName: req.recipientName ?? "",
            parcel: req.parcel,
            senderAccountId: req.senderAccountId,
            oneOfAKind: req.oneOfAKind === true,
            requestId: req.requestId,
            system: req.system === true,
        });
        callback(null, { result });
    }

    @Method("returnParcel")
    async returnParcel(call: GrpcCall<ReturnParcelRequest>, callback: GrpcCallback<ReturnParcelReply>) {
        const req = call.request;
        const result = await this.parcelService.returnParcels(req.worldId, req.characterId, req.requestId, req.parcels);
        callback(null, { result });
    }

    @Method("abortParcelSend")
    async abortParcelSend(call: GrpcCall<AbortParcelSendRequest>, callback: GrpcCallback<AbortParcelSendReply>) {
        const req = call.request;
        callback(null, {
            sent: await this.parcelService.abortParcelSend(req.worldId, req.receiverId, req.requestId),
        });
    }

    @Method("findParcelRecipient")
    async findParcelRecipient(
        call: GrpcCall<FindParcelRecipientRequest>,
        callback: GrpcCallback<FindParcelRecipientReply>,
    ) {
        const req = call.request;
        callback(
            null,
            await this.parcelService.findParcelRecipient(req.worldId, req.recipientName, req.senderAccountId, false),
        );
    }

    @Method("claimParcel")
    async claimParcel(call: GrpcCall<ClaimParcelRequest>, callback: GrpcCallback<ClaimParcelReply>) {
        const req = call.request;
        callback(null, await this.parcelService.claimParcel(req.worldId, req.characterId, req.parcelId, req.requestId));
    }

    @Method("deleteParcel")
    async deleteParcel(call: GrpcCall<DeleteParcelRequest>, callback: GrpcCallback<DeleteParcelReply>) {
        const req = call.request;
        const result = await this.parcelService.deleteParcel(req.worldId, req.characterId, req.parcelId);
        callback(null, { result });
    }

    @Method("checkParcelArrivals")
    async checkParcelArrivals(
        call: GrpcCall<CheckParcelArrivalsRequest>,
        callback: GrpcCallback<CheckParcelArrivalsReply>,
    ) {
        const req = call.request;
        callback(null, await this.parcelService.checkParcelArrivals(req.worldId, req.characterId));
    }
}
