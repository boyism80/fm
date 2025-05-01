package resp

import (
	"github.com/boyism80/fm/common/stream"
)

type DeleteCharacter struct {
	Id      uint32
	Success bool
}

func (a *DeleteCharacter) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU16(0x07)
	writer.WriteU32(a.Id)
	writer.WriteBoolean(!a.Success)

	return nil
}

func (a *DeleteCharacter) Deserialize(reader *stream.StreamReader) error {
	return nil
}
