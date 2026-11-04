import { AllianceErrorCode, GuildErrorCode } from "@fm/shared/protobuf/generated/fminternal/internal_service";
import type {
    AcceptGuildInviteReply,
    AcceptGuildInviteRequest,
    ChangeGuildEmblemReply,
    ChangeGuildEmblemRequest,
    ChangeGuildNoticeReply,
    ChangeGuildNoticeRequest,
    IncreaseGuildCapacityReply,
    IncreaseGuildCapacityRequest,
    GainGuildGPReply,
    GainGuildGPRequest,
    SendGuildMessageReply,
    SendGuildMessageRequest,
    GetGuildRankingReply,
    GetGuildRankingRequest,
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
    UpdateGuildMemberReply,
    UpdateGuildMemberRequest,
} from "@fm/shared/protobuf/generated/fminternal/internal_service";
import type {
    GuildService,
} from "../../services/guild-service";
import type {
    GuildBulletinBoardService,
} from "../../services/guild-bulletin-board-service";
import type {
    AllianceService,
} from "../../services/alliance-service";
import { okReply } from "./types";
import type {
    GrpcCall,
    GrpcCallback,
} from "./types";
import { Controller, Method } from "../grpc-method-decorator";

@Controller("guildController")
export class GuildGrpcController {
    private readonly guildService: GuildService;
    private readonly guildBulletinBoardService: GuildBulletinBoardService;
    private readonly allianceService: AllianceService;

    constructor(
        guildService: GuildService,
        guildBulletinBoardService: GuildBulletinBoardService,
        allianceService: AllianceService
    ) {
        this.guildService = guildService;
        this.guildBulletinBoardService = guildBulletinBoardService;
        this.allianceService = allianceService;
    }

    @Method("createGuild")
    async createGuild(call: GrpcCall<CreateGuildRequest>, callback: GrpcCallback<CreateGuildReply>) {
        const req = call.request;
        const result = await this.guildService.createGuild(
            req.worldId,
            req.guildName,
            req.leader,
            req.requestId,
        );
        callback(null, okReply(result, GuildErrorCode.GUILD_ERROR_NONE, GuildErrorCode.GUILD_ERROR_UNKNOWN));
    }

    @Method("getGuild")
    async getGuild(call: GrpcCall<GetGuildRequest>, callback: GrpcCallback<GetGuildReply>) {
        const worldId = call.request.worldId;
        const result = await this.guildService.getGuild(worldId, call.request.guildId);
        if (!result.guild) {
            callback(null, { found: false, guild: undefined });
            return;
        }
        const guild = await this.guildService.guildToPb(worldId, result.guild, result.members ?? []);
        callback(null, { found: true, guild });
    }

    @Method("updateGuildMember")
    async updateGuildMember(call: GrpcCall<UpdateGuildMemberRequest>, callback: GrpcCallback<UpdateGuildMemberReply>) {
        await this.guildService.updateMember(call.request);
        callback(null, {});
    }

    @Method("acceptGuildInvite")
    async acceptGuildInvite(call: GrpcCall<AcceptGuildInviteRequest>, callback: GrpcCallback<AcceptGuildInviteReply>) {
        const req = call.request;
        const result = await this.guildService.acceptGuildInvite(
            req.worldId,
            req.guildId,
            req.member,
        );
        callback(null, okReply(result, GuildErrorCode.GUILD_ERROR_NONE, GuildErrorCode.GUILD_ERROR_UNKNOWN));
    }

    @Method("leaveGuild")
    async leaveGuild(call: GrpcCall<LeaveGuildRequest>, callback: GrpcCallback<LeaveGuildReply>) {
        const result = await this.guildService.leaveGuild(
            call.request.worldId,
            call.request.characterId,
        );
        callback(null, okReply(result, GuildErrorCode.GUILD_ERROR_NONE, GuildErrorCode.GUILD_ERROR_UNKNOWN));
    }

