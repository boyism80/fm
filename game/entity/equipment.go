package entity

import (
	"github.com/boyism80/fm/common/stream"
	"github.com/boyism80/fm/common/util"
	"github.com/boyism80/fm/game/data"
)

type Equipment struct {
	*BaseItem
	EnchantChance uint8
	OwnerName     string
	Flag          uint16
	SkillBonus    uint16
}

func (equipment *Equipment) Serialize(writer *stream.StreamWriter, zeroPosition, leaveOut, trade bool, slot int16, itemType ItemType) {
	template, ok := equipment.Template.(*data.EquipmentTemplate)
	if !ok {
		return // TODO: return error
	}

	if zeroPosition {
		if !leaveOut {
			writer.WriteU8(0)
		}
	} else {
		if slot <= -1 {
			slot *= -1
			if slot > 100 && slot < 1000 {
				slot -= 100
			}
		}
		if !trade && itemType == ItemTypeEquipment { // equipment
			writer.WriteU16(uint16(slot))
		} else {
			writer.WriteU8(uint8(slot))
		}
	}

	writer.WriteU8(uint8(itemType))
	writer.WriteU32(template.Id)

	hasUID := equipment.UniqueId > 0
	writer.WriteBoolean(hasUID)
	if hasUID {
		writer.Write64(equipment.UniqueId)
	}

	writer.WriteU64(util.GetTime(equipment.Expiration))
	writer.WriteU8(equipment.EnchantChance)
	writer.WriteU8(template.Required.Level)
	writer.WriteU16(template.Ability.Str)
	writer.WriteU16(template.Ability.Dex)
	writer.WriteU16(template.Ability.Int)
	writer.WriteU16(template.Ability.Luk)
	writer.WriteU16(template.Ability.MaxHP)
	writer.WriteU16(template.Ability.MaxMP)
	writer.WriteU16(template.Ability.PAD)
	writer.WriteU16(template.Ability.MAD)
	writer.WriteU16(template.Ability.PDD)
	writer.WriteU16(template.Ability.MDD)
	writer.WriteU16(template.Ability.ACC)
	writer.WriteU16(template.Ability.Avoid)
	writer.WriteU16(template.Ability.Hands)
	writer.WriteU16(template.Ability.Speed)
	writer.WriteU16(template.Ability.Jump)
	writer.WriteStr16(equipment.OwnerName)
	writer.WriteU16(equipment.Flag)
	writer.WriteBoolean(equipment.SkillBonus > 0)
	writer.WriteU8(1)  // item level?
	writer.WriteU32(0) // item exp percent?
	if equipment.UniqueId <= 0 {
		inventoryId := 0 // TODO: tracking
		if inventoryId > 0 {
			writer.WriteU64(uint64(inventoryId))
		} else {
			writer.Write64(-1)
		}
	}
	writer.WriteU64(util.GetTime(-2))
	writer.Write32(-1)
}
