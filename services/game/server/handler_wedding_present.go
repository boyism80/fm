package server

import (
	"fmt"

	"github.com/boyism80/fm/core"
	pconst "github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
	"github.com/boyism80/fm/services/game/constant"
)

type WeddingPresent struct {
	gs *GameServer
}

func (WeddingPresent) New(gs *GameServer) *WeddingPresent {
	return &WeddingPresent{
		gs: gs,
	}
}

func (h *WeddingPresent) Handle(ctx *core.ClientContext, req *request.WeddingPresent) error {
	client, ok := ctx.Client.(*client.GameClient)
	if ok == false {
		return fmt.Errorf("client is not a GameClient")
	}

	ch := client.GetCharacter()
	if ch == nil {
		return nil
	}

	var err error
	switch req.Mode {
	case pconst.WeddingPresentGive:
		err = ch.Wedding.GiveGift(ctx.ActorContext, req.Slot, req.ItemID, req.Count)
	case pconst.WeddingPresentReceive:
		err = ch.Wedding.ReceiveGift(ctx.ActorContext, constant.InventoryType(req.InventoryType), int(req.Index))
	case pconst.WeddingPresentClose:
		ch.Wedding.CloseGift()
	}
	if err != nil {
		ch.Listener.OnUnlockAction(ch)
	}
	return nil
}
