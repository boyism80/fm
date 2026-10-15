package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
)

type MagicAttack struct {
	gs *GameServer
}

func (MagicAttack) New(gs *GameServer) *MagicAttack {
	return &MagicAttack{
		gs: gs,
	}
}

func (h *MagicAttack) Handle(ctx *core.ClientContext, req *request.MagicAttack) error {
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

	if character.Spectating() {
		return nil
	}

	if character.GetHp() <= 0 {
		return nil
	}

	if character.GetMap() == nil {
		log.Printf("Character is not in a map")
		return fmt.Errorf("character is not in a map")
	}

	if req.Skill == 0 {
		character.Listener.OnUpdateStats(character, nil, true)
		return nil
	}

	skill := character.Skills.Get(req.Skill)
	var skillLevel uint8
	attacked, err := character.Attack(req.ToAttackInfo(), true, func() bool {
		return skill != nil && character.Skills.CallHook(ctx.ActorContext, skill, "on_activating")
	}, func() {
		skillLevel = uint8(character.Skills.Level(req.Skill))
		CallOnAttackHooks(ctx.ActorContext, character, req.Damages, skill, true, false)
	})
	if err != nil {
		return err
	}
	if !attacked {
		character.Listener.OnUpdateStats(character, nil, true)
		return nil
	}

	character.Listener.OnMagicAttack(character, req, skillLevel)
	character.Skills.CallHook(ctx.ActorContext, skill, "on_activated")
	return nil
}
