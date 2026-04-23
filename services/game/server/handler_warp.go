package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
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

		currentMap := character.GetMap()
		if currentMap == nil {
			return fmt.Errorf("current map not found")
		}

		wz := currentMap.Wz
		if wz == nil {
			return fmt.Errorf("map model not found")
		}

		targetMapId = uint32(wz.ReturnMapId)
		spawnPoint = 0
		stats[constant.STAT_HP] = int32(character.GetHp())

		character.Listener.OnUpdateStats(character, stats, true)
	} else {
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

		targetMapId = uint32(portal.TargetMapId)
		spawnPoint = targetPortal.ID
	}

	targetMap := h.gs.GetMap(targetMapId)
	if targetMap == nil {
		return fmt.Errorf("target map not found")
	}
	if err := character.Warp(targetMap, spawnPoint); err != nil {
		return err
	}
	return nil
}
