export type RankingKind = "overall" | "fame" | "guild";

// Jobs follow the official ranking page: explorers by branch, every other family as one job. job_group is the base
// class of the branch or family (0, 100..500, 1000, 2000); job_id is the base class of the line, and there is none
// before the class has picked its line (beginners, the first explorer class, Noblesse, Legend).
export const JOB_GROUP_SQL = "CASE WHEN class_id < 1000 THEN class_id / 100 * 100 ELSE class_id / 1000 * 1000 END";
export const JOB_ID_SQL =
    "CASE WHEN class_id < 1000 THEN CASE WHEN class_id % 100 = 0 THEN NULL ELSE class_id - class_id % 10 END " +
    "ELSE CASE WHEN class_id % 1000 = 0 THEN NULL ELSE class_id - class_id % 100 END END";
// GM classes (900, 910) are never ranked.
export const RANKED_CLASS_SQL = "class_id / 100 <> 9";

// A world belongs to one group, and each group has its own "all worlds" ranking.
export type WorldGroups = ReadonlyMap<number, number>;

export type RankingSubject = {
    worldId: number;
    jobGroup: number | null;
    jobId: number | null;
};

// Filter keys: g{group} for a group's all worlds, w{world} for one world, then :j{job_group} and :l{job_id}.
export function rankingFilters(kind: RankingKind, groups: WorldGroups, subject: RankingSubject): string[] {
    const worlds = [`g${groups.get(subject.worldId) ?? 0}`, `w${subject.worldId}`];
    if (kind !== "overall" || subject.jobGroup === null) {
        return worlds;
    }
    const filters: string[] = [];
    for (const world of worlds) {
        filters.push(world, `${world}:j${subject.jobGroup}`);
        if (subject.jobId !== null) {
            filters.push(`${world}:j${subject.jobGroup}:l${subject.jobId}`);
        }
    }
    return filters;
}
