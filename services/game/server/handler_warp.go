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
	character := h.gs.character(ctx)
	if character == nil {
		return nil
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
		portal := currentMap.FindPortalByName(req.PortalName)
		if portal == nil || portal.Wz() == nil || character.Occupied() {
			character.Listener().OnUpdateStats(character, nil, true)
			return nil
		}
		return character.EnterPortal(ctx.ActorContext, portal)
	}

	return character.Revive(ctx.ActorContext)
}
