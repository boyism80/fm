import type { SavedLocation } from "../protobuf/generated/fminternal/internal_service";
import type { SavedLocationModel } from "../repos/saved-location-repository";
import { createMap, forMember, mapFrom } from "@automapper/core";
import { grpcMapper } from "./mappers";

export const SAVED_LOCATION_MODEL = "SavedLocationModel";
export const SAVED_LOCATION_PROTO = "SavedLocationProto";

createMap(
    grpcMapper,
    SAVED_LOCATION_MODEL,
    SAVED_LOCATION_PROTO,
    forMember((destination: any) => destination.characterId, mapFrom((source: SavedLocationModel) => source.characterId >>> 0)),
    forMember((destination: any) => destination.name, mapFrom((source: SavedLocationModel) => source.name ?? "")),
    forMember((destination: any) => destination.mapId, mapFrom((source: SavedLocationModel) => source.mapId >>> 0))
);

createMap(
    grpcMapper,
    SAVED_LOCATION_PROTO,
    SAVED_LOCATION_MODEL,
    forMember((destination: any) => destination.characterId, mapFrom((source: SavedLocation) => source.characterId >>> 0)),
    forMember((destination: any) => destination.name, mapFrom((source: SavedLocation) => source.name ?? "")),
    forMember((destination: any) => destination.mapId, mapFrom((source: SavedLocation) => source.mapId >>> 0))
);
