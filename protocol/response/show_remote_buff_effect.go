package response

import (
	"github.com/boyism80/fm/stream"
)

type ShowRemoteBuffEffect struct {
	CharacterID uint32
	EffectID    uint8
	SkillID     uint32
	SkillLevel  uint8
	Additional  *uint8
}

func (p *ShowRemoteBuffEffect) Opcode() uint16 {
	return 0x8F
}

func (p *ShowRemoteBuffEffect) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU32(p.CharacterID)
	writer.WriteU8(p.EffectID)
	writer.WriteU32(p.SkillID)
	writer.WriteU8(p.SkillLevel)
	if p.Additional != nil {
		writer.WriteU8(*p.Additional)
	}
	return nil
}

func (p *ShowRemoteBuffEffect) Deserialize(reader *stream.StreamReader) {
}
