package request

import (
	"github.com/boyism80/fm/stream"
)

type PetFood struct {
	Tick   uint32
	Slot   int16
	ItemID uint32
}

func (*PetFood) Opcode() byte { return 0x3B }

func (p *PetFood) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU32(p.Tick)
	writer.Write16(p.Slot)
	writer.WriteU32(p.ItemID)
	return nil
}

func (p *PetFood) Deserialize(reader *stream.StreamReader) {
	p.Tick = reader.ReadU32()
	p.Slot = reader.Read16()
	p.ItemID = reader.ReadU32()
}
