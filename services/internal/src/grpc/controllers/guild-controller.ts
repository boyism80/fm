import { GuildErrorCode } from "../../protobuf/generated/fminternal/internal_service";
import type {
    AcceptGuildInviteReply,
    AcceptGuildInviteRequest,
    ChangeGuildEmblemReply,
    ChangeGuildEmblemRequest,
    ChangeGuildNoticeReply,
    ChangeGuildNoticeRequest,
    ChangeGuildMemberRankReply,
    ChangeGuildMemberRankRequest,
    ChangeGuildRankTitlesReply,
    ChangeGuildRankTitlesRequest,
    CreateGuildReply,
    CreateGuildRequest,
    DisbandGuildReply,
    DisbandGuildRequest,
    ExpelGuildReply,
    ExpelGuildRequest,
    GetGuildReply,
    GetGuildRequest,
    LeaveGuildReply,
    LeaveGuildRequest,
} from "../../protobuf/generated/fminternal/internal_service";
import type {
    AcceptGuildInviteResult,
    ChangeGuildEmblemResult,
    ChangeGuildNoticeResult,
    ChangeGuildMemberRankResult,
    ChangeGuildRankTitlesResult,
    CreateGuildResult,
    DisbandGuildResult,
    ExpelGuildResult,
    GetGuildResult,
    GuildService,
    LeaveGuildResult,
} from "../../services/guild-service";
import type { GrpcCall, GrpcCallback, GrpcErrorHandler } from "./types";
import { Controller, Method } from "../grpc-method-decorator";

@Controller("guildController")
export class GuildGrpcController {
    private readonly guildService: GuildService;
    private readonly grpcError: GrpcErrorHandler;

    constructor(guildService: GuildService, grpcError: GrpcErrorHandler) {
        this.guildService = guildService;
        this.grpcError = grpcError;
    }

    @Method("createGuild")
    async createGuild(call: GrpcCall<CreateGuildRequest>, callback: GrpcCallback<CreateGuildReply>) {
        try {
            const req = call.request;
            const result = await this.guildService.createGuild(req.worldId, req.guildName, req.leader) as CreateGuildResult;
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
            this.grpcError(err, callback);
        }
    }

