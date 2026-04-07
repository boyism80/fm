package request

import (
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/stream"
)

type SortInventory struct {
	Tick          uint32
	InventoryType constant.InventoryType
}

func (p *SortInventory) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (p *SortInventory) Deserialize(reader *stream.StreamReader) {
	p.Tick = reader.ReadU32()
	p.InventoryType = constant.InventoryType(reader.ReadU8())
}
