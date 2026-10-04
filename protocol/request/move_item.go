package request

import (
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/stream"
)

type MoveItem struct {
	Tick          uint32
	InventoryType constant.InventoryType
	Source        int16
	Dest          int16
	Count         uint16
}

func (*MoveItem) Opcode() byte { return 0x36 }

func (p *MoveItem) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU32(p.Tick)
	writer.WriteU8(uint8(p.InventoryType))
	writer.Write16(p.Source)
	writer.Write16(p.Dest)
	writer.WriteU16(p.Count)
	return nil
}

func (p *MoveItem) Deserialize(reader *stream.StreamReader) {
	p.Tick = reader.ReadU32()
	p.InventoryType = constant.InventoryType(reader.ReadU8())
	p.Source = reader.Read16()
	p.Dest = reader.Read16()
	p.Count = reader.ReadU16()
}
