package request

import (
	"github.com/boyism80/fm/stream"
)

type UseMountFood struct {
	Tick   uint32
	Slot   int16
	ItemID uint32
}

func (*UseMountFood) Opcode() byte { return 0x3C }

func (p *UseMountFood) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU32(p.Tick)
	writer.Write16(p.Slot)
	writer.WriteU32(p.ItemID)
	return nil
}

func (p *UseMountFood) Deserialize(reader *stream.StreamReader) {
	p.Tick = reader.ReadU32()
	p.Slot = reader.Read16()
	p.ItemID = reader.ReadU32()
}
