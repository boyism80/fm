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

	resources := h.gs.GetResources()
	if resources == nil {
		log.Printf("Resources not available")
		return fmt.Errorf("resources not available")
	}

	shopID := character.CurrentShopID
	if shopID == 0 {
		return nil
	}

	shop := resources.GetShop(shopID)
	if shop == nil {
		return nil
	}

	if req.Transaction == nil {
		character.CurrentShopID = 0
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
		character.CurrentShopID = 0
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
	if !constant.IsRechargeable(tx.ItemID) {
		price = shopItem.Price * int(tx.Quantity)
	}
	if ch.Inventory.Meso < int32(price) {
		return
	}

	inventoryType := constant.GetInventoryTypeByItemID(tx.ItemID)
	if ch.Inventory.Tabs[inventoryType] == nil {
		return
	}

	quantity := tx.Quantity
	if constant.IsRechargeable(tx.ItemID) {
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

	inventoryType := constant.GetInventoryTypeByItemID(tx.ItemID)
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

	if constant.IsRechargeable(tx.ItemID) {
		quantity = item.GetCount()
	}

	itemQuantity := item.GetCount()
	if itemQuantity == constant.ShopUnspecifiedQuantity {
		itemQuantity = 1
	}
	if quantity > itemQuantity || itemQuantity <= 0 {
		return
	}
	if constant.ItemCategoryOf(tx.ItemID) == constant.ItemCategoryPet {
		return
	}

	itemModel := item.GetModel()
	price := itemModel.GetPrice()
	if constant.IsRechargeable(tx.ItemID) {
		wholePrice := itemModel.GetPrice()
		slotMax := itemModel.GetCapacity()
		if slotMax > 0 {
			price = int(math.Ceil(float64(wholePrice) / float64(slotMax)))
		}
	}

	recvMesos := int32(math.Max(math.Ceil(float64(price)*float64(quantity)), 0))
	if price == -1 || recvMesos <= 0 {
		return
	}

	spec := entity.ExchangeSpec{
		Cost: entity.ExchangeSide{
			Items: map[uint32]uint16{
				tx.ItemID: quantity,
			},
		},
		Reward: entity.ExchangeSide{
			Meso: recvMesos,
		},
	}
	if spec.Valid(ch) != entity.ExchangeOK {
		return
	}

	if quantity >= itemQuantity {
		delete(inventory.Items, tx.Slot)
		ch.Listener.OnRemoveInventorySlot(ch, inventoryType, tx.Slot)
	} else {
		item.Reduce(quantity)
		ch.Listener.OnUpdateInventorySlot(ch, inventoryType, tx.Slot, item)
	}

	ch.Inventory.GainMeso(recvMesos)
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

	itemID := item.GetModel().GetID()
	if !constant.IsRechargeable(itemID) {
		return
	}

	itemModel := item.GetModel()
	slotMax := itemModel.GetCapacity()
	if item.GetCount() >= slotMax {
		return
	}

	price := int(math.Round(float64(itemModel.GetPrice()) * float64(slotMax-item.GetCount())))
	if ch.Inventory.Meso < int32(price) {
		return
	}

	item.SetCount(slotMax)
	ch.Listener.OnUpdateInventorySlot(ch, constant.InventoryTypeConsume, tx.Slot, item)
	ch.Inventory.RemoveMeso(int32(price))
	ch.Listener.OnConfirmShopTransaction(ch, constant.ShopTransactionUpdateOK)
}
