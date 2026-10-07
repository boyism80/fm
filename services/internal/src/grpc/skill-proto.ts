import type { Skill } from "../protobuf/generated/fminternal/internal_service";
import type { SkillModel } from "../repos/skill-repository";
import { createMap, forMember, mapFrom } from "@automapper/core";
import { grpcMapper } from "./mappers";

export const SKILL_MODEL = "SkillModel";
export const SKILL_PROTO = "SkillProto";

createMap(
    grpcMapper,
    SKILL_PROTO,
    SKILL_MODEL,
    forMember((destination: any) => destination.characterId, mapFrom((source: Skill) => source.characterId >>> 0)),
    forMember((destination: any) => destination.skillId, mapFrom((source: Skill) => source.skillId >>> 0)),
    forMember((destination: any) => destination.level, mapFrom((source: Skill) => source.level | 0)),
    forMember((destination: any) => destination.masterLevel, mapFrom((source: Skill) => source.masterLevel | 0)),
    forMember((destination: any) => destination.cooldownEndUnixMs, mapFrom((source: Skill) => source.cooldownEndUnixMs > 0 ? source.cooldownEndUnixMs : null))
);

createMap(
    grpcMapper,
    SKILL_MODEL,
    SKILL_PROTO,
    forMember((destination: any) => destination.characterId, mapFrom((source: SkillModel) => source.characterId >>> 0)),
    forMember((destination: any) => destination.skillId, mapFrom((source: SkillModel) => source.skillId >>> 0)),
    forMember((destination: any) => destination.level, mapFrom((source: SkillModel) => source.level | 0)),
    forMember((destination: any) => destination.masterLevel, mapFrom((source: SkillModel) => source.masterLevel | 0)),
    forMember((destination: any) => destination.cooldownEndUnixMs, mapFrom((source: SkillModel) => source.cooldownEndUnixMs ?? 0))
);
