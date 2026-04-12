package request

import (
	"github.com/boyism80/fm/stream"
)

type DistributeAP struct {
	Tick     uint32
	StatType uint32
}

func (p *DistributeAP) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (p *DistributeAP) Deserialize(reader *stream.StreamReader) {
	p.Tick = reader.ReadU32()
	p.StatType = reader.ReadU32()
}
