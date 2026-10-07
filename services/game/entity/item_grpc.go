package entity

import (
	"fmt"
	"slices"
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

func buildInventoryProto(item Item, ownerID uint32, slot int32, uniqueID *uint64, ownerName string, flag uint16, enhanceChance uint8, enhanceCount uint8, skillBonus uint16) *internal.Inventory {
	if item == nil {
		return nil
	}
	out := &internal.Inventory{
		OwnerId:          ownerID,
		ItemId:           item.GetModel().GetID(),
		Slot:             slot,
		Count:            uint32(item.GetCount()),
		ExpirationUnixMs: expirationUnixMs(item.GetExpiration()),
		EnhanceChance:    uint32(enhanceChance),
		EnhanceCount:     uint32(enhanceCount),
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

func (e *EquipmentCore) ToProto(ownerID uint32, slot int32) *internal.Inventory {
	pb := buildInventoryProto(e, ownerID, slot, e.UniqueId, e.OwnerName, e.Flag, e.EnhanceChance, e.EnhanceCount, e.SkillBonus)
	if bonus := e.BonusStats.ToProto(); bonus != nil {
		pb.EquipBonusStats = bonus
	}
	return pb
}

func (item *Consume) ToProto(ownerID uint32, slot int32) *internal.Inventory {
	return buildInventoryProto(item, ownerID, slot, nil, item.OwnerName, item.Flags, 0, 0, 0)
}

func (item *Installation) ToProto(ownerID uint32, slot int32) *internal.Inventory {
	return buildInventoryProto(item, ownerID, slot, nil, item.OwnerName, item.Flags, 0, 0, 0)
}

func (item *MiscItem) ToProto(ownerID uint32, slot int32) *internal.Inventory {
	return buildInventoryProto(item, ownerID, slot, nil, item.OwnerName, item.Flags, 0, 0, 0)
}

func (item *CashItem) ToProto(ownerID uint32, slot int32) *internal.Inventory {
	return buildInventoryProto(item, ownerID, slot, item.UniqueId, item.OwnerName, item.Flags, 0, 0, 0)
}

func (item *Pet) ToProto(ownerID uint32, slot int32) *internal.Inventory {
	pb := buildInventoryProto(item, ownerID, slot, item.UniqueId, "", 0, 0, 0, 0)
	pb.Pet = &internal.Pet{
		Name:        item.Name,
		Level:       uint32(item.Level),
		Closeness:   uint32(item.Closeness),
		Fullness:    uint32(item.Fullness),
		Speed:       uint32(item.Speed),
		Skills:      uint32(item.Skills),
		SecondsLeft: item.SecondsLeft,
		Exceptions:  slices.Clone(item.Exceptions),
	}
	return pb
}

func NewItemFromInternalProto(pb *internal.Inventory, gw GameWorld) (Item, error) {
	if pb == nil {
		return nil, fmt.Errorf("nil Inventory")
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
	enhanceChance := uint8(pb.GetEnhanceChance())
	enhanceCount := uint8(pb.GetEnhanceCount())
	skillBonus := uint16(pb.GetSkillBonus())

	switch m := model.(type) {
	case wz.Equipment:
		if !constant.IsEquipment(itemID) {
			return nil, fmt.Errorf("item %d: not equippable", itemID)
		}
		return &EquipmentCore{
			ItemCore: &ItemCore{
				Wz:         m,
				Count:      1,
				Expiration: expiration,
			},
			EnhanceChance: enhanceChance,
			EnhanceCount:  enhanceCount,
			OwnerName:     ownerName,
			Flag:          flag,
			SkillBonus:    skillBonus,
			UniqueId:      uniqueID,
			BonusStats:    EquipmentBonusStatsFromProto(pb.GetEquipBonusStats()),
		}, nil
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
		if uniqueID == nil {
			v := gw.NewUniqueID()
			uniqueID = &v
		}
		pet := pb.GetPet()
		if pet == nil {
			pet = &internal.Pet{
				Name:        gw.GetResources().GetItemName(itemID),
				Level:       1,
				Fullness:    constant.PetMaxFullness,
				Skills:      uint32(m.Skills),
				SecondsLeft: uint32(m.LimitedLife),
			}
		}
		return &Pet{
			ItemCore:    &ItemCore{Wz: m, Count: 1, Expiration: expiration},
			UniqueId:    uniqueID,
			Name:        pet.GetName(),
			Level:       uint8(pet.GetLevel()),
			Closeness:   uint16(pet.GetCloseness()),
			Fullness:    uint8(pet.GetFullness()),
			Speed:       uint16(pet.GetSpeed()),
			Skills:      constant.PetSkill(pet.GetSkills()),
			SecondsLeft: pet.GetSecondsLeft(),
			Exceptions:  slices.Clone(pet.GetExceptions()),
		}, nil
	default:
		return nil, fmt.Errorf("item %d: unsupported item type", itemID)
	}
}
