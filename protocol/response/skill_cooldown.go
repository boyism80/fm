package response

import (
	"github.com/boyism80/fm/stream"
)

type SkillCooldown struct {
	SkillID      uint32
	RemainingSec uint32
}

func (p *SkillCooldown) Opcode() uint16 {
	return 0xA7
}

func (p *SkillCooldown) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU32(p.SkillID)
	writer.WriteU32(p.RemainingSec)
	return nil
}

func (p *SkillCooldown) Deserialize(reader *stream.StreamReader) {
}
