import type Redis from "ioredis";
import type { InternalContext } from "../context/internal-context";
import {
    DistributedLock,
    DistributedLockMultiGuard,
    redisLockKey,
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

    acquireWithClient(client: Redis, lockKey: string, options?: DistributedLockOptions): Promise<DistributedLockGuard> {
        return this.lock.acquire(client, lockKey, options);
    }

    tryAcquireWithClient(client: Redis, lockKey: string, options?: DistributedLockOptions): Promise<DistributedLockGuard | null> {
        return this.lock.tryAcquire(client, lockKey, options);
    }

    acquireWorldDataLock(worldId: number, keyRest: string, options?: DistributedLockOptions): Promise<DistributedLockGuard> {
        const lockKey = redisWorldLockKey(worldId, keyRest);
        const hash = redisLockShardHash(lockKey);
        const { client } = this.ctx.getRedisDataAccess(worldId, hash);
        return this.lock.acquire(client, lockKey, options);
    }

    tryAcquireWorldDataLock(worldId: number, keyRest: string, options?: DistributedLockOptions): Promise<DistributedLockGuard | null> {
        const lockKey = redisWorldLockKey(worldId, keyRest);
        const hash = redisLockShardHash(lockKey);
        const { client } = this.ctx.getRedisDataAccess(worldId, hash);
        return this.lock.tryAcquire(client, lockKey, options);
    }

    acquireWorldGlobalLock(worldId: number, keyRest: string, options?: DistributedLockOptions): Promise<DistributedLockGuard> {
        const lockKey = redisWorldLockKey(worldId, keyRest);
        const { client } = this.ctx.getRedisGlobalAccess(worldId);
        return this.lock.acquire(client, lockKey, options);
    }

    tryAcquireWorldGlobalLock(worldId: number, keyRest: string, options?: DistributedLockOptions): Promise<DistributedLockGuard | null> {
        const lockKey = redisWorldLockKey(worldId, keyRest);
        const { client } = this.ctx.getRedisGlobalAccess(worldId);
        return this.lock.tryAcquire(client, lockKey, options);
    }

    acquireWorldGlobalLocks(worldId: number, keyRests: string[], options?: DistributedLockOptions): Promise<DistributedLockMultiGuard> {
        if (keyRests.length === 0) {
            return Promise.resolve(new DistributedLockMultiGuard([]));
        }
        const lockKeys = keyRests.map((keyRest) => {
            return redisWorldLockKey(worldId, keyRest);
        });
        const { client } = this.ctx.getRedisGlobalAccess(worldId);
        return this.lock.acquireAll(client, lockKeys, options);
    }

    acquireUnifiedLock(keyRest: string, options?: DistributedLockOptions): Promise<DistributedLockGuard> {
        const lockKey = redisLockKey(keyRest);
        const { client } = this.ctx.getRedisUnifiedAccess();
        return this.lock.acquire(client, lockKey, options);
    }

    tryAcquireUnifiedLock(keyRest: string, options?: DistributedLockOptions): Promise<DistributedLockGuard | null> {
        const lockKey = redisLockKey(keyRest);
        const { client } = this.ctx.getRedisUnifiedAccess();
        return this.lock.tryAcquire(client, lockKey, options);
    }
}
