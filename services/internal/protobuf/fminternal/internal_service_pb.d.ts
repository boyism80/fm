// package: fm.internal
// file: fminternal/internal_service.proto

/* tslint:disable */
/* eslint-disable */

import * as jspb from "google-protobuf";

export class PingRequest extends jspb.Message { 

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): PingRequest.AsObject;
    static toObject(includeInstance: boolean, msg: PingRequest): PingRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: PingRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): PingRequest;
    static deserializeBinaryFromReader(message: PingRequest, reader: jspb.BinaryReader): PingRequest;
}

export namespace PingRequest {
    export type AsObject = {
    }
}

export class PingReply extends jspb.Message { 
    getMessage(): string;
    setMessage(value: string): PingReply;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): PingReply.AsObject;
    static toObject(includeInstance: boolean, msg: PingReply): PingReply.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: PingReply, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): PingReply;
    static deserializeBinaryFromReader(message: PingReply, reader: jspb.BinaryReader): PingReply;
}

export namespace PingReply {
    export type AsObject = {
        message: string,
    }
}

export class GetServerCatalogRequest extends jspb.Message { 

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): GetServerCatalogRequest.AsObject;
    static toObject(includeInstance: boolean, msg: GetServerCatalogRequest): GetServerCatalogRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: GetServerCatalogRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): GetServerCatalogRequest;
    static deserializeBinaryFromReader(message: GetServerCatalogRequest, reader: jspb.BinaryReader): GetServerCatalogRequest;
}

export namespace GetServerCatalogRequest {
    export type AsObject = {
    }
}

export class ChannelCatalog extends jspb.Message { 
    getChannelId(): number;
    setChannelId(value: number): ChannelCatalog;
    getHost(): string;
    setHost(value: string): ChannelCatalog;
    getPort(): number;
    setPort(value: number): ChannelCatalog;
    getName(): string;
    setName(value: string): ChannelCatalog;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): ChannelCatalog.AsObject;
    static toObject(includeInstance: boolean, msg: ChannelCatalog): ChannelCatalog.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: ChannelCatalog, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): ChannelCatalog;
    static deserializeBinaryFromReader(message: ChannelCatalog, reader: jspb.BinaryReader): ChannelCatalog;
}

export namespace ChannelCatalog {
    export type AsObject = {
        channelId: number,
        host: string,
        port: number,
        name: string,
    }
}

export class WorldCatalog extends jspb.Message { 
    getWorldId(): number;
    setWorldId(value: number): WorldCatalog;
    getWorldName(): string;
    setWorldName(value: string): WorldCatalog;
    getFlag(): number;
    setFlag(value: number): WorldCatalog;
    getEventMessage(): string;
    setEventMessage(value: string): WorldCatalog;
    clearChannelsList(): void;
    getChannelsList(): Array<ChannelCatalog>;
    setChannelsList(value: Array<ChannelCatalog>): WorldCatalog;
    addChannels(value?: ChannelCatalog, index?: number): ChannelCatalog;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): WorldCatalog.AsObject;
    static toObject(includeInstance: boolean, msg: WorldCatalog): WorldCatalog.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: WorldCatalog, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): WorldCatalog;
    static deserializeBinaryFromReader(message: WorldCatalog, reader: jspb.BinaryReader): WorldCatalog;
}

export namespace WorldCatalog {
    export type AsObject = {
        worldId: number,
        worldName: string,
        flag: number,
        eventMessage: string,
        channelsList: Array<ChannelCatalog.AsObject>,
    }
}

export class GetServerCatalogReply extends jspb.Message { 
    clearWorldsList(): void;
    getWorldsList(): Array<WorldCatalog>;
    setWorldsList(value: Array<WorldCatalog>): GetServerCatalogReply;
    addWorlds(value?: WorldCatalog, index?: number): WorldCatalog;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): GetServerCatalogReply.AsObject;
    static toObject(includeInstance: boolean, msg: GetServerCatalogReply): GetServerCatalogReply.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: GetServerCatalogReply, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): GetServerCatalogReply;
    static deserializeBinaryFromReader(message: GetServerCatalogReply, reader: jspb.BinaryReader): GetServerCatalogReply;
}

export namespace GetServerCatalogReply {
    export type AsObject = {
        worldsList: Array<WorldCatalog.AsObject>,
    }
}

export class EnterGameRequest extends jspb.Message { 
    getWorldId(): number;
    setWorldId(value: number): EnterGameRequest;
    getCharacterId(): number;
    setCharacterId(value: number): EnterGameRequest;
    getChannelId(): number;
    setChannelId(value: number): EnterGameRequest;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): EnterGameRequest.AsObject;
    static toObject(includeInstance: boolean, msg: EnterGameRequest): EnterGameRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: EnterGameRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): EnterGameRequest;
    static deserializeBinaryFromReader(message: EnterGameRequest, reader: jspb.BinaryReader): EnterGameRequest;
}

export namespace EnterGameRequest {
    export type AsObject = {
        worldId: number,
        characterId: number,
        channelId: number,
    }
}

export class CharacterPersisted extends jspb.Message { 
    getCharacterId(): number;
    setCharacterId(value: number): CharacterPersisted;
    getWorldId(): number;
    setWorldId(value: number): CharacterPersisted;
    getName(): string;
    setName(value: string): CharacterPersisted;
    getGender(): number;
    setGender(value: number): CharacterPersisted;
    getSkinColor(): number;
    setSkinColor(value: number): CharacterPersisted;
    getFace(): number;
    setFace(value: number): CharacterPersisted;
    getHair(): number;
    setHair(value: number): CharacterPersisted;
    getLevel(): number;
    setLevel(value: number): CharacterPersisted;
    getClassId(): number;
    setClassId(value: number): CharacterPersisted;
    getStr(): number;
    setStr(value: number): CharacterPersisted;
    getDex(): number;
    setDex(value: number): CharacterPersisted;
    getIntStat(): number;
    setIntStat(value: number): CharacterPersisted;
    getLuk(): number;
    setLuk(value: number): CharacterPersisted;
    getHp(): number;
    setHp(value: number): CharacterPersisted;
    getMaxHp(): number;
    setMaxHp(value: number): CharacterPersisted;
    getMp(): number;
    setMp(value: number): CharacterPersisted;
    getMaxMp(): number;
    setMaxMp(value: number): CharacterPersisted;
    getAbilityPoint(): number;
    setAbilityPoint(value: number): CharacterPersisted;
    getExp(): number;
    setExp(value: number): CharacterPersisted;
    getMapId(): number;
    setMapId(value: number): CharacterPersisted;
    getSpawnPoint(): number;
    setSpawnPoint(value: number): CharacterPersisted;
    getPositionX(): number;
    setPositionX(value: number): CharacterPersisted;
    getPositionY(): number;
    setPositionY(value: number): CharacterPersisted;
    getStance(): number;
    setStance(value: number): CharacterPersisted;
    getMeso(): number;
    setMeso(value: number): CharacterPersisted;
    getSkillPoint(): number;
    setSkillPoint(value: number): CharacterPersisted;
    getAccountId(): number;
    setAccountId(value: number): CharacterPersisted;
    getRole(): number;
    setRole(value: number): CharacterPersisted;
    getHidden(): boolean;
    setHidden(value: boolean): CharacterPersisted;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): CharacterPersisted.AsObject;
    static toObject(includeInstance: boolean, msg: CharacterPersisted): CharacterPersisted.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: CharacterPersisted, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): CharacterPersisted;
    static deserializeBinaryFromReader(message: CharacterPersisted, reader: jspb.BinaryReader): CharacterPersisted;
}

