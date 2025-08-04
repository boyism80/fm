package req

import (
	"github.com/boyism80/fm/core/stream"
	"github.com/boyism80/fm/game/constant"
)

type SortInventory struct {
	Tick          uint32
	InventoryType constant.InventoryType
}

func (p *SortInventory) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (p *SortInventory) Deserialize(reader *stream.StreamReader) error {
	p.Tick, _ = reader.ReadU32()
	inventoryType, _ := reader.ReadU8()
	p.InventoryType = constant.InventoryType(inventoryType)
	return nil
}
