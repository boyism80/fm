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

	m := character.GetMap()
	if m == nil {
		character.Listener().OnUpdateStats(character, nil, true)
		return nil
	}
	door := m.FindDoorOwner(req.OwnerCharacterID)
	if door == nil || !door.Use(ctx.ActorContext, character) {
		character.Listener().OnUpdateStats(character, nil, true)
	}
	return nil
}
