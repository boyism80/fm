package resp

import (
	"github.com/boyism80/fm/common/stream"
)

type DeleteCharacter struct {
	ID      uint32
	Success bool
}

func (a *DeleteCharacter) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU16(0x07)
	writer.WriteU32(a.ID)
	writer.WriteBoolean(!a.Success)

	return nil
}

func (a *DeleteCharacter) Deserialize(reader *stream.StreamReader) error {
	return nil
}
