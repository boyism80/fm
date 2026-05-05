import amqp from "amqplib";
import type { AppConfiguration } from "../config/app-configuration";
import type { InternalConfig } from "../types/internal-config";

export class RabbitMQService {
    private readonly cfg: InternalConfig["rabbitmq"];
    private connection: amqp.ChannelModel | null;
    private channel: amqp.Channel | null;
    private started: boolean;

    constructor(appConfiguration: AppConfiguration) {
        this.cfg = appConfiguration.raw.rabbitmq;
        this.connection = null;
        this.channel = null;
        this.started = false;
    }

    private amqpUrl() {
        const user = encodeURIComponent(this.cfg.uid);
        const pass = encodeURIComponent(this.cfg.pwd);
        const vhost = encodeURIComponent(this.cfg.vhost || "/");
        return `amqp://${user}:${pass}@${this.cfg.ip}:${this.cfg.port}/${vhost}`;
    }

    async start() {
        if (this.started) {
            return;
        }
        this.connection = await amqp.connect(this.amqpUrl());
        this.channel = await this.connection.createChannel();
        this.started = true;
    }

    async assertTopicExchange(exchange: string) {
        if (!this.started || !this.channel) {
            throw new Error("rabbitmq service is not started");
        }
        await this.channel.assertExchange(exchange, "topic", { durable: true });
    }

    async assertDirectExchange(exchange: string) {
        if (!this.started || !this.channel) {
            throw new Error("rabbitmq service is not started");
        }
        await this.channel.assertExchange(exchange, "direct", { durable: true });
    }

    async publish(exchange: string, routingKey: string, eventType: string, payload: Record<string, unknown> = {}) {
        if (!this.started || !this.channel) {
            throw new Error("rabbitmq service is not started");
        }
        if (typeof eventType !== "string" || eventType.trim() === "") {
            throw new Error("rabbitmq publish: eventType is required (non-empty string)");
        }
        const bodyObj = { ...payload, event_type: eventType };
        const body = Buffer.from(JSON.stringify(bodyObj));
        return this.channel.publish(exchange, routingKey, body, {
            contentType: "application/json",
            persistent: true,
        });
    }

    async close() {
        this.started = false;
        if (this.channel) {
            try {
                await this.channel.close();
            } catch {
            }
            this.channel = null;
        }
        if (this.connection) {
            try {
                await this.connection.close();
            } catch {
            }
            this.connection = null;
        }
    }
}