export namespace CharacterPersisted {
    export type AsObject = {
        characterId: number,
        worldId: number,
        name: string,
        gender: number,
        skinColor: number,
        face: number,
        hair: number,
        level: number,
        classId: number,
        str: number,
        dex: number,
        intStat: number,
        luk: number,
        hp: number,
        maxHp: number,
        mp: number,
        maxMp: number,
        abilityPoint: number,
        exp: number,
        mapId: number,
        spawnPoint: number,
        positionX: number,
        positionY: number,
        stance: number,
        meso: number,
        skillPoint: number,
        accountId: number,
        role: number,
        hidden: boolean,
    }
}

export class KeyLayoutBinding extends jspb.Message { 
    getSlot(): number;
    setSlot(value: number): KeyLayoutBinding;
    getType(): number;
    setType(value: number): KeyLayoutBinding;
    getAction(): number;
    setAction(value: number): KeyLayoutBinding;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): KeyLayoutBinding.AsObject;
    static toObject(includeInstance: boolean, msg: KeyLayoutBinding): KeyLayoutBinding.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: KeyLayoutBinding, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): KeyLayoutBinding;
    static deserializeBinaryFromReader(message: KeyLayoutBinding, reader: jspb.BinaryReader): KeyLayoutBinding;
}

export namespace KeyLayoutBinding {
    export type AsObject = {
        slot: number,
        type: number,
        action: number,
    }
}

export class EnterGameReply extends jspb.Message { 
    getFound(): boolean;
    setFound(value: boolean): EnterGameReply;

    hasCharacter(): boolean;
    clearCharacter(): void;
    getCharacter(): CharacterPersisted | undefined;
    setCharacter(value?: CharacterPersisted): EnterGameReply;
    clearInventoryList(): void;
    getInventoryList(): Array<InventoryPersisted>;
    setInventoryList(value: Array<InventoryPersisted>): EnterGameReply;
    addInventory(value?: InventoryPersisted, index?: number): InventoryPersisted;
    clearSkillsList(): void;
    getSkillsList(): Array<SkillPersisted>;
    setSkillsList(value: Array<SkillPersisted>): EnterGameReply;
    addSkills(value?: SkillPersisted, index?: number): SkillPersisted;
    clearKeyLayoutList(): void;
    getKeyLayoutList(): Array<KeyLayoutBinding>;
    setKeyLayoutList(value: Array<KeyLayoutBinding>): EnterGameReply;
    addKeyLayout(value?: KeyLayoutBinding, index?: number): KeyLayoutBinding;

    hasPartyId(): boolean;
    clearPartyId(): void;
    getPartyId(): number | undefined;
    setPartyId(value: number): EnterGameReply;
    getGuildId(): number;
    setGuildId(value: number): EnterGameReply;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): EnterGameReply.AsObject;
    static toObject(includeInstance: boolean, msg: EnterGameReply): EnterGameReply.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: EnterGameReply, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): EnterGameReply;
    static deserializeBinaryFromReader(message: EnterGameReply, reader: jspb.BinaryReader): EnterGameReply;
}

export namespace EnterGameReply {
    export type AsObject = {
        found: boolean,
        character?: CharacterPersisted.AsObject,
        inventoryList: Array<InventoryPersisted.AsObject>,
        skillsList: Array<SkillPersisted.AsObject>,
        keyLayoutList: Array<KeyLayoutBinding.AsObject>,
        partyId?: number,
        guildId: number,
    }
}

export class BeginGameTransitionRequest extends jspb.Message { 
    getWorldId(): number;
    setWorldId(value: number): BeginGameTransitionRequest;
    getAccountId(): number;
    setAccountId(value: number): BeginGameTransitionRequest;
    getCharacterId(): number;
    setCharacterId(value: number): BeginGameTransitionRequest;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): BeginGameTransitionRequest.AsObject;
    static toObject(includeInstance: boolean, msg: BeginGameTransitionRequest): BeginGameTransitionRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: BeginGameTransitionRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): BeginGameTransitionRequest;
    static deserializeBinaryFromReader(message: BeginGameTransitionRequest, reader: jspb.BinaryReader): BeginGameTransitionRequest;
}

export namespace BeginGameTransitionRequest {
    export type AsObject = {
        worldId: number,
        accountId: number,
        characterId: number,
    }
}

export class BeginGameTransitionReply extends jspb.Message { 
    getOk(): boolean;
    setOk(value: boolean): BeginGameTransitionReply;
    getErrorCode(): SessionErrorCode;
    setErrorCode(value: SessionErrorCode): BeginGameTransitionReply;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): BeginGameTransitionReply.AsObject;
    static toObject(includeInstance: boolean, msg: BeginGameTransitionReply): BeginGameTransitionReply.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: BeginGameTransitionReply, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): BeginGameTransitionReply;
    static deserializeBinaryFromReader(message: BeginGameTransitionReply, reader: jspb.BinaryReader): BeginGameTransitionReply;
}

export namespace BeginGameTransitionReply {
    export type AsObject = {
        ok: boolean,
        errorCode: SessionErrorCode,
    }
}

export class SaveCharacterRequest extends jspb.Message { 

    hasCharacter(): boolean;
    clearCharacter(): void;
    getCharacter(): CharacterPersisted | undefined;
    setCharacter(value?: CharacterPersisted): SaveCharacterRequest;

    getBaseLooksMap(): jspb.Map<number, number>;
    clearBaseLooksMap(): void;

    getOverlaysMap(): jspb.Map<number, number>;
    clearOverlaysMap(): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): SaveCharacterRequest.AsObject;
    static toObject(includeInstance: boolean, msg: SaveCharacterRequest): SaveCharacterRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: SaveCharacterRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): SaveCharacterRequest;
    static deserializeBinaryFromReader(message: SaveCharacterRequest, reader: jspb.BinaryReader): SaveCharacterRequest;
}

export namespace SaveCharacterRequest {
    export type AsObject = {
        character?: CharacterPersisted.AsObject,

        baseLooksMap: Array<[number, number]>,

        overlaysMap: Array<[number, number]>,
    }
}

export class SaveCharacterReply extends jspb.Message { 
    getOk(): boolean;
    setOk(value: boolean): SaveCharacterReply;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): SaveCharacterReply.AsObject;
    static toObject(includeInstance: boolean, msg: SaveCharacterReply): SaveCharacterReply.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: SaveCharacterReply, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): SaveCharacterReply;
    static deserializeBinaryFromReader(message: SaveCharacterReply, reader: jspb.BinaryReader): SaveCharacterReply;
}

export namespace SaveCharacterReply {
    export type AsObject = {
        ok: boolean,
    }
}

export class CharacterSaveEntry extends jspb.Message { 

    hasCharacter(): boolean;
    clearCharacter(): void;
    getCharacter(): CharacterPersisted | undefined;
    setCharacter(value?: CharacterPersisted): CharacterSaveEntry;

    getBaseLooksMap(): jspb.Map<number, number>;
    clearBaseLooksMap(): void;

