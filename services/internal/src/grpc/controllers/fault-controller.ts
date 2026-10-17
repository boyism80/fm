import type { RabbitMQService } from "../../services/rabbitmq-service";
import type {
    SetPublishFaultReply,
    SetPublishFaultRequest,
} from "../../protobuf/generated/fminternal/internal_service";
import type { GrpcCall, GrpcCallback } from "./types";
import { Controller, Method } from "../grpc-method-decorator";

@Controller("faultController")
export class FaultGrpcController {
    private readonly rabbitmqService: RabbitMQService;

    constructor(rabbitmqService: RabbitMQService) {
        this.rabbitmqService = rabbitmqService;
    }

    @Method("setPublishFault")
    async setPublishFault(call: GrpcCall<SetPublishFaultRequest>, callback: GrpcCallback<SetPublishFaultReply>) {
        const enabled = call.request.enabled === true;
        this.rabbitmqService.setPublishFault(enabled);
        callback(null, { enabled });
    }
}