    @Method("expelGuild")
    async expelGuild(call: GrpcCall<ExpelGuildRequest>, callback: GrpcCallback<ExpelGuildReply>) {
        const req = call.request;
        const result = await this.guildService.expelGuild(
            req.worldId,
            req.requesterCharacterId,
            req.targetCharacterId,
        );
        callback(null, okReply(result, GuildErrorCode.GUILD_ERROR_NONE, GuildErrorCode.GUILD_ERROR_UNKNOWN));
    }

    @Method("changeGuildRankTitles")
    async changeGuildRankTitles(
        call: GrpcCall<ChangeGuildRankTitlesRequest>,
        callback: GrpcCallback<ChangeGuildRankTitlesReply>,
    ) {
        const req = call.request;
        const result = await this.guildService.changeGuildRankTitles(
            req.worldId,
            req.characterId,
            req.rankTitles,
        );
        callback(null, okReply(result, GuildErrorCode.GUILD_ERROR_NONE, GuildErrorCode.GUILD_ERROR_UNKNOWN));
    }

    @Method("changeGuildMemberRank")
    async changeGuildMemberRank(
        call: GrpcCall<ChangeGuildMemberRankRequest>,
        callback: GrpcCallback<ChangeGuildMemberRankReply>,
    ) {
        const req = call.request;
        const result = await this.guildService.changeGuildMemberRank(
            req.worldId,
            req.requesterCharacterId,
            req.targetCharacterId,
            req.newRank,
        );
        callback(null, okReply(result, GuildErrorCode.GUILD_ERROR_NONE, GuildErrorCode.GUILD_ERROR_UNKNOWN));
    }

    @Method("changeGuildEmblem")
    async changeGuildEmblem(call: GrpcCall<ChangeGuildEmblemRequest>, callback: GrpcCallback<ChangeGuildEmblemReply>) {
        const req = call.request;
        const result = await this.guildService.changeGuildEmblem(
            req.worldId,
            req.characterId,
            req.logo,
            req.requestId,
            req.guildId,
        );
        callback(null, okReply(result, GuildErrorCode.GUILD_ERROR_NONE, GuildErrorCode.GUILD_ERROR_UNKNOWN));
    }

    @Method("changeGuildNotice")
    async changeGuildNotice(call: GrpcCall<ChangeGuildNoticeRequest>, callback: GrpcCallback<ChangeGuildNoticeReply>) {
        const req = call.request;
        const result = await this.guildService.changeGuildNotice(
            req.worldId,
            req.characterId,
            req.notice,
        );
        callback(null, okReply(result, GuildErrorCode.GUILD_ERROR_NONE, GuildErrorCode.GUILD_ERROR_UNKNOWN));
    }

    @Method("increaseGuildCapacity")
    async increaseGuildCapacity(
        call: GrpcCall<IncreaseGuildCapacityRequest>,
        callback: GrpcCallback<IncreaseGuildCapacityReply>,
    ) {
        const req = call.request;
        const result = await this.guildService.increaseGuildCapacity(
            req.worldId,
            req.characterId,
            req.extendedCap,
            req.requestId,
            req.guildId,
        );
        if (result.ok) {
            callback(null, {
                ok: true,
                errorCode: GuildErrorCode.GUILD_ERROR_NONE,
                guildId: result.guildId,
                capacity: result.capacity ?? 0,
                gp: result.gp ?? 0,
            });
        } else {
            callback(null, {
                ok: false,
                errorCode: result.code ?? GuildErrorCode.GUILD_ERROR_UNKNOWN,
                guildId: undefined,
                capacity: 0,
                gp: 0,
            });
        }
    }

    @Method("gainGuildGp")
    async gainGuildGP(call: GrpcCall<GainGuildGPRequest>, callback: GrpcCallback<GainGuildGPReply>) {
        const req = call.request;
        const result = await this.guildService.gainGuildGP(
            req.worldId,
            req.guildId,
            req.amount,
            req.requestId,
        );
        if (result.ok) {
            callback(null, {
                ok: true,
                errorCode: GuildErrorCode.GUILD_ERROR_NONE,
                guildId: result.guildId,
                gp: result.gp ?? 0,
            });
        } else {
            callback(null, {
                ok: false,
                errorCode: result.code ?? GuildErrorCode.GUILD_ERROR_UNKNOWN,
                guildId: undefined,
                gp: 0,
            });
        }
    }