    getOverlaysMap(): jspb.Map<number, number>;
    clearOverlaysMap(): void;
    clearInventoryList(): void;
    getInventoryList(): Array<InventoryPersisted>;
    setInventoryList(value: Array<InventoryPersisted>): CharacterSaveEntry;
    addInventory(value?: InventoryPersisted, index?: number): InventoryPersisted;
    clearSkillsList(): void;
    getSkillsList(): Array<SkillPersisted>;
    setSkillsList(value: Array<SkillPersisted>): CharacterSaveEntry;
    addSkills(value?: SkillPersisted, index?: number): SkillPersisted;
    clearKeyLayoutList(): void;
    getKeyLayoutList(): Array<KeyLayoutBinding>;
    setKeyLayoutList(value: Array<KeyLayoutBinding>): CharacterSaveEntry;
    addKeyLayout(value?: KeyLayoutBinding, index?: number): KeyLayoutBinding;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): CharacterSaveEntry.AsObject;
    static toObject(includeInstance: boolean, msg: CharacterSaveEntry): CharacterSaveEntry.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: CharacterSaveEntry, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): CharacterSaveEntry;
    static deserializeBinaryFromReader(message: CharacterSaveEntry, reader: jspb.BinaryReader): CharacterSaveEntry;
}

export namespace CharacterSaveEntry {
    export type AsObject = {
        character?: CharacterPersisted.AsObject,

        baseLooksMap: Array<[number, number]>,

        overlaysMap: Array<[number, number]>,
        inventoryList: Array<InventoryPersisted.AsObject>,
        skillsList: Array<SkillPersisted.AsObject>,
        keyLayoutList: Array<KeyLayoutBinding.AsObject>,
    }
}

export class SaveCharactersRequest extends jspb.Message { 
    clearEntriesList(): void;
    getEntriesList(): Array<CharacterSaveEntry>;
    setEntriesList(value: Array<CharacterSaveEntry>): SaveCharactersRequest;
    addEntries(value?: CharacterSaveEntry, index?: number): CharacterSaveEntry;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): SaveCharactersRequest.AsObject;
    static toObject(includeInstance: boolean, msg: SaveCharactersRequest): SaveCharactersRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: SaveCharactersRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): SaveCharactersRequest;
    static deserializeBinaryFromReader(message: SaveCharactersRequest, reader: jspb.BinaryReader): SaveCharactersRequest;
}

export namespace SaveCharactersRequest {
    export type AsObject = {
        entriesList: Array<CharacterSaveEntry.AsObject>,
    }
}

export class SaveCharactersReply extends jspb.Message { 
    getOk(): boolean;
    setOk(value: boolean): SaveCharactersReply;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): SaveCharactersReply.AsObject;
    static toObject(includeInstance: boolean, msg: SaveCharactersReply): SaveCharactersReply.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: SaveCharactersReply, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): SaveCharactersReply;
    static deserializeBinaryFromReader(message: SaveCharactersReply, reader: jspb.BinaryReader): SaveCharactersReply;
}

export namespace SaveCharactersReply {
    export type AsObject = {
        ok: boolean,
    }
}

export class EquipmentBonusStatsPersisted extends jspb.Message { 
    getStr(): number;
    setStr(value: number): EquipmentBonusStatsPersisted;
    getDex(): number;
    setDex(value: number): EquipmentBonusStatsPersisted;
    getIntStat(): number;
    setIntStat(value: number): EquipmentBonusStatsPersisted;
    getLuk(): number;
    setLuk(value: number): EquipmentBonusStatsPersisted;
    getMaxHp(): number;
    setMaxHp(value: number): EquipmentBonusStatsPersisted;
    getMaxMp(): number;
    setMaxMp(value: number): EquipmentBonusStatsPersisted;
    getPad(): number;
    setPad(value: number): EquipmentBonusStatsPersisted;
    getMad(): number;
    setMad(value: number): EquipmentBonusStatsPersisted;
    getPdd(): number;
    setPdd(value: number): EquipmentBonusStatsPersisted;
    getMdd(): number;
    setMdd(value: number): EquipmentBonusStatsPersisted;
    getAcc(): number;
    setAcc(value: number): EquipmentBonusStatsPersisted;
    getAvoid(): number;
    setAvoid(value: number): EquipmentBonusStatsPersisted;
    getHands(): number;
    setHands(value: number): EquipmentBonusStatsPersisted;
    getSpeed(): number;
    setSpeed(value: number): EquipmentBonusStatsPersisted;
    getJump(): number;
    setJump(value: number): EquipmentBonusStatsPersisted;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): EquipmentBonusStatsPersisted.AsObject;
    static toObject(includeInstance: boolean, msg: EquipmentBonusStatsPersisted): EquipmentBonusStatsPersisted.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: EquipmentBonusStatsPersisted, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): EquipmentBonusStatsPersisted;
    static deserializeBinaryFromReader(message: EquipmentBonusStatsPersisted, reader: jspb.BinaryReader): EquipmentBonusStatsPersisted;
}

export namespace EquipmentBonusStatsPersisted {
    export type AsObject = {
        str: number,
        dex: number,
        intStat: number,
        luk: number,
        maxHp: number,
        maxMp: number,
        pad: number,
        mad: number,
        pdd: number,
        mdd: number,
        acc: number,
        avoid: number,
        hands: number,
        speed: number,
        jump: number,
    }
}

export class InventoryPersisted extends jspb.Message { 

    hasUniqueId(): boolean;
    clearUniqueId(): void;
    getUniqueId(): number | undefined;
    setUniqueId(value: number): InventoryPersisted;
    getOwnerId(): number;
    setOwnerId(value: number): InventoryPersisted;
    getItemId(): number;
    setItemId(value: number): InventoryPersisted;
    getSlot(): number;
    setSlot(value: number): InventoryPersisted;
    getCount(): number;
    setCount(value: number): InventoryPersisted;
    getExpirationUnixMs(): number;
    setExpirationUnixMs(value: number): InventoryPersisted;
    getEnhanceChance(): number;
    setEnhanceChance(value: number): InventoryPersisted;
    getFlag(): number;
    setFlag(value: number): InventoryPersisted;
    getSkillBonus(): number;
    setSkillBonus(value: number): InventoryPersisted;
    getOwnerName(): string;
    setOwnerName(value: string): InventoryPersisted;
    getInventoryType(): number;
    setInventoryType(value: number): InventoryPersisted;

    hasEquipBonusStats(): boolean;
    clearEquipBonusStats(): void;
    getEquipBonusStats(): EquipmentBonusStatsPersisted | undefined;
    setEquipBonusStats(value?: EquipmentBonusStatsPersisted): InventoryPersisted;
    getEnhanceCount(): number;
    setEnhanceCount(value: number): InventoryPersisted;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): InventoryPersisted.AsObject;
    static toObject(includeInstance: boolean, msg: InventoryPersisted): InventoryPersisted.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: InventoryPersisted, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): InventoryPersisted;
    static deserializeBinaryFromReader(message: InventoryPersisted, reader: jspb.BinaryReader): InventoryPersisted;
}

export namespace InventoryPersisted {
    export type AsObject = {
        uniqueId?: number,
        ownerId: number,
        itemId: number,
        slot: number,
        count: number,
        expirationUnixMs: number,
        enhanceChance: number,
        flag: number,
        skillBonus: number,
        ownerName: string,
        inventoryType: number,
        equipBonusStats?: EquipmentBonusStatsPersisted.AsObject,
        enhanceCount: number,
    }
}

export class SkillPersisted extends jspb.Message { 
    getCharacterId(): number;
    setCharacterId(value: number): SkillPersisted;
    getSkillId(): number;
    setSkillId(value: number): SkillPersisted;
    getLevel(): number;
    setLevel(value: number): SkillPersisted;
    getMasterLevel(): number;
    setMasterLevel(value: number): SkillPersisted;
    getCooldownEndUnixMs(): number;
    setCooldownEndUnixMs(value: number): SkillPersisted;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): SkillPersisted.AsObject;
    static toObject(includeInstance: boolean, msg: SkillPersisted): SkillPersisted.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: SkillPersisted, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): SkillPersisted;
    static deserializeBinaryFromReader(message: SkillPersisted, reader: jspb.BinaryReader): SkillPersisted;
}

