import { PartyErrorCode } from "../../protobuf/generated/fminternal/internal_service";
import type {
    BroadcastMultiChatReply,
    BroadcastMultiChatRequest,
} from "../../protobuf/generated/fminternal/internal_service";
import type { BroadcastBuddyMultiChatResult, BuddyService } from "../../services/buddy-service";
import type { BroadcastMultiChatResult, PartyService } from "../../services/party-service";
import type { GrpcCall, GrpcCallback, GrpcErrorHandler } from "./types";
import { GrpcController, GrpcMethod } from "../grpc-method-decorator";

const MULTI_CHAT_MODE_BUDDY = 0;
const MULTI_CHAT_MODE_PARTY = 1;

export function createChatHandlers(partyService: PartyService, buddyService: BuddyService, grpcError: GrpcErrorHandler) {
    return {
        async broadcastMultiChat(call: GrpcCall<BroadcastMultiChatRequest>, callback: GrpcCallback<BroadcastMultiChatReply>) {
            try {
                const req = call.request;
                const mode = req.chatMode;
                if (mode === MULTI_CHAT_MODE_BUDDY) {
                    const recipients = (req.recipientCharacterIds ?? []).map((id) => id >>> 0);
                    const result = (await buddyService.broadcastBuddyMultiChat(
                        req.worldId,
                        req.senderCharacterId,
                        recipients,
                        req.senderName,
                        req.message
                    )) as BroadcastBuddyMultiChatResult;
                    callback(null, {
                        ok: result.ok,
                        errorCode: result.ok ? PartyErrorCode.NONE : PartyErrorCode.UNKNOWN,
                        deliveredCount: result.deliveredCount ?? 0,
                    });
                    return;
                }
                if (mode === MULTI_CHAT_MODE_PARTY) {
                    const result = (await partyService.broadcastMultiChat(
                        req.worldId,
                        req.memberId,
                        req.senderCharacterId,
                        req.chatMode,
                        req.senderName,
                        req.message
                    )) as BroadcastMultiChatResult;
                    callback(null, {
                        ok: result.ok,
                        errorCode: result.code ?? PartyErrorCode.UNKNOWN,
                        deliveredCount: result.deliveredCount ?? 0,
                    });
                    return;
                }
                callback(null, {
                    ok: false,
                    errorCode: PartyErrorCode.UNKNOWN,
                    deliveredCount: 0,
                });
            } catch (err) {
                grpcError(err, callback);
            }
        },
    };
}

@GrpcController("chatController")
export class ChatGrpcController {
    private readonly handlers: ReturnType<typeof createChatHandlers>;

    constructor(partyService: PartyService, buddyService: BuddyService, grpcError: GrpcErrorHandler) {
        this.handlers = createChatHandlers(partyService, buddyService, grpcError);
    }

    @GrpcMethod("broadcastMultiChat")
    async broadcastMultiChat(call: GrpcCall<BroadcastMultiChatRequest>, callback: GrpcCallback<BroadcastMultiChatReply>) {
        return this.handlers.broadcastMultiChat(call, callback);
    }
}
