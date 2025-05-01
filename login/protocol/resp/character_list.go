package resp

import (
	"github.com/boyism80/fm/common/stream"
	"github.com/boyism80/fm/game/entity"
)

type CharacterList struct {
	SecondPw   string
	Characters []entity.Character
	SlotCount  uint32
}

func (e *CharacterList) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU16(0x03)
	writer.WriteU8(0)
	writer.WriteU32(0)
	writer.WriteU8(uint8(len(e.Characters)))

	for _, ch := range e.Characters {
		ch.Serialize(writer)
	}

	if e.SecondPw != "" {
		writer.WriteU8(1)
	} else {
		writer.WriteU8(2)
	}
	writer.WriteU8(0)
	writer.WriteU32(e.SlotCount)

	return nil
}

func (e *CharacterList) Deserialize(reader *stream.StreamReader) error {
	return nil
}
