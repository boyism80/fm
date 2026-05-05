import type { InternalConfig } from "../types/internal-config";

export class AppConfiguration {
    private readonly config: InternalConfig;

    constructor(internalConfig: InternalConfig) {
        if (!internalConfig || typeof internalConfig !== "object") {
            throw new Error("AppConfiguration requires internalConfig from loadConfig()");
        }
        this.config = internalConfig;
    }

    get raw(): InternalConfig {
        return this.config;
    }

    get configPath(): string {
        return this.config.configPath;
    }

    get app() {
        return this.config.app;
    }

    get grpc() {
        return this.config.grpc;
    }

    get postgresql() {
        return this.config.postgresql;
    }

    get redis() {
        return this.config.redis;
    }

    get rabbitmq() {
        return this.config.rabbitmq;
    }

    get cache() {
        return this.config.cache;
    }

    get sequelize() {
        return this.config.sequelize;
    }

    get resources() {
        return this.config.resources;
    }

    getCharacterCacheTtlSeconds(): number {
        const n = this.config.cache?.character_ttl_seconds;
        return Number.isFinite(n) && n > 0 ? n : 300;
    }

    getItemCacheTtlSeconds(): number {
        const n = this.config.cache?.item_ttl_seconds;
        return Number.isFinite(n) && n > 0 ? n : 300;
    }

    getSection(
        section: "app" | "grpc" | "postgresql" | "redis" | "rabbitmq" | "cache" | "sequelize" | "resources"
    ) {
        return this.config[section];
    }
}
