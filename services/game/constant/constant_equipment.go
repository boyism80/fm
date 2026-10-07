package constant

import "slices"

type EquipmentPartsType int16

const (
	EquipmentPartsCap       EquipmentPartsType = -1
	EquipmentPartsFace      EquipmentPartsType = -2
	EquipmentPartsEye       EquipmentPartsType = -3
	EquipmentPartsEar       EquipmentPartsType = -4
	EquipmentPartsTop       EquipmentPartsType = -5
	EquipmentPartsPants     EquipmentPartsType = -6
	EquipmentPartsShoes     EquipmentPartsType = -7
	EquipmentPartsGlove     EquipmentPartsType = -8
	EquipmentPartsCape      EquipmentPartsType = -9
	EquipmentPartsShield    EquipmentPartsType = -10
	EquipmentPartsWeapon    EquipmentPartsType = -11
	EquipmentPartsRing      EquipmentPartsType = -12
	EquipmentPartsRing2     EquipmentPartsType = -13
	EquipmentPartsPetEquip  EquipmentPartsType = -14
	EquipmentPartsRing3     EquipmentPartsType = -15
	EquipmentPartsRing4     EquipmentPartsType = -16
	EquipmentPartsPendant   EquipmentPartsType = -17
	EquipmentPartsTamingMob EquipmentPartsType = -18
	EquipmentPartsSaddle    EquipmentPartsType = -19
	EquipmentPartsMobEquip  EquipmentPartsType = -20
	EquipmentPartsMedal     EquipmentPartsType = -21
	EquipmentPartsBelt      EquipmentPartsType = -22
)

const equipmentPartsCashOffset EquipmentPartsType = -100

func (p EquipmentPartsType) Cash() EquipmentPartsType {
	return p + equipmentPartsCashOffset
}

var equipmentPartsByCategory = map[uint32][]EquipmentPartsType{
	100: {EquipmentPartsCap},
	101: {EquipmentPartsFace},
	102: {EquipmentPartsEye},
	103: {EquipmentPartsEar},
	104: {EquipmentPartsTop},
	105: {EquipmentPartsTop},
	106: {EquipmentPartsPants},
	107: {EquipmentPartsShoes},
	108: {EquipmentPartsGlove},
	109: {EquipmentPartsShield},
	110: {EquipmentPartsCape},
	111: {EquipmentPartsRing, EquipmentPartsRing2, EquipmentPartsRing3, EquipmentPartsRing4},
	112: {EquipmentPartsPendant},
	113: {EquipmentPartsBelt},
	114: {EquipmentPartsMedal},
	119: {EquipmentPartsShield},
	180: {EquipmentPartsPetEquip},
	190: {EquipmentPartsTamingMob},
	191: {EquipmentPartsSaddle},
	192: {EquipmentPartsMobEquip},
}

func EquipmentParts(itemID uint32) []EquipmentPartsType {
	category := ItemCategoryOf(itemID)
	if parts, ok := equipmentPartsByCategory[category]; ok {
		return parts
	}
	switch category / 10 {
	case 13, 14, 16, 17:
		return []EquipmentPartsType{EquipmentPartsWeapon}
	}
	return nil
}

func CanEquipAt(itemID uint32, parts EquipmentPartsType) bool {
	if parts <= 2*equipmentPartsCashOffset || parts >= 0 {
		return false
	}
	if parts <= equipmentPartsCashOffset {
		parts -= equipmentPartsCashOffset
	}
	if parts >= 0 {
		return false
	}

	return slices.Contains(EquipmentParts(itemID), parts)
}
