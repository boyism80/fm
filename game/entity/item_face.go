package entity

import (
	"github.com/boyism80/fm/game/constant"
)

type Face struct {
	*EquipmentCore
}

func (item *Face) GetInventoryType() constant.InventoryType { return constant.INVENTORY_TYPE_EQUIPMENT }
func (item *Face) GetCount() uint16                         { return 1 }
func (item *Face) Reduce(count uint16) uint16               { return 0 }
func (item *Face) Increase(count uint16) uint16             { return 0 }
func (item *Face) Clone(count uint16) Item {
	return &Face{EquipmentCore: cloneEquipmentCore(item.EquipmentCore, count)}
}
