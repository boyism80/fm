package request

import "github.com/boyism80/fm/stream"

type UseSummonBag struct {
	Tick   uint32
	Slot   int16
	ItemID uint32
}

func (*UseSummonBag) Opcode() byte { return 0x3A }

func (p *UseSummonBag) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU32(p.Tick)
	writer.Write16(p.Slot)
	writer.WriteU32(p.ItemID)
	return nil
}

func (p *UseSummonBag) Deserialize(reader *stream.StreamReader) {
	p.Tick = reader.ReadU32()
	p.Slot = reader.Read16()
	p.ItemID = reader.ReadU32()
}
