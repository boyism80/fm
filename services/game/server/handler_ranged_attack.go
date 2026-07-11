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

	mapInstance := character.GetMap()
	if mapInstance == nil {
		log.Printf("Character is not in a map")
		return fmt.Errorf("character is not in a map")
	}

	attackHandler := (&Attack{}).New(h.gs)
	var skillLevel uint8 = 0
	skillID := req.Skill
	if skillID != 0 {
		if !attackHandler.validateSkillForAttack(character, skillID) {
			return nil
		}
		skillLevel = uint8(character.GetTotalSkillLevel(skillID))
		if !CallSkillHook(character, skillID, "on_activating") {
			character.Listener.OnUpdateStats(character, nil, true)
			return nil
		}
	}
	return h.finishRangedAttack(ctx, character, mapInstance, req, skillLevel, skillID)
}

func (h *RangedAttack) finishRangedAttack(ctx *core.ClientContext, character *entity.Character, mapInstance *entity.Map, req *request.RangedAttack, skillLevel uint8, skillID uint32) error {
	damages := req.Damages
	CallOnAttackHooks(character, damages, skillID, false, true, req.Slot)
	character.DamageTo(damages)
	character.Listener.OnRangedAttack(character, req, skillLevel)
	if skillID != 0 {
		CallSkillHook(character, skillID, "on_activated")
	}
	return nil
}
