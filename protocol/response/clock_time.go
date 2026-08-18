package response

import (
	"github.com/boyism80/fm/stream"
)

type ClockTime struct {
	Hour   uint8
	Minute uint8
	Second uint8
}

func (p *ClockTime) Opcode() uint16 {
	return 0x65
}

func (p *ClockTime) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(1)
	writer.WriteU8(p.Hour)
	writer.WriteU8(p.Minute)
	writer.WriteU8(p.Second)
	return nil
}

func (p *ClockTime) Deserialize(reader *stream.StreamReader) {
}
