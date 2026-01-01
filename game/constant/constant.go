// Package constant provides MapleStory game constants and enumerations.
// This file contains type definitions and constants used throughout the game server.
package constant

import "time"

// ItemType represents the general category of an item.
type ItemType int

const (
	ITEM_TYPE_EQUIPMENT ItemType = 1 // Equippable items
	ITEM_TYPE_ETC       ItemType = 2 // Miscellaneous items
)

// EquipmentPartsType represents the equipment slot position.
type EquipmentPartsType int16

const (
	EQUIPMENT_PARTS_WEAPON EquipmentPartsType = -11 // Weapon slot
	EQUIPMENT_PARTS_SHIELD EquipmentPartsType = -10 // Shield slot
	EQUIPMENT_PARTS_TOP    EquipmentPartsType = -5  // Top clothing slot
	EQUIPMENT_PARTS_PANTS  EquipmentPartsType = -6  // Pants slot
)

// InventoryType represents the different inventory tabs.
type InventoryType int8

const (
	INVENTORY_TYPE_EQUIPMENT    InventoryType = 1 // Equipment inventory
	INVENTORY_TYPE_CONSUME      InventoryType = 2 // Consumable items
	INVENTORY_TYPE_INSTALLATION InventoryType = 3 // Installation items
	INVENTORY_TYPE_ETC          InventoryType = 4 // Miscellaneous items
	INVENTORY_TYPE_CASH         InventoryType = 5 // Cash shop items
)

// DropItemAnimationType represents the animation when items are dropped.
type DropItemAnimationType uint8

const (
	DROP_ITEM_ANIMATION_TYPE_LOOTING        DropItemAnimationType = 1 // Item being looted
	DROP_ITEM_ANIMATION_TYPE_NONE           DropItemAnimationType = 2 // No animation
	DROP_ITEM_ANIMATION_TYPE_DISAPPEAR_FADE DropItemAnimationType = 3 // Fade out disappear
	DROP_ITEM_ANIMATION_TYPE_DISAPPEAR      DropItemAnimationType = 4 // Instant disappear
)

// DropType represents who can pick up a dropped item.
type DropType uint8

const (
	DROP_TYPE_FFA   DropType = 0 // Free for all
	DROP_TYPE_PARTY DropType = 1 // Party members only
	DROP_TYPE_OWNED DropType = 2 // Owner only
)

// Stat represents character statistics as bit flags.
type Stat uint32

const (
	STAT_SKIN         Stat = 0x1     // Skin color
	STAT_FACE         Stat = 0x2     // Face ID
	STAT_HAIR         Stat = 0x4     // Hair ID
	STAT_PET          Stat = 0x8     // Pet data
	STAT_LEVEL        Stat = 0x10    // Character level
	STAT_JOB          Stat = 0x20    // Job class
	STAT_STR          Stat = 0x40    // Strength
	STAT_DEX          Stat = 0x80    // Dexterity
	STAT_INT          Stat = 0x100   // Intelligence
	STAT_LUK          Stat = 0x200   // Luck
	STAT_HP           Stat = 0x400   // Current HP
	STAT_MAX_HP       Stat = 0x800   // Maximum HP
	STAT_MP           Stat = 0x1000  // Current MP
	STAT_MAX_MP       Stat = 0x2000  // Maximum MP
	STAT_AVAILABLE_AP Stat = 0x4000  // Available ability points
	STAT_AVAILABLE_SP Stat = 0x8000  // Available skill points
	STAT_EXP          Stat = 0x10000 // Experience points
	STAT_FAME         Stat = 0x20000 // Fame points
	STAT_MESO         Stat = 0x40000 // Meso currency
)

// DialogType represents the type of NPC dialog interaction.
type DialogType uint8

