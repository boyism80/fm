package entity

import (
	"time"

	"github.com/boyism80/fm/services/game/constant"
)

type Pet struct {
	*ItemCore
	UniqueId    *uint64
	Level       uint8
	Closeness   uint16
	Fullness    uint8
	Speed       uint16
	Flags       uint16
	SecondsLeft uint32
	Expiration  time.Time
}

func (item *Pet) GetInventoryType() constant.InventoryType { return constant.InventoryTypeCash }
func (item *Pet) GetCount() uint16                         { return 1 }
func (item *Pet) Reduce(count uint16) uint16               { return 0 }
func (item *Pet) Increase(count uint16) uint16             { return 0 }
func (item *Pet) Clone(count uint16) Item {
	return &Pet{
		ItemCore: &ItemCore{
			FieldPlacement: nil,
			Count:          count,
			Wz:             item.Wz,
			Expiration:     item.ItemCore.Expiration,
		},
		UniqueId:    copyUint64Ptr(item.UniqueId),
		Flags:       item.Flags,
		Level:       item.Level,
		Closeness:   item.Closeness,
		Fullness:    item.Fullness,
		Speed:       item.Speed,
		SecondsLeft: item.SecondsLeft,
		Expiration:  item.Expiration,
	}
}
