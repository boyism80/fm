package server

import (
	"github.com/boyism80/fm/core"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game"
)

type Messenger struct {
	gs *GameServer
}

func (Messenger) New(gs *GameServer) *Messenger {
	return &Messenger{
		gs: gs,
	}
}

func (h *Messenger) Handle(ctx *core.ClientContext, req *request.Messenger) error {
	if ctx.ActorContext == nil {
		return nil
	}

	ch := h.gs.character(ctx)
	if ch == nil {
		return nil
	}

	switch req.Action {
	case internal.MessengerAction_MESSENGER_ACTION_ENTER, internal.MessengerAction_MESSENGER_ACTION_LEAVE:
	case internal.MessengerAction_MESSENGER_ACTION_INVITE, internal.MessengerAction_MESSENGER_ACTION_DECLINE:
		if !game.ClientText(req.Text).Valid(13) {
			return nil
		}
	case internal.MessengerAction_MESSENGER_ACTION_CHAT:
		if !game.ClientText(req.Text).Valid(200) {
			return nil
		}
	default:
		return nil
	}
	ch.Messenger().RequestAsync(ctx.ActorContext, req.Action, req.RoomID, req.Text)
	return nil
}
