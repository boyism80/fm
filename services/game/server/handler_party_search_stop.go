package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
)

type PartySearchStop struct {
	gs *GameServer
}

func (PartySearchStop) New(gs *GameServer) *PartySearchStop {
	return &PartySearchStop{
		gs: gs,
	}
}

func (h *PartySearchStop) Handle(ctx *core.ClientContext, req *request.PartySearchStop) error {
	client, ok := ctx.Client.(*client.GameClient)
	if !ok {
		log.Printf("Client is not a GameClient")
		return fmt.Errorf("client is not a GameClient")
	}

	ch := client.GetCharacter()
	if ch == nil {
		return nil
	}

	ch.SetPartySearchConfig(nil)
	return nil
}
