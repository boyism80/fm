package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
)

type DropMeso struct {
	gs *GameServer
}

func (DropMeso) New(gs *GameServer) *DropMeso {
	return &DropMeso{
		gs: gs,
	}
}

func (h *DropMeso) Handle(ctx *core.ClientContext, req *request.DropMeso) error {
	character, err := h.gs.character(ctx)
	if err != nil {
		return err
	}

	if !character.DropMeso(req.Count) {
		character.Listener().OnUpdateStats(character, nil, true)
	}
	return nil
}
