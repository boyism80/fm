"use strict";

/** @param {import("../../protobuf/fminternal/ping_pb.js").CharacterPersisted} msg */
function fillMessageFromPersisted(msg, p) {
    msg.setCharacterId(p.characterId >>> 0);
    msg.setWorldId(p.worldId >>> 0);
    msg.setName(p.name ?? "");
    msg.setGender(p.gender >>> 0);
    msg.setSkinColor(p.skinColor >>> 0);
    msg.setFace(p.face >>> 0);
    msg.setHair(p.hair >>> 0);
    msg.setLevel(p.level >>> 0);
    msg.setClassId(p.classId >>> 0);
    msg.setRole((p.role ?? 0) >>> 0);
    msg.setStr(p.str >>> 0);
    msg.setDex(p.dex >>> 0);
    msg.setIntStat(p.intStat >>> 0);
    msg.setLuk(p.luk >>> 0);
    msg.setHp(p.hp >>> 0);
    msg.setMaxHp(p.maxHp >>> 0);
    msg.setMp(p.mp >>> 0);
    msg.setMaxMp(p.maxMp >>> 0);
    msg.setAbilityPoint(p.abilityPoint >>> 0);
    msg.setExp(p.exp >>> 0);
    msg.setMapId(p.mapId >>> 0);
    msg.setSpawnPoint(p.spawnPoint >>> 0);
    msg.setPositionX(p.positionX | 0);
    msg.setPositionY(p.positionY | 0);
    msg.setStance(p.stance >>> 0);
    msg.setMeso(p.meso | 0);
    msg.setSkillPoint(p.skillPoint >>> 0);
    msg.setUpdatedAtUnixMs(p.updatedAt ? new Date(p.updatedAt).getTime() : 0);
    msg.setAccountId(p.accountId >>> 0);
}

/** @param {import("../../protobuf/fminternal/ping_pb.js").CharacterPersisted} msg */
function persistedFromMessage(msg) {
    return {
        characterId: msg.getCharacterId() >>> 0,
        worldId: msg.getWorldId() >>> 0,
        name: msg.getName(),
        gender: msg.getGender() >>> 0,
        skinColor: msg.getSkinColor() >>> 0,
        face: msg.getFace() >>> 0,
        hair: msg.getHair() >>> 0,
        level: msg.getLevel() >>> 0,
        classId: msg.getClassId() >>> 0,
        role: msg.getRole() >>> 0,
        str: msg.getStr() >>> 0,
        dex: msg.getDex() >>> 0,
        intStat: msg.getIntStat() >>> 0,
        luk: msg.getLuk() >>> 0,
        hp: msg.getHp() >>> 0,
        maxHp: msg.getMaxHp() >>> 0,
        mp: msg.getMp() >>> 0,
        maxMp: msg.getMaxMp() >>> 0,
        abilityPoint: msg.getAbilityPoint() >>> 0,
        exp: msg.getExp() >>> 0,
        mapId: msg.getMapId() >>> 0,
        spawnPoint: msg.getSpawnPoint() >>> 0,
        positionX: msg.getPositionX(),
        positionY: msg.getPositionY(),
        stance: msg.getStance() >>> 0,
        meso: msg.getMeso(),
        skillPoint: msg.getSkillPoint() >>> 0,
        accountId: msg.getAccountId() >>> 0,
    };
}

module.exports = { fillMessageFromPersisted, persistedFromMessage };
