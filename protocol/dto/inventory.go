package dto

import (
	"github.com/boyism80/fm/maple/constant"
)

type Inventory struct {
	Tabs     map[constant.InventoryType]*InventoryTab
	Equipped map[constant.EquipmentPartsType]*Equipment
	Rings    Rings
	Meso     int32
}
