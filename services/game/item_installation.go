package game

import (
	"github.com/boyism80/fm/maple/constant"
)

type Installation struct {
	*itemCore
	ownerName string
}

func (item *Installation) GetInventoryType() constant.InventoryType {
	return constant.InventoryTypeInstallation
}
func (item *Installation) GetCount() uint16             { return 1 }
func (item *Installation) Reduce(count uint16) uint16   { return 0 }
func (item *Installation) Increase(count uint16) uint16 { return 0 }
func (item *Installation) Clone(count uint16) Item {
	return &Installation{
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
