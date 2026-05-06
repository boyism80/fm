package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
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
	client, ok := ctx.Client.(*client.GameClient)
	if !ok {
		log.Printf("Client is not a GameClient")
		return fmt.Errorf("client is not a GameClient")
	}

	character := client.GetCharacter()
	if character == nil {
		return fmt.Errorf("character not found")
	}

	currentMap := character.GetMap()
	if currentMap == nil {
		return fmt.Errorf("current map not found")
	}

	wz := currentMap.Wz
	if wz == nil {
		return fmt.Errorf("map model not found")
	}

	portal, ok := wz.FindPortal(req.PortalName)
	if !ok {
		character.Listener.OnUpdateStats(character, nil, true)
		return nil
	}

	targetMapWz, ok := h.gs.resources.Maps[uint32(portal.TargetMapId)]
	if !ok {
		character.Listener.OnUpdateStats(character, nil, true)
		return nil
	}

	targetPortal, ok := targetMapWz.FindPortal(portal.Target)
	if !ok {
		character.Listener.OnUpdateStats(character, nil, true)
		return nil
	}

	targetMap := h.gs.GetMap(uint32(portal.TargetMapId))
	if targetMap == nil {
		return fmt.Errorf("target map not found")
	}

	if err := character.Warp(targetMap, targetPortal.ID); err != nil {
		return err
	}
	return nil
}
