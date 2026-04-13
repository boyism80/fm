package entity

import (
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/services/game/wz"
)

func ItemToDTO(item Item) dto.Item {
	if item == nil {
		return nil
	}
	return item.ToDTO()
}

func (item *Consume) ToDTO() dto.Item {
	model := item.GetModel()
	cons, ok := model.(*wz.Consume)
	if !ok {
		panic("Consume.GetModel() must be *wz.Consume")
	}
	return &dto.ConsumeItem{
		ItemId:     cons.GetID(),
		UniqueId:   item.UniqueId,
		Count:      item.GetCount(),
		Expiration: item.GetExpiration(),
		OwnerName:  item.OwnerName,
		Flags:      item.Flags,
	}
}

func (item *MiscItem) ToDTO() dto.Item {
	return &dto.MiscItem{
		ItemId:     item.GetModel().GetID(),
		UniqueId:   item.UniqueId,
		Count:      item.GetCount(),
		Expiration: item.GetExpiration(),
		OwnerName:  item.OwnerName,
		Flags:      item.Flags,
	}
}

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

func (item *Installation) ToDTO() dto.Item {
	return &dto.InstallationItem{
		ItemId:     item.GetModel().GetID(),
		UniqueId:   item.UniqueId,
		Expiration: item.GetExpiration(),
		OwnerName:  item.OwnerName,
		Flags:      item.Flags,
	}
}

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
		Expiration:     pet.GetExpiration(),
		PetName:        petName,
		PetLevel:       pet.Level,
		PetCloseness:   pet.Closeness,
		PetFullness:    pet.Fullness,
		PetSpeed:       pet.Speed,
		PetFlags:       pet.Flags,
		PetExpiration:  pet.Expiration,
		PetSecondsLeft: pet.SecondsLeft,
	}
}

func ToEquipmentDTOFromCore(core *EquipmentCore, model wz.Equipment) *dto.Equipment {
	if core == nil || model == nil {
		return nil
	}
	ability := model.GetAbility()
	maxEnchantChance := model.GetEnchantChance()
	return &dto.Equipment{
		ItemId:        model.GetID(),
		UniqueId:      core.UniqueId,
		Expiration:    core.Expiration,
		EnchantChance: maxEnchantChance,
		EnchantCount:  maxEnchantChance - model.GetEnchantChance(),
		Str:           ability.Str,
		Dex:           ability.Dex,
		Int:           ability.Int,
		Luk:           ability.Luk,
		MaxHP:         ability.MaxHP,
		MaxMP:         ability.MaxMP,
		PAD:           ability.PAD,
		MAD:           ability.MAD,
		PDD:           ability.PDD,
		MDD:           ability.MDD,
		ACC:           ability.ACC,
		Avoid:         ability.Avoid,
		Hands:         ability.Hands,
		Speed:         ability.Speed,
		Jump:          ability.Jump,
		OwnerName:     core.OwnerName,
		Flag:          core.Flag,
		SkillBonus:    uint8(core.SkillBonus),
	}
}

func (e *Weapon) ToDTO() dto.Item {
	model, ok := e.GetModel().(wz.Equipment)
	if !ok {
		return nil
	}
	return ToEquipmentDTOFromCore(e.EquipmentCore, model)
}
func (e *Weapon) ToEquipmentDTO() *dto.Equipment {
	model, ok := e.GetModel().(wz.Equipment)
	if !ok {
		return nil
	}
	return ToEquipmentDTOFromCore(e.EquipmentCore, model)
}

func equipToDTO(e Equipment) dto.Item {
	model, ok := e.GetModel().(wz.Equipment)
	if !ok {
		return nil
	}
	return ToEquipmentDTOFromCore(e.GetEquipmentCore(), model)
}

func equipToEquipmentDTO(e Equipment) *dto.Equipment {
	model, ok := e.GetModel().(wz.Equipment)
	if !ok {
		return nil
	}
	return ToEquipmentDTOFromCore(e.GetEquipmentCore(), model)
}

func (e *Shield) ToDTO() dto.Item                   { return equipToDTO(e) }
func (e *Shield) ToEquipmentDTO() *dto.Equipment    { return equipToEquipmentDTO(e) }
func (e *Cap) ToDTO() dto.Item                      { return equipToDTO(e) }
func (e *Cap) ToEquipmentDTO() *dto.Equipment       { return equipToEquipmentDTO(e) }
func (e *Face) ToDTO() dto.Item                     { return equipToDTO(e) }
func (e *Face) ToEquipmentDTO() *dto.Equipment      { return equipToEquipmentDTO(e) }
func (e *Accessory) ToDTO() dto.Item                { return equipToDTO(e) }
func (e *Accessory) ToEquipmentDTO() *dto.Equipment { return equipToEquipmentDTO(e) }
func (e *Top) ToDTO() dto.Item                      { return equipToDTO(e) }
func (e *Top) ToEquipmentDTO() *dto.Equipment       { return equipToEquipmentDTO(e) }
func (e *Pants) ToDTO() dto.Item                    { return equipToDTO(e) }
func (e *Pants) ToEquipmentDTO() *dto.Equipment     { return equipToEquipmentDTO(e) }
func (e *Shoes) ToDTO() dto.Item                    { return equipToDTO(e) }
func (e *Shoes) ToEquipmentDTO() *dto.Equipment     { return equipToEquipmentDTO(e) }
func (e *Glove) ToDTO() dto.Item                    { return equipToDTO(e) }
func (e *Glove) ToEquipmentDTO() *dto.Equipment     { return equipToEquipmentDTO(e) }
func (e *Cape) ToDTO() dto.Item                     { return equipToDTO(e) }
func (e *Cape) ToEquipmentDTO() *dto.Equipment      { return equipToEquipmentDTO(e) }
func (e *RingEquip) ToDTO() dto.Item                { return equipToDTO(e) }
func (e *RingEquip) ToEquipmentDTO() *dto.Equipment { return equipToEquipmentDTO(e) }
