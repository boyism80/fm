export const REDIS_SESSION_KEY_PREFIX = "fm:session:";

export function redisSessionKey(rest: string): string {
    return `${REDIS_SESSION_KEY_PREFIX}${rest}`;
}
