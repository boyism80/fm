package response

import (
	"github.com/boyism80/fm/stream"
)

type Clock struct {
	Seconds int32
}

func (p *Clock) Opcode() uint16 {
	return 0x65
}

func (p *Clock) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(2)
	writer.Write32(p.Seconds)
	return nil
}

func (p *Clock) Deserialize(reader *stream.StreamReader) {
}
