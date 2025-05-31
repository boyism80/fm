package resp

import (
	"github.com/boyism80/fm/common/stream"
	"github.com/boyism80/fm/game/entity"
)

type CreateCharacter struct {
	Character *entity.Character
	Success   bool
}

func (a *CreateCharacter) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU16(0x06)
	writer.WriteBoolean(!a.Success)
	a.Character.SerializeOverview(writer)
	return nil
}

func (a *CreateCharacter) Deserialize(reader *stream.StreamReader) error {
	return nil
}
