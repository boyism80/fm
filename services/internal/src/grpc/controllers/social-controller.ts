import type {
    CompleteWhisperReply,
    CompleteWhisperRequest,
    FamilyReply,
    FamilyRequest,
    MessengerReply,
    MessengerRequest,
    NotesReply,
    NotesRequest,
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

    @Method("family")
    async family(call: GrpcCall<FamilyRequest>, callback: GrpcCallback<FamilyReply>) {
        try {
            callback(null, await this.familyService.family(call.request));
        } catch (err) {
            this.grpcError(err, callback);
        }
    }
}
