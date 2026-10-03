package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
)

type CancelChair struct {
	gs *GameServer
}

func (CancelChair) New(gs *GameServer) *CancelChair {
	return &CancelChair{
		gs: gs,
	}
}

func (h *CancelChair) Handle(ctx *core.ClientContext, req *request.CancelChair) error {
	client, ok := ctx.Client.(*client.GameClient)
	if !ok {
		log.Printf("Client is not a GameClient")
		return fmt.Errorf("client is not a GameClient")
	}

	character := client.GetCharacter()
	if character == nil {
		return nil
	}

	if req.ChairID == -1 {
		character.StandUp()
	} else {
		character.SitOnMapSeat(req.ChairID)
	}

	character.Listener.OnUpdateStats(character, nil, true)

	return nil
}
