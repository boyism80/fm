package entity

import (
	"github.com/boyism80/fm/services/game/constant"
)

type Inventory struct {
	owner      *Character
	Containers map[constant.InventoryType]*ItemContainer
	Equipped   map[constant.EquipmentPartsType]Equipment
	Rings      RingContainer
	Meso       int32
}

func NewInventory(owner *Character) *Inventory {
	return &Inventory{
		owner: owner,
		Containers: map[constant.InventoryType]*ItemContainer{
			constant.InventoryTypeEquipment:    NewItemContainer(constant.InventoryTypeEquipment),
			constant.InventoryTypeConsume:      NewItemContainer(constant.InventoryTypeConsume),
			constant.InventoryTypeInstallation: NewItemContainer(constant.InventoryTypeInstallation),
			constant.InventoryTypeETC:          NewItemContainer(constant.InventoryTypeETC),
			constant.InventoryTypeCash:         NewItemContainer(constant.InventoryTypeCash),
		},
		Equipped: map[constant.EquipmentPartsType]Equipment{},
		Rings: RingContainer{
			Left:  []*Ring{},
			Right: []*Ring{},
			Mid:   []*Ring{},
		},
	}
}
