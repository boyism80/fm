import type { RabbitMQService } from "./rabbitmq-service";

const AMQ_DIRECT_EXCHANGE = "amq.direct";
const SERVER_DATETIME_EVENT = "server_datetime_changed";
const NOTICE_EVENT = "notice";

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

    async broadcastNotice(worldId: number, sourceChannelId: number, messageType: number, message: string, megaEar: boolean): Promise<boolean> {
        if (!Number.isFinite(worldId) || worldId < 0) {
            return false;
        }
        if (!Number.isFinite(sourceChannelId) || sourceChannelId < 0) {
            return false;
        }
        if (!Number.isInteger(messageType) || messageType < 0 || messageType > 255) {
            return false;
        }
        const text = message;
        if (text.length <= 0 || text.length > 200) {
            return false;
        }
        await this.rabbitmqService.assertDirectExchange(AMQ_DIRECT_EXCHANGE);
        const routingKey = `fm.${worldId}.all.global`;
        await this.rabbitmqService.publish(AMQ_DIRECT_EXCHANGE, routingKey, NOTICE_EVENT, {
            event_id: `${Date.now()}-${Math.random().toString(36).slice(2, 8)}`,
            world_id: worldId,
            source_channel_id: sourceChannelId,
            message_type: messageType,
            message: text,
            mega_ear: megaEar,
            occurred_at: new Date().toISOString(),
        });
        return true;
    }
}
