package wz

import (
	"slices"

	"github.com/boyism80/fm/services/game/constant"
)

type EquipmentCore struct {
	*ItemCore
	EquipmentType   constant.EquipmentType
	Parts           []constant.EquipmentPartsType
	Required        RequiredStats
	Ability         AbilityStats
	EnhanceChance   uint8
	EquipTradeBlock bool
	RoyalSpecial    bool
	MasterSpecial   bool
	Hide            bool
	AttackSpeed     int
	Pets            []uint32
	TamingMob       uint32
}

var equipmentPartsByCategory = map[uint32][]constant.EquipmentPartsType{
	100: {constant.EquipmentPartsCap},
	101: {constant.EquipmentPartsFace},
	102: {constant.EquipmentPartsEye},
	103: {constant.EquipmentPartsEar},
	104: {constant.EquipmentPartsTop},
	105: {constant.EquipmentPartsTop},
	106: {constant.EquipmentPartsPants},
	107: {constant.EquipmentPartsShoes},
	108: {constant.EquipmentPartsGlove},
	109: {constant.EquipmentPartsShield},
	110: {constant.EquipmentPartsCape},
	111: {constant.EquipmentPartsRing, constant.EquipmentPartsRing2, constant.EquipmentPartsRing3, constant.EquipmentPartsRing4},
	112: {constant.EquipmentPartsPendant},
	113: {constant.EquipmentPartsBelt},
	114: {constant.EquipmentPartsMedal},
	119: {constant.EquipmentPartsShield},
	180: {constant.EquipmentPartsPetEquip},
	190: {constant.EquipmentPartsTamingMob},
	191: {constant.EquipmentPartsSaddle},
	192: {constant.EquipmentPartsMobEquip},
}

func (core *EquipmentCore) setEquipmentType() {
	id := core.ID
	switch {
	case id >= 1300000 && id < 1800000:
		core.EquipmentType = constant.EquipmentTypeWeapon
	case id >= 1092000 && id < 1100000:
		core.EquipmentType = constant.EquipmentTypeShield
	case id >= 1000000 && id < 1010000:
		core.EquipmentType = constant.EquipmentTypeCap
	case id >= 1040000 && id < 1050000:
		core.EquipmentType = constant.EquipmentTypeCoat
	case id >= 1050000 && id < 1060000:
		core.EquipmentType = constant.EquipmentTypeLongcoat
	case id >= 1060000 && id < 1070000:
		core.EquipmentType = constant.EquipmentTypePants
	case id >= 1070000 && id < 1080000:
		core.EquipmentType = constant.EquipmentTypeShoes
	case id >= 1080000 && id < 1090000:
		core.EquipmentType = constant.EquipmentTypeGlove
	case id >= 1102000 && id < 1104000:
		core.EquipmentType = constant.EquipmentTypeCape
	case id >= 1112000 && id < 1120000:
		core.EquipmentType = constant.EquipmentTypeRing
	case id >= 20000 && id < 22000:
		core.EquipmentType = constant.EquipmentTypeFace
	case (id >= 1132000 && id < 1183000) || (id >= 1010000 && id < 1040000) || (id >= 1122000 && id < 1123000):
		core.EquipmentType = constant.EquipmentTypeAccessory
	case id >= 30000 && id < 35000:
		core.EquipmentType = constant.EquipmentTypeHair
	case id >= 1172000 && id < 1173000:
		core.EquipmentType = constant.EquipmentTypeMonsterBook
	case id >= 1662000 && id < 1680000:
		core.EquipmentType = constant.EquipmentTypeAndroid
	case id >= 1920000 && id < 2000000:
		core.EquipmentType = constant.EquipmentTypeDragon
	case id >= 1802000 && id < 1820000:
		core.EquipmentType = constant.EquipmentTypePetEquip
	case id >= 1900000 && id < 1920000:
		core.EquipmentType = constant.EquipmentTypeTaming
	case id >= 1610000 && id < 1660000:
		core.EquipmentType = constant.EquipmentTypeMechanic
	default:
		core.EquipmentType = constant.EquipmentTypeUnknown
	}

	category := id / 10000
	if parts, ok := equipmentPartsByCategory[category]; ok {
		core.Parts = parts
		return
	}
	switch category / 10 {
	case 13, 14, 16, 17:
		core.Parts = []constant.EquipmentPartsType{constant.EquipmentPartsWeapon}
	}
}

func (core *EquipmentCore) CanEquipAt(parts constant.EquipmentPartsType) bool {
	if parts <= 2*constant.EquipmentPartsCashOffset || parts >= 0 {
		return false
	}
	if parts <= constant.EquipmentPartsCashOffset {
		parts -= constant.EquipmentPartsCashOffset
	}
	if parts >= 0 {
		return false
	}

	return slices.Contains(core.Parts, parts)
}

func (core *EquipmentCore) GetEquipmentType() constant.EquipmentType { return core.EquipmentType }
func (core *EquipmentCore) GetParts() []constant.EquipmentPartsType  { return core.Parts }
func (core *EquipmentCore) GetEnhanceChance() uint8                  { return core.EnhanceChance }
func (core *EquipmentCore) GetRequired() RequiredStats {
	return core.Required
}
func (core *EquipmentCore) GetAbility() AbilityStats { return core.Ability }
func (core *EquipmentCore) IsEquipTradeBlock() bool  { return core.EquipTradeBlock }
