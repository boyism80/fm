package request

import (
	"github.com/boyism80/fm/stream"
)

type LieDetectorAnswer struct {
	Answer string
}

func (*LieDetectorAnswer) Opcode() byte { return 0x58 }

func (p *LieDetectorAnswer) Serialize(writer *stream.StreamWriter) error {
	writer.WriteStr16(p.Answer)
	return nil
}

func (p *LieDetectorAnswer) Deserialize(reader *stream.StreamReader) {
	p.Answer = reader.ReadStr16()
}
