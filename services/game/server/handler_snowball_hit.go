package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
)

type SnowballHit struct {
	gs *GameServer
}

func (SnowballHit) New(gs *GameServer) *SnowballHit {
	return &SnowballHit{gs: gs}
}

// The channel event judges the hit by where the character stands, not by the damage the client reports.
func (h *SnowballHit) Handle(ctx *core.ClientContext, req *request.SnowballHit) error {
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
	ch.StateMachine().CallHook("on_snowball_hit", ch)
	return nil
}
