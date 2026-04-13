package entity

import (
	"github.com/boyism80/fm/services/game/constant"
)

type Glove struct {
	*EquipmentCore
}

func (item *Glove) GetInventoryType() constant.InventoryType {
	return constant.INVENTORY_TYPE_EQUIPMENT
}
func (item *Glove) GetCount() uint16             { return 1 }
func (item *Glove) Reduce(count uint16) uint16   { return 0 }
func (item *Glove) Increase(count uint16) uint16 { return 0 }
func (item *Glove) Clone(count uint16) Item {
	return &Glove{EquipmentCore: cloneEquipmentCore(item.EquipmentCore, count)}
}
