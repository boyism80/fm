package server

import (
	"fmt"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
)

type DirectWarp struct {
	gs *GameServer
}

func (DirectWarp) New(gs *GameServer) *DirectWarp {
	return &DirectWarp{
		gs: gs,
	}
}

func (h *DirectWarp) Handle(ctx *core.ClientContext, req *request.DirectWarp) error {
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

	if currentMap.Wz() == nil {
		return fmt.Errorf("map model not found")
	}

	return character.EnterPortalByName(ctx.ActorContext, req.PortalName)
}
