// Package data provides MapleStory game data specifications and types.
// This file contains item-related data structures and interfaces.
package data

// ItemSpec defines the common interface for all item types in MapleStory.
type ItemSpec interface {
	GetID() uint32
	GetName() string
	GetPrice() int
	IsCash() bool
	IsQuest() bool
	GetCapacity() uint16
	IsTradeAvailable() int
}

// BasicStats represents the four basic character stats in MapleStory.
type BasicStats struct {
	Str uint16 // Strength
	Dex uint16 // Dexterity
	Int uint16 // Intelligence
	Luk uint16 // Luck
}

// RequiredStats represents the requirements to equip an item.
type RequiredStats struct {
	BasicStats
	Level uint8 // Required character level
	Class int   // Required job class
}

// AbilityStats represents the stat bonuses provided by equipment.
type AbilityStats struct {
	BasicStats
	PAD       uint16 // Physical attack damage
	MAD       uint16 // Magical attack damage
	PDD       uint16 // Physical defense damage
	MDD       uint16 // Magical defense damage
	Speed     uint16
	Jump      uint16
	ACC       uint16
	EVA       int
	MaxHP     uint16
	MaxMP     uint16
	PVPDamage int

	Avoid uint16
	Hands uint16
}

// ItemCoreSpec contains the common properties shared by all item types.
type ItemCoreSpec struct {
	ID             uint32
	Name           string
	Price          int
	Cash           bool
	Quest          bool
	SlotMax        uint16
	TradeAvailable int
}

// EquipmentSpec represents equippable items like weapons and armor.
type EquipmentSpec struct {
	*ItemCoreSpec
	Required        RequiredStats
	Ability         AbilityStats
	TUC             uint8 // Total upgrade count
	EquipTradeBlock bool
	RoyalSpecial    bool
	MasterSpecial   bool
	Hide            bool
	AttackSpeed     int
}

// PetSpec represents pet items that can be summoned.
type PetSpec struct {
	*ItemCoreSpec
}

// GeneralItemSpec represents general use items.
type GeneralItemSpec struct {
	*ItemCoreSpec
}

// CashItemSpec represents cash shop items.
type CashItemSpec struct {
	*ItemCoreSpec
}

// ActiveEffect represents HP/MP restoration effects.
type ActiveEffect struct {
	HP int
	MP int
}

// ConsumeSpec represents consumable items like potions.
type ConsumeSpec struct {
	*ItemCoreSpec
	ActiveEffect ActiveEffect
}

// InstallationSpec represents setup items like chairs.
type InstallationSpec struct {
	*ItemCoreSpec
}

// SpecialItemSpec represents special event or quest items.
type SpecialItemSpec struct {
	*ItemCoreSpec
	ActiveEffect ActiveEffect
}

// GetID returns the item's unique identifier.
func (spec *ItemCoreSpec) GetID() uint32 {
	return spec.ID
}

// GetName returns the item's display name.
func (spec *ItemCoreSpec) GetName() string {
	return spec.Name
}

// GetPrice returns the item's NPC shop price.
func (spec *ItemCoreSpec) GetPrice() int {
	return spec.Price
}

// IsCash returns true if this is a cash shop item.
func (spec *ItemCoreSpec) IsCash() bool {
	return spec.Cash
}

// IsQuest returns true if this is a quest-related item.
func (spec *ItemCoreSpec) IsQuest() bool {
	return spec.Quest
}

// GetCapacity returns the maximum stack size for this item.
func (spec *ItemCoreSpec) GetCapacity() uint16 {
	return max(1, spec.SlotMax)
}

// IsTradeAvailable returns the trade availability status.
func (spec *ItemCoreSpec) IsTradeAvailable() int {
	return spec.TradeAvailable
}
