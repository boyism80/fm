package request

import (
	"github.com/boyism80/fm/stream"
)

// DistributeAP represents the DISTRIBUTE_AP packet (0x46)
// Used to distribute 1 AP to a single stat (STR, DEX, INT, LUK, HP, MP)
type DistributeAP struct {
	Tick     uint32
	StatType uint32 // 64=STR, 128=DEX, 256=INT, 512=LUK, 2048=HP, 8192=MP
}

func (p *DistributeAP) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (p *DistributeAP) Deserialize(reader *stream.StreamReader) {
	p.Tick = reader.ReadU32()
	p.StatType = reader.ReadU32()
}
