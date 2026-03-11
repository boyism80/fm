package server

import (
	"fmt"
	"log"
	"math"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/game/client"
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/game/entity"
	"github.com/boyism80/fm/game/wz"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/types"
)

// NpcShop handles NPC shop transaction packet requests
type NpcShop struct {
	gs     *GameServer
	opcode byte
}

func (NpcShop) New(gs *GameServer) *NpcShop {
	return &NpcShop{
		gs:     gs,
		opcode: 0x2C,
	}
}

func (h *NpcShop) GetOpcode() byte {
	return h.opcode
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
		return h.handleBuy(character, shop, buyTx, resources)
	case 1:
		sellTx, ok := req.Transaction.(*request.SellTransaction)
		if !ok {
			return nil
		}
		return h.handleSell(character, shop, sellTx, resources)
	case 2:
		rechargeTx, ok := req.Transaction.(*request.RechargeTransaction)
		if !ok {
			return nil
		}
		return h.handleRecharge(character, shop, rechargeTx, resources)
	default:
		character.CurrentShopID = 0
		return nil
	}
}

func (h *NpcShop) handleBuy(character *entity.Character, shop *wz.Shop, tx *request.BuyTransaction, resources *wz.Resources) error {
	if tx.Quantity <= 0 {
		return nil
	}

	var shopItem *wz.ShopItem
	for i := range shop.Items {
		if shop.Items[i].ItemID == tx.ItemID {
			shopItem = &shop.Items[i]
			break
		}
	}

	if shopItem == nil {
		return nil
	}

	itemModel, ok := resources.Items[tx.ItemID]
	if !ok {
		return nil
	}

	if shopItem.Price <= 0 {
		return nil
	}

	price := shopItem.Price
	if h.isRechargable(tx.ItemID) {
		price = shopItem.Price
	} else {
		price = shopItem.Price * int(tx.Quantity)
	}

	if character.Meso < int32(price) {
		return nil
	}

	inventoryType := h.getItemInventoryType(tx.ItemID, itemModel)
	inventory := character.Inventory[inventoryType]
	if inventory == nil {
		return nil
	}

	if !inventory.IsFree(itemModel, tx.Quantity) {
		character.Message("인벤토리 공간이 부족합니다.")
		return nil
	}

	quantity := tx.Quantity
	if h.isRechargable(tx.ItemID) {
		quantity = itemModel.GetCapacity()
	}

	item, err := entity.NewItem(tx.ItemID, quantity, h.gs)
	if err != nil {
		return nil
	}

	_, err = character.AddItem(item, true) // allOrNothing = true for shop purchases
	if err != nil {
		character.Message("인벤토리 공간이 부족합니다.")
		return nil
	}

	character.RemoveMeso(int32(price))

	confirmPacket := &response.ConfirmShopTransaction{
		Code: 0,
	}
	character.Send(confirmPacket, types.SEND_POLICY_ENCRYPT)

	return nil
}

func (h *NpcShop) handleSell(character *entity.Character, shop *wz.Shop, tx *request.SellTransaction, resources *wz.Resources) error {
	quantity := tx.Quantity
	if quantity == 0xFFFF || quantity == 0 {
		quantity = 1
	}

	inventoryType := h.getItemInventoryType(tx.ItemID, nil)
	inventory := character.Inventory[inventoryType]
	if inventory == nil {
		return nil
	}

	item := inventory.Items[tx.Slot]
	if item == nil {
		return nil
	}

	if item.GetModel().GetID() != tx.ItemID {
		return nil
	}

	if h.isShuriken(tx.ItemID) || h.isBullet(tx.ItemID) {
		quantity = item.GetCount()
	}

	itemQuantity := item.GetCount()
	if itemQuantity == 0xFFFF {
		itemQuantity = 1
	}

	if quantity > itemQuantity || itemQuantity <= 0 {
		return nil
	}

	if h.isPet(tx.ItemID) {
		return nil
	}

	itemModel := item.GetModel()
	price := itemModel.GetPrice()

	if h.isShuriken(tx.ItemID) || h.isBullet(tx.ItemID) {
		wholePrice := itemModel.GetPrice()
		slotMax := itemModel.GetCapacity()
		if slotMax > 0 {
			price = int(math.Ceil(float64(wholePrice) / float64(slotMax)))
		}
	}

	recvMesos := int32(math.Max(math.Ceil(float64(price)*float64(quantity)), 0))
	if price == -1 || recvMesos <= 0 {
		return nil
	}

	if quantity >= itemQuantity {
		delete(inventory.Items, tx.Slot)
		character.Listener.OnRemoveInventorySlot(inventoryType, tx.Slot)
	} else {
		item.Reduce(quantity)
		character.Listener.OnUpdateInventorySlot(inventoryType, tx.Slot, item)
	}

	character.GainMeso(recvMesos)

	confirmPacket := &response.ConfirmShopTransaction{
		Code: 0x8,
	}
	character.Send(confirmPacket, types.SEND_POLICY_ENCRYPT)

	return nil
}