export namespace SkillPersisted {
    export type AsObject = {
        characterId: number,
        skillId: number,
        level: number,
        masterLevel: number,
        cooldownEndUnixMs: number,
    }
}

export class LoginAccountRequest extends jspb.Message { 
    getLoginId(): string;
    setLoginId(value: string): LoginAccountRequest;
    getPassword(): string;
    setPassword(value: string): LoginAccountRequest;
    getMacAddress(): string;
    setMacAddress(value: string): LoginAccountRequest;
    getIpAddress(): string;
    setIpAddress(value: string): LoginAccountRequest;
    getInitialRole(): number;
    setInitialRole(value: number): LoginAccountRequest;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): LoginAccountRequest.AsObject;
    static toObject(includeInstance: boolean, msg: LoginAccountRequest): LoginAccountRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: LoginAccountRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): LoginAccountRequest;
    static deserializeBinaryFromReader(message: LoginAccountRequest, reader: jspb.BinaryReader): LoginAccountRequest;
}

export namespace LoginAccountRequest {
    export type AsObject = {
        loginId: string,
        password: string,
        macAddress: string,
        ipAddress: string,
        initialRole: number,
    }
}

export class LoginAccountReply extends jspb.Message { 
    getStatus(): LoginAccountReply.Status;
    setStatus(value: LoginAccountReply.Status): LoginAccountReply;
    getAccountId(): number;
    setAccountId(value: number): LoginAccountReply;
    getGender(): number;
    setGender(value: number): LoginAccountReply;
    getIsChatBlocked(): boolean;
    setIsChatBlocked(value: boolean): LoginAccountReply;
    getChatBlockedUntilUnixMs(): number;
    setChatBlockedUntilUnixMs(value: number): LoginAccountReply;
    getCharacterSlotCount(): number;
    setCharacterSlotCount(value: number): LoginAccountReply;
    getRole(): number;
    setRole(value: number): LoginAccountReply;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): LoginAccountReply.AsObject;
    static toObject(includeInstance: boolean, msg: LoginAccountReply): LoginAccountReply.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: LoginAccountReply, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): LoginAccountReply;
    static deserializeBinaryFromReader(message: LoginAccountReply, reader: jspb.BinaryReader): LoginAccountReply;
}

export namespace LoginAccountReply {
    export type AsObject = {
        status: LoginAccountReply.Status,
        accountId: number,
        gender: number,
        isChatBlocked: boolean,
        chatBlockedUntilUnixMs: number,
        characterSlotCount: number,
        role: number,
    }

    export enum Status {
    SUCCESS = 0,
    WRONG_PASSWORD = 1,
    BANNED = 2,
    REGISTERED = 3,
    ALREADY_LOGGED_IN = 4,
    }

}

export class CharacterOverview extends jspb.Message { 
    getCharacterId(): number;
    setCharacterId(value: number): CharacterOverview;
    getName(): string;
    setName(value: string): CharacterOverview;
    getGender(): number;
    setGender(value: number): CharacterOverview;
    getSkinColor(): number;
    setSkinColor(value: number): CharacterOverview;
    getFace(): number;
    setFace(value: number): CharacterOverview;
    getHair(): number;
    setHair(value: number): CharacterOverview;
    getLevel(): number;
    setLevel(value: number): CharacterOverview;
    getClassId(): number;
    setClassId(value: number): CharacterOverview;
    getMapId(): number;
    setMapId(value: number): CharacterOverview;
    getSpawnPoint(): number;
    setSpawnPoint(value: number): CharacterOverview;
    getAccountId(): number;
    setAccountId(value: number): CharacterOverview;
    getWorldId(): number;
    setWorldId(value: number): CharacterOverview;
    getRank(): number;
    setRank(value: number): CharacterOverview;
    getRankDiff(): number;
    setRankDiff(value: number): CharacterOverview;
    getClassRank(): number;
    setClassRank(value: number): CharacterOverview;
    getClassRankDiff(): number;
    setClassRankDiff(value: number): CharacterOverview;

    getBaseLooksMap(): jspb.Map<number, number>;
    clearBaseLooksMap(): void;

    getOverlaysMap(): jspb.Map<number, number>;
    clearOverlaysMap(): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): CharacterOverview.AsObject;
    static toObject(includeInstance: boolean, msg: CharacterOverview): CharacterOverview.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: CharacterOverview, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): CharacterOverview;
    static deserializeBinaryFromReader(message: CharacterOverview, reader: jspb.BinaryReader): CharacterOverview;
}

export namespace CharacterOverview {
    export type AsObject = {
        characterId: number,
        name: string,
        gender: number,
        skinColor: number,
        face: number,
        hair: number,
        level: number,
        classId: number,
        mapId: number,
        spawnPoint: number,
        accountId: number,
        worldId: number,
        rank: number,
        rankDiff: number,
        classRank: number,
        classRankDiff: number,

        baseLooksMap: Array<[number, number]>,

        overlaysMap: Array<[number, number]>,
    }
}

export class GetCharacterListRequest extends jspb.Message { 
    getAccountId(): number;
    setAccountId(value: number): GetCharacterListRequest;
    getWorldId(): number;
    setWorldId(value: number): GetCharacterListRequest;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): GetCharacterListRequest.AsObject;
    static toObject(includeInstance: boolean, msg: GetCharacterListRequest): GetCharacterListRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: GetCharacterListRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): GetCharacterListRequest;
    static deserializeBinaryFromReader(message: GetCharacterListRequest, reader: jspb.BinaryReader): GetCharacterListRequest;
}

export namespace GetCharacterListRequest {
    export type AsObject = {
        accountId: number,
        worldId: number,
    }
}

export class GetCharacterListReply extends jspb.Message { 
    clearCharactersList(): void;
    getCharactersList(): Array<CharacterOverview>;
    setCharactersList(value: Array<CharacterOverview>): GetCharacterListReply;
    addCharacters(value?: CharacterOverview, index?: number): CharacterOverview;
    getSlotCount(): number;
    setSlotCount(value: number): GetCharacterListReply;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): GetCharacterListReply.AsObject;
    static toObject(includeInstance: boolean, msg: GetCharacterListReply): GetCharacterListReply.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: GetCharacterListReply, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): GetCharacterListReply;
    static deserializeBinaryFromReader(message: GetCharacterListReply, reader: jspb.BinaryReader): GetCharacterListReply;
}

export namespace GetCharacterListReply {
    export type AsObject = {
        charactersList: Array<CharacterOverview.AsObject>,
        slotCount: number,
    }
}

export class CheckCharacterNameRequest extends jspb.Message { 
    getName(): string;
    setName(value: string): CheckCharacterNameRequest;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): CheckCharacterNameRequest.AsObject;
    static toObject(includeInstance: boolean, msg: CheckCharacterNameRequest): CheckCharacterNameRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: CheckCharacterNameRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): CheckCharacterNameRequest;
    static deserializeBinaryFromReader(message: CheckCharacterNameRequest, reader: jspb.BinaryReader): CheckCharacterNameRequest;
}

export namespace CheckCharacterNameRequest {
    export type AsObject = {
        name: string,
    }
}

