package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
)

type MonsterCarnival struct {
	gs *GameServer
}

func (MonsterCarnival) New(gs *GameServer) *MonsterCarnival {
	return &MonsterCarnival{gs: gs}
}

func (h *MonsterCarnival) Handle(ctx *core.ClientContext, req *request.MonsterCarnival) error {
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
	if !match.HandleTab(ch, req.Tab, req.Num) {
		ch.Listener.OnUpdateStats(ch, nil, true)
	}
	return nil
}
