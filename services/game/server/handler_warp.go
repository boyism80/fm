package server

import (
	"fmt"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
)

type Warp struct {
	gs *GameServer
}

func (Warp) New(gs *GameServer) *Warp {
	return &Warp{
		gs: gs,
	}
}

func (h *Warp) Handle(ctx *core.ClientContext, req *request.Warp) error {
	character, err := h.gs.character(ctx)
	if err != nil {
		return err
	}
	if req.FieldKey != character.FieldKey() {
		character.Listener().OnUpdateStats(character, nil, true)
		return nil
	}

	currentMap := character.GetMap()
	if currentMap == nil {
		return fmt.Errorf("current map not found")
	}

	if req.Target == 0xFFFFFFFF {
		return character.EnterPortalByName(ctx.ActorContext, req.PortalName)
	}

	return character.Revive(ctx.ActorContext)
}
