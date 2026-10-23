import { randomUUID } from "node:crypto";
import type { Pool } from "pg";
import type { InternalContext } from "@fm/shared/context/internal-context";
import { guildWeek, guildWeekStart } from "@fm/shared/guild-week";
import {
    JOB_GROUP_SQL,
    JOB_ID_SQL,
    RANKED_CLASS_SQL,
    rankingFilters,
    type RankingKind,
    type WorldGroups,
} from "@fm/shared/ranking-filters";
import {
    KEPT_ROUNDS,
    RankingRepository,
    type RankingBounds,
    type RankingKeys,
    type RankingRow,
} from "@fm/shared/repos/ranking-repository";
import { RangeReader } from "./range-reader";
import { WorkTable, type RangeStep, type RunningRound } from "./work-table";

const HOUR_MS = 60 * 60 * 1000;
const DAY_MS = 24 * HOUR_MS;
const KST_OFFSET_MS = 9 * HOUR_MS;
const POLL_MS = 60 * 1000;
const BUSY_POLL_MS = 1000;
const HEARTBEAT_MS = 10 * 1000;
const RANGE_LEASE_MS = 15 * 60 * 1000;
const SAMPLES_PER_RANGE = 32;
const WRITE_BATCH = 5000;
const OVERALL_MIN_LEVEL = 10;

export type RankingSchedulerOptions = {
    dailyHour: number;
    rangeMultiplier: number;
    workers: number;
    worldGroups: WorldGroups;
};

type KindRun = {
    kind: RankingKind;
    period: string;
    // Columns in ranking_inputs order. Every key is ascending, so (k1, k2, k3, k4) order is ranking order.
    select: (roundId: number) => string;
};

// A round freezes every shard's inputs, splits the key space into ranges, counts every range per filter, gives
// each range its starting rank per filter, ranks the ranges into the round table, checks it, and attaches it.
// Any number of scheduler processes share the range work; starting, advancing, and publishing a round run under
// a per-kind lock.
export class RankingScheduler {
    private readonly ctx: InternalContext;
    private readonly rankings: RankingRepository;
    private readonly work: WorkTable;
    private readonly options: RankingSchedulerOptions;
    private readonly nodeId = randomUUID();
    private stopped = false;
    private loop: Promise<void> = Promise.resolve();
    private timer: NodeJS.Timeout | undefined;
    private heartbeatTimer: NodeJS.Timeout | undefined;
    private wake: (() => void) | undefined;

    constructor(
        internalContext: InternalContext,
        rankingRepository: RankingRepository,
        options: RankingSchedulerOptions,
    ) {
        this.ctx = internalContext;
        this.rankings = rankingRepository;
        this.work = new WorkTable(internalContext.getRedisUnifiedAccess().client);
        this.options = options;
    }

    async start(): Promise<void> {
        await this.work.heartbeat(this.nodeId, this.options.workers);
        this.heartbeatTimer = setInterval(() => {
            this.work
                .heartbeat(this.nodeId, this.options.workers)
                .catch((error: unknown) => console.error("[ranking] heartbeat failed", error));
        }, HEARTBEAT_MS);
        this.loop = this.poll();
    }

    async stop(): Promise<void> {
        this.stopped = true;
        clearInterval(this.heartbeatTimer);
        clearTimeout(this.timer);
        this.wake?.();
        await this.loop;
    }

    // Every poll moves each due round forward, so a scheduler that was down catches up on start.
    private async poll(): Promise<void> {
        while (!this.stopped) {
            let busy = false;
            for (const run of this.runs(Date.now())) {
                if (this.stopped) {
                    return;
                }
                const unfinished = await this.step(run).catch((error: unknown) => {
                    console.error(`[ranking] kind=${run.kind} period=${run.period} failed`, error);
                    return false;
                });
                busy = busy || unfinished;
            }

            await new Promise<void>((resolve) => {
                this.wake = resolve;
                this.timer = setTimeout(resolve, busy ? BUSY_POLL_MS : POLL_MS);
            });
        }
    }

