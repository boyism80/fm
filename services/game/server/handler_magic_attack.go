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
	activated := character.UseAttackSkill(req.Skill, func() bool {
		return skill != nil && character.CallSkillHook(ctx.ActorContext, skill, "on_activating")
	})
	if !activated {
		character.Listener.OnUpdateStats(character, nil, true)
		return nil
	}
	skillLevel := uint8(character.GetTotalSkillLevel(req.Skill))

	CallOnAttackHooks(ctx.ActorContext, character, req.Damages, skill, true, false, 0)
	character.DamageTo(req.Damages)
	character.Listener.OnMagicAttack(character, req, skillLevel)
	character.CallSkillHook(ctx.ActorContext, skill, "on_activated")
	return nil
}
