package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
)

type NpcControl struct {
	gs *GameServer
}

func (NpcControl) New(gs *GameServer) *NpcControl {
	return &NpcControl{
		gs: gs,
	}
}

func (h *NpcControl) Handle(ctx *core.ClientContext, req *request.NpcAction) error {
	character, err := h.gs.character(ctx)
	if err != nil {
		return err
	}

	character.Listener().OnNpcAction(character, req.Bytes)

	return nil
}
