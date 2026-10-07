package server

import (
	"fmt"

	"github.com/boyism80/fm/core"
	pconst "github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
)

type RingAction struct {
	gs *GameServer
}

func (RingAction) New(gs *GameServer) *RingAction {
	return &RingAction{
		gs: gs,
	}
}

func (h *RingAction) Handle(ctx *core.ClientContext, req *request.RingAction) error {
	client, ok := ctx.Client.(*client.GameClient)
	if ok == false {
		return fmt.Errorf("client is not a GameClient")
	}

	ch := client.GetCharacter()
	if ch == nil {
		return nil
	}

	switch req.Mode {
	case pconst.RingActionPropose:
		ch.Propose(req.Name, req.ItemID)
		ch.Listener.OnUnlockAction(ch)
	case pconst.RingActionCancelProposal:
		ch.CancelProposal()
	case pconst.RingActionAnswer:
		ch.AnswerProposal(ctx.ActorContext, req.Accepted, req.Name, req.CharacterID)
		ch.Listener.OnUnlockAction(ch)
	case pconst.RingActionDropRing:
		ch.DropMarriageItem(ctx.ActorContext, req.ItemID)
		ch.Listener.OnUnlockAction(ch)
	case pconst.RingActionInviteGuest:
		ch.InviteWeddingGuest(ctx.ActorContext, req.Name, req.MarriageID, int16(req.Slot))
		ch.Listener.OnUnlockAction(ch)
	case pconst.RingActionOpenInvitation:
		ch.OpenWeddingInvitation(ctx.ActorContext, int16(req.Slot), req.ItemID)
	case pconst.RingActionWeddingWishlist:
		ch.SubmitWeddingWishlist(ctx.ActorContext, req.Wishes)
	}
	return nil
}
