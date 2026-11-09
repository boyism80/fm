package response

import (
	"github.com/boyism80/fm/stream"
)

type ShowPopulationGain struct {
	Amount int32
}

func (p *ShowPopulationGain) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(4)
	writer.Write32(p.Amount)
	return nil
}

func (p *ShowPopulationGain) Deserialize(reader *stream.StreamReader) {
}

func (p *ShowPopulationGain) Opcode() uint16 {
	return 0x1C
}