    @Method("getGuild")
    async getGuild(call: GrpcCall<GetGuildRequest>, callback: GrpcCallback<GetGuildReply>) {
        try {
            const worldId = call.request.worldId;
            const result = await this.guildService.getGuild(worldId, call.request.guildId) as GetGuildResult;
            if (result.found && result.guild) {
                const guild = await this.guildService.buildGuildMessage(worldId, result.guild, result.members ?? []);
                callback(null, { found: true, guild });
            } else {
                callback(null, { found: false, guild: undefined });
            }
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("acceptGuildInvite")
    async acceptGuildInvite(call: GrpcCall<AcceptGuildInviteRequest>, callback: GrpcCallback<AcceptGuildInviteReply>) {
        try {
            const req = call.request;
            const result = await this.guildService.acceptGuildInvite(req.worldId, req.guildId, req.member) as AcceptGuildInviteResult;
            if (result.ok) {
                callback(null, {
                    ok: true,
                    errorCode: GuildErrorCode.GUILD_ERROR_NONE,
                    guild: result.guild,
                });
            } else {
                callback(null, {
                    ok: false,
                    errorCode: result.code ?? GuildErrorCode.GUILD_ERROR_UNKNOWN,
                    guild: undefined,
                });
            }
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("leaveGuild")
    async leaveGuild(call: GrpcCall<LeaveGuildRequest>, callback: GrpcCallback<LeaveGuildReply>) {
        try {
            const result = await this.guildService.leaveGuild(call.request.worldId, call.request.characterId) as LeaveGuildResult;
            if (result.ok) {
                callback(null, {
                    ok: true,
                    errorCode: GuildErrorCode.GUILD_ERROR_NONE,
                    guildId: result.guildId,
                    revision: result.revision ?? 0,
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
            this.grpcError(err, callback);
        }
    }

    @Method("expelGuild")
    async expelGuild(call: GrpcCall<ExpelGuildRequest>, callback: GrpcCallback<ExpelGuildReply>) {
        try {
            const req = call.request;
            const result = await this.guildService.expelGuild(
                req.worldId,
                req.requesterCharacterId,
                req.targetCharacterId
            ) as ExpelGuildResult;
            if (result.ok) {
                callback(null, {
                    ok: true,
                    errorCode: GuildErrorCode.GUILD_ERROR_NONE,
                    guildId: result.guildId,
                    revision: result.revision ?? 0,
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
            this.grpcError(err, callback);
        }
    }

    @Method("changeGuildRankTitles")
    async changeGuildRankTitles(
        call: GrpcCall<ChangeGuildRankTitlesRequest>,
        callback: GrpcCallback<ChangeGuildRankTitlesReply>
    ) {
        try {
            const req = call.request;
            const result = await this.guildService.changeGuildRankTitles(
                req.worldId,
                req.characterId,
                req.rankTitles
            ) as ChangeGuildRankTitlesResult;
            if (result.ok) {
                callback(null, {
                    ok: true,
                    errorCode: GuildErrorCode.GUILD_ERROR_NONE,
                    guildId: result.guildId,
                    revision: result.revision ?? 0,
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
            this.grpcError(err, callback);
        }
    }

    @Method("changeGuildMemberRank")
    async changeGuildMemberRank(
        call: GrpcCall<ChangeGuildMemberRankRequest>,
        callback: GrpcCallback<ChangeGuildMemberRankReply>
    ) {
        try {
            const req = call.request;
            const result = await this.guildService.changeGuildMemberRank(
                req.worldId,
                req.requesterCharacterId,
                req.targetCharacterId,
                req.newRank
            ) as ChangeGuildMemberRankResult;
            if (result.ok) {
                callback(null, {
                    ok: true,
                    errorCode: GuildErrorCode.GUILD_ERROR_NONE,
                    guildId: result.guildId,
                    revision: result.revision ?? 0,
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
            this.grpcError(err, callback);
        }
    }

    @Method("changeGuildEmblem")
    async changeGuildEmblem(call: GrpcCall<ChangeGuildEmblemRequest>, callback: GrpcCallback<ChangeGuildEmblemReply>) {
        try {
            const req = call.request;
            const result = await this.guildService.changeGuildEmblem(
                req.worldId,
                req.characterId,
                req.logo
            ) as ChangeGuildEmblemResult;
            if (result.ok) {
                callback(null, {
                    ok: true,
                    errorCode: GuildErrorCode.GUILD_ERROR_NONE,
                    guildId: result.guildId,
                    revision: result.revision ?? 0,
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
            this.grpcError(err, callback);
        }
    }

    @Method("changeGuildNotice")
    async changeGuildNotice(call: GrpcCall<ChangeGuildNoticeRequest>, callback: GrpcCallback<ChangeGuildNoticeReply>) {
        try {
            const req = call.request;
            const result = await this.guildService.changeGuildNotice(
                req.worldId,
                req.characterId,
                req.notice
            ) as ChangeGuildNoticeResult;
            if (result.ok) {
                callback(null, {
                    ok: true,
                    errorCode: GuildErrorCode.GUILD_ERROR_NONE,
                    guildId: result.guildId,
                    revision: result.revision ?? 0,
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
            this.grpcError(err, callback);
        }
    }

    @Method("disbandGuild")
    async disbandGuild(call: GrpcCall<DisbandGuildRequest>, callback: GrpcCallback<DisbandGuildReply>) {
        try {
            const result = await this.guildService.disbandGuild(
                call.request.worldId,
                call.request.characterId
            ) as DisbandGuildResult;
            if (result.ok) {
                callback(null, {
                    ok: true,
                    errorCode: GuildErrorCode.GUILD_ERROR_NONE,
                    guildId: result.guildId,
                    revision: result.revision ?? 0,
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
            this.grpcError(err, callback);
        }
    }
}
