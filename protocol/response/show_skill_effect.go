package response

import (
	pconst "github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/stream"
)

type ShowSkillEffect struct {
	CharacterID uint32
	Type        pconst.SkillEffectType
	SkillID     uint32
	SkillLevel  uint8
	Additional  *uint8
}

func (p *ShowSkillEffect) Opcode() uint16 {
	return 0x8F
}

func (p *ShowSkillEffect) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU32(p.CharacterID)
	writer.WriteU8(uint8(p.Type))
	writer.WriteU32(p.SkillID)
	writer.WriteU8(p.SkillLevel)
	if p.Additional != nil {
		writer.WriteU8(*p.Additional)
	}
	return nil
}

func (p *ShowSkillEffect) Deserialize(reader *stream.StreamReader) {
}
