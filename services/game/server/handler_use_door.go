package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
	"github.com/boyism80/fm/services/game/entity"
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

	door := mapInstance.FindDoorOwner(req.OwnerCharacterID)
	if door == nil || door.Map != mapInstance {
		character.Listener.OnUpdateStats(character, nil, true)
		return nil
	}
	if !door.UsableBy(character) {
		character.Listener.OnUpdateStats(character, nil, true)
		return nil
	}

	var targetMap *entity.Map
	var spawnID uint8
	switch mapInstance {
	case door.Field.Map:
		targetMap = door.Return.Map
		spawnID = door.Return.PortalID
	case door.Return.Map:
		targetMap = door.Field.Map
		spawnID = door.Field.PortalID
	}
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
