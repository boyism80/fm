package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/clock"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
)

type HealOverTime struct{}

func (HealOverTime) New(_ *GameServer) *HealOverTime {
	return &HealOverTime{}
}

func (*HealOverTime) Handle(ctx *core.ClientContext, req *request.HealOverTime) error {
	client, ok := ctx.Client.(*client.GameClient)
	if !ok {
		log.Printf("Client is not a GameClient")
		return fmt.Errorf("client is not a GameClient")
	}

	character := client.GetCharacter()
	if character == nil {
		return nil
	}

	character.HealOverTime(req.HealHP, req.HealMP, req.PRate&1 == 1, clock.Now())
	return nil
}
