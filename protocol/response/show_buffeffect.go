package response

import (
	"github.com/boyism80/fm/stream"
)

type ShowBuffeffect struct {
	CharacterID uint32
	EffectID    uint8
	SkillID     uint32
	SkillLevel  uint8
	Additional  *uint8
}

func (p *ShowBuffeffect) Opcode() uint16 {
	return 0x8F
}

func (p *ShowBuffeffect) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU32(p.CharacterID)
	writer.WriteU8(p.EffectID)
	writer.WriteU32(p.SkillID)
	writer.WriteU8(p.SkillLevel)
	if p.Additional != nil {
		writer.WriteU8(*p.Additional)
	}
	return nil
}

func (p *ShowBuffeffect) Deserialize(reader *stream.StreamReader) error {
	return nil
}
