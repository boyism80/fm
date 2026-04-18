package entity

import (
	"github.com/boyism80/fm/services/game/constant"
)

type CashItem struct {
	*ItemCore
	UniqueId  *uint64
	OwnerName string
	Flags     uint16
}

func (item *CashItem) GetInventoryType() constant.InventoryType { return constant.INVENTORY_TYPE_CASH }
func (item *CashItem) GetCount() uint16                         { return item.Count }
func (item *CashItem) Reduce(count uint16) uint16 {
	if count > item.Count {
		item.Count = 0
	} else {
		item.Count -= count
	}
	return item.Count
}
func (item *CashItem) Increase(count uint16) uint16 {
	item.Count += count
	return item.Count
}
func (item *CashItem) Clone(count uint16) Item {
	return &CashItem{
		ItemCore: &ItemCore{
			Drop:       nil,
			Count:      count,
			Wz:         item.Wz,
			Expiration: item.Expiration,
		},
		UniqueId:  copyUint64Ptr(item.UniqueId),
		OwnerName: item.OwnerName,
		Flags:     item.Flags,
	}
}
