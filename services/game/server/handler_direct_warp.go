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

	if currentMap.Wz() == nil {
		return fmt.Errorf("map model not found")
	}

	portal := currentMap.FindPortalByName(req.PortalName)
	if portal == nil || portal.Wz() == nil || character.Occupied() {
		character.Listener().OnUpdateStats(character, nil, true)
		return nil
	}
	return character.EnterPortal(ctx.ActorContext, portal)
}
