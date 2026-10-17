export const REDIS_ALIVE_KEY_PREFIX = "fm:alive:";

export function redisAliveKey(rest: string): string {
    return `${REDIS_ALIVE_KEY_PREFIX}${rest}`;
}
