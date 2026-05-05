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
    is_banned: boolean | number;
    ban_reason?: string | null;
    is_chat_blocked: boolean | number;
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
    hidden: boolean | number;
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
    updatedAt?: Date;
}

export type CharacterRealtimeStateRow = {
    world_id: number;
    character_id: number;
    party_id: number | null;
    guild_id: number | null;
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
    equipBonusStats?: Record<string, unknown>;
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
    equip_bonus_stats: string | Record<string, unknown>;
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
    key_layout_json: string | Record<string, unknown>;
    updated_at?: Date | string;
    deleted?: boolean;
};

export type KeyLayoutDeleteModel = { worldId: number; characterId: number };

export interface PartyModel {
    worldId: number;
    partyId: number;
    leaderCharacterId: number;
    state: string;
    revision: number;
    disbandedAt?: Date | null;
    createdAt?: Date;
    updatedAt?: Date;
}

export type PartyRow = {
    world_id: number;
    party_id: number;
    leader_character_id: number;
    state: string;
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
    role: string;
    mapId: number;
    channelIndex: number;
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
    role: string | null;
    map_id: number | null;
    channel_index: number | null;
    door: string | null;
    joined_at: Date | string;
    updated_at: Date | string;
    deleted?: boolean;
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
