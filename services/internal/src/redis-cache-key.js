"use strict";

/** Fixed Redis key namespace for FM cache keys (not configurable). */
const REDIS_CACHE_KEY_PREFIX = "fm:cache:";

/** @param {string} rest Path after `fm:cache:`, e.g. `w0:character:1` */
function redisCacheKey(rest) {
    return `${REDIS_CACHE_KEY_PREFIX}${rest}`;
}

module.exports = { REDIS_CACHE_KEY_PREFIX, redisCacheKey };
