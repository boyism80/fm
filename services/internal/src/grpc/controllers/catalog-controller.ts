import type Redis from "ioredis";
import type { AppConfiguration } from "@fm/shared/config/app-configuration";
import type { InternalContext } from "@fm/shared/context/internal-context";
import type {
    FindCashShopReply,
    FindCashShopRequest,
    GetGameChannelStatusReply,
    GetGameChannelStatusRequest,
    GetServerCatalogReply,
    GetServerCatalogRequest,
    PingReply,
    PingRequest,
} from "@fm/shared/protobuf/generated/fminternal/internal_service";
import { ServerRole } from "@fm/shared/protobuf/generated/fminternal/internal_service";
import { redisAliveKey } from "@fm/shared/redis-alive-key";
import type { InternalConfig } from "@fm/shared/types/internal-config";
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
    private readonly grpcError: GrpcErrorHandler;

    constructor(
        internalConfig: InternalConfig,
        internalContext: InternalContext,
        appConfiguration: AppConfiguration,
        grpcError: GrpcErrorHandler,
    ) {
        this.internalContext = internalContext;
        this.appConfiguration = appConfiguration;
        this.internalConfig = internalConfig;
        this.grpcError = grpcError;
    }

    private findChannel(worldId: number, channelId: number): ChannelConfig | undefined {
        const world = this.internalConfig.game_servers?.worlds?.[String(worldId)];
        const list = world?.channels ?? [];
        return list.find((c) => c.channel_id === channelId);
    }

    private heartbeatIndex(scope: string) {
        return redisAliveKey(`${scope}:index`);
    }

    private heartbeatKey(scope: string, member: string) {
        return redisAliveKey(`${scope}:${member}`);
    }

    // The index lives only while some server keeps pinging, so it disappears once every server is down.
    private async heartbeat(client: Redis, scope: string, member: string, online: number) {
        const ttl = this.appConfiguration.getServerAliveTtlSeconds();
        const now = Date.now();
        const index = this.heartbeatIndex(scope);
        await client
            .multi()
            .set(this.heartbeatKey(scope, member), String(online), "EX", ttl)
            .zadd(index, now, member)
            .zremrangebyscore(index, "-inf", now - ttl * 1000)
            .expire(index, ttl)
            .exec();
    }

    private async liveServers(worldId: number, prefix: string): Promise<string[]> {
        const { client } = this.internalContext.getRedisGlobalAccess(worldId);
        const index = this.heartbeatIndex(`w${worldId}`);
        const expired = Date.now() - this.appConfiguration.getServerAliveTtlSeconds() * 1000;
        await client.zremrangebyscore(index, "-inf", expired);
        const members = await client.zrangebyscore(index, expired, "+inf");
        return members.filter((member) => member.startsWith(prefix));
    }

    private async onlineUsers(worldId: number, member: string): Promise<number | null> {
        const { client } = this.internalContext.getRedisGlobalAccess(worldId);
        const raw = await client.get(this.heartbeatKey(`w${worldId}`, member));
        if (raw == null) {
            return null;
        }
        return Number(raw);
    }

    @Method("ping")
    async ping(call: GrpcCall<PingRequest>, callback: GrpcCallback<PingReply>) {
        const req = call.request;
        try {
            if (req.role === ServerRole.SERVER_ROLE_GAME) {
                await this.heartbeat(this.internalContext.getRedisGlobalAccess(req.worldId).client, `w${req.worldId}`, `game:c${req.channelId}`, req.onlineUsers);
            } else if (req.role === ServerRole.SERVER_ROLE_CASH_SHOP) {
                await this.heartbeat(this.internalContext.getRedisGlobalAccess(req.worldId).client, `w${req.worldId}`, `cashshop:cs${req.cashShopId}`, req.onlineUsers);
            } else if (req.role === ServerRole.SERVER_ROLE_LOGIN) {
                const id = (req.loginInstanceId ?? "").trim();
                if (!id) {
                    this.grpcError({ code: "INVALID_PAYLOAD", message: "login ping requires login_instance_id" }, callback);
                    return;
                }
                await this.heartbeat(this.internalContext.getRedisUnifiedAccess().client, "login", id, 0);
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
            const online = await this.onlineUsers(req.worldId, `game:c${req.channelId}`);
            const maxConcurrentUsers = row.max_concurrent_users ?? 0;
            const channelFull = online != null && maxConcurrentUsers > 0 && online >= maxConcurrentUsers;
            callback(null, {
                alive: online != null,
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
            let best: { cashShopId: number; host: string; port: number; users: number } | null = null;
            for (const member of await this.liveServers(worldId, "cashshop:cs")) {
                const cs = cashShops.find((c) => member === `cashshop:cs${c.cash_shop_id}`);
                const users = await this.onlineUsers(worldId, member);
                if (cs == null || users == null) {
                    continue;
                }
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
