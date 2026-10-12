import type {
    CompleteWhisperReply,
    CompleteWhisperRequest,
    FamilyReply,
    FinishFamilyMoveRequest,
    GetFamilyReply,
    GetFamilyRequest,
    JoinFamilyRequest,
    MessengerReply,
    MessengerRequest,
    NotesReply,
    NotesRequest,
    SeparateFamilyRequest,
    SetFamilyPreceptRequest,
    UseFamilyFeatureRequest,
    WhisperReply,
    WhisperRequest,
} from "../../protobuf/generated/fminternal/internal_service";
import type { FamilyService } from "../../services/family-service";
import type { SocialService } from "../../services/social-service";
import type { GrpcCall, GrpcCallback, GrpcErrorHandler } from "./types";
import { Controller, Method } from "../grpc-method-decorator";

@Controller("socialController")
export class SocialGrpcController {
    private readonly socialService: SocialService;
    private readonly familyService: FamilyService;
    private readonly grpcError: GrpcErrorHandler;

    constructor(socialService: SocialService, familyService: FamilyService, grpcError: GrpcErrorHandler) {
        this.socialService = socialService;
        this.familyService = familyService;
        this.grpcError = grpcError;
    }

    @Method("whisper")
    async whisper(call: GrpcCall<WhisperRequest>, callback: GrpcCallback<WhisperReply>) {
        try {
            callback(null, await this.socialService.whisper(call.request));
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("completeWhisper")
    async completeWhisper(call: GrpcCall<CompleteWhisperRequest>, callback: GrpcCallback<CompleteWhisperReply>) {
        try {
            callback(null, await this.socialService.completeWhisper(call.request));
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("messenger")
    async messenger(call: GrpcCall<MessengerRequest>, callback: GrpcCallback<MessengerReply>) {
        try {
            callback(null, await this.socialService.messenger(call.request));
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("notes")
    async notes(call: GrpcCall<NotesRequest>, callback: GrpcCallback<NotesReply>) {
        try {
            callback(null, await this.socialService.notes(call.request));
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("getFamily")
    async getFamily(call: GrpcCall<GetFamilyRequest>, callback: GrpcCallback<GetFamilyReply>) {
        try {
            callback(null, await this.familyService.get(call.request));
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("joinFamily")
    async joinFamily(call: GrpcCall<JoinFamilyRequest>, callback: GrpcCallback<FamilyReply>) {
        try {
            callback(null, await this.familyService.join(call.request));
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("separateFamily")
    async separateFamily(call: GrpcCall<SeparateFamilyRequest>, callback: GrpcCallback<FamilyReply>) {
        try {
            callback(null, await this.familyService.separate(call.request));
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("setFamilyPrecept")
    async setFamilyPrecept(call: GrpcCall<SetFamilyPreceptRequest>, callback: GrpcCallback<FamilyReply>) {
        try {
            callback(null, await this.familyService.setPrecept(call.request));
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("useFamilyFeature")
    async useFamilyFeature(call: GrpcCall<UseFamilyFeatureRequest>, callback: GrpcCallback<FamilyReply>) {
        try {
            callback(null, await this.familyService.use(call.request));
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("finishFamilyMove")
    async finishFamilyMove(call: GrpcCall<FinishFamilyMoveRequest>, callback: GrpcCallback<FamilyReply>) {
        try {
            callback(null, await this.familyService.finishMove(call.request));
        } catch (err) {
            this.grpcError(err, callback);
        }
    }
}
