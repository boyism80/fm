package wz

import "github.com/boyism80/fm/game/constant"

type Weapon struct {
	*EquipmentCore
}

func (w *Weapon) WeaponType() constant.WeaponType {
	return constant.GetWeaponType(w.EquipmentCore.ItemCore.ID)
}
