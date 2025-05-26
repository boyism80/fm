package entity

import (
	"errors"

	"github.com/boyism80/fm/common/stream"
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/game/data"
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

func (m *Inventory) NextSlot() (uint8, bool) {
	for i := 1; i <= int(m.SlotLimit); i++ {
		if _, ok := m.Items[int16(i)]; !ok {
			return uint8(i), true
		}
	}
	return 0, false
}

func (m *Inventory) FindSlot(spec data.ItemSpec) (uint8, bool) {
	for i := 1; i <= int(m.SlotLimit); i++ {
		item, ok := m.Items[int16(i)]
		if !ok {
			continue
		}

		if item.GetSpec().GetID() != spec.GetID() {
			continue
		}

		cap := spec.GetCapacity() - item.GetCount()
		if cap == 0 {
			continue
		}

		return uint8(i), true
	}

	return m.NextSlot()
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
