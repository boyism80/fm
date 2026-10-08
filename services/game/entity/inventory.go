package entity

import (
	"github.com/boyism80/fm/services/game/constant"
)

type Inventory struct {
	owner    *Character
	Tabs     map[constant.InventoryType]*InventoryTab
	Equipped map[constant.EquipmentPartsType]Equipment
	Rings    Rings
	Meso     int32
}

func NewInventory(owner *Character) *Inventory {
	return &Inventory{
		owner: owner,
		Tabs: map[constant.InventoryType]*InventoryTab{
			constant.InventoryTypeEquipment:    NewInventoryTab(),
			constant.InventoryTypeConsume:      NewInventoryTab(),
			constant.InventoryTypeInstallation: NewInventoryTab(),
			constant.InventoryTypeETC:          NewInventoryTab(),
			constant.InventoryTypeCash:         NewInventoryTab(),
		},
		Equipped: map[constant.EquipmentPartsType]Equipment{},
		Rings: Rings{
			Left: []*Ring{},
			Mid:  []*Ring{},
		},
	}
}
