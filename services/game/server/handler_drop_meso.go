package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
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

	if !character.DropMeso(req.Count) {
		character.Listener.OnUpdateStats(character, nil, true)
	}
	return nil
}