export class CheckCharacterNameReply extends jspb.Message { 
    getExists(): boolean;
    setExists(value: boolean): CheckCharacterNameReply;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): CheckCharacterNameReply.AsObject;
    static toObject(includeInstance: boolean, msg: CheckCharacterNameReply): CheckCharacterNameReply.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: CheckCharacterNameReply, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): CheckCharacterNameReply;
    static deserializeBinaryFromReader(message: CheckCharacterNameReply, reader: jspb.BinaryReader): CheckCharacterNameReply;
}

export namespace CheckCharacterNameReply {
    export type AsObject = {
        exists: boolean,
    }
}

export class CreateCharacterRequest extends jspb.Message { 
    getAccountId(): number;
    setAccountId(value: number): CreateCharacterRequest;
    getWorldId(): number;
    setWorldId(value: number): CreateCharacterRequest;
    getName(): string;
    setName(value: string): CreateCharacterRequest;
    getFace(): number;
    setFace(value: number): CreateCharacterRequest;
    getHair(): number;
    setHair(value: number): CreateCharacterRequest;
    getSkinColor(): number;
    setSkinColor(value: number): CreateCharacterRequest;
    getTopItemId(): number;
    setTopItemId(value: number): CreateCharacterRequest;
    getBottomItemId(): number;
    setBottomItemId(value: number): CreateCharacterRequest;
    getShoesItemId(): number;
    setShoesItemId(value: number): CreateCharacterRequest;
    getWeaponItemId(): number;
    setWeaponItemId(value: number): CreateCharacterRequest;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): CreateCharacterRequest.AsObject;
    static toObject(includeInstance: boolean, msg: CreateCharacterRequest): CreateCharacterRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: CreateCharacterRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): CreateCharacterRequest;
    static deserializeBinaryFromReader(message: CreateCharacterRequest, reader: jspb.BinaryReader): CreateCharacterRequest;
}

export namespace CreateCharacterRequest {
    export type AsObject = {
        accountId: number,
        worldId: number,
        name: string,
        face: number,
        hair: number,
        skinColor: number,
        topItemId: number,
        bottomItemId: number,
        shoesItemId: number,
        weaponItemId: number,
    }
}

export class CreateCharacterReply extends jspb.Message { 
    getSuccess(): boolean;
    setSuccess(value: boolean): CreateCharacterReply;
    getErrorMsg(): string;
    setErrorMsg(value: string): CreateCharacterReply;

    hasCharacter(): boolean;
    clearCharacter(): void;
    getCharacter(): CharacterOverview | undefined;
    setCharacter(value?: CharacterOverview): CreateCharacterReply;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): CreateCharacterReply.AsObject;
    static toObject(includeInstance: boolean, msg: CreateCharacterReply): CreateCharacterReply.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: CreateCharacterReply, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): CreateCharacterReply;
    static deserializeBinaryFromReader(message: CreateCharacterReply, reader: jspb.BinaryReader): CreateCharacterReply;
}

export namespace CreateCharacterReply {
    export type AsObject = {
        success: boolean,
        errorMsg: string,
        character?: CharacterOverview.AsObject,
    }
}

export class DeleteCharacterRequest extends jspb.Message { 
    getAccountId(): number;
    setAccountId(value: number): DeleteCharacterRequest;
    getCharacterId(): number;
    setCharacterId(value: number): DeleteCharacterRequest;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): DeleteCharacterRequest.AsObject;
    static toObject(includeInstance: boolean, msg: DeleteCharacterRequest): DeleteCharacterRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: DeleteCharacterRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): DeleteCharacterRequest;
    static deserializeBinaryFromReader(message: DeleteCharacterRequest, reader: jspb.BinaryReader): DeleteCharacterRequest;
}

export namespace DeleteCharacterRequest {
    export type AsObject = {
        accountId: number,
        characterId: number,
    }
}

export class DeleteCharacterReply extends jspb.Message { 
    getSuccess(): boolean;
    setSuccess(value: boolean): DeleteCharacterReply;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): DeleteCharacterReply.AsObject;
    static toObject(includeInstance: boolean, msg: DeleteCharacterReply): DeleteCharacterReply.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: DeleteCharacterReply, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): DeleteCharacterReply;
    static deserializeBinaryFromReader(message: DeleteCharacterReply, reader: jspb.BinaryReader): DeleteCharacterReply;
}

export namespace DeleteCharacterReply {
    export type AsObject = {
        success: boolean,
    }
}

export class RefreshSessionRequest extends jspb.Message { 
    getWorldId(): number;
    setWorldId(value: number): RefreshSessionRequest;
    getAccountId(): number;
    setAccountId(value: number): RefreshSessionRequest;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): RefreshSessionRequest.AsObject;
    static toObject(includeInstance: boolean, msg: RefreshSessionRequest): RefreshSessionRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: RefreshSessionRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): RefreshSessionRequest;
    static deserializeBinaryFromReader(message: RefreshSessionRequest, reader: jspb.BinaryReader): RefreshSessionRequest;
}

export namespace RefreshSessionRequest {
    export type AsObject = {
        worldId: number,
        accountId: number,
    }
}

export class RefreshSessionReply extends jspb.Message { 
    getOk(): boolean;
    setOk(value: boolean): RefreshSessionReply;
    getErrorCode(): SessionErrorCode;
    setErrorCode(value: SessionErrorCode): RefreshSessionReply;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): RefreshSessionReply.AsObject;
    static toObject(includeInstance: boolean, msg: RefreshSessionReply): RefreshSessionReply.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: RefreshSessionReply, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): RefreshSessionReply;
    static deserializeBinaryFromReader(message: RefreshSessionReply, reader: jspb.BinaryReader): RefreshSessionReply;
}

export namespace RefreshSessionReply {
    export type AsObject = {
        ok: boolean,
        errorCode: SessionErrorCode,
    }
}

export class LogoutSessionRequest extends jspb.Message { 
    getWorldId(): number;
    setWorldId(value: number): LogoutSessionRequest;
    getAccountId(): number;
    setAccountId(value: number): LogoutSessionRequest;
    getDisconnectSource(): SessionDisconnectSource;
    setDisconnectSource(value: SessionDisconnectSource): LogoutSessionRequest;
    getTransferDisconnect(): boolean;
    setTransferDisconnect(value: boolean): LogoutSessionRequest;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): LogoutSessionRequest.AsObject;
    static toObject(includeInstance: boolean, msg: LogoutSessionRequest): LogoutSessionRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: LogoutSessionRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): LogoutSessionRequest;
    static deserializeBinaryFromReader(message: LogoutSessionRequest, reader: jspb.BinaryReader): LogoutSessionRequest;
}

export namespace LogoutSessionRequest {
    export type AsObject = {
        worldId: number,
        accountId: number,
        disconnectSource: SessionDisconnectSource,
        transferDisconnect: boolean,
    }
}

export class LogoutSessionReply extends jspb.Message { 
    getOk(): boolean;
    setOk(value: boolean): LogoutSessionReply;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): LogoutSessionReply.AsObject;
    static toObject(includeInstance: boolean, msg: LogoutSessionReply): LogoutSessionReply.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: LogoutSessionReply, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): LogoutSessionReply;
    static deserializeBinaryFromReader(message: LogoutSessionReply, reader: jspb.BinaryReader): LogoutSessionReply;
}

export namespace LogoutSessionReply {
    export type AsObject = {
        ok: boolean,
    }
}

export class PartyDoor extends jspb.Message { 
    getTown(): number;
    setTown(value: number): PartyDoor;
    getTarget(): number;
    setTarget(value: number): PartyDoor;
    getX(): number;
    setX(value: number): PartyDoor;
    getY(): number;
    setY(value: number): PartyDoor;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): PartyDoor.AsObject;
    static toObject(includeInstance: boolean, msg: PartyDoor): PartyDoor.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: PartyDoor, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): PartyDoor;
    static deserializeBinaryFromReader(message: PartyDoor, reader: jspb.BinaryReader): PartyDoor;
}

