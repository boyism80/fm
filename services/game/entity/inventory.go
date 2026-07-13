package entity

import (
	"github.com/boyism80/fm/services/game/constant"
)

type Inventory struct {
	owner    *Character
	Tabs     map[constant.InventoryType]*ItemContainer
	Equipped map[constant.EquipmentPartsType]Equipment
	Rings    RingContainer
	Meso     int32
}

func NewInventory(owner *Character) *Inventory {
	return &Inventory{
		owner: owner,
		Tabs: map[constant.InventoryType]*ItemContainer{
			constant.InventoryTypeEquipment:    NewItemContainer(),
			constant.InventoryTypeConsume:      NewItemContainer(),
			constant.InventoryTypeInstallation: NewItemContainer(),
			constant.InventoryTypeETC:          NewItemContainer(),
			constant.InventoryTypeCash:         NewItemContainer(),
		},
		Equipped: map[constant.EquipmentPartsType]Equipment{},
		Rings: RingContainer{
			Left:  []*Ring{},
			Right: []*Ring{},
			Mid:   []*Ring{},
		},
	}
}
