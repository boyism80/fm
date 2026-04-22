package response

import (
	"github.com/boyism80/fm/stream"
)

type ShowDragonBloodEffect struct {
	CharacterID uint32
	SkillID     uint32
	SkillLevel  uint8
}

func (p *ShowDragonBloodEffect) Opcode() uint16 {
	return 0x8F
}

func (p *ShowDragonBloodEffect) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU32(p.CharacterID)
	writer.WriteU8(5)
	writer.WriteU32(p.SkillID)
	writer.WriteU8(p.SkillLevel)
	return nil
}

func (p *ShowDragonBloodEffect) Deserialize(reader *stream.StreamReader) {
}
