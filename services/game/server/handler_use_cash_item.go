package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game"
)

type UseCashItem struct {
	gs *GameServer
}

func (UseCashItem) New(gs *GameServer) *UseCashItem {
	return &UseCashItem{gs: gs}
}

func (h *UseCashItem) Handle(ctx *core.ClientContext, req *request.UseCashItem) error {
	ch, err := h.gs.character(ctx)
	if err != nil {
		return nil
	}

	model, ok := h.gs.GetResources().Items[req.ItemID]
	if ok && model.GetTraits()&constant.ItemTraitCashTeleportStone != 0 {
		ch.TeleportStones().Use(ctx.ActorContext, constant.InventoryTypeCash, int16(req.Slot), req.ItemID, req.Target.MapID, req.Target.Name)
		return nil
	}
	if ok && model.GetTraits()&constant.ItemTraitShopScanner != 0 {
		ch.RemoteShops().Scan(ctx.ActorContext, constant.InventoryTypeCash, int16(req.Slot), req.ItemID, req.SearchID, req.HighFirst)
		return nil
	}
	if req.Malformed {
		ch.Listener().OnUnlockAction(ch)
		return nil
	}
	ch.UseCashItem(ctx.ActorContext, int16(req.Slot), req.ItemID, game.CashItemInput{
		Text:         req.Text,
		Ear:          req.Ear,
		PetSN:        req.PetSN,
		Recipient:    req.Recipient,
		Lines:        req.Lines,
		Attached:     req.Attached,
		AttachedType: req.AttachedType,
		AttachedSlot: req.AttachedSlot,
		From:         req.From,
		To:           req.To,
		TargetType:   req.TargetType,
		TargetSlot:   req.TargetSlot,
		EquipSlot:    req.EquipSlot,
	})
	return nil
}
