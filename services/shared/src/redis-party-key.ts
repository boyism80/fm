export const REDIS_PARTY_KEY_PREFIX = "fm:party:";

export function redisPartyKey(rest: string): string {
    return `${REDIS_PARTY_KEY_PREFIX}${rest}`;
}
