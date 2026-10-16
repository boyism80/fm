package game

import (
	"github.com/boyism80/fm/maple/constant"
)

type CashItem struct {
	*itemCore
	uniqueId  *uint64
	ownerName string
}

func (item *CashItem) GetInventoryType() constant.InventoryType { return constant.InventoryTypeCash }
func (item *CashItem) GetCount() uint16                         { return item.count }
func (item *CashItem) Reduce(count uint16) uint16 {
	if count > item.count {
		item.count = 0
	} else {
		item.count -= count
	}
	return item.count
}
func (item *CashItem) Increase(count uint16) uint16 {
	item.count += count
	return item.count
}
func (item *CashItem) Clone(count uint16) Item {
	return &CashItem{
		itemCore: &ItemCore{
			fieldPlacement: nil,
			count:          count,
			wz:             item.wz,
			expiration:     item.expiration,
			flags:          item.flags,
		},
		uniqueId:  copyUint64Ptr(item.uniqueId),
		ownerName: item.ownerName,
	}
}
