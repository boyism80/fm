package entity

import (
	"errors"

	"github.com/boyism80/fm/common/stream"
	"github.com/boyism80/fm/game/constant"
)

var (
	ErrInventoryFull       = errors.New("inventory is full")
	ErrSourceSlotEmpty     = errors.New("source slot is empty")
	ErrSlotNotFound        = errors.New("slot not found")
	ErrSlotAlreadyOccupied = errors.New("slot already occupied")
)

type Inventory struct {
	Items     map[int16]Item
	SlotLimit uint8
	Type      constant.InventoryType
}

func NewInventory(typ constant.InventoryType) *Inventory {
	return &Inventory{
		Type:      typ,
		SlotLimit: 32,
		Items:     map[int16]Item{},
	}
}

func (m *Inventory) NextSlot() (uint8, error) {
	for i := 1; i <= int(m.SlotLimit); i++ {
		if _, ok := m.Items[int16(i)]; !ok {
			return uint8(i), nil
		}
	}
	return 0, errors.New("inventory is full")
}

func (m *Inventory) Serialize(sw *stream.StreamWriter) {

	for slot, item := range m.Items {
		if item == nil {
			continue
		}

		item.Serialize(sw, true, slot)
	}
	sw.WriteU8(0)
}
