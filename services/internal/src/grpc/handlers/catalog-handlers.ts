import type { AppConfiguration } from "../../config/app-configuration";
import type { InternalContext } from "../../context/internal-context";
import type {
    CheckGameChannelAliveReply,
    CheckGameChannelAliveRequest,
    GetServerCatalogReply,
    GetServerCatalogRequest,
    PingReply,
    PingRequest,
} from "../../protobuf/generated/fminternal/internal_service";
import { ServerRole } from "../../protobuf/generated/fminternal/internal_service";
import { redisAliveKey } from "../../redis-alive-key";
import type { GrpcCall, GrpcCallback, GrpcErrorHandler } from "./types";
import { GrpcController, GrpcMethod } from "../grpc-method-decorator";

type ChannelConfig = { channel_id: number; host: string; port: number; name?: string };
type WorldConfig = { world_name?: string; flag?: number; event_message?: string; channels?: ChannelConfig[] };
type GameServersConfig = { worlds?: Record<string, WorldConfig> };
type InternalConfig = { game_servers?: GameServersConfig };

export function createCatalogHandlers(internalConfig: InternalConfig, grpcError: GrpcErrorHandler) {
    return {
        getServerCatalog(_call: GrpcCall<GetServerCatalogRequest>, callback: GrpcCallback<GetServerCatalogReply>) {
            try {
                const worlds = internalConfig.game_servers?.worlds ?? {};
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
                grpcError(err, callback);
            }
        },
    };
}

@GrpcController("catalogController")
export class CatalogGrpcController {
    private readonly handlers: ReturnType<typeof createCatalogHandlers>;
    private readonly internalContext: InternalContext;
    private readonly appConfiguration: AppConfiguration;
    private readonly grpcError: GrpcErrorHandler;

    constructor(
        internalConfig: InternalConfig,
        internalContext: InternalContext,
        appConfiguration: AppConfiguration,
        grpcError: GrpcErrorHandler,
    ) {
        this.handlers = createCatalogHandlers(internalConfig, grpcError);
        this.internalContext = internalContext;
        this.appConfiguration = appConfiguration;
        this.grpcError = grpcError;
    }

    @GrpcMethod("ping")
    async ping(call: GrpcCall<PingRequest>, callback: GrpcCallback<PingReply>) {
        const req = call.request;
        const ttl = this.appConfiguration.getServerAliveTtlSeconds();
        try {
            if (req.role === ServerRole.SERVER_ROLE_GAME) {
                const { client } = this.internalContext.getRedisGlobalAccess(req.worldId);
                const key = redisAliveKey(`game:w${req.worldId}:c${req.channelId}`);
                await client.set(key, String(Date.now()), "EX", ttl);
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

    @GrpcMethod("checkGameChannelAlive")
    async checkGameChannelAlive(
        call: GrpcCall<CheckGameChannelAliveRequest>,
        callback: GrpcCallback<CheckGameChannelAliveReply>,
    ) {
        const req = call.request;
        try {
            const { client } = this.internalContext.getRedisGlobalAccess(req.worldId);
            const key = redisAliveKey(`game:w${req.worldId}:c${req.channelId}`);
            const v = await client.get(key);
            const alive = v != null && v !== "";
            callback(null, { alive });
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @GrpcMethod("getServerCatalog")
    async getServerCatalog(call: GrpcCall<GetServerCatalogRequest>, callback: GrpcCallback<GetServerCatalogReply>) {
        return this.handlers.getServerCatalog(call, callback);
    }
}
