export const REDIS_CACHE_KEY_PREFIX = "fm:cache:";

export function redisCacheKey(rest: string): string {
    return `${REDIS_CACHE_KEY_PREFIX}${rest}`;
}
