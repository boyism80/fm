package request

import (
	"github.com/boyism80/fm/stream"
)

type LieDetectorSkill struct {
	Target string
}

func (*LieDetectorSkill) Opcode() byte { return 0x57 }

func (p *LieDetectorSkill) Serialize(writer *stream.StreamWriter) error {
	writer.WriteStr16(p.Target)
	return nil
}

func (p *LieDetectorSkill) Deserialize(reader *stream.StreamReader) {
	p.Target = reader.ReadStr16()
}
