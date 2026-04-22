package response

import (
	"github.com/boyism80/fm/stream"
)

type ShowSelfDiceEffect struct {
	EffectID   int32
	SkillID    uint32
	SkillLevel uint8
}

func (p *ShowSelfDiceEffect) Opcode() uint16 {
	return 0x97
}

func (p *ShowSelfDiceEffect) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(0x03)
	writer.Write32(p.EffectID)
	writer.WriteU32(p.SkillID)
	writer.WriteU8(p.SkillLevel)
	writer.WriteU8(0)
	return nil
}

func (p *ShowSelfDiceEffect) Deserialize(reader *stream.StreamReader) {
}
