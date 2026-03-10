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
	if req.AttackInfo.Skill != 0 {
		if !attackHandler.validateAndConsumeSkill(character, req.AttackInfo.Skill) {
			return nil
		}
		skillLevel = uint8(character.GetTotalSkillLevel(req.AttackInfo.Skill))
	}

	attackHandler.callOnAttackScript(ctx, character, mapInstance, req.AttackInfo.Damages, req.AttackInfo.Skill)
	attackHandler.applyDamageToMobs(character, mapInstance, req.AttackInfo.Damages)
	character.Listener.OnAttack(character, req.AttackInfo, skillLevel)

	return nil
}
