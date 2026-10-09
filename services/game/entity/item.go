package entity

import (
	"fmt"
	"time"

	"github.com/boyism80/fm/core/clock"
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/protocol/dto"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/wz"
	"github.com/boyism80/fm/types"
	"github.com/boyism80/fm/util"
)

type FieldPlaceable interface {
	GetFieldPlacement() *FieldPlacement
	BindFieldPlacement(placement *FieldPlacement)
	GetCount32() int32
	IsMeso() bool
}

type Item interface {
	luax.Luable
	GetFieldPlacement() *FieldPlacement
	GetModel() wz.Item
	GetInventoryType() constant.InventoryType
	GetExpiration() time.Time
	SetExpiration(expiration time.Time)
	GetCount() uint16
	GetCount32() int32
	SetCount(count uint16)
	GetFlags() constant.ItemFlag
	SetFlags(flags constant.ItemFlag)
	Increase(count uint16) uint16
	Reduce(count uint16) uint16
	Clone(count uint16) Item
	BindFieldPlacement(placement *FieldPlacement)
	ToDTO() dto.Item
	ToProto(ownerID uint32, slot int32) *internal.Inventory
}

type FieldPlacement struct {
	*ObjectCore
	Owner        uint32
	SpawnedPoint types.Point[int16]
	DropType     constant.DropType
	PlayerDrop   bool
	Quest        uint32
	nextFFA      time.Time
	nextExpiry   time.Time
}

func (fp *FieldPlacement) GetObjectType() constant.ObjectType {
	return constant.ObjectTypeItem
}

type ItemCore struct {
	*FieldPlacement
	Wz         wz.Item
	Expiration time.Time
	Count      uint16
}

type EquipmentCore struct {
	*ItemCore
	EnhanceChance uint8
	EnhanceCount  uint8
	OwnerName     string
	Flag          uint16
	SkillBonus    uint16
	UniqueId      *uint64
	BonusStats    *EquipmentBonusStats
}

func (e *EquipmentCore) GetEquipmentCore() *EquipmentCore { return e }
func (e *EquipmentCore) GetInventoryType() constant.InventoryType {
	return constant.InventoryTypeEquipment
}
func (e *EquipmentCore) GetCount() uint16            { return 1 }
func (e *EquipmentCore) GetFlags() constant.ItemFlag { return constant.ItemFlag(e.Flag) }
func (e *EquipmentCore) SetFlags(flags constant.ItemFlag) {
	e.Flag = uint16(flags)
}
func (e *EquipmentCore) Reduce(count uint16) uint16   { return 0 }
func (e *EquipmentCore) Increase(count uint16) uint16 { return 0 }
func (e *EquipmentCore) Clone(count uint16) Item      { return cloneEquipmentCore(e, count) }
func (e *EquipmentCore) IsOverall() bool              { return e.GetModel().GetID()/10000 == 105 }

type Equipment interface {
	Item
	GetEquipmentCore() *EquipmentCore
	ToEquipmentDTO() *dto.Equipment
}

func copyUint64Ptr(p *uint64) *uint64 {
	if p == nil {
		return nil
	}
	v := *p
	return &v
}

func cloneEquipmentCore(c *EquipmentCore, count uint16) *EquipmentCore {
	var bonus *EquipmentBonusStats
	if c.BonusStats != nil {
		v := *c.BonusStats
		bonus = &v
	}
	return &EquipmentCore{
		ItemCore: &ItemCore{
			FieldPlacement: nil,
			Count:          count,
			Wz:             c.Wz,
			Expiration:     c.Expiration,
		},
		EnhanceChance: c.EnhanceChance,
		EnhanceCount:  c.EnhanceCount,
		OwnerName:     c.OwnerName,
		Flag:          c.Flag,
		SkillBonus:    c.SkillBonus,
		UniqueId:      copyUint64Ptr(c.UniqueId),
		BonusStats:    bonus,
	}
}

func (item *ItemCore) GetObjectType() constant.ObjectType {
	return constant.ObjectTypeItem
}

