"use strict";

const messages = require("../../protobuf/fminternal/internal_service_pb");

function plainObjectToBonusProto(o) {
    if (o == null || typeof o !== "object") {
        return null;
    }
    const b = new messages.EquipmentBonusStatsPersisted();
    let any = false;
    const setIfHas = (key, fn) => {
        if (!Object.prototype.hasOwnProperty.call(o, key)) {
            return;
        }
        any = true;
        fn(Number(o[key]) | 0);
    };
    setIfHas("str", (v) => b.setStr(v));
    setIfHas("dex", (v) => b.setDex(v));
    setIfHas("int", (v) => b.setIntStat(v));
    setIfHas("luk", (v) => b.setLuk(v));
    setIfHas("maxHp", (v) => b.setMaxHp(v));
    setIfHas("maxMp", (v) => b.setMaxMp(v));
    setIfHas("pad", (v) => b.setPad(v));
    setIfHas("mad", (v) => b.setMad(v));
    setIfHas("pdd", (v) => b.setPdd(v));
    setIfHas("mdd", (v) => b.setMdd(v));
    setIfHas("acc", (v) => b.setAcc(v));
    setIfHas("avoid", (v) => b.setAvoid(v));
    setIfHas("hands", (v) => b.setHands(v));
    setIfHas("speed", (v) => b.setSpeed(v));
    setIfHas("jump", (v) => b.setJump(v));
    return any ? b : null;
}

function bonusProtoToPlainObject(b) {
    if (b == null) {
        return null;
    }
    const o = {};
    if (b.getStr() !== 0) {
        o.str = b.getStr();
    }
    if (b.getDex() !== 0) {
        o.dex = b.getDex();
    }
    if (b.getIntStat() !== 0) {
        o.int = b.getIntStat();
    }
    if (b.getLuk() !== 0) {
        o.luk = b.getLuk();
    }
    if (b.getMaxHp() !== 0) {
        o.maxHp = b.getMaxHp();
    }
    if (b.getMaxMp() !== 0) {
        o.maxMp = b.getMaxMp();
    }
    if (b.getPad() !== 0) {
        o.pad = b.getPad();
    }
    if (b.getMad() !== 0) {
        o.mad = b.getMad();
    }
    if (b.getPdd() !== 0) {
        o.pdd = b.getPdd();
    }
    if (b.getMdd() !== 0) {
        o.mdd = b.getMdd();
    }
    if (b.getAcc() !== 0) {
        o.acc = b.getAcc();
    }
    if (b.getAvoid() !== 0) {
        o.avoid = b.getAvoid();
    }
    if (b.getHands() !== 0) {
        o.hands = b.getHands();
    }
    if (b.getSpeed() !== 0) {
        o.speed = b.getSpeed();
    }
    if (b.getJump() !== 0) {
        o.jump = b.getJump();
    }
    return Object.keys(o).length > 0 ? o : null;
}

function fillMessageFromPersisted(msg, model) {
    const uid = model.uniqueId;
    if (uid != null && uid !== "" && uid !== 0 && uid !== "0") {
        msg.setUniqueId(Number(uid));
    } else if (typeof msg.clearUniqueId === "function") {
        msg.clearUniqueId();
    }
    msg.setOwnerId(model.ownerId >>> 0);
    msg.setInventoryType((model.inventoryType ?? 0) >>> 0);
    msg.setItemId(model.itemId >>> 0);
    msg.setSlot(model.slot | 0);
    msg.setCount(model.count >>> 0);
    msg.setExpirationUnixMs(model.expiration ? model.expiration.getTime() : 0);
    msg.setEnchantChance((model.enchantChance ?? 0) >>> 0);
    msg.setFlag((model.flag ?? 0) >>> 0);
    msg.setSkillBonus((model.skillBonus ?? 0) >>> 0);
    msg.setOwnerName(model.ownerName ?? "");
    const bonus = plainObjectToBonusProto(model.equipBonusStats);
    if (bonus != null) {
        msg.setEquipBonusStats(bonus);
    } else if (typeof msg.clearEquipBonusStats === "function") {
        msg.clearEquipBonusStats();
    }
}

function persistedFromMessage(msg) {
    const expirationMs = msg.getExpirationUnixMs();
    const hasUid = typeof msg.hasUniqueId === "function" && msg.hasUniqueId();
    const slot = msg.getSlot() | 0;
    const itemId = msg.getItemId() >>> 0;
    return {
        uniqueId:      hasUid ? Number(msg.getUniqueId()) : null,
        ownerId:       msg.getOwnerId() >>> 0,
        inventoryType: msg.getInventoryType() >>> 0,
        itemId,
        slot,
        count:         msg.getCount() >>> 0,
        expiration:    expirationMs > 0 ? new Date(expirationMs) : null,
        enchantChance: msg.getEnchantChance() || null,
        flag:          msg.getFlag() || null,
        skillBonus:    msg.getSkillBonus() || null,
        ownerName:     msg.getOwnerName() || null,
        equipBonusStats: bonusProtoToPlainObject(msg.getEquipBonusStats()),
    };
}

function makeInventoryMessage(model) {
    const msg = new messages.InventoryPersisted();
    fillMessageFromPersisted(msg, model);
    return msg;
}

module.exports = { fillMessageFromPersisted, persistedFromMessage, makeInventoryMessage };
