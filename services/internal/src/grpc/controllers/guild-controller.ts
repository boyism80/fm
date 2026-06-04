import { AllianceErrorCode, GuildErrorCode } from "../../protobuf/generated/fminternal/internal_service";
import type {
    AcceptGuildInviteReply,
    AcceptGuildInviteRequest,
    ChangeGuildEmblemReply,
    ChangeGuildEmblemRequest,
    ChangeGuildNoticeReply,
    ChangeGuildNoticeRequest,
    IncreaseGuildCapacityReply,
    IncreaseGuildCapacityRequest,
    CreateGuildBulletinBoardReplyReply,
    CreateGuildBulletinBoardReplyRequest,
    CreateGuildBulletinBoardThreadReply,
    CreateGuildBulletinBoardThreadRequest,
    DeleteGuildBulletinBoardReplyReply,
    DeleteGuildBulletinBoardReplyRequest,
    DeleteGuildBulletinBoardThreadReply,
    DeleteGuildBulletinBoardThreadRequest,
    ListGuildBulletinBoardThreadsReply,
    ListGuildBulletinBoardThreadsRequest,
    ShowGuildBulletinBoardThreadReply,
    ShowGuildBulletinBoardThreadRequest,
    UpdateGuildBulletinBoardThreadReply,
    UpdateGuildBulletinBoardThreadRequest,
    ChangeGuildMemberRankReply,
    ChangeGuildMemberRankRequest,
    ChangeGuildRankTitlesReply,
    ChangeGuildRankTitlesRequest,
    AcceptAllianceInviteReply,
    AcceptAllianceInviteRequest,
    IncreaseAllianceCapacityReply,
    IncreaseAllianceCapacityRequest,
    ChangeAllianceRankTitlesReply,
    ChangeAllianceRankTitlesRequest,
    ChangeAllianceMemberRankReply,
    ChangeAllianceMemberRankRequest,
    ChangeAllianceLeaderReply,
    ChangeAllianceLeaderRequest,
    ChangeAllianceNoticeReply,
    ChangeAllianceNoticeRequest,
    CreateAllianceReply,
    CreateAllianceRequest,
    CreateGuildReply,
    CreateGuildRequest,
    DisbandAllianceReply,
    DisbandAllianceRequest,
    LeaveAllianceReply,
    LeaveAllianceRequest,
    ExpelAllianceGuildReply,
    ExpelAllianceGuildRequest,
    GetAllianceReply,
    GetAllianceRequest,
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
    CreateGuildBulletinBoardReplyResult,
    CreateGuildBulletinBoardThreadResult,
    DeleteGuildBulletinBoardReplyResult,
    DeleteGuildBulletinBoardThreadResult,
    ListGuildBulletinBoardThreadsResult,
    ShowGuildBulletinBoardThreadResult,
    UpdateGuildBulletinBoardThreadResult,
    ChangeGuildMemberRankResult,
    ChangeGuildRankTitlesResult,
    AcceptAllianceInviteResult,
    IncreaseAllianceCapacityResult,
    ChangeAllianceRankTitlesResult,
    ChangeAllianceMemberRankResult,
    ChangeAllianceLeaderResult,
    ChangeAllianceNoticeResult,
    CreateAllianceResult,
    CreateGuildResult,
    DisbandAllianceResult,
    LeaveAllianceResult,
    ExpelAllianceGuildResult,
    DisbandGuildResult,
    GetAllianceResult,
    ExpelGuildResult,
    GetGuildResult,
    GuildService,
    LeaveGuildResult,
    IncreaseGuildCapacityResult,
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
            if (!result.guild) {
                callback(null, { found: false, guild: undefined });
                return;
            }
            const guild = await this.guildService.buildGuildMessage(worldId, result.guild, result.members ?? []);
            callback(null, { found: true, guild });
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

    @Method("increaseGuildCapacity")
    async increaseGuildCapacity(
        call: GrpcCall<IncreaseGuildCapacityRequest>,
        callback: GrpcCallback<IncreaseGuildCapacityReply>
    ) {
        try {
            const req = call.request;
            const result = await this.guildService.increaseGuildCapacity(
                req.worldId,
                req.characterId,
                req.extendedCap
            ) as IncreaseGuildCapacityResult;
            if (result.ok) {
                callback(null, {
                    ok: true,
                    errorCode: GuildErrorCode.GUILD_ERROR_NONE,
                    guildId: result.guildId,
                    revision: result.revision ?? 0,
                    capacity: result.capacity ?? 0,
                    gp: result.gp ?? 0,
                });
            } else {
                callback(null, {
                    ok: false,
                    errorCode: result.code ?? GuildErrorCode.GUILD_ERROR_UNKNOWN,
                    guildId: undefined,
                    revision: 0,
                    capacity: 0,
                    gp: 0,
                });
            }
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("listGuildBulletinBoardThreads")
    async listGuildBulletinBoardThreads(
        call: GrpcCall<ListGuildBulletinBoardThreadsRequest>,
        callback: GrpcCallback<ListGuildBulletinBoardThreadsReply>
    ) {
        try {
            const req = call.request;
            const result = await this.guildService.listGuildBulletinBoardThreads(
                req.worldId,
                req.characterId,
                req.page
            ) as ListGuildBulletinBoardThreadsResult;
            if (result.ok) {
                callback(null, {
                    ok: true,
                    errorCode: GuildErrorCode.GUILD_ERROR_NONE,
                    threads: result.threads ?? [],
                    listStart: result.listStart ?? 0,
                    threadCount: result.threadCount ?? 0,
                    notice: result.notice,
                });
            } else {
                callback(null, {
                    ok: false,
                    errorCode: result.code ?? GuildErrorCode.GUILD_ERROR_UNKNOWN,
                    threads: [],
                    listStart: 0,
                    threadCount: 0,
                    notice: undefined,
                });
            }
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("showGuildBulletinBoardThread")
    async showGuildBulletinBoardThread(
        call: GrpcCall<ShowGuildBulletinBoardThreadRequest>,
        callback: GrpcCallback<ShowGuildBulletinBoardThreadReply>
    ) {
        try {
            const req = call.request;
            const result = await this.guildService.showGuildBulletinBoardThread(
                req.worldId,
                req.characterId,
                req.localThreadId
            ) as ShowGuildBulletinBoardThreadResult;
            if (result.ok) {
                callback(null, {
                    ok: true,
                    errorCode: GuildErrorCode.GUILD_ERROR_NONE,
                    thread: result.thread,
                });
            } else {
                callback(null, {
                    ok: false,
                    errorCode: result.code ?? GuildErrorCode.GUILD_ERROR_UNKNOWN,
                    thread: undefined,
                });
            }
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("createGuildBulletinBoardThread")
    async createGuildBulletinBoardThread(
        call: GrpcCall<CreateGuildBulletinBoardThreadRequest>,
        callback: GrpcCallback<CreateGuildBulletinBoardThreadReply>
    ) {
        try {
            const req = call.request;
            const result = await this.guildService.createGuildBulletinBoardThread(
                req.worldId,
                req.characterId,
                req.notice,
                req.title,
                req.body,
                req.icon
            ) as CreateGuildBulletinBoardThreadResult;
            if (result.ok) {
                callback(null, {
                    ok: true,
                    errorCode: GuildErrorCode.GUILD_ERROR_NONE,
                    thread: result.thread,
                    threads: result.threads ?? [],
                    listStart: result.listStart ?? 0,
                    threadCount: result.threadCount ?? 0,
                    notice: result.notice,
                });
            } else {
                callback(null, {
                    ok: false,
                    errorCode: result.code ?? GuildErrorCode.GUILD_ERROR_UNKNOWN,
                    thread: undefined,
                    threads: [],
                    listStart: 0,
                    threadCount: 0,
                    notice: undefined,
                });
            }
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("updateGuildBulletinBoardThread")
    async updateGuildBulletinBoardThread(
        call: GrpcCall<UpdateGuildBulletinBoardThreadRequest>,
        callback: GrpcCallback<UpdateGuildBulletinBoardThreadReply>
    ) {
        try {
            const req = call.request;
            const result = await this.guildService.updateGuildBulletinBoardThread(
                req.worldId,
                req.characterId,
                req.localThreadId,
                req.title,
                req.body,
                req.icon
            ) as UpdateGuildBulletinBoardThreadResult;
            if (result.ok) {
                callback(null, {
                    ok: true,
                    errorCode: GuildErrorCode.GUILD_ERROR_NONE,
                    thread: result.thread,
                });
            } else {
                callback(null, {
                    ok: false,
                    errorCode: result.code ?? GuildErrorCode.GUILD_ERROR_UNKNOWN,
                    thread: undefined,
                });
            }
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("deleteGuildBulletinBoardThread")
    async deleteGuildBulletinBoardThread(
        call: GrpcCall<DeleteGuildBulletinBoardThreadRequest>,
        callback: GrpcCallback<DeleteGuildBulletinBoardThreadReply>
    ) {
        try {
            const req = call.request;
            const result = await this.guildService.deleteGuildBulletinBoardThread(
                req.worldId,
                req.characterId,
                req.localThreadId
            ) as DeleteGuildBulletinBoardThreadResult;
            if (result.ok) {
                callback(null, {
                    ok: true,
                    errorCode: GuildErrorCode.GUILD_ERROR_NONE,
                });
            } else {
                callback(null, {
                    ok: false,
                    errorCode: result.code ?? GuildErrorCode.GUILD_ERROR_UNKNOWN,
                });
            }
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("createGuildBulletinBoardReply")
    async createGuildBulletinBoardReply(
        call: GrpcCall<CreateGuildBulletinBoardReplyRequest>,
        callback: GrpcCallback<CreateGuildBulletinBoardReplyReply>
    ) {
        try {
            const req = call.request;
            const result = await this.guildService.createGuildBulletinBoardReply(
                req.worldId,
                req.characterId,
                req.localThreadId,
                req.content
            ) as CreateGuildBulletinBoardReplyResult;
            if (result.ok) {
                callback(null, {
                    ok: true,
                    errorCode: GuildErrorCode.GUILD_ERROR_NONE,
                    thread: result.thread,
                });
            } else {
                callback(null, {
                    ok: false,
                    errorCode: result.code ?? GuildErrorCode.GUILD_ERROR_UNKNOWN,
                    thread: undefined,
                });
            }
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("deleteGuildBulletinBoardReply")
    async deleteGuildBulletinBoardReply(
        call: GrpcCall<DeleteGuildBulletinBoardReplyRequest>,
        callback: GrpcCallback<DeleteGuildBulletinBoardReplyReply>
    ) {
        try {
            const req = call.request;
            const result = await this.guildService.deleteGuildBulletinBoardReply(
                req.worldId,
                req.characterId,
                req.localThreadId,
                req.replyId
            ) as DeleteGuildBulletinBoardReplyResult;
            if (result.ok) {
                callback(null, {
                    ok: true,
                    errorCode: GuildErrorCode.GUILD_ERROR_NONE,
                    thread: result.thread,
                });
            } else {
                callback(null, {
                    ok: false,
                    errorCode: result.code ?? GuildErrorCode.GUILD_ERROR_UNKNOWN,
                    thread: undefined,
                });
            }
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("createAlliance")
    async createAlliance(call: GrpcCall<CreateAllianceRequest>, callback: GrpcCallback<CreateAllianceReply>) {
        try {
            const req = call.request;
            const result = await this.guildService.createAlliance(
                req.worldId,
                req.allianceName,
                req.leaderCharacterId,
                req.partnerCharacterId
            ) as CreateAllianceResult;
            if (result.ok) {
                callback(null, {
                    ok: true,
                    errorCode: AllianceErrorCode.ALLIANCE_ERROR_NONE,
                    allianceId: result.allianceId ?? 0,
                    revision: result.revision ?? 0,
                    alliance: result.alliance,
                });
            } else {
                callback(null, {
                    ok: false,
                    errorCode: result.code ?? AllianceErrorCode.ALLIANCE_ERROR_UNKNOWN,
                    allianceId: 0,
                    revision: 0,
                    alliance: undefined,
                });
            }
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("leaveAlliance")
    async leaveAlliance(call: GrpcCall<LeaveAllianceRequest>, callback: GrpcCallback<LeaveAllianceReply>) {
        try {
            const result = await this.guildService.leaveAlliance(
                call.request.worldId,
                call.request.characterId
            ) as LeaveAllianceResult;
            if (result.ok) {
                callback(null, {
                    ok: true,
                    errorCode: AllianceErrorCode.ALLIANCE_ERROR_NONE,
                    allianceId: result.allianceId ?? 0,
                    revision: result.revision ?? 0,
                    disbanded: result.disbanded ?? false,
                    removedGuildId: result.removedGuildId ?? 0,
                    alliance: result.alliance,
                });
            } else {
                callback(null, {
                    ok: false,
                    errorCode: result.code ?? AllianceErrorCode.ALLIANCE_ERROR_UNKNOWN,
                    allianceId: 0,
                    revision: 0,
                    disbanded: false,
                    removedGuildId: 0,
                    alliance: undefined,
                });
            }
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("expelAllianceGuild")
    async expelAllianceGuild(
        call: GrpcCall<ExpelAllianceGuildRequest>,
        callback: GrpcCallback<ExpelAllianceGuildReply>
    ) {
        try {
            const req = call.request;
            const result = await this.guildService.expelAllianceGuild(
                req.worldId,
                req.characterId,
                req.targetGuildId,
                req.allianceId > 0 ? req.allianceId : undefined
            ) as ExpelAllianceGuildResult;
            if (result.ok) {
                callback(null, {
                    ok: true,
                    errorCode: AllianceErrorCode.ALLIANCE_ERROR_NONE,
                    allianceId: result.allianceId ?? 0,
                    revision: result.revision ?? 0,
                    disbanded: result.disbanded ?? false,
                    removedGuildId: result.removedGuildId ?? 0,
                    alliance: result.alliance,
                });
            } else {
                callback(null, {
                    ok: false,
                    errorCode: result.code ?? AllianceErrorCode.ALLIANCE_ERROR_UNKNOWN,
                    allianceId: 0,
                    revision: 0,
                    disbanded: false,
                    removedGuildId: 0,
                    alliance: undefined,
                });
            }
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("disbandAlliance")
    async disbandAlliance(call: GrpcCall<DisbandAllianceRequest>, callback: GrpcCallback<DisbandAllianceReply>) {
        try {
            const result = await this.guildService.disbandAlliance(
                call.request.worldId,
                call.request.characterId
            ) as DisbandAllianceResult;
            if (result.ok) {
                callback(null, {
                    ok: true,
                    errorCode: AllianceErrorCode.ALLIANCE_ERROR_NONE,
                    allianceId: result.allianceId ?? 0,
                    revision: result.revision ?? 0,
                });
            } else {
                callback(null, {
                    ok: false,
                    errorCode: result.code ?? AllianceErrorCode.ALLIANCE_ERROR_UNKNOWN,
                    allianceId: 0,
                    revision: 0,
                });
            }
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("acceptAllianceInvite")
    async acceptAllianceInvite(
        call: GrpcCall<AcceptAllianceInviteRequest>,
        callback: GrpcCallback<AcceptAllianceInviteReply>
    ) {
        try {
            const req = call.request;
            const result = await this.guildService.acceptAllianceInvite(
                req.worldId,
                req.characterId,
                req.allianceId,
                req.guildId
            ) as AcceptAllianceInviteResult;
            if (result.ok) {
                callback(null, {
                    ok: true,
                    errorCode: AllianceErrorCode.ALLIANCE_ERROR_NONE,
                    allianceId: result.allianceId ?? 0,
                    revision: result.revision ?? 0,
                    alliance: result.alliance,
                });
            } else {
                callback(null, {
                    ok: false,
                    errorCode: result.code ?? AllianceErrorCode.ALLIANCE_ERROR_UNKNOWN,
                    allianceId: 0,
                    revision: 0,
                    alliance: undefined,
                });
            }
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("increaseAllianceCapacity")
    async increaseAllianceCapacity(
        call: GrpcCall<IncreaseAllianceCapacityRequest>,
        callback: GrpcCallback<IncreaseAllianceCapacityReply>
    ) {
        try {
            const result = await this.guildService.increaseAllianceCapacity(
                call.request.worldId,
                call.request.characterId
            ) as IncreaseAllianceCapacityResult;
            if (result.ok) {
                callback(null, {
                    ok: true,
                    errorCode: AllianceErrorCode.ALLIANCE_ERROR_NONE,
                    allianceId: result.allianceId ?? 0,
                    revision: result.revision ?? 0,
                    alliance: result.alliance,
                });
            } else {
                callback(null, {
                    ok: false,
                    errorCode: result.code ?? AllianceErrorCode.ALLIANCE_ERROR_UNKNOWN,
                    allianceId: 0,
                    revision: 0,
                    alliance: undefined,
                });
            }
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("changeAllianceRankTitles")
    async changeAllianceRankTitles(
        call: GrpcCall<ChangeAllianceRankTitlesRequest>,
        callback: GrpcCallback<ChangeAllianceRankTitlesReply>
    ) {
        try {
            const req = call.request;
            const result = await this.guildService.changeAllianceRankTitles(
                req.worldId,
                req.characterId,
                req.rankTitles
            ) as ChangeAllianceRankTitlesResult;
            if (result.ok) {
                callback(null, {
                    ok: true,
                    errorCode: AllianceErrorCode.ALLIANCE_ERROR_NONE,
                    allianceId: result.allianceId ?? 0,
                    revision: result.revision ?? 0,
                    alliance: result.alliance,
                });
            } else {
                callback(null, {
                    ok: false,
                    errorCode: result.code ?? AllianceErrorCode.ALLIANCE_ERROR_UNKNOWN,
                    allianceId: 0,
                    revision: 0,
                    alliance: undefined,
                });
            }
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("changeAllianceLeader")
    async changeAllianceLeader(
        call: GrpcCall<ChangeAllianceLeaderRequest>,
        callback: GrpcCallback<ChangeAllianceLeaderReply>
    ) {
        try {
            const req = call.request;
            const result = await this.guildService.changeAllianceLeader(
                req.worldId,
                req.characterId,
                req.newLeaderCharacterId
            ) as ChangeAllianceLeaderResult;
            if (result.ok) {
                callback(null, {
                    ok: true,
                    errorCode: AllianceErrorCode.ALLIANCE_ERROR_NONE,
                    allianceId: result.allianceId ?? 0,
                    revision: result.revision ?? 0,
                    oldLeaderCharacterId: result.oldLeaderCharacterId ?? 0,
                    newLeaderCharacterId: result.newLeaderCharacterId ?? 0,
                    alliance: result.alliance,
                });
            } else {
                callback(null, {
                    ok: false,
                    errorCode: result.code ?? AllianceErrorCode.ALLIANCE_ERROR_UNKNOWN,
                    allianceId: 0,
                    revision: 0,
                    oldLeaderCharacterId: 0,
                    newLeaderCharacterId: 0,
                    alliance: undefined,
                });
            }
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("changeAllianceNotice")
    async changeAllianceNotice(
        call: GrpcCall<ChangeAllianceNoticeRequest>,
        callback: GrpcCallback<ChangeAllianceNoticeReply>
    ) {
        try {
            const req = call.request;
            const result = await this.guildService.changeAllianceNotice(
                req.worldId,
                req.characterId,
                req.notice
            ) as ChangeAllianceNoticeResult;
            if (result.ok) {
                callback(null, {
                    ok: true,
                    errorCode: AllianceErrorCode.ALLIANCE_ERROR_NONE,
                    allianceId: result.allianceId ?? 0,
                    revision: result.revision ?? 0,
                    alliance: result.alliance,
                });
            } else {
                callback(null, {
                    ok: false,
                    errorCode: result.code ?? AllianceErrorCode.ALLIANCE_ERROR_UNKNOWN,
                    allianceId: 0,
                    revision: 0,
                    alliance: undefined,
                });
            }
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("changeAllianceMemberRank")
    async changeAllianceMemberRank(
        call: GrpcCall<ChangeAllianceMemberRankRequest>,
        callback: GrpcCallback<ChangeAllianceMemberRankReply>
    ) {
        try {
            const req = call.request;
            const result = await this.guildService.changeAllianceMemberRank(
                req.worldId,
                req.requesterCharacterId,
                req.targetCharacterId,
                req.promote
            ) as ChangeAllianceMemberRankResult;
            if (result.ok) {
                callback(null, {
                    ok: true,
                    errorCode: AllianceErrorCode.ALLIANCE_ERROR_NONE,
                    allianceId: result.allianceId ?? 0,
                    revision: result.revision ?? 0,
                    targetCharacterId: result.targetCharacterId ?? 0,
                    newAllianceRank: result.newAllianceRank ?? 0,
                    alliance: result.alliance,
                });
            } else {
                callback(null, {
                    ok: false,
                    errorCode: result.code ?? AllianceErrorCode.ALLIANCE_ERROR_UNKNOWN,
                    allianceId: 0,
                    revision: 0,
                    targetCharacterId: 0,
                    newAllianceRank: 0,
                    alliance: undefined,
                });
            }
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("getAlliance")
    async getAlliance(call: GrpcCall<GetAllianceRequest>, callback: GrpcCallback<GetAllianceReply>) {
        try {
            const worldId = call.request.worldId;
            const result = await this.guildService.getAlliance(worldId, call.request.allianceId) as GetAllianceResult;
            if (!result.alliance) {
                callback(null, { found: false, alliance: undefined });
                return;
            }
            callback(null, { found: true, alliance: result.alliance });
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
