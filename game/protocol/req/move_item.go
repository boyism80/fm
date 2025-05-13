package req

import (
	"github.com/boyism80/fm/common/stream"
	"github.com/boyism80/fm/game/constant"
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

func (p *MoveItem) Deserialize(reader *stream.StreamReader) error {
	p.Tick, _ = reader.ReadU32()
	inventoryType, _ := reader.ReadU8()
	p.InventoryType = constant.InventoryType(inventoryType)
	p.Source, _ = reader.Read16()
	p.Dest, _ = reader.Read16()
	p.Count, _ = reader.ReadU16()
	return nil
}
