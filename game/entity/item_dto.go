package entity

import (
	"github.com/boyism80/fm/game/wz"
	"github.com/boyism80/fm/protocol/dto"
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

func ToEquipmentDTOFromCore(core *EquipmentCore, model wz.EquipmentModel) *dto.Equipment {
	if core == nil || model == nil {
		return nil
	}
	eq := model.GetEquipment()
	maxEnchantChance := model.GetEquipment().EnchantChance
	return &dto.Equipment{
		ItemId:        eq.ItemCore.ID,
		UniqueId:      core.UniqueId,
		Expiration:    core.Expiration,
		EnchantChance: maxEnchantChance,
		EnchantCount:  maxEnchantChance - eq.EnchantChance,
		Str:           eq.Ability.Str,
		Dex:           eq.Ability.Dex,
		Int:           eq.Ability.Int,
		Luk:           eq.Ability.Luk,
		MaxHP:         eq.Ability.MaxHP,
		MaxMP:         eq.Ability.MaxMP,
		PAD:           eq.Ability.PAD,
		MAD:           eq.Ability.MAD,
		PDD:           eq.Ability.PDD,
		MDD:           eq.Ability.MDD,
		ACC:           eq.Ability.ACC,
		Avoid:         eq.Ability.Avoid,
		Hands:         eq.Ability.Hands,
		Speed:         eq.Ability.Speed,
		Jump:          eq.Ability.Jump,
		OwnerName:     core.OwnerName,
		Flag:          core.Flag,
		SkillBonus:    uint8(core.SkillBonus),
	}
}

func (e *Weapon) ToDTO() dto.Item {
	model, ok := e.GetModel().(wz.EquipmentModel)
	if !ok {
		return nil
	}
	return ToEquipmentDTOFromCore(e.EquipmentCore, model)
}
func (e *Weapon) ToEquipmentDTO() *dto.Equipment {
	model, ok := e.GetModel().(wz.EquipmentModel)
	if !ok {
		return nil
	}
	return ToEquipmentDTOFromCore(e.EquipmentCore, model)
}

func equipToDTO(e Equipment) dto.Item {
	model, ok := e.GetModel().(wz.EquipmentModel)
	if !ok {
		return nil
	}
	return ToEquipmentDTOFromCore(e.GetEquipmentCore(), model)
}

func equipToEquipmentDTO(e Equipment) *dto.Equipment {
	model, ok := e.GetModel().(wz.EquipmentModel)
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
