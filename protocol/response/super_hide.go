package response

import (
	"github.com/boyism80/fm/stream"
)

type SuperHide struct {
	Hidden bool
}

func (p *SuperHide) Opcode() uint16 {
	return 0x62
}

func (p *SuperHide) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(12)
	v := byte(0)
	if p.Hidden {
		v = 1
	}
	writer.WriteU8(v)
	return nil
}

func (p *SuperHide) Deserialize(reader *stream.StreamReader) {
}
