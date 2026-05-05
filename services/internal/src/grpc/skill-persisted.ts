import type { SkillPersisted } from "../protobuf/generated/fminternal/internal_service";

export function persistedFromMessage(msg: SkillPersisted) {
    const cooldownMs = msg.cooldownEndUnixMs;
    return {
        characterId: msg.characterId >>> 0,
        skillId: msg.skillId >>> 0,
        level: msg.level | 0,
        masterLevel: msg.masterLevel | 0,
        cooldownEndUnixMs: cooldownMs > 0 ? cooldownMs : null,
    };
}

export function makeSkillMessage(model: any) {
    return {
        characterId: model.characterId >>> 0,
        skillId: model.skillId >>> 0,
        level: model.level | 0,
        masterLevel: model.masterLevel | 0,
        cooldownEndUnixMs: model.cooldownEndUnixMs ?? 0,
    } as SkillPersisted;
}
