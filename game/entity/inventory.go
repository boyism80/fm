package entity

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

type Inventory struct {
	Items     map[int16]Item
	SlotLimit uint8
	Type      MapleInventoryType
}

func NewMapleInventory(typ MapleInventoryType) *Inventory {
	return &Inventory{
		Type:      typ,
		SlotLimit: 32,
		Items:     map[int16]Item{},
	}
}

func (m *Inventory) NewList() []Item {
	ret := []Item{}
	for _, item := range m.Items {
		ret = append(ret, item)
	}
	return ret
}

func (m *Inventory) ListIDs() []uint32 {
	idSet := make(map[uint32]struct{})
	for _, item := range m.Items {
		idSet[item.GetTemplate().GetID()] = struct{}{}
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
