package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
)

type QuestPotFeed struct{}

func (QuestPotFeed) New(_ *GameServer) *QuestPotFeed {
	return &QuestPotFeed{}
}

func (*QuestPotFeed) Handle(ctx *core.ClientContext, req *request.QuestPotFeed) error {
	gameClient, ok := ctx.Client.(*client.GameClient)
	if !ok {
		log.Printf("Client is not a GameClient")
		return fmt.Errorf("client is not a GameClient")
	}

	ch := gameClient.GetCharacter()
	if ch == nil {
		return nil
	}
	defer ch.Listener().OnUnlockAction(ch)

	ch.Quests().FeedPot(req.Slot, req.ItemID, req.QuestID)
	return nil
}
