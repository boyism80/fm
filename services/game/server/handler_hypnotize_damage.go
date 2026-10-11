package server

import (
	"fmt"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
)

type HypnotizeDamage struct {
	gs *GameServer
}

func (HypnotizeDamage) New(gs *GameServer) *HypnotizeDamage {
	return &HypnotizeDamage{gs: gs}
}

func (h *HypnotizeDamage) Handle(ctx *core.ClientContext, req *request.HypnotizeDamage) error {
	client, ok := ctx.Client.(*client.GameClient)
	if !ok {
		return fmt.Errorf("client is not a GameClient")
	}
	character := client.GetCharacter()
	if character == nil {
		return fmt.Errorf("character is nil")
	}
	if req.CharacterID != character.GetID() {
		return nil
	}
	character.HypnotizeDamage(req.SourceOID, req.TargetOID, req.Damage)
	return nil
}