    private runs(now: number): KindRun[] {
        const kstNow = now + KST_OFFSET_MS;
        let dailyAt = Math.floor(kstNow / DAY_MS) * DAY_MS + this.options.dailyHour * HOUR_MS - KST_OFFSET_MS;
        if (dailyAt > now) {
            dailyAt -= DAY_MS;
        }
        const day = new Date(dailyAt + KST_OFFSET_MS).toISOString().slice(0, 10);

        // Thursday shows last week's final guild ranking; the current week shows from the first refresh a full day in.
        const week = guildWeek(now);
        const progress = dailyAt - guildWeekStart(week) >= DAY_MS;
        const guildWeekShown = progress ? week : week - 1;
        const guildPeriod = progress ? `w${week}:${day}` : `w${week - 1}`;

        const worlds = Object.keys(this.ctx.appConfiguration.postgresql.worlds).map(Number).join(", ");
        const ranked = `world_id IN (${worlds}) AND NOT deleted AND NOT hidden AND role = 0 AND ${RANKED_CLASS_SQL}`;
        return [
            {
                kind: "overall",
                period: day,
                select: (roundId) =>
                    `SELECT 'overall', ${roundId}, -level, -exp, (EXTRACT(EPOCH FROM exp_reached_at) * 1000000)::bigint, id,
                            id, name, world_id, class_id, ${JOB_GROUP_SQL}, ${JOB_ID_SQL}, level, exp
                     FROM characters WHERE ${ranked} AND level >= ${OVERALL_MIN_LEVEL}`,
            },
            {
                kind: "fame",
                period: day,
                select: (roundId) =>
                    `SELECT 'fame', ${roundId}, -population, -level, -exp, id,
                            id, name, world_id, class_id, NULL, NULL, level, population
                     FROM characters WHERE ${ranked}`,
            },
            {
                kind: "guild",
                period: guildPeriod,
                select: (roundId) =>
                    `SELECT 'guild', ${roundId}, -weekly_gp, guild_id, 0, 0,
                            guild_id, name, world_id, NULL, NULL, NULL, NULL, weekly_gp
                     FROM guilds
                     WHERE world_id IN (${worlds}) AND disbanded_at IS NULL AND gp_week = ${guildWeekShown} AND weekly_gp > 0`,
            },
        ];
    }

    private async step(run: KindRun): Promise<boolean> {
        const published = await this.rankings.publishedRounds(run.kind);
        if (published[0]?.periodToken === RankingRepository.periodToken(run.period)) {
            return false;
        }
        const running = await this.work.running(run.kind);
        if (running === null || running.period !== run.period || running.phase === "freezing") {
            await this.withLock(run.kind, () => this.startRound(run));
            return true;
        }

        await Promise.all(Array.from({ length: this.options.workers }, () => this.runRanges(run, running)));
        await this.withLock(run.kind, () => this.advance(run));
        return true;
    }

    // Several schedulers may run at once; whoever holds the lock re-reads the round before acting on it.
    private async withLock(kind: RankingKind, work: () => Promise<void>): Promise<void> {
        const token = await this.work.lock(kind);
        if (token === null) {
            return;
        }
        const renew = setInterval(() => {
            this.work
                .renewLock(kind, token)
                .catch((error: unknown) => console.error("[ranking] lock renewal failed", error));
        }, WorkTable.lockRenewMs());
        try {
            await work();
        } finally {
            clearInterval(renew);
            await this.work.unlock(kind, token);
        }
    }

