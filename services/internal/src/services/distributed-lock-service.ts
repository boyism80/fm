import type { InternalContext } from "../context/internal-context";
import {
    DistributedLock,
    DistributedLockMultiGuard,
    redisLockShardHash,
    redisWorldLockKey,
    type DistributedLockGuard,
    type DistributedLockOptions,
} from "../system/distributed-lock";

export class DistributedLockService {
    private readonly ctx: InternalContext;
    private readonly lock: DistributedLock;

    constructor(internalContext: InternalContext, distributedLock: DistributedLock) {
        this.ctx = internalContext;
        this.lock = distributedLock;
    }

    acquireWorldDataLock(worldId: number, keyRest: string, options?: DistributedLockOptions): Promise<DistributedLockGuard> {
        const lockKey = redisWorldLockKey(worldId, keyRest);
        const hash = redisLockShardHash(lockKey);
        const { client } = this.ctx.getRedisDataAccess(worldId, hash);
        return this.lock.acquire(client, lockKey, options);
    }

    async acquireWorldDataLocks(worldId: number, keyRests: string[], options?: DistributedLockOptions): Promise<DistributedLockMultiGuard> {
        const sorted = [...new Set(keyRests)].sort();
        if (sorted.length === 0) {
            return new DistributedLockMultiGuard([]);
        }
        const guards: DistributedLockGuard[] = [];
        try {
            for (const keyRest of sorted) {
                guards.push(await this.acquireWorldDataLock(worldId, keyRest, options));
            }
            return new DistributedLockMultiGuard(guards);
        } catch (err) {
            for (let i = guards.length - 1; i >= 0; i--) {
                await guards[i]!.release();
            }
            throw err;
        }
    }
}
