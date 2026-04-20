package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
	"github.com/boyism80/fm/services/game/wz"
)

type MagicAttack struct {
	gs     *GameServer
	opcode byte
}

func (MagicAttack) New(gs *GameServer) *MagicAttack {
	return &MagicAttack{
		gs:     gs,
		opcode: 0x1D,
	}
}

func (h *MagicAttack) GetOpcode() byte {
	return h.opcode
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
	var wzSkill *wz.Skill
	if character.GameWorld != nil {
		resources := character.GameWorld.GetResources()
		if resources != nil {
			wzSkill = resources.GetSkill(skillID)
		}
	}

	if wzSkill == nil {
		log.Printf("Skill not found: %d", skillID)
		character.Listener.OnUpdateStats(character, nil, true)
		return nil
	}

	skillLevel := character.GetTotalSkillLevel(skillID)
	if skillLevel <= 0 {
		log.Printf("Character does not have skill %d or skill level is 0", skillID)
		character.Listener.OnUpdateStats(character, nil, true)
		return nil
	}

	levelData := wzSkill.GetLevelData(skillLevel)
	if levelData == nil {
		character.Listener.OnUpdateStats(character, nil, true)
		return nil
	}

	if levelData.Cooldown > 0 {
		skillEntry := character.Skills.Get(skillID)
		if skillEntry == nil || skillEntry.IsCooling() {
			log.Printf("Skill %d is on cooldown", skillID)
			character.Listener.OnUpdateStats(character, nil, true)
			return nil
		}
		skillEntry.StartCooldown(levelData.Cooldown)
	}

	if !CallSkillHook(ctx, character, uint32(skillID), "on_activating") {
		character.Listener.OnUpdateStats(character, nil, true)
		return nil
	}

	damages := req.Damages
	CallOnAttackHooks(ctx, character, mapInstance, damages, uint32(skillID), false, 0)
	ApplyDamageToMobs(character, mapInstance, damages)

	character.Listener.OnMagicAttack(character, req, uint8(skillLevel))

	CallSkillHook(ctx, character, uint32(skillID), "on_activated")

	return nil
}
