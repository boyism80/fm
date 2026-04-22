package response

import (
	"github.com/boyism80/fm/stream"
)

type ShowDiceEffect struct {
	CharacterID uint32
	EffectID    int32
	SkillID     uint32
	SkillLevel  uint8
}

func (p *ShowDiceEffect) Opcode() uint16 {
	return 0x8F
}

func (p *ShowDiceEffect) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU32(p.CharacterID)
	writer.WriteU8(0x03)
	writer.Write32(p.EffectID)
	writer.WriteU32(p.SkillID)
	writer.WriteU8(p.SkillLevel)
	writer.WriteU8(0)
	return nil
}

func (p *ShowDiceEffect) Deserialize(reader *stream.StreamReader) {
}
