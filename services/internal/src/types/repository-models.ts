import type { PartyMemberRole, PartyState, GuildMemberRank } from "../protobuf/generated/fminternal/internal_service";
import type { EquipmentBonusStatsJson } from "./equipment-bonus-stats";
import type { KeyLayoutJsonRecord } from "./key-layout-json";
import type { GuildLogo as GuildLogoModel, GuildRankTitles as GuildRankTitlesModel } from "./guild-json";

export interface AccountModel {
    accountId: number;
    loginId: string;
    passwordHash: string;
    gender: number;
    role?: number;
    isBanned: boolean;
    banReason?: string | null;
    isChatBlocked: boolean;
    chatBlockedUntil?: string | null;
    characterSlotCount: number;
    lastLoginIp?: string | null;
    macAddress?: string | null;
}

export type AccountRow = {
    id: number;
    login_id: string;
    password_hash: string;
    gender: number;
    role: number;
    is_banned: boolean;
    ban_reason?: string | null;
    is_chat_blocked: boolean;
    chat_blocked_until?: string | null;
    character_slot_count: number;
    last_login_ip?: string | null;
    mac_address?: string | null;
};

export type AccountDeleteModel = { accountId: number };

export interface CharacterModel {
    characterId: number;
    accountId: number;
    worldId: number;
    name: string;
    gender: number;
    skinColor: number;
    face: number;
    hair: number;
    level: number;
    classId: number;
    role?: number;
    str?: number;
    dex?: number;
    intStat?: number;
    luk?: number;
    hp?: number;
    maxHp?: number;
    mp?: number;
    maxMp?: number;
    abilityPoint?: number;
    exp?: number;
    mapId: number;
    spawnPoint: number;
    positionX?: number;
    positionY?: number;
    stance?: number;
    meso?: number;
    skillPoint?: number;
    hidden?: boolean;
    updatedAt?: Date;
}

export type CharacterRow = {
    id: number;
    account_id: number;
    world_id: number;
    name: string;
    gender: number;
    skin_color: number;
    face: number;
    hair: number;
    level: number;
    class_id: number;
    role: number;
    str: number;
    dex: number;
    int_stat: number;
    luk: number;
    hp: number;
    max_hp: number;
    mp: number;
    max_mp: number;
    ability_point: number;
    exp: number;
    map_id: number;
    spawn_point: number;
    pos_x: number;
    pos_y: number;
    stance: number;
    meso: number;
    skill_point: number;
    hidden: boolean;
    deleted?: boolean;
    created_at?: Date | string;
    updated_at?: Date | string;
};

export type CharacterDeleteModel = { worldId: number; characterId: number; accountId: number };

export interface CharacterOverviewModel {
    characterId: number;
    accountId: number;
    worldId: number;
    name: string;
    gender: number;
    skinColor: number;
    face: number;
    hair: number;
    level: number;
    classId: number;
    mapId: number;
    spawnPoint: number;
    rank: number;
    rankDiff: number;
    classRank: number;
    classRankDiff: number;
    baseLooks: Record<string, number> | Record<number, number>;
    overlays: Record<string, number> | Record<number, number>;
}

export type CharacterOverviewRow = {
    character_id: number;
    account_id: number;
    world_id: number;
    name: string;
    gender: number;
    skin_color: number;
    face: number;
    hair: number;
    level: number;
    class_id: number;
    map_id: number;
    spawn_point: number;
    rank: number;
    rank_diff: number;
    class_rank: number;
    class_rank_diff: number;
    base_looks: string | Record<string, number> | Record<number, number>;
    overlays: string | Record<string, number> | Record<number, number>;
    deleted?: boolean;
    updated_at?: Date | string;
};

export interface CharacterRealtimeStateModel {
    worldId: number;
    characterId: number;
    partyId?: number | null;
    guildId?: number | null;
    buddyCapacity?: number;
    updatedAt?: Date;
}

export type CharacterRealtimeStateRow = {
    world_id: number;
    character_id: number;
    party_id: number | null;
    guild_id: number | null;
    buddy_capacity?: number | null;
    updated_at?: Date | string;
};

