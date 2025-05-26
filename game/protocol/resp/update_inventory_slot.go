package resp

import (
	"github.com/boyism80/fm/common/stream"
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/game/entity"
)

type InventoryMode uint8

const (
	InventoryModeAdd    InventoryMode = iota // 0: Add (새 슬롯 추가)
	InventoryModeUpdate                      // 1: Update (기존 슬롯 수량 갱신)
	InventoryModeMove                        // 2: Move (슬롯 간 아이템 이동)
	InventoryModeRemove                      // 3: Remove (슬롯 아이템 제거)
)

type SlotItem struct {
	Slot int16
	Item entity.Item
}

type UpdateInventorySlot struct {
	InventoryType constant.InventoryType
	Mode          InventoryMode
	IsDrop        bool
	Items         []SlotItem
}

func (p *UpdateInventorySlot) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU16(0x12)
	writer.WriteBoolean(p.IsDrop)
	writer.WriteU8(uint8(len(p.Items)))
	writer.WriteU8(uint8(p.Mode))
	writer.WriteU8(uint8(p.InventoryType))

	switch p.Mode {
	case InventoryModeAdd:
		for _, v := range p.Items {
			writer.WriteU8(uint8(v.Slot))
			v.Item.Serialize(writer, false, 0)
		}

	case InventoryModeUpdate:
		for _, v := range p.Items {
			writer.Write16(v.Slot)
			writer.WriteU16(v.Item.GetCount())
		}

	case InventoryModeMove:
		// TODO: item move

	case InventoryModeRemove:
		for _, v := range p.Items {
			writer.Write16(v.Slot)
		}
	}
	return nil
}

func (s *UpdateInventorySlot) Deserialize(reader *stream.StreamReader) error {
	return nil
}
