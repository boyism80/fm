import grpc, { type handleUnaryCall } from "@grpc/grpc-js";
import * as awilix from "awilix";
import { InternalService } from "./protobuf/generated/fminternal/internal_service";
import { createAppContainer } from "./container";
import { autoMigrateAllIfEnabled } from "./auto-migrate";
import { CatalogGrpcController } from "./grpc/handlers/catalog-handlers";
import { AuthGrpcController } from "./grpc/handlers/auth-handlers";
import { SessionGrpcController } from "./grpc/handlers/session-handlers";
import { CharacterGrpcController } from "./grpc/handlers/character-handlers";
import { PartyGrpcController } from "./grpc/handlers/party-handlers";
import { GuildGrpcController } from "./grpc/handlers/guild-handlers";
import { BuddyGrpcController } from "./grpc/handlers/buddy-handlers";
import { ChatGrpcController } from "./grpc/handlers/chat-handlers";
import { getGrpcRoutes } from "./grpc/grpc-method-decorator";
import type { AppConfiguration } from "./config/app-configuration";
import type { InternalContext } from "./context/internal-context";
import type { RabbitMQService } from "./services/rabbitmq-service";
import type { WzService } from "./services/wz-service";
type ServerContainerCradle = {
    catalogController: CatalogGrpcController;
    authController: AuthGrpcController;
    sessionController: SessionGrpcController;
    characterController: CharacterGrpcController;
    partyController: PartyGrpcController;
    guildController: GuildGrpcController;
    buddyController: BuddyGrpcController;
    chatController: ChatGrpcController;
    appConfiguration: AppConfiguration;
    internalContext: InternalContext;
    rabbitmqService: RabbitMQService;
    wzService: WzService;
    grpcError: typeof grpcError;
};

const INVALID_CODES = new Set([
    "UNKNOWN_WORLD", "INVALID_CHARACTER_ID",
    "INVALID_PAYLOAD", "UNKNOWN_CHANNEL",
]);

function grpcError(err: { code?: string; message?: string }, callback: (error: { code: number; message: string }) => void) {
    const code = err && err.code && INVALID_CODES.has(err.code)
        ? grpc.status.INVALID_ARGUMENT
        : grpc.status.INTERNAL;
    callback({ code, message: err.message || String(err) });
}

function validateGrpcRouteCoverage(routes: Array<{ grpcMethod: string }>) {
    const expectedMethods = new Set(Object.keys(InternalService));
    const registeredMethods = routes.map((route) => route.grpcMethod);
    const duplicateMethods = registeredMethods.filter((method, index) => registeredMethods.indexOf(method) !== index);
    const uniqueRegisteredMethods = new Set(registeredMethods);

    const missingMethods = [...expectedMethods].filter((method) => !uniqueRegisteredMethods.has(method));
    const unknownMethods = [...uniqueRegisteredMethods].filter((method) => !expectedMethods.has(method));

    if (missingMethods.length > 0 || unknownMethods.length > 0 || duplicateMethods.length > 0) {
        const parts: string[] = [];
        if (missingMethods.length > 0) {
            parts.push(`missing=[${missingMethods.sort().join(", ")}]`);
        }
        if (unknownMethods.length > 0) {
            parts.push(`unknown=[${unknownMethods.sort().join(", ")}]`);
        }
        if (duplicateMethods.length > 0) {
            parts.push(`duplicate=[${[...new Set(duplicateMethods)].sort().join(", ")}]`);
        }
        throw new Error(`gRPC route registration mismatch: ${parts.join(" | ")}`);
    }
}

