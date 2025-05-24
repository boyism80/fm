package resp

import (
	"github.com/boyism80/fm/common/stream"
	"github.com/boyism80/fm/game/constant"
)

type UpdateItemCount struct {
	InventoryType constant.InventoryType
	Slot          uint16
	Count         uint16
}

func (p *UpdateItemCount) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU16(0x12)
	writer.WriteU8(1)
	writer.WriteU8(1)
	if p.Count > 0 {
		writer.WriteU8(1)
	} else {
		writer.WriteU8(3)
	}
	writer.WriteU8(uint8(p.InventoryType))
	writer.WriteU16(p.Slot)
	writer.WriteU16(p.Count)
	return nil
}

func (s *UpdateItemCount) Deserialize(reader *stream.StreamReader) error {
	return nil
}
