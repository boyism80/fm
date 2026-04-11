package entity

import (
	"github.com/boyism80/fm/game/constant"
)

type Pants struct {
	*EquipmentCore
}

func (item *Pants) GetInventoryType() constant.InventoryType {
	return constant.INVENTORY_TYPE_EQUIPMENT
}
func (item *Pants) GetCount() uint16             { return 1 }
func (item *Pants) Reduce(count uint16) uint16   { return 0 }
func (item *Pants) Increase(count uint16) uint16 { return 0 }
func (item *Pants) Clone(count uint16) Item {
	return &Pants{EquipmentCore: cloneEquipmentCore(item.EquipmentCore, count)}
}
