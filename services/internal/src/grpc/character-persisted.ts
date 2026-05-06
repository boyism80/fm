import type { CharacterPersisted } from "../protobuf/generated/fminternal/internal_service";
import type { CharacterModel } from "../repos/character-repository";
import { createMap, forMember, mapFrom } from "@automapper/core";
import { grpcMapper } from "./mappers";

export const CHARACTER_MODEL = "CharacterModel";
export const CHARACTER_PERSISTED = "CharacterPersisted";

createMap(
    grpcMapper,
    CHARACTER_MODEL,
    CHARACTER_PERSISTED,
    forMember((destination: any) => destination.characterId, mapFrom((source: CharacterModel) => source.characterId >>> 0)),
    forMember((destination: any) => destination.worldId, mapFrom((source: CharacterModel) => source.worldId >>> 0)),
    forMember((destination: any) => destination.name, mapFrom((source: CharacterModel) => source.name ?? "")),
    forMember((destination: any) => destination.gender, mapFrom((source: CharacterModel) => source.gender >>> 0)),
    forMember((destination: any) => destination.skinColor, mapFrom((source: CharacterModel) => source.skinColor >>> 0)),
    forMember((destination: any) => destination.face, mapFrom((source: CharacterModel) => source.face >>> 0)),
    forMember((destination: any) => destination.hair, mapFrom((source: CharacterModel) => source.hair >>> 0)),
    forMember((destination: any) => destination.level, mapFrom((source: CharacterModel) => source.level >>> 0)),
    forMember((destination: any) => destination.classId, mapFrom((source: CharacterModel) => source.classId >>> 0)),
    forMember((destination: any) => destination.role, mapFrom((source: CharacterModel) => (source.role ?? 0) >>> 0)),
    forMember((destination: any) => destination.str, mapFrom((source: CharacterModel) => (source.str ?? 0) >>> 0)),
    forMember((destination: any) => destination.dex, mapFrom((source: CharacterModel) => (source.dex ?? 0) >>> 0)),
    forMember((destination: any) => destination.intStat, mapFrom((source: CharacterModel) => (source.intStat ?? 0) >>> 0)),
    forMember((destination: any) => destination.luk, mapFrom((source: CharacterModel) => (source.luk ?? 0) >>> 0)),
    forMember((destination: any) => destination.hp, mapFrom((source: CharacterModel) => (source.hp ?? 0) >>> 0)),
    forMember((destination: any) => destination.maxHp, mapFrom((source: CharacterModel) => (source.maxHp ?? 0) >>> 0)),
    forMember((destination: any) => destination.mp, mapFrom((source: CharacterModel) => (source.mp ?? 0) >>> 0)),
    forMember((destination: any) => destination.maxMp, mapFrom((source: CharacterModel) => (source.maxMp ?? 0) >>> 0)),
    forMember((destination: any) => destination.abilityPoint, mapFrom((source: CharacterModel) => (source.abilityPoint ?? 0) >>> 0)),
    forMember((destination: any) => destination.exp, mapFrom((source: CharacterModel) => (source.exp ?? 0) >>> 0)),
    forMember((destination: any) => destination.mapId, mapFrom((source: CharacterModel) => source.mapId >>> 0)),
    forMember((destination: any) => destination.spawnPoint, mapFrom((source: CharacterModel) => source.spawnPoint >>> 0)),
    forMember((destination: any) => destination.positionX, mapFrom((source: CharacterModel) => (source.positionX ?? 0) | 0)),
    forMember((destination: any) => destination.positionY, mapFrom((source: CharacterModel) => (source.positionY ?? 0) | 0)),
    forMember((destination: any) => destination.stance, mapFrom((source: CharacterModel) => (source.stance ?? 0) >>> 0)),
    forMember((destination: any) => destination.meso, mapFrom((source: CharacterModel) => (source.meso ?? 0) | 0)),
    forMember((destination: any) => destination.skillPoint, mapFrom((source: CharacterModel) => (source.skillPoint ?? 0) >>> 0)),
    forMember((destination: any) => destination.accountId, mapFrom((source: CharacterModel) => source.accountId >>> 0)),
    forMember((destination: any) => destination.hidden, mapFrom((source: CharacterModel) => source.hidden ?? false))
);

