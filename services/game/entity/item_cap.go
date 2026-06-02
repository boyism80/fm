package entity

import (
	"github.com/boyism80/fm/services/game/constant"
)

type Cap struct {
	*EquipmentCore
}

func (item *Cap) GetInventoryType() constant.InventoryType { return constant.InventoryTypeEquipment }
func (item *Cap) GetCount() uint16                         { return 1 }
func (item *Cap) Reduce(count uint16) uint16               { return 0 }
func (item *Cap) Increase(count uint16) uint16             { return 0 }
func (item *Cap) Clone(count uint16) Item {
	return &Cap{EquipmentCore: cloneEquipmentCore(item.EquipmentCore, count)}
}