export type CharacterRealtimeStateDeleteRow = { worldId: number; characterId: number };

export interface InventoryModel {
    uniqueId: number | null;
    ownerId: number;
    inventoryType: number;
    itemId: number;
    slot: number;
    count: number;
    expiration: Date | null;
    enhanceChance: number | null;
    enhanceCount: number | null;
    flag: number | null;
    skillBonus: number | null;
    ownerName: string | null;
    equipBonusStats?: EquipmentBonusStatsJson;
    updatedAt?: Date;
}

export type InventoryRow = {
    unique_id: number | null;
    owner_id: number;
    inventory_type: number;
    item_id: number;
    slot: number;
    count: number;
    expiration: Date | string | null;
    enhance_chance: number | null;
    enhance_count: number | null;
    flag: number | null;
    skill_bonus: number | null;
    owner_name: string | null;
    equip_bonus_stats: string | EquipmentBonusStatsJson;
    created_at?: Date | string;
    updated_at?: Date | string;
    deleted?: boolean;
};

export interface KeyLayoutModel {
    characterId: number;
    worldId: number;
    keyLayoutJson: string;
    updatedAt?: Date;
}

export type KeyLayoutRow = {
    character_id: number;
    world_id: number;
    key_layout_json: string | KeyLayoutJsonRecord;
    updated_at?: Date | string;
    deleted?: boolean;
};

export type KeyLayoutDeleteModel = { worldId: number; characterId: number };

export interface PartyModel {
    worldId: number;
    partyId: number;
    leaderCharacterId: number;
    state: PartyState;
    revision: number;
    disbandedAt?: Date | null;
    createdAt?: Date;
    updatedAt?: Date;
}

export type PartyRow = {
    world_id: number;
    party_id: number;
    leader_character_id: number;
    state: number | null;
    revision: number;
    disbanded_at?: Date | string | null;
    created_at?: Date | string;
    updated_at?: Date | string;
    deleted?: boolean;
};

export type PartyDeleteRow = { worldId: number; partyId: number };

export interface PartyMemberModel {
    worldId: number;
    partyId: number;
    characterId: number;
    characterName: string;
    level: number;
    classId: number;
    role: PartyMemberRole;
    mapId: number;
    door: { town: number; target: number; x: number; y: number } | null;
    joinedAt?: Date;
    updatedAt?: Date;
}

export type PartyMemberRow = {
    world_id: number;
    party_id: number;
    character_id: number;
    character_name: string;
    level: number;
    class_id: number;
    role: PartyMemberRole | null;
    map_id: number | null;
    door: { town: number; target: number; x: number; y: number } | null;
    joined_at: Date | string;
    updated_at: Date | string;
    deleted?: boolean;
};

export interface GuildModel {
    worldId: number;
    guildId: number;
    name: string;
    leaderCharacterId: number;
    gp: number;
    capacity: number;
    notice: string;
    logo: GuildLogoModel;
    rankTitles: GuildRankTitlesModel;
    allianceId?: number | null;
    revision: number;
    disbandedAt?: Date | null;
    createdAt?: Date;
    updatedAt?: Date;
}

export type GuildRow = {
    world_id: number;
    guild_id: number;
    name: string;
    leader_character_id: number;
    gp: number;
    capacity: number;
    notice: string;
    logo: GuildLogoModel | null;
    rank_titles: GuildRankTitlesModel;
    alliance_id?: number | null;
    revision: number;
    disbanded_at?: Date | string | null;
    created_at?: Date | string;
    updated_at?: Date | string;
    deleted?: boolean;
};

export type GuildDeleteRow = { worldId: number; guildId: number };

export interface GuildMemberModel {
    worldId: number;
    guildId: number;
    characterId: number;
    characterName: string;
    level: number;
    classId: number;
    guildRank: GuildMemberRank;
    allianceRank?: number | null;
    joinedAt?: Date;
    updatedAt?: Date;
}

export type GuildMemberRow = {
    world_id: number;
    guild_id: number;
    character_id: number;
    character_name: string;
    level: number;
    class_id: number;
    guild_rank: GuildMemberRank | number;
    alliance_rank?: number | null;
    joined_at: Date | string;
    updated_at: Date | string;
    deleted?: boolean;
};

