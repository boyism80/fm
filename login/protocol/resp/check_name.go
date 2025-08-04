package resp

import (
	"github.com/boyism80/fm/core/stream"
)

type CheckName struct {
	Name   string
	Exists bool
}

// Opcode returns the packet opcode for CheckName
func (a *CheckName) Opcode() uint16 {
	return 0x05
}

func (a *CheckName) Serialize(writer *stream.StreamWriter) error {
	writer.WriteStr16(a.Name)
	writer.WriteBoolean(a.Exists)
	return nil
}

func (a *CheckName) Deserialize(reader *stream.StreamReader) error {
	return nil
}
