package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
)

type CoconutHit struct {
	gs *GameServer
}

func (CoconutHit) New(gs *GameServer) *CoconutHit {
	return &CoconutHit{gs: gs}
}

// The channel event decides what the struck coconut does.
func (h *CoconutHit) Handle(ctx *core.ClientContext, req *request.CoconutHit) error {
	client, ok := ctx.Client.(*client.GameClient)
	if !ok {
		log.Printf("Client is not a GameClient")
		return fmt.Errorf("client is not a GameClient")
	}
	ch := client.GetCharacter()
	if ch == nil {
		return fmt.Errorf("character is nil")
	}

	if ch.StateMachine() == nil {
		return nil
	}
	ch.StateMachine().CallHook("on_coconut_hit", ch, int(req.ID))
	return nil
}
