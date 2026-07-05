import type { ServerTimeService } from "../../services/server-time-service";
import type {
    SetServerDateTimeReply,
    SetServerDateTimeRequest,
} from "../../protobuf/generated/fminternal/internal_service";
import type { GrpcCall, GrpcCallback, GrpcErrorHandler } from "./types";
import { Controller, Method } from "../grpc-method-decorator";

@Controller("serverTimeController")
export class ServerTimeGrpcController {
    private readonly serverTimeService: ServerTimeService;
    private readonly grpcError: GrpcErrorHandler;

    constructor(serverTimeService: ServerTimeService, grpcError: GrpcErrorHandler) {
        this.serverTimeService = serverTimeService;
        this.grpcError = grpcError;
    }

    @Method("setServerDateTime")
    async setServerDateTime(
        call: GrpcCall<SetServerDateTimeRequest>,
        callback: GrpcCallback<SetServerDateTimeReply>,
    ) {
        const req = call.request;
        try {
            const ok = await this.serverTimeService.setServerDateTime(
                req.worldId >>> 0,
                req.datetime ?? "",
                req.reset === true,
            );
            callback(null, { ok });
        } catch (err) {
            this.grpcError(err, callback);
        }
    }
}
