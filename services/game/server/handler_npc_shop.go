package server

import (
	"fmt"
	"log"
	"math"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/entity"
	"github.com/boyism80/fm/services/game/wz"
)

type NpcShop struct {
	gs *GameServer
}

func (NpcShop) New(gs *GameServer) *NpcShop {
	return &NpcShop{
		gs: gs,
	}
}

func (h *NpcShop) Handle(ctx *core.ClientContext, req *request.NpcShop) error {
	client, ok := ctx.Client.(*client.GameClient)
	if !ok {
		log.Printf("Client is not a GameClient")
		return fmt.Errorf("client is not a GameClient")
	}

	character := client.GetCharacter()
	if character == nil {
		return nil
	}
	if character.UsingCashItem() {
		character.Listener.OnUpdateStats(character, nil, true)
		return nil
	}

	resources := h.gs.GetResources()
	if resources == nil {
		log.Printf("Resources not available")
		return fmt.Errorf("resources not available")
	}

	shopID := character.Dialog.ShopID
	if shopID == 0 {
		return nil
	}

	shop := resources.GetShop(shopID)
	if shop == nil {
		return nil
	}

	if req.Transaction == nil {
		character.Dialog.ShopID = 0
		return nil
	}

	switch req.Mode {
	case 0:
		buyTx, ok := req.Transaction.(*request.BuyTransaction)
		if !ok {
			return nil
		}
		h.buy(character, shop, buyTx, resources)
		return nil
	case 1:
		sellTx, ok := req.Transaction.(*request.SellTransaction)
		if !ok {
			return nil
		}
		h.sell(character, sellTx)
		return nil
	case 2:
		rechargeTx, ok := req.Transaction.(*request.RechargeTransaction)
		if !ok {
			return nil
		}
		h.recharge(character, rechargeTx)
		return nil
	default:
		character.Dialog.ShopID = 0
		return nil
	}
}

func (h *NpcShop) buy(ch *entity.Character, shop *wz.Shop, tx *request.BuyTransaction, resources *wz.Resources) {
	if shop == nil || tx.Quantity == 0 {
		return
	}

	var shopItem *wz.ShopItem
	for i := range shop.Items {
		if shop.Items[i].ItemID == tx.ItemID {
			shopItem = &shop.Items[i]
			break
		}
	}
	if shopItem == nil || shopItem.Price <= 0 {
		return
	}

	itemModel, ok := resources.Items[tx.ItemID]
	if !ok {
		return
	}

	price := shopItem.Price
	if itemModel.GetTraits()&constant.ItemTraitRechargeable == 0 {
		price = shopItem.Price * int(tx.Quantity)
	}
	if ch.Inventory.Meso < int32(price) {
		return
	}

	if ch.Inventory.Tabs[itemModel.GetInventoryType()] == nil {
		return
	}

	quantity := tx.Quantity
	if itemModel.GetTraits()&constant.ItemTraitRechargeable != 0 {
		quantity = itemModel.GetCapacity()
	}

	spec := entity.ExchangeSpec{
		Cost: entity.ExchangeSide{
			Meso: int32(price),
		},
		Reward: entity.ExchangeSide{
			Items: map[uint32]uint16{
				tx.ItemID: quantity,
			},
		},
	}
	if ch.Exchange(spec) != entity.ExchangeOK {
		ch.Message("인벤토리 공간이 부족합니다.")
		return
	}

	ch.Listener.OnConfirmShopTransaction(ch, constant.ShopTransactionBuyOK)
}

func (h *NpcShop) sell(ch *entity.Character, tx *request.SellTransaction) {
	quantity := tx.Quantity
	if quantity == constant.ShopUnspecifiedQuantity || quantity == 0 {
		quantity = 1
	}

	itemModel, ok := h.gs.GetResources().Items[tx.ItemID]
	if !ok {
		return
	}
	inventoryType := itemModel.GetInventoryType()
	inventory := ch.Inventory.Tabs[inventoryType]
	if inventory == nil {
		return
	}

	item := inventory.Items[tx.Slot]
	if item == nil {
		return
	}
	if item.GetModel().GetID() != tx.ItemID {
		return
	}

	if itemModel.GetTraits()&constant.ItemTraitRechargeable != 0 {
		quantity = item.GetCount()
	}

	itemQuantity := item.GetCount()
	if itemQuantity == constant.ShopUnspecifiedQuantity {
		itemQuantity = 1
	}
	if quantity > itemQuantity || itemQuantity <= 0 {
		return
	}
	if itemModel.GetTraits()&constant.ItemTraitPet != 0 {
		return
	}

	if itemModel.GetPrice() == -1 {
		return
	}
	recvMesos := itemModel.SellPrice(quantity)
	if recvMesos <= 0 || recvMesos > math.MaxInt32-int(ch.Inventory.Meso) {
		return
	}
	if !ch.Inventory.RemoveItem(inventoryType, tx.Slot, quantity) {
		return
	}

	ch.Inventory.GainMeso(int32(recvMesos))
	ch.Listener.OnConfirmShopTransaction(ch, constant.ShopTransactionUpdateOK)
}

func (h *NpcShop) recharge(ch *entity.Character, tx *request.RechargeTransaction) {
	inventory := ch.Inventory.Tabs[constant.InventoryTypeConsume]
	if inventory == nil {
		return
	}

	item := inventory.Items[tx.Slot]
	if item == nil {
		return
	}

	itemModel, ok := item.GetModel().(*wz.Consume)
	if !ok || itemModel.GetTraits()&constant.ItemTraitRechargeable == 0 {
		return
	}

	slotMax := itemModel.GetCapacity()
	if item.GetCount() >= slotMax {
		return
	}

	price := itemModel.RechargePrice(slotMax - item.GetCount())
	if ch.Inventory.Meso < int32(price) {
		return
	}

	item.SetCount(slotMax)
	ch.Listener.OnInventorySlotUpdated(ch, constant.InventoryTypeConsume, tx.Slot, item)
	ch.Inventory.RemoveMeso(int32(price))
	ch.Listener.OnConfirmShopTransaction(ch, constant.ShopTransactionUpdateOK)
}
