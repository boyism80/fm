package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
	"github.com/boyism80/fm/services/game/entity"
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

	if character.GetHp() <= 0 {
		return nil
	}

	mapInstance := character.GetMap()
	if mapInstance == nil {
		log.Printf("Character is not in a map")
		return fmt.Errorf("character is not in a map")
	}

	if req.Skill == 0 {
		character.Listener.OnUpdateStats(character, nil, true)
		return nil
	}

	skillID := req.Skill
	activated := character.UseAttackSkill(skillID, func() bool {
		return CallSkillHook(character, skillID, "on_activating")
	})
	if activated == false {
		character.Listener.OnUpdateStats(character, nil, true)
		return nil
	}
	skillLevel := character.GetTotalSkillLevel(skillID)
	return h.finishMagicAttack(ctx, character, mapInstance, req, uint8(skillLevel), skillID)
}

func (h *MagicAttack) finishMagicAttack(ctx *core.ClientContext, character *entity.Character, mapInstance *entity.Map, req *request.MagicAttack, skillLevel uint8, skillID uint32) error {
	damages := req.Damages
	CallOnAttackHooks(character, damages, skillID, true, false, 0)
	character.DamageTo(damages)
	character.Listener.OnMagicAttack(character, req, skillLevel)
	CallSkillHook(character, skillID, "on_activated")
	return nil
}
