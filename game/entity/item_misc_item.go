package entity

import (
	"github.com/boyism80/fm/game/constant"
)

type MiscItem struct {
	*ItemCore
	OwnerName string
	Flags     uint16
}

func (item *MiscItem) GetInventoryType() constant.InventoryType {
	return constant.INVENTORY_TYPE_ETC
}
func (item *MiscItem) GetCount() uint16 { return item.Count }
func (item *MiscItem) Reduce(count uint16) uint16 {
	if count > item.Count {
		item.Count = 0
	} else {
		item.Count -= count
	}
	return item.Count
}
func (item *MiscItem) Increase(count uint16) uint16 {
	item.Count += count
	return item.Count
}
func (item *MiscItem) Clone(count uint16) Item {
	return &MiscItem{
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
