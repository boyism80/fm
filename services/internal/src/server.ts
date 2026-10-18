import grpc, { type handleUnaryCall } from "@grpc/grpc-js";
import * as awilix from "awilix";
import { InternalService } from "@fm/shared/protobuf/generated/fminternal/internal_service";
import { createAppContainer } from "./container";
import { autoMigrateAllIfEnabled } from "@fm/shared/auto-migrate";
import { CatalogGrpcController } from "./grpc/controllers/catalog-controller";
import { AuthGrpcController } from "./grpc/controllers/auth-controller";
import { SessionGrpcController } from "./grpc/controllers/session-controller";
import { CharacterGrpcController } from "./grpc/controllers/character-controller";
import { PartyGrpcController } from "./grpc/controllers/party-controller";
import { GuildGrpcController } from "./grpc/controllers/guild-controller";
import { BuddyGrpcController } from "./grpc/controllers/buddy-controller";
import { ChatGrpcController } from "./grpc/controllers/chat-controller";
import { ServerTimeGrpcController } from "./grpc/controllers/server-time-controller";
import { FaultGrpcController } from "./grpc/controllers/fault-controller";
import { PlayerNpcGrpcController } from "./grpc/controllers/player-npc-controller";
import { OperationLogGrpcController } from "./grpc/controllers/operation-log-controller";
import { ParcelGrpcController } from "./grpc/controllers/parcel-controller";
import { ShopGrpcController } from "./grpc/controllers/shop-controller";
import { CashShopGrpcController } from "./grpc/controllers/cash-shop-controller";
import { SocialGrpcController } from "./grpc/controllers/social-controller";
import { MarriageGrpcController } from "./grpc/controllers/marriage-controller";
import { RankingGrpcController } from "./grpc/controllers/ranking-controller";
import { MedalGrpcController } from "./grpc/controllers/medal-controller";
import { getGrpcRoutes } from "./grpc/grpc-method-decorator";
import { DailyLogFile } from "@fm/shared/log-file";
import type { AppConfiguration } from "@fm/shared/config/app-configuration";
import type { InternalContext } from "@fm/shared/context/internal-context";
import type { RabbitMQService } from "./services/rabbitmq-service";
import type { WzService } from "./services/wz-service";
import type { GuildService } from "./services/guild-service";
import type { GuildBulletinBoardService } from "./services/guild-bulletin-board-service";
import type { AllianceService } from "./services/alliance-service";
import type { FamilyService } from "./services/family-service";
import type { CashShopService } from "./services/cash-shop-service";
import type { AccountService } from "./services/account-service";
import type { SocialService } from "./services/social-service";
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
    guildService: GuildService;
    guildBulletinBoardService: GuildBulletinBoardService;
    allianceService: AllianceService;
    familyService: FamilyService;
    cashShopService: CashShopService;
    accountService: AccountService;
    socialService: SocialService;
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
    new DailyLogFile("internal").captureConsole();
    const container = createAppContainer() as awilix.AwilixContainer<ServerContainerCradle>;
    const appConfiguration = container.resolve("appConfiguration");
    const internalConfig = appConfiguration.raw;
    await autoMigrateAllIfEnabled(internalConfig);
    const internalContext = container.resolve("internalContext");
    const rabbitmqService = container.resolve("rabbitmqService");
    const wzService = container.resolve("wzService");

    const worldIds = Object.keys(internalConfig.postgresql.worlds).map(Number);
    const pgUnified = internalConfig.postgresql.unified;
    console.log(
        `fm internal: loaded ${internalConfig.configPath} | grpc ${internalConfig.grpc.host}:${internalConfig.grpc.port} | worlds ${worldIds.join(",")} | ` +
            `character_cache_ttl_s ${appConfiguration.getCharacterCacheTtlSeconds()} | item_cache_ttl_s ${appConfiguration.getItemCacheTtlSeconds()}` +
            (pgUnified ? ` | pg_unified ${pgUnified.host}:${pgUnified.port}/${pgUnified.database}` : "") +
            ` | game_catalog_worlds ${Object.keys(internalConfig.game_servers?.worlds ?? {}).length}` +
            ` | rabbitmq ${internalConfig.rabbitmq?.ip}:${internalConfig.rabbitmq?.port}/${internalConfig.rabbitmq?.vhost}`
    );

    await container.resolve("accountService").moveAccountsToUnified();
    await wzService.preload();
    await rabbitmqService.start();

    for (const worldId of worldIds) {
        try {
            await container.resolve("socialService").pruneMessengers(worldId);
        } catch (err) {
            console.error(`[social] messenger prune failed world=${worldId}`, err);
        }
    }

    const cashShopService = container.resolve("cashShopService");
    let settlingCashTransfers = false;
    const settleCashTransfers = async () => {
        if (settlingCashTransfers) {
            return;
        }
        settlingCashTransfers = true;
        for (const worldId of worldIds) {
            try {
                await cashShopService.settleTransfers(worldId);
            } catch (err) {
                console.error(`[cash-shop] transfer scan failed world=${worldId}`, err);
            }
        }
        settlingCashTransfers = false;
    };
    const cashTransferTimer = setInterval(settleCashTransfers, 15_000);
    cashTransferTimer.unref();
    void settleCashTransfers();

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
        serverTimeController: awilix.asClass(ServerTimeGrpcController).scoped(),
        faultController: awilix.asClass(FaultGrpcController).scoped(),
        playerNpcController: awilix.asClass(PlayerNpcGrpcController).scoped(),
        operationLogController: awilix.asClass(OperationLogGrpcController).scoped(),
        parcelController: awilix.asClass(ParcelGrpcController).scoped(),
        shopController: awilix.asClass(ShopGrpcController).scoped(),
        cashShopController: awilix.asClass(CashShopGrpcController).scoped(),
        marriageController: awilix.asClass(MarriageGrpcController).scoped(),
        socialController: awilix.asClass(SocialGrpcController).scoped(),
        rankingController: awilix.asClass(RankingGrpcController).scoped(),
        medalController: awilix.asClass(MedalGrpcController).scoped(),
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
                await internalContext.withDatabaseOperation(route.grpcMethod, "", async () => {
                    await method.call(controller, call, callback);
                });
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
        clearInterval(cashTransferTimer);
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
