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

type EquipmentActionType uint8

const (
	EquipmentActionTypeNone EquipmentActionType = 0
	EquipmentActionTypeOff  EquipmentActionType = 1
	EquipmentActionTypeOn   EquipmentActionType = 2
)

type SlotItem struct {
	Slot int16
	Item entity.Item
}

type UpdateInventorySlot struct {
	InventoryType constant.InventoryType
	Slot          int16
	Item          entity.Item
}

type AddInventorySlot struct {
	InventoryType constant.InventoryType
	Slot          int16
	Item          entity.Item
}

type RemoveInventorySlot struct {
	InventoryType constant.InventoryType
	Slot          int16
}

type SwapInventorySlot struct {
	InventoryType   constant.InventoryType
	Source, Dest    int16
	EquipmentAction EquipmentActionType
}

type PartialMergeInventorySlot struct {
	InventoryType          constant.InventoryType
	Source, Dest           int16
	SourceCount, DestCount uint16
}

type FullMergeInventorySlot struct {
	InventoryType constant.InventoryType
	Source, Dest  int16
	Count         uint16
}

func (p *UpdateInventorySlot) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU16(0x12)
	writer.WriteBoolean(true)
	writer.WriteU8(1)
	writer.WriteU8(uint8(InventoryModeUpdate))
	writer.WriteU8(uint8(p.InventoryType))
	writer.Write16(p.Slot)
	writer.WriteU16(p.Item.GetCount())
	return nil
}

func (p *UpdateInventorySlot) Deserialize(reader *stream.StreamReader) error {
	return nil
}

func (p *AddInventorySlot) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU16(0x12)
	writer.WriteBoolean(true)
	writer.WriteU8(1)
	writer.WriteU8(uint8(InventoryModeAdd))
	writer.WriteU8(uint8(p.InventoryType))
	writer.WriteU8(uint8(p.Slot))
	p.Item.Serialize(writer, false, 0)
	return nil
}

func (p *AddInventorySlot) Deserialize(reader *stream.StreamReader) error {
	return nil
}

func (p *RemoveInventorySlot) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU16(0x12)
	writer.WriteBoolean(true)
	writer.WriteU8(1)
	writer.WriteU8(uint8(InventoryModeRemove))
	writer.WriteU8(uint8(p.InventoryType))
	writer.Write16(p.Slot)
	return nil
}

func (p *RemoveInventorySlot) Deserialize(reader *stream.StreamReader) error {
	return nil
}

func (p *SwapInventorySlot) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU16(0x12)
	writer.WriteBoolean(true)
	writer.WriteU8(1)
	writer.WriteU8(uint8(InventoryModeMove))
	writer.WriteU8(uint8(p.InventoryType))
	writer.Write16(p.Source)
	writer.Write16(p.Dest)
	if p.EquipmentAction != EquipmentActionTypeNone {
		writer.WriteU8(uint8(p.EquipmentAction))
	}
	return nil
}

func (p *SwapInventorySlot) Deserialize(reader *stream.StreamReader) error {
	return nil
}

func (p *PartialMergeInventorySlot) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU16(0x12)
	writer.WriteBoolean(true)
	writer.WriteU8(2)
	writer.WriteU8(uint8(InventoryModeUpdate))
	writer.WriteU8(uint8(p.InventoryType))
	writer.Write16(p.Source)
	writer.WriteU16(p.SourceCount)
	writer.WriteU8(1)
	writer.WriteU8(uint8(p.InventoryType))
	writer.Write16(p.Dest)
	writer.WriteU16(p.DestCount)
	return nil
}

func (p *PartialMergeInventorySlot) Deserialize(reader *stream.StreamReader) error {
	return nil
}

func (p *FullMergeInventorySlot) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU16(0x12)
	writer.WriteBoolean(true)
	writer.WriteU8(2)
	writer.WriteU8(uint8(InventoryModeRemove))
	writer.WriteU8(uint8(p.InventoryType))
	writer.Write16(p.Source)
	writer.WriteU8(1)
	writer.WriteU8(uint8(p.InventoryType))
	writer.Write16(p.Dest)
	writer.WriteU16(p.Count)
	return nil
}

func (p *FullMergeInventorySlot) Deserialize(reader *stream.StreamReader) error {
	return nil
}
