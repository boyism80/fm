package request

import (
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/stream"
)

type MoveItem struct {
	Tick          uint32
	InventoryType constant.InventoryType
	Source        int16
	Dest          int16
	Count         uint16
}

func (p *MoveItem) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (p *MoveItem) Deserialize(reader *stream.StreamReader) {
	p.Tick = reader.ReadU32()
	p.InventoryType = constant.InventoryType(reader.ReadU8())
	p.Source = reader.Read16()
	p.Dest = reader.Read16()
	p.Count = reader.ReadU16()
}
