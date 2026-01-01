package response

import (
	"github.com/boyism80/fm/stream"
)

// UpdateSkills represents the UPDATE_SKILLS packet (0x19)
// Sent when a skill level or master level changes
type UpdateSkills struct {
	SkillID     uint32
	Level       int32
	MasterLevel int32
}

func (p *UpdateSkills) Opcode() uint16 {
	return 0x19
}

func (p *UpdateSkills) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(1)                      // Flag
	writer.WriteU16(1)                     // SkillCount
	writer.WriteU32(p.SkillID)             // SkillID
	writer.WriteU32(uint32(p.Level))       // Level
	writer.WriteU32(uint32(p.MasterLevel)) // MasterLevel
	writer.WriteU8(1)                      // Unknown
	return nil
}

func (p *UpdateSkills) Deserialize(reader *stream.StreamReader) error {
	return nil
}
