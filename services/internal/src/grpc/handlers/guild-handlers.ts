import { GuildErrorCode } from "../../protobuf/generated/fminternal/internal_service";
import type {
    CreateGuildReply,
    CreateGuildRequest,
    GetGuildReply,
    GetGuildRequest,
} from "../../protobuf/generated/fminternal/internal_service";
import type { CreateGuildResult, GetGuildResult, GuildService } from "../../services/guild-service";
import type { GrpcCall, GrpcCallback, GrpcErrorHandler } from "./types";
import { GrpcController, GrpcMethod } from "../grpc-method-decorator";

export function createGuildHandlers(guildService: GuildService, grpcError: GrpcErrorHandler) {
    return {
        async createGuild(call: GrpcCall<CreateGuildRequest>, callback: GrpcCallback<CreateGuildReply>) {
            try {
                const req = call.request;
                const result = await guildService.createGuild(req.worldId, req.guildName, req.leader) as CreateGuildResult;
                if (result.ok) {
                    callback(null, {
                        ok: true,
                        errorCode: GuildErrorCode.GUILD_ERROR_NONE,
                        guildId: result.guildId,
                        revision: result.revision ?? 0,
                        guild: result.guild,
                    });
                } else {
                    callback(null, {
                        ok: false,
                        errorCode: result.code ?? GuildErrorCode.GUILD_ERROR_UNKNOWN,
                        guildId: undefined,
                        revision: 0,
                    });
                }
            } catch (err) {
                grpcError(err, callback);
            }
        },
        async getGuild(call: GrpcCall<GetGuildRequest>, callback: GrpcCallback<GetGuildReply>) {
            try {
                const worldId = call.request.worldId;
                const result = await guildService.getGuild(worldId, call.request.guildId) as GetGuildResult;
                if (result.found && result.guild) {
                    const guild = await guildService.buildGuildMessage(worldId, result.guild, result.members ?? []);
                    callback(null, { found: true, guild });
                } else {
                    callback(null, { found: false, guild: undefined });
                }
            } catch (err) {
                grpcError(err, callback);
            }
        },
    };
}

@GrpcController("guildController")
export class GuildGrpcController {
    private readonly handlers: ReturnType<typeof createGuildHandlers>;

    constructor(guildService: GuildService, grpcError: GrpcErrorHandler) {
        this.handlers = createGuildHandlers(guildService, grpcError);
    }

    @GrpcMethod("createGuild")
    async createGuild(call: GrpcCall<CreateGuildRequest>, callback: GrpcCallback<CreateGuildReply>) {
        return this.handlers.createGuild(call, callback);
    }

    @GrpcMethod("getGuild")
    async getGuild(call: GrpcCall<GetGuildRequest>, callback: GrpcCallback<GetGuildReply>) {
        return this.handlers.getGuild(call, callback);
    }
}
