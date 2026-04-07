package response

import (
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/stream"
)

type EndSortInventory struct {
	InventoryType constant.InventoryType
}

func (p *EndSortInventory) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(1)
	writer.WriteU8(uint8(p.InventoryType))
	return nil
}

func (p *EndSortInventory) Deserialize(reader *stream.StreamReader) {
}

func (p *EndSortInventory) Opcode() uint16 {
	return 0x29
}