export namespace PartyDoor {
    export type AsObject = {
        town: number,
        target: number,
        x: number,
        y: number,
    }
}

export class PartyMember extends jspb.Message { 
    getWorldId(): number;
    setWorldId(value: number): PartyMember;
    getCharacterId(): number;
    setCharacterId(value: number): PartyMember;
    getCharacterName(): string;
    setCharacterName(value: string): PartyMember;
    getLevel(): number;
    setLevel(value: number): PartyMember;
    getClassId(): number;
    setClassId(value: number): PartyMember;
    getRole(): string;
    setRole(value: string): PartyMember;
    getMapId(): number;
    setMapId(value: number): PartyMember;

    hasChannelIndex(): boolean;
    clearChannelIndex(): void;
    getChannelIndex(): number | undefined;
    setChannelIndex(value: number): PartyMember;

    hasDoor(): boolean;
    clearDoor(): void;
    getDoor(): PartyDoor | undefined;
    setDoor(value?: PartyDoor): PartyMember;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): PartyMember.AsObject;
    static toObject(includeInstance: boolean, msg: PartyMember): PartyMember.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: PartyMember, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): PartyMember;
    static deserializeBinaryFromReader(message: PartyMember, reader: jspb.BinaryReader): PartyMember;
}

export namespace PartyMember {
    export type AsObject = {
        worldId: number,
        characterId: number,
        characterName: string,
        level: number,
        classId: number,
        role: string,
        mapId: number,
        channelIndex?: number,
        door?: PartyDoor.AsObject,
    }
}

export class CreatePartyRequest extends jspb.Message { 
    getWorldId(): number;
    setWorldId(value: number): CreatePartyRequest;

    hasLeader(): boolean;
    clearLeader(): void;
    getLeader(): PartyMember | undefined;
    setLeader(value?: PartyMember): CreatePartyRequest;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): CreatePartyRequest.AsObject;
    static toObject(includeInstance: boolean, msg: CreatePartyRequest): CreatePartyRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: CreatePartyRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): CreatePartyRequest;
    static deserializeBinaryFromReader(message: CreatePartyRequest, reader: jspb.BinaryReader): CreatePartyRequest;
}

export namespace CreatePartyRequest {
    export type AsObject = {
        worldId: number,
        leader?: PartyMember.AsObject,
    }
}

export class CreatePartyReply extends jspb.Message { 
    getOk(): boolean;
    setOk(value: boolean): CreatePartyReply;
    getErrorCode(): PartyErrorCode;
    setErrorCode(value: PartyErrorCode): CreatePartyReply;

    hasPartyId(): boolean;
    clearPartyId(): void;
    getPartyId(): number | undefined;
    setPartyId(value: number): CreatePartyReply;
    getRevision(): number;
    setRevision(value: number): CreatePartyReply;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): CreatePartyReply.AsObject;
    static toObject(includeInstance: boolean, msg: CreatePartyReply): CreatePartyReply.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: CreatePartyReply, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): CreatePartyReply;
    static deserializeBinaryFromReader(message: CreatePartyReply, reader: jspb.BinaryReader): CreatePartyReply;
}

export namespace CreatePartyReply {
    export type AsObject = {
        ok: boolean,
        errorCode: PartyErrorCode,
        partyId?: number,
        revision: number,
    }
}

export class JoinPartyRequest extends jspb.Message { 
    getWorldId(): number;
    setWorldId(value: number): JoinPartyRequest;
    getPartyId(): number;
    setPartyId(value: number): JoinPartyRequest;

    hasMember(): boolean;
    clearMember(): void;
    getMember(): PartyMember | undefined;
    setMember(value?: PartyMember): JoinPartyRequest;
    getSkipInvitePendingCheck(): boolean;
    setSkipInvitePendingCheck(value: boolean): JoinPartyRequest;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): JoinPartyRequest.AsObject;
    static toObject(includeInstance: boolean, msg: JoinPartyRequest): JoinPartyRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: JoinPartyRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): JoinPartyRequest;
    static deserializeBinaryFromReader(message: JoinPartyRequest, reader: jspb.BinaryReader): JoinPartyRequest;
}

export namespace JoinPartyRequest {
    export type AsObject = {
        worldId: number,
        partyId: number,
        member?: PartyMember.AsObject,
        skipInvitePendingCheck: boolean,
    }
}

export class JoinPartyReply extends jspb.Message { 
    getOk(): boolean;
    setOk(value: boolean): JoinPartyReply;
    getErrorCode(): PartyErrorCode;
    setErrorCode(value: PartyErrorCode): JoinPartyReply;

    hasPartyId(): boolean;
    clearPartyId(): void;
    getPartyId(): number | undefined;
    setPartyId(value: number): JoinPartyReply;
    getRevision(): number;
    setRevision(value: number): JoinPartyReply;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): JoinPartyReply.AsObject;
    static toObject(includeInstance: boolean, msg: JoinPartyReply): JoinPartyReply.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: JoinPartyReply, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): JoinPartyReply;
    static deserializeBinaryFromReader(message: JoinPartyReply, reader: jspb.BinaryReader): JoinPartyReply;
}

export namespace JoinPartyReply {
    export type AsObject = {
        ok: boolean,
        errorCode: PartyErrorCode,
        partyId?: number,
        revision: number,
    }
}

export class LeavePartyRequest extends jspb.Message { 
    getWorldId(): number;
    setWorldId(value: number): LeavePartyRequest;
    getCharacterId(): number;
    setCharacterId(value: number): LeavePartyRequest;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): LeavePartyRequest.AsObject;
    static toObject(includeInstance: boolean, msg: LeavePartyRequest): LeavePartyRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: LeavePartyRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): LeavePartyRequest;
    static deserializeBinaryFromReader(message: LeavePartyRequest, reader: jspb.BinaryReader): LeavePartyRequest;
}

export namespace LeavePartyRequest {
    export type AsObject = {
        worldId: number,
        characterId: number,
    }
}

export class LeavePartyReply extends jspb.Message { 
    getOk(): boolean;
    setOk(value: boolean): LeavePartyReply;
    getErrorCode(): PartyErrorCode;
    setErrorCode(value: PartyErrorCode): LeavePartyReply;

    hasPartyId(): boolean;
    clearPartyId(): void;
    getPartyId(): number | undefined;
    setPartyId(value: number): LeavePartyReply;
    getRevision(): number;
    setRevision(value: number): LeavePartyReply;
    getDisbanded(): boolean;
    setDisbanded(value: boolean): LeavePartyReply;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): LeavePartyReply.AsObject;
    static toObject(includeInstance: boolean, msg: LeavePartyReply): LeavePartyReply.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: LeavePartyReply, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): LeavePartyReply;
    static deserializeBinaryFromReader(message: LeavePartyReply, reader: jspb.BinaryReader): LeavePartyReply;
}

export namespace LeavePartyReply {
    export type AsObject = {
        ok: boolean,
        errorCode: PartyErrorCode,
        partyId?: number,
        revision: number,
        disbanded: boolean,
    }
}

