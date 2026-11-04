import type {
    MarriageService,
} from "../../services/marriage-service";
import type {
    AbortMarriageChangeReply,
    AbortMarriageChangeRequest,
    AbortWeddingGiftReply,
    AbortWeddingGiftRequest,
    BreakEngagementRequest,
    CancelDivorceRequest,
    ClaimWeddingGiftReply,
    ClaimWeddingGiftRequest,
    CreateMarriageRequest,
    FinishWeddingRequest,
    GetMarriageRequest,
    GiveWeddingGiftReply,
    GiveWeddingGiftRequest,
    InviteWeddingGuestReply,
    InviteWeddingGuestRequest,
    LoadMarriageRequest,
    LoadWeddingGiftsReply,
    LoadWeddingGiftsRequest,
    MarriageReply,
    NotifySpouseMapReply,
    NotifySpouseMapRequest,
    RequestDivorceRequest,
    ReserveWeddingRequest,
    SetWeddingWishlistRequest,
} from "@fm/shared/protobuf/generated/fminternal/internal_service";
import type {
    GrpcCall,
    GrpcCallback,
} from "./types";
import { Controller, Method } from "../grpc-method-decorator";

@Controller("marriageController")
export class MarriageGrpcController {
    private readonly marriageService: MarriageService;

    constructor(marriageService: MarriageService) {
        this.marriageService = marriageService;
    }

    @Method("createMarriage")
    async createMarriage(call: GrpcCall<CreateMarriageRequest>, callback: GrpcCallback<MarriageReply>) {
        const req = call.request;
        callback(
            null,
            await this.marriageService.createMarriage(
                req.worldId,
                req.groomId,
                req.groomName,
                req.brideId,
                req.brideName,
                req.ringItemId,
                req.requestId,
            ),
        );
    }

    @Method("loadMarriage")
    async loadMarriage(call: GrpcCall<LoadMarriageRequest>, callback: GrpcCallback<MarriageReply>) {
        callback(null, await this.marriageService.loadMarriage(call.request.worldId, call.request.characterId));
    }

    @Method("getMarriage")
    async getMarriage(call: GrpcCall<GetMarriageRequest>, callback: GrpcCallback<MarriageReply>) {
        callback(null, await this.marriageService.getMarriage(call.request.worldId, call.request.marriageId));
    }

    @Method("breakEngagement")
    async breakEngagement(call: GrpcCall<BreakEngagementRequest>, callback: GrpcCallback<MarriageReply>) {
        const req = call.request;
        callback(null, await this.marriageService.breakEngagement(req.worldId, req.marriageId, req.characterId));
    }

    @Method("reserveWedding")
    async reserveWedding(call: GrpcCall<ReserveWeddingRequest>, callback: GrpcCallback<MarriageReply>) {
        const req = call.request;
        callback(
            null,
            await this.marriageService.reserveWedding(
                req.worldId,
                req.marriageId,
                req.characterId,
                req.ticketItemId,
                req.requestId,
            ),
        );
    }

    @Method("abortMarriageChange")
    async abortMarriageChange(
        call: GrpcCall<AbortMarriageChangeRequest>,
        callback: GrpcCallback<AbortMarriageChangeReply>,
    ) {
        const req = call.request;
        callback(null, { changed: await this.marriageService.abortMarriageChange(req.worldId, req.requestId) });
    }

    @Method("setWeddingWishlist")
    async setWeddingWishlist(call: GrpcCall<SetWeddingWishlistRequest>, callback: GrpcCallback<MarriageReply>) {
        const req = call.request;
        callback(
            null,
            await this.marriageService.setWeddingWishlist(
                req.worldId,
                req.marriageId,
                req.characterId,
                req.wishes ?? [],
                req.requestId,
            ),
        );
    }

    @Method("inviteWeddingGuest")
    async inviteWeddingGuest(
        call: GrpcCall<InviteWeddingGuestRequest>,
        callback: GrpcCallback<InviteWeddingGuestReply>,
    ) {
        const req = call.request;
        if (!req.parcel) {
            throw Object.assign(new Error("parcel is required"), { code: "INVALID_PAYLOAD" });
        }
        callback(
            null,
            await this.marriageService.inviteWeddingGuest(
                req.worldId,
                req.marriageId,
                req.guestName,
                req.parcel,
                req.requestId,
            ),
        );
    }

    @Method("finishWedding")
    async finishWedding(call: GrpcCall<FinishWeddingRequest>, callback: GrpcCallback<MarriageReply>) {
        callback(null, await this.marriageService.finishWedding(call.request.worldId, call.request.marriageId));
    }

    @Method("requestDivorce")
    async requestDivorce(call: GrpcCall<RequestDivorceRequest>, callback: GrpcCallback<MarriageReply>) {
        const req = call.request;
        callback(
            null,
            await this.marriageService.requestDivorce(
                req.worldId,
                req.marriageId,
                req.characterId,
                Number(req.nowUnixMs) || Date.now(),
            ),
        );
    }

    @Method("cancelDivorce")
    async cancelDivorce(call: GrpcCall<CancelDivorceRequest>, callback: GrpcCallback<MarriageReply>) {
        const req = call.request;
        callback(null, await this.marriageService.cancelDivorce(req.worldId, req.marriageId, req.characterId));
    }

    @Method("giveWeddingGift")
    async giveWeddingGift(call: GrpcCall<GiveWeddingGiftRequest>, callback: GrpcCallback<GiveWeddingGiftReply>) {
        const req = call.request;
        if (!req.gift) {
            throw Object.assign(new Error("gift is required"), { code: "INVALID_PAYLOAD" });
        }
        const result = await this.marriageService.giveWeddingGift({
            worldId: req.worldId,
            receiverId: req.receiverId,
            gift: req.gift,
            requestId: req.requestId,
        });
        callback(null, { result });
    }

    @Method("abortWeddingGift")
    async abortWeddingGift(call: GrpcCall<AbortWeddingGiftRequest>, callback: GrpcCallback<AbortWeddingGiftReply>) {
        const req = call.request;
        callback(null, {
            given: await this.marriageService.abortWeddingGift(req.worldId, req.receiverId, req.requestId),
        });
    }

    @Method("loadWeddingGifts")
    async loadWeddingGifts(call: GrpcCall<LoadWeddingGiftsRequest>, callback: GrpcCallback<LoadWeddingGiftsReply>) {
        callback(null, await this.marriageService.loadWeddingGifts(call.request.worldId, call.request.characterId));
    }

    @Method("claimWeddingGift")
    async claimWeddingGift(call: GrpcCall<ClaimWeddingGiftRequest>, callback: GrpcCallback<ClaimWeddingGiftReply>) {
        const req = call.request;
        callback(
            null,
            await this.marriageService.claimWeddingGift(req.worldId, req.characterId, req.giftId, req.requestId),
        );
    }

    @Method("notifySpouseMap")
    async notifySpouseMap(call: GrpcCall<NotifySpouseMapRequest>, callback: GrpcCallback<NotifySpouseMapReply>) {
        const req = call.request;
        await this.marriageService.notifySpouseMap(
            req.worldId,
            req.characterId,
            req.spouseId,
            req.mapId,
            req.reply,
        );
        callback(null, {});
    }
}
