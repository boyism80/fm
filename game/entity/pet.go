package entity

import (
	"errors"
	"time"

	"github.com/boyism80/fm/common/stream"
	"github.com/boyism80/fm/common/util"
	"github.com/boyism80/fm/game/data"
)

type Pet struct {
	*BaseItem
	Level       uint8
	Closeness   uint16
	Fullness    uint8
	Speed       uint16
	Flags       uint16
	PetItemId   uint32
	SecondsLeft uint32
	Expiration  int64
}

func (pet *Pet) serializePet(writer *stream.StreamWriter) error {
	if pet == nil {
		return errors.New("serializePet: pet is nil")
	}

	template := pet.GetTemplate()

	writer.WriteU64(util.GetTime(pet.Expiration))
	writer.WriteStaticStr(template.GetName(), 13)
	writer.WriteU8(pet.Level)
	writer.WriteU16(pet.Closeness)
	writer.WriteU8(pet.Fullness)

	if pet == nil {
		writer.WriteU64(util.GetKoreanTimestamp(int64(float64(time.Now().UnixNano()) / float64(time.Millisecond) * 1.5)))
	} else {
		writer.WriteU64(util.GetTime(pet.Expiration))
	}

	writer.WriteU16(pet.Speed)
	writer.WriteU16(pet.Flags)

	if pet.PetItemId == 5000054 && pet.SecondsLeft > 0 {
		writer.WriteU32(pet.SecondsLeft)
	} else {
		writer.WriteU32(0)
	}

	return nil
}

func (pet *Pet) Serialize(writer *stream.StreamWriter, zeroPosition, leaveOut, trade, bagSlot bool) {
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
