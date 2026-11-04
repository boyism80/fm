import type {
    MedalRankingRow,
} from "@fm/shared/repos/medal-ranking-repository";
import type {
    MedalService,
} from "../../services/medal-service";
import type {
    ListMedalRankingsReply,
    ListMedalRankingsRequest,
    MedalRanking,
    SubmitMedalScoreReply,
    SubmitMedalScoreRequest,
} from "@fm/shared/protobuf/generated/fminternal/internal_service";
import type {
    GrpcCall,
    GrpcCallback,
} from "./types";
import { Controller, Method } from "../grpc-method-decorator";

function makeMedalRanking(ranking: MedalRankingRow): MedalRanking {
    return {
        kind: ranking.kind,
        month: ranking.month,
        entries: ranking.entries.map((entry) => ({
            characterId: entry.characterId,
            name: entry.name,
            score: entry.score,
        })),
    };
}

@Controller("medalController")
export class MedalGrpcController {
    private readonly medalService: MedalService;

    constructor(medalService: MedalService) {
        this.medalService = medalService;
    }

    @Method("submitMedalScore")
    async submitMedalScore(call: GrpcCall<SubmitMedalScoreRequest>, callback: GrpcCallback<SubmitMedalScoreReply>) {
        const { worldId, kind, month, characterId, name, score } = call.request;
        const result = await this.medalService.submit(
            worldId >>> 0,
            kind,
            month >>> 0,
            characterId >>> 0,
            name,
            score,
        );
        callback(null, { rank: result.rank, ranking: makeMedalRanking(result.ranking) });
    }

    @Method("listMedalRankings")
    async listMedalRankings(call: GrpcCall<ListMedalRankingsRequest>, callback: GrpcCallback<ListMedalRankingsReply>) {
        const rankings = await this.medalService.list(call.request.worldId >>> 0);
        callback(null, { rankings: rankings.map(makeMedalRanking) });
    }
}