    // Each shard copies its inputs in one statement: a consistent snapshot whose cost grows only with the shard.
    private async startRound(run: KindRun): Promise<void> {
        const running = await this.work.running(run.kind);
        if (running !== null && running.period === run.period && running.phase !== "freezing") {
            return;
        }
        if (running !== null && running.period !== run.period) {
            await this.abandon(run.kind, running);
        }
        const roundId =
            running !== null && running.period === run.period ? running.roundId : await this.rankings.nextRoundId();
        const table = RankingRepository.roundTable(run.kind, roundId, run.period);
        await this.work.setRunning(run.kind, { roundId, period: run.period, phase: "freezing", ranges: 0 });
        await this.dropUnattached(run.kind, roundId);
        await this.rankings.createRoundTable(table, roundId);

        const shards = this.shardPools();
        const counts = await Promise.all(
            shards.map(async (shard) => {
                await this.rankings.clearInputs(shard, run.kind);
                return this.rankings.freezeInputs(shard, run.select(roundId));
            }),
        );
        const total = counts.reduce((sum, n) => sum + n, 0);

        const rangeCount = Math.max(1, Math.min(total, (await this.work.liveWorkers()) * this.options.rangeMultiplier));
        const splits: RankingKeys[] = [];
        if (rangeCount > 1) {
            const fraction = Math.min(1, (SAMPLES_PER_RANGE * rangeCount) / total);
            const samples = (
                await Promise.all(shards.map((shard) => this.rankings.sampleKeys(shard, roundId, fraction)))
            )
                .flat()
                .sort(RankingScheduler.compare);
            for (let i = 1; i < rangeCount; i += 1) {
                const split = samples[Math.floor((i * samples.length) / rangeCount)];
                if (
                    split !== undefined &&
                    (splits.length === 0 || RankingScheduler.compare(split, splits[splits.length - 1]!) > 0)
                ) {
                    splits.push(split);
                }
            }
        }
        const bounds: RankingBounds[] = [null, ...splits].map((lo, i) => ({ lo, hi: splits[i] ?? null }));
        await this.work.writeRanges(roundId, bounds);
        await this.work.setRunning(run.kind, { roundId, period: run.period, phase: "counting", ranges: bounds.length });
        console.log(
            `[ranking] kind=${run.kind} period=${run.period} round=${roundId} frozen=${total} ranges=${bounds.length}`,
        );
    }

    private async runRanges(run: KindRun, running: RunningRound): Promise<void> {
        while (!this.stopped) {
            const step: RangeStep = running.phase === "counting" ? "count" : "rank";
            const index = await this.work.claim(running.roundId, step, RANGE_LEASE_MS);
            if (index === null) {
                return;
            }
            if (step === "count") {
                await this.count(run, running, index);
            } else {
                await this.rank(run, running, index);
            }
        }
    }

    private async count(run: KindRun, running: RunningRound, index: number): Promise<void> {
        const task = await this.work.range(running.roundId, index);
        const counts: Record<string, number> = {};
        const shardCounts = await Promise.all(
            this.shardPools().map((shard) => this.rankings.countRange(shard, running.roundId, task.bounds)),
        );
        for (const combination of shardCounts.flat()) {
            for (const filter of rankingFilters(run.kind, this.options.worldGroups, combination)) {
                counts[filter] = (counts[filter] ?? 0) + combination.count;
            }
        }
        await this.work.finish(running.roundId, "count", index, counts);
    }

    private async rank(run: KindRun, running: RunningRound, index: number): Promise<void> {
        const task = await this.work.range(running.roundId, index);
        const table = RankingRepository.roundTable(run.kind, running.roundId, running.period);
        const readers = this.shardPools().map(
            (shard) => new RangeReader(this.rankings, shard, running.roundId, task.bounds),
        );
        const written: Record<string, number> = {};
        let batch: RankingRow[] = [];

        while (true) {
            let next: RangeReader | null = null;
            let nextKeys: RankingKeys | null = null;
            for (const reader of readers) {
                const head = await reader.peek();
                if (head === undefined) {
                    continue;
                }
                if (nextKeys === null || RankingScheduler.compare(head.keys, nextKeys) < 0) {
                    next = reader;
                    nextKeys = head.keys;
                }
            }
            if (next === null) {
                break;
            }

            const input = next.take();
            for (const filter of rankingFilters(run.kind, this.options.worldGroups, input)) {
                written[filter] = (written[filter] ?? 0) + 1;
                batch.push({ filter, rankNo: (task.offsets[filter] ?? 0) + written[filter], input });
            }
            if (batch.length >= WRITE_BATCH) {
                await this.rankings.insertRows(table, running.roundId, batch);
                batch = [];
            }
        }
        if (batch.length > 0) {
            await this.rankings.insertRows(table, running.roundId, batch);
        }

        for (const filter of new Set([...Object.keys(task.counts), ...Object.keys(written)])) {
            if ((task.counts[filter] ?? 0) !== (written[filter] ?? 0)) {
                throw new Error(
                    `range ${index} wrote ${written[filter] ?? 0} rows for ${filter}, counted ${task.counts[filter] ?? 0}`,
                );
            }
        }
        await this.work.finish(running.roundId, "rank", index, null);
    }

