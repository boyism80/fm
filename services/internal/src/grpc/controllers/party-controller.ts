import { PartyErrorCode } from "@fm/shared/protobuf/generated/fminternal/internal_service";
import type { PartyService } from "../../services/party-service";
import type {
    ChangePartyLeaderReply,
    ChangePartyLeaderRequest,
    CreatePartyReply,
    CreatePartyRequest,
    DenyPartyReply,
    DenyPartyRequest,
    ExpelPartyReply,
    ExpelPartyRequest,
    GetPartyReply,
    GetPartyRequest,
    InvitePartyReply,
    InvitePartyRequest,
    JoinPartyReply,
    JoinPartyRequest,
    LeavePartyReply,
    LeavePartyRequest,
    UpdatePartyMemberReply,
    UpdatePartyMemberRequest,
} from "@fm/shared/protobuf/generated/fminternal/internal_service";
import { okReply } from "./types";
import type { GrpcCall, GrpcCallback } from "./types";
import { Controller, Method } from "../grpc-method-decorator";

@Controller("partyController")
export class PartyGrpcController {
    private readonly partyService: PartyService;

    constructor(partyService: PartyService) {
        this.partyService = partyService;
    }

    @Method("createParty")
    async createParty(call: GrpcCall<CreatePartyRequest>, callback: GrpcCallback<CreatePartyReply>) {
        const req = call.request;
        const result = await this.partyService.createParty(req.worldId, req.leader);
        callback(null, { ...okReply(result, PartyErrorCode.NONE, PartyErrorCode.UNKNOWN), party: result.party });
    }

    @Method("joinParty")
    async joinParty(call: GrpcCall<JoinPartyRequest>, callback: GrpcCallback<JoinPartyReply>) {
        const req = call.request;
        const result = await this.partyService.joinParty(
            req.worldId,
            req.partyId,
            req.member,
            req.skipInvitePendingCheck,
        );
        console.log(
            `[party] join character=${req.member?.characterId} party=${req.partyId} skip_pending=${req.skipInvitePendingCheck}` +
                ` ok=${result.ok} code=${PartyErrorCode[result.code ?? PartyErrorCode.NONE]}`,
        );
        callback(null, {
            ...okReply(result, PartyErrorCode.NONE, PartyErrorCode.UNKNOWN),
            partyId: result.partyId,
            revision: result.revision ?? 0,
        });
    }

    @Method("inviteParty")
    async inviteParty(call: GrpcCall<InvitePartyRequest>, callback: GrpcCallback<InvitePartyReply>) {
        const result = await this.partyService.inviteParty(
            call.request.worldId,
            call.request.inviterCharacterId,
            call.request.targetCharacterName,
        );
        console.log(
            `[party] invite inviter=${call.request.inviterCharacterId} target="${call.request.targetCharacterName}"` +
                ` ok=${result.ok} code=${PartyErrorCode[result.code ?? PartyErrorCode.NONE]} party=${result.partyId}` +
                ` target_id=${result.targetCharacterId} target_channel=${result.targetChannelId}`,
        );
        callback(null, {
            ok: result.ok,
            errorCode: result.code ?? PartyErrorCode.UNKNOWN,
            targetCharacterId: result.targetCharacterId ?? 0,
            targetChannelId: result.targetChannelId ?? 0,
            partyId: result.partyId,
        });
    }

    @Method("leaveParty")
    async leaveParty(call: GrpcCall<LeavePartyRequest>, callback: GrpcCallback<LeavePartyReply>) {
        const result = await this.partyService.leaveParty(call.request.worldId, call.request.characterId);
        callback(null, {
            ...okReply(result, PartyErrorCode.NONE, PartyErrorCode.UNKNOWN),
            partyId: result.partyId,
            disbanded: result.disbanded ?? false,
            revision: result.revision ?? 0,
        });
    }

    @Method("expelParty")
    async expelParty(call: GrpcCall<ExpelPartyRequest>, callback: GrpcCallback<ExpelPartyReply>) {
        const result = await this.partyService.expelParty(
            call.request.worldId,
            call.request.requesterCharacterId,
            call.request.targetCharacterId,
        );
        callback(null, {
            ...okReply(result, PartyErrorCode.NONE, PartyErrorCode.UNKNOWN),
            partyId: result.partyId,
            disbanded: result.disbanded ?? false,
            revision: result.revision ?? 0,
        });
    }

    @Method("changePartyLeader")
    async changePartyLeader(call: GrpcCall<ChangePartyLeaderRequest>, callback: GrpcCallback<ChangePartyLeaderReply>) {
        const result = await this.partyService.changePartyLeader(
            call.request.worldId,
            call.request.partyId,
            call.request.requesterCharacterId,
            call.request.newLeaderCharacterId,
        );
        callback(null, okReply(result, PartyErrorCode.NONE, PartyErrorCode.UNKNOWN));
    }

    @Method("getParty")
    async getParty(call: GrpcCall<GetPartyRequest>, callback: GrpcCallback<GetPartyReply>) {
        const worldId = call.request.worldId;
        const result = await this.partyService.getParty(worldId, call.request.partyId);
        if (!result.party) {
            callback(null, { found: false, party: undefined });
            return;
        }
        const party = await this.partyService.partyToPb(worldId, result.party, result.members ?? []);
        callback(null, { found: true, party });
    }

    @Method("updatePartyMember")
    async updatePartyMember(call: GrpcCall<UpdatePartyMemberRequest>, callback: GrpcCallback<UpdatePartyMemberReply>) {
        const result = await this.partyService.updatePartyMember(call.request.member);
        callback(null, okReply(result, PartyErrorCode.NONE, PartyErrorCode.UNKNOWN));
    }

    @Method("denyParty")
    async denyParty(call: GrpcCall<DenyPartyRequest>, callback: GrpcCallback<DenyPartyReply>) {
        const result = await this.partyService.denyParty(
            call.request.worldId,
            call.request.deniedCharacterId,
            call.request.inviterName,
            call.request.action,
        );
        console.log(
            `[party] deny denied=${call.request.deniedCharacterId} inviter="${call.request.inviterName}" action=${call.request.action}` +
                ` ok=${result.ok} code=${PartyErrorCode[result.code ?? PartyErrorCode.NONE]}`,
        );
        callback(null, { ok: result.ok, errorCode: result.code ?? PartyErrorCode.UNKNOWN });
    }
}
