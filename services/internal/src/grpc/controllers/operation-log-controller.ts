import type { OperationLogService } from "../../services/operation-log-service";
import type {
    WriteOperationLogReply,
    WriteOperationLogRequest,
} from "../../protobuf/generated/fminternal/internal_service";
import type { GrpcCall, GrpcCallback, GrpcErrorHandler } from "./types";
import { Controller, Method } from "../grpc-method-decorator";

@Controller("operationLogController")
export class OperationLogGrpcController {
    private readonly operationLogService: OperationLogService;
    private readonly grpcError: GrpcErrorHandler;

    constructor(operationLogService: OperationLogService, grpcError: GrpcErrorHandler) {
        this.operationLogService = operationLogService;
        this.grpcError = grpcError;
    }

    @Method("writeOperationLog")
    async writeOperationLog(
        call: GrpcCall<WriteOperationLogRequest>,
        callback: GrpcCallback<WriteOperationLogReply>,
    ) {
        const req = call.request;
        try {
            const ok = await this.operationLogService.write(
                req.worldId >>> 0,
                req.channelId >>> 0,
                req.characterId >>> 0,
                req.kind ?? "",
                req.meso ?? 0,
                req.detail ?? "",
            );
            callback(null, { ok });
        } catch (err) {
            this.grpcError(err, callback);
        }
    }
}
