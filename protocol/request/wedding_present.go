package request

import (
	pconst "github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/stream"
)

type WeddingPresent struct {
	Mode          pconst.WeddingPresentMode
	Slot          int16
	ItemID        uint32
	Count         uint16
	InventoryType uint8
	Index         uint8
}

func (*WeddingPresent) Opcode() byte { return 0x74 }

func (p *WeddingPresent) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(uint8(p.Mode))
	switch p.Mode {
	case pconst.WeddingPresentGive:
		writer.Write16(p.Slot)
		writer.WriteU32(p.ItemID)
		writer.WriteU16(p.Count)
	case pconst.WeddingPresentReceive:
		writer.WriteU8(p.InventoryType)
		writer.WriteU8(p.Index)
	}
	return nil
}

func (p *WeddingPresent) Deserialize(reader *stream.StreamReader) {
	p.Mode = pconst.WeddingPresentMode(reader.ReadU8())
	switch p.Mode {
	case pconst.WeddingPresentGive:
		p.Slot = reader.Read16()
		p.ItemID = reader.ReadU32()
		p.Count = reader.ReadU16()
	case pconst.WeddingPresentReceive:
		p.InventoryType = reader.ReadU8()
		p.Index = reader.ReadU8()
	}
}
