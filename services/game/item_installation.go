package game

import (
	"github.com/boyism80/fm/maple/constant"
)

type Installation struct {
	*ItemCore
	OwnerName string
}

func (item *Installation) GetInventoryType() constant.InventoryType {
	return constant.InventoryTypeInstallation
}
func (item *Installation) GetCount() uint16             { return 1 }
func (item *Installation) Reduce(count uint16) uint16   { return 0 }
func (item *Installation) Increase(count uint16) uint16 { return 0 }
func (item *Installation) Clone(count uint16) Item {
	return &Installation{
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