export class ExpelPartyRequest extends jspb.Message { 
    getWorldId(): number;
    setWorldId(value: number): ExpelPartyRequest;
    getRequesterCharacterId(): number;
    setRequesterCharacterId(value: number): ExpelPartyRequest;
    getTargetCharacterId(): number;
    setTargetCharacterId(value: number): ExpelPartyRequest;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): ExpelPartyRequest.AsObject;
    static toObject(includeInstance: boolean, msg: ExpelPartyRequest): ExpelPartyRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: ExpelPartyRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): ExpelPartyRequest;
    static deserializeBinaryFromReader(message: ExpelPartyRequest, reader: jspb.BinaryReader): ExpelPartyRequest;
}

export namespace ExpelPartyRequest {
    export type AsObject = {
        worldId: number,
        requesterCharacterId: number,
        targetCharacterId: number,
    }
}

export class ExpelPartyReply extends jspb.Message { 
    getOk(): boolean;
    setOk(value: boolean): ExpelPartyReply;
    getErrorCode(): PartyErrorCode;
    setErrorCode(value: PartyErrorCode): ExpelPartyReply;

    hasPartyId(): boolean;
    clearPartyId(): void;
    getPartyId(): number | undefined;
    setPartyId(value: number): ExpelPartyReply;
    getRevision(): number;
    setRevision(value: number): ExpelPartyReply;
    getDisbanded(): boolean;
    setDisbanded(value: boolean): ExpelPartyReply;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): ExpelPartyReply.AsObject;
    static toObject(includeInstance: boolean, msg: ExpelPartyReply): ExpelPartyReply.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: ExpelPartyReply, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): ExpelPartyReply;
    static deserializeBinaryFromReader(message: ExpelPartyReply, reader: jspb.BinaryReader): ExpelPartyReply;
}

export namespace ExpelPartyReply {
    export type AsObject = {
        ok: boolean,
        errorCode: PartyErrorCode,
        partyId?: number,
        revision: number,
        disbanded: boolean,
    }
}

export class ChangePartyLeaderRequest extends jspb.Message { 
    getWorldId(): number;
    setWorldId(value: number): ChangePartyLeaderRequest;
    getPartyId(): number;
    setPartyId(value: number): ChangePartyLeaderRequest;
    getRequesterCharacterId(): number;
    setRequesterCharacterId(value: number): ChangePartyLeaderRequest;
    getNewLeaderCharacterId(): number;
    setNewLeaderCharacterId(value: number): ChangePartyLeaderRequest;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): ChangePartyLeaderRequest.AsObject;
    static toObject(includeInstance: boolean, msg: ChangePartyLeaderRequest): ChangePartyLeaderRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: ChangePartyLeaderRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): ChangePartyLeaderRequest;
    static deserializeBinaryFromReader(message: ChangePartyLeaderRequest, reader: jspb.BinaryReader): ChangePartyLeaderRequest;
}

export namespace ChangePartyLeaderRequest {
    export type AsObject = {
        worldId: number,
        partyId: number,
        requesterCharacterId: number,
        newLeaderCharacterId: number,
    }
}

export class ChangePartyLeaderReply extends jspb.Message { 
    getOk(): boolean;
    setOk(value: boolean): ChangePartyLeaderReply;
    getErrorCode(): PartyErrorCode;
    setErrorCode(value: PartyErrorCode): ChangePartyLeaderReply;

    hasPartyId(): boolean;
    clearPartyId(): void;
    getPartyId(): number | undefined;
    setPartyId(value: number): ChangePartyLeaderReply;
    getRevision(): number;
    setRevision(value: number): ChangePartyLeaderReply;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): ChangePartyLeaderReply.AsObject;
    static toObject(includeInstance: boolean, msg: ChangePartyLeaderReply): ChangePartyLeaderReply.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: ChangePartyLeaderReply, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): ChangePartyLeaderReply;
    static deserializeBinaryFromReader(message: ChangePartyLeaderReply, reader: jspb.BinaryReader): ChangePartyLeaderReply;
}

export namespace ChangePartyLeaderReply {
    export type AsObject = {
        ok: boolean,
        errorCode: PartyErrorCode,
        partyId?: number,
        revision: number,
    }
}

export class Party extends jspb.Message { 
    getWorldId(): number;
    setWorldId(value: number): Party;
    getPartyId(): number;
    setPartyId(value: number): Party;
    getLeaderCharacterId(): number;
    setLeaderCharacterId(value: number): Party;
    getRevision(): number;
    setRevision(value: number): Party;
    getState(): string;
    setState(value: string): Party;
    clearMembersList(): void;
    getMembersList(): Array<PartyMember>;
    setMembersList(value: Array<PartyMember>): Party;
    addMembers(value?: PartyMember, index?: number): PartyMember;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Party.AsObject;
    static toObject(includeInstance: boolean, msg: Party): Party.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: Party, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Party;
    static deserializeBinaryFromReader(message: Party, reader: jspb.BinaryReader): Party;
}

export namespace Party {
    export type AsObject = {
        worldId: number,
        partyId: number,
        leaderCharacterId: number,
        revision: number,
        state: string,
        membersList: Array<PartyMember.AsObject>,
    }
}

export class GetPartyRequest extends jspb.Message { 
    getWorldId(): number;
    setWorldId(value: number): GetPartyRequest;
    getPartyId(): number;
    setPartyId(value: number): GetPartyRequest;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): GetPartyRequest.AsObject;
    static toObject(includeInstance: boolean, msg: GetPartyRequest): GetPartyRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: GetPartyRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): GetPartyRequest;
    static deserializeBinaryFromReader(message: GetPartyRequest, reader: jspb.BinaryReader): GetPartyRequest;
}

export namespace GetPartyRequest {
    export type AsObject = {
        worldId: number,
        partyId: number,
    }
}

export class GetPartyReply extends jspb.Message { 
    getFound(): boolean;
    setFound(value: boolean): GetPartyReply;

    hasParty(): boolean;
    clearParty(): void;
    getParty(): Party | undefined;
    setParty(value?: Party): GetPartyReply;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): GetPartyReply.AsObject;
    static toObject(includeInstance: boolean, msg: GetPartyReply): GetPartyReply.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: GetPartyReply, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): GetPartyReply;
    static deserializeBinaryFromReader(message: GetPartyReply, reader: jspb.BinaryReader): GetPartyReply;
}

export namespace GetPartyReply {
    export type AsObject = {
        found: boolean,
        party?: Party.AsObject,
    }
}

export class UpdatePartyMemberRequest extends jspb.Message { 

    hasMember(): boolean;
    clearMember(): void;
    getMember(): PartyMember | undefined;
    setMember(value?: PartyMember): UpdatePartyMemberRequest;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): UpdatePartyMemberRequest.AsObject;
    static toObject(includeInstance: boolean, msg: UpdatePartyMemberRequest): UpdatePartyMemberRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: UpdatePartyMemberRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): UpdatePartyMemberRequest;
    static deserializeBinaryFromReader(message: UpdatePartyMemberRequest, reader: jspb.BinaryReader): UpdatePartyMemberRequest;
}

export namespace UpdatePartyMemberRequest {
    export type AsObject = {
        member?: PartyMember.AsObject,
    }
}

export class UpdatePartyMemberReply extends jspb.Message { 
    getOk(): boolean;
    setOk(value: boolean): UpdatePartyMemberReply;
    getErrorCode(): PartyErrorCode;
    setErrorCode(value: PartyErrorCode): UpdatePartyMemberReply;

    hasPartyId(): boolean;
    clearPartyId(): void;
    getPartyId(): number | undefined;
    setPartyId(value: number): UpdatePartyMemberReply;
    getRevision(): number;
    setRevision(value: number): UpdatePartyMemberReply;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): UpdatePartyMemberReply.AsObject;
    static toObject(includeInstance: boolean, msg: UpdatePartyMemberReply): UpdatePartyMemberReply.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: UpdatePartyMemberReply, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): UpdatePartyMemberReply;
    static deserializeBinaryFromReader(message: UpdatePartyMemberReply, reader: jspb.BinaryReader): UpdatePartyMemberReply;
}

