package resp

import (
	"github.com/boyism80/fm/common/stream"
	"github.com/boyism80/fm/game/entity"
)

type UpdateCharacterLook struct {
	Character *entity.Character
}

func (p *UpdateCharacterLook) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU32(p.Character.ID)
	writer.WriteU8(1)
	p.Character.SerializeLook(writer)

	writer.WriteU8(0) // size of left rings
	writer.WriteU8(0) // size of right rings
	writer.WriteU8(0) // size of mid rings
	return nil
}

func (p *UpdateCharacterLook) Deserialize(reader *stream.StreamReader) error {
	return nil
}

func (p *UpdateCharacterLook) Opcode() uint16 {
	return 0x8E
}
