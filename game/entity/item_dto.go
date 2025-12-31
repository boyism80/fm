package entity

import (
	"github.com/boyism80/fm/game/wz"
	"github.com/boyism80/fm/protocol/dto"
)

// ItemToDTO converts entity Item to dto Item interface
func ItemToDTO(item Item) dto.Item {
	if item == nil {
		return nil
	}
	return item.ToDTO()
}

// ToDTO converts Consume to dto.ConsumeItem
func (item *Consume) ToDTO() dto.Item {
	model := item.GetModel()
	itemId := model.GetID()
	isThrowingStart := itemId/10000 == 207
	isBullet := itemId/10000 == 233
	isWhat := itemId/10000 == 287
	return &dto.ConsumeItem{
		ItemId:          itemId,
		UniqueId:        item.UniqueId,
		Count:           item.GetCount(),
		Expiration:      item.GetExpiration(),
		OwnerName:       item.OwnerName,
		Flags:           item.Flags,
		IsThrowingStart: isThrowingStart,
		IsBullet:        isBullet,
		IsWhat:          isWhat,
	}
}

// ToDTO converts GeneralItem to dto.GeneralItem
func (item *GeneralItem) ToDTO() dto.Item {
	return &dto.GeneralItem{
		ItemId:     item.GetModel().GetID(),
		UniqueId:   item.UniqueId,
		Count:      item.GetCount(),
		Expiration: item.GetExpiration(),
		OwnerName:  item.OwnerName,
		Flags:      item.Flags,
	}
}

// ToDTO converts CashItem to dto.CashItem
func (item *CashItem) ToDTO() dto.Item {
	return &dto.CashItem{
		ItemId:     item.GetModel().GetID(),
		UniqueId:   item.UniqueId,
		Count:      item.GetCount(),
		Expiration: item.GetExpiration(),
		OwnerName:  item.OwnerName,
		Flags:      item.Flags,
	}
}

// ToDTO converts Installation to dto.InstallationItem
func (item *Installation) ToDTO() dto.Item {
	return &dto.InstallationItem{
		ItemId:     item.GetModel().GetID(),
		UniqueId:   item.UniqueId,
		Expiration: item.GetExpiration(),
		OwnerName:  item.OwnerName,
		Flags:      item.Flags,
	}
}

// ToDTO converts Pet to dto.PetItem
func (pet *Pet) ToDTO() dto.Item {
	model := pet.GetModel()
	petModel, ok := model.(*wz.Pet)
	petName := ""
	if ok {
		petName = petModel.Name
	}
	return &dto.PetItem{
		ItemId:         model.GetID(),
		UniqueId:       pet.UniqueId,
		Expiration:     pet.GetExpiration(), // ItemCore.Expiration
		PetName:        petName,
		PetLevel:       pet.Level,
		PetCloseness:   pet.Closeness,
		PetFullness:    pet.Fullness,
		PetSpeed:       pet.Speed,
		PetFlags:       pet.Flags,
		PetExpiration:  pet.Expiration, // Pet.Expiration
		PetSecondsLeft: pet.SecondsLeft,
	}
}

// ToDTO converts Equipment to dto.Equipment
func (equipment *Equipment) ToDTO() dto.Item {
	return equipment.ToEquipmentDTO()
}

// ToEquipmentDTO converts Equipment to *dto.Equipment
func (equipment *Equipment) ToEquipmentDTO() *dto.Equipment {
	model, ok := equipment.Wz.(*wz.Equipment)
	if !ok {
		return nil
	}
	return &dto.Equipment{
		ItemId:        model.ID,
		UniqueId:      equipment.UniqueId,
		Expiration:    equipment.Expiration,
		EnchantChance: equipment.EnchantChance,
		Level:         model.Required.Level,
		Str:           model.Ability.Str,
		Dex:           model.Ability.Dex,
		Int:           model.Ability.Int,
		Luk:           model.Ability.Luk,
		MaxHP:         model.Ability.MaxHP,
		MaxMP:         model.Ability.MaxMP,
		PAD:           model.Ability.PAD,
		MAD:           model.Ability.MAD,
		PDD:           model.Ability.PDD,
		MDD:           model.Ability.MDD,
		ACC:           model.Ability.ACC,
		Avoid:         model.Ability.Avoid,
		Hands:         model.Ability.Hands,
		Speed:         model.Ability.Speed,
		Jump:          model.Ability.Jump,
		OwnerName:     equipment.OwnerName,
		Flag:          equipment.Flag,
		SkillBonus:    uint8(equipment.SkillBonus),
	}
}
