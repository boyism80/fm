package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/game/client"
	"github.com/boyism80/fm/protocol/request"
)

type RangedAttack struct {
	gs     *GameServer
	opcode byte
}

func (RangedAttack) New(gs *GameServer) *RangedAttack {
	return &RangedAttack{
		gs:     gs,
		opcode: 0x1C,
	}
}

func (h *RangedAttack) GetOpcode() byte {
	return h.opcode
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
	skillID := req.AttackInfo.Skill
	if skillID != 0 {
		if !attackHandler.validateAndConsumeSkill(character, skillID) {
			return nil
		}
		skillLevel = uint8(character.GetTotalSkillLevel(skillID))
		// Phase 1: per-skill on_activating (pre-attack hook).
		attackHandler.callSkillHook(ctx, character, skillID, "on_activating")
	}

	attackHandler.callOnAttackScript(ctx, character, mapInstance, req.AttackInfo.Damages, skillID, true, req.AttackInfo.Slot)
	attackHandler.applyDamageToMobs(character, mapInstance, req.AttackInfo.Damages)

	// Phase 2: per-skill on_attack (post-damage).
	if skillID != 0 {
		attackHandler.callSkillHook(ctx, character, skillID, "on_attack")
	}

	character.Listener.OnAttack(character, req.AttackInfo, skillLevel)

	// Phase 3: per-skill on_activated (finalization).
	if skillID != 0 {
		attackHandler.callSkillHook(ctx, character, skillID, "on_activated")
	}

	return nil
}
