package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game"
	"github.com/boyism80/fm/services/game/client"
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

	if character.Spectating() {
		return nil
	}

	if !character.IsAlive() {
		return nil
	}

	if character.GetMap() == nil {
		log.Printf("Character is not in a map")
		return fmt.Errorf("character is not in a map")
	}

	var skill *game.SkillEntry
	if req.Skill != 0 {
		skill = character.Skills.Get(req.Skill)
	}
	var skillLevel uint8
	attacked, err := character.Attack(req.ToAttackInfo(), false, func() bool {
		return skill != nil && character.Skills.CallHook(ctx.ActorContext, skill, "on_activating")
	}, func() {
		if req.Skill != 0 {
			skillLevel = uint8(character.Skills.Level(req.Skill))
		}
		CallOnAttackHooks(ctx.ActorContext, character, req.Damages, skill, false, true)
	})
	if err != nil {
		return err
	}
	if !attacked {
		character.Listener.OnUpdateStats(character, nil, true)
		return nil
	}

	character.Listener.OnRangedAttack(character, req, skillLevel)
	if skill != nil {
		character.Skills.CallHook(ctx.ActorContext, skill, "on_activated")
	}
	return nil
}
