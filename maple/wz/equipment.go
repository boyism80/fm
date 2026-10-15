package wz

import "github.com/boyism80/fm/maple/constant"

type Equipment interface {
	Item
	GetEquipmentType() constant.EquipmentType
	GetParts() []constant.EquipmentPartsType
	CanEquipAt(parts constant.EquipmentPartsType) bool
	GetEnhanceChance() uint8
	GetRequired() RequiredStats
	GetAbility() AbilityStats
}
