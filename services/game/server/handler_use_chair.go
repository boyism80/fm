package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
)

type UseChair struct {
	gs *GameServer
}

func (UseChair) New(gs *GameServer) *UseChair {
	return &UseChair{
		gs: gs,
	}
}

func (h *UseChair) Handle(ctx *core.ClientContext, req *request.UseChair) error {
	client, ok := ctx.Client.(*client.GameClient)
	if !ok {
		log.Printf("Client is not a GameClient")
		return fmt.Errorf("client is not a GameClient")
	}

	character := client.GetCharacter()
	if character == nil {
		return nil
	}

	if character.GetMap() == nil {
		return nil
	}

	if err := character.SitOnChair(req.ItemID); err != nil {
		log.Printf("Chair %d: %v", req.ItemID, err)
	}
	character.Listener.OnUpdateStats(character, nil, true)

	return nil
}
