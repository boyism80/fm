package constant

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

const EquipmentPartsCashOffset EquipmentPartsType = -100

func (p EquipmentPartsType) Cash() EquipmentPartsType {
	return p + EquipmentPartsCashOffset
}
