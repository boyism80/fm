import type {
    EquipmentBonusStatsPersisted,
    InventoryPersisted,
} from "../protobuf/generated/fminternal/internal_service";
import type { InventoryModel } from "../repos/inventory-repository";

function plainObjectToBonusProto(input: unknown): EquipmentBonusStatsPersisted | undefined {
    const o = input as Record<string, unknown>;
    if (o == null || typeof o !== "object") {
        return undefined;
    }
    const b: EquipmentBonusStatsPersisted = {
        str: 0, dex: 0, intStat: 0, luk: 0, maxHp: 0, maxMp: 0, pad: 0, mad: 0, pdd: 0, mdd: 0,
        acc: 0, avoid: 0, hands: 0, speed: 0, jump: 0,
    };
    let hasAny = false;
    const setIfHas = (key: string, setter: (v: number) => void) => {
        if (!Object.prototype.hasOwnProperty.call(o, key)) {
            return;
        }
        hasAny = true;
        setter(Number(o[key]) | 0);
    };
    setIfHas("str", (v) => { b.str = v; });
    setIfHas("dex", (v) => { b.dex = v; });
    setIfHas("int", (v) => { b.intStat = v; });
    setIfHas("luk", (v) => { b.luk = v; });
    setIfHas("maxHp", (v) => { b.maxHp = v; });
    setIfHas("maxMp", (v) => { b.maxMp = v; });
    setIfHas("pad", (v) => { b.pad = v; });
    setIfHas("mad", (v) => { b.mad = v; });
    setIfHas("pdd", (v) => { b.pdd = v; });
    setIfHas("mdd", (v) => { b.mdd = v; });
    setIfHas("acc", (v) => { b.acc = v; });
    setIfHas("avoid", (v) => { b.avoid = v; });
    setIfHas("hands", (v) => { b.hands = v; });
    setIfHas("speed", (v) => { b.speed = v; });
    setIfHas("jump", (v) => { b.jump = v; });
    return hasAny ? b : undefined;
}

function bonusProtoToPlainObject(b: EquipmentBonusStatsPersisted | undefined): Record<string, number> | undefined {
    if (b == null) {
        return undefined;
    }
    const o: Record<string, number> = {};
    if (b.str !== 0) o.str = b.str;
    if (b.dex !== 0) o.dex = b.dex;
    if (b.intStat !== 0) o.int = b.intStat;
    if (b.luk !== 0) o.luk = b.luk;
    if (b.maxHp !== 0) o.maxHp = b.maxHp;
    if (b.maxMp !== 0) o.maxMp = b.maxMp;
    if (b.pad !== 0) o.pad = b.pad;
    if (b.mad !== 0) o.mad = b.mad;
    if (b.pdd !== 0) o.pdd = b.pdd;
    if (b.mdd !== 0) o.mdd = b.mdd;
    if (b.acc !== 0) o.acc = b.acc;
    if (b.avoid !== 0) o.avoid = b.avoid;
    if (b.hands !== 0) o.hands = b.hands;
    if (b.speed !== 0) o.speed = b.speed;
    if (b.jump !== 0) o.jump = b.jump;
    return Object.keys(o).length > 0 ? o : undefined;
}

export function fillMessageFromPersisted(model: InventoryModel): InventoryPersisted {
    const uid = model.uniqueId;
    return {
        uniqueId: uid != null && uid !== 0 ? uid : undefined,
        ownerId: model.ownerId >>> 0,
        inventoryType: (model.inventoryType ?? 0) >>> 0,
        itemId: model.itemId >>> 0,
        slot: model.slot | 0,
        count: model.count >>> 0,
        expirationUnixMs: model.expiration ? model.expiration.getTime() : 0,
        enhanceChance: (model.enhanceChance ?? 0) >>> 0,
        enhanceCount: (model.enhanceCount ?? 0) >>> 0,
        flag: (model.flag ?? 0) >>> 0,
        skillBonus: (model.skillBonus ?? 0) >>> 0,
        ownerName: model.ownerName ?? "",
        equipBonusStats: plainObjectToBonusProto(model.equipBonusStats),
    };
}

export { fillMessageFromPersisted as makeInventoryMessage };

export function persistedFromMessage(msg: InventoryPersisted) {
    const expirationMs = msg.expirationUnixMs;
    return {
        uniqueId: msg.uniqueId ?? null,
        ownerId: msg.ownerId >>> 0,
        inventoryType: msg.inventoryType >>> 0,
        itemId: msg.itemId >>> 0,
        slot: msg.slot | 0,
        count: msg.count >>> 0,
        expiration: expirationMs > 0 ? new Date(expirationMs) : null,
        enhanceChance: msg.enhanceChance || null,
        enhanceCount: msg.enhanceCount || null,
        flag: msg.flag || null,
        skillBonus: msg.skillBonus || null,
        ownerName: msg.ownerName || null,
        equipBonusStats: bonusProtoToPlainObject(msg.equipBonusStats),
    };
}