    private async advance(run: KindRun): Promise<void> {
        const running = await this.work.running(run.kind);
        if (running === null || running.period !== run.period) {
            return;
        }
        if (!(await this.work.hasRanges(running.roundId))) {
            console.error(`[ranking] kind=${run.kind} round=${running.roundId} lost its ranges; round dropped`);
            await this.abandon(run.kind, running);
            return;
        }

        if (running.phase === "counting") {
            if ((await this.work.done(running.roundId, "count")) !== running.ranges) {
                return;
            }
            // Each range starts after every subject of the same filter in the ranges before it.
            const offsets: Array<Record<string, number>> = [];
            const totals: Record<string, number> = {};
            for (const counts of await this.work.allCounts(running.roundId, running.ranges)) {
                offsets.push({ ...totals });
                for (const [filter, n] of Object.entries(counts)) {
                    totals[filter] = (totals[filter] ?? 0) + n;
                }
            }
            await this.work.startRanking(running.roundId, offsets);
            await this.work.setRunning(run.kind, { ...running, phase: "ranking" });
            console.log(
                `[ranking] kind=${run.kind} period=${run.period} round=${running.roundId} counted filters=${Object.keys(totals).length}`,
            );
            return;
        }

        if ((await this.work.done(running.roundId, "rank")) !== running.ranges) {
            return;
        }
        await this.publish(run, running);
    }

    private async publish(run: KindRun, running: RunningRound): Promise<void> {
        const table = RankingRepository.roundTable(run.kind, running.roundId, running.period);
        const expected: Record<string, number> = {};
        for (const counts of await this.work.allCounts(running.roundId, running.ranges)) {
            for (const [filter, n] of Object.entries(counts)) {
                expected[filter] = (expected[filter] ?? 0) + n;
            }
        }
        const totals = await this.rankings.filterTotals(table);
        for (const filter of new Set([...Object.keys(expected), ...totals.keys()])) {
            const total = totals.get(filter);
            if (total === undefined || total.count !== expected[filter] || total.last !== expected[filter]) {
                console.error(
                    `[ranking] kind=${run.kind} round=${running.roundId} filter=${filter} expected=${expected[filter] ?? 0} got=${JSON.stringify(total)}; round dropped`,
                );
                await this.abandon(run.kind, running);
                return;
            }
        }

        await this.rankings.attach(table, running.roundId);
        const published = await this.rankings.publishedRounds(run.kind);
        await this.rankings.setCurrentRounds(
            run.kind,
            published.slice(0, KEPT_ROUNDS).map((round) => round.roundId),
        );
        await this.work.clearRound(run.kind, running);
        console.log(
            `[ranking] kind=${run.kind} period=${run.period} round=${running.roundId} published filters=${totals.size}`,
        );

        for (const old of published.slice(KEPT_ROUNDS)) {
            await this.rankings.detachAndDrop(old.table);
        }
        for (const shard of this.shardPools()) {
            await this.rankings.clearInputs(shard, run.kind);
        }
    }

    private async abandon(kind: RankingKind, running: RunningRound): Promise<void> {
        await this.rankings.dropTable(RankingRepository.roundTable(kind, running.roundId, running.period));
        await this.work.clearRound(kind, running);
    }

    // Tables left by a round that never finished (a scheduler died, Redis lost the round) are dropped.
    private async dropUnattached(kind: RankingKind, keepRoundId: number): Promise<void> {
        for (const round of await this.rankings.unattachedTables(kind)) {
            if (round.roundId !== keepRoundId) {
                await this.rankings.dropTable(round.table);
            }
        }
    }

    // Shards of every world, once per database: a unified dev setup lists one database several times.
    private shardPools(): Pool[] {
        const pools: Pool[] = [];
        const seen = new Set<string>();
        for (const [worldKey, world] of Object.entries(this.ctx.appConfiguration.postgresql.worlds)) {
            const worldPools = this.ctx.getPgDataPools(Number(worldKey));
            world.data.forEach((endpoint, i) => {
                const key = `${endpoint.host}:${endpoint.port}/${endpoint.database}`;
                if (!seen.has(key)) {
                    seen.add(key);
                    pools.push(worldPools[i]!);
                }
            });
        }
        return pools;
    }

    private static compare(a: RankingKeys, b: RankingKeys): number {
        for (let i = 0; i < a.length; i += 1) {
            const diff = a[i]! - b[i]!;
            if (diff !== 0) {
                return diff;
            }
        }
        return 0;
    }
}
