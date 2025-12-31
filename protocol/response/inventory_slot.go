package response

import (
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/stream"
)

type InventoryMode uint8

const (
	INVENTORY_MODE_ADD    InventoryMode = iota // 0: Add (new slot added)
	INVENTORY_MODE_UPDATE                      // 1: Update (existing slot quantity updated)
	INVENTORY_MODE_MOVE                        // 2: Move (slot item moved)
	INVENTORY_MODE_REMOVE                      // 3: Remove (slot item removed)
)

type EquipmentActionType uint8

const (
	EQUIPMENT_ACTION_TYPE_NONE EquipmentActionType = 0
	EQUIPMENT_ACTION_TYPE_OFF  EquipmentActionType = 1
	EQUIPMENT_ACTION_TYPE_ON   EquipmentActionType = 2
)

type SlotItem struct {
	Slot int16
	Item dto.Item
}

type UpdateInventorySlot struct {
	InventoryType constant.InventoryType
	Slot          int16
	Item          dto.Item
}

type AddInventorySlot struct {
	InventoryType constant.InventoryType
	Slot          int16
	Item          dto.Item
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
	writer.WriteBoolean(true)
	writer.WriteU8(1)
	writer.WriteU8(uint8(INVENTORY_MODE_UPDATE))
	writer.WriteU8(uint8(p.InventoryType))
	writer.Write16(p.Slot)
	if p.Item != nil {
		writer.WriteU16(p.Item.GetCount())
	} else {
		writer.WriteU16(0)
	}
	return nil
}

func (p *UpdateInventorySlot) Deserialize(reader *stream.StreamReader) error {
	return nil
}

func (p *UpdateInventorySlot) Opcode() uint16 {
	return 0x12
}

func (p *AddInventorySlot) Serialize(writer *stream.StreamWriter) error {
	writer.WriteBoolean(true)
	writer.WriteU8(1)
	writer.WriteU8(uint8(INVENTORY_MODE_ADD))
	writer.WriteU8(uint8(p.InventoryType))
	writer.WriteU8(uint8(p.Slot))
	p.Item.Serialize(writer, false, 0)
	return nil
}

func (p *AddInventorySlot) Deserialize(reader *stream.StreamReader) error {
	return nil
}

func (p *AddInventorySlot) Opcode() uint16 {
	return 0x12
}

func (p *RemoveInventorySlot) Serialize(writer *stream.StreamWriter) error {
	writer.WriteBoolean(true)
	writer.WriteU8(1)
	writer.WriteU8(uint8(INVENTORY_MODE_REMOVE))
	writer.WriteU8(uint8(p.InventoryType))
	writer.Write16(p.Slot)
	return nil
}

func (p *RemoveInventorySlot) Deserialize(reader *stream.StreamReader) error {
	return nil
}

func (p *RemoveInventorySlot) Opcode() uint16 {
	return 0x12
}

func (p *SwapInventorySlot) Serialize(writer *stream.StreamWriter) error {
	writer.WriteBoolean(true)
	writer.WriteU8(1)
	writer.WriteU8(uint8(INVENTORY_MODE_MOVE))
	writer.WriteU8(uint8(p.InventoryType))
	writer.Write16(p.Source)
	writer.Write16(p.Dest)
	if p.EquipmentAction != EQUIPMENT_ACTION_TYPE_NONE {
		writer.WriteU8(uint8(p.EquipmentAction))
	}
	return nil
}

func (p *SwapInventorySlot) Deserialize(reader *stream.StreamReader) error {
	return nil
}

func (p *SwapInventorySlot) Opcode() uint16 {
	return 0x12
}

func (p *PartialMergeInventorySlot) Serialize(writer *stream.StreamWriter) error {
	writer.WriteBoolean(true)
	writer.WriteU8(2)
	writer.WriteU8(uint8(INVENTORY_MODE_UPDATE))
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

func (p *PartialMergeInventorySlot) Opcode() uint16 {
	return 0x12
}

func (p *FullMergeInventorySlot) Serialize(writer *stream.StreamWriter) error {
	writer.WriteBoolean(true)
	writer.WriteU8(2)
	writer.WriteU8(uint8(INVENTORY_MODE_REMOVE))
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

func (p *FullMergeInventorySlot) Opcode() uint16 {
	return 0x12
}
