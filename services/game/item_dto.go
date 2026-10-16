package game

import (
	"github.com/boyism80/fm/maple/wz"
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/util"
)

func (item *Consume) BuildS2C() dto.Item {
	model := item.GetModel()
	cons, ok := model.(*wz.Consume)
	if !ok {
		panic("Consume.GetModel() must be *wz.Consume")
	}
	return &dto.ConsumeItem{
		ItemId:     cons.GetID(),
		Count:      item.GetCount(),
		Expiration: item.GetExpiration(),
		OwnerName:  item.ownerName,
		Flags:      item.flags,
	}
}

func (item *MiscItem) BuildS2C() dto.Item {
	return &dto.MiscItem{
		ItemId:     item.GetModel().GetID(),
		Count:      item.GetCount(),
		Expiration: item.GetExpiration(),
		OwnerName:  item.ownerName,
		Flags:      item.flags,
	}
}

func (item *CashItem) BuildS2C() dto.Item {
	return &dto.CashItem{
		ItemId:     item.GetModel().GetID(),
		UniqueId:   copyUint64Ptr(item.uniqueId),
		Count:      item.GetCount(),
		Expiration: item.GetExpiration(),
		OwnerName:  item.ownerName,
		Flags:      item.flags,
	}
}

func (item *Installation) BuildS2C() dto.Item {
	return &dto.InstallationItem{
		ItemId:     item.GetModel().GetID(),
		Expiration: item.GetExpiration(),
		OwnerName:  item.ownerName,
		Flags:      item.flags,
	}
}

func (pet *Pet) BuildS2C() dto.Item {
	return &dto.PetItem{
		ItemId:         pet.GetModel().GetID(),
		UniqueId:       copyUint64Ptr(pet.uniqueId),
		Expiration:     pet.GetExpiration(),
		PetName:        pet.name,
		PetLevel:       pet.level,
		PetCloseness:   pet.closeness,
		PetFullness:    pet.fullness,
		PetSpeed:       pet.speed,
		PetFlags:       uint16(pet.skills),
		PetExpiration:  pet.GetExpiration(),
		PetSecondsLeft: pet.secondsLeft,
	}
}

func (e *EquipmentCore) BuildS2C() dto.Item { return e.BuildEquipmentS2C() }

func (e *EquipmentCore) BuildEquipmentS2C() *dto.Equipment {
	if e == nil || e.wz == nil {
		return nil
	}
	model, ok := e.wz.(wz.Equipment)
	if !ok || model == nil {
		return nil
	}
	ability := model.GetAbility()
	b := e.bonusStats
	if b == nil {
		b = &EquipmentBonusStats{}
	}
	return &dto.Equipment{
		ItemId:        model.GetID(),
		UniqueId:      copyUint64Ptr(e.uniqueId),
		Expiration:    e.expiration,
		EnhanceChance: e.enhanceChance,
		EnhanceCount:  e.enhanceCount,
		Str:           uint16(util.ClampInt32(int32(ability.Str)+int32(b.Str), 0, 0xffff)),
		Dex:           uint16(util.ClampInt32(int32(ability.Dex)+int32(b.Dex), 0, 0xffff)),
		Int:           uint16(util.ClampInt32(int32(ability.Int)+int32(b.Int), 0, 0xffff)),
		Luk:           uint16(util.ClampInt32(int32(ability.Luk)+int32(b.Luk), 0, 0xffff)),
		MaxHP:         uint16(util.ClampInt32(int32(ability.MaxHP)+int32(b.MaxHP), 0, 0xffff)),
		MaxMP:         uint16(util.ClampInt32(int32(ability.MaxMP)+int32(b.MaxMP), 0, 0xffff)),
		PAD:           uint16(util.ClampInt32(int32(ability.PAD)+int32(b.PAD), 0, 0xffff)),
		MAD:           uint16(util.ClampInt32(int32(ability.MAD)+int32(b.MAD), 0, 0xffff)),
		PDD:           uint16(util.ClampInt32(int32(ability.PDD)+int32(b.PDD), 0, 0xffff)),
		MDD:           uint16(util.ClampInt32(int32(ability.MDD)+int32(b.MDD), 0, 0xffff)),
		ACC:           uint16(util.ClampInt32(int32(ability.ACC)+int32(b.ACC), 0, 0xffff)),
		Avoid:         uint16(util.ClampInt32(int32(ability.Avoid)+int32(b.Avoid), 0, 0xffff)),
		Hands:         uint16(util.ClampInt32(int32(ability.Hands)+int32(b.Hands), 0, 0xffff)),
		Speed:         uint16(util.ClampInt32(int32(ability.Speed)+int32(b.Speed), 0, 0xffff)),
		Jump:          uint16(util.ClampInt32(int32(ability.Jump)+int32(b.Jump), 0, 0xffff)),
		OwnerName:     e.ownerName,
		Flag:          e.flags,
		SkillBonus:    uint8(e.skillBonus),
	}
}
