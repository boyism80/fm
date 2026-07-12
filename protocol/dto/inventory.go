package dto

import (
	"github.com/boyism80/fm/services/game/constant"
)

type Inventory struct {
	Containers map[constant.InventoryType]*ItemContainer
	Equipped   map[constant.EquipmentPartsType]*Equipment
	Rings      RingContainer
	Meso       int32
}
