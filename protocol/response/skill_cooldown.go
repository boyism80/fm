package response

import (
	"github.com/boyism80/fm/stream"
)

// SkillCooldown notifies the client of a skill cooldown. RemainingSec 0 clears the cooldown UI.
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
