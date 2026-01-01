package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/game/client"
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/game/entity"
	"github.com/boyism80/fm/protocol/request"
)

// Warp handles warp packet requests
type Warp struct {
	gameServer *GameServer
	opcode     byte
}

func (Warp) New(gameServer *GameServer) *Warp {
	return &Warp{
		gameServer: gameServer,
		opcode:     0x15,
	}
}

func (h *Warp) GetOpcode() byte {
	return h.opcode
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
		if character.Hp == 0 {
			character.Hp = 50
			character.Stance = 0

			currentMap := h.gameServer.GetMap(character.Map)
			if currentMap == nil {
				return fmt.Errorf("current map not found")
			}

			mapSpec := currentMap.GetSpec()
			if mapSpec == nil {
				return fmt.Errorf("map model not found")
			}

			targetMapId = uint32(mapSpec.ReturnMapId)
			spawnPoint = 0
			stats[constant.STAT_HP] = int32(character.Hp)

			character.Listener.OnUpdateStats(stats, true)
		} else {
			character.Listener.OnUpdateStats(nil, true)
			return nil
		}
	} else {
		currentMap := h.gameServer.GetMap(character.Map)
		if currentMap == nil {
			return fmt.Errorf("current map not found")
		}

		mapSpec := currentMap.GetSpec()
		if mapSpec == nil {
			return fmt.Errorf("map model not found")
		}

		portal, ok := mapSpec.FindPortal(req.PortalName)
		if !ok {
			character.Listener.OnUpdateStats(nil, true)
			return nil
		}

		targetMapSpec, ok := h.gameServer.resources.Maps[uint32(portal.TargetMapId)]
		if !ok {
			character.Listener.OnUpdateStats(nil, true)
			return nil
		}

		targetPortal, ok := targetMapSpec.FindPortal(portal.Target)
		if !ok {
			character.Listener.OnUpdateStats(nil, true)
			return nil
		}

		targetMapId = uint32(portal.TargetMapId)
		spawnPoint = targetPortal.ID
	}

	// Perform the warp
	if err := h.performWarp(client, character, targetMapId, spawnPoint); err != nil {
		return fmt.Errorf("failed to perform warp: %v", err)
	}

	return nil
}

func (h *Warp) performWarp(client *client.GameClient, character *entity.Character, targetMapId uint32, spawnPoint uint8) error {
	currentMap := h.gameServer.GetMap(character.Map)
	if currentMap == nil {
		return fmt.Errorf("current map not found")
	}

	targetMap := h.gameServer.GetMap(targetMapId)
	if targetMap == nil {
		return fmt.Errorf("target map %d not found", targetMapId)
	}

	currentMap.RemovePlayer(character.GetID())

	character.Map = targetMapId
	character.SpawnPoint = spawnPoint

	mapSpec := targetMap.GetSpec()
	if mapSpec != nil && len(mapSpec.Portals) > 0 {
		for portalId, portal := range mapSpec.Portals {
			if portalId == spawnPoint {
				character.Position = portal.Position
				break
			}
		}
	}

	task := &core.LogicTask{
		Predicate: func() bool {
			return character != nil && client.GetConnection() != nil
		},
		Logic: func() error {
			if err := targetMap.AddPlayer(character.GetID(), character, false); err != nil {
				return fmt.Errorf("failed to add character to target map: %v", err)
			}
			return nil
		},
		Callback: func(success bool, err error) {
			if err != nil {
				log.Printf("Failed to add character to target map: %v", err)
			}
		},
		Object:     client,
		MaxRetries: 3,
	}

	if err := h.gameServer.server.SubmitLogicTaskForObject(client, task); err != nil {
		return fmt.Errorf("failed to submit warp task: %v", err)
	}

	return nil
}
