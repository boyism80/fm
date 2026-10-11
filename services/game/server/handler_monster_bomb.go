package server

import (
	"fmt"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
)

type MonsterBomb struct {
	gs *GameServer
}

func (MonsterBomb) New(gs *GameServer) *MonsterBomb {
	return &MonsterBomb{gs: gs}
}

func (h *MonsterBomb) Handle(ctx *core.ClientContext, req *request.MonsterBomb) error {
	client, ok := ctx.Client.(*client.GameClient)
	if !ok {
		return fmt.Errorf("client is not a GameClient")
	}
	character := client.GetCharacter()
	if character == nil {
		return fmt.Errorf("character is nil")
	}
	character.SelfDestructMob(req.OID)
	return nil
}
