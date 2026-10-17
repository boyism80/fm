export type RedisEvalOk = [1, number];
export type RedisEvalFail = [0, number];
export type RedisEvalResult = RedisEvalOk | RedisEvalFail;

export function isRedisEvalArray(raw: Array<number | string>): raw is RedisEvalResult {
    return raw.length >= 2 && (Number(raw[0]) === 0 || Number(raw[0]) === 1);
}
