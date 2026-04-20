"use strict";

const { Pool } = require("pg");
const Redis = require("ioredis");
const {
    getPostgresqlWorld,
    getRedisWorld,
} = require("../config");

function shardIndex(hash, shardCount) {
    if (shardCount <= 1) {
        return 0;
    }
    const h = Number(hash);
    if (!Number.isFinite(h)) {
        return 0;
    }
    return Math.abs(Math.trunc(h)) % shardCount;
}

/**
 * Holds pg Pools and ioredis clients per world (global + data shards).
 * Character traffic uses data shards only; global is available for future use.
 */
class InternalContext {
    /**
     * @param {import("../config/app-configuration").AppConfiguration} appConfiguration
     */
    constructor(appConfiguration) {
        const ic = appConfiguration.raw;
        /** @type {import("../config/app-configuration").AppConfiguration} */
        this.appConfiguration = appConfiguration;
        this.configPath = ic.configPath;
        this.config = {
            app: ic.app,
            grpc: ic.grpc,
            postgresql: ic.postgresql,
            redis: ic.redis,
            cache: ic.cache,
            sequelize: ic.sequelize,
        };

        /** @type {import('pg').Pool | null} */
        this._pgUnified = this.config.postgresql.unified
            ? InternalContext._createPool(this.config.postgresql.unified)
            : null;
        /** @type {Record<string, import('pg').Pool>} */
        this._pgGlobal = {};
        /** @type {Record<string, import('pg').Pool[]>} */
        this._pgData = {};
        /** @type {Record<string, Redis[]>} */
        this._redisGlobal = {};
        /** @type {Record<string, Redis[]>} */
        this._redisData = {};

        for (const wid of Object.keys(this.config.postgresql.worlds)) {
            const worldPg = getPostgresqlWorld(this.config, wid);
            this._pgGlobal[wid] = InternalContext._createPool(worldPg.global);
            this._pgData[wid] = worldPg.data.map((ep) => InternalContext._createPool(ep));
        }

        for (const wid of Object.keys(this.config.redis.worlds)) {
            const worldRedis = getRedisWorld(this.config, wid);
            this._redisGlobal[wid] = [InternalContext._createRedis(worldRedis.global)];
            this._redisData[wid] = worldRedis.data.map((ep) => InternalContext._createRedis(ep));
        }
    }

    static _createPool(ep) {
        return new Pool({
            host: ep.host,
            port: ep.port,
            database: ep.database,
            user: ep.username,
            password: ep.password,
            ssl: ep.ssl ? { rejectUnauthorized: false } : false,
            max: ep.pool.max,
            min: ep.pool.min,
        });
    }

    static _createRedis(ep) {
        const opts = {
            host: ep.host,
            port: ep.port,
            password: ep.password || undefined,
            db: ep.db,
            maxRetriesPerRequest: 2,
        };
        if (ep.tls) {
            opts.tls = {};
        }
        return new Redis(opts);
    }

    /**
     * Data shard pool for entity routing (e.g. character id).
     * @param {number|string} worldId
     * @param {number|string} hash
     */
    getPgDataPool(worldId, hash) {
        const wid = String(worldId);
        const worldPg = getPostgresqlWorld(this.config, wid);
        const idx = shardIndex(hash, worldPg.data.length);
        return this._pgData[wid][idx];
    }

    async withPgDataTransaction(worldId, hash, fn) {
        const pool = this.getPgDataPool(worldId, hash);
        const client = await pool.connect();
        try {
            await client.query("BEGIN");
            const result = await fn(client);
            await client.query("COMMIT");
            return result;
        } catch (err) {
            try {
                await client.query("ROLLBACK");
            } catch {
            }
            throw err;
        } finally {
            client.release();
        }
    }

    async withPgGlobalTransaction(worldId, fn) {
        const wid = String(worldId);
        const pool = this._pgGlobal[wid];
        if (!pool) {
            throw new Error(`Unknown world_id for global transaction: ${worldId}`);
        }
        const client = await pool.connect();
        try {
            await client.query("BEGIN");
            const result = await fn(client);
            await client.query("COMMIT");
            return result;
        } catch (err) {
            try {
                await client.query("ROLLBACK");
            } catch {
            }
            throw err;
        } finally {
            client.release();
        }
    }

    getPgUnifiedPool() {
        if (!this._pgUnified) {
            throw new Error("Unified PostgreSQL pool is not configured (postgresql.unified missing)");
        }
        return this._pgUnified;
    }

    /**
     * Returns the ioredis client for the data shard that owns `hash`.
     * @param {number|string} worldId
     * @param {number|string} hash  – used only for shard routing (e.g. character id)
     * @returns {{ client: import('ioredis') }}
     */
    getRedisDataAccess(worldId, hash) {
        const wid = String(worldId);
        const worldRedis = getRedisWorld(this.config, wid);
        const idx = shardIndex(hash, worldRedis.data.length);
        const client = this._redisData[wid][idx];
        return { client };
    }

    /**
     * Global redis is used for account/session coordination keys.
     * @param {number|string} worldId
     * @returns {{ client: import('ioredis') }}
     */
    getRedisGlobalAccess(worldId) {
        const wid = String(worldId);
        const worldRedis = getRedisWorld(this.config, wid);
        const client = this._redisGlobal[wid][0];
        return { client };
    }

    async close() {
        const pools = [];
        if (this._pgUnified) {
            pools.push(this._pgUnified);
        }
        for (const w of Object.keys(this._pgGlobal)) {
            pools.push(this._pgGlobal[w]);
            for (const p of this._pgData[w]) {
                pools.push(p);
            }
        }
        await Promise.all(
            pools.map((p) =>
                p
                    .end()
                    .catch(() => {})
            )
        );

        const clients = [];
        for (const w of Object.keys(this._redisData)) {
            for (const c of this._redisGlobal[w]) {
                clients.push(c);
            }
            for (const c of this._redisData[w]) {
                clients.push(c);
            }
        }
        for (const c of clients) {
            try {
                c.disconnect();
            } catch {
                /* ignore */
            }
        }
    }
}

module.exports = { InternalContext };
