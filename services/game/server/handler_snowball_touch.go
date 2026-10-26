package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
)

type SnowballTouch struct {
	gs *GameServer
}

func (SnowballTouch) New(gs *GameServer) *SnowballTouch {
	return &SnowballTouch{gs: gs}
}

// Running into a snowball knocks the character back.
func (h *SnowballTouch) Handle(ctx *core.ClientContext, req *request.SnowballTouch) error {
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
	ch.StateMachine().CallHook("on_snowball_touch", ch)
	return nil
}
