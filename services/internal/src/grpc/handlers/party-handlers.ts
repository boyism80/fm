import { PartyErrorCode, type Party, type PartyMember } from "../../protobuf/generated/fminternal/internal_service";
import type {
    BroadcastMultiChatResult,
    DenyPartyResult,
    ExpelPartyResult,
    GetPartyResult,
    InvitePartyResult,
    LeavePartyResult,
    PartyMutationResult,
    PartyService,
} from "../../services/party-service";
import type {
    BroadcastMultiChatReply,
    BroadcastMultiChatRequest,
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

export function createPartyHandlers(partyService: PartyService, grpcError: GrpcErrorHandler) {
    return {
        async createParty(call: GrpcCall<CreatePartyRequest>, callback: GrpcCallback<CreatePartyReply>) {
            try {
                const req = call.request;
                const result = await partyService.createParty(req.worldId, req.leader) as PartyMutationResult;
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
            } catch (err) { grpcError(err, callback); }
        },
        async joinParty(call: GrpcCall<JoinPartyRequest>, callback: GrpcCallback<JoinPartyReply>) {
            try {
                const req = call.request;
                const result = await partyService.joinParty(req.worldId, req.partyId, req.member, req.skipInvitePendingCheck) as PartyMutationResult;
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
            } catch (err) { grpcError(err, callback); }
        },
        async inviteParty(call: GrpcCall<InvitePartyRequest>, callback: GrpcCallback<InvitePartyReply>) {
            try {
                const result = await partyService.inviteParty(call.request.worldId, call.request.inviterCharacterId, call.request.targetCharacterName) as InvitePartyResult;
                callback(null, {
                    ok: Boolean(result.ok),
                    errorCode: result.code ?? PartyErrorCode.UNKNOWN,
                    targetCharacterId: result.targetCharacterId ?? 0,
                    targetChannelId: result.targetChannelId ?? 0,
                    partyId: result.partyId,
                });
            } catch (err) { grpcError(err, callback); }
        },
        async leaveParty(call: GrpcCall<LeavePartyRequest>, callback: GrpcCallback<LeavePartyReply>) {
            try {
                const result = await partyService.leaveParty(call.request.worldId, call.request.characterId) as LeavePartyResult;
                if (result.ok) {
                    callback(null, {
                        ok: true,
                        errorCode: PartyErrorCode.NONE,
                        partyId: result.partyId,
                        revision: result.revision ?? 0,
                        disbanded: Boolean(result.disbanded),
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
            } catch (err) { grpcError(err, callback); }
        },
        async expelParty(call: GrpcCall<ExpelPartyRequest>, callback: GrpcCallback<ExpelPartyReply>) {
            try {
                const result = await partyService.expelParty(call.request.worldId, call.request.requesterCharacterId, call.request.targetCharacterId) as ExpelPartyResult;
                if (result.ok) {
                    callback(null, {
                        ok: true,
                        errorCode: PartyErrorCode.NONE,
                        partyId: result.partyId,
                        revision: result.revision ?? 0,
                        disbanded: Boolean(result.disbanded),
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
            } catch (err) { grpcError(err, callback); }
        },
        async changePartyLeader(call: GrpcCall<ChangePartyLeaderRequest>, callback: GrpcCallback<ChangePartyLeaderReply>) {
            try {
                const result = await partyService.changePartyLeader(
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
            } catch (err) { grpcError(err, callback); }
        },
        async getParty(call: GrpcCall<GetPartyRequest>, callback: GrpcCallback<GetPartyReply>) {
            try {
                const worldId = call.request.worldId;
                const result = await partyService.getParty(worldId, call.request.partyId) as GetPartyResult;
                if (result.found && result.party) {
                    const party: Party = {
                        worldId: result.party.worldId,
                        partyId: result.party.partyId,
                        leaderCharacterId: result.party.leaderCharacterId,
                        revision: result.party.revision,
                        state: result.party.state,
                        members: (result.members ?? []).map((m): PartyMember => ({
                            worldId,
                            characterId: m.characterId,
                            characterName: m.characterName,
                            level: m.level,
                            classId: m.classId,
                            role: m.role,
                            mapId: m.mapId ?? 0,
                            channelIndex: m.channelIndex ?? -2,
                            door: m.door ? {
                                town: m.door.town,
                                target: m.door.target,
                                x: m.door.x,
                                y: m.door.y,
                            } : undefined,
                        })),
                    };
                    callback(null, { found: true, party });
                } else {
                    callback(null, { found: false, party: undefined });
                }
            } catch (err) { grpcError(err, callback); }
        },
        async updatePartyMember(call: GrpcCall<UpdatePartyMemberRequest>, callback: GrpcCallback<UpdatePartyMemberReply>) {
            try {
                const result = await partyService.updatePartyMember(call.request.member) as PartyMutationResult;
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
            } catch (err) { grpcError(err, callback); }
        },
        async denyParty(call: GrpcCall<DenyPartyRequest>, callback: GrpcCallback<DenyPartyReply>) {
            try {
                const result = await partyService.denyParty(call.request.worldId, call.request.deniedCharacterId, call.request.inviterName, call.request.action) as DenyPartyResult;
                callback(null, { ok: Boolean(result.ok), errorCode: result.code ?? PartyErrorCode.UNKNOWN });
            } catch (err) { grpcError(err, callback); }
        },
        async broadcastMultiChat(call: GrpcCall<BroadcastMultiChatRequest>, callback: GrpcCallback<BroadcastMultiChatReply>) {
            try {
                const req = call.request;
                const result = await partyService.broadcastMultiChat(req.worldId, req.memberId, req.senderCharacterId, req.chatMode, req.senderName, req.message) as BroadcastMultiChatResult;
                callback(null, {
                    ok: Boolean(result.ok),
                    errorCode: result.code ?? PartyErrorCode.UNKNOWN,
                    deliveredCount: result.deliveredCount ?? 0,
                });
            } catch (err) { grpcError(err, callback); }
        },
    };
}
