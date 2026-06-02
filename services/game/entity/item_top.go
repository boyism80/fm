package entity

import (
	"github.com/boyism80/fm/services/game/constant"
)

type Top struct {
	*EquipmentCore
}

func (item *Top) GetInventoryType() constant.InventoryType { return constant.InventoryTypeEquipment }
func (item *Top) GetCount() uint16                         { return 1 }
func (item *Top) Reduce(count uint16) uint16               { return 0 }
func (item *Top) Increase(count uint16) uint16             { return 0 }
func (item *Top) Clone(count uint16) Item {
	return &Top{EquipmentCore: cloneEquipmentCore(item.EquipmentCore, count)}
}

func (item *Top) IsOverall() bool { return item.GetModel().GetID()/10000 == 105 }
