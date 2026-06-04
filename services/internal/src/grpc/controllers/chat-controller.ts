import { PartyErrorCode } from "../../protobuf/generated/fminternal/internal_service";
import type {
    BroadcastMultiChatReply,
    BroadcastMultiChatRequest,
} from "../../protobuf/generated/fminternal/internal_service";
import type { BroadcastBuddyMultiChatResult, BuddyService } from "../../services/buddy-service";
import type { GuildService } from "../../services/guild-service";
import type { BroadcastMultiChatResult, PartyService } from "../../services/party-service";
import type { GrpcCall, GrpcCallback, GrpcErrorHandler } from "./types";
import { Controller, Method } from "../grpc-method-decorator";

const MULTI_CHAT_MODE_BUDDY = 0;
const MULTI_CHAT_MODE_PARTY = 1;
const MULTI_CHAT_MODE_GUILD = 2;
const MULTI_CHAT_MODE_ALLIANCE = 3;

@Controller("chatController")
export class ChatGrpcController {
    private readonly partyService: PartyService;
    private readonly buddyService: BuddyService;
    private readonly guildService: GuildService;
    private readonly grpcError: GrpcErrorHandler;

    constructor(
        partyService: PartyService,
        buddyService: BuddyService,
        guildService: GuildService,
        grpcError: GrpcErrorHandler
    ) {
        this.partyService = partyService;
        this.buddyService = buddyService;
        this.guildService = guildService;
        this.grpcError = grpcError;
    }

    @Method("broadcastMultiChat")
    async broadcastMultiChat(call: GrpcCall<BroadcastMultiChatRequest>, callback: GrpcCallback<BroadcastMultiChatReply>) {
        try {
            const req = call.request;
            const mode = req.chatMode;
            if (mode === MULTI_CHAT_MODE_BUDDY) {
                const recipients = (req.recipientCharacterIds ?? []).map((id) => id >>> 0);
                const result = (await this.buddyService.broadcastBuddyMultiChat(
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
                const result = (await this.partyService.broadcastMultiChat(
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
            if (mode === MULTI_CHAT_MODE_GUILD) {
                const result = await this.guildService.broadcastGuildMultiChat(
                    req.worldId,
                    req.memberId,
                    req.senderCharacterId,
                    req.senderName,
                    req.message
                );
                callback(null, {
                    ok: result.ok,
                    errorCode: result.ok ? PartyErrorCode.NONE : PartyErrorCode.UNKNOWN,
                    deliveredCount: result.deliveredCount ?? 0,
                });
                return;
            }
            if (mode === MULTI_CHAT_MODE_ALLIANCE) {
                const result = await this.guildService.broadcastAllianceMultiChat(
                    req.worldId,
                    req.memberId,
                    req.senderCharacterId,
                    req.senderName,
                    req.message
                );
                callback(null, {
                    ok: result.ok,
                    errorCode: result.ok ? PartyErrorCode.NONE : PartyErrorCode.UNKNOWN,
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
            this.grpcError(err, callback);
        }
    }
}
