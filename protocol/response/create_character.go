package response

import (
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/stream"
)

type CreateCharacter struct {
	Character *dto.Character
	Success   bool
}

func (a *CreateCharacter) Opcode() uint16 {
	return 0x06
}

func (a *CreateCharacter) Serialize(writer *stream.StreamWriter) error {
	writer.WriteBoolean(!a.Success)
	a.Character.SerializeOverview(writer)
	return nil
}

func (a *CreateCharacter) Deserialize(reader *stream.StreamReader) {
}
