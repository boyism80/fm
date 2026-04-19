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

type Dropable interface {
	GetDrop() *Drop
	BindDrop(drop *Drop)
	GetCount32() int32
	IsMeso() bool
}

type Item interface {
	luax.Luable
	GetDrop() *Drop
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
	ToGrpcDTO(ownerID uint32, slot int32) *internal.InventoryPersisted
}

type Drop struct {
	*ObjectCore
	Owner        uint32
	SpawnedPoint types.Point[int16]
	DropType     constant.DropType
	nextFFA      time.Time
	nextExpiry   time.Time
}

func (d *Drop) GetObjectType() constant.ObjectType {
	return constant.ObjectTypeItem
}

type ItemCore struct {
	*Drop
	Wz         wz.Item
	Expiration time.Time
	Count      uint16
}

type EquipmentCore struct {
	*ItemCore
	EnchantChance uint8
	OwnerName     string
	Flag          uint16
	SkillBonus    uint16
	UniqueId      *uint64
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
	return &EquipmentCore{
		ItemCore: &ItemCore{
			Drop:       nil,
			Count:      count,
			Wz:         c.Wz,
			Expiration: c.Expiration,
		},
		EnchantChance: c.EnchantChance,
		OwnerName:     c.OwnerName,
		Flag:          c.Flag,
		SkillBonus:    c.SkillBonus,
		UniqueId:      copyUint64Ptr(c.UniqueId),
	}
}

func (item *ItemCore) GetObjectType() constant.ObjectType {
	return constant.ObjectTypeItem
}

func (item *ItemCore) SendSpawnSyncToViewer(viewer *Character) {
	if item == nil || viewer == nil {
		return
	}
	drop := item.GetDrop()
	if drop == nil {
		return
	}
	viewer.Send(&response.SpawnItem{
		ID:           drop.OID,
		Animation:    constant.DROP_ITEM_ANIMATION_TYPE_NONE,
		DropType:     drop.DropType,
		ItemModel:    item.GetModel(),
		Expiration:   item.GetExpiration(),
		Position:     drop.Position,
		OwnerID:      drop.Owner,
		SpawnedPoint: drop.SpawnedPoint,
		IsPlayerDrop: true,
	}, types.SEND_POLICY_ENCRYPT)
}

func (item *ItemCore) GetDrop() *Drop           { return item.Drop }
func (item *ItemCore) Getcount() uint16         { return item.Count }
func (item *ItemCore) GetCount32() int32        { return int32(item.Count) }
func (item *ItemCore) IsMeso() bool             { return false }
func (item *ItemCore) SetCount(count uint16)    { item.Count = count }
func (item *ItemCore) GetModel() wz.Item        { return item.Wz }
func (item *ItemCore) GetExpiration() time.Time { return item.Expiration }
func (item *ItemCore) BindDrop(drop *Drop)      { item.Drop = drop }

func (drop *Drop) RegisterExpire(duration time.Duration) { drop.nextExpiry = time.Now().Add(duration) }
func (drop *Drop) RegisterFFA(duration time.Duration)    { drop.nextFFA = time.Now().Add(duration) }
func (drop *Drop) ShouldExpire(now time.Time) bool {
	return !drop.nextExpiry.IsZero() && now.After(drop.nextExpiry)
}
func (drop *Drop) ShouldFFA(now time.Time) bool {
	return drop.DropType != constant.DROP_TYPE_FFA && !drop.nextFFA.IsZero() && now.After(drop.nextFFA)
}

func NewItem(itemId uint32, count uint16, context GameContext) (Item, error) {
	model, ok := context.GetResources().Items[itemId]
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
			EnchantChance: m.GetEnchantChance(),
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
