package resp

import (
	"errors"

	"github.com/boyism80/fm/common/stream"
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/game/entity"
)

type AddItem struct {
	IsDrop        bool
	Slot          uint8
	InventoryType constant.InventoryType
	Item          entity.Item
}

func (p *AddItem) Serialize(writer *stream.StreamWriter) error {
	if p.Item == nil {
		return errors.New("item is empty")
	}

	writer.WriteU16(0x12)
	writer.WriteBoolean(p.IsDrop)
	writer.WriteU8(1)
	writer.WriteU8(0)
	writer.WriteU8(uint8(p.InventoryType))
	writer.WriteU8(p.Slot)
	p.Item.Serialize(writer, false, 0)
	return nil
}

func (p *AddItem) Deserialize(reader *stream.StreamReader) error {
	return nil
}
