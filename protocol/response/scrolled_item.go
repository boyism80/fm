package response

import (
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/stream"
)

type ScrolledItem struct {
	ScrollInventoryType constant.InventoryType
	ScrollSlot          int16
	ScrollCount         uint16
	UpgradedSlot        int16
	Destroyed           bool
	Potential           bool
	UpgradedItem        dto.Item
}

func (p *ScrolledItem) Opcode() uint16 {
	return 0x12
}

func (p *ScrolledItem) Serialize(writer *stream.StreamWriter) error {
	writer.WriteBoolean(true)
	if p.Destroyed {
		writer.WriteU8(2)
	} else {
		writer.WriteU8(3)
	}
	if p.ScrollCount > 0 {
		writer.WriteU8(1)
	} else {
		writer.WriteU8(3)
	}
	writer.WriteU8(uint8(p.ScrollInventoryType))
	writer.Write16(p.ScrollSlot)
	if p.ScrollCount > 0 {
		writer.WriteU16(p.ScrollCount)
	}
	writer.WriteU8(3)
	if !p.Destroyed {
		writer.WriteU8(uint8(constant.InventoryTypeEquipment))
		writer.Write16(p.UpgradedSlot)
		writer.WriteU8(0)
	}
	writer.WriteU8(uint8(constant.InventoryTypeEquipment))
	writer.Write16(p.UpgradedSlot)
	if !p.Destroyed && p.UpgradedItem != nil {
		p.UpgradedItem.Serialize(writer, dto.ItemSerializeOption{
			Trade:    false,
			Slot:     0,
			SlotMode: dto.SlotEncodeOmit,
		})
	}
	if !p.Potential {
		writer.WriteU8(1)
	}
	return nil
}

func (p *ScrolledItem) Deserialize(reader *stream.StreamReader) {
}
