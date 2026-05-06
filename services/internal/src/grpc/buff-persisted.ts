import { BuffKind } from "../protobuf/generated/fminternal/internal_service";
import type { BuffPersisted } from "../protobuf/generated/fminternal/internal_service";
import type { BuffModel } from "../repos/buff-repository";
import { createMap, forMember, mapFrom } from "@automapper/core";
import { grpcMapper } from "./mappers";

export const BUFF_MODEL = "BuffModel";
export const BUFF_PERSISTED = "BuffPersisted";

createMap(
    grpcMapper,
    BUFF_PERSISTED,
    BUFF_MODEL,
    forMember((destination: any) => destination.characterId, mapFrom((source: BuffPersisted) => source.characterId >>> 0)),
    forMember((destination: any) => destination.buffSourceId, mapFrom((source: BuffPersisted) => source.buffSourceId | 0)),
    forMember((destination: any) => destination.kind, mapFrom((source: BuffPersisted) => source.kind)),
    forMember((destination: any) => destination.flagValues, mapFrom((source: BuffPersisted) => (source.flagValues ?? []).map((f) => ({
        mask: f.mask >>> 0,
        position: f.position | 0,
        value: f.value | 0,
    })))),
    forMember((destination: any) => destination.remainingDurationMs, mapFrom((source: BuffPersisted) => {
        const remaining = source.remainingDurationMs;
        return remaining !== undefined && Number.isFinite(remaining) && remaining > 0 ? Math.floor(remaining) : null;
    })),
    forMember((destination: any) => destination.skillLevel, mapFrom((source: BuffPersisted) => source.kind === BuffKind.BUFF_KIND_ITEM ? null : source.skillLevel >>> 0)),
    forMember((destination: any) => destination.causerId, mapFrom((source: BuffPersisted) => source.kind === BuffKind.BUFF_KIND_ITEM ? null : source.causerId >>> 0))
);

createMap(
    grpcMapper,
    BUFF_MODEL,
    BUFF_PERSISTED,
    forMember((destination: any) => destination.characterId, mapFrom((source: BuffModel) => source.characterId >>> 0)),
    forMember((destination: any) => destination.buffSourceId, mapFrom((source: BuffModel) => source.buffSourceId | 0)),
    forMember((destination: any) => destination.kind, mapFrom((source: BuffModel) => source.kind === BuffKind.BUFF_KIND_ITEM ? BuffKind.BUFF_KIND_ITEM : BuffKind.BUFF_KIND_SKILL)),
    forMember((destination: any) => destination.flagValues, mapFrom((source: BuffModel) => source.flagValues.map((f) => ({
        mask: f.mask >>> 0,
        position: f.position | 0,
        value: f.value | 0,
    })))),
    forMember((destination: any) => destination.skillLevel, mapFrom((source: BuffModel) => source.skillLevel ?? 0)),
    forMember((destination: any) => destination.causerId, mapFrom((source: BuffModel) => source.causerId ?? 0)),
    forMember((destination: any) => destination.remainingDurationMs, mapFrom((source: BuffModel) => {
        if (source.remainingDurationMs != null && source.remainingDurationMs > 0) {
            return Math.floor(source.remainingDurationMs);
        }
        return undefined;
    }))
);
