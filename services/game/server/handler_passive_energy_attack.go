package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
)

type PassiveEnergyAttack struct {
	gs *GameServer
}

func (PassiveEnergyAttack) New(gs *GameServer) *PassiveEnergyAttack {
	return &PassiveEnergyAttack{gs: gs}
}

func (h *PassiveEnergyAttack) Handle(ctx *core.ClientContext, req *request.PassiveEnergyAttack) error {
	client, ok := ctx.Client.(*client.GameClient)
	if !ok {
		log.Printf("Client is not a GameClient")
		return fmt.Errorf("client is not a GameClient")
	}
	character := client.GetCharacter()
	if character == nil {
		return fmt.Errorf("character is nil")
	}

	attacked, err := character.EnergyAttack(req.ToAttackInfo(), func() {
		CallOnAttackHooks(ctx.ActorContext, character, req.Damages, nil, false, false)
	})
	if err != nil {
		return err
	}
	if !attacked {
		character.Listener.OnUpdateStats(character, nil, true)
		return nil
	}
	character.Listener.OnEnergyAttack(character, req)
	return nil
}
