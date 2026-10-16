package game

import (
	"errors"

	"github.com/boyism80/fm/maple/constant"
)

var (
	ErrInventoryFull         = errors.New("inventory is full")
	ErrSourceSlotEmpty       = errors.New("source slot is empty")
	ErrSlotNotFound          = errors.New("slot not found")
	ErrSlotAlreadyOccupied   = errors.New("slot already occupied")
	ErrInvalidEquipmentPart  = errors.New("invalid equipment part")
	ErrEquipmentRequirements = errors.New("equipment requirements not met")
	ErrItemNotEquipment      = errors.New("item is not equipment")
)

type InventoryTab struct {
	items     map[int16]Item
	slotLimit uint8
}

func NewInventoryTab() *InventoryTab {
	return &InventoryTab{
		slotLimit: 32,
		items:     map[int16]Item{},
	}
}

func (m *InventoryTab) NextSlot() (uint8, bool) {
	for i := 1; i <= int(m.slotLimit); i++ {
		item := m.items[int16(i)]
		if item == nil {
			return uint8(i), true
		}
	}
	return 0, false
}

// Stackable reports whether a and b may share one slot. Rechargeables keep their own slot even when partly used.
func (m *InventoryTab) Stackable(a, b Item) bool {
	model := a.GetModel()
	return model.GetID() == b.GetModel().GetID() &&
		model.GetCapacity() > 1 &&
		model.GetTraits()&constant.ItemTraitRechargeable == 0 &&
		a.GetExpiration().Equal(b.GetExpiration())
}

func (m *InventoryTab) FindSlot(incoming Item) (uint8, bool) {
	for i := 1; i <= int(m.slotLimit); i++ {
		item := m.items[int16(i)]
		if item == nil || !m.Stackable(item, incoming) {
			continue
		}
		if item.GetCount() >= item.GetModel().GetCapacity() {
			continue
		}

		return uint8(i), true
	}

	return m.NextSlot()
}

func (m *InventoryTab) EmptySlotCount() uint16 {
	count := uint16(0)
	for i := 1; i <= int(m.slotLimit); i++ {
		if _, ok := m.items[int16(i)]; !ok {
			count += 1
		}
	}

	return count
}

func (m *InventoryTab) Get(slot uint8) Item {
	return m.items[int16(slot)]
}

func (m *InventoryTab) Add(slot uint8, item Item) error {
	if slot < 1 || slot > m.slotLimit {
		return ErrSlotNotFound
	}

	if m.items[int16(slot)] != nil {
		return ErrSlotAlreadyOccupied
	}

	m.items[int16(slot)] = item
	return nil
}

func (m *InventoryTab) Remove(slot uint8) error {
	if slot < 1 || slot > m.slotLimit {
		return ErrSlotNotFound
	}

	if m.items[int16(slot)] == nil {
		return ErrSourceSlotEmpty
	}

	delete(m.items, int16(slot))
	return nil
}

func (m *InventoryTab) Find(itemID uint32) Item {
	for _, item := range m.items {
		if item != nil && item.GetModel().GetID() == itemID {
			return item
		}
	}
	return nil
}
