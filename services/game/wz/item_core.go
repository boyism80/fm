package wz

import (
	"slices"

	"github.com/boyism80/fm/services/game/constant"
)

type ItemCore struct {
	ID              uint32
	Name            string
	Price           int
	Cash            bool
	Quest           bool
	SlotMax         uint16
	TradeAvailable  int
	TradeBlock      bool
	Only            bool
	AccountSharable bool
	MakerLevel      int
	MakerBonus      map[string]int16
	InventoryType   constant.InventoryType
	Traits          constant.ItemTrait
}

func (model *ItemCore) setTradeFlags(info *node) {
	model.TradeBlock = info.Int("tradeBlock", 0) == 1
	model.Only = info.Int("only", 0) == 1
	model.AccountSharable = info.Int("accountSharable", 0) == 1
}

var crushRings = []uint32{1112001, 1112002, 1112003, 1112005, 1112006, 1112007, 1112012, 1112015, 1048000, 1048001, 1048002}

var friendshipRings = []uint32{1112800, 1112801, 1112802, 1112810, 1112811, 1112812, 1112816, 1112817, 1049000}

func (model *ItemCore) classify() {
	category := model.ID / 10000
	switch {
	case category >= 100 && category < 200:
		model.InventoryType = constant.InventoryTypeEquipment
	case category >= 200 && category < 300:
		model.InventoryType = constant.InventoryTypeConsume
	case category >= 300 && category < 400:
		model.InventoryType = constant.InventoryTypeInstallation
	case category >= 500 && category < 600:
		model.InventoryType = constant.InventoryTypeCash
	default:
		model.InventoryType = constant.InventoryTypeETC
	}

	switch category {
	case constant.ItemCategoryShuriken, constant.ItemCategoryBullet:
		model.Traits |= constant.ItemTraitRechargeable
	case constant.ItemCategoryPet:
		model.Traits |= constant.ItemTraitPet
	case constant.ItemCategoryMarriage:
		model.Traits |= constant.ItemTraitMarriage
	case 231, 523:
		model.Traits |= constant.ItemTraitShopScanner
	case 504:
		model.Traits |= constant.ItemTraitCashTeleportStone
		if model.ID/1000 != 5040 {
			model.Traits |= constant.ItemTraitVipTeleportStone
		}
	}
	if model.ID >= constant.EngagementRingFirst && model.ID <= constant.EngagementRingLast {
		model.Traits |= constant.ItemTraitEngagementRing
	}
	if model.ID >= constant.WeddingRingFirst && model.ID <= constant.WeddingRingLast {
		model.Traits |= constant.ItemTraitWeddingRing
	}
	if slices.Contains(crushRings, model.ID) {
		model.Traits |= constant.ItemTraitCrushRing
	}
	if slices.Contains(friendshipRings, model.ID) {
		model.Traits |= constant.ItemTraitFriendshipRing
	}
	if model.ID/1000 == 2388 {
		model.Traits |= constant.ItemTraitSpecialMonsterCard
	}
}

func (model *ItemCore) GetID() uint32 {
	return model.ID
}

func (model *ItemCore) GetName() string {
	return model.Name
}

func (model *ItemCore) GetPrice() int {
	return model.Price
}

func (model *ItemCore) IsCash() bool {
	return model.Cash
}

func (model *ItemCore) IsQuest() bool {
	return model.Quest
}

func (model *ItemCore) GetCapacity() uint16 {
	return max(1, model.SlotMax)
}

func (model *ItemCore) IsTradeAvailable() int {
	return model.TradeAvailable
}

func (model *ItemCore) IsTradeBlock() bool {
	return model.TradeBlock
}

func (model *ItemCore) IsOnly() bool {
	return model.Only
}

func (model *ItemCore) IsAccountSharable() bool {
	return model.AccountSharable
}

func (model *ItemCore) GetMakerLevel() int { return model.MakerLevel }

func (model *ItemCore) GetMakerBonus() map[string]int16 { return model.MakerBonus }

func (model *ItemCore) GetInventoryType() constant.InventoryType { return model.InventoryType }
func (model *ItemCore) GetTraits() constant.ItemTrait            { return model.Traits }

func (model *ItemCore) IsConsumeOnPickup() bool {
	return false
}
