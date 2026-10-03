export const REDIS_NAME_KEY_PREFIX = "fm:name:";

export function redisNameKey(rest: string): string {
    return `${REDIS_NAME_KEY_PREFIX}${rest}`;
}
