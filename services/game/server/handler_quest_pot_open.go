package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
)

type QuestPotOpen struct{}

func (QuestPotOpen) New(_ *GameServer) *QuestPotOpen {
	return &QuestPotOpen{}
}

func (*QuestPotOpen) Handle(ctx *core.ClientContext, req *request.QuestPotOpen) error {
	gameClient, ok := ctx.Client.(*client.GameClient)
	if !ok {
		log.Printf("Client is not a GameClient")
		return fmt.Errorf("client is not a GameClient")
	}

	ch := gameClient.GetCharacter()
	if ch == nil {
		return nil
	}

	ch.Quests().OpenPot(uint32(req.QuestID))
	return nil
}
