package game

import (
	"github.com/boyism80/fm/maple/constant"
)

type Consume struct {
	*ItemCore
	OwnerName string
}

func (item *Consume) GetInventoryType() constant.InventoryType {
	return constant.InventoryTypeConsume
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
			FieldPlacement: nil,
			Count:          count,
			Wz:             item.Wz,
			Expiration:     item.Expiration,
			Flags:          item.Flags,
		},
		OwnerName: item.OwnerName,
	}
}