    @Method("sendGuildMessage")
    async sendGuildMessage(call: GrpcCall<SendGuildMessageRequest>, callback: GrpcCallback<SendGuildMessageReply>) {
        const req = call.request;
        const ok = await this.guildService.sendGuildMessage(req.worldId, req.guildId, req.messageType, req.message);
        callback(null, { ok });
    }

    @Method("getGuildRanking")
    async getGuildRanking(call: GrpcCall<GetGuildRankingRequest>, callback: GrpcCallback<GetGuildRankingReply>) {
        const entries = await this.guildService.getGuildRanking(call.request.worldId);
        callback(null, {
            entries: entries.map((entry) => {
                return {
                    guildId: entry.guildId,
                    name: entry.name,
                    gp: entry.gp,
                    logo: {
                        logo: entry.logo.logo,
                        logoColor: entry.logo.logoColor,
                        logoBg: entry.logo.logoBG,
                        logoBgColor: entry.logo.logoBGColor,
                    },
                };
            }),
        });
    }

    @Method("listGuildBulletinBoardThreads")
    async listGuildBulletinBoardThreads(
        call: GrpcCall<ListGuildBulletinBoardThreadsRequest>,
        callback: GrpcCallback<ListGuildBulletinBoardThreadsReply>,
    ) {
        const req = call.request;
        const result = await this.guildBulletinBoardService.listGuildBulletinBoardThreads(
            req.worldId,
            req.characterId,
            req.page,
        );
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
    }

    @Method("showGuildBulletinBoardThread")
    async showGuildBulletinBoardThread(
        call: GrpcCall<ShowGuildBulletinBoardThreadRequest>,
        callback: GrpcCallback<ShowGuildBulletinBoardThreadReply>,
    ) {
        const req = call.request;
        const result = await this.guildBulletinBoardService.showGuildBulletinBoardThread(
            req.worldId,
            req.characterId,
            req.localThreadId,
        );
        callback(null, { ...okReply(result, GuildErrorCode.GUILD_ERROR_NONE, GuildErrorCode.GUILD_ERROR_UNKNOWN), thread: result.thread });
    }

    @Method("createGuildBulletinBoardThread")
    async createGuildBulletinBoardThread(
        call: GrpcCall<CreateGuildBulletinBoardThreadRequest>,
        callback: GrpcCallback<CreateGuildBulletinBoardThreadReply>,
    ) {
        const req = call.request;
        const result = await this.guildBulletinBoardService.createGuildBulletinBoardThread(
            req.worldId,
            req.characterId,
            req.notice,
            req.title,
            req.body,
            req.icon,
        );
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
    }

    @Method("updateGuildBulletinBoardThread")
    async updateGuildBulletinBoardThread(
        call: GrpcCall<UpdateGuildBulletinBoardThreadRequest>,
        callback: GrpcCallback<UpdateGuildBulletinBoardThreadReply>,
    ) {
        const req = call.request;
        const result = await this.guildBulletinBoardService.updateGuildBulletinBoardThread(
            req.worldId,
            req.characterId,
            req.localThreadId,
            req.title,
            req.body,
            req.icon,
        );
        callback(null, { ...okReply(result, GuildErrorCode.GUILD_ERROR_NONE, GuildErrorCode.GUILD_ERROR_UNKNOWN), thread: result.thread });
    }

    @Method("deleteGuildBulletinBoardThread")
    async deleteGuildBulletinBoardThread(
        call: GrpcCall<DeleteGuildBulletinBoardThreadRequest>,
        callback: GrpcCallback<DeleteGuildBulletinBoardThreadReply>,
    ) {
        const req = call.request;
        const result = await this.guildBulletinBoardService.deleteGuildBulletinBoardThread(
            req.worldId,
            req.characterId,
            req.localThreadId,
        );
        callback(null, okReply(result, GuildErrorCode.GUILD_ERROR_NONE, GuildErrorCode.GUILD_ERROR_UNKNOWN));
    }

