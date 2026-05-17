import {
    BuddyErrorCode,
    type AcceptBuddyReply,
    type AcceptBuddyRequest,
    type BuddyEntry,
    type RemoveBuddyReply,
    type RemoveBuddyRequest,
    type RequestBuddyReply,
    type RequestBuddyRequest,
} from "../../protobuf/generated/fminternal/internal_service";
import { BUDDY_ENTRY, BUDDY_LIST_ENTRY } from "../buddy-persisted";
import { grpcMapper } from "../mappers";
import type { BuddyListEntry, BuddyService } from "../../services/buddy-service";
import type { GrpcCall, GrpcCallback, GrpcErrorHandler } from "./types";
import { GrpcController, GrpcMethod } from "../grpc-method-decorator";

export function createBuddyHandlers(buddyService: BuddyService, grpcError: GrpcErrorHandler) {
    return {
        async requestBuddy(call: GrpcCall<RequestBuddyRequest>, callback: GrpcCallback<RequestBuddyReply>) {
            try {
                const req = call.request;
                const result = await buddyService.requestBuddy(
                    req.worldId,
                    req.requesterCharacterId,
                    req.targetCharacterName,
                    req.groupName
                );
                callback(null, {
                    ok: result.ok,
                    errorCode: result.ok
                        ? BuddyErrorCode.BUDDY_ERROR_NONE
                        : result.code ?? BuddyErrorCode.BUDDY_ERROR_UNKNOWN,
                    targetCharacterId: result.targetCharacterId ?? 0,
                    targetChannelId: result.targetChannelId ?? 0,
                    requesterView: result.requesterView
                        ? grpcMapper.map<BuddyListEntry, BuddyEntry>(
                              result.requesterView,
                              BUDDY_LIST_ENTRY,
                              BUDDY_ENTRY
                          )
                        : undefined,
                });
            } catch (err) {
                grpcError(err, callback);
            }
        },
        async acceptBuddy(call: GrpcCall<AcceptBuddyRequest>, callback: GrpcCallback<AcceptBuddyReply>) {
            try {
                const req = call.request;
                const result = await buddyService.acceptBuddy(
                    req.worldId,
                    req.accepterCharacterId,
                    req.requesterCharacterId
                );
                callback(null, {
                    ok: result.ok,
                    errorCode: result.ok
                        ? BuddyErrorCode.BUDDY_ERROR_NONE
                        : result.code ?? BuddyErrorCode.BUDDY_ERROR_UNKNOWN,
                    requesterCharacterId: result.requesterCharacterId ?? 0,
                    requesterChannelId: result.requesterChannelId ?? 0,
                    accepterView: result.accepterView
                        ? grpcMapper.map<BuddyListEntry, BuddyEntry>(
                              result.accepterView,
                              BUDDY_LIST_ENTRY,
                              BUDDY_ENTRY
                          )
                        : undefined,
                    requesterView: result.requesterView
                        ? grpcMapper.map<BuddyListEntry, BuddyEntry>(
                              result.requesterView,
                              BUDDY_LIST_ENTRY,
                              BUDDY_ENTRY
                          )
                        : undefined,
                });
            } catch (err) {
                grpcError(err, callback);
            }
        },
        async removeBuddy(call: GrpcCall<RemoveBuddyRequest>, callback: GrpcCallback<RemoveBuddyReply>) {
            try {
                const req = call.request;
                const result = await buddyService.removeBuddy(
                    req.worldId,
                    req.characterId,
                    req.buddyCharacterId
                );
                callback(null, {
                    ok: result.ok,
                    errorCode: result.ok
                        ? BuddyErrorCode.BUDDY_ERROR_NONE
                        : result.code ?? BuddyErrorCode.BUDDY_ERROR_UNKNOWN,
                    buddyCharacterId: result.buddyCharacterId ?? 0,
                    buddyChannelId: result.buddyChannelId ?? 0,
                    buddyWasAccepted: result.buddyWasAccepted ?? false,
                    removedFromOwner: result.removedFromOwner ?? false,
                });
            } catch (err) {
                grpcError(err, callback);
            }
        },
    };
}

@GrpcController("buddyController")
export class BuddyGrpcController {
    private readonly handlers: ReturnType<typeof createBuddyHandlers>;

    constructor(buddyService: BuddyService, grpcError: GrpcErrorHandler) {
        this.handlers = createBuddyHandlers(buddyService, grpcError);
    }

    @GrpcMethod("requestBuddy")
    async requestBuddy(call: GrpcCall<RequestBuddyRequest>, callback: GrpcCallback<RequestBuddyReply>) {
        return this.handlers.requestBuddy(call, callback);
    }

    @GrpcMethod("acceptBuddy")
    async acceptBuddy(call: GrpcCall<AcceptBuddyRequest>, callback: GrpcCallback<AcceptBuddyReply>) {
        return this.handlers.acceptBuddy(call, callback);
    }

    @GrpcMethod("removeBuddy")
    async removeBuddy(call: GrpcCall<RemoveBuddyRequest>, callback: GrpcCallback<RemoveBuddyReply>) {
        return this.handlers.removeBuddy(call, callback);
    }
}
