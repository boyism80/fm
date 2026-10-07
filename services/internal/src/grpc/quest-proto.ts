import type { Quest } from "../protobuf/generated/fminternal/internal_service";
import type { QuestModel } from "../repos/quest-repository";
import { createMap, forMember, mapFrom } from "@automapper/core";
import { grpcMapper } from "./mappers";

export const QUEST_MODEL = "QuestModel";
export const QUEST_PROTO = "QuestProto";

function normalizeNumberMap(input: Record<number, number> | Record<string, number> | undefined): Record<string, number> {
    const out: Record<string, number> = {};
    for (const [key, value] of Object.entries(input ?? {})) {
        out[key] = Number(value) >>> 0;
    }
    return out;
}

function normalizeStringMap(input: Record<string, string> | undefined): Record<string, string> {
    const out: Record<string, string> = {};
    for (const [key, value] of Object.entries(input ?? {})) {
        out[key] = String(value);
    }
    return out;
}

createMap(
    grpcMapper,
    QUEST_MODEL,
    QUEST_PROTO,
    forMember((destination: any) => destination.characterId, mapFrom((source: QuestModel) => source.characterId >>> 0)),
    forMember((destination: any) => destination.questId, mapFrom((source: QuestModel) => source.questId >>> 0)),
    forMember((destination: any) => destination.status, mapFrom((source: QuestModel) => source.status >>> 0)),
    forMember((destination: any) => destination.mobKills, mapFrom((source: QuestModel) => normalizeNumberMap(source.mobKills))),
    forMember((destination: any) => destination.statusRecord, mapFrom((source: QuestModel) => source.statusRecord ?? "")),
    forMember((destination: any) => destination.recordEx, mapFrom((source: QuestModel) => normalizeStringMap(source.recordEx))),
    forMember((destination: any) => destination.completionTimeUnixMs, mapFrom((source: QuestModel) => source.completionTimeUnixMs ?? 0)),
    forMember((destination: any) => destination.forfeited, mapFrom((source: QuestModel) => source.forfeited >>> 0)),
    forMember((destination: any) => destination.deadlineUnixMs, mapFrom((source: QuestModel) => source.deadlineUnixMs ?? 0)),
    forMember((destination: any) => destination.startTimeUnixMs, mapFrom((source: QuestModel) => source.startTimeUnixMs ?? 0))
);

createMap(
    grpcMapper,
    QUEST_PROTO,
    QUEST_MODEL,
    forMember((destination: any) => destination.characterId, mapFrom((source: Quest) => source.characterId >>> 0)),
    forMember((destination: any) => destination.questId, mapFrom((source: Quest) => source.questId >>> 0)),
    forMember((destination: any) => destination.status, mapFrom((source: Quest) => source.status >>> 0)),
    forMember((destination: any) => destination.mobKills, mapFrom((source: Quest) => normalizeNumberMap(source.mobKills))),
    forMember((destination: any) => destination.statusRecord, mapFrom((source: Quest) => source.statusRecord ?? "")),
    forMember((destination: any) => destination.recordEx, mapFrom((source: Quest) => normalizeStringMap(source.recordEx))),
    forMember((destination: any) => destination.completionTimeUnixMs, mapFrom((source: Quest) => source.completionTimeUnixMs || 0)),
    forMember((destination: any) => destination.forfeited, mapFrom((source: Quest) => source.forfeited >>> 0)),
    forMember((destination: any) => destination.deadlineUnixMs, mapFrom((source: Quest) => source.deadlineUnixMs || 0)),
    forMember((destination: any) => destination.startTimeUnixMs, mapFrom((source: Quest) => source.startTimeUnixMs || 0))
);
