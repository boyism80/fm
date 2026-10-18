import { randomUUID } from "node:crypto";
import type Redis from "ioredis";
import type { RankingKind } from "@fm/shared/ranking-filters";
import type { RankingBounds, RankingKeys } from "@fm/shared/repos/ranking-repository";

const ROUND_TTL_SEC = 48 * 60 * 60;
const LOCK_TTL_SEC = 10 * 60;
const NODE_TTL_MS = 30 * 1000;
const NODES_KEY_TTL_SEC = 60;
const NODES_KEY = "fm:ranking:nodes";

export type RankingPhase = "freezing" | "counting" | "ranking";
export type RangeStep = "count" | "rank";

export type RunningRound = {
    roundId: number;
    period: string;
    phase: RankingPhase;
    ranges: number;
};

export type RangeTask = {
    index: number;
    bounds: RankingBounds;
    counts: Record<string, number>;
    offsets: Record<string, number>;
};

// A lease that ran out belonged to a worker that died; its range goes back to the queue unless it is already done.
const CLAIM_SCRIPT = `
local todo, leased, done = KEYS[1], KEYS[2], KEYS[3]
local now, lease, ttl = tonumber(ARGV[1]), tonumber(ARGV[2]), tonumber(ARGV[3])
for _, i in ipairs(redis.call("ZRANGEBYSCORE", leased, "-inf", now)) do
  redis.call("ZREM", leased, i)
  if redis.call("SISMEMBER", done, i) == 0 then
    redis.call("RPUSH", todo, i)
  end
end
while true do
  local i = redis.call("LPOP", todo)
  if not i then
    return false
  end
  if redis.call("SISMEMBER", done, i) == 0 then
    redis.call("ZADD", leased, now + lease, i)
    redis.call("EXPIRE", todo, ttl)
    redis.call("EXPIRE", leased, ttl)
    return i
  end
end
`;

// The first worker to finish a range records its result; a late duplicate of the same range is ignored.
const FINISH_SCRIPT = `
local leased, done, range = KEYS[1], KEYS[2], KEYS[3]
local i, field, value, ttl = ARGV[1], ARGV[2], ARGV[3], tonumber(ARGV[4])
redis.call("ZREM", leased, i)
if redis.call("SADD", done, i) == 1 and field ~= "" then
  redis.call("HSET", range, field, value)
end
redis.call("EXPIRE", done, ttl)
redis.call("EXPIRE", range, ttl)
return 1
`;

const RENEW_SCRIPT = `
if redis.call("GET", KEYS[1]) == ARGV[1] then
  return redis.call("EXPIRE", KEYS[1], ARGV[2])
end
return 0
`;

const RELEASE_SCRIPT = `
if redis.call("GET", KEYS[1]) == ARGV[1] then
  return redis.call("DEL", KEYS[1])
end
return 0
`;

// Everything here is coordination for rounds in progress; every key has a TTL, so nothing outlives a stopped
// scheduler for long. Losing it only restarts the round in progress.
export class WorkTable {
    private readonly redis: Redis;

    constructor(redis: Redis) {
        this.redis = redis;
    }

    private runningKey(kind: RankingKind) {
        return `fm:ranking:${kind}:running`;
    }

    private lockKey(kind: RankingKind) {
        return `fm:lock:ranking:${kind}`;
    }

    // One hash tag per round keeps its keys on one Redis Cluster node, so the scripts can touch them together.
    private roundKey(roundId: number, rest: string) {
        return `fm:ranking:{r:${roundId}}:${rest}`;
    }

    async running(kind: RankingKind): Promise<RunningRound | null> {
        const hash = await this.redis.hgetall(this.runningKey(kind));
        if (!hash.round_id) {
            return null;
        }
        return {
            roundId: Number(hash.round_id),
            period: hash.period ?? "",
            phase: hash.phase as RankingPhase,
            ranges: Number(hash.ranges ?? 0),
        };
    }

    async setRunning(kind: RankingKind, running: RunningRound): Promise<void> {
        const key = this.runningKey(kind);
        await this.redis
            .multi()
            .hset(key, { round_id: running.roundId, period: running.period, phase: running.phase, ranges: running.ranges })
            .expire(key, ROUND_TTL_SEC)
            .exec();
    }

    async clearRound(kind: RankingKind, running: RunningRound): Promise<void> {
        const keys = [this.runningKey(kind)];
        for (const rest of ["todo:count", "todo:rank", "leased:count", "leased:rank", "done:count", "done:rank"]) {
            keys.push(this.roundKey(running.roundId, rest));
        }
        for (let i = 0; i < running.ranges; i += 1) {
            keys.push(this.roundKey(running.roundId, `range:${i}`));
        }
        for (let i = 0; i < keys.length; i += 1000) {
            await this.redis.del(...keys.slice(i, i + 1000));
        }
    }

