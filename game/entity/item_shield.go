package entity

import (
	"github.com/boyism80/fm/game/constant"
)

type Shield struct {
	*EquipmentCore
}

func (item *Shield) GetInventoryType() constant.InventoryType {
	return constant.INVENTORY_TYPE_EQUIPMENT
}
func (item *Shield) GetCount() uint16             { return 1 }
func (item *Shield) Reduce(count uint16) uint16   { return 0 }
func (item *Shield) Increase(count uint16) uint16 { return 0 }
func (item *Shield) Clone(count uint16) Item {
	return &Shield{EquipmentCore: cloneEquipmentCore(item.EquipmentCore, count)}
}
