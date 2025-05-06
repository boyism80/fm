package entity

import (
	"github.com/boyism80/fm/common/stream"
	"github.com/boyism80/fm/game/data"
)

type Equipment struct {
	*BaseItem
}

func (pet *Equipment) Serialize(writer *stream.StreamWriter, zeroPosition, leaveOut, trade, bagSlot bool) {
	if zeroPosition {
		writer.WriteU8(0)
	} else {
		parts := 0 // pet은 parts가 뭐지?
		if parts <= -1 {
			parts *= -1
			if parts > 100 && parts < 1000 {
				parts -= 100
			}
		}

		if bagSlot {
			writer.WriteU32(uint32((parts % 100) - 1))
		} else if !trade {
			writer.WriteU16(uint16(parts))
		} else {
			writer.WriteU8(uint8(parts))
		}
	}

	template, ok := pet.GetTemplate().(*data.PetTemplate)
	if !ok {
		return
	}

	writer.WriteU8(3)
	writer.WriteU32(template.Id)

	hasUID := pet.UniqueId > 0
	writer.WriteBoolean(hasUID)
	if hasUID {
		writer.Write64(pet.UniqueId)
	}
	pet.serializePet(writer)
}
