package entity

import (
	"errors"

	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/game/wz"
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
		item := m.Items[int16(i)]
		if item == nil {
			return uint8(i), true
		}
	}
	return 0, false
}

func (m *Inventory) FindSlot(model wz.Item) (uint8, bool) {
	for i := 1; i <= int(m.SlotLimit); i++ {
		item := m.Items[int16(i)]
		if item == nil {
			continue
		}

		if item.GetModel().GetID() != model.GetID() {
			continue
		}

		cap := model.GetCapacity() - item.GetCount()
		if cap == 0 {
			continue
		}

		return uint8(i), true
	}

	return m.NextSlot()
}

func (m *Inventory) EmptySlotCount() uint16 {
	count := uint16(0)
	for i := 1; i <= int(m.SlotLimit); i++ {
		if _, ok := m.Items[int16(i)]; !ok {
			count += 1
		}
	}

	return count
}

func (m *Inventory) IsFree(model wz.Item, count uint16) bool {
	space := uint16(0)
	slot, ok := m.FindSlot(model)
	if ok {
		exists, ok := m.Items[int16(slot)]
		if ok {
			if model.GetCapacity() > exists.GetCount() {
				space += model.GetCapacity() - exists.GetCount()
			}
		}
	}

	space += m.EmptySlotCount() * model.GetCapacity()
	return space >= count
}

// GetItem returns the item at the specified slot
func (m *Inventory) GetItem(slot uint8) Item {
	return m.Items[int16(slot)]
}

// AddItem adds an item to the specified slot
func (m *Inventory) AddItem(slot uint8, item Item) error {
	if slot < 1 || slot > m.SlotLimit {
		return ErrSlotNotFound
	}

	if m.Items[int16(slot)] != nil {
		return ErrSlotAlreadyOccupied
	}

	m.Items[int16(slot)] = item
	return nil
}

// RemoveItem removes an item from the specified slot
func (m *Inventory) RemoveItem(slot uint8) error {
	if slot < 1 || slot > m.SlotLimit {
		return ErrSlotNotFound
	}

	if m.Items[int16(slot)] == nil {
		return ErrSourceSlotEmpty
	}

	delete(m.Items, int16(slot))
	return nil
}

// FindById finds an item by its ID in the inventory
func (m *Inventory) FindById(itemID uint32) Item {
	for _, item := range m.Items {
		if item != nil && item.GetModel().GetID() == itemID {
			return item
		}
	}
	return nil
}
