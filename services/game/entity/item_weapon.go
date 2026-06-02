package entity

import (
	"github.com/boyism80/fm/services/game/constant"
)

type Weapon struct {
	*EquipmentCore
}

func (item *Weapon) GetInventoryType() constant.InventoryType {
	return constant.InventoryTypeEquipment
}
func (item *Weapon) GetCount() uint16             { return 1 }
func (item *Weapon) Reduce(count uint16) uint16   { return 0 }
func (item *Weapon) Increase(count uint16) uint16 { return 0 }
func (item *Weapon) Clone(count uint16) Item {
	return &Weapon{EquipmentCore: cloneEquipmentCore(item.EquipmentCore, count)}
}

func (w *Weapon) WeaponType() constant.WeaponType {
	return constant.GetWeaponType(w.GetModel().GetID())
}