func (h *NpcShop) handleRecharge(character *entity.Character, shop *wz.Shop, tx *request.RechargeTransaction, resources *wz.Resources) error {
	inventory := character.Inventory[constant.INVENTORY_TYPE_CONSUME]
	if inventory == nil {
		return nil
	}

	item := inventory.Items[tx.Slot]
	if item == nil {
		return nil
	}

	itemID := item.GetModel().GetID()
	if !h.isShuriken(itemID) && !h.isBullet(itemID) {
		return nil
	}

	itemModel := item.GetModel()
	slotMax := itemModel.GetCapacity()

	if item.GetCount() >= slotMax {
		return nil
	}

	price := int(math.Round(float64(itemModel.GetPrice()) * float64(slotMax-item.GetCount())))

	if character.Meso < int32(price) {
		return nil
	}

	item.SetCount(slotMax)
	character.Listener.OnUpdateInventorySlot(constant.INVENTORY_TYPE_CONSUME, tx.Slot, item)

	character.RemoveMeso(int32(price))

	confirmPacket := &response.ConfirmShopTransaction{
		Code: 0x8,
	}
	character.Send(confirmPacket, types.SEND_POLICY_ENCRYPT)

	return nil
}

func (h *NpcShop) isShuriken(itemID uint32) bool {
	return itemID/10000 == 207
}

func (h *NpcShop) isBullet(itemID uint32) bool {
	return itemID/10000 == 233
}

func (h *NpcShop) isRechargable(itemID uint32) bool {
	return h.isShuriken(itemID) || h.isBullet(itemID)
}

func (h *NpcShop) isPet(itemID uint32) bool {
	return itemID/10000 == 500
}

func (h *NpcShop) getItemInventoryType(itemID uint32, itemModel wz.Item) constant.InventoryType {
	if itemModel != nil {
		switch itemModel.(type) {
		case *wz.Weapon, *wz.Armor:
			return constant.INVENTORY_TYPE_EQUIPMENT
		case *wz.Consume:
			return constant.INVENTORY_TYPE_CONSUME
		case *wz.Installation:
			return constant.INVENTORY_TYPE_INSTALLATION
		case *wz.GeneralItem:
			return constant.INVENTORY_TYPE_ETC
		case *wz.CashItem:
			return constant.INVENTORY_TYPE_CASH
		}
	}

	itemType := itemID / 10000
	switch {
	case itemType >= 100 && itemType < 200:
		return constant.INVENTORY_TYPE_EQUIPMENT
	case itemType >= 200 && itemType < 300:
		return constant.INVENTORY_TYPE_CONSUME
	case itemType >= 300 && itemType < 400:
		return constant.INVENTORY_TYPE_INSTALLATION
	case itemType >= 400 && itemType < 500:
		return constant.INVENTORY_TYPE_ETC
	case itemType >= 500 && itemType < 600:
		return constant.INVENTORY_TYPE_CASH
	default:
		return constant.INVENTORY_TYPE_ETC
	}
}
