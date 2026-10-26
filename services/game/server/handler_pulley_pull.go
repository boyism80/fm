package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
)

type PulleyPull struct {
	gs *GameServer
}

func (PulleyPull) New(gs *GameServer) *PulleyPull {
	return &PulleyPull{gs: gs}
}

func (h *PulleyPull) Handle(ctx *core.ClientContext, req *request.PulleyPull) error {
	client, ok := ctx.Client.(*client.GameClient)
	if !ok {
		log.Printf("Client is not a GameClient")
		return fmt.Errorf("client is not a GameClient")
	}
	ch := client.GetCharacter()
	if ch == nil {
		return fmt.Errorf("character is nil")
	}

	m := ch.GetMap()
	if m == nil {
		return nil
	}
	m.PullPulley(ch)
	return nil
}