    @Method("createGuildBulletinBoardReply")
    async createGuildBulletinBoardReply(
        call: GrpcCall<CreateGuildBulletinBoardReplyRequest>,
        callback: GrpcCallback<CreateGuildBulletinBoardReplyReply>,
    ) {
        const req = call.request;
        const result = await this.guildBulletinBoardService.createGuildBulletinBoardReply(
            req.worldId,
            req.characterId,
            req.localThreadId,
            req.content,
        );
        callback(null, { ...okReply(result, GuildErrorCode.GUILD_ERROR_NONE, GuildErrorCode.GUILD_ERROR_UNKNOWN), thread: result.thread });
    }

    @Method("deleteGuildBulletinBoardReply")
    async deleteGuildBulletinBoardReply(
        call: GrpcCall<DeleteGuildBulletinBoardReplyRequest>,
        callback: GrpcCallback<DeleteGuildBulletinBoardReplyReply>,
    ) {
        const req = call.request;
        const result = await this.guildBulletinBoardService.deleteGuildBulletinBoardReply(
            req.worldId,
            req.characterId,
            req.localThreadId,
            req.replyId,
        );
        callback(null, { ...okReply(result, GuildErrorCode.GUILD_ERROR_NONE, GuildErrorCode.GUILD_ERROR_UNKNOWN), thread: result.thread });
    }

    @Method("createAlliance")
    async createAlliance(call: GrpcCall<CreateAllianceRequest>, callback: GrpcCallback<CreateAllianceReply>) {
        const req = call.request;
        const result = await this.allianceService.createAlliance(
            req.worldId,
            req.allianceName,
            req.leaderCharacterId,
            req.partnerCharacterId,
            req.requestId,
            req.guildId,
        );
        if (result.ok) {
            callback(null, {
                ok: true,
                errorCode: AllianceErrorCode.ALLIANCE_ERROR_NONE,
                allianceId: result.allianceId ?? 0,
                alliance: result.alliance,
            });
        } else {
            callback(null, {
                ok: false,
                errorCode: result.code ?? AllianceErrorCode.ALLIANCE_ERROR_UNKNOWN,
                allianceId: 0,
                alliance: undefined,
            });
        }
    }

    @Method("leaveAlliance")
    async leaveAlliance(call: GrpcCall<LeaveAllianceRequest>, callback: GrpcCallback<LeaveAllianceReply>) {
        const result = await this.allianceService.leaveAlliance(
            call.request.worldId,
            call.request.characterId,
        );
        if (result.ok) {
            callback(null, {
                ok: true,
                errorCode: AllianceErrorCode.ALLIANCE_ERROR_NONE,
                allianceId: result.allianceId ?? 0,
                disbanded: result.disbanded ?? false,
                removedGuildId: result.removedGuildId ?? 0,
                alliance: result.alliance,
            });
        } else {
            callback(null, {
                ok: false,
                errorCode: result.code ?? AllianceErrorCode.ALLIANCE_ERROR_UNKNOWN,
                allianceId: 0,
                disbanded: false,
                removedGuildId: 0,
                alliance: undefined,
            });
        }
    }

    @Method("expelAllianceGuild")
    async expelAllianceGuild(
        call: GrpcCall<ExpelAllianceGuildRequest>,
        callback: GrpcCallback<ExpelAllianceGuildReply>,
    ) {
        const req = call.request;
        const result = await this.allianceService.expelAllianceGuild(
            req.worldId,
            req.characterId,
            req.targetGuildId,
            req.allianceId > 0 ? req.allianceId : undefined,
        );
        if (result.ok) {
            callback(null, {
                ok: true,
                errorCode: AllianceErrorCode.ALLIANCE_ERROR_NONE,
                allianceId: result.allianceId ?? 0,
                disbanded: result.disbanded ?? false,
                removedGuildId: result.removedGuildId ?? 0,
                alliance: result.alliance,
            });
        } else {
            callback(null, {
                ok: false,
                errorCode: result.code ?? AllianceErrorCode.ALLIANCE_ERROR_UNKNOWN,
                allianceId: 0,
                disbanded: false,
                removedGuildId: 0,
                alliance: undefined,
            });
        }
    }

