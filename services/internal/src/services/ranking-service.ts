import type { RankingKind } from "@fm/shared/ranking-filters";
import type { RankingEntry, RankingRepository } from "@fm/shared/repos/ranking-repository";

const MAX_TOP = 100;

export type CharacterRank = {
    rank: number;
    rankDiff: number;
    classRank: number;
    classRankDiff: number;
};

export class RankingService {
    private readonly rankingRepo: RankingRepository;

    constructor(rankingRepository: RankingRepository) {
        this.rankingRepo = rankingRepository;
    }

    async top(kind: RankingKind, worldId: number, count: number): Promise<RankingEntry[]> {
        const [current] = await this.rankingRepo.currentRounds(kind);
        if (current === undefined) {
            return [];
        }
        return this.rankingRepo.top(current, `w${worldId}`, Math.min(count, MAX_TOP));
    }

    // The character select screen keeps the legacy ranks: the rank among all worlds of the character's group, and
    // the class rank by class_id / 100, which is the branch for explorers, the line for Cygnus Knights and Aran,
    // and the family for Noblesse and Legend. A character missing from the previous round has no rank change.
    async characterRanks(characterIds: number[]): Promise<Map<number, CharacterRank>> {
        const ranks = new Map<number, CharacterRank>();
        if (characterIds.length === 0) {
            return ranks;
        }
        const roundIds = await this.rankingRepo.currentRounds("overall");
        const current = roundIds[0];
        if (current === undefined) {
            return ranks;
        }

        const entries = await this.rankingRepo.subjectEntries(roundIds, characterIds);
        const previous = new Map<string, number>();
        const groupEntries = new Map<number, RankingEntry[]>();
        for (const entry of entries) {
            if (!entry.filter.startsWith("g")) {
                continue;
            }
            if (entry.roundId !== current) {
                previous.set(`${entry.subjectId}|${entry.filter}`, entry.rankNo);
                continue;
            }
            const list = groupEntries.get(entry.subjectId) ?? [];
            list.push(entry);
            groupEntries.set(entry.subjectId, list);
        }

        for (const [characterId, list] of groupEntries) {
            const all = list.find((entry) => !entry.filter.includes(":"));
            if (all === undefined) {
                continue;
            }
            const line = (all.classId ?? 0) >= 1000 ? list.find((entry) => entry.filter.includes(":l")) : undefined;
            const job = line ?? list.find((entry) => !entry.filter.includes(":l") && entry.filter.includes(":j"));
            const allBefore = previous.get(`${characterId}|${all.filter}`);
            const jobBefore = job === undefined ? undefined : previous.get(`${characterId}|${job.filter}`);
            ranks.set(characterId, {
                rank: all.rankNo,
                rankDiff: allBefore === undefined ? 0 : allBefore - all.rankNo,
                classRank: job?.rankNo ?? 0,
                classRankDiff: job === undefined || jobBefore === undefined ? 0 : jobBefore - job.rankNo,
            });
        }
        return ranks;
    }
}
