package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
)

type UseSkillBook struct{}

func (UseSkillBook) New(_ *GameServer) *UseSkillBook {
	return &UseSkillBook{}
}

func (*UseSkillBook) Handle(ctx *core.ClientContext, req *request.UseSkillBook) error {
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

	ch.UseSkillBook(req.Slot, req.ItemID)
	return nil
}