    @Method("disbandAlliance")
    async disbandAlliance(call: GrpcCall<DisbandAllianceRequest>, callback: GrpcCallback<DisbandAllianceReply>) {
        const result = await this.allianceService.disbandAlliance(
            call.request.worldId,
            call.request.characterId,
        );
        if (result.ok) {
            callback(null, {
                ok: true,
                errorCode: AllianceErrorCode.ALLIANCE_ERROR_NONE,
                allianceId: result.allianceId ?? 0,
            });
        } else {
            callback(null, {
                ok: false,
                errorCode: result.code ?? AllianceErrorCode.ALLIANCE_ERROR_UNKNOWN,
                allianceId: 0,
            });
        }
    }

    @Method("acceptAllianceInvite")
    async acceptAllianceInvite(
        call: GrpcCall<AcceptAllianceInviteRequest>,
        callback: GrpcCallback<AcceptAllianceInviteReply>,
    ) {
        const req = call.request;
        const result = await this.allianceService.acceptAllianceInvite(
            req.worldId,
            req.characterId,
            req.allianceId,
            req.guildId,
        );
        if (result.ok) {
            callback(null, {
                ok: true,
                errorCode: AllianceErrorCode.ALLIANCE_ERROR_NONE,
                allianceId: result.allianceId ?? 0,
                alliance: result.alliance,
            });
        } else {
            callback(null, {
                ok: false,
                errorCode: result.code ?? AllianceErrorCode.ALLIANCE_ERROR_UNKNOWN,
                allianceId: 0,
                alliance: undefined,
            });
        }
    }

    @Method("increaseAllianceCapacity")
    async increaseAllianceCapacity(
        call: GrpcCall<IncreaseAllianceCapacityRequest>,
        callback: GrpcCallback<IncreaseAllianceCapacityReply>,
    ) {
        const result = await this.allianceService.increaseAllianceCapacity(
            call.request.worldId,
            call.request.characterId,
            call.request.requestId,
            call.request.allianceId,
        );
        if (result.ok) {
            callback(null, {
                ok: true,
                errorCode: AllianceErrorCode.ALLIANCE_ERROR_NONE,
                allianceId: result.allianceId ?? 0,
                alliance: result.alliance,
            });
        } else {
            callback(null, {
                ok: false,
                errorCode: result.code ?? AllianceErrorCode.ALLIANCE_ERROR_UNKNOWN,
                allianceId: 0,
                alliance: undefined,
            });
        }
    }

    @Method("changeAllianceRankTitles")
    async changeAllianceRankTitles(
        call: GrpcCall<ChangeAllianceRankTitlesRequest>,
        callback: GrpcCallback<ChangeAllianceRankTitlesReply>,
    ) {
        const req = call.request;
        const result = await this.allianceService.changeAllianceRankTitles(
            req.worldId,
            req.characterId,
            req.rankTitles,
        );
        if (result.ok) {
            callback(null, {
                ok: true,
                errorCode: AllianceErrorCode.ALLIANCE_ERROR_NONE,
                allianceId: result.allianceId ?? 0,
                alliance: result.alliance,
            });
        } else {
            callback(null, {
                ok: false,
                errorCode: result.code ?? AllianceErrorCode.ALLIANCE_ERROR_UNKNOWN,
                allianceId: 0,
                alliance: undefined,
            });
        }
    }

