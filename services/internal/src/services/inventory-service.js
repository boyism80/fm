"use strict";

class InventoryService {
    constructor(inventoryRepository, appConfiguration) {
        this.repo = inventoryRepository;
        this.app = appConfiguration;
    }

    _assertWorld(worldId) {
        const wid = String(worldId);
        if (!this.app.postgresql.worlds[wid]) {
            const err = new Error(`Unknown world_id: ${worldId}`);
            err.code = "UNKNOWN_WORLD";
            throw err;
        }
    }

    _assertOwnerId(ownerId) {
        const n = Number(ownerId);
        if (!Number.isInteger(n) || n <= 0 || n > 0xffffffff) {
            const err = new Error("owner_id must be a positive uint32");
            err.code = "INVALID_OWNER_ID";
            throw err;
        }
    }

    _assertUniqueId(uniqueId) {
        const s = String(uniqueId);
        const n = BigInt(s);
        if (n <= 0n) {
            const err = new Error("unique_id must be a positive int64");
            err.code = "INVALID_UNIQUE_ID";
            throw err;
        }
    }

    _validateItem(item) {
        if (item.count == null || item.count < 1) {
            const err = new Error("count must be >= 1");
            err.code = "INVALID_PAYLOAD";
            throw err;
        }
        if (item.itemId == null || item.itemId <= 0) {
            const err = new Error("item_id must be a positive integer");
            err.code = "INVALID_PAYLOAD";
            throw err;
        }
    }

    async getInventory(worldId, ownerId) {
        this._assertWorld(worldId);
        this._assertOwnerId(ownerId);
        const map = await this.repo.getInventory(worldId, ownerId);
        return [...map.values()];
    }

    async saveInventory(worldId, items) {
        if (!items.length) return [];
        this._assertWorld(worldId);
        const ownerId = items[0].ownerId;
        this._assertOwnerId(ownerId);
        for (const item of items) {
            if (Number(item.ownerId) !== Number(ownerId)) {
                const err = new Error("all items in a single SaveInventory call must share the same owner_id");
                err.code = "INVALID_PAYLOAD";
                throw err;
            }
            this._assertUniqueId(item.uniqueId);
            this._validateItem(item);
        }
        return this.repo.saveAll(worldId, items);
    }

    async deleteInventory(worldId, ownerId, uniqueIds) {
        if (!uniqueIds.length) return;
        this._assertWorld(worldId);
        this._assertOwnerId(ownerId);
        for (const uid of uniqueIds) {
            this._assertUniqueId(uid);
        }
        return this.repo.deleteAll(worldId, ownerId, uniqueIds);
    }
}

module.exports = { InventoryService };
