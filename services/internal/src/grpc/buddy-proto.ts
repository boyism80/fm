import type { BuddyEntry } from "../protobuf/generated/fminternal/internal_service";
import type { BuddyListEntry } from "../services/buddy-service";
import { createMap, forMember, mapFrom } from "@automapper/core";
import { grpcMapper } from "./mappers";

export const BUDDY_LIST_ENTRY = "BuddyListEntry";
export const BUDDY_ENTRY = "BuddyEntry";

createMap(
    grpcMapper,
    BUDDY_LIST_ENTRY,
    BUDDY_ENTRY,
    forMember((destination: any) => destination.characterId, mapFrom((source: BuddyListEntry) => source.characterId >>> 0)),
    forMember((destination: any) => destination.name, mapFrom((source: BuddyListEntry) => source.name ?? "")),
    forMember((destination: any) => destination.groupName, mapFrom((source: BuddyListEntry) => source.groupName ?? "")),
    forMember((destination: any) => destination.pending, mapFrom((source: BuddyListEntry) => source.pending)),
    forMember((destination: any) => destination.channelIndex, mapFrom((source: BuddyListEntry) => source.channelIndex | 0))
);

createMap(
    grpcMapper,
    BUDDY_ENTRY,
    BUDDY_LIST_ENTRY,
    forMember((destination: any) => destination.characterId, mapFrom((source: BuddyEntry) => source.characterId >>> 0)),
    forMember((destination: any) => destination.name, mapFrom((source: BuddyEntry) => source.name ?? "")),
    forMember((destination: any) => destination.groupName, mapFrom((source: BuddyEntry) => source.groupName ?? "")),
    forMember((destination: any) => destination.pending, mapFrom((source: BuddyEntry) => source.pending)),
    forMember((destination: any) => destination.channelIndex, mapFrom((source: BuddyEntry) => source.channelIndex | 0))
);