func (item *ItemCore) SendSpawnSyncToViewer(viewer *Character) {
	if item == nil || viewer == nil {
		return
	}
	fp := item.GetFieldPlacement()
	if fp == nil {
		return
	}
	if fp.Quest > 0 && !viewer.NeedsQuestItem(fp.Quest, item.GetModel().GetID()) {
		return
	}
	viewer.Send(&response.SpawnItem{
		ID:           fp.OID,
		Animation:    constant.DropItemAnimationTypeNone,
		DropType:     fp.DropType,
		ItemModel:    item.GetModel(),
		Expiration:   item.GetExpiration(),
		Position:     fp.Position,
		OwnerID:      fp.Owner,
		SpawnedPoint: fp.SpawnedPoint,
		IsPlayerDrop: fp.PlayerDrop,
	}, types.SEND_POLICY_ENCRYPT)
}

func (item *ItemCore) SendDestroySyncToViewer(viewer *Character) {
	if item == nil || viewer == nil {
		return
	}
	fp := item.GetFieldPlacement()
	if fp == nil {
		return
	}
	viewer.Send(&response.RemoveItem{
		Mode:        constant.RemoveItemTypeNoAnimated,
		OID:         fp.OID,
		CharacterId: 0,
	}, types.SEND_POLICY_ENCRYPT)
}

func (item *ItemCore) GetFieldPlacement() *FieldPlacement           { return item.FieldPlacement }
func (item *ItemCore) Getcount() uint16                             { return item.Count }
func (item *ItemCore) GetCount32() int32                            { return int32(item.Count) }
func (item *ItemCore) IsMeso() bool                                 { return false }
func (item *ItemCore) SetCount(count uint16)                        { item.Count = count }
func (item *ItemCore) GetFlags() constant.ItemFlag                  { return 0 }
func (item *ItemCore) SetFlags(flags constant.ItemFlag)             {}
func (item *ItemCore) GetModel() wz.Item                            { return item.Wz }
func (item *ItemCore) GetExpiration() time.Time                     { return item.Expiration }
func (item *ItemCore) SetExpiration(expiration time.Time)           { item.Expiration = expiration }
func (item *ItemCore) BindFieldPlacement(placement *FieldPlacement) { item.FieldPlacement = placement }

func (fp *FieldPlacement) RegisterExpire(duration time.Duration) {
	fp.nextExpiry = clock.Now().Add(duration)
}
func (fp *FieldPlacement) RegisterFFA(duration time.Duration) {
	fp.nextFFA = clock.Now().Add(duration)
}
func (fp *FieldPlacement) ShouldExpire(now time.Time) bool {
	return !fp.nextExpiry.IsZero() && now.After(fp.nextExpiry)
}
func (fp *FieldPlacement) ShouldFFA(now time.Time) bool {
	return fp.DropType != constant.DropTypeFFA && !fp.nextFFA.IsZero() && now.After(fp.nextFFA)
}

func NewItem(itemId uint32, count uint16, gw ItemWorld) (Item, error) {
	model, ok := gw.GetResources().Items[itemId]
	if !ok {
		return nil, fmt.Errorf("item model not found for ID: %d", itemId)
	}
	switch m := model.(type) {
	case wz.Equipment:
		if !constant.IsEquipment(itemId) {
			return nil, fmt.Errorf("item %d: model is equipment but GetEquipmentType is not equippable", itemId)
		}
		return &EquipmentCore{
			ItemCore: &ItemCore{
				Wz:         m,
				Count:      1,
				Expiration: util.TimeMax,
			},
			EnhanceChance: m.GetEnhanceChance(),
		}, nil
	case *wz.Consume:
		return &Consume{
			ItemCore: &ItemCore{Wz: m, Count: count, Expiration: util.TimeMax},
		}, nil
	case *wz.Installation:
		return &Installation{
			ItemCore: &ItemCore{Wz: m, Count: 1, Expiration: util.TimeMax},
		}, nil
	case *wz.MiscItem:
		return &MiscItem{
			ItemCore: &ItemCore{Wz: m, Count: count, Expiration: util.TimeMax},
		}, nil
	case *wz.CashItem:
		return &CashItem{
			ItemCore: &ItemCore{Wz: m, Count: count, Expiration: util.TimeMax},
		}, nil
	case *wz.Pet:
		uniqueID := gw.NewUniqueID()
		return &Pet{
			ItemCore:    &ItemCore{Wz: m, Count: 1, Expiration: time.Now().AddDate(0, 0, m.Life)},
			UniqueId:    &uniqueID,
			Name:        gw.GetResources().GetItemName(itemId),
			Level:       1,
			Fullness:    constant.PetMaxFullness,
			Skills:      m.Skills,
			SecondsLeft: uint32(m.LimitedLife),
		}, nil
	default:
		return nil, fmt.Errorf("not supported item type")
	}
}
