"use strict";

function jsonStringToBindings(jsonStr) {
    if (jsonStr == null || jsonStr === "" || jsonStr === "{}") {
        return [];
    }
    let obj;
    try {
        obj = typeof jsonStr === "string" ? JSON.parse(jsonStr) : jsonStr;
    } catch {
        return [];
    }
    if (!obj || typeof obj !== "object") {
        return [];
    }
    const out = [];
    for (const [k, v] of Object.entries(obj)) {
        const slot = Number(k);
        if (!Number.isInteger(slot)) {
            continue;
        }
        out.push({
            slot,
            type: Number(v.type) >>> 0,
            action: Number(v.action) | 0,
        });
    }
    out.sort((a, b) => a.slot - b.slot);
    return out;
}

function bindingsToJsonString(bindings) {
    if (!bindings || bindings.length === 0) {
        return "{}";
    }
    const obj = {};
    for (const b of bindings) {
        obj[String(b.slot)] = {
            type: b.type >>> 0,
            action: b.action | 0,
        };
    }
    return JSON.stringify(obj);
}

function bindingsFromProtoList(list) {
    if (!list || list.length === 0) {
        return [];
    }
    return list.map((msg) => ({
        slot: msg.getSlot() | 0,
        type: msg.getType() >>> 0,
        action: msg.getAction() | 0,
    }));
}

function makeKeyLayoutProtoList(messages, bindings) {
    if (!bindings || bindings.length === 0) {
        return [];
    }
    return bindings.map((b) => {
        const m = new messages.KeyLayoutBinding();
        m.setSlot(b.slot | 0);
        m.setType(b.type >>> 0);
        m.setAction(b.action | 0);
        return m;
    });
}

module.exports = {
    jsonStringToBindings,
    bindingsToJsonString,
    bindingsFromProtoList,
    makeKeyLayoutProtoList,
};
