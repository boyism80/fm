package game

import (
	"fmt"
	"slices"
	"time"

	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/maple/wz"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
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

func (e *EquipmentCore) BuildRPC(ownerID uint32, slot int32) *internal.Inventory {
	pb := buildInventoryProto(e, ownerID, slot, e.uniqueId, e.ownerName, e.flags, e.enhanceChance, e.enhanceCount, e.skillBonus)
	if bonus := e.bonusStats.BuildRPC(); bonus != nil {
		pb.EquipBonusStats = bonus
	}
	return pb
}

func (item *Consume) BuildRPC(ownerID uint32, slot int32) *internal.Inventory {
	return buildInventoryProto(item, ownerID, slot, nil, item.ownerName, item.flags, 0, 0, 0)
}

func (item *Installation) BuildRPC(ownerID uint32, slot int32) *internal.Inventory {
	return buildInventoryProto(item, ownerID, slot, nil, item.ownerName, item.flags, 0, 0, 0)
}

func (item *MiscItem) BuildRPC(ownerID uint32, slot int32) *internal.Inventory {
	pb := buildInventoryProto(item, ownerID, slot, nil, item.ownerName, item.flags, 0, 0, 0)
	pb.MarriageId = item.marriageID
	return pb
}

func (item *CashItem) BuildRPC(ownerID uint32, slot int32) *internal.Inventory {
	return buildInventoryProto(item, ownerID, slot, item.uniqueId, item.ownerName, item.flags, 0, 0, 0)
}

func (item *Pet) BuildRPC(ownerID uint32, slot int32) *internal.Inventory {
	pb := buildInventoryProto(item, ownerID, slot, item.uniqueId, "", 0, 0, 0, 0)
	pb.Pet = &internal.Pet{
		Name:        item.name,
		Level:       uint32(item.level),
		Closeness:   uint32(item.closeness),
		Fullness:    uint32(item.fullness),
		Speed:       uint32(item.speed),
		Skills:      uint32(item.skills),
		SecondsLeft: item.secondsLeft,
		Exceptions:  slices.Clone(item.exceptions),
	}
	return pb
}

func NewItemFromInternalProto(pb *internal.Inventory, gw ItemWorld) (Item, error) {
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
		if !m.GetEquipmentType().Equippable() {
			return nil, fmt.Errorf("item %d: not equippable", itemID)
		}
		return &EquipmentCore{
			itemCore: &ItemCore{
				wz:         m,
				count:      1,
				expiration: expiration,
				flags:      flag,
			},
			enhanceChance: enhanceChance,
			enhanceCount:  enhanceCount,
			ownerName:     ownerName,
			skillBonus:    skillBonus,
			uniqueId:      uniqueID,
			bonusStats:    EquipmentBonusStatsFromProto(pb.GetEquipBonusStats()),
		}, nil
	case *wz.Consume:
		return &Consume{
			itemCore:  &ItemCore{wz: m, count: count, expiration: expiration, flags: flag},
			ownerName: ownerName,
		}, nil
	case *wz.Installation:
		return &Installation{
			itemCore:  &ItemCore{wz: m, count: count, expiration: expiration, flags: flag},
			ownerName: ownerName,
		}, nil
	case *wz.MiscItem:
		return &MiscItem{
			itemCore:   &ItemCore{wz: m, count: count, expiration: expiration, flags: flag},
			ownerName:  ownerName,
			marriageID: pb.GetMarriageId(),
		}, nil
	case *wz.CashItem:
		return &CashItem{
			itemCore:  &ItemCore{wz: m, count: count, expiration: expiration, flags: flag},
			uniqueId:  uniqueID,
			ownerName: ownerName,
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
			itemCore:    &ItemCore{wz: m, count: 1, expiration: expiration},
			uniqueId:    uniqueID,
			name:        pet.GetName(),
			level:       uint8(pet.GetLevel()),
			closeness:   uint16(pet.GetCloseness()),
			fullness:    uint8(pet.GetFullness()),
			speed:       uint16(pet.GetSpeed()),
			skills:      constant.PetSkill(pet.GetSkills()),
			secondsLeft: pet.GetSecondsLeft(),
			exceptions:  slices.Clone(pet.GetExceptions()),
		}, nil
	default:
		return nil, fmt.Errorf("item %d: unsupported item type", itemID)
	}
}
