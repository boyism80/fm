"use strict";

const grpc = require("@grpc/grpc-js");
const { InternalService } = require("../protobuf/fminternal/ping_grpc_pb.js");
const messages = require("../protobuf/fminternal/ping_pb.js");
const { loadConfig } = require("./config");

function main() {
    const cfg = loadConfig();
    console.log(
        `fm internal: loaded ${cfg.configPath} | grpc ${cfg.grpc.host}:${cfg.grpc.port} | world ${cfg.app.world_id} | pg ${cfg.postgresql.host}:${cfg.postgresql.port}/${cfg.postgresql.database} | redis ${cfg.redis.host}:${cfg.redis.port}`
    );
    
    const server = new grpc.Server();
    
    server.addService(InternalService, {
        ping(_call, callback) {
            const reply = new messages.PingReply();
            reply.setMessage("pong");
            callback(null, reply);
        },
    });
    
    const addr = `${cfg.grpc.host}:${cfg.grpc.port}`;
    
    server.bindAsync(addr, grpc.ServerCredentials.createInsecure(), (err, boundPort) => {
        if (err) {
            console.error(err);
            process.exit(1);
        }
        server.start();
        console.log(`fm internal gRPC listening on ${addr} (bound port ${boundPort})`);
    });
}

main();
