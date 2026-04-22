package response

import (
	"github.com/boyism80/fm/stream"
)

type ShowSelfDragonBloodEffect struct {
	SkillID    uint32
	SkillLevel uint8
}

func (p *ShowSelfDragonBloodEffect) Opcode() uint16 {
	return 0x97
}

func (p *ShowSelfDragonBloodEffect) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(5)
	writer.WriteU32(p.SkillID)
	writer.WriteU8(p.SkillLevel)
	return nil
}

func (p *ShowSelfDragonBloodEffect) Deserialize(reader *stream.StreamReader) {
}