    async writeRanges(roundId: number, bounds: RankingBounds[]): Promise<void> {
        const multi = this.redis.multi();
        const todo = this.roundKey(roundId, "todo:count");
        multi.del(todo);
        bounds.forEach((range, i) => {
            const key = this.roundKey(roundId, `range:${i}`);
            multi.del(key);
            multi.hset(key, { lo: JSON.stringify(range.lo), hi: JSON.stringify(range.hi) });
            multi.expire(key, ROUND_TTL_SEC);
            multi.rpush(todo, String(i));
        });
        multi.expire(todo, ROUND_TTL_SEC);
        await multi.exec();
    }

    async claim(roundId: number, step: RangeStep, leaseMs: number): Promise<number | null> {
        const index = (await this.redis.eval(
            CLAIM_SCRIPT,
            3,
            this.roundKey(roundId, `todo:${step}`),
            this.roundKey(roundId, `leased:${step}`),
            this.roundKey(roundId, `done:${step}`),
            String(Date.now()),
            String(leaseMs),
            String(ROUND_TTL_SEC)
        )) as string | null;
        return index === null ? null : Number(index);
    }

    async range(roundId: number, index: number): Promise<RangeTask> {
        const hash = await this.redis.hgetall(this.roundKey(roundId, `range:${index}`));
        return {
            index,
            bounds: { lo: JSON.parse(hash.lo ?? "null") as RankingKeys | null, hi: JSON.parse(hash.hi ?? "null") as RankingKeys | null },
            counts: JSON.parse(hash.counts ?? "{}") as Record<string, number>,
            offsets: JSON.parse(hash.offsets ?? "{}") as Record<string, number>,
        };
    }

    async finish(roundId: number, step: RangeStep, index: number, counts: Record<string, number> | null): Promise<void> {
        await this.redis.eval(
            FINISH_SCRIPT,
            3,
            this.roundKey(roundId, `leased:${step}`),
            this.roundKey(roundId, `done:${step}`),
            this.roundKey(roundId, `range:${index}`),
            String(index),
            counts === null ? "" : "counts",
            counts === null ? "" : JSON.stringify(counts),
            String(ROUND_TTL_SEC)
        );
    }

    async hasRanges(roundId: number): Promise<boolean> {
        return (await this.redis.exists(this.roundKey(roundId, "range:0"))) === 1;
    }

    async done(roundId: number, step: RangeStep): Promise<number> {
        return this.redis.scard(this.roundKey(roundId, `done:${step}`));
    }

    async allCounts(roundId: number, ranges: number): Promise<Array<Record<string, number>>> {
        const multi = this.redis.multi();
        for (let i = 0; i < ranges; i += 1) {
            multi.hget(this.roundKey(roundId, `range:${i}`), "counts");
        }
        const results = (await multi.exec()) ?? [];
        return results.map(([, value]) => JSON.parse((value as string | null) ?? "{}") as Record<string, number>);
    }

    async startRanking(roundId: number, offsets: Array<Record<string, number>>): Promise<void> {
        const multi = this.redis.multi();
        const todo = this.roundKey(roundId, "todo:rank");
        multi.del(todo);
        offsets.forEach((offset, i) => {
            const key = this.roundKey(roundId, `range:${i}`);
            multi.hset(key, "offsets", JSON.stringify(offset));
            multi.expire(key, ROUND_TTL_SEC);
            multi.rpush(todo, String(i));
        });
        multi.expire(todo, ROUND_TTL_SEC);
        await multi.exec();
    }

    async lock(kind: RankingKind): Promise<string | null> {
        const token = randomUUID();
        const ok = await this.redis.set(this.lockKey(kind), token, "EX", LOCK_TTL_SEC, "NX");
        return ok === "OK" ? token : null;
    }

    async renewLock(kind: RankingKind, token: string): Promise<void> {
        await this.redis.eval(RENEW_SCRIPT, 1, this.lockKey(kind), token, String(LOCK_TTL_SEC));
    }

    async unlock(kind: RankingKind, token: string): Promise<void> {
        await this.redis.eval(RELEASE_SCRIPT, 1, this.lockKey(kind), token);
    }

    static lockRenewMs(): number {
        return (LOCK_TTL_SEC * 1000) / 3;
    }

    async heartbeat(nodeId: string, workers: number): Promise<void> {
        const now = Date.now();
        await this.redis
            .multi()
            .zadd(NODES_KEY, now + NODE_TTL_MS, `${nodeId}:${workers}`)
            .zremrangebyscore(NODES_KEY, "-inf", now)
            .expire(NODES_KEY, NODES_KEY_TTL_SEC)
            .exec();
    }

    async liveWorkers(): Promise<number> {
        const members = await this.redis.zrangebyscore(NODES_KEY, Date.now(), "+inf");
        return members.reduce((sum, member) => sum + Number(member.slice(member.lastIndexOf(":") + 1)), 0);
    }
}
