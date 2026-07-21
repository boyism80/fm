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

	if currentMap.Wz == nil {
		return fmt.Errorf("map model not found")
	}

	portal := currentMap.FindPortalByName(req.PortalName)
	if portal == nil || portal.Wz == nil {
		character.Listener.OnUpdateStats(character, nil, true)
		return nil
	}
	scriptName := portal.Script()
	if scriptName != "" {
		root := currentMap.GetLuaRoot()
		if root == nil {
			character.Listener.OnUnlockAction(character)
			return fmt.Errorf("root lua state not found")
		}
		scriptPath := fmt.Sprintf("script/portal/%s.lua", scriptName)
		thread, err := luax.NewThread(root, scriptPath)
		if err != nil {
			character.Listener.OnUnlockAction(character)
			return fmt.Errorf("portal script thread: %w", err)
		}
		luax.SetConfiguration(thread, luax.Configuration{
			ActorContext: ctx.ActorContext,
			ActorPID:     currentMap.LogicActorPID(),
		})
		luax.CallAsync(root, thread, "on_enter", character).Then(func(_ interface{}) (interface{}, error) {
			if character.GetDialog() == nil {
				character.Listener.OnUnlockAction(character)
			}
			return nil, nil
		}).OnError(func(err error) {
			log.Printf("portal script %s failed: %v", scriptPath, err)
			if character.GetDialog() == nil {
				character.Listener.OnUnlockAction(character)
			}
		})
		return nil
	} else {
		targetMap := h.gs.GetMapSystem().Get(uint32(portal.Wz.TargetMapId))
		if targetMap == nil {
			character.Listener.OnUpdateStats(character, nil, true)
			return nil
		}

		targetPortal := targetMap.FindPortalByName(portal.Wz.Target)
		if targetPortal == nil || targetPortal.Wz == nil {
			character.Listener.OnUpdateStats(character, nil, true)
			return nil
		}

		if err := character.Warp(ctx.ActorContext, targetMap, targetPortal.Wz.ID); err != nil {
			return err
		}
	}
	return nil
}
