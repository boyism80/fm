import type { RabbitMQService } from "./rabbitmq-service";

const AMQ_DIRECT_EXCHANGE = "amq.direct";
const SERVER_DATETIME_EVENT = "server_datetime_changed";

export class ServerTimeService {
    private readonly rabbitmqService: RabbitMQService;

    constructor(rabbitmqService: RabbitMQService) {
        this.rabbitmqService = rabbitmqService;
    }

    async setServerDateTime(worldId: number, datetime: string, reset: boolean): Promise<boolean> {
        if (!Number.isFinite(worldId) || worldId < 0) {
            return false;
        }
        await this.rabbitmqService.assertDirectExchange(AMQ_DIRECT_EXCHANGE);
        const routingKey = `fm.${worldId}.all.global`;
        await this.rabbitmqService.publish(AMQ_DIRECT_EXCHANGE, routingKey, SERVER_DATETIME_EVENT, {
            event_id: `${Date.now()}-${Math.random().toString(36).slice(2, 8)}`,
            world_id: worldId,
            reset,
            datetime: reset ? "" : datetime,
            occurred_at: new Date().toISOString(),
        });
        return true;
    }
}
