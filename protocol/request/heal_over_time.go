package request

import (
	"github.com/boyism80/fm/stream"
)

// HealOverTime represents the HEAL_OVER_TIME packet (0x48)
type HealOverTime struct {
	Tick   uint32
	HealHP uint16
	HealMP uint16
	PRate  uint8
}

func (p *HealOverTime) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (p *HealOverTime) Deserialize(reader *stream.StreamReader) {
	if reader.Remaining() >= 8 {
		reader.Skip(4)
	}

	p.HealHP = reader.ReadU16()
	p.HealMP = reader.ReadU16()
	p.PRate = reader.ReadU8()
}
