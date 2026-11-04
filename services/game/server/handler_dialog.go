package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
)

type Dialog struct {
	gs *GameServer
}

func (Dialog) New(gs *GameServer) *Dialog {
	return &Dialog{
		gs: gs,
	}
}

func (h *Dialog) Handle(ctx *core.ClientContext, req *request.Dialog) error {
	character, err := h.gs.character(ctx)
	if err != nil {
		return err
	}

	return character.Dialog().Resume(ctx.ActorContext, req.DialogType, req.Next, req.Selected, req.Text)
}
