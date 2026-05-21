import type { KeyLayoutBinding } from "../protobuf/generated/fminternal/internal_service";
import type { KeyLayoutJsonRecord, KeyLayoutSlotJson } from "../types/key-layout-json";

export type KeyLayoutBindingModel = { slot: number; type: number; action: number };

export function jsonStringToBindings(jsonStr: string): KeyLayoutBindingModel[] {
    if (jsonStr === "" || jsonStr === "{}") {
        return [];
    }
    let obj: KeyLayoutJsonRecord;
    try {
        obj = JSON.parse(jsonStr) as KeyLayoutJsonRecord;
    } catch {
        return [];
    }
    const out: KeyLayoutBindingModel[] = [];
    for (const [k, v] of Object.entries(obj)) {
        const slot = Number(k);
        if (!Number.isInteger(slot)) continue;
        const value: KeyLayoutSlotJson = v;
        out.push({ slot, type: Number(value.type) >>> 0, action: Number(value.action) | 0 });
    }
    out.sort((a, b) => a.slot - b.slot);
    return out;
}

export function bindingsToJsonString(bindings: KeyLayoutBindingModel[]) {
    if (!bindings || bindings.length === 0) {
        return "{}";
    }
    const obj: Record<string, { type: number; action: number }> = {};
    for (const b of bindings) {
        obj[String(b.slot)] = { type: b.type >>> 0, action: b.action | 0 };
    }
    return JSON.stringify(obj);
}

export function bindingsFromProtoList(list: KeyLayoutBinding[]): KeyLayoutBindingModel[] {
    if (!list || list.length === 0) {
        return [];
    }
    return list.map((msg) => ({ slot: msg.slot | 0, type: msg.type >>> 0, action: msg.action | 0 }));
}

export function makeKeyLayoutProtoList(bindings: KeyLayoutBindingModel[]): KeyLayoutBinding[] {
    if (!bindings || bindings.length === 0) {
        return [];
    }
    return bindings.map((b) => ({ slot: b.slot | 0, type: b.type >>> 0, action: b.action | 0 }));
}
