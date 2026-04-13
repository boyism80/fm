package constant

const (
	ITEM_CATEGORY_SHURIKEN uint32 = 207
	ITEM_CATEGORY_BULLET   uint32 = 233
)

const (
	ITEM_SHURIKEN_BASE uint32 = 2070000
	ITEM_BULLET_BASE   uint32 = 2330000
)

var RechargeableShurikens = []uint32{
	2070000, 2070001, 2070002, 2070003, 2070004, 2070005,
	2070006, 2070007, 2070008, 2070009, 2070010, 2070011,
	2070012, 2070013,
}

var RechargeableBullets = []uint32{
	2330000, 2330001, 2330002, 2330003, 2330004, 2330005,
	2331000, 2332000,
}

func GetEquipmentType(itemID uint32) EquipmentType {
	switch {
	case itemID >= 1300000 && itemID < 1800000:
		return EquipmentTypeWeapon
	case itemID >= 1092000 && itemID < 1100000:
		return EquipmentTypeShield
	case itemID >= 1000000 && itemID < 1010000:
		return EquipmentTypeCap
	case itemID >= 1040000 && itemID < 1050000:
		return EquipmentTypeCoat
	case itemID >= 1050000 && itemID < 1060000:
		return EquipmentTypeLongcoat
	case itemID >= 1060000 && itemID < 1070000:
		return EquipmentTypePants
	case itemID >= 1070000 && itemID < 1080000:
		return EquipmentTypeShoes
	case itemID >= 1080000 && itemID < 1090000:
		return EquipmentTypeGlove
	case itemID >= 1102000 && itemID < 1104000:
		return EquipmentTypeCape
	case itemID >= 1112000 && itemID < 1120000:
		return EquipmentTypeRing
	case itemID >= 20000 && itemID < 22000:
		return EquipmentTypeFace
	case (itemID >= 1132000 && itemID < 1183000) || (itemID >= 1010000 && itemID < 1040000) || (itemID >= 1122000 && itemID < 1123000):
		return EquipmentTypeAccessory
	case itemID >= 30000 && itemID < 35000:
		return EquipmentTypeHair
	case itemID >= 1172000 && itemID < 1173000:
		return EquipmentTypeMonsterBook
	case itemID >= 1662000 && itemID < 1680000:
		return EquipmentTypeAndroid
	case itemID >= 1920000 && itemID < 2000000:
		return EquipmentTypeDragon
	case itemID >= 1802000 && itemID < 1820000:
		return EquipmentTypePetEquip
	case itemID >= 1900000 && itemID < 1920000:
		return EquipmentTypeTaming
	case itemID >= 1610000 && itemID < 1660000:
		return EquipmentTypeMechanic
	default:
		return EquipmentTypeUnknown
	}
}

func IsEquipment(itemID uint32) bool {
	t := GetEquipmentType(itemID)
	return t != EquipmentTypeUnknown && t != EquipmentTypeMonsterBook && t != EquipmentTypeAndroid &&
		t != EquipmentTypeDragon && t != EquipmentTypePetEquip && t != EquipmentTypeTaming && t != EquipmentTypeMechanic && t != EquipmentTypeHair
}

func GetEquipmentPartsType(itemID uint32) EquipmentPartsType {
	switch GetEquipmentType(itemID) {
	case EquipmentTypeWeapon:
		return EQUIPMENT_PARTS_WEAPON
	case EquipmentTypeShield:
		return EQUIPMENT_PARTS_SHIELD
	case EquipmentTypeCap:
		return EQUIPMENT_PARTS_CAP
	case EquipmentTypeCoat, EquipmentTypeLongcoat:
		return EQUIPMENT_PARTS_TOP
	case EquipmentTypePants:
		return EQUIPMENT_PARTS_PANTS
	case EquipmentTypeShoes:
		return EQUIPMENT_PARTS_SHOES
	case EquipmentTypeGlove:
		return EQUIPMENT_PARTS_GLOVE
	case EquipmentTypeCape:
		return EQUIPMENT_PARTS_CAPE
	case EquipmentTypeRing:
		return EQUIPMENT_PARTS_RING
	case EquipmentTypeFace:
		return EQUIPMENT_PARTS_FACE
	case EquipmentTypeAccessory:
		return EQUIPMENT_PARTS_EYE
	default:
		return 0
	}
}

func GetWeaponType(itemID uint32) WeaponType {
	if itemID < 1300000 || itemID >= 1800000 {
		return WeaponTypeNone
	}
	cat := (itemID / 10000) % 100
	switch cat {
	case 30:
		return WeaponTypeSword1H
	case 31:
		return WeaponTypeAxe1H
	case 32:
		return WeaponTypeBlunt1H
	case 33:
		return WeaponTypeDagger
	case 34:
		return WeaponTypeKatara
	case 35:
		return WeaponTypeMagicArrow
	case 37:
		return WeaponTypeWand
	case 38:
		return WeaponTypeStaff
	case 40:
		return WeaponTypeSword2H
	case 41:
		return WeaponTypeAxe2H
	case 42:
		return WeaponTypeBlunt2H
	case 43:
		return WeaponTypeSpear
	case 44:
		return WeaponTypePoleArm
	case 45:
		return WeaponTypeBow
	case 46:
		return WeaponTypeCrossbow
	case 47:
		return WeaponTypeClaw
	case 48:
		return WeaponTypeKnuckle
	case 49:
		return WeaponTypeGun
	case 52:
		return WeaponTypeDualBow
	case 53:
		return WeaponTypeCannon
	default:
		return WeaponTypeNone
	}
}

func GetConsumeType(itemID uint32) ConsumeType {
	cat := itemID / 10000
	switch cat {
	case 203:
		return ConsumeTypeReturnScroll
	case 204:
		return ConsumeTypeEnhanceScroll
	case 206:
		if itemID >= 2061000 && itemID < 2062000 {
			return ConsumeTypeArrowCrossBow
		}
		return ConsumeTypeArrowBow
	case 207:
		return ConsumeTypeShuriken
	case 210:
		return ConsumeTypeMonsterSummon
	case 212:
		return ConsumeTypePetFood
	case 224:
		return ConsumeTypeWeddingRing
	case 226:
		return ConsumeTypeRidingFood
	case 228:
		return ConsumeTypeSkillBookGet
	case 229:
		return ConsumeTypeSkillBookUnlock
	case 231:
		return ConsumeTypeMinervaOwl
	case 232:
		return ConsumeTypeTeleportRock
	case 233:
		return ConsumeTypeBullet
	case 238:
		return ConsumeTypeCard
	default:
		return ConsumeTypeNone
	}
}

func ConsumeNeedsAmmoInventoryID(itemID uint32) bool {
	ct := GetConsumeType(itemID)
	return ct == ConsumeTypeShuriken || ct == ConsumeTypeBullet
}
