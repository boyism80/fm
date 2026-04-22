package response

import (
	"github.com/boyism80/fm/stream"
)

type ShowHPHealedEffect struct {
	CharacterID uint32
	Amount      int32
}

func (p *ShowHPHealedEffect) Opcode() uint16 {
	return 0x8F
}

func (p *ShowHPHealedEffect) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU32(p.CharacterID)
	writer.WriteU8(0x0A)
	writer.Write32(p.Amount)
	return nil
}

func (p *ShowHPHealedEffect) Deserialize(reader *stream.StreamReader) {
}
