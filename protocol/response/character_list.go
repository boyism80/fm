package response

import (
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/stream"
)

type CharacterList struct {
	SecondPw   string
	Characters []dto.Character
	SlotCount  uint32
}

func (e *CharacterList) Opcode() uint16 {
	return 0x03
}

func (e *CharacterList) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(0)
	writer.WriteU32(0)
	writer.WriteU8(uint8(len(e.Characters)))

	for _, ch := range e.Characters {
		ch.SerializeOverview(writer)
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

func (e *CharacterList) Deserialize(reader *stream.StreamReader) {
}
