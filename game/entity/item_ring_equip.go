package entity

import (
	"github.com/boyism80/fm/game/constant"
)

type RingEquip struct {
	*EquipmentCore
}

func (item *RingEquip) GetInventoryType() constant.InventoryType {
	return constant.INVENTORY_TYPE_EQUIPMENT
}
func (item *RingEquip) GetCount() uint16             { return 1 }
func (item *RingEquip) Reduce(count uint16) uint16   { return 0 }
func (item *RingEquip) Increase(count uint16) uint16 { return 0 }
func (item *RingEquip) Clone(count uint16) Item {
	return &RingEquip{EquipmentCore: cloneEquipmentCore(item.EquipmentCore, count)}
}
