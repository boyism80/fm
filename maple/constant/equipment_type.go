package constant

type EquipmentType uint8

const (
	EquipmentTypeUnknown EquipmentType = iota
	EquipmentTypeWeapon
	EquipmentTypeShield
	EquipmentTypeCap
	EquipmentTypeCoat
	EquipmentTypeLongcoat
	EquipmentTypePants
	EquipmentTypeShoes
	EquipmentTypeGlove
	EquipmentTypeCape
	EquipmentTypeRing
	EquipmentTypeFace
	EquipmentTypeAccessory
	EquipmentTypeHair
	EquipmentTypeMonsterBook
	EquipmentTypeAndroid
	EquipmentTypeDragon
	EquipmentTypePetEquip
	EquipmentTypeTaming
	EquipmentTypeMechanic
)

func (t EquipmentType) Equippable() bool {
	return t != EquipmentTypeUnknown && t != EquipmentTypeMonsterBook && t != EquipmentTypeAndroid &&
		t != EquipmentTypeDragon && t != EquipmentTypeMechanic && t != EquipmentTypeHair
}
