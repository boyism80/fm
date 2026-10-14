package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/domain"
)

type UseCashItem struct {
	gs *GameServer
}

func (UseCashItem) New(gs *GameServer) *UseCashItem {
	return &UseCashItem{gs: gs}
}

func (h *UseCashItem) Handle(ctx *core.ClientContext, req *request.UseCashItem) error {
	client, ok := ctx.Client.(*client.GameClient)
	if !ok {
		log.Printf("Client is not a GameClient")
		return fmt.Errorf("client is not a GameClient")
	}

	ch := client.GetCharacter()
	if ch == nil {
		return nil
	}

	model, ok := h.gs.GetResources().Items[req.ItemID]
	if ok && model.GetTraits()&constant.ItemTraitCashTeleportStone != 0 {
		ch.TeleportStones.Use(ctx.ActorContext, constant.InventoryTypeCash, int16(req.Slot), req.ItemID, req.Target.MapID, req.Target.Name)
		return nil
	}
	if ok && model.GetTraits()&constant.ItemTraitShopScanner != 0 {
		ch.UseShopScanner(ctx.ActorContext, constant.InventoryTypeCash, int16(req.Slot), req.ItemID, req.SearchID, req.HighFirst)
		return nil
	}
	if req.Malformed {
		ch.Listener.OnUnlockAction(ch)
		return nil
	}
	ch.UseCashItem(ctx.ActorContext, int16(req.Slot), req.ItemID, domain.CashItemInput{
		Text:         req.Text,
		Ear:          req.Ear,
		PetSN:        req.PetSN,
		Recipient:    req.Recipient,
		Lines:        req.Lines,
		Attached:     req.Attached,
		AttachedType: req.AttachedType,
		AttachedSlot: req.AttachedSlot,
	})
	return nil
}
