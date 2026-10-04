package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
)

type EnhanceEquipment struct{}

func (EnhanceEquipment) New(_ *GameServer) *EnhanceEquipment {
	return &EnhanceEquipment{}
}

func (h *EnhanceEquipment) Handle(ctx *core.ClientContext, req *request.EnhanceEquipment) error {
	gameClient, ok := ctx.Client.(*client.GameClient)
	if !ok {
		log.Printf("Client is not a GameClient")
		return fmt.Errorf("client is not a GameClient")
	}

	ch := gameClient.GetCharacter()
	if ch == nil {
		return nil
	}

	if err := ch.Scroll(req.ScrollSlot, req.TargetSlot); err != nil {
		ch.Listener.OnUpdateStats(ch, nil, true)
	}
	return nil
}
