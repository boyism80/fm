package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
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

	shopID := character.Dialog.Shop()
	if shopID == 0 {
		return nil
	}

	shop := resources.GetShop(shopID)
	if shop == nil {
		return nil
	}

	if req.Transaction == nil {
		character.Dialog.CloseShop()
		return nil
	}

	switch tx := req.Transaction.(type) {
	case *request.BuyTransaction:
		character.Dialog.Buy(shop, tx.ItemID, tx.Quantity)
	case *request.SellTransaction:
		character.Dialog.Sell(tx.Slot, tx.ItemID, tx.Quantity)
	case *request.RechargeTransaction:
		character.Dialog.Recharge(tx.Slot)
	default:
		character.Dialog.CloseShop()
	}
	return nil
}
