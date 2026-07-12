package entity

import (
	"errors"

	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/wz"
)

var (
	ErrInventoryFull        = errors.New("inventory is full")
	ErrSourceSlotEmpty      = errors.New("source slot is empty")
	ErrSlotNotFound         = errors.New("slot not found")
	ErrSlotAlreadyOccupied  = errors.New("slot already occupied")
	ErrInvalidEquipmentPart = errors.New("invalid equipment part")
	ErrItemNotEquipment     = errors.New("item is not equipment")
)

type ItemContainer struct {
	Items     map[int16]Item
	SlotLimit uint8
	Type      constant.InventoryType
}

func NewItemContainer(typ constant.InventoryType) *ItemContainer {
	return &ItemContainer{
		Type:      typ,
		SlotLimit: 32,
		Items:     map[int16]Item{},
	}
}

func (m *ItemContainer) NextSlot() (uint8, bool) {
	for i := 1; i <= int(m.SlotLimit); i++ {
		item := m.Items[int16(i)]
		if item == nil {
			return uint8(i), true
		}
	}
	return 0, false
}

func (m *ItemContainer) FindSlot(model wz.Item) (uint8, bool) {
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

func (m *ItemContainer) EmptySlotCount() uint16 {
	count := uint16(0)
	for i := 1; i <= int(m.SlotLimit); i++ {
		if _, ok := m.Items[int16(i)]; !ok {
			count += 1
		}
	}

	return count
}

func (m *ItemContainer) IsFree(model wz.Item, count uint16) bool {
	if model == nil || count == 0 {
		return true
	}
	modelOf := func(id uint32) wz.Item {
		if id == model.GetID() {
			return model
		}
		return nil
	}
	reward := map[uint32]uint16{model.GetID(): count}
	return m.validateItemExchange(nil, reward, modelOf) == ExchangeOK
}

func (m *ItemContainer) Get(slot uint8) Item {
	return m.Items[int16(slot)]
}

func (m *ItemContainer) Add(slot uint8, item Item) error {
	if slot < 1 || slot > m.SlotLimit {
		return ErrSlotNotFound
	}

	if m.Items[int16(slot)] != nil {
		return ErrSlotAlreadyOccupied
	}

	m.Items[int16(slot)] = item
	return nil
}

func (m *ItemContainer) Remove(slot uint8) error {
	if slot < 1 || slot > m.SlotLimit {
		return ErrSlotNotFound
	}

	if m.Items[int16(slot)] == nil {
		return ErrSourceSlotEmpty
	}

	delete(m.Items, int16(slot))
	return nil
}

func (m *ItemContainer) Find(itemID uint32) Item {
	for _, item := range m.Items {
		if item != nil && item.GetModel().GetID() == itemID {
			return item
		}
	}
	return nil
}