export interface AllianceModel {
    worldId: number;
    allianceId: number;
    name: string;
    leaderCharacterId: number;
    guildIds: number[];
    rankTitles: import("./alliance-json").AllianceRankTitles;
    capacity: number;
    notice: string;
    revision: number;
    disbandedAt?: Date | null;
    createdAt?: Date;
    updatedAt?: Date;
}

export type AllianceRow = {
    world_id: number;
    alliance_id: number | null;
    name: string;
    leader_character_id: number;
    guild_ids: number[];
    rank_titles: import("./alliance-json").AllianceRankTitles;
    capacity: number;
    notice: string;
    revision: number;
    disbanded_at?: Date | string | null;
    created_at?: Date | string;
    updated_at?: Date | string;
    deleted?: boolean;
};

export interface GuildBulletinBoardThreadModel {
    worldId: number;
    guildId: number;
    localThreadId: number;
    posterCharacterId: number;
    title: string;
    body: string;
    icon: number;
    createdAt: Date;
    updatedAt: Date;
    replyCount?: number;
}

export type GuildBulletinBoardThreadRow = {
    world_id: number;
    guild_id: number;
    local_thread_id: number;
    poster_character_id: number;
    title: string;
    body: string;
    icon: number;
    created_at: Date | string;
    updated_at: Date | string;
    reply_count?: number | string;
};

export interface GuildBulletinBoardReplyModel {
    worldId: number;
    guildId: number;
    localThreadId: number;
    replyId: number;
    posterCharacterId: number;
    content: string;
    createdAt: Date;
}

export interface GuildBulletinBoardThreadWithRepliesModel {
    thread: GuildBulletinBoardThreadModel;
    replies: GuildBulletinBoardReplyModel[];
}

export type GuildBulletinBoardReplyRow = {
    world_id: number;
    guild_id: number;
    local_thread_id: number;
    reply_id: number;
    poster_character_id: number;
    content: string;
    created_at: Date | string;
};

export interface SkillModel {
    characterId: number;
    skillId: number;
    level: number;
    masterLevel: number;
    cooldownEndUnixMs: number | null;
    updatedAt?: Date;
}

export type SkillRow = {
    character_id: number;
    skill_id: number;
    level: number;
    master_level: number;
    cooldown_end_unix_ms: number | null;
    updated_at?: Date | string;
    deleted?: boolean;
};

export type BuffFlagValueModel = {
    mask: number;
    position: number;
    value: number;
};

export interface BuffModel {
    characterId: number;
    buffSourceId: number;
    kind: number;
    flagValues: BuffFlagValueModel[];
    remainingDurationMs: number | null;
    skillLevel: number | null;
    causerId: number | null;
    updatedAt?: Date;
}

export type BuffRow = {
    character_id: number;
    buff_source_id: number;
    kind: number;
    remaining_duration_ms: number | null;
    skill_level: number | null;
    causer_id: number | null;
    flag_values: BuffFlagValueRow[] | string;
    updated_at?: Date | string;
    deleted?: boolean;
};

export interface QuestModel {
    characterId: number;
    questId: number;
    status: number;
    mobKills: Record<string, number>;
    statusRecord: string;
    unknown2: Record<string, string>;
    completionTimeUnixMs: number;
    forfeited: number;
    updatedAt?: Date;
}

export type QuestRow = {
    character_id: number;
    quest_id: number;
    status: number;
    mob_kills: Record<string, number> | string;
    status_record: string;
    unknown2: Record<string, string> | string;
    completion_time_unix_ms: number | null;
    forfeited: number;
    updated_at?: Date | string;
    deleted?: boolean;
};

export type BuffFlagValueRow = {
    mask: number;
    position: number;
    value: number;
};

export interface CharacterBuddyModel {
    characterId: number;
    buddyCharacterId: number;
    groupName: string;
    pending: boolean;
    createdAt?: Date;
    updatedAt?: Date;
}

export type CharacterBuddyRow = {
    character_id: number;
    buddy_character_id: number;
    group_name: string;
    pending: boolean;
    created_at?: Date | string;
    updated_at?: Date | string;
};
