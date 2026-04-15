"use strict";

const { persistedFromMessage: inventoryFromMessage, makeInventoryMessage } = require("../inventory-persisted");

function createInventoryHandlers(
    inventoryService,
    messages,
    grpcError
) {
    return {
        async getInventory(call, callback) {
            try {
                const items = await inventoryService.getInventory(
                    call.request.getWorldId(),
                    call.request.getOwnerId()
                );
                const reply = new messages.GetInventoryReply();
                reply.setItemsList(items.map(makeInventoryMessage));
                callback(null, reply);
            } catch (err) {
                grpcError(err, callback);
            }
        },

        async saveInventory(call, callback) {
            try {
                const worldId = call.request.getWorldId();
                const msgItems = call.request.getItemsList();
                if (msgItems.length) {
                    await inventoryService.saveInventory(worldId, msgItems.map(inventoryFromMessage));
                }
                const reply = new messages.SaveInventoryReply();
                reply.setOk(true);
                callback(null, reply);
            } catch (err) {
                grpcError(err, callback);
            }
        },

        async deleteInventory(call, callback) {
            try {
                await inventoryService.deleteInventory(
                    call.request.getWorldId(),
                    call.request.getOwnerId(),
                    call.request.getUniqueIdsList().map(String)
                );
                const reply = new messages.DeleteInventoryReply();
                reply.setOk(true);
                callback(null, reply);
            } catch (err) {
                grpcError(err, callback);
            }
        },
    };
}

module.exports = { createInventoryHandlers };
