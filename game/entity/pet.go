package entity

import (
	"time"

	"github.com/boyism80/fm/common/stream"
	"github.com/boyism80/fm/game/data"
)

type Pet struct {
	*BaseItem
	Level       uint8
	Closeness   uint16
	Fullness    uint8
	Speed       uint16
	Flags       uint16
	SecondsLeft uint32
	Expiration  time.Time
}

func (pet *Pet) Serialize(writer *stream.StreamWriter, zeroPosition, leaveOut, trade bool, slot int16, itemType ItemType) {

	template, ok := pet.Template.(*data.PetTemplate)
	if !ok {
		return // TODO: return error
	}

	if zeroPosition {
		if !leaveOut {
			writer.WriteU8(0)
		}
	} else {
		writer.WriteU8(uint8(slot))
	}

	writer.WriteU8(3) // pet
	writer.WriteU32(template.Id)

	hasUID := pet.UniqueId > 0
	writer.WriteBoolean(hasUID)
	if hasUID {
		writer.Write64(pet.UniqueId)
	}

	writer.WriteDateTime(pet.BaseItem.Expiration)
	writer.WriteStaticStr(template.Name, 13)
	writer.WriteU8(pet.Level)
	writer.WriteU16(pet.Closeness)
	writer.WriteU8(pet.Fullness)
	writer.WriteDateTime(pet.Expiration)
	writer.WriteU16(pet.Speed)
	writer.WriteU16(pet.Flags)
	if template.Id == 5000054 && pet.SecondsLeft > 0 {
		writer.WriteU32(pet.SecondsLeft)
	} else {
		writer.WriteU32(0)
	}
}
