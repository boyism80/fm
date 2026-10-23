import amqp from "amqplib";
import type { AppConfiguration } from "@fm/shared/config/app-configuration";
import type { InternalConfig } from "@fm/shared/types/internal-config";

const RECONNECT_MIN_DELAY_MS = 1000;
const RECONNECT_MAX_DELAY_MS = 30000;
const DIRECT_EXCHANGE = "amq.direct";

export class RabbitMQService {
    private readonly cfg: InternalConfig["rabbitmq"];
    private connection: amqp.ChannelModel | null;
    private channel: amqp.Channel | null;
    private started: boolean;
    private reconnectTimer: NodeJS.Timeout | null;
    private reconnectDelayMs: number;
    private exchangeAsserted: boolean;
    private publishFault: boolean;

    constructor(appConfiguration: AppConfiguration) {
        this.cfg = appConfiguration.raw.rabbitmq;
        this.connection = null;
        this.channel = null;
        this.started = false;
        this.reconnectTimer = null;
        this.reconnectDelayMs = RECONNECT_MIN_DELAY_MS;
        this.exchangeAsserted = false;
        this.publishFault = false;
    }

    setPublishFault(enabled: boolean) {
        this.publishFault = enabled;
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
        this.started = true;
        try {
            await this.connect();
        } catch (error) {
            this.started = false;
            throw error;
        }
    }

    private async connect() {
        const connection = await amqp.connect(this.amqpUrl());
        connection.on("error", (error) => {
            console.error("[rabbitmq] connection error", error);
        });
        connection.on("close", () => {
            this.disconnected(connection);
        });

        let channel: amqp.Channel;
        try {
            channel = await connection.createChannel();
        } catch (error) {
            await connection.close().catch(() => {});
            throw error;
        }
        channel.on("error", (error) => {
            console.error("[rabbitmq] channel error", error);
        });
        channel.on("close", () => {
            this.disconnected(connection);
        });

        if (this.started === false) {
            await connection.close().catch(() => {});
            return;
        }
        this.connection = connection;
        this.channel = channel;
        this.reconnectDelayMs = RECONNECT_MIN_DELAY_MS;
        this.exchangeAsserted = false;
    }

    private disconnected(connection: amqp.ChannelModel) {
        if (this.connection !== connection) {
            return;
        }
        this.connection = null;
        this.channel = null;
        this.exchangeAsserted = false;
        connection.close().catch(() => {});
        if (this.started === false) {
            return;
        }
        console.error("[rabbitmq] disconnected, reconnecting");
        this.scheduleReconnect();
    }

    private scheduleReconnect() {
        if (this.reconnectTimer !== null) {
            return;
        }
        const delay = this.reconnectDelayMs;
        this.reconnectDelayMs = Math.min(delay * 2, RECONNECT_MAX_DELAY_MS);
        this.reconnectTimer = setTimeout(() => {
            this.reconnectTimer = null;
            if (this.started === false) {
                return;
            }
            this.connect().catch((error) => {
                console.error("[rabbitmq] reconnect failed", error);
                this.scheduleReconnect();
            });
        }, delay);
    }

    async publish(routingKey: string, eventType: string, payload: Record<string, unknown> = {}) {
        if (!this.started || !this.channel || this.publishFault) {
            throw new Error("rabbitmq service is not started");
        }
        if (typeof eventType !== "string" || eventType.length === 0) {
            throw new Error("rabbitmq publish: eventType must be a non-empty string");
        }
        if (!this.exchangeAsserted) {
            await this.channel.assertExchange(DIRECT_EXCHANGE, "direct", { durable: true });
            this.exchangeAsserted = true;
        }
        const bodyObj = { ...payload, event_type: eventType };
        const body = Buffer.from(JSON.stringify(bodyObj));
        return this.channel.publish(DIRECT_EXCHANGE, routingKey, body, {
            contentType: "application/json",
            persistent: true,
        });
    }

    async close() {
        this.started = false;
        if (this.reconnectTimer !== null) {
            clearTimeout(this.reconnectTimer);
            this.reconnectTimer = null;
        }
        this.exchangeAsserted = false;
        if (this.channel) {
            try {
                await this.channel.close();
            } catch {}
            this.channel = null;
        }
        if (this.connection) {
            try {
                await this.connection.close();
            } catch {}
            this.connection = null;
        }
    }
}
