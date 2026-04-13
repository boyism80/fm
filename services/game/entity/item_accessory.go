package entity

import (
	"github.com/boyism80/fm/services/game/constant"
)

type Accessory struct {
	*EquipmentCore
}

func (item *Accessory) GetInventoryType() constant.InventoryType {
	return constant.INVENTORY_TYPE_EQUIPMENT
}
func (item *Accessory) GetCount() uint16             { return 1 }
func (item *Accessory) Reduce(count uint16) uint16   { return 0 }
func (item *Accessory) Increase(count uint16) uint16 { return 0 }
func (item *Accessory) Clone(count uint16) Item {
	return &Accessory{EquipmentCore: cloneEquipmentCore(item.EquipmentCore, count)}
}
