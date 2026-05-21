export type RedisEvalOk = [1, number];
export type RedisEvalFail = [0, number];
export type RedisEvalResult = RedisEvalOk | RedisEvalFail;

export type RedisRefreshEvalOk = [1, number, number, number];
export type RedisRefreshEvalFail = [0, number, number, number];
export type RedisRefreshEvalResult = RedisRefreshEvalOk | RedisRefreshEvalFail;

export type RedisLogoutEvalFail = [0, number, number];
export type RedisLogoutEvalOk = [1, number, number];
export type RedisLogoutEvalResult = RedisLogoutEvalOk | RedisLogoutEvalFail;

export function isRedisEvalArray(raw: Array<number | string>): raw is RedisEvalResult {
    return raw.length >= 2 && (Number(raw[0]) === 0 || Number(raw[0]) === 1);
}

export function isRedisRefreshEvalArray(raw: Array<number | string>): raw is RedisRefreshEvalResult {
    return raw.length >= 4 && (Number(raw[0]) === 0 || Number(raw[0]) === 1);
}

export function isRedisLogoutEvalArray(raw: Array<number | string>): raw is RedisLogoutEvalResult {
    return raw.length >= 3 && (Number(raw[0]) === 0 || Number(raw[0]) === 1);
}