async function main() {
    const container = createAppContainer() as awilix.AwilixContainer<ServerContainerCradle>;
    const appConfiguration = container.resolve("appConfiguration");
    const internalConfig = appConfiguration.raw;
    await autoMigrateAllIfEnabled(internalConfig);
    const internalContext = container.resolve("internalContext");
    const rabbitmqService = container.resolve("rabbitmqService");
    const wzService = container.resolve("wzService");

    const wid = String(internalConfig.app.world_id);
    const pgw = internalConfig.postgresql.worlds[wid];
    const rgw = internalConfig.redis.worlds[wid];
    if (!pgw || !rgw) {
        throw new Error(`world configuration not found for world_id=${wid}`);
    }
    const pgUnified = internalConfig.postgresql.unified;
    console.log(
        `fm internal: loaded ${internalConfig.configPath} | grpc ${internalConfig.grpc.host}:${internalConfig.grpc.port} | world ${wid} | ` +
            `pg global ${pgw.global.host}:${pgw.global.port}/${pgw.global.database} | pg_data_shards ${pgw.data.length} | ` +
            `redis global ${rgw.global.host}:${rgw.global.port} | redis_data_shards ${rgw.data.length} | ` +
            `character_cache_ttl_s ${appConfiguration.getCharacterCacheTtlSeconds()} | item_cache_ttl_s ${appConfiguration.getItemCacheTtlSeconds()}` +
            (pgUnified ? ` | pg_unified ${pgUnified.host}:${pgUnified.port}/${pgUnified.database}` : "") +
            ` | game_catalog_worlds ${Object.keys(internalConfig.game_servers?.worlds ?? {}).length}` +
            ` | rabbitmq ${internalConfig.rabbitmq?.ip}:${internalConfig.rabbitmq?.port}/${internalConfig.rabbitmq?.vhost}`
    );

    await wzService.preload();
    await rabbitmqService.start();

    const server = new grpc.Server();

    container.register({
        grpcError: awilix.asValue(grpcError),
        catalogController: awilix.asClass(CatalogGrpcController).scoped(),
        authController: awilix.asClass(AuthGrpcController).scoped(),
        sessionController: awilix.asClass(SessionGrpcController).scoped(),
        characterController: awilix.asClass(CharacterGrpcController).scoped(),
        partyController: awilix.asClass(PartyGrpcController).scoped(),
        guildController: awilix.asClass(GuildGrpcController).scoped(),
        buddyController: awilix.asClass(BuddyGrpcController).scoped(),
        chatController: awilix.asClass(ChatGrpcController).scoped(),
    });

    const serviceImplementation: Record<string, handleUnaryCall<unknown, unknown>> = {};
    const routes = getGrpcRoutes();
    validateGrpcRouteCoverage(routes);
    for (const route of routes) {
        serviceImplementation[route.grpcMethod] = async (call, callback) => {
            const scope = container.createScope();
            try {
                const controller = scope.resolve(route.resolverName) as Record<string, handleUnaryCall<unknown, unknown>>;
                const method = controller[route.methodName];
                if (typeof method !== "function") {
                    throw new Error(`gRPC method not found: ${route.resolverName}.${route.methodName}`);
                }
                await method.call(controller, call, callback);
            } finally {
                await scope.dispose();
            }
        };
    }

    server.addService(InternalService, serviceImplementation);

    const addr = `${internalConfig.grpc.host}:${internalConfig.grpc.port}`;

    server.bindAsync(addr, grpc.ServerCredentials.createInsecure(), (err: Error | null, boundPort: number) => {
        if (err) {
            console.error(err);
            process.exit(1);
        }
        server.start();
        console.log(`fm internal gRPC listening on ${addr} (bound port ${boundPort})`);
    });

    async function shutdown() {
        server.tryShutdown(async () => {
            await rabbitmqService.close().catch(() => {});
            await internalContext.close().catch(() => {});
            process.exit(0);
        });
        setTimeout(() => {
            rabbitmqService.close().catch(() => {});
            internalContext.close().catch(() => {});
            process.exit(1);
        }, 10_000).unref();
    }

    process.on("SIGINT", shutdown);
    process.on("SIGTERM", shutdown);
}

main().catch((err: unknown) => {
    console.error("[startup] fatal:", err);
    process.exit(1);
});
