import type { SkillPersisted } from "../protobuf/generated/fminternal/internal_service";
import type { SkillModel } from "../repos/skill-repository";
import { createMap, forMember, mapFrom } from "@automapper/core";
import { grpcMapper } from "./mappers";

export const SKILL_MODEL = "SkillModel";
export const SKILL_PERSISTED = "SkillPersisted";

createMap(
    grpcMapper,
    SKILL_PERSISTED,
    SKILL_MODEL,
    forMember((destination: any) => destination.characterId, mapFrom((source: SkillPersisted) => source.characterId >>> 0)),
    forMember((destination: any) => destination.skillId, mapFrom((source: SkillPersisted) => source.skillId >>> 0)),
    forMember((destination: any) => destination.level, mapFrom((source: SkillPersisted) => source.level | 0)),
    forMember((destination: any) => destination.masterLevel, mapFrom((source: SkillPersisted) => source.masterLevel | 0)),
    forMember((destination: any) => destination.cooldownEndUnixMs, mapFrom((source: SkillPersisted) => source.cooldownEndUnixMs > 0 ? source.cooldownEndUnixMs : null))
);

createMap(
    grpcMapper,
    SKILL_MODEL,
    SKILL_PERSISTED,
    forMember((destination: any) => destination.characterId, mapFrom((source: SkillModel) => source.characterId >>> 0)),
    forMember((destination: any) => destination.skillId, mapFrom((source: SkillModel) => source.skillId >>> 0)),
    forMember((destination: any) => destination.level, mapFrom((source: SkillModel) => source.level | 0)),
    forMember((destination: any) => destination.masterLevel, mapFrom((source: SkillModel) => source.masterLevel | 0)),
    forMember((destination: any) => destination.cooldownEndUnixMs, mapFrom((source: SkillModel) => source.cooldownEndUnixMs ?? 0))
);
