package item

import (
	"fmt"
	"time"

	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/maple/wz"
	"github.com/boyism80/fm/protocol/dto"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/util"
)

func NewInventory(itemID uint32, count uint16, ownerID uint32, slot int32, resources *wz.Resources, newUniqueID func() uint64) (*internal.Inventory, error) {
	model := resources.Items[itemID]
	if model == nil {
		return nil, fmt.Errorf("item model not found for ID: %d", itemID)
	}

	item := &internal.Inventory{
		OwnerId:       ownerID,
		ItemId:        itemID,
		Slot:          slot,
		Count:         uint32(count),
		InventoryType: uint32(model.GetInventoryType()),
	}
	switch m := model.(type) {
	case wz.Equipment:
		if !m.GetEquipmentType().Equippable() {
			return nil, fmt.Errorf("item %d: model is equipment but not equippable", itemID)
		}
		item.Count = 1
		item.EnhanceChance = uint32(m.GetEnhanceChance())
	case *wz.Consume, *wz.MiscItem, *wz.CashItem:
	case *wz.Installation:
		item.Count = 1
	case *wz.Pet:
		uniqueID := newUniqueID()
		item.Count = 1
		item.UniqueId = &uniqueID
		item.ExpirationUnixMs = time.Now().AddDate(0, 0, m.Life).UnixMilli()
		item.Pet = &internal.Pet{
			Name:        resources.GetItemName(itemID),
			Level:       1,
			Fullness:    constant.PetMaxFullness,
			Skills:      uint32(m.Skills),
			SecondsLeft: uint32(m.LimitedLife),
		}
	default:
		return nil, fmt.Errorf("not supported item type")
	}
	return item, nil
}

func BuildS2C(pb *internal.Inventory, resources *wz.Resources, newUniqueID func() uint64) (dto.Item, error) {
	if pb == nil {
		return nil, fmt.Errorf("nil Inventory")
	}
	model := resources.Items[pb.GetItemId()]
	if model == nil {
		return nil, fmt.Errorf("item model not found for ID: %d", pb.GetItemId())
	}

	expiration := util.TimeMax
	if ms := pb.GetExpirationUnixMs(); ms > 0 {
		expiration = time.UnixMilli(ms)
	}
	var uniqueID *uint64
	if pb.UniqueId != nil {
		value := *pb.UniqueId
		uniqueID = &value
	}
	ownerName := pb.GetOwnerName()
	flags := uint16(pb.GetFlag())
	count := uint16(pb.GetCount())

	switch m := model.(type) {
	case wz.Equipment:
		if !m.GetEquipmentType().Equippable() {
			return nil, fmt.Errorf("item %d: not equippable", pb.GetItemId())
		}
		ability := m.GetAbility()
		bonus := pb.GetEquipBonusStats()
		return &dto.Equipment{
			ItemId:        pb.GetItemId(),
			UniqueId:      uniqueID,
			Expiration:    expiration,
			EnhanceChance: uint8(pb.GetEnhanceChance()),
			EnhanceCount:  uint8(pb.GetEnhanceCount()),
			Str:           equipmentStat(ability.Str, bonus.GetStr()),
			Dex:           equipmentStat(ability.Dex, bonus.GetDex()),
			Int:           equipmentStat(ability.Int, bonus.GetIntStat()),
			Luk:           equipmentStat(ability.Luk, bonus.GetLuk()),
			MaxHP:         equipmentStat(ability.MaxHP, bonus.GetMaxHp()),
			MaxMP:         equipmentStat(ability.MaxMP, bonus.GetMaxMp()),
			PAD:           equipmentStat(ability.PAD, bonus.GetPad()),
			MAD:           equipmentStat(ability.MAD, bonus.GetMad()),
			PDD:           equipmentStat(ability.PDD, bonus.GetPdd()),
			MDD:           equipmentStat(ability.MDD, bonus.GetMdd()),
			ACC:           equipmentStat(ability.ACC, bonus.GetAcc()),
			Avoid:         equipmentStat(ability.Avoid, bonus.GetAvoid()),
			Hands:         equipmentStat(ability.Hands, bonus.GetHands()),
			Speed:         equipmentStat(ability.Speed, bonus.GetSpeed()),
			Jump:          equipmentStat(ability.Jump, bonus.GetJump()),
			OwnerName:     ownerName,
			Flag:          flags,
			SkillBonus:    uint8(pb.GetSkillBonus()),
		}, nil
	case *wz.Consume:
		return &dto.ConsumeItem{ItemId: pb.GetItemId(), Count: count, Expiration: expiration, OwnerName: ownerName, Flags: flags}, nil
	case *wz.Installation:
		return &dto.InstallationItem{ItemId: pb.GetItemId(), Expiration: expiration, OwnerName: ownerName, Flags: flags}, nil
	case *wz.MiscItem:
		return &dto.MiscItem{ItemId: pb.GetItemId(), Count: count, Expiration: expiration, OwnerName: ownerName, Flags: flags}, nil
	case *wz.CashItem:
		return &dto.CashItem{ItemId: pb.GetItemId(), UniqueId: uniqueID, Count: count, Expiration: expiration, OwnerName: ownerName, Flags: flags}, nil
	case *wz.Pet:
		if uniqueID == nil {
			value := newUniqueID()
			uniqueID = &value
		}
		pet := pb.GetPet()
		if pet == nil {
			pet = &internal.Pet{
				Name:        resources.GetItemName(pb.GetItemId()),
				Level:       1,
				Fullness:    constant.PetMaxFullness,
				Skills:      uint32(m.Skills),
				SecondsLeft: uint32(m.LimitedLife),
			}
		}
		return &dto.PetItem{
			ItemId:         pb.GetItemId(),
			UniqueId:       uniqueID,
			Expiration:     expiration,
			PetName:        pet.GetName(),
			PetLevel:       uint8(pet.GetLevel()),
			PetCloseness:   uint16(pet.GetCloseness()),
			PetFullness:    uint8(pet.GetFullness()),
			PetSpeed:       uint16(pet.GetSpeed()),
			PetFlags:       uint16(pet.GetSkills()),
			PetExpiration:  expiration,
			PetSecondsLeft: pet.GetSecondsLeft(),
		}, nil
	default:
		return nil, fmt.Errorf("item %d: unsupported item type", pb.GetItemId())
	}
}

func equipmentStat(base uint16, bonus int32) uint16 {
	return uint16(util.ClampInt32(int32(base)+int32(int16(util.ClampInt32(bonus, -32768, 32767))), 0, 0xffff))
}
