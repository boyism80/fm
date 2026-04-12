package response

import (
	"github.com/boyism80/fm/stream"
)

type NormalChat struct {
	CharacterId       uint32
	Highlight         bool
	Message           string
	DontRecordHistory bool
}

func (a *NormalChat) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU32(a.CharacterId)
	writer.WriteBoolean(a.Highlight)
	writer.WriteStr16(a.Message)
	writer.WriteBoolean(a.DontRecordHistory)
	return nil
}

func (a *NormalChat) Deserialize(reader *stream.StreamReader) {
}

func (a *NormalChat) Opcode() uint16 {
	return 0x70
}
