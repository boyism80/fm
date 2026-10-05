package request

import (
	"github.com/boyism80/fm/stream"
)

type UseCatchItem struct {
	Tick   uint32
	Slot   uint16
	ItemID uint32
	MobOID uint32
}

func (*UseCatchItem) Opcode() byte { return 0x40 }

func (p *UseCatchItem) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU32(p.Tick)
	writer.WriteU16(p.Slot)
	writer.WriteU32(p.ItemID)
	writer.WriteU32(p.MobOID)
	return nil
}

func (p *UseCatchItem) Deserialize(reader *stream.StreamReader) {
	p.Tick = reader.ReadU32()
	p.Slot = reader.ReadU16()
	p.ItemID = reader.ReadU32()
	p.MobOID = reader.ReadU32()
}