    @Method("changeAllianceLeader")
    async changeAllianceLeader(
        call: GrpcCall<ChangeAllianceLeaderRequest>,
        callback: GrpcCallback<ChangeAllianceLeaderReply>,
    ) {
        const req = call.request;
        const result = await this.allianceService.changeAllianceLeader(
            req.worldId,
            req.characterId,
            req.newLeaderCharacterId,
        );
        if (result.ok) {
            callback(null, {
                ok: true,
                errorCode: AllianceErrorCode.ALLIANCE_ERROR_NONE,
                allianceId: result.allianceId ?? 0,
                oldLeaderCharacterId: result.oldLeaderCharacterId ?? 0,
                newLeaderCharacterId: result.newLeaderCharacterId ?? 0,
                alliance: result.alliance,
            });
        } else {
            callback(null, {
                ok: false,
                errorCode: result.code ?? AllianceErrorCode.ALLIANCE_ERROR_UNKNOWN,
                allianceId: 0,
                oldLeaderCharacterId: 0,
                newLeaderCharacterId: 0,
                alliance: undefined,
            });
        }
    }

    @Method("changeAllianceNotice")
    async changeAllianceNotice(
        call: GrpcCall<ChangeAllianceNoticeRequest>,
        callback: GrpcCallback<ChangeAllianceNoticeReply>,
    ) {
        const req = call.request;
        const result = await this.allianceService.changeAllianceNotice(
            req.worldId,
            req.characterId,
            req.notice,
        );
        if (result.ok) {
            callback(null, {
                ok: true,
                errorCode: AllianceErrorCode.ALLIANCE_ERROR_NONE,
                allianceId: result.allianceId ?? 0,
                alliance: result.alliance,
            });
        } else {
            callback(null, {
                ok: false,
                errorCode: result.code ?? AllianceErrorCode.ALLIANCE_ERROR_UNKNOWN,
                allianceId: 0,
                alliance: undefined,
            });
        }
    }

    @Method("changeAllianceMemberRank")
    async changeAllianceMemberRank(
        call: GrpcCall<ChangeAllianceMemberRankRequest>,
        callback: GrpcCallback<ChangeAllianceMemberRankReply>,
    ) {
        const req = call.request;
        const result = await this.allianceService.changeAllianceMemberRank(
            req.worldId,
            req.requesterCharacterId,
            req.targetCharacterId,
            req.promote,
        );
        if (result.ok) {
            callback(null, {
                ok: true,
                errorCode: AllianceErrorCode.ALLIANCE_ERROR_NONE,
                allianceId: result.allianceId ?? 0,
                targetCharacterId: result.targetCharacterId ?? 0,
                newAllianceRank: result.newAllianceRank ?? 0,
                alliance: result.alliance,
            });
        } else {
            callback(null, {
                ok: false,
                errorCode: result.code ?? AllianceErrorCode.ALLIANCE_ERROR_UNKNOWN,
                allianceId: 0,
                targetCharacterId: 0,
                newAllianceRank: 0,
                alliance: undefined,
            });
        }
    }

    @Method("getAlliance")
    async getAlliance(call: GrpcCall<GetAllianceRequest>, callback: GrpcCallback<GetAllianceReply>) {
        const worldId = call.request.worldId;
        const result = await this.allianceService.getAlliance(
            worldId,
            call.request.allianceId,
        );
        if (!result.alliance) {
            callback(null, { found: false, alliance: undefined });
            return;
        }
        callback(null, { found: true, alliance: result.alliance });
    }

    @Method("disbandGuild")
    async disbandGuild(call: GrpcCall<DisbandGuildRequest>, callback: GrpcCallback<DisbandGuildReply>) {
        const result = await this.guildService.disbandGuild(
            call.request.worldId,
            call.request.characterId,
        );
        if (result.ok) {
            if (result.allianceId != null && result.guildId != null) {
                await this.allianceService.removeDisbandedGuild(
                    call.request.worldId,
                    result.allianceId,
                    result.guildId,
                );
            }
            callback(null, {
                ok: true,
                errorCode: GuildErrorCode.GUILD_ERROR_NONE,
                guildId: result.guildId,
            });
        } else {
            callback(null, {
                ok: false,
                errorCode: result.code ?? GuildErrorCode.GUILD_ERROR_UNKNOWN,
                guildId: undefined,
            });
        }
    }
}
