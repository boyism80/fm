package entity

import (
	"errors"

	"github.com/boyism80/fm/common/stream"
)

var (
	ErrInventoryFull       = errors.New("inventory is full")
	ErrSourceSlotEmpty     = errors.New("source slot is empty")
	ErrSlotNotFound        = errors.New("slot not found")
	ErrSlotAlreadyOccupied = errors.New("slot already occupied")
)

type InventoryType uint8

const (
	InventoryTypeEquip        InventoryType = 1
	InventoryTypeConsume      InventoryType = 2
	InventoryTypeInstallation InventoryType = 3
	InventoryTypeEtc          InventoryType = 4
	InventoryTypeCash         InventoryType = 5
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

func NewInventory(typ InventoryType) *Inventory {
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
