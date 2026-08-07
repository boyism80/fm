package response

import (
	"github.com/boyism80/fm/stream"
)

type ShipSpecialEffect struct {
	Effect uint16
}

func (p *ShipSpecialEffect) Opcode() uint16 {
	return 0x66
}

func (p *ShipSpecialEffect) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU16(p.Effect)
	return nil
}

func (p *ShipSpecialEffect) Deserialize(reader *stream.StreamReader) {
	p.Effect = reader.ReadU16()
}
