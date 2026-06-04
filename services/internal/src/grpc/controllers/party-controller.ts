import { PartyErrorCode } from "../../protobuf/generated/fminternal/internal_service";
import type {
    DenyPartyResult,
    ExpelPartyResult,
    GetPartyResult,
    InvitePartyResult,
    LeavePartyResult,
    PartyMutationResult,
    PartyService,
} from "../../services/party-service";
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
} from "../../protobuf/generated/fminternal/internal_service";
import type { GrpcCall, GrpcCallback, GrpcErrorHandler } from "./types";
import { Controller, Method } from "../grpc-method-decorator";

@Controller("partyController")
export class PartyGrpcController {
    private readonly partyService: PartyService;
    private readonly grpcError: GrpcErrorHandler;

    constructor(partyService: PartyService, grpcError: GrpcErrorHandler) {
        this.partyService = partyService;
        this.grpcError = grpcError;
    }

    @Method("createParty")
    async createParty(call: GrpcCall<CreatePartyRequest>, callback: GrpcCallback<CreatePartyReply>) {
        try {
            const req = call.request;
            const result = await this.partyService.createParty(req.worldId, req.leader) as PartyMutationResult;
            if (result.ok) {
                callback(null, {
                    ok: true,
                    errorCode: PartyErrorCode.NONE,
                    partyId: result.partyId,
                    revision: result.revision ?? 0,
                });
            } else {
                callback(null, {
                    ok: false,
                    errorCode: result.code ?? PartyErrorCode.UNKNOWN,
                    partyId: undefined,
                    revision: 0,
                });
            }
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("joinParty")
    async joinParty(call: GrpcCall<JoinPartyRequest>, callback: GrpcCallback<JoinPartyReply>) {
        try {
            const req = call.request;
            const result = await this.partyService.joinParty(req.worldId, req.partyId, req.member, req.skipInvitePendingCheck) as PartyMutationResult;
            if (result.ok) {
                callback(null, {
                    ok: true,
                    errorCode: PartyErrorCode.NONE,
                    partyId: result.partyId,
                    revision: result.revision ?? 0,
                });
            } else {
                callback(null, {
                    ok: false,
                    errorCode: result.code ?? PartyErrorCode.UNKNOWN,
                    partyId: undefined,
                    revision: 0,
                });
            }
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("inviteParty")
    async inviteParty(call: GrpcCall<InvitePartyRequest>, callback: GrpcCallback<InvitePartyReply>) {
        try {
            const result = await this.partyService.inviteParty(call.request.worldId, call.request.inviterCharacterId, call.request.targetCharacterName) as InvitePartyResult;
            callback(null, {
                ok: result.ok,
                errorCode: result.code ?? PartyErrorCode.UNKNOWN,
                targetCharacterId: result.targetCharacterId ?? 0,
                targetChannelId: result.targetChannelId ?? 0,
                partyId: result.partyId,
            });
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("leaveParty")
    async leaveParty(call: GrpcCall<LeavePartyRequest>, callback: GrpcCallback<LeavePartyReply>) {
        try {
            const result = await this.partyService.leaveParty(call.request.worldId, call.request.characterId) as LeavePartyResult;
            if (result.ok) {
                callback(null, {
                    ok: true,
                    errorCode: PartyErrorCode.NONE,
                    partyId: result.partyId,
                    revision: result.revision ?? 0,
                    disbanded: result.disbanded ?? false,
                });
            } else {
                callback(null, {
                    ok: false,
                    errorCode: result.code ?? PartyErrorCode.UNKNOWN,
                    partyId: undefined,
                    revision: 0,
                    disbanded: false,
                });
            }
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("expelParty")
    async expelParty(call: GrpcCall<ExpelPartyRequest>, callback: GrpcCallback<ExpelPartyReply>) {
        try {
            const result = await this.partyService.expelParty(call.request.worldId, call.request.requesterCharacterId, call.request.targetCharacterId) as ExpelPartyResult;
            if (result.ok) {
                callback(null, {
                    ok: true,
                    errorCode: PartyErrorCode.NONE,
                    partyId: result.partyId,
                    revision: result.revision ?? 0,
                    disbanded: result.disbanded ?? false,
                });
            } else {
                callback(null, {
                    ok: false,
                    errorCode: result.code ?? PartyErrorCode.UNKNOWN,
                    partyId: undefined,
                    revision: 0,
                    disbanded: false,
                });
            }
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("changePartyLeader")
    async changePartyLeader(call: GrpcCall<ChangePartyLeaderRequest>, callback: GrpcCallback<ChangePartyLeaderReply>) {
        try {
            const result = await this.partyService.changePartyLeader(
                call.request.worldId,
                call.request.partyId,
                call.request.requesterCharacterId,
                call.request.newLeaderCharacterId
            ) as PartyMutationResult;
            if (result.ok) {
                callback(null, {
                    ok: true,
                    errorCode: PartyErrorCode.NONE,
                    partyId: result.partyId,
                    revision: result.revision ?? 0,
                });
            } else {
                callback(null, {
                    ok: false,
                    errorCode: result.code ?? PartyErrorCode.UNKNOWN,
                    partyId: undefined,
                    revision: 0,
                });
            }
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("getParty")
    async getParty(call: GrpcCall<GetPartyRequest>, callback: GrpcCallback<GetPartyReply>) {
        try {
            const worldId = call.request.worldId;
            const result = await this.partyService.getParty(worldId, call.request.partyId) as GetPartyResult;
            if (!result.party) {
                callback(null, { found: false, party: undefined });
                return;
            }
            const party = await this.partyService.buildPartyMessage(worldId, result.party, result.members ?? []);
            callback(null, { found: true, party });
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("updatePartyMember")
    async updatePartyMember(call: GrpcCall<UpdatePartyMemberRequest>, callback: GrpcCallback<UpdatePartyMemberReply>) {
        try {
            const result = await this.partyService.updatePartyMember(call.request.member) as PartyMutationResult;
            if (result.ok) {
                callback(null, {
                    ok: true,
                    errorCode: PartyErrorCode.NONE,
                    partyId: result.partyId,
                    revision: result.revision ?? 0,
                });
            } else {
                callback(null, {
                    ok: false,
                    errorCode: result.code ?? PartyErrorCode.UNKNOWN,
                    partyId: undefined,
                    revision: 0,
                });
            }
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("denyParty")
    async denyParty(call: GrpcCall<DenyPartyRequest>, callback: GrpcCallback<DenyPartyReply>) {
        try {
            const result = await this.partyService.denyParty(call.request.worldId, call.request.deniedCharacterId, call.request.inviterName, call.request.action) as DenyPartyResult;
            callback(null, { ok: result.ok, errorCode: result.code ?? PartyErrorCode.UNKNOWN });
        } catch (err) {
            this.grpcError(err, callback);
        }
    }
}
