package entity

import (
	"github.com/boyism80/fm/services/game/constant"
)

type MiscItem struct {
	*ItemCore
	OwnerName  string
	Flags      uint16
	MarriageID uint32
}

func (item *MiscItem) GetInventoryType() constant.InventoryType {
	return constant.InventoryTypeETC
}
func (item *MiscItem) GetCount() uint16 { return item.Count }
func (item *MiscItem) GetFlags() constant.ItemFlag {
	return constant.ItemFlag(item.Flags)
}
func (item *MiscItem) SetFlags(flags constant.ItemFlag) {
	item.Flags = uint16(flags)
}
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
			FieldPlacement: nil,
			Count:          count,
			Wz:             item.Wz,
			Expiration:     item.Expiration,
		},
		OwnerName:  item.OwnerName,
		Flags:      item.Flags,
		MarriageID: item.MarriageID,
	}
}
