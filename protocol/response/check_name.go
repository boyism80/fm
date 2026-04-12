package response

import (
	"github.com/boyism80/fm/stream"
)

type CheckName struct {
	Name   string
	Exists bool
}

func (a *CheckName) Opcode() uint16 {
	return 0x05
}

func (a *CheckName) Serialize(writer *stream.StreamWriter) error {
	writer.WriteStr16(a.Name)
	writer.WriteBoolean(a.Exists)
	return nil
}

func (a *CheckName) Deserialize(reader *stream.StreamReader) {
}
