package resp

import (
	"github.com/boyism80/fm/common/stream"
)

type NormalChat struct {
	CharacterId       uint32
	Highlight         bool
	Message           string
	DontRecordHistory bool
}

func (a *NormalChat) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU16(0x70)
	writer.WriteU32(a.CharacterId)
	writer.WriteBoolean(a.Highlight)
	writer.WriteStr16(a.Message)
	writer.WriteBoolean(a.DontRecordHistory)
	return nil
}

func (a *NormalChat) Deserialize(reader *stream.StreamReader) error {
	return nil
}
