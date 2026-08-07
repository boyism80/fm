package response

import (
	"github.com/boyism80/fm/stream"
)

const (
	ShipStateDocked   uint16 = 1
	ShipStateLeaving  uint16 = 3
	ShipSpecialBalrog uint16 = 1034
)

type ShipState struct {
	State uint16
}

func (p *ShipState) Opcode() uint16 {
	return 0x67
}

func (p *ShipState) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU16(p.State)
	return nil
}

func (p *ShipState) Deserialize(reader *stream.StreamReader) {
	p.State = reader.ReadU16()
}
