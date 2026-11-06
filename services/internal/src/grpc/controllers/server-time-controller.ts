import type { ServerTimeService } from "../../services/server-time-service";
import type {
    BroadcastNoticeReply,
    BroadcastNoticeRequest,
    SetServerDateTimeReply,
    SetServerDateTimeRequest,
} from "@fm/shared/protobuf/generated/fminternal/internal_service";
import type { GrpcCall, GrpcCallback } from "./types";
import { Controller, Method } from "../grpc-method-decorator";

@Controller("serverTimeController")
export class ServerTimeGrpcController {
    private readonly serverTimeService: ServerTimeService;

    constructor(serverTimeService: ServerTimeService) {
        this.serverTimeService = serverTimeService;
    }

    @Method("setServerDateTime")
    async setServerDateTime(call: GrpcCall<SetServerDateTimeRequest>, callback: GrpcCallback<SetServerDateTimeReply>) {
        const req = call.request;
        const ok = await this.serverTimeService.setServerDateTime(
            req.worldId >>> 0,
            req.datetime ?? "",
            req.reset === true,
        );
        callback(null, { ok });
    }

    @Method("broadcastNotice")
    async broadcastNotice(call: GrpcCall<BroadcastNoticeRequest>, callback: GrpcCallback<BroadcastNoticeReply>) {
        const req = call.request;
        const ok = await this.serverTimeService.broadcastNotice(
            req.worldId >>> 0,
            req.sourceChannelId >>> 0,
            req.messageType >>> 0,
            req.message ?? "",
            req.megaEar === true,
            req.nextLines ?? [],
            req.item ?? Buffer.alloc(0),
            req.adminsOnly === true,
        );
        callback(null, { ok });
    }
}
