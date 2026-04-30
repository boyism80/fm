package entity

import (
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/services/game/wz"
	"github.com/boyism80/fm/util"
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
		Count:      item.GetCount(),
		Expiration: item.GetExpiration(),
		OwnerName:  item.OwnerName,
		Flags:      item.Flags,
	}
}

func (item *MiscItem) ToDTO() dto.Item {
	return &dto.MiscItem{
		ItemId:     item.GetModel().GetID(),
		Count:      item.GetCount(),
		Expiration: item.GetExpiration(),
		OwnerName:  item.OwnerName,
		Flags:      item.Flags,
	}
}

func (item *CashItem) ToDTO() dto.Item {
	return &dto.CashItem{
		ItemId:     item.GetModel().GetID(),
		UniqueId:   copyUint64Ptr(item.UniqueId),
		Count:      item.GetCount(),
		Expiration: item.GetExpiration(),
		OwnerName:  item.OwnerName,
		Flags:      item.Flags,
	}
}

func (item *Installation) ToDTO() dto.Item {
	return &dto.InstallationItem{
		ItemId:     item.GetModel().GetID(),
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
		UniqueId:       copyUint64Ptr(pet.UniqueId),
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

func ToEquipmentDTOFromCore(core *EquipmentCore) *dto.Equipment {
	if core == nil || core.Wz == nil {
		return nil
	}
	model, ok := core.Wz.(wz.Equipment)
	if !ok || model == nil {
		return nil
	}
	ability := model.GetAbility()
	b := core.BonusStats
	if b == nil {
		b = &EquipmentBonusStats{}
	}
	return &dto.Equipment{
		ItemId:        model.GetID(),
		UniqueId:      copyUint64Ptr(core.UniqueId),
		Expiration:    core.Expiration,
		EnhanceChance: core.EnhanceChance,
		EnhanceCount:  core.EnhanceCount,
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
		OwnerName:     core.OwnerName,
		Flag:          core.Flag,
		SkillBonus:    uint8(core.SkillBonus),
	}
}

func (e *Weapon) ToDTO() dto.Item {
	return ToEquipmentDTOFromCore(e.EquipmentCore)
}
func (e *Weapon) ToEquipmentDTO() *dto.Equipment {
	return ToEquipmentDTOFromCore(e.EquipmentCore)
}

func equipToDTO(e Equipment) dto.Item {
	return ToEquipmentDTOFromCore(e.GetEquipmentCore())
}

func equipToEquipmentDTO(e Equipment) *dto.Equipment {
	return ToEquipmentDTOFromCore(e.GetEquipmentCore())
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
