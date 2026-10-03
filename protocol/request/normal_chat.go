package request

import (
	"github.com/boyism80/fm/stream"
)

type NormalChat struct {
	Message           string
	DontRecordHistory bool
}

func (*NormalChat) Opcode() byte { return 0x20 }

func (m *NormalChat) Serialize(writer *stream.StreamWriter) error {
	writer.WriteStr16(m.Message)
	writer.WriteBoolean(m.DontRecordHistory)
	return nil
}

func (m *NormalChat) Deserialize(reader *stream.StreamReader) {
	m.Message = reader.ReadStr16()
	m.DontRecordHistory = reader.ReadBool()
}
