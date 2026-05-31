import type { AppConfiguration } from "../../config/app-configuration";
import type { InternalContext } from "../../context/internal-context";
import type {
    GetGameChannelStatusReply,
    GetGameChannelStatusRequest,
    GetServerCatalogReply,
    GetServerCatalogRequest,
    PingReply,
    PingRequest,
} from "../../protobuf/generated/fminternal/internal_service";
import { ServerRole } from "../../protobuf/generated/fminternal/internal_service";
import { redisAliveKey } from "../../redis-alive-key";
import type { SessionRepository } from "../../repos/session-repository";
import type { InternalConfig } from "../../types/internal-config";
import type { GrpcCall, GrpcCallback, GrpcErrorHandler } from "./types";
import { Controller, Method } from "../grpc-method-decorator";

type ChannelConfig = {
    channel_id: number;
    host: string;
    port: number;
    name?: string;
    max_concurrent_users?: number;
};
type WorldConfig = { world_name?: string; flag?: number; event_message?: string; channels?: ChannelConfig[] };

@Controller("catalogController")
export class CatalogGrpcController {
    private readonly internalContext: InternalContext;
    private readonly appConfiguration: AppConfiguration;
    private readonly internalConfig: InternalConfig;
    private readonly sessionRepository: SessionRepository;
    private readonly grpcError: GrpcErrorHandler;

    constructor(
        internalConfig: InternalConfig,
        internalContext: InternalContext,
        appConfiguration: AppConfiguration,
        sessionRepository: SessionRepository,
        grpcError: GrpcErrorHandler,
    ) {
        this.internalContext = internalContext;
        this.appConfiguration = appConfiguration;
        this.internalConfig = internalConfig;
        this.sessionRepository = sessionRepository;
        this.grpcError = grpcError;
    }

    private resolveChannel(worldId: number, channelId: number): ChannelConfig | undefined {
        const world = this.internalConfig.game_servers?.worlds?.[String(worldId)];
        const list = world?.channels ?? [];
        return list.find((c) => c.channel_id === channelId);
    }

    private resolveChannelMaxConcurrentUsers(worldId: number, channelId: number): number {
        const row = this.resolveChannel(worldId, channelId);
        if (!row) {
            return 0;
        }
        return row.max_concurrent_users ?? 0;
    }

    @Method("ping")
    async ping(call: GrpcCall<PingRequest>, callback: GrpcCallback<PingReply>) {
        const req = call.request;
        const ttl = this.appConfiguration.getServerAliveTtlSeconds();
        try {
            if (req.role === ServerRole.SERVER_ROLE_GAME) {
                const { client } = this.internalContext.getRedisGlobalAccess(req.worldId);
                const key = redisAliveKey(`game:w${req.worldId}:c${req.channelId}`);
                await client.set(key, String(Date.now()), "EX", ttl);
                await this.sessionRepository.touchChannelOnlineUsersTtl(req.worldId, req.channelId);
            } else if (req.role === ServerRole.SERVER_ROLE_LOGIN) {
                const id = (req.loginInstanceId ?? "").trim();
                if (!id) {
                    this.grpcError({ code: "INVALID_PAYLOAD", message: "login ping requires login_instance_id" }, callback);
                    return;
                }
                const redisWorld = this.appConfiguration.raw.app.world_id;
                const { client } = this.internalContext.getRedisGlobalAccess(redisWorld);
                const key = redisAliveKey(`login:${id}`);
                await client.set(key, String(Date.now()), "EX", ttl);
            } else {
                this.grpcError({ code: "INVALID_PAYLOAD", message: "ping requires role LOGIN or GAME" }, callback);
                return;
            }
            callback(null, { message: "pong" });
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("getGameChannelStatus")
    async getGameChannelStatus(
        call: GrpcCall<GetGameChannelStatusRequest>,
        callback: GrpcCallback<GetGameChannelStatusReply>,
    ) {
        const req = call.request;
        try {
            const row = this.resolveChannel(req.worldId, req.channelId);
            if (!row) {
                callback(null, {
                    alive: false,
                    channelFull: false,
                    found: false,
                    host: "",
                    port: 0,
                });
                return;
            }
            const { client } = this.internalContext.getRedisGlobalAccess(req.worldId);
            const key = redisAliveKey(`game:w${req.worldId}:c${req.channelId}`);
            const v = await client.get(key);
            const alive = v != null && v !== "";
            const onlineUserCount = await this.sessionRepository.getChannelOnlineUserCount(req.worldId, req.channelId);
            const maxConcurrentUsers = this.resolveChannelMaxConcurrentUsers(req.worldId, req.channelId);
            const channelFull = maxConcurrentUsers > 0 && onlineUserCount >= maxConcurrentUsers;
            callback(null, {
                alive,
                channelFull,
                found: true,
                host: row.host,
                port: row.port,
            });
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("getServerCatalog")
    getServerCatalog(_call: GrpcCall<GetServerCatalogRequest>, callback: GrpcCallback<GetServerCatalogReply>) {
        try {
            const worlds = this.internalConfig.game_servers?.worlds ?? {};
            const worldMsgs = [];
            for (const [worldId, world] of Object.entries(worlds)) {
                const channels = Array.isArray(world.channels) ? world.channels : [];
                worldMsgs.push({
                    worldId: Number(worldId),
                    worldName: world.world_name || `World-${worldId}`,
                    flag: world.flag ?? 0,
                    eventMessage: world.event_message || "",
                    channels: channels.map((ch) => ({
                        channelId: ch.channel_id,
                        host: ch.host,
                        port: ch.port,
                        name: ch.name || `Channel ${ch.channel_id + 1}`,
                    })),
                });
            }
            callback(null, { worlds: worldMsgs });
        } catch (err) {
            this.grpcError(err, callback);
        }
    }
}
