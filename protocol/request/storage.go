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
