package entity

import (
	"fmt"
	"time"

	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/game/wz"
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/types"
	"github.com/boyism80/fm/util"
)

type Dropable interface {
	GetDrop() *Drop
	BindDrop(drop *Drop)
	GetCount32() int32
	IsMeso() bool
}

type Item interface {
	GetDrop() *Drop
	GetObject() *Object
	GetModel() wz.Item
	GetInventoryType() constant.InventoryType
	GetExpiration() time.Time
	GetCount() uint16
	GetCount32() int32
	SetCount(count uint16)
	Increase(count uint16) uint16
	Reduce(count uint16) uint16
	Clone(count uint16) Item
	BindDrop(drop *Drop)
	ToDTO() dto.Item
}

type Drop struct {
	*Object
	Owner        uint32
	SpawnedPoint types.Point[int16]
	DropType     constant.DropType
	nextFFA      time.Time // Time when item becomes FFA
	nextExpiry   time.Time // Time when item expires
}

type ItemCore struct {
	*Drop
	Wz         wz.Item
	UniqueId   int64
	Expiration time.Time
	Count      uint16
}

type Meso struct {
	*Drop
	Count int32
}

type CashItem struct {
	*ItemCore
	OwnerName string
	Flags     uint16
}

type Installation struct {
	*ItemCore
	OwnerName string
	Flags     uint16
}

type GeneralItem struct {
	*ItemCore
	OwnerName string
	Flags     uint16
}

type Equipment struct {
	*ItemCore
	EnchantChance uint8
	OwnerName     string
	Flag          uint16
	SkillBonus    uint16
}

type Consume struct {
	*ItemCore
	OwnerName string
	Flags     uint16
}

type Pet struct {
	*ItemCore
	Level       uint8
	Closeness   uint16
	Fullness    uint8
	Speed       uint16
	Flags       uint16
	SecondsLeft uint32
	Expiration  time.Time
}

func (item *ItemCore) GetDrop() *Drop {
	return item.Drop
}

func (item *ItemCore) Getcount() uint16 {
	return item.Count
}

func (item *ItemCore) GetCount32() int32 {
	return int32(item.Count)
}

func (item *ItemCore) IsMeso() bool {
	return false
}

func (item *ItemCore) SetCount(count uint16) {
	item.Count = count
}

func (item *ItemCore) GetObject() *Object {
	return item.Object
}

func (item *ItemCore) GetModel() wz.Item {
	return item.Wz
}

func (item *ItemCore) GetExpiration() time.Time {
	return item.Expiration
}

func (item *ItemCore) BindDrop(drop *Drop) {
	item.Drop = drop
}

func (meso *Meso) GetCount32() int32 {
	return meso.Count
}

func (meso *Meso) GetDrop() *Drop {
	return meso.Drop
}

func (meso *Meso) BindDrop(drop *Drop) {
	meso.Drop = drop
}

func (meso *Meso) IsMeso() bool {
	return true
}

// RegisterExpire sets the expiry time for the drop
func (drop *Drop) RegisterExpire(duration time.Duration) {
	drop.nextExpiry = time.Now().Add(duration)
}

// RegisterFFA sets the FFA time for the drop
func (drop *Drop) RegisterFFA(duration time.Duration) {
	drop.nextFFA = time.Now().Add(duration)
}

// ShouldExpire checks if the drop should expire
func (drop *Drop) ShouldExpire(now time.Time) bool {
	if drop.nextExpiry.IsZero() {
		return false
	}

	return now.After(drop.nextExpiry)
}

// ShouldFFA checks if the drop should become FFA
func (drop *Drop) ShouldFFA(now time.Time) bool {
	if drop.DropType == constant.DROP_TYPE_FFA {
		return false
	}

	if drop.nextFFA.IsZero() {
		return false
	}

	return now.After(drop.nextFFA)
}

func (item *CashItem) GetInventoryType() constant.InventoryType {
	return constant.INVENTORY_TYPE_CASH
}

func (item *CashItem) GetCount() uint16 {
	return item.Count
}

func (item *CashItem) Reduce(count uint16) uint16 {
	if count > item.Count {
		item.Count = 0
	} else {
		item.Count -= count
	}
	return item.Count
}

func (item *CashItem) Increase(count uint16) uint16 {
	item.Count += count
	return item.Count
}

func (item *CashItem) Clone(count uint16) Item {
	return &CashItem{
		ItemCore: &ItemCore{
			Drop:       nil,
			Count:      count,
			Wz:         item.Wz,
			UniqueId:   item.UniqueId,
			Expiration: item.Expiration,
		},
		OwnerName: item.OwnerName,
		Flags:     item.Flags,
	}
}

func (item *Installation) GetInventoryType() constant.InventoryType {
	return constant.INVENTORY_TYPE_INSTALLATION
}

func (item *Installation) GetCount() uint16 {
	return 1
}

func (item *Installation) Reduce(count uint16) uint16 {
	return 0
}

func (item *Installation) Increase(count uint16) uint16 {
	return 0
}

func (item *Installation) Clone(count uint16) Item {
	return &Installation{
		ItemCore: &ItemCore{
			Drop:       nil,
			Count:      count,
			Wz:         item.Wz,
			UniqueId:   item.UniqueId,
			Expiration: item.Expiration,
		},
		OwnerName: item.OwnerName,
		Flags:     item.Flags,
	}
}

