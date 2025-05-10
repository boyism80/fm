package entity

import (
	"errors"
	"sort"

	"github.com/boyism80/fm/common/stream"
)

var (
	ErrInventoryFull       = errors.New("inventory is full")
	ErrSourceSlotEmpty     = errors.New("source slot is empty")
	ErrSlotNotFound        = errors.New("slot not found")
	ErrSlotAlreadyOccupied = errors.New("slot already occupied")
)

type InventoryType int

const (
	InventoryTypeEquip InventoryType = iota
	InventoryTypeUse
	InventoryTypeSetUp
	InventoryTypeEtc
	InventoryTypeCash
)

type ItemType int

const (
	ItemTypeEquipment ItemType = 1
	ItemTypeEtc       ItemType = 2
)

type EquipmentPartsType int

const (
	EquipmentPartsWeapon EquipmentPartsType = -11
	EquipmentPartsShield EquipmentPartsType = -10
)

type Inventory struct {
	Items     map[int16]Item
	SlotLimit uint8
	Type      InventoryType
}

func NewMapleInventory(typ InventoryType) *Inventory {
	return &Inventory{
		Type:      typ,
		SlotLimit: 32,
		Items:     map[int16]Item{},
	}
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

func (m *Inventory) Serialize(sw *stream.StreamWriter) {

	itemType := ItemTypeEtc
	if m.Type == InventoryTypeEquip {
		itemType = ItemTypeEquipment
	}

	for slot, item := range m.Items {
		if item == nil {
			continue
		}

		item.Serialize(sw, false, false, true, slot, itemType)
	}
	sw.WriteU8(0)
}
