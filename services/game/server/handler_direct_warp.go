package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/luax"
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
	if portal.ScriptName != "" {
		root := currentMap.GetLuaRoot()
		if root == nil {
			return fmt.Errorf("root lua state not found")
		}
		scriptPath := fmt.Sprintf("script/portal/%s.lua", portal.ScriptName)
		if _, err := luax.InlineCall(root, scriptPath, "on_enter", character); err != nil {
			return err
		}
		character.Listener.OnUpdateStats(character, nil, true)
		return nil
	} else {
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

		targetMap := h.gs.GetMapSystem().Get(uint32(portal.TargetMapId))
		if targetMap == nil {
			return fmt.Errorf("target map not found")
		}

		if err := character.Warp(targetMap, targetPortal.ID); err != nil {
			return err
		}
	}
	return nil
}
