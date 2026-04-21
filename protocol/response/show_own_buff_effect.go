package response

import (
	"github.com/boyism80/fm/stream"
)

type ShowOwnBuffEffect struct {
	EffectID   uint8
	SkillID    uint32
	SkillLevel uint8
	Additional *uint8
}

func (p *ShowOwnBuffEffect) Opcode() uint16 {
	return 0x97
}

func (p *ShowOwnBuffEffect) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(p.EffectID)
	writer.WriteU32(p.SkillID)
	writer.WriteU8(p.SkillLevel)
	if p.Additional != nil {
		writer.WriteU8(*p.Additional)
	}
	return nil
}

func (p *ShowOwnBuffEffect) Deserialize(reader *stream.StreamReader) {
}
