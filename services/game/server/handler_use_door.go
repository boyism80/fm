package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
)

type UseDoor struct {
	gs *GameServer
}

func (UseDoor) New(gs *GameServer) *UseDoor {
	return &UseDoor{
		gs: gs,
	}
}

func (h *UseDoor) Handle(ctx *core.ClientContext, req *request.UseDoor) error {
	client, ok := ctx.Client.(*client.GameClient)
	if !ok {
		log.Printf("Client is not a GameClient")
		return fmt.Errorf("client is not a GameClient")
	}

	character := client.GetCharacter()
	if character == nil {
		return fmt.Errorf("character not found")
	}

	mapInstance := character.GetMap()
	if mapInstance == nil || mapInstance.Wz == nil {
		character.Listener.OnUpdateStats(character, nil, true)
		return nil
	}

	door := mapInstance.FindDoorByOwner(req.OwnerCharacterID)
	if door == nil || door.Map != mapInstance {
		character.Listener.OnUpdateStats(character, nil, true)
		return nil
	}
	if !door.UsableBy(character) {
		character.Listener.OnUpdateStats(character, nil, true)
		return nil
	}

	if door.Return.MapID == 0 || door.Field.MapID == 0 {
		character.Listener.OnUpdateStats(character, nil, true)
		return nil
	}

	fieldWZID := mapInstance.Wz.ID
	var targetMapID uint32
	var spawnID uint8
	switch fieldWZID {
	case door.Field.MapID:
		targetMapID = door.Return.MapID
		spawnID = door.Return.PortalID
	case door.Return.MapID:
		targetMapID = door.Field.MapID
		spawnID = door.Field.PortalID
	default:
		character.Listener.OnUpdateStats(character, nil, true)
		return nil
	}

	targetMap := h.gs.GetMapSystem().Get(targetMapID)
	if targetMap == nil || targetMap.Wz == nil {
		character.Listener.OnUpdateStats(character, nil, true)
		return nil
	}

	if err := character.Warp(ctx.ActorContext, targetMap, spawnID); err != nil {
		log.Printf("use door warp: %v", err)
		character.Listener.OnUpdateStats(character, nil, true)
		return nil
	}
	return nil
}
