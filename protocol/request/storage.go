package request

import (
	"github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/stream"
)

type Storage struct {
	Mode          constant.StorageMode
	InventoryType uint8
	Index         uint8
	Slot          int16
	ItemID        uint32
	Count         uint16
	Meso          int32
}

func (*Storage) Opcode() byte { return 0x2D }

func (p *Storage) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(uint8(p.Mode))
	switch p.Mode {
	case constant.StorageTakeOut:
		writer.WriteU8(p.InventoryType)
		writer.WriteU8(p.Index)
	case constant.StorageStore:
		writer.Write16(p.Slot)
		writer.WriteU32(p.ItemID)
		writer.WriteU16(p.Count)
	case constant.StorageMeso:
		writer.Write32(p.Meso)
	}
	return nil
}

func (p *Storage) Deserialize(reader *stream.StreamReader) {
	p.Mode = constant.StorageMode(reader.ReadU8())
	switch p.Mode {
	case constant.StorageTakeOut:
		p.InventoryType = reader.ReadU8()
		p.Index = reader.ReadU8()
	case constant.StorageStore:
		p.Slot = reader.Read16()
		p.ItemID = reader.ReadU32()
		p.Count = reader.ReadU16()
	case constant.StorageMeso:
		p.Meso = reader.Read32()
	}
}
