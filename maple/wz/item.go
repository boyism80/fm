package wz

import "github.com/boyism80/fm/maple/constant"

type Item interface {
	GetID() uint32
	GetName() string
	GetPrice() int
	SellPrice(quantity uint16) int
	GetCapacity() uint16
	GetMakerLevel() int
	GetMakerBonus() map[string]int16
	GetInventoryType() constant.InventoryType
	GetTraits() constant.ItemTrait
}
