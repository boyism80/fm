package game

import (
	"math"

	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/maple/wz"
)

func (d *Dialog) Buy(shop *wz.Shop, itemID uint32, count uint16) {
	if count == 0 {
		return
	}
	shopItem := shop.Find(itemID)
	if shopItem == nil || shopItem.Price <= 0 {
		return
	}
	model := d.owner.itemModel(itemID)
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
	switch d.owner.Exchange(spec) {
	case ExchangeOK:
		d.owner.listener.OnConfirmShopTransaction(d.owner, constant.ShopTransactionBuyOK)
	case ExchangeLackCapacity:
		d.owner.Message("인벤토리 공간이 부족합니다.")
	}
}

func (d *Dialog) Sell(slot int16, itemID uint32, count uint16) {
	model := d.owner.itemModel(itemID)
	if model == nil {
		return
	}
	if model.GetTraits()&(constant.ItemTraitPet|constant.ItemTraitNotSale) != 0 || model.GetPrice() == -1 {
		return
	}
	invType := model.GetInventoryType()
	item := d.owner.inventory.GetItem(invType, slot)
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
	if meso <= 0 || meso > math.MaxInt32-int(d.owner.inventory.meso) {
		return
	}

	if !d.owner.inventory.RemoveItem(invType, slot, count) {
		return
	}
	d.owner.inventory.GainMeso(int32(meso))
	d.owner.listener.OnConfirmShopTransaction(d.owner, constant.ShopTransactionUpdateOK)
}

func (d *Dialog) Recharge(slot int16) {
	item := d.owner.inventory.GetItem(constant.InventoryTypeConsume, slot)
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
	if price > 0 && !d.owner.inventory.RemoveMeso(price) {
		return
	}
	item.SetCount(capacity)
	d.owner.listener.OnInventorySlotUpdated(d.owner, constant.InventoryTypeConsume, slot, item)
	d.owner.listener.OnConfirmShopTransaction(d.owner, constant.ShopTransactionUpdateOK)
}
