package game

import (
	"github.com/boyism80/fm/maple/constant"
)

type Inventory struct {
	owner    *Character
	tabs     map[constant.InventoryType]*InventoryTab
	equipped map[constant.EquipmentPartsType]Equipment
	rings    Rings
	meso     int32
}

func NewInventory(owner *Character) *Inventory {
	return &Inventory{
		owner: owner,
		tabs: map[constant.InventoryType]*InventoryTab{
			constant.InventoryTypeEquipment:    NewInventoryTab(),
			constant.InventoryTypeConsume:      NewInventoryTab(),
			constant.InventoryTypeInstallation: NewInventoryTab(),
			constant.InventoryTypeETC:          NewInventoryTab(),
			constant.InventoryTypeCash:         NewInventoryTab(),
		},
		equipped: map[constant.EquipmentPartsType]Equipment{},
		rings: Rings{
			left: []*Ring{},
			mid:  []*Ring{},
		},
	}
}

func (inv *Inventory) look() equippedLook {
	look := make(equippedLook, len(inv.equipped))
	for parts, equipment := range inv.equipped {
		if equipment != nil {
			look[parts] = equipment.GetModel().GetID()
		}
	}
	return look
}
