package resp

import (
	"github.com/boyism80/fm/common/stream"
	"github.com/boyism80/fm/game/constant"
)

type EndSortInventory struct {
	InventoryType constant.InventoryType
}

func (p *EndSortInventory) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU16(0x29)
	writer.WriteU8(1)
	writer.WriteU8(uint8(p.InventoryType))
	return nil
}

func (p *EndSortInventory) Deserialize(reader *stream.StreamReader) error {
	return nil
}
