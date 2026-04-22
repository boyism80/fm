package response

import (
	"github.com/boyism80/fm/stream"
)

type ShowSelfHPHealedEffect struct {
	Amount int32
}

func (p *ShowSelfHPHealedEffect) Opcode() uint16 {
	return 0x97
}

func (p *ShowSelfHPHealedEffect) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(0x0A)
	writer.Write32(p.Amount)
	return nil
}

func (p *ShowSelfHPHealedEffect) Deserialize(reader *stream.StreamReader) {
}
