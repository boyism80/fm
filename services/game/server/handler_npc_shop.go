package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
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
	character, err := h.gs.character(ctx)
	if err != nil {
		return nil
	}
	if character.UsingCashItem() {
		character.Listener().OnUpdateStats(character, nil, true)
		return nil
	}

	shop := character.Dialog().Shop()
	if shop == nil {
		return nil
	}

	if req.Transaction == nil {
		character.Dialog().CloseShop()
		return nil
	}

	switch tx := req.Transaction.(type) {
	case *request.BuyTransaction:
		character.Dialog().Buy(shop, tx.ItemID, tx.Quantity)
	case *request.SellTransaction:
		character.Dialog().Sell(tx.Slot, tx.ItemID, tx.Quantity)
	case *request.RechargeTransaction:
		character.Dialog().Recharge(tx.Slot)
	default:
		character.Dialog().CloseShop()
	}
	return nil
}
