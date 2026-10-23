import type { RankingKind as RankingKindName } from "@fm/shared/ranking-filters";
import type { RankingService } from "../../services/ranking-service";
import {
    RankingKind,
    type GetRankingReply,
    type GetRankingRequest,
} from "@fm/shared/protobuf/generated/fminternal/internal_service";
import type { GrpcCall, GrpcCallback, GrpcErrorHandler } from "./types";
import { Controller, Method } from "../grpc-method-decorator";

const KIND_NAMES = new Map<RankingKind, RankingKindName>([
    [RankingKind.RANKING_KIND_OVERALL, "overall"],
    [RankingKind.RANKING_KIND_FAME, "fame"],
    [RankingKind.RANKING_KIND_GUILD, "guild"],
]);

@Controller("rankingController")
export class RankingGrpcController {
    private readonly rankingService: RankingService;
    private readonly grpcError: GrpcErrorHandler;

    constructor(rankingService: RankingService, grpcError: GrpcErrorHandler) {
        this.rankingService = rankingService;
        this.grpcError = grpcError;
    }

    @Method("getRanking")
    async getRanking(call: GrpcCall<GetRankingRequest>, callback: GrpcCallback<GetRankingReply>) {
        try {
            const kind = KIND_NAMES.get(call.request.kind);
            if (kind === undefined) {
                throw Object.assign(new Error(`unknown ranking kind ${call.request.kind}`), {
                    code: "INVALID_PAYLOAD",
                });
            }
            const entries = await this.rankingService.top(kind, call.request.worldId >>> 0, call.request.count >>> 0);
            callback(null, {
                entries: entries.map((entry) => ({
                    rank: entry.rankNo,
                    subjectId: entry.subjectId,
                    name: entry.name,
                    worldId: entry.worldId,
                    classId: entry.classId ?? 0,
                    level: entry.level ?? 0,
                    value: entry.value,
                })),
            });
        } catch (err) {
            this.grpcError(err, callback);
        }
    }
}
