package req

import (
	"github.com/boyism80/fm/common/stream"
	"github.com/boyism80/fm/game/constant"
)

type SortItem struct {
	Tick          uint32
	InventoryType constant.InventoryType
}

func (p *SortItem) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (p *SortItem) Deserialize(reader *stream.StreamReader) error {
	p.Tick, _ = reader.ReadU32()
	inventoryType, _ := reader.ReadU8()
	p.InventoryType = constant.InventoryType(inventoryType)
	return nil
}
