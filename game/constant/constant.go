// Package constant provides MapleStory game constants and enumerations.
// This file contains type definitions and constants used throughout the game server.
package constant

import "time"

// CharacterRole represents the character's permission role.
// Values are ordered: higher value means higher privilege. Use HasRoleAtLeast to compare.
type CharacterRole uint8

const (
	RoleUser  CharacterRole = iota // 일반유저
	RoleAdmin                      // 관리자
)

// ItemType represents the general category of an item.
type ItemType int

const (
	ITEM_TYPE_EQUIPMENT ItemType = 1 // Equippable items
	ITEM_TYPE_ETC       ItemType = 2 // Miscellaneous items
)

// EquipmentPartsType represents the equipment slot position.
type EquipmentPartsType int16

const (
	EQUIPMENT_PARTS_CAP    EquipmentPartsType = -1
	EQUIPMENT_PARTS_FACE   EquipmentPartsType = -2
	EQUIPMENT_PARTS_EYE    EquipmentPartsType = -3
	EQUIPMENT_PARTS_EAR    EquipmentPartsType = -4
	EQUIPMENT_PARTS_TOP    EquipmentPartsType = -5
	EQUIPMENT_PARTS_PANTS  EquipmentPartsType = -6
	EQUIPMENT_PARTS_SHOES  EquipmentPartsType = -7
	EQUIPMENT_PARTS_GLOVE  EquipmentPartsType = -8
	EQUIPMENT_PARTS_CAPE   EquipmentPartsType = -9
	EQUIPMENT_PARTS_SHIELD EquipmentPartsType = -10
	EQUIPMENT_PARTS_WEAPON EquipmentPartsType = -11
	EQUIPMENT_PARTS_RING   EquipmentPartsType = -12
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
	STAT_CLASS        Stat = 0x20    // Job class
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

// Item cleanup times
const (
	ITEM_EXPIRE_TIME = 120 * time.Second // Time before item/meso expires
	ITEM_FFA_TIME    = 30 * time.Second  // Time before owned/party drop becomes FFA
)

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

// Class ranges for stat calculation
const (
	CLASS_BEGINNER_MIN uint16 = 0
	CLASS_BEGINNER_1   uint16 = 1000
	CLASS_BEGINNER_2   uint16 = 2000
	CLASS_WARRIOR_MIN  uint16 = 100
	CLASS_WARRIOR_MAX  uint16 = 132
	CLASS_MAGICIAN_MIN uint16 = 200
	CLASS_MAGICIAN_MAX uint16 = 232
	CLASS_BOWMAN_MIN   uint16 = 300
	CLASS_BOWMAN_MAX   uint16 = 322
	CLASS_THIEF_MIN    uint16 = 400
	CLASS_THIEF_MAX    uint16 = 422
)

// Item category IDs (itemID / 10000)
const (
	ITEM_CATEGORY_SHURIKEN uint32 = 207
	ITEM_CATEGORY_BULLET   uint32 = 233
)

// Rechargeable item IDs
const (
	ITEM_SHURIKEN_BASE uint32 = 2070000
	ITEM_BULLET_BASE   uint32 = 2330000
)

// Rechargeable item ID ranges
var RechargeableShurikens = []uint32{
	2070000, 2070001, 2070002, 2070003, 2070004, 2070005,
	2070006, 2070007, 2070008, 2070009, 2070010, 2070011,
	2070012, 2070013,
}

var RechargeableBullets = []uint32{
	2330000, 2330001, 2330002, 2330003, 2330004, 2330005,
	2331000, 2332000,
}

// GetInventoryTypeByItemID returns the inventory type for an item by its ID.
// Uses MapleStory item ID ranges (itemID / 10000).
func GetInventoryTypeByItemID(itemID uint32) InventoryType {
	itemType := itemID / 10000
	switch {
	case itemType >= 100 && itemType < 200:
		return INVENTORY_TYPE_EQUIPMENT
	case itemType >= 200 && itemType < 300:
		return INVENTORY_TYPE_CONSUME
	case itemType >= 300 && itemType < 400:
		return INVENTORY_TYPE_INSTALLATION
	case itemType >= 400 && itemType < 500:
		return INVENTORY_TYPE_ETC
	case itemType >= 500 && itemType < 600:
		return INVENTORY_TYPE_CASH
	default:
		return INVENTORY_TYPE_ETC
	}
}
