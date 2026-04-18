"use strict";

/**
 * ASP.NET IConfiguration-style access: inject this instead of raw `loadConfig()` result
 * when services need typed getters. Root document remains on `.raw`.
 */
class AppConfiguration {
    /** @param {object} internalConfig normalized root from loadConfig() */
    constructor(internalConfig) {
        if (!internalConfig || typeof internalConfig !== "object") {
            throw new Error("AppConfiguration requires internalConfig from loadConfig()");
        }
        this._c = internalConfig;
    }

    /** Full normalized root (includes configPath, app, grpc, postgresql, redis, rabbitmq, cache, sequelize). */
    get raw() {
        return this._c;
    }

    get configPath() {
        return this._c.configPath;
    }

    get app() {
        return this._c.app;
    }

    get grpc() {
        return this._c.grpc;
    }

    get postgresql() {
        return this._c.postgresql;
    }

    get redis() {
        return this._c.redis;
    }

    get rabbitmq() {
        return this._c.rabbitmq;
    }

    get cache() {
        return this._c.cache;
    }

    get sequelize() {
        return this._c.sequelize;
    }

    getCharacterCacheTtlSeconds() {
        const v = this._c.cache?.character_ttl_seconds;
        const n = Number(v);
        return Number.isFinite(n) && n > 0 ? n : 300;
    }

    getItemCacheTtlSeconds() {
        const v = this._c.cache?.item_ttl_seconds;
        const n = Number(v);
        return Number.isFinite(n) && n > 0 ? n : 300;
    }

    /**
     * @param {"app"|"grpc"|"postgresql"|"redis"|"rabbitmq"|"cache"|"sequelize"} section
     */
    getSection(section) {
        return this._c[section];
    }

}

module.exports = { AppConfiguration };
