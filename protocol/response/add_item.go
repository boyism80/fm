package response

import (
	"errors"

	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/stream"
)

type AddItem struct {
	IsDrop        bool
	Slot          uint8
	InventoryType constant.InventoryType
	Item          dto.Item
}

func (p *AddItem) Serialize(writer *stream.StreamWriter) error {
	if p.Item == nil {
		return errors.New("item is empty")
	}

	writer.WriteBoolean(p.IsDrop)
	writer.WriteU8(1)
	writer.WriteU8(0)
	writer.WriteU8(uint8(p.InventoryType))
	writer.WriteU8(p.Slot)
	p.Item.Serialize(writer, false, 0)
	return nil
}

func (p *AddItem) Deserialize(reader *stream.StreamReader) {
}

func (p *AddItem) Opcode() uint16 {
	return 0x12
}
