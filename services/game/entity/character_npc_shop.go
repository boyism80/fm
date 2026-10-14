package entity

import (
	"math"

	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/wz"
)

func (ch *Character) BuyShopItem(shop *wz.Shop, itemID uint32, count uint16) {
	if count == 0 {
		return
	}
	shopItem := shop.Find(itemID)
	if shopItem == nil || shopItem.Price <= 0 {
		return
	}
	model := ch.itemModel(itemID)
	if model == nil {
		return
	}

	price := shopItem.Price * int(count)
	if model.GetTraits()&constant.ItemTraitRechargeable != 0 {
		price = shopItem.Price
		count = model.GetCapacity()
	}
	spec := ExchangeSpec{
		Cost:   ExchangeSide{Meso: int32(price)},
		Reward: ExchangeSide{Items: map[uint32]uint16{itemID: count}},
	}
	switch ch.Exchange(spec) {
	case ExchangeOK:
		ch.Listener.OnConfirmShopTransaction(ch, constant.ShopTransactionBuyOK)
	case ExchangeLackCapacity:
		ch.Message("인벤토리 공간이 부족합니다.")
	}
}

func (ch *Character) SellShopItem(slot int16, itemID uint32, count uint16) {
	model := ch.itemModel(itemID)
	if model == nil {
		return
	}
	if model.GetTraits()&(constant.ItemTraitPet|constant.ItemTraitNotSale) != 0 || model.GetPrice() == -1 {
		return
	}
	invType := model.GetInventoryType()
	item := ch.Inventory.GetItem(invType, slot)
	if item == nil || item.GetModel().GetID() != itemID {
		return
	}

	if count == constant.ShopUnspecifiedQuantity || count == 0 {
		count = 1
	}
	if model.GetTraits()&constant.ItemTraitRechargeable != 0 {
		count = item.GetCount()
	}
	if count > item.GetCount() {
		return
	}
	meso := model.SellPrice(count)
	if meso <= 0 || meso > math.MaxInt32-int(ch.Inventory.Meso) {
		return
	}

	if !ch.Inventory.RemoveItem(invType, slot, count) {
		return
	}
	ch.Inventory.GainMeso(int32(meso))
	ch.Listener.OnConfirmShopTransaction(ch, constant.ShopTransactionUpdateOK)
}

func (ch *Character) RechargeShopItem(slot int16) {
	item := ch.Inventory.GetItem(constant.InventoryTypeConsume, slot)
	if item == nil {
		return
	}
	model, ok := item.GetModel().(*wz.Consume)
	if !ok || model.GetTraits()&constant.ItemTraitRechargeable == 0 {
		return
	}
	capacity := model.GetCapacity()
	if item.GetCount() >= capacity {
		return
	}

	price := int32(model.RechargePrice(capacity - item.GetCount()))
	if price > 0 && !ch.Inventory.RemoveMeso(price) {
		return
	}
	item.SetCount(capacity)
	ch.Listener.OnInventorySlotUpdated(ch, constant.InventoryTypeConsume, slot, item)
	ch.Listener.OnConfirmShopTransaction(ch, constant.ShopTransactionUpdateOK)
}
