package request

import (
	"github.com/boyism80/fm/stream"
)

type PetAutoPotion struct {
	Tick   uint32
	Slot   int16
	ItemID uint32
}

func (*PetAutoPotion) Opcode() byte { return 0x88 }

func (p *PetAutoPotion) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(0)
	writer.WriteU32(p.Tick)
	writer.Write16(p.Slot)
	writer.WriteU32(p.ItemID)
	return nil
}

func (p *PetAutoPotion) Deserialize(reader *stream.StreamReader) {
	reader.Skip(1)
	p.Tick = reader.ReadU32()
	p.Slot = reader.Read16()
	p.ItemID = reader.ReadU32()
}
