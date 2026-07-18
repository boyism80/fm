package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
	"github.com/boyism80/fm/services/game/constant"
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

	var targetMapId uint32
	var spawnPoint uint8
	stats := map[constant.Stat]int32{}

	if req.Target != 0xFFFFFFFF {
		if character.GetHp() > 0 {
			character.Listener.OnUpdateStats(character, nil, true)
			return nil
		}

		character.SetHp(50, false)
		character.Stance = constant.StanceDefaultValue

		targetMapId = uint32(wz.ReturnMapId)
		spawnPoint = 0
		stats[constant.StatHP] = int32(character.GetHp())

		character.Listener.OnUpdateStats(character, stats, true)
	} else {
		portal := currentMap.FindPortalByName(req.PortalName)
		if portal == nil || portal.Wz == nil {
			character.Listener.OnUpdateStats(character, nil, true)
			return nil
		}

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

		targetMapId = uint32(portal.Wz.TargetMapId)
		spawnPoint = targetPortal.Wz.ID

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
				MapActorPID:  currentMap.GetActorPID(),
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
		}
	}

	targetMap := h.gs.GetMapSystem().Get(targetMapId)
	if targetMap == nil {
		return fmt.Errorf("target map not found")
	}
	if err := character.Warp(targetMap, spawnPoint); err != nil {
		return err
	}

	return nil
}
