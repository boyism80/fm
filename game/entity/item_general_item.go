package entity

import (
	"github.com/boyism80/fm/game/constant"
)

type GeneralItem struct {
	*ItemCore
	OwnerName string
	Flags     uint16
}

func (item *GeneralItem) GetInventoryType() constant.InventoryType {
	return constant.INVENTORY_TYPE_ETC
}
func (item *GeneralItem) GetCount() uint16 { return item.Count }
func (item *GeneralItem) Reduce(count uint16) uint16 {
	if count > item.Count {
		item.Count = 0
	} else {
		item.Count -= count
	}
	return item.Count
}
func (item *GeneralItem) Increase(count uint16) uint16 {
	item.Count += count
	return item.Count
}
func (item *GeneralItem) Clone(count uint16) Item {
	return &GeneralItem{
		ItemCore: &ItemCore{
			Drop:       nil,
			Count:      count,
			Wz:         item.Wz,
			UniqueId:   item.UniqueId,
			Expiration: item.Expiration,
		},
		OwnerName: item.OwnerName,
		Flags:     item.Flags,
	}
}
