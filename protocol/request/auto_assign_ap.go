package request

import (
	"github.com/boyism80/fm/stream"
)

type AutoAssignAP struct {
	Tick          uint32
	Unknown       uint32
	PrimaryStat   uint32
	Amount        uint32
	SecondaryStat uint32
	Amount2       uint32
}

func (p *AutoAssignAP) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (p *AutoAssignAP) Deserialize(reader *stream.StreamReader) {
	p.Tick = reader.ReadU32()
	p.Unknown = reader.ReadU32()

	if reader.Remaining() < 16 {
		return
	}

	p.PrimaryStat = reader.ReadU32()
	p.Amount = reader.ReadU32()
	p.SecondaryStat = reader.ReadU32()
	p.Amount2 = reader.ReadU32()
}
