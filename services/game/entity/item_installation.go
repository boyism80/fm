package entity

import (
	"github.com/boyism80/fm/services/game/constant"
)

type Installation struct {
	*ItemCore
	OwnerName string
	Flags     uint16
}

func (item *Installation) GetInventoryType() constant.InventoryType {
	return constant.INVENTORY_TYPE_INSTALLATION
}
func (item *Installation) GetCount() uint16             { return 1 }
func (item *Installation) Reduce(count uint16) uint16   { return 0 }
func (item *Installation) Increase(count uint16) uint16 { return 0 }
func (item *Installation) Clone(count uint16) Item {
	return &Installation{
		ItemCore: &ItemCore{
			Drop:       nil,
			Count:      count,
			Wz:         item.Wz,
			Expiration: item.Expiration,
		},
		OwnerName: item.OwnerName,
		Flags:     item.Flags,
	}
}
