package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
	"github.com/boyism80/fm/services/game/entity"
)

type RangedAttack struct {
	gs *GameServer
}

func (RangedAttack) New(gs *GameServer) *RangedAttack {
	return &RangedAttack{
		gs: gs,
	}
}

func (h *RangedAttack) Handle(ctx *core.ClientContext, req *request.RangedAttack) error {
	gameClient, ok := ctx.Client.(*client.GameClient)
	if !ok {
		log.Printf("Client is not a GameClient")
		return fmt.Errorf("client is not a GameClient")
	}

	character := gameClient.GetCharacter()
	if character == nil {
		log.Printf("Character is nil for client")
		return fmt.Errorf("character is nil")
	}

	if character.IsAlive() == false {
		return nil
	}

	if character.GetMap() == nil {
		log.Printf("Character is not in a map")
		return fmt.Errorf("character is not in a map")
	}

	var skillLevel uint8 = 0
	var skill *entity.SkillEntry
	if req.Skill != 0 {
		skill = character.Skills.Get(req.Skill)
		activated := character.UseAttackSkill(req.Skill, func() bool {
			return skill != nil && character.CallSkillHook(ctx.ActorContext, skill, "on_activating")
		})
		if activated == false {
			character.Listener.OnUpdateStats(character, nil, true)
			return nil
		}
		skillLevel = uint8(character.GetTotalSkillLevel(req.Skill))
	}

	CallOnAttackHooks(ctx.ActorContext, character, req.Damages, skill, false, true, req.Slot)
	character.DamageTo(req.Damages)
	character.Listener.OnRangedAttack(character, req, skillLevel)
	if skill != nil {
		character.CallSkillHook(ctx.ActorContext, skill, "on_activated")
	}
	return nil
}
