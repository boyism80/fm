import type { MedalRankingRepository, MedalRankingRow } from "@fm/shared/repos/medal-ranking-repository";
import type { RabbitMQService } from "./rabbitmq-service";

const MEDAL_RANKING_CHANGED_EVENT = "medal_ranking_changed";
const MAX_KIND_LENGTH = 32;

export class MedalService {
    private readonly medalRankingRepository: MedalRankingRepository;
    private readonly rabbitmqService: RabbitMQService;

    constructor(medalRankingRepository: MedalRankingRepository, rabbitmqService: RabbitMQService) {
        this.medalRankingRepository = medalRankingRepository;
        this.rabbitmqService = rabbitmqService;
    }

    async list(worldId: number): Promise<MedalRankingRow[]> {
        return this.medalRankingRepository.list(worldId);
    }

    // A lost notice leaves the other channels on the old board until they restart, so it does not undo the write.
    async submit(worldId: number, kind: string, month: number, characterId: number, name: string, score: number): Promise<{ rank: number; ranking: MedalRankingRow }> {
        if (kind === "" || kind.length > MAX_KIND_LENGTH || name === "" || characterId === 0 || month < 200001 || month % 100 < 1 || month % 100 > 12) {
            throw Object.assign(new Error("invalid medal score"), { code: "INVALID_PAYLOAD" });
        }
        const { previousHolderId, ranking } = await this.medalRankingRepository.submit(worldId, kind, month, characterId, name, score);
        const holderId = ranking.entries[0]?.characterId ?? 0;
        try {
            await this.rabbitmqService.publish(`fm.${worldId}.all.global`, MEDAL_RANKING_CHANGED_EVENT, {
                ranking: {
                    kind: ranking.kind,
                    month: ranking.month,
                    entries: ranking.entries.map((entry) => ({ character_id: entry.characterId, name: entry.name, score: entry.score })),
                },
                previous_holder_id: previousHolderId !== holderId ? previousHolderId : 0,
            });
        } catch (err) {
            console.error("[medal] ranking notice failed", err);
        }
        const rank = ranking.entries.findIndex((entry) => entry.characterId === characterId) + 1;
        return { rank, ranking };
    }
}
