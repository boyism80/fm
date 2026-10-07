package entity

import (
	"slices"
	"time"

	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/wz"
)

type Pet struct {
	*ItemCore
	UniqueId    *uint64
	Name        string
	Level       uint8
	Closeness   uint16
	Fullness    uint8
	Speed       uint16
	Skills      constant.PetSkill
	SecondsLeft uint32
	Exceptions  []uint32
}

func (item *Pet) GetInventoryType() constant.InventoryType { return constant.InventoryTypeCash }
func (item *Pet) GetCount() uint16                         { return 1 }
func (item *Pet) Reduce(count uint16) uint16               { return 0 }
func (item *Pet) Increase(count uint16) uint16             { return 0 }
func (item *Pet) Alive(now time.Time) bool {
	if now.After(item.GetExpiration()) {
		return false
	}
	return item.GetModel().(*wz.Pet).LimitedLife == 0 || item.SecondsLeft > 0
}

func (item *Pet) AddCloseness(n int) bool {
	item.Closeness = uint16(min(int(item.Closeness)+n, constant.PetMaxCloseness))
	level := item.Level
	for item.Level < constant.PetMaxLevel && item.Closeness >= constant.PetClosenessByLevel[item.Level] {
		item.Level++
	}
	return item.Level > level
}

func (item *Pet) Clone(count uint16) Item {
	return &Pet{
		ItemCore: &ItemCore{
			FieldPlacement: nil,
			Count:          count,
			Wz:             item.Wz,
			Expiration:     item.ItemCore.Expiration,
		},
		UniqueId:    copyUint64Ptr(item.UniqueId),
		Name:        item.Name,
		Level:       item.Level,
		Closeness:   item.Closeness,
		Fullness:    item.Fullness,
		Speed:       item.Speed,
		Skills:      item.Skills,
		SecondsLeft: item.SecondsLeft,
		Exceptions:  slices.Clone(item.Exceptions),
	}
}