export namespace UpdatePartyMemberReply {
    export type AsObject = {
        ok: boolean,
        errorCode: PartyErrorCode,
        partyId?: number,
        revision: number,
    }
}

export class InvitePartyRequest extends jspb.Message { 
    getWorldId(): number;
    setWorldId(value: number): InvitePartyRequest;
    getInviterCharacterId(): number;
    setInviterCharacterId(value: number): InvitePartyRequest;
    getTargetCharacterName(): string;
    setTargetCharacterName(value: string): InvitePartyRequest;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): InvitePartyRequest.AsObject;
    static toObject(includeInstance: boolean, msg: InvitePartyRequest): InvitePartyRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: InvitePartyRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): InvitePartyRequest;
    static deserializeBinaryFromReader(message: InvitePartyRequest, reader: jspb.BinaryReader): InvitePartyRequest;
}

export namespace InvitePartyRequest {
    export type AsObject = {
        worldId: number,
        inviterCharacterId: number,
        targetCharacterName: string,
    }
}

export class InvitePartyReply extends jspb.Message { 
    getOk(): boolean;
    setOk(value: boolean): InvitePartyReply;
    getErrorCode(): PartyErrorCode;
    setErrorCode(value: PartyErrorCode): InvitePartyReply;
    getTargetCharacterId(): number;
    setTargetCharacterId(value: number): InvitePartyReply;
    getTargetChannelId(): number;
    setTargetChannelId(value: number): InvitePartyReply;

    hasPartyId(): boolean;
    clearPartyId(): void;
    getPartyId(): number | undefined;
    setPartyId(value: number): InvitePartyReply;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): InvitePartyReply.AsObject;
    static toObject(includeInstance: boolean, msg: InvitePartyReply): InvitePartyReply.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: InvitePartyReply, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): InvitePartyReply;
    static deserializeBinaryFromReader(message: InvitePartyReply, reader: jspb.BinaryReader): InvitePartyReply;
}

export namespace InvitePartyReply {
    export type AsObject = {
        ok: boolean,
        errorCode: PartyErrorCode,
        targetCharacterId: number,
        targetChannelId: number,
        partyId?: number,
    }
}

export class DenyPartyRequest extends jspb.Message { 
    getWorldId(): number;
    setWorldId(value: number): DenyPartyRequest;
    getDeniedCharacterId(): number;
    setDeniedCharacterId(value: number): DenyPartyRequest;
    getInviterName(): string;
    setInviterName(value: string): DenyPartyRequest;
    getAction(): number;
    setAction(value: number): DenyPartyRequest;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): DenyPartyRequest.AsObject;
    static toObject(includeInstance: boolean, msg: DenyPartyRequest): DenyPartyRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: DenyPartyRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): DenyPartyRequest;
    static deserializeBinaryFromReader(message: DenyPartyRequest, reader: jspb.BinaryReader): DenyPartyRequest;
}

export namespace DenyPartyRequest {
    export type AsObject = {
        worldId: number,
        deniedCharacterId: number,
        inviterName: string,
        action: number,
    }
}

export class DenyPartyReply extends jspb.Message { 
    getOk(): boolean;
    setOk(value: boolean): DenyPartyReply;
    getErrorCode(): PartyErrorCode;
    setErrorCode(value: PartyErrorCode): DenyPartyReply;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): DenyPartyReply.AsObject;
    static toObject(includeInstance: boolean, msg: DenyPartyReply): DenyPartyReply.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: DenyPartyReply, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): DenyPartyReply;
    static deserializeBinaryFromReader(message: DenyPartyReply, reader: jspb.BinaryReader): DenyPartyReply;
}

export namespace DenyPartyReply {
    export type AsObject = {
        ok: boolean,
        errorCode: PartyErrorCode,
    }
}

export class BroadcastMultiChatRequest extends jspb.Message { 
    getWorldId(): number;
    setWorldId(value: number): BroadcastMultiChatRequest;
    getMemberId(): number;
    setMemberId(value: number): BroadcastMultiChatRequest;
    getSenderCharacterId(): number;
    setSenderCharacterId(value: number): BroadcastMultiChatRequest;
    getChatMode(): number;
    setChatMode(value: number): BroadcastMultiChatRequest;
    getSenderName(): string;
    setSenderName(value: string): BroadcastMultiChatRequest;
    getMessage(): string;
    setMessage(value: string): BroadcastMultiChatRequest;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): BroadcastMultiChatRequest.AsObject;
    static toObject(includeInstance: boolean, msg: BroadcastMultiChatRequest): BroadcastMultiChatRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: BroadcastMultiChatRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): BroadcastMultiChatRequest;
    static deserializeBinaryFromReader(message: BroadcastMultiChatRequest, reader: jspb.BinaryReader): BroadcastMultiChatRequest;
}

export namespace BroadcastMultiChatRequest {
    export type AsObject = {
        worldId: number,
        memberId: number,
        senderCharacterId: number,
        chatMode: number,
        senderName: string,
        message: string,
    }
}

export class BroadcastMultiChatReply extends jspb.Message { 
    getOk(): boolean;
    setOk(value: boolean): BroadcastMultiChatReply;
    getErrorCode(): PartyErrorCode;
    setErrorCode(value: PartyErrorCode): BroadcastMultiChatReply;
    getDeliveredCount(): number;
    setDeliveredCount(value: number): BroadcastMultiChatReply;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): BroadcastMultiChatReply.AsObject;
    static toObject(includeInstance: boolean, msg: BroadcastMultiChatReply): BroadcastMultiChatReply.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: BroadcastMultiChatReply, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): BroadcastMultiChatReply;
    static deserializeBinaryFromReader(message: BroadcastMultiChatReply, reader: jspb.BinaryReader): BroadcastMultiChatReply;
}

export namespace BroadcastMultiChatReply {
    export type AsObject = {
        ok: boolean,
        errorCode: PartyErrorCode,
        deliveredCount: number,
    }
}

export enum SessionErrorCode {
    SESSION_NONE = 0,
    SESSION_UNKNOWN = 1,
    SESSION_ALREADY_LOGGED_IN = 2,
    SESSION_NOT_FOUND = 3,
    SESSION_LOGOUT_FAILED = 4,
}

export enum SessionDisconnectSource {
    SESSION_DISCONNECT_SOURCE_UNSPECIFIED = 0,
    SESSION_DISCONNECT_SOURCE_LOGIN_SERVER = 1,
    SESSION_DISCONNECT_SOURCE_GAME_SERVER = 2,
}

export enum PartyErrorCode {
    NONE = 0,
    UNKNOWN = 1,
    ALREADY_IN_PARTY = 2,
    CHARACTER_NOT_FOUND = 3,
    PARTY_NOT_FOUND = 4,
    PARTY_FULL = 5,
    NOT_IN_PARTY = 6,
    NOT_PARTY_LEADER = 7,
    TARGET_NOT_IN_PARTY = 8,
    TARGET_ALREADY_LEADER = 9,
    CANNOT_EXPEL_SELF = 10,
    INVITE_EXPIRED_OR_INVALID = 11,
    TARGET_OFFLINE = 12,
    INVITER_NOT_IN_PARTY = 13,
    CANNOT_INVITE_SELF = 14,
    TARGET_ALREADY_IN_PARTY = 15,
}
