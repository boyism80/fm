package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
)

type Carnival struct {
	gs *GameServer
}

func (Carnival) New(gs *GameServer) *Carnival {
	return &Carnival{gs: gs}
}

func (h *Carnival) Handle(ctx *core.ClientContext, req *request.Carnival) error {
	client, ok := ctx.Client.(*client.GameClient)
	if !ok {
		log.Printf("Client is not a GameClient")
		return fmt.Errorf("client is not a GameClient")
	}
	ch := client.GetCharacter()
	if ch == nil {
		return fmt.Errorf("character is nil")
	}
	match := ch.CarnivalMatch()
	if match == nil {
		ch.Listener.OnUpdateStats(ch, nil, true)
		return nil
	}
	if !match.Summon(ch, req.Tab, req.Num) {
		ch.Listener.OnUpdateStats(ch, nil, true)
	}
	return nil
}
