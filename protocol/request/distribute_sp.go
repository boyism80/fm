package request

import (
	"github.com/boyism80/fm/stream"
)

type DistributeSP struct {
	Tick    uint32
	SkillID uint32
}

func (*DistributeSP) Opcode() byte { return 0x49 }

func (p *DistributeSP) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (p *DistributeSP) Deserialize(reader *stream.StreamReader) {
	p.Tick = reader.ReadU32()
	p.SkillID = reader.ReadU32()
}