const (
	DIALOG_TYPE_DEFAULT       DialogType = 0  // Simple message dialog
	DIALOG_TYPE_YES_NO        DialogType = 1  // Yes/No confirmation
	DIALOG_TYPE_INPUT         DialogType = 2  // Text input dialog
	DIALOG_TYPE_LIST          DialogType = 4  // Selection list
	DIALOG_TYPE_ACCEPT_ESCAPE DialogType = 11 // Accept or escape
	DIALOG_TYPE_ACCEPT        DialogType = 12 // Accept only
)

// ObjectType represents the type of game object as bit flags.
type ObjectType uint8

const (
	OBJECT_TYPE_ITEM      ObjectType = 1                                                                            // Dropped items
	OBJECT_TYPE_NPC       ObjectType = 2                                                                            // Non-player characters
	OBJECT_TYPE_MOB       ObjectType = 4                                                                            // Monsters
	OBJECT_TYPE_CHARACTER ObjectType = 8                                                                            // Player characters
	OBJECT_TYPE_LIFE      ObjectType = OBJECT_TYPE_MOB | OBJECT_TYPE_CHARACTER                                      // Living entities
	OBJECT_TYPE_ALL       ObjectType = OBJECT_TYPE_ITEM | OBJECT_TYPE_NPC | OBJECT_TYPE_MOB | OBJECT_TYPE_CHARACTER // All object types
)

// MobSpawnType represents how a monster spawns.
type MobSpawnType int8

const (
	MOB_SPAWN_TYPE_NONE    MobSpawnType = -1 // No spawn animation
	MOB_SPAWN_TYPE_ANIMATE MobSpawnType = -2 // With spawn animation
)

// MobDieAnimationType represents how a monster dies.
type MobDieAnimationType uint8

const (
	MOB_DIE_ANIMATION_TYPE_DISAPPEAR MobDieAnimationType = 0 // Instant disappear
	MOB_DIE_ANIMATION_TYPE_FADE_OUT  MobDieAnimationType = 1 // Fade out animation
)

// DefaultMobSpawnTime is the default respawn time for monsters.
var DefaultMobSpawnTime = 5 * time.Second

type RemoveItemType uint8

const (
	REMOVE_ITEM_TYPE_EXPIRED RemoveItemType = iota
	REMOVE_ITEM_TYPE_NO_ANIMATED
	REMOVE_ITEM_TYPE_ANIMATED
	REMOVE_ITEM_TYPE_EXPLOSION
	REMOVE_ITEM_TYPE_LOOT_BY_PET
)

// StatType represents the stat type for AP distribution
type StatType uint32

const (
	STAT_TYPE_STR StatType = 64   // Strength
	STAT_TYPE_DEX StatType = 128  // Dexterity
	STAT_TYPE_INT StatType = 256  // Intelligence
	STAT_TYPE_LUK StatType = 512  // Luck
	STAT_TYPE_HP  StatType = 2048 // Maximum HP
	STAT_TYPE_MP  StatType = 8192 // Maximum MP
)

// Stat limits
const (
	STAT_MAX_STR_DEX_INT_LUK uint16 = 999   // Maximum value for STR, DEX, INT, LUK
	STAT_MAX_HP_MP           uint16 = 30000 // Maximum value for HP and MP
	HP_AP_USED_MAX           uint16 = 10000 // Maximum HP/MP AP usage count
)

// Job ranges for stat calculation
const (
	JOB_BEGINNER_MIN uint16 = 0
	JOB_BEGINNER_1   uint16 = 1000
	JOB_BEGINNER_2   uint16 = 2000
	JOB_WARRIOR_MIN  uint16 = 100
	JOB_WARRIOR_MAX  uint16 = 132
	JOB_MAGICIAN_MIN uint16 = 200
	JOB_MAGICIAN_MAX uint16 = 232
	JOB_BOWMAN_MIN   uint16 = 300
	JOB_BOWMAN_MAX   uint16 = 322
	JOB_THIEF_MIN    uint16 = 400
	JOB_THIEF_MAX    uint16 = 422
)
