import type {
    GetServerCatalogReply,
    GetServerCatalogRequest,
    PingReply,
    PingRequest,
} from "../../protobuf/generated/fminternal/internal_service";
import type { GrpcCall, GrpcCallback, GrpcErrorHandler } from "./types";
import { GrpcController, GrpcMethod } from "../grpc-method-decorator";

type ChannelConfig = { channel_id: number; host: string; port: number; name?: string };
type WorldConfig = { world_name?: string; flag?: number; event_message?: string; channels?: ChannelConfig[] };
type GameServersConfig = { worlds?: Record<string, WorldConfig> };
type InternalConfig = { game_servers?: GameServersConfig };

export function createCatalogHandlers(internalConfig: InternalConfig, grpcError: GrpcErrorHandler) {
    return {
        ping(_call: GrpcCall<PingRequest>, callback: GrpcCallback<PingReply>) {
            try {
                callback(null, { message: "pong" });
            } catch (err) {
                grpcError(err, callback);
            }
        },
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

    constructor(internalConfig: InternalConfig, grpcError: GrpcErrorHandler) {
        this.handlers = createCatalogHandlers(internalConfig, grpcError);
    }

    @GrpcMethod("ping")
    async ping(call: GrpcCall<PingRequest>, callback: GrpcCallback<PingReply>) {
        return this.handlers.ping(call, callback);
    }

    @GrpcMethod("getServerCatalog")
    async getServerCatalog(call: GrpcCall<GetServerCatalogRequest>, callback: GrpcCallback<GetServerCatalogReply>) {
        return this.handlers.getServerCatalog(call, callback);
    }
}
