package wz

import "github.com/boyism80/fm/maple/constant"

type Weapon struct {
	*EquipmentCore
	WeaponType constant.WeaponType
	TwoHanded  bool
}

func (model *Weapon) setWeaponType() {
	if model.ID < 1300000 || model.ID >= 1800000 {
		model.WeaponType = constant.WeaponTypeNone
		return
	}

	switch (model.ID / 10000) % 100 {
	case 30:
		model.WeaponType = constant.WeaponTypeSword1H
	case 31:
		model.WeaponType = constant.WeaponTypeAxe1H
	case 32:
		model.WeaponType = constant.WeaponTypeBlunt1H
	case 33:
		model.WeaponType = constant.WeaponTypeDagger
	case 34:
		model.WeaponType = constant.WeaponTypeKatara
	case 35:
		model.WeaponType = constant.WeaponTypeMagicArrow
	case 37:
		model.WeaponType = constant.WeaponTypeWand
	case 38:
		model.WeaponType = constant.WeaponTypeStaff
	case 40:
		model.WeaponType = constant.WeaponTypeSword2H
	case 41:
		model.WeaponType = constant.WeaponTypeAxe2H
	case 42:
		model.WeaponType = constant.WeaponTypeBlunt2H
	case 43:
		model.WeaponType = constant.WeaponTypeSpear
	case 44:
		model.WeaponType = constant.WeaponTypePoleArm
	case 45:
		model.WeaponType = constant.WeaponTypeBow
	case 46:
		model.WeaponType = constant.WeaponTypeCrossbow
	case 47:
		model.WeaponType = constant.WeaponTypeClaw
	case 48:
		model.WeaponType = constant.WeaponTypeKnuckle
	case 49:
		model.WeaponType = constant.WeaponTypeGun
	case 52:
		model.WeaponType = constant.WeaponTypeDualBow
	case 53:
		model.WeaponType = constant.WeaponTypeCannon
	default:
		model.WeaponType = constant.WeaponTypeNone
	}

	switch model.WeaponType {
	case constant.WeaponTypeSword2H, constant.WeaponTypeAxe2H, constant.WeaponTypeBlunt2H, constant.WeaponTypeSpear,
		constant.WeaponTypePoleArm, constant.WeaponTypeBow, constant.WeaponTypeCrossbow, constant.WeaponTypeClaw,
		constant.WeaponTypeKnuckle, constant.WeaponTypeGun, constant.WeaponTypeCannon:
		model.TwoHanded = true
	}
}
