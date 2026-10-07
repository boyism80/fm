import type { AppConfiguration } from "../../config/app-configuration";
import type { InternalContext } from "../../context/internal-context";
import type {
    FindCashShopReply,
    FindCashShopRequest,
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

    private findChannel(worldId: number, channelId: number): ChannelConfig | undefined {
        const world = this.internalConfig.game_servers?.worlds?.[String(worldId)];
        const list = world?.channels ?? [];
        return list.find((c) => c.channel_id === channelId);
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
            } else if (req.role === ServerRole.SERVER_ROLE_CASH_SHOP) {
                const { client } = this.internalContext.getRedisGlobalAccess(req.worldId);
                const key = redisAliveKey(`cashshop:w${req.worldId}:cs${req.cashShopId}`);
                await client.set(key, String(Date.now()), "EX", ttl);
                await this.sessionRepository.touchCashShopOnlineUsersTtl(req.worldId, req.cashShopId);
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
                this.grpcError({ code: "INVALID_PAYLOAD", message: "ping requires role LOGIN, GAME or CASH_SHOP" }, callback);
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
            const row = this.findChannel(req.worldId, req.channelId);
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
            const maxConcurrentUsers = row?.max_concurrent_users ?? 0;
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

    @Method("findCashShop")
    async findCashShop(call: GrpcCall<FindCashShopRequest>, callback: GrpcCallback<FindCashShopReply>) {
        const worldId = call.request.worldId;
        try {
            const cashShops = this.internalConfig.game_servers?.worlds?.[String(worldId)]?.cash_shops ?? [];
            const { client } = this.internalContext.getRedisGlobalAccess(worldId);
            let best: { cashShopId: number; host: string; port: number; users: number } | null = null;
            for (const cs of cashShops) {
                const alive = await client.get(redisAliveKey(`cashshop:w${worldId}:cs${cs.cash_shop_id}`));
                if (alive == null || alive === "") {
                    continue;
                }
                const users = await this.sessionRepository.getCashShopOnlineUserCount(worldId, cs.cash_shop_id);
                if (best == null || users < best.users) {
                    best = { cashShopId: cs.cash_shop_id, host: cs.host, port: cs.port, users };
                }
            }
            if (best == null) {
                callback(null, { found: false, cashShopId: 0, host: "", port: 0 });
                return;
            }
            callback(null, { found: true, cashShopId: best.cashShopId, host: best.host, port: best.port });
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
