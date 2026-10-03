package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
	"github.com/boyism80/fm/services/game/constant"
)

type DistributeSP struct {
	gs *GameServer
}

func (DistributeSP) New(gs *GameServer) *DistributeSP {
	return &DistributeSP{
		gs: gs,
	}
}

func (h *DistributeSP) Handle(ctx *core.ClientContext, req *request.DistributeSP) error {
	client, ok := ctx.Client.(*client.GameClient)
	if !ok {
		log.Printf("Client is not a GameClient")
		return fmt.Errorf("client is not a GameClient")
	}

	character := client.GetCharacter()
	if character == nil {
		return nil
	}

	if character.DistributeSP(req.SkillID) == false {
		character.Listener.OnUpdateStats(character, map[constant.Stat]int32{}, true)
	}
	return nil
}
