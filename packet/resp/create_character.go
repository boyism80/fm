package resp

import (
	"github.com/boyism80/fm/dao"
	"github.com/boyism80/fm/stream"
)

type CreateCharacter struct {
	Character *dao.Character
	Success   bool
}

func (a *CreateCharacter) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU16(0x06)
	writer.WriteBoolean(!a.Success)
	a.Character.Serialize(writer)
	return nil
}

func (a *CreateCharacter) Deserialize(reader *stream.StreamReader) error {
	return nil
}
