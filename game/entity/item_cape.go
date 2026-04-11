package entity

import (
	"github.com/boyism80/fm/game/constant"
)

type Cape struct {
	*EquipmentCore
}

func (item *Cape) GetInventoryType() constant.InventoryType { return constant.INVENTORY_TYPE_EQUIPMENT }
func (item *Cape) GetCount() uint16                         { return 1 }
func (item *Cape) Reduce(count uint16) uint16               { return 0 }
func (item *Cape) Increase(count uint16) uint16             { return 0 }
func (item *Cape) Clone(count uint16) Item {
	return &Cape{EquipmentCore: cloneEquipmentCore(item.EquipmentCore, count)}
}
