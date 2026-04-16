"use strict";

const messages = require("../../protobuf/fminternal/internal_service_pb");

function fillMessageFromPersisted(msg, model) {
    msg.setUniqueId(model.uniqueId);
    msg.setOwnerId(model.ownerId >>> 0);
    msg.setItemId(model.itemId >>> 0);
    msg.setSlot(model.slot | 0);
    msg.setCount(model.count >>> 0);
    msg.setExpirationUnixMs(model.expiration ? model.expiration.getTime() : 0);
    msg.setEnchantChance((model.enchantChance ?? 0) >>> 0);
    msg.setFlag((model.flag ?? 0) >>> 0);
    msg.setSkillBonus((model.skillBonus ?? 0) >>> 0);
    msg.setOwnerName(model.ownerName ?? "");
}

function persistedFromMessage(msg) {
    const expirationMs = msg.getExpirationUnixMs();
    return {
        uniqueId:      String(msg.getUniqueId()),
        ownerId:       msg.getOwnerId() >>> 0,
        itemId:        msg.getItemId() >>> 0,
        slot:          msg.getSlot() | 0,
        count:         msg.getCount() >>> 0,
        expiration:    expirationMs > 0 ? new Date(expirationMs) : null,
        enchantChance: msg.getEnchantChance() || null,
        flag:          msg.getFlag() || null,
        skillBonus:    msg.getSkillBonus() || null,
        ownerName:     msg.getOwnerName() || null,
    };
}

function makeInventoryMessage(model) {
    const msg = new messages.InventoryPersisted();
    fillMessageFromPersisted(msg, model);
    return msg;
}

module.exports = { fillMessageFromPersisted, persistedFromMessage, makeInventoryMessage };
