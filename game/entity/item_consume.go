package entity

import (
	"github.com/boyism80/fm/game/constant"
)

type Consume struct {
	*ItemCore
	OwnerName string
	Flags     uint16
}

func (item *Consume) GetInventoryType() constant.InventoryType {
	return constant.INVENTORY_TYPE_CONSUME
}
func (item *Consume) GetCount() uint16 { return item.Count }
func (item *Consume) Reduce(count uint16) uint16 {
	if count > item.Count {
		item.Count = 0
	} else {
		item.Count -= count
	}
	return item.Count
}
func (item *Consume) Increase(count uint16) uint16 {
	item.Count += count
	return item.Count
}
func (item *Consume) Clone(count uint16) Item {
	return &Consume{
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
