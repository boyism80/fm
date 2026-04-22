package response

import (
	pconst "github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/stream"
)

type ShowSelfSkillEffect struct {
	Type       pconst.SkillEffectType
	SkillID    uint32
	SkillLevel uint8
	Additional *uint8
}

func (p *ShowSelfSkillEffect) Opcode() uint16 {
	return 0x97
}

func (p *ShowSelfSkillEffect) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(uint8(p.Type))
	writer.WriteU32(p.SkillID)
	writer.WriteU8(p.SkillLevel)
	if p.Additional != nil {
		writer.WriteU8(*p.Additional)
	}
	return nil
}

func (p *ShowSelfSkillEffect) Deserialize(reader *stream.StreamReader) {
}
