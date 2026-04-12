package response

import (
	"github.com/boyism80/fm/stream"
)

type UpdateSkills struct {
	SkillID     uint32
	Level       int32
	MasterLevel int32
}

func (p *UpdateSkills) Opcode() uint16 {
	return 0x19
}

func (p *UpdateSkills) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(1)
	writer.WriteU16(1)
	writer.WriteU32(p.SkillID)
	writer.WriteU32(uint32(p.Level))
	writer.WriteU32(uint32(p.MasterLevel))
	writer.WriteU8(1)
	return nil
}

func (p *UpdateSkills) Deserialize(reader *stream.StreamReader) {
}
