package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
	"github.com/boyism80/fm/services/game/constant"
)

type DropMeso struct {
	gs *GameServer
}

func (DropMeso) New(gs *GameServer) *DropMeso {
	return &DropMeso{
		gs: gs,
	}
}

func (h *DropMeso) Handle(ctx *core.ClientContext, req *request.DropMeso) error {
	client, ok := ctx.Client.(*client.GameClient)
	if !ok {
		log.Printf("Client is not a GameClient")
		return fmt.Errorf("client is not a GameClient")
	}

	character := client.GetCharacter()
	if character == nil {
		log.Printf("Character is nil for client")
		return fmt.Errorf("character is nil")
	}

	if req.Count < 10 || req.Count > 50000 {
		log.Printf("Invalid meso count: %d", req.Count)
		character.Listener.OnUpdateStats(character, nil, true)
		return nil
	}

	if req.Count > character.Meso {
		log.Printf("Character doesn't have enough meso: %d < %d", character.Meso, req.Count)
		character.Listener.OnUpdateStats(character, nil, true)
		return nil
	}

	character.SetMeso(character.Meso - req.Count)
	character.Listener.OnUpdateStats(character, map[constant.Stat]int32{
		constant.StatMeso: character.Meso,
	}, true)

	mapInstance := character.GetMap()
	if mapInstance != nil {
		if _, err := mapInstance.SpawnMeso(req.Count, character.Position, character.GetID(), constant.DropTypeFFA); err != nil {
			log.Printf("Failed to spawn meso on map: %v", err)
		}
	}

	return nil
}
