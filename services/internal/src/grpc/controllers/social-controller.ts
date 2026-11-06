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
    UpdateFamilyMemberReply,
    UpdateFamilyMemberRequest,
    UseFamilyFeatureRequest,
    WhisperReply,
    WhisperRequest,
} from "@fm/shared/protobuf/generated/fminternal/internal_service";
import type { FamilyService } from "../../services/family-service";
import type { SocialService } from "../../services/social-service";
import type { GrpcCall, GrpcCallback } from "./types";
import { Controller, Method } from "../grpc-method-decorator";

@Controller("socialController")
export class SocialGrpcController {
    private readonly socialService: SocialService;
    private readonly familyService: FamilyService;

    constructor(socialService: SocialService, familyService: FamilyService) {
        this.socialService = socialService;
        this.familyService = familyService;
    }

    @Method("whisper")
    async whisper(call: GrpcCall<WhisperRequest>, callback: GrpcCallback<WhisperReply>) {
        callback(null, await this.socialService.whisper(call.request));
    }

    @Method("completeWhisper")
    async completeWhisper(call: GrpcCall<CompleteWhisperRequest>, callback: GrpcCallback<CompleteWhisperReply>) {
        callback(null, await this.socialService.completeWhisper(call.request));
    }

    @Method("messenger")
    async messenger(call: GrpcCall<MessengerRequest>, callback: GrpcCallback<MessengerReply>) {
        callback(null, await this.socialService.messenger(call.request));
    }

    @Method("notes")
    async notes(call: GrpcCall<NotesRequest>, callback: GrpcCallback<NotesReply>) {
        callback(null, await this.socialService.notes(call.request));
    }

    @Method("getFamily")
    async getFamily(call: GrpcCall<GetFamilyRequest>, callback: GrpcCallback<GetFamilyReply>) {
        callback(null, await this.familyService.get(call.request));
    }

    @Method("joinFamily")
    async joinFamily(call: GrpcCall<JoinFamilyRequest>, callback: GrpcCallback<FamilyReply>) {
        callback(null, await this.familyService.join(call.request));
    }

    @Method("separateFamily")
    async separateFamily(call: GrpcCall<SeparateFamilyRequest>, callback: GrpcCallback<FamilyReply>) {
        callback(null, await this.familyService.separate(call.request));
    }

    @Method("setFamilyPrecept")
    async setFamilyPrecept(call: GrpcCall<SetFamilyPreceptRequest>, callback: GrpcCallback<FamilyReply>) {
        callback(null, await this.familyService.setPrecept(call.request));
    }

    @Method("useFamilyFeature")
    async useFamilyFeature(call: GrpcCall<UseFamilyFeatureRequest>, callback: GrpcCallback<FamilyReply>) {
        callback(null, await this.familyService.use(call.request));
    }

    @Method("finishFamilyMove")
    async finishFamilyMove(call: GrpcCall<FinishFamilyMoveRequest>, callback: GrpcCallback<FamilyReply>) {
        callback(null, await this.familyService.finishMove(call.request));
    }

    @Method("updateFamilyMember")
    async updateFamilyMember(
        call: GrpcCall<UpdateFamilyMemberRequest>,
        callback: GrpcCallback<UpdateFamilyMemberReply>,
    ) {
        await this.familyService.updateMember(call.request);
        callback(null, {});
    }
}
