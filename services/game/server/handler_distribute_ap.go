package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game"
	"github.com/boyism80/fm/services/game/client"
)

type DistributeAP struct {
	gs *GameServer
}

func (DistributeAP) New(gs *GameServer) *DistributeAP {
	return &DistributeAP{
		gs: gs,
	}
}

func (h *DistributeAP) Handle(ctx *core.ClientContext, req *request.DistributeAP) error {
	client, ok := ctx.Client.(*client.GameClient)
	if !ok {
		log.Printf("Client is not a GameClient")
		return fmt.Errorf("client is not a GameClient")
	}

	character := client.GetCharacter()
	if character == nil {
		log.Printf("Character is nil for client")
		return fmt.Errorf("character is nil")
	}

	stat := constant.StatType(req.StatType)
	assigned := false
	switch stat {
	case constant.StatTypeHP, constant.StatTypeMP:
		assigned = character.Points.AssignHPMP(stat)
	default:
		assigned = character.Points.AssignStats([]game.APEntry{{Stat: stat, Amount: 1}})
	}
	if !assigned {
		character.Listener.OnUpdateStats(character, map[constant.Stat]int32{}, true)
	}
	return nil
}
