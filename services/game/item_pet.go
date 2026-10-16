package game

import (
	"slices"
	"time"

	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/maple/wz"
)

type Pet struct {
	*itemCore
	uniqueId    *uint64
	name        string
	level       uint8
	closeness   uint16
	fullness    uint8
	speed       uint16
	skills      constant.PetSkill
	secondsLeft uint32
	exceptions  []uint32
}

func (item *Pet) GetInventoryType() constant.InventoryType { return constant.InventoryTypeCash }
func (item *Pet) GetCount() uint16                         { return 1 }
func (item *Pet) Reduce(count uint16) uint16               { return 0 }
func (item *Pet) Increase(count uint16) uint16             { return 0 }
func (item *Pet) Alive(now time.Time) bool {
	if now.After(item.GetExpiration()) {
		return false
	}
	return item.GetModel().(*wz.Pet).LimitedLife == 0 || item.secondsLeft > 0
}

func (item *Pet) SetCloseness(n int) bool {
	level := item.level
	item.closeness = uint16(min(max(n, 0), constant.PetMaxCloseness))
	item.level = 1
	for item.level < constant.PetMaxLevel && item.closeness >= constant.PetClosenessByLevel[item.level] {
		item.level++
	}
	return item.level > level
}

func (item *Pet) Clone(count uint16) Item {
	return &Pet{
		itemCore: &ItemCore{
			fieldPlacement: nil,
			count:          count,
			wz:             item.wz,
			expiration:     item.itemCore.expiration,
		},
		uniqueId:    copyUint64Ptr(item.uniqueId),
		name:        item.name,
		level:       item.level,
		closeness:   item.closeness,
		fullness:    item.fullness,
		speed:       item.speed,
		skills:      item.skills,
		secondsLeft: item.secondsLeft,
		exceptions:  slices.Clone(item.exceptions),
	}
}
