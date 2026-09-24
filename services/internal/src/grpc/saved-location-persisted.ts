import type { SavedLocationPersisted } from "../protobuf/generated/fminternal/internal_service";
import type { SavedLocationModel } from "../repos/saved-location-repository";
import { createMap, forMember, mapFrom } from "@automapper/core";
import { grpcMapper } from "./mappers";

export const SAVED_LOCATION_MODEL = "SavedLocationModel";
export const SAVED_LOCATION_PERSISTED = "SavedLocationPersisted";

createMap(
    grpcMapper,
    SAVED_LOCATION_MODEL,
    SAVED_LOCATION_PERSISTED,
    forMember((destination: any) => destination.characterId, mapFrom((source: SavedLocationModel) => source.characterId >>> 0)),
    forMember((destination: any) => destination.name, mapFrom((source: SavedLocationModel) => source.name ?? "")),
    forMember((destination: any) => destination.mapId, mapFrom((source: SavedLocationModel) => source.mapId >>> 0))
);

createMap(
    grpcMapper,
    SAVED_LOCATION_PERSISTED,
    SAVED_LOCATION_MODEL,
    forMember((destination: any) => destination.characterId, mapFrom((source: SavedLocationPersisted) => source.characterId >>> 0)),
    forMember((destination: any) => destination.name, mapFrom((source: SavedLocationPersisted) => source.name ?? "")),
    forMember((destination: any) => destination.mapId, mapFrom((source: SavedLocationPersisted) => source.mapId >>> 0))
);
