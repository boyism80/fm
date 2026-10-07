import type { MarriageService } from "../../services/marriage-service";
import type {
    BreakEngagementRequest,
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
} from "../../protobuf/generated/fminternal/internal_service";
import type { GrpcCall, GrpcCallback, GrpcErrorHandler } from "./types";
import { Controller, Method } from "../grpc-method-decorator";
import { makeSaveCharacterEntry } from "./character-controller";

@Controller("marriageController")
export class MarriageGrpcController {
    private readonly marriageService: MarriageService;
    private readonly grpcError: GrpcErrorHandler;

    constructor(marriageService: MarriageService, grpcError: GrpcErrorHandler) {
        this.marriageService = marriageService;
        this.grpcError = grpcError;
    }

    @Method("createMarriage")
    async createMarriage(call: GrpcCall<CreateMarriageRequest>, callback: GrpcCallback<MarriageReply>) {
        try {
            const req = call.request;
            callback(
                null,
                await this.marriageService.createMarriage(req.worldId, req.groomId, req.groomName, req.brideId, req.brideName, req.ringItemId)
            );
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("loadMarriage")
    async loadMarriage(call: GrpcCall<LoadMarriageRequest>, callback: GrpcCallback<MarriageReply>) {
        try {
            callback(null, await this.marriageService.loadMarriage(call.request.worldId, call.request.characterId));
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("getMarriage")
    async getMarriage(call: GrpcCall<GetMarriageRequest>, callback: GrpcCallback<MarriageReply>) {
        try {
            callback(null, await this.marriageService.getMarriage(call.request.worldId, call.request.marriageId));
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("breakEngagement")
    async breakEngagement(call: GrpcCall<BreakEngagementRequest>, callback: GrpcCallback<MarriageReply>) {
        try {
            const req = call.request;
            callback(null, await this.marriageService.breakEngagement(req.worldId, req.marriageId, req.characterId));
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("reserveWedding")
    async reserveWedding(call: GrpcCall<ReserveWeddingRequest>, callback: GrpcCallback<MarriageReply>) {
        try {
            const req = call.request;
            callback(null, await this.marriageService.reserveWedding(req.worldId, req.marriageId, req.ticketItemId));
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("setWeddingWishlist")
    async setWeddingWishlist(call: GrpcCall<SetWeddingWishlistRequest>, callback: GrpcCallback<MarriageReply>) {
        try {
            const req = call.request;
            callback(null, await this.marriageService.setWeddingWishlist(req.worldId, req.marriageId, req.characterId, req.wishes ?? []));
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("inviteWeddingGuest")
    async inviteWeddingGuest(call: GrpcCall<InviteWeddingGuestRequest>, callback: GrpcCallback<InviteWeddingGuestReply>) {
        try {
            const req = call.request;
            callback(null, await this.marriageService.inviteWeddingGuest(req.worldId, req.marriageId, req.guestName));
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("finishWedding")
    async finishWedding(call: GrpcCall<FinishWeddingRequest>, callback: GrpcCallback<MarriageReply>) {
        try {
            callback(null, await this.marriageService.finishWedding(call.request.worldId, call.request.marriageId));
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("requestDivorce")
    async requestDivorce(call: GrpcCall<RequestDivorceRequest>, callback: GrpcCallback<MarriageReply>) {
        try {
            const req = call.request;
            callback(null, await this.marriageService.requestDivorce(req.worldId, req.marriageId, req.characterId, Number(req.nowUnixMs) || Date.now()));
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("giveWeddingGift")
    async giveWeddingGift(call: GrpcCall<GiveWeddingGiftRequest>, callback: GrpcCallback<GiveWeddingGiftReply>) {
        try {
            const req = call.request;
            if (!req.gift) {
                throw Object.assign(new Error("gift is required"), { code: "INVALID_PAYLOAD" });
            }
            const result = await this.marriageService.giveWeddingGift({
                worldId: req.worldId,
                receiverId: req.receiverId,
                gift: req.gift,
                sender: req.sender ? makeSaveCharacterEntry(req.sender) : undefined,
            });
            callback(null, { result });
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("loadWeddingGifts")
    async loadWeddingGifts(call: GrpcCall<LoadWeddingGiftsRequest>, callback: GrpcCallback<LoadWeddingGiftsReply>) {
        try {
            callback(null, await this.marriageService.loadWeddingGifts(call.request.worldId, call.request.characterId));
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("claimWeddingGift")
    async claimWeddingGift(call: GrpcCall<ClaimWeddingGiftRequest>, callback: GrpcCallback<ClaimWeddingGiftReply>) {
        try {
            const req = call.request;
            callback(null, await this.marriageService.claimWeddingGift(req.worldId, req.characterId, req.giftId));
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("notifySpouseMap")
    async notifySpouseMap(call: GrpcCall<NotifySpouseMapRequest>, callback: GrpcCallback<NotifySpouseMapReply>) {
        try {
            const req = call.request;
            await this.marriageService.notifySpouseMap(req.worldId, req.characterId, req.spouseId, req.mapId, req.reply);
            callback(null, {});
        } catch (err) {
            this.grpcError(err, callback);
        }
    }
}
