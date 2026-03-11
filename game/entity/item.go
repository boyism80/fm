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
	nextFFA      time.Time
	nextExpiry   time.Time
}

type ItemCore struct {
	*Drop
	Wz         wz.Item
	UniqueId   int64
	Expiration time.Time
	Count      uint16
}

type EquipmentCore struct {
	*ItemCore
	EnchantChance uint8
	OwnerName     string
	Flag          uint16
	SkillBonus    uint16
}

func (e *EquipmentCore) GetEquipmentCore() *EquipmentCore { return e }

type Equipment interface {
	Item
	GetEquipmentCore() *EquipmentCore
	ToEquipmentDTO() *dto.Equipment
}

func cloneEquipmentCore(c *EquipmentCore, count uint16) *EquipmentCore {
	return &EquipmentCore{
		ItemCore: &ItemCore{
			Drop:       nil,
			Count:      count,
			Wz:         c.Wz,
			UniqueId:   c.UniqueId,
			Expiration: c.Expiration,
		},
		EnchantChance: c.EnchantChance,
		OwnerName:     c.OwnerName,
		Flag:          c.Flag,
		SkillBonus:    c.SkillBonus,
	}
}

func (item *ItemCore) GetDrop() *Drop           { return item.Drop }
func (item *ItemCore) Getcount() uint16         { return item.Count }
func (item *ItemCore) GetCount32() int32        { return int32(item.Count) }
func (item *ItemCore) IsMeso() bool             { return false }
func (item *ItemCore) SetCount(count uint16)    { item.Count = count }
func (item *ItemCore) GetObject() *Object       { return item.Object }
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
	case wz.EquipmentModel:
		if !constant.IsEquipment(itemId) {
			return nil, fmt.Errorf("item %d: model is equipment but GetEquipmentType is not equippable", itemId)
		}
		eq := m.GetEquipment()
		core := &EquipmentCore{
			ItemCore: &ItemCore{
				Wz:         m,
				Count:      1,
				Expiration: util.TimeMax,
			},
			EnchantChance: eq.TUC,
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
	case *wz.GeneralItem:
		return &GeneralItem{
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
