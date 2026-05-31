import { randomUUID } from "node:crypto";
import type Redis from "ioredis";

export const REDIS_LOCK_KEY_PREFIX = "fm:lock:";

const FNV_OFFSET_BASIS = 0xcbf29ce484222325n;
const FNV_PRIME = 0x100000001b3n;
const MASK64 = 0xffffffffffffffffn;

const RELEASE_SCRIPT = `
if redis.call("GET", KEYS[1]) == ARGV[1] then
  return redis.call("DEL", KEYS[1])
end
return 0
`;

export type DistributedLockOptions = {
    leaseTtlSec?: number;
    waitTimeoutMs?: number;
    retryIntervalMs?: number;
    retryBackoffFactor?: number;
    maxRetryIntervalMs?: number;
};

const DEFAULT_LOCK_LEASE_TTL_SEC = 30;
const DEFAULT_LOCK_WAIT_TIMEOUT_MS = 5000;
const DEFAULT_LOCK_RETRY_INTERVAL_MS = 50;
const DEFAULT_LOCK_RETRY_BACKOFF_FACTOR = 2;
const DEFAULT_LOCK_MAX_RETRY_INTERVAL_MS = 500;

export class DistributedLockAcquireError extends Error {
    readonly lockKey: string;

    constructor(lockKey: string) {
        super(`Failed to acquire distributed lock: ${lockKey}`);
        this.name = "DistributedLockAcquireError";
        this.lockKey = lockKey;
    }
}

export function redisLockKey(rest: string): string {
    return `${REDIS_LOCK_KEY_PREFIX}${rest}`;
}

export function redisWorldLockKey(worldId: number, rest: string): string {
    return redisLockKey(`w${worldId}:${rest}`);
}

export function redisLockShardHash(lockKey: string): number {
    let hash = FNV_OFFSET_BASIS;
    for (let i = 0; i < lockKey.length; i++) {
        hash ^= BigInt(lockKey.charCodeAt(i));
        hash = (hash * FNV_PRIME) & MASK64;
    }
    return Number(hash & 0xffffffffn);
}

export class DistributedLockGuard implements AsyncDisposable {
    private readonly client: Redis;
    private readonly lockKey: string;
    private readonly token: string;
    private released = false;

    constructor(client: Redis, lockKey: string, token: string) {
        this.client = client;
        this.lockKey = lockKey;
        this.token = token;
    }

    async release(): Promise<void> {
        if (this.released) {
            return;
        }
        this.released = true;
        await this.client.eval(RELEASE_SCRIPT, 1, this.lockKey, this.token).catch(() => {});
    }

    async [Symbol.asyncDispose](): Promise<void> {
        await this.release();
    }
}

export class DistributedLockMultiGuard implements AsyncDisposable {
    private readonly guards: DistributedLockGuard[];

    constructor(guards: DistributedLockGuard[]) {
        this.guards = guards;
    }

    async release(): Promise<void> {
        for (let i = this.guards.length - 1; i >= 0; i--) {
            await this.guards[i]!.release();
        }
    }

    async [Symbol.asyncDispose](): Promise<void> {
        await this.release();
    }
}

export class DistributedLock {
    async tryAcquire(client: Redis, lockKey: string, options: DistributedLockOptions = {}): Promise<DistributedLockGuard | null> {
        const leaseTtlSec = options.leaseTtlSec ?? DEFAULT_LOCK_LEASE_TTL_SEC;
        const token = randomUUID();
        const ok = await client.set(lockKey, token, "EX", leaseTtlSec, "NX");
        if (ok !== "OK") {
            return null;
        }
        return new DistributedLockGuard(client, lockKey, token);
    }

    async acquire(client: Redis, lockKey: string, options: DistributedLockOptions = {}): Promise<DistributedLockGuard> {
        if (options.waitTimeoutMs === 0) {
            const guard = await this.tryAcquire(client, lockKey, options);
            if (!guard) {
                throw new DistributedLockAcquireError(lockKey);
            }
            return guard;
        }

        const waitTimeoutMs = options.waitTimeoutMs ?? DEFAULT_LOCK_WAIT_TIMEOUT_MS;
        const deadline = Date.now() + waitTimeoutMs;
        let retryIntervalMs = options.retryIntervalMs ?? DEFAULT_LOCK_RETRY_INTERVAL_MS;
        const retryBackoffFactor = options.retryBackoffFactor ?? DEFAULT_LOCK_RETRY_BACKOFF_FACTOR;
        const maxRetryIntervalMs = options.maxRetryIntervalMs ?? DEFAULT_LOCK_MAX_RETRY_INTERVAL_MS;

        while (true) {
            const guard = await this.tryAcquire(client, lockKey, options);
            if (guard) {
                return guard;
            }
            if (Date.now() >= deadline) {
                throw new DistributedLockAcquireError(lockKey);
            }
            await DistributedLock.sleep(retryIntervalMs);
            retryIntervalMs = Math.min(Math.trunc(retryIntervalMs * retryBackoffFactor), maxRetryIntervalMs);
        }
    }

    async acquireAll(client: Redis, lockKeys: string[], options: DistributedLockOptions = {}): Promise<DistributedLockMultiGuard> {
        const uniqueSorted = [...new Set(lockKeys)].sort();
        if (uniqueSorted.length === 0) {
            return new DistributedLockMultiGuard([]);
        }

        const guards: DistributedLockGuard[] = [];
        try {
            await Promise.all(
                uniqueSorted.map(async (lockKey) => {
                    const guard = await this.acquire(client, lockKey, options);
                    guards.push(guard);
                })
            );
            return new DistributedLockMultiGuard(guards);
        } catch (err) {
            for (let i = guards.length - 1; i >= 0; i--) {
                await guards[i]!.release();
            }
            throw err;
        }
    }

    private static sleep(ms: number): Promise<void> {
        return new Promise((resolve) => {
            setTimeout(resolve, ms);
        });
    }
}
