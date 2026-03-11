package entity

import (
	"github.com/boyism80/fm/game/constant"
)

type Shoes struct {
	*EquipmentCore
}

func (item *Shoes) GetInventoryType() constant.InventoryType {
	return constant.INVENTORY_TYPE_EQUIPMENT
}
func (item *Shoes) GetCount() uint16             { return 1 }
func (item *Shoes) Reduce(count uint16) uint16   { return 0 }
func (item *Shoes) Increase(count uint16) uint16 { return 0 }
func (item *Shoes) Clone(count uint16) Item {
	return &Shoes{EquipmentCore: cloneEquipmentCore(item.EquipmentCore, count)}
}
