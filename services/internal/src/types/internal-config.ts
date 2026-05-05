export interface PgPoolConfig {
    min: number;
    max: number;
}

export interface PgEndpoint {
    host: string;
    port: number;
    database: string;
    username: string;
    password: string;
    ssl: boolean;
    pool: PgPoolConfig;
}

export interface RedisEndpoint {
    host: string;
    port: number;
    password: string;
    db: number;
    tls: boolean;
}

export interface WorldPgConfig {
    global: PgEndpoint;
    data: PgEndpoint[];
}

export interface WorldRedisConfig {
    global: RedisEndpoint;
    data: RedisEndpoint[];
}

export type SequelizeDefineValue = string | number | boolean | null;
export type SequelizeDefineConfig = Record<string, SequelizeDefineValue>;

export interface InternalConfig {
    configPath: string;
    app: {
        log_level: string;
        world_id: number;
    };
    grpc: {
        host: string;
        port: number;
    };
    postgresql: {
        unified: PgEndpoint | null;
        worlds: Record<string, WorldPgConfig>;
    };
    redis: {
        unified: RedisEndpoint | null;
        worlds: Record<string, WorldRedisConfig>;
    };
    rabbitmq: {
        ip: string;
        port: number;
        uid: string;
        pwd: string;
        vhost: string;
    };
    cache: {
        write_strategy: string;
        character_ttl_seconds: number;
        item_ttl_seconds: number;
    };
    sequelize: {
        dialect: string;
        migration_storage: string;
        auto_migrate_on_startup: boolean;
        define: SequelizeDefineConfig;
    };
    resources: {
        wz_root: string;
    };
    game_servers: GameServersConfig;
}

export interface GameServerChannelConfig {
    channel_id: number;
    host: string;
    port: number;
    name: string;
}

export interface GameServerWorldConfig {
    world_name: string;
    flag: number;
    event_message: string;
    channels: GameServerChannelConfig[];
}

export interface GameServersConfig {
    worlds: Record<string, GameServerWorldConfig>;
}
