import type {
    EquipmentBonusStatsPersisted,
    InventoryPersisted,
} from "../protobuf/generated/fminternal/internal_service";
import type { InventoryModel } from "../repos/inventory-repository";
import { createMap, forMember, mapFrom } from "@automapper/core";
import { grpcMapper } from "./mappers";

export const INVENTORY_MODEL = "InventoryModel";
export const INVENTORY_PERSISTED = "InventoryPersisted";

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

createMap(
    grpcMapper,
    INVENTORY_MODEL,
    INVENTORY_PERSISTED,
    forMember((destination: any) => destination.uniqueId, mapFrom((source: InventoryModel) => {
        const uid = source.uniqueId;
        return uid != null && uid !== 0 ? uid : undefined;
    })),
    forMember((destination: any) => destination.ownerId, mapFrom((source: InventoryModel) => source.ownerId >>> 0)),
    forMember((destination: any) => destination.inventoryType, mapFrom((source: InventoryModel) => (source.inventoryType ?? 0) >>> 0)),
    forMember((destination: any) => destination.itemId, mapFrom((source: InventoryModel) => source.itemId >>> 0)),
    forMember((destination: any) => destination.slot, mapFrom((source: InventoryModel) => source.slot | 0)),
    forMember((destination: any) => destination.count, mapFrom((source: InventoryModel) => source.count >>> 0)),
    forMember((destination: any) => destination.expirationUnixMs, mapFrom((source: InventoryModel) => source.expiration ? source.expiration.getTime() : 0)),
    forMember((destination: any) => destination.enhanceChance, mapFrom((source: InventoryModel) => (source.enhanceChance ?? 0) >>> 0)),
    forMember((destination: any) => destination.enhanceCount, mapFrom((source: InventoryModel) => (source.enhanceCount ?? 0) >>> 0)),
    forMember((destination: any) => destination.flag, mapFrom((source: InventoryModel) => (source.flag ?? 0) >>> 0)),
    forMember((destination: any) => destination.skillBonus, mapFrom((source: InventoryModel) => (source.skillBonus ?? 0) >>> 0)),
    forMember((destination: any) => destination.ownerName, mapFrom((source: InventoryModel) => source.ownerName ?? "")),
    forMember((destination: any) => destination.equipBonusStats, mapFrom((source: InventoryModel) => plainObjectToBonusProto(source.equipBonusStats)))
);

createMap(
    grpcMapper,
    INVENTORY_PERSISTED,
    INVENTORY_MODEL,
    forMember((destination: any) => destination.uniqueId, mapFrom((source: InventoryPersisted) => source.uniqueId ?? null)),
    forMember((destination: any) => destination.ownerId, mapFrom((source: InventoryPersisted) => source.ownerId >>> 0)),
    forMember((destination: any) => destination.inventoryType, mapFrom((source: InventoryPersisted) => source.inventoryType >>> 0)),
    forMember((destination: any) => destination.itemId, mapFrom((source: InventoryPersisted) => source.itemId >>> 0)),
    forMember((destination: any) => destination.slot, mapFrom((source: InventoryPersisted) => source.slot | 0)),
    forMember((destination: any) => destination.count, mapFrom((source: InventoryPersisted) => source.count >>> 0)),
    forMember((destination: any) => destination.expiration, mapFrom((source: InventoryPersisted) => source.expirationUnixMs > 0 ? new Date(source.expirationUnixMs) : null)),
    forMember((destination: any) => destination.enhanceChance, mapFrom((source: InventoryPersisted) => source.enhanceChance || null)),
    forMember((destination: any) => destination.enhanceCount, mapFrom((source: InventoryPersisted) => source.enhanceCount || null)),
    forMember((destination: any) => destination.flag, mapFrom((source: InventoryPersisted) => source.flag || null)),
    forMember((destination: any) => destination.skillBonus, mapFrom((source: InventoryPersisted) => source.skillBonus || null)),
    forMember((destination: any) => destination.ownerName, mapFrom((source: InventoryPersisted) => source.ownerName || null)),
    forMember((destination: any) => destination.equipBonusStats, mapFrom((source: InventoryPersisted) => bonusProtoToPlainObject(source.equipBonusStats)))
);

