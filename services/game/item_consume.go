package game

import (
	"github.com/boyism80/fm/maple/constant"
)

type Consume struct {
	*itemCore
	ownerName string
}

func (item *Consume) GetInventoryType() constant.InventoryType {
	return constant.InventoryTypeConsume
}
func (item *Consume) GetCount() uint16 { return item.count }
func (item *Consume) Reduce(count uint16) uint16 {
	if count > item.count {
		item.count = 0
	} else {
		item.count -= count
	}
	return item.count
}
func (item *Consume) Increase(count uint16) uint16 {
	item.count += count
	return item.count
}
func (item *Consume) Clone(count uint16) Item {
	return &Consume{
		itemCore: &ItemCore{
			fieldPlacement: nil,
			count:          count,
			wz:             item.wz,
			expiration:     item.expiration,
			flags:          item.flags,
		},
		ownerName: item.ownerName,
	}
}
