package resp

import (
	"github.com/boyism80/fm/common/stream"
)

type CheckName struct {
	Name   string
	Exists bool
}

func (a *CheckName) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU16(0x05)
	writer.WriteStr16(a.Name)
	writer.WriteBoolean(a.Exists)
	return nil
}

func (a *CheckName) Deserialize(reader *stream.StreamReader) error {
	return nil
}
