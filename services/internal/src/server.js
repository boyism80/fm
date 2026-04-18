"use strict";

const grpc = require("@grpc/grpc-js");
const awilix = require("awilix");
const { InternalService } = require("../protobuf/fminternal/internal_service_grpc_pb");
const messages = require("../protobuf/fminternal/internal_service_pb");
const { createAppContainer } = require("./container");
const { autoMigrateAllIfEnabled } = require("./auto-migrate");
const { createCatalogHandlers } = require("./grpc/handlers/catalog-handlers");
const { createAuthHandlers } = require("./grpc/handlers/auth-handlers");
const { createSessionHandlers } = require("./grpc/handlers/session-handlers");
const { createCharacterHandlers } = require("./grpc/handlers/character-handlers");
const { createPartyHandlers } = require("./grpc/handlers/party-handlers");

const INVALID_CODES = new Set([
    "UNKNOWN_WORLD", "INVALID_CHARACTER_ID",
    "INVALID_PAYLOAD", "UNKNOWN_CHANNEL",
]);

function grpcError(err, callback) {
    const code = err && INVALID_CODES.has(err.code)
        ? grpc.status.INVALID_ARGUMENT
        : grpc.status.INTERNAL;
    callback({ code, message: err.message || String(err) });
}

async function main() {
    const container = createAppContainer();
    const appConfiguration = container.resolve("appConfiguration");
    const internalConfig = appConfiguration.raw;
    await autoMigrateAllIfEnabled(internalConfig);
    const internalContext = container.resolve("internalContext");
    const rabbitmqService = container.resolve("rabbitmqService");
    const partyEventPublisher = container.resolve("partyEventPublisher");

    const wid = String(internalConfig.app.world_id);
    const pgw = internalConfig.postgresql.worlds[wid];
    const rgw = internalConfig.redis.worlds[wid];
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

    await rabbitmqService.start();
    await partyEventPublisher.initialize();

    const server = new grpc.Server();

    container.register({
        messages: awilix.asValue(messages),
        grpcError: awilix.asValue(grpcError),
        catalogHandlers: awilix.asFunction(createCatalogHandlers).singleton(),
        authHandlers: awilix.asFunction(createAuthHandlers).singleton(),
        sessionHandlers: awilix.asFunction(createSessionHandlers).singleton(),
        characterHandlers: awilix.asFunction(createCharacterHandlers).singleton(),
        partyHandlers: awilix.asFunction(createPartyHandlers).singleton(),
    });

    const catalogHandlers = container.resolve("catalogHandlers");
    const authHandlers = container.resolve("authHandlers");
    const sessionHandlers = container.resolve("sessionHandlers");
    const characterHandlers = container.resolve("characterHandlers");
    const partyHandlers = container.resolve("partyHandlers");

    server.addService(InternalService, {
        ...catalogHandlers,
        ...authHandlers,
        ...sessionHandlers,
        ...characterHandlers,
        ...partyHandlers,
    });

    const addr = `${internalConfig.grpc.host}:${internalConfig.grpc.port}`;

    server.bindAsync(addr, grpc.ServerCredentials.createInsecure(), (err, boundPort) => {
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

main().catch((err) => {
    console.error("[startup] fatal:", err);
    process.exit(1);
});
