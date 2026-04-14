"use strict";

const messages = require("../../protobuf/fminternal/ping_pb");

function fillMessageFromPersisted(msg, model) {
    msg.setCharacterId(model.characterId >>> 0);
    msg.setSkillId(model.skillId >>> 0);
    msg.setLevel(model.level | 0);
    msg.setMasterLevel(model.masterLevel | 0);
    msg.setCooldownEndUnixMs(model.cooldownEndUnixMs ?? 0);
    msg.setUpdatedAtUnixMs(model.updatedAt ? model.updatedAt.getTime() : 0);
}

function persistedFromMessage(msg) {
    const cooldownMs = msg.getCooldownEndUnixMs();
    return {
        characterId:       msg.getCharacterId() >>> 0,
        skillId:           msg.getSkillId() >>> 0,
        level:             msg.getLevel() | 0,
        masterLevel:       msg.getMasterLevel() | 0,
        cooldownEndUnixMs: cooldownMs > 0 ? cooldownMs : null,
        updatedAt:         new Date(msg.getUpdatedAtUnixMs()),
    };
}

function makeSkillMessage(model) {
    const msg = new messages.SkillPersisted();
    fillMessageFromPersisted(msg, model);
    return msg;
}

module.exports = { fillMessageFromPersisted, persistedFromMessage, makeSkillMessage };
