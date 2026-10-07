package request

import (
	"github.com/boyism80/fm/stream"
)

type SummonPet struct {
	Tick uint32
	Slot int16
}

func (*SummonPet) Opcode() byte { return 0x51 }

func (p *SummonPet) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU32(p.Tick)
	writer.Write16(p.Slot)
	return nil
}

func (p *SummonPet) Deserialize(reader *stream.StreamReader) {
	p.Tick = reader.ReadU32()
	p.Slot = reader.Read16()
}
