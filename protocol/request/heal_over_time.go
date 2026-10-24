package request

import (
	"github.com/boyism80/fm/stream"
)

type HealOverTime struct {
	Tick   uint32
	HealHP uint16
	HealMP uint16
	PRate  uint8
}

func (*HealOverTime) Opcode() byte { return 0x48 }

func (p *HealOverTime) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU32(p.Tick)
	writer.WriteU16(p.HealHP)
	writer.WriteU16(p.HealMP)
	writer.WriteU8(p.PRate)
	return nil
}

func (p *HealOverTime) Deserialize(reader *stream.StreamReader) {
	reader.Skip(4)
	p.HealHP = reader.ReadU16()
	p.HealMP = reader.ReadU16()
	p.PRate = reader.ReadU8()
}