func (item *GeneralItem) GetInventoryType() constant.InventoryType {
	return constant.INVENTORY_TYPE_ETC
}

func (item *GeneralItem) GetCount() uint16 {
	return item.Count
}

func (item *GeneralItem) Reduce(count uint16) uint16 {
	if count > item.Count {
		item.Count = 0
	} else {
		item.Count -= count
	}
	return item.Count
}

func (item *GeneralItem) Increase(count uint16) uint16 {
	item.Count += count
	return item.Count
}

func (item *GeneralItem) Clone(count uint16) Item {
	return &GeneralItem{
		ItemCore: &ItemCore{
			Drop:       nil,
			Count:      count,
			Wz:         item.Wz,
			UniqueId:   item.UniqueId,
			Expiration: item.Expiration,
		},
		OwnerName: item.OwnerName,
		Flags:     item.Flags,
	}
}

func (item *Equipment) GetInventoryType() constant.InventoryType {
	return constant.INVENTORY_TYPE_EQUIPMENT
}

func (item *Equipment) GetCount() uint16 {
	return 1
}

func (item *Equipment) Reduce(count uint16) uint16 {
	return 0
}

func (item *Equipment) Increase(count uint16) uint16 {
	return 0
}

func (item *Equipment) Clone(count uint16) Item {
	return &Equipment{
		ItemCore: &ItemCore{
			Drop:       nil,
			Count:      count,
			Wz:         item.Wz,
			UniqueId:   item.UniqueId,
			Expiration: item.Expiration,
		},
		OwnerName:     item.OwnerName,
		EnchantChance: item.EnchantChance,
		Flag:          item.Flag,
		SkillBonus:    item.SkillBonus,
	}
}

func (item *Equipment) IsOverall() bool {
	model := item.Wz
	return (model.GetID() / 10000) == 105
}

func (item *Consume) GetInventoryType() constant.InventoryType {
	return constant.INVENTORY_TYPE_CONSUME
}

func (item *Consume) GetCount() uint16 {
	return item.Count
}

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
			Drop:       nil,
			Count:      count,
			Wz:         item.Wz,
			UniqueId:   item.UniqueId,
			Expiration: item.Expiration,
		},
		OwnerName: item.OwnerName,
		Flags:     item.Flags,
	}
}

func (item *Pet) GetInventoryType() constant.InventoryType {
	return constant.INVENTORY_TYPE_CASH
}

func (item *Pet) GetCount() uint16 {
	return 1
}

func (item *Pet) Reduce(count uint16) uint16 {
	return 0
}

func (item *Pet) Increase(count uint16) uint16 {
	return 0
}

func (item *Pet) Clone(count uint16) Item {
	return &Pet{
		ItemCore: &ItemCore{
			Drop:       nil,
			Count:      count,
			Wz:         item.Wz,
			UniqueId:   item.UniqueId,
			Expiration: item.ItemCore.Expiration,
		},
		Flags:       item.Flags,
		Level:       item.Level,
		Closeness:   item.Closeness,
		Fullness:    item.Fullness,
		Speed:       item.Speed,
		SecondsLeft: item.SecondsLeft,
		Expiration:  item.Expiration,
	}
}

func NewMeso(count int32, position types.Point[int16], ownerID uint32, dropType constant.DropType, sequence uint32, context GameContext, mapInstance *Map) *Meso {
	obj := &Object{
		OID:      sequence,
		Position: position,
		Context:  context,
		Map:      mapInstance,
	}
	meso := &Meso{
		Drop: &Drop{
			Object:       obj,
			SpawnedPoint: position,
			DropType:     dropType,
			Owner:        ownerID,
		},
		Count: count,
	}

	return meso
}

// NewItem creates an item from item ID using GameContext
func NewItem(itemId uint32, count uint16, context GameContext) (Item, error) {
	// Get item model from resources
	model, ok := context.GetResources().Items[itemId]
	if !ok {
		return nil, fmt.Errorf("item model not found for ID: %d", itemId)
	}

	switch wz := model.(type) {
	case *wz.Equipment:
		return &Equipment{
			ItemCore: &ItemCore{
				Wz:         wz,
				Count:      1,
				Expiration: util.TimeMax,
			},
			EnchantChance: wz.TUC,
		}, nil

	case *wz.Consume:
		return &Consume{
			ItemCore: &ItemCore{
				Wz:         wz,
				Count:      count,
				Expiration: util.TimeMax,
			},
		}, nil

	case *wz.Installation:
		return &Installation{
			ItemCore: &ItemCore{
				Wz:         wz,
				Count:      1,
				Expiration: util.TimeMax,
			},
		}, nil

	case *wz.GeneralItem:
		return &GeneralItem{
			ItemCore: &ItemCore{
				Wz:         wz,
				Count:      count,
				Expiration: util.TimeMax,
			},
		}, nil

	case *wz.CashItem:
		return &CashItem{
			ItemCore: &ItemCore{
				Wz:         wz,
				Count:      count,
				Expiration: util.TimeMax,
			},
		}, nil

	case *wz.Pet:
		petExpiration, err := time.ParseInLocation("2006-01-02 15:04:05", "2025-05-30 09:30:00", util.KST)
		if err != nil {
			fmt.Println(err)
		}
		return &Pet{
			ItemCore: &ItemCore{
				Wz:    wz,
				Count: 1,
			},
			Expiration: petExpiration,
		}, nil

	default:
		return nil, fmt.Errorf("not supported item type")
	}
}
