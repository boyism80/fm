package response

import (
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/stream"
)

type UpdateCharacterLook struct {
	Character *dto.Character
}

func (p *UpdateCharacterLook) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU32(p.Character.ID)
	writer.WriteU8(1)
	p.Character.SerializeLook(writer)

	writer.WriteU8(0)
	writer.WriteU8(0)
	writer.WriteU8(0)
	return nil
}

func (p *UpdateCharacterLook) Deserialize(reader *stream.StreamReader) {
}

func (p *UpdateCharacterLook) Opcode() uint16 {
	return 0x8E
}
