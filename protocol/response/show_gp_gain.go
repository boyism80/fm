package response

import (
	"github.com/boyism80/fm/stream"
)

type ShowGPGain struct {
	Amount int32
}

func (p *ShowGPGain) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(6)
	writer.Write32(p.Amount)
	return nil
}

func (p *ShowGPGain) Deserialize(reader *stream.StreamReader) {
}

func (p *ShowGPGain) Opcode() uint16 {
	return 0x1C
}
