package entity

import (
	"fmt"
	"time"

	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/wz"
	"github.com/boyism80/fm/util"
)

func expirationUnixMs(t time.Time) int64 {
	if t.IsZero() || t.Equal(util.TimeMax) {
		return 0
	}
	return t.UnixMilli()
}

func buildInventoryProto(item Item, ownerID uint32, slot int32, uniqueID *uint64, ownerName string, flag uint16, enchantChance uint8, skillBonus uint16) *internal.InventoryPersisted {
	if item == nil {
		return nil
	}
	out := &internal.InventoryPersisted{
		OwnerId:          ownerID,
		ItemId:           item.GetModel().GetID(),
		Slot:             slot,
		Count:            uint32(item.GetCount()),
		ExpirationUnixMs: expirationUnixMs(item.GetExpiration()),
		EnchantChance:    uint32(enchantChance),
		Flag:             uint32(flag),
		SkillBonus:       uint32(skillBonus),
		OwnerName:        ownerName,
		InventoryType:    uint32(item.GetInventoryType()),
	}
	if uniqueID != nil {
		v := *uniqueID
		out.UniqueId = &v
	}
	return out
}

func equipmentToInventoryProto(e Equipment, ownerID uint32, slot int32) *internal.InventoryPersisted {
	if e == nil {
		return nil
	}
	c := e.GetEquipmentCore()
	if c == nil {
		return nil
	}
	return buildInventoryProto(e, ownerID, slot, c.UniqueId, c.OwnerName, c.Flag, c.EnchantChance, c.SkillBonus)
}

func (item *Weapon) ToProto(ownerID uint32, slot int32) *internal.InventoryPersisted {
	return equipmentToInventoryProto(item, ownerID, slot)
}

func (item *Shield) ToProto(ownerID uint32, slot int32) *internal.InventoryPersisted {
	return equipmentToInventoryProto(item, ownerID, slot)
}

func (item *Cap) ToProto(ownerID uint32, slot int32) *internal.InventoryPersisted {
	return equipmentToInventoryProto(item, ownerID, slot)
}

func (item *Face) ToProto(ownerID uint32, slot int32) *internal.InventoryPersisted {
	return equipmentToInventoryProto(item, ownerID, slot)
}

func (item *Accessory) ToProto(ownerID uint32, slot int32) *internal.InventoryPersisted {
	return equipmentToInventoryProto(item, ownerID, slot)
}

func (item *Top) ToProto(ownerID uint32, slot int32) *internal.InventoryPersisted {
	return equipmentToInventoryProto(item, ownerID, slot)
}

func (item *Pants) ToProto(ownerID uint32, slot int32) *internal.InventoryPersisted {
	return equipmentToInventoryProto(item, ownerID, slot)
}

func (item *Shoes) ToProto(ownerID uint32, slot int32) *internal.InventoryPersisted {
	return equipmentToInventoryProto(item, ownerID, slot)
}

func (item *Glove) ToProto(ownerID uint32, slot int32) *internal.InventoryPersisted {
	return equipmentToInventoryProto(item, ownerID, slot)
}

func (item *Cape) ToProto(ownerID uint32, slot int32) *internal.InventoryPersisted {
	return equipmentToInventoryProto(item, ownerID, slot)
}

func (item *RingEquip) ToProto(ownerID uint32, slot int32) *internal.InventoryPersisted {
	return equipmentToInventoryProto(item, ownerID, slot)
}

func (item *Consume) ToProto(ownerID uint32, slot int32) *internal.InventoryPersisted {
	return buildInventoryProto(item, ownerID, slot, nil, item.OwnerName, item.Flags, 0, 0)
}

func (item *Installation) ToProto(ownerID uint32, slot int32) *internal.InventoryPersisted {
	return buildInventoryProto(item, ownerID, slot, nil, item.OwnerName, item.Flags, 0, 0)
}

func (item *MiscItem) ToProto(ownerID uint32, slot int32) *internal.InventoryPersisted {
	return buildInventoryProto(item, ownerID, slot, nil, item.OwnerName, item.Flags, 0, 0)
}

func (item *CashItem) ToProto(ownerID uint32, slot int32) *internal.InventoryPersisted {
	return buildInventoryProto(item, ownerID, slot, item.UniqueId, item.OwnerName, item.Flags, 0, 0)
}

func (item *Pet) ToProto(ownerID uint32, slot int32) *internal.InventoryPersisted {
	return buildInventoryProto(item, ownerID, slot, item.UniqueId, "", item.Flags, 0, 0)
}

func NewItemFromInternalProto(pb *internal.InventoryPersisted, gw GameWorld) (Item, error) {
	if pb == nil {
		return nil, fmt.Errorf("nil InventoryPersisted")
	}
	if gw == nil {
		return nil, fmt.Errorf("nil GameWorld")
	}
	itemID := pb.GetItemId()
	model, ok := gw.GetResources().Items[itemID]
	if !ok {
		return nil, fmt.Errorf("item model not found for ID: %d", itemID)
	}

	expiration := util.TimeMax
	if exp := pb.GetExpirationUnixMs(); exp > 0 {
		expiration = time.UnixMilli(exp)
	}

	count := uint16(pb.GetCount())
	var uniqueID *uint64
	if pb.UniqueId != nil {
		v := *pb.UniqueId
		uniqueID = &v
	}
	ownerName := pb.GetOwnerName()
	flag := uint16(pb.GetFlag())
	enchantChance := uint8(pb.GetEnchantChance())
	skillBonus := uint16(pb.GetSkillBonus())

	switch m := model.(type) {
	case wz.Equipment:
		if !constant.IsEquipment(itemID) {
			return nil, fmt.Errorf("item %d: not equippable", itemID)
		}
		core := &EquipmentCore{
			ItemCore: &ItemCore{
				Wz:         m,
				Count:      1,
				Expiration: expiration,
			},
			EnchantChance: enchantChance,
			OwnerName:     ownerName,
			Flag:          flag,
			SkillBonus:    skillBonus,
			UniqueId:      uniqueID,
		}
		switch constant.GetEquipmentType(itemID) {
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
			return nil, fmt.Errorf("item %d: unsupported equipment type", itemID)
		}
	case *wz.Consume:
		return &Consume{
			ItemCore:  &ItemCore{Wz: m, Count: count, Expiration: expiration},
			OwnerName: ownerName,
			Flags:     flag,
		}, nil
	case *wz.Installation:
		return &Installation{
			ItemCore:  &ItemCore{Wz: m, Count: count, Expiration: expiration},
			OwnerName: ownerName,
			Flags:     flag,
		}, nil
	case *wz.MiscItem:
		return &MiscItem{
			ItemCore:  &ItemCore{Wz: m, Count: count, Expiration: expiration},
			OwnerName: ownerName,
			Flags:     flag,
		}, nil
	case *wz.CashItem:
		return &CashItem{
			ItemCore:  &ItemCore{Wz: m, Count: count, Expiration: expiration},
			UniqueId:  uniqueID,
			OwnerName: ownerName,
			Flags:     flag,
		}, nil
	case *wz.Pet:
		petExpiration, _ := time.ParseInLocation("2006-01-02 15:04:05", "2025-05-30 09:30:00", util.KST)
		if exp := pb.GetExpirationUnixMs(); exp > 0 {
			petExpiration = time.UnixMilli(exp)
		}
		return &Pet{
			ItemCore:   &ItemCore{Wz: m, Count: 1, Expiration: expiration},
			UniqueId:   uniqueID,
			Expiration: petExpiration,
		}, nil
	default:
		return nil, fmt.Errorf("item %d: unsupported item type", itemID)
	}
}
