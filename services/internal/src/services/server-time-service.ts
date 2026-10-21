import type { RabbitMQService } from "./rabbitmq-service";

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
        const routingKey = `fm.${worldId}.all.global`;
        await this.rabbitmqService.publish(routingKey, SERVER_DATETIME_EVENT, {
            world_id: worldId,
            reset,
            datetime: reset ? "" : datetime,
        });
        return true;
    }

    async broadcastNotice(
        worldId: number,
        sourceChannelId: number,
        messageType: number,
        message: string,
        megaEar: boolean,
        nextLines: string[],
        item: Buffer,
        adminsOnly: boolean,
    ): Promise<boolean> {
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
        if (nextLines.length > 2 || nextLines.some((line) => line.length <= 0 || line.length > 200)) {
            return false;
        }
        const routingKey = `fm.${worldId}.all.global`;
        await this.rabbitmqService.publish(routingKey, NOTICE_EVENT, {
            world_id: worldId,
            source_channel_id: sourceChannelId,
            message_type: messageType,
            message: text,
            mega_ear: megaEar,
            next_lines: nextLines,
            item: item.toString("base64"),
            admins_only: adminsOnly,
        });
        return true;
    }
}