createMap(
    grpcMapper,
    CHARACTER_PERSISTED,
    CHARACTER_MODEL,
    forMember((destination: any) => destination.characterId, mapFrom((source: CharacterPersisted) => source.characterId >>> 0)),
    forMember((destination: any) => destination.worldId, mapFrom((source: CharacterPersisted) => source.worldId >>> 0)),
    forMember((destination: any) => destination.name, mapFrom((source: CharacterPersisted) => source.name)),
    forMember((destination: any) => destination.gender, mapFrom((source: CharacterPersisted) => source.gender >>> 0)),
    forMember((destination: any) => destination.skinColor, mapFrom((source: CharacterPersisted) => source.skinColor >>> 0)),
    forMember((destination: any) => destination.face, mapFrom((source: CharacterPersisted) => source.face >>> 0)),
    forMember((destination: any) => destination.hair, mapFrom((source: CharacterPersisted) => source.hair >>> 0)),
    forMember((destination: any) => destination.level, mapFrom((source: CharacterPersisted) => source.level >>> 0)),
    forMember((destination: any) => destination.classId, mapFrom((source: CharacterPersisted) => source.classId >>> 0)),
    forMember((destination: any) => destination.role, mapFrom((source: CharacterPersisted) => source.role >>> 0)),
    forMember((destination: any) => destination.str, mapFrom((source: CharacterPersisted) => source.str >>> 0)),
    forMember((destination: any) => destination.dex, mapFrom((source: CharacterPersisted) => source.dex >>> 0)),
    forMember((destination: any) => destination.intStat, mapFrom((source: CharacterPersisted) => source.intStat >>> 0)),
    forMember((destination: any) => destination.luk, mapFrom((source: CharacterPersisted) => source.luk >>> 0)),
    forMember((destination: any) => destination.hp, mapFrom((source: CharacterPersisted) => source.hp >>> 0)),
    forMember((destination: any) => destination.maxHp, mapFrom((source: CharacterPersisted) => source.maxHp >>> 0)),
    forMember((destination: any) => destination.mp, mapFrom((source: CharacterPersisted) => source.mp >>> 0)),
    forMember((destination: any) => destination.maxMp, mapFrom((source: CharacterPersisted) => source.maxMp >>> 0)),
    forMember((destination: any) => destination.abilityPoint, mapFrom((source: CharacterPersisted) => source.abilityPoint >>> 0)),
    forMember((destination: any) => destination.exp, mapFrom((source: CharacterPersisted) => source.exp >>> 0)),
    forMember((destination: any) => destination.mapId, mapFrom((source: CharacterPersisted) => source.mapId >>> 0)),
    forMember((destination: any) => destination.spawnPoint, mapFrom((source: CharacterPersisted) => source.spawnPoint >>> 0)),
    forMember((destination: any) => destination.positionX, mapFrom((source: CharacterPersisted) => source.positionX | 0)),
    forMember((destination: any) => destination.positionY, mapFrom((source: CharacterPersisted) => source.positionY | 0)),
    forMember((destination: any) => destination.stance, mapFrom((source: CharacterPersisted) => source.stance >>> 0)),
    forMember((destination: any) => destination.meso, mapFrom((source: CharacterPersisted) => source.meso | 0)),
    forMember((destination: any) => destination.skillPoint, mapFrom((source: CharacterPersisted) => source.skillPoint >>> 0)),
    forMember((destination: any) => destination.accountId, mapFrom((source: CharacterPersisted) => source.accountId >>> 0)),
    forMember((destination: any) => destination.hidden, mapFrom((source: CharacterPersisted) => source.hidden))
);
