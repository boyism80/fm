package wz

import "github.com/boyism80/fm/services/game/constant"

type Item interface {
	GetID() uint32
	GetName() string
	GetPrice() int
	SellPrice(quantity uint16) int
	IsCash() bool
	IsQuest() bool
	GetCapacity() uint16
	IsTradeAvailable() int
	IsTradeBlock() bool
	IsOnly() bool
	IsAccountSharable() bool
	GetMakerLevel() int
	GetMakerBonus() map[string]int16
	GetInventoryType() constant.InventoryType
	GetTraits() constant.ItemTrait

	IsConsumeOnPickup() bool
}
