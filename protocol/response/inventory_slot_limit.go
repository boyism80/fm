package response

import (
	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/stream"
)

type InventorySlotLimit struct {
	InventoryType constant.InventoryType
	Limit         uint8
}

func (p *InventorySlotLimit) Opcode() uint16 { return 0x13 }

func (p *InventorySlotLimit) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(uint8(p.InventoryType))
	writer.WriteU8(p.Limit)
	return nil
}

func (p *InventorySlotLimit) Deserialize(reader *stream.StreamReader) {
	p.InventoryType = constant.InventoryType(reader.ReadU8())
	p.Limit = reader.ReadU8()
}
