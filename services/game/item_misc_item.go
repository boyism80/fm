package game

import (
	"github.com/boyism80/fm/maple/constant"
)

type MiscItem struct {
	*itemCore
	ownerName  string
	marriageID uint32
}

func (item *MiscItem) GetInventoryType() constant.InventoryType {
	return constant.InventoryTypeETC
}
func (item *MiscItem) GetCount() uint16 { return item.count }
func (item *MiscItem) Reduce(count uint16) uint16 {
	if count > item.count {
		item.count = 0
	} else {
		item.count -= count
	}
	return item.count
}
func (item *MiscItem) Increase(count uint16) uint16 {
	item.count += count
	return item.count
}
func (item *MiscItem) Clone(count uint16) Item {
	return &MiscItem{
		itemCore: &ItemCore{
			fieldPlacement: nil,
			count:          count,
			wz:             item.wz,
			expiration:     item.expiration,
			flags:          item.flags,
		},
		ownerName:  item.ownerName,
		marriageID: item.marriageID,
	}
}
