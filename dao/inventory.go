package dao

import (
	"errors"
	"sort"
)

var (
	ErrInventoryFull       = errors.New("inventory is full")
	ErrSourceSlotEmpty     = errors.New("source slot is empty")
	ErrSlotNotFound        = errors.New("slot not found")
	ErrSlotAlreadyOccupied = errors.New("slot already occupied")
)

type MapleInventoryType int

const (
	InventoryTypeEquip MapleInventoryType = iota
	InventoryTypeUse
	InventoryTypeSetUp
	InventoryTypeEtc
	InventoryTypeCash
	InventoryTypeEquipped
)

type MapleInventory struct {
	inventory map[int16]*Item
	SlotLimit uint8
	Type      MapleInventoryType
}

func NewMapleInventory(typ MapleInventoryType) *MapleInventory {
	return &MapleInventory{
		inventory: make(map[int16]*Item),
		SlotLimit: 0,
		Type:      typ,
	}
}

func (m *MapleInventory) AddSlot(slot uint8) {
	m.SlotLimit += slot
	if m.SlotLimit > 96 {
		m.SlotLimit = 96
	}
}

func (m *MapleInventory) SetSlotLimit(slot uint8) {
	if slot > 96 {
		slot = 96
	}
	m.SlotLimit = slot
}

func (m *MapleInventory) GetSlotLimit() uint8 {
	return m.SlotLimit
}

func (m *MapleInventory) FindByID(itemID uint32) *Item {
	for _, item := range m.inventory {
		if item.Id == itemID {
			return item
		}
	}
	return nil
}

func (m *MapleInventory) FindByUniqueID(uniqueID uint64) *Item {
	for _, item := range m.inventory {
		if item.UniqueID == uniqueID {
			return item
		}
	}
	return nil
}

func (m *MapleInventory) FindByInventoryID(invID uint64, itemID uint32) *Item {
	for _, item := range m.inventory {
		if item.InventoryID == invID && item.Id == itemID {
			return item
		}
	}
	return m.FindByID(itemID)
}

func (m *MapleInventory) FindByInventoryIDOnly(invID uint64, itemID uint32) *Item {
	for _, item := range m.inventory {
		if item.InventoryID == invID && item.Id == itemID {
			return item
		}
	}
	return nil
}

func (m *MapleInventory) CountByID(itemID uint32) int {
	count := 0
	for _, item := range m.inventory {
		if item.Id == itemID {
			count += int(item.Quantity)
		}
	}
	return count
}

func (m *MapleInventory) ListByID(itemID uint32) []*Item {
	var ret []*Item
	for _, item := range m.inventory {
		if item.Id == itemID {
			ret = append(ret, item)
		}
	}
	if len(ret) > 1 {
		sort.Slice(ret, func(i, j int) bool {
			return ret[i].Parts < ret[j].Parts
		})
	}
	return ret
}

func (m *MapleInventory) NewList() []*Item {
	ret := []*Item{}
	for _, item := range m.inventory {
		ret = append(ret, item)
	}
	return ret
}

func (m *MapleInventory) ListIDs() []uint32 {
	idSet := make(map[uint32]struct{})
	for _, item := range m.inventory {
		idSet[item.Id] = struct{}{}
	}
	ret := make([]uint32, 0, len(idSet))
	for id := range idSet {
		ret = append(ret, id)
	}
	sort.Slice(ret, func(i, j int) bool {
		return ret[i] < ret[j]
	})
	return ret
}

func (m *MapleInventory) AddItem(item *Item) (int16, error) {
	slotID, err := m.GetNextFreeSlot()
	if err != nil {
		return 0, err
	}
	m.inventory[slotID] = item
	item.Parts = slotID
	return slotID, nil
}

func (m *MapleInventory) AddFromDB(item *Item) error {
	if item.Parts != 0 && m.Type == InventoryTypeEquipped {
		return nil
	}
	if _, ok := m.inventory[item.Parts]; ok {
		return ErrSlotAlreadyOccupied
	}
	m.inventory[item.Parts] = item
	return nil
}

func (m *MapleInventory) Move(sSlot, dSlot int16, slotMax uint16) error {
	source := m.inventory[sSlot]
	target := m.inventory[dSlot]

	if source == nil {
		return ErrSourceSlotEmpty
	}

	if target == nil {
		if dSlot != 0 && m.Type != InventoryTypeEquipped {
			return nil
		}

		source.Parts = dSlot
		m.inventory[dSlot] = source
		delete(m.inventory, sSlot)
	} else if source.Id == target.Id && source.Owner == target.Owner && source.Expiration == target.Expiration {
		if m.Type == InventoryTypeEquip || m.Type == InventoryTypeCash {
			m.swap(source, target)
		} else if source.Quantity+target.Quantity > slotMax {
			source.Quantity = (source.Quantity + target.Quantity) - slotMax
			target.Quantity = slotMax
		} else {
			target.Quantity += source.Quantity
			delete(m.inventory, sSlot)
		}
	} else {
		m.swap(source, target)
	}
	return nil
}

func (m *MapleInventory) swap(a, b *Item) {
	delete(m.inventory, a.Parts)
	delete(m.inventory, b.Parts)
	a.Parts, b.Parts = b.Parts, a.Parts
	m.inventory[a.Parts] = a
	m.inventory[b.Parts] = b
}

func (m *MapleInventory) GetItem(slot int16) *Item {
	return m.inventory[slot]
}

func (m *MapleInventory) RemoveItem(slot int16) {
	m.RemoveItemWithQuantity(slot, 1, false)
}

func (m *MapleInventory) RemoveItemWithQuantity(slot int16, quantity uint16, allowZero bool) error {
	item := m.inventory[slot]
	if item == nil {
		return ErrSlotNotFound
	}
	item.Quantity -= quantity
	if item.Quantity == 0 && !allowZero {
		m.RemoveSlot(slot)
	}
	return nil
}

func (m *MapleInventory) RemoveSlot(slot int16) {
	delete(m.inventory, slot)
}

func (m *MapleInventory) IsFull() bool {
	return uint8(len(m.inventory)) >= m.SlotLimit
}

func (m *MapleInventory) IsFullWithMargin(margin uint8) bool {
	return uint8(len(m.inventory))+margin >= m.SlotLimit
}

func (m *MapleInventory) GetNextFreeSlot() (int16, error) {
	if m.IsFull() {
		return 0, ErrInventoryFull
	}
	for i := int16(1); i <= int16(m.SlotLimit); i++ {
		if _, ok := m.inventory[i]; !ok {
			return i, nil
		}
	}
	return 0, ErrInventoryFull
}

func (m *MapleInventory) GetNumFreeSlot() int16 {
	if m.IsFull() {
		return 0
	}
	free := int16(0)
	for i := int16(1); i <= int16(m.SlotLimit); i++ {
		if _, ok := m.inventory[i]; !ok {
			free++
		}
	}
	return free
}

func (m *MapleInventory) GetType() MapleInventoryType {
	return m.Type
}
