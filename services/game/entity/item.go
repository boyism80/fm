package entity

import (
	"fmt"
	"time"

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
	GetCount() uint16
	GetCount32() int32
	SetCount(count uint16)
	Increase(count uint16) uint16
	Reduce(count uint16) uint16
	Clone(count uint16) Item
	BindFieldPlacement(placement *FieldPlacement)
	ToDTO() dto.Item
	ToProto(ownerID uint32, slot int32) *internal.InventoryPersisted
}

type FieldPlacement struct {
	*ObjectCore
	Owner        uint32
	SpawnedPoint types.Point[int16]
	DropType     constant.DropType
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

type Equipment interface {
	Item
	GetEquipmentCore() *EquipmentCore
	ToEquipmentDTO() *dto.Equipment
}

func copyStringPtr(p *string) *string {
	if p == nil {
		return nil
	}
	v := *p
	return &v
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
	viewer.Send(&response.SpawnItem{
		ID:           fp.OID,
		Animation:    constant.DropItemAnimationTypeNone,
		DropType:     fp.DropType,
		ItemModel:    item.GetModel(),
		Expiration:   item.GetExpiration(),
		Position:     fp.Position,
		OwnerID:      fp.Owner,
		SpawnedPoint: fp.SpawnedPoint,
		IsPlayerDrop: true,
	}, types.SEND_POLICY_ENCRYPT)
}

func (item *ItemCore) GetFieldPlacement() *FieldPlacement           { return item.FieldPlacement }
func (item *ItemCore) Getcount() uint16                             { return item.Count }
func (item *ItemCore) GetCount32() int32                            { return int32(item.Count) }
func (item *ItemCore) IsMeso() bool                                 { return false }
func (item *ItemCore) SetCount(count uint16)                        { item.Count = count }
func (item *ItemCore) GetModel() wz.Item                            { return item.Wz }
func (item *ItemCore) GetExpiration() time.Time                     { return item.Expiration }
func (item *ItemCore) BindFieldPlacement(placement *FieldPlacement) { item.FieldPlacement = placement }

func (fp *FieldPlacement) RegisterExpire(duration time.Duration) {
	fp.nextExpiry = time.Now().Add(duration)
}
func (fp *FieldPlacement) RegisterFFA(duration time.Duration) {
	fp.nextFFA = time.Now().Add(duration)
}
func (fp *FieldPlacement) ShouldExpire(now time.Time) bool {
	return !fp.nextExpiry.IsZero() && now.After(fp.nextExpiry)
}
func (fp *FieldPlacement) ShouldFFA(now time.Time) bool {
	return fp.DropType != constant.DropTypeFFA && !fp.nextFFA.IsZero() && now.After(fp.nextFFA)
}

func NewItem(itemId uint32, count uint16, gw GameWorld) (Item, error) {
	model, ok := gw.GetResources().Items[itemId]
	if !ok {
		return nil, fmt.Errorf("item model not found for ID: %d", itemId)
	}
	switch m := model.(type) {
	case wz.Equipment:
		if !constant.IsEquipment(itemId) {
			return nil, fmt.Errorf("item %d: model is equipment but GetEquipmentType is not equippable", itemId)
		}
		core := &EquipmentCore{
			ItemCore: &ItemCore{
				Wz:         m,
				Count:      1,
				Expiration: util.TimeMax,
			},
			EnhanceChance: m.GetEnhanceChance(),
		}
		switch constant.GetEquipmentType(itemId) {
		case constant.EquipmentTypeWeapon:
			return &Weapon{EquipmentCore: core}, nil
		case constant.EquipmentTypeShield:
			return &Shield{EquipmentCore: core}, nil
		case constant.EquipmentTypeCap:
			return &Cap{EquipmentCore: core}, nil
		case constant.EquipmentTypeCoat, constant.EquipmentTypeLongcoat:
			return &Top{EquipmentCore: core}, nil
		case constant.EquipmentTypePants:
			return &Pants{EquipmentCore: core}, nil
		case constant.EquipmentTypeShoes:
			return &Shoes{EquipmentCore: core}, nil
		case constant.EquipmentTypeGlove:
			return &Glove{EquipmentCore: core}, nil
		case constant.EquipmentTypeCape:
			return &Cape{EquipmentCore: core}, nil
		case constant.EquipmentTypeRing:
			return &RingEquip{EquipmentCore: core}, nil
		case constant.EquipmentTypeFace:
			return &Face{EquipmentCore: core}, nil
		case constant.EquipmentTypeAccessory:
			return &Accessory{EquipmentCore: core}, nil
		default:
			return nil, fmt.Errorf("item %d: unsupported equipment type", itemId)
		}
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
		petExpiration, _ := time.ParseInLocation("2006-01-02 15:04:05", "2025-05-30 09:30:00", util.KST)
		return &Pet{
			ItemCore:   &ItemCore{Wz: m, Count: 1},
			Expiration: petExpiration,
		}, nil
	default:
		return nil, fmt.Errorf("not supported item type")
	}
}
