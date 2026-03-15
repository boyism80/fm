package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/game/client"
	"github.com/boyism80/fm/game/entity"
	"github.com/boyism80/fm/game/wz"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/protocol/response"
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

	if character.Hp <= 0 {
		return nil
	}

	mapInstance := character.GetMap()
	if mapInstance == nil {
		log.Printf("Character is not in a map")
		return fmt.Errorf("character is not in a map")
	}

	if req.AttackInfo.Skill == 0 {
		character.Listener.OnUpdateStats(nil, true)
		return nil
	}

	skillID := req.AttackInfo.Skill
	var wzSkill *wz.Skill
	if character.Context != nil {
		resources := character.Context.GetResources()
		if resources != nil {
			wzSkill = resources.GetSkill(skillID)
		}
	}

	if wzSkill == nil {
		log.Printf("Skill not found: %d", skillID)
		character.Listener.OnUpdateStats(nil, true)
		return nil
	}

	skillLevel := character.GetTotalSkillLevel(skillID)
	if skillLevel <= 0 {
		log.Printf("Character does not have skill %d or skill level is 0", skillID)
		character.Listener.OnUpdateStats(nil, true)
		return nil
	}

	levelData := wzSkill.GetLevelData(skillLevel)
	if levelData == nil {
		character.Listener.OnUpdateStats(nil, true)
		return nil
	}

	if levelData.Cooldown > 0 {
		skillEntry := character.Skills[skillID]
		if skillEntry == nil || skillEntry.IsCooling() {
			log.Printf("Skill %d is on cooldown", skillID)
			character.Listener.OnUpdateStats(nil, true)
			return nil
		}
		skillEntry.StartCooldown(levelData.Cooldown)
	}

	if levelData.MPCon > 0 {
		mpCon := uint16(levelData.MPCon)
		if !character.ConsumeMP(mpCon) {
			log.Printf("Not enough MP for skill %d (required: %d, current: %d)", skillID, mpCon, character.Mp)
			character.Listener.OnUpdateStats(nil, true)
			return nil
		}
	}
	if levelData.HPCon > 0 {
		hpCon := uint16(levelData.HPCon)
		if !character.ConsumeHP(hpCon) {
			log.Printf("Not enough HP for skill %d (required: %d, current: %d)", skillID, hpCon, character.Hp)
			character.Listener.OnUpdateStats(nil, true)
			return nil
		}
	}

	// Phase 1: per-skill on_activating (pre-attack hook).
	h.callSkillHook(ctx, character, uint32(skillID), "on_activating")

	damages := req.AttackInfo.Damages
	CallSkillOnAttack(ctx, character, uint32(skillID), mapInstance, damages)
	CallOnAttackScript(ctx, character, mapInstance, damages, uint32(skillID), false, 0)
	ApplyDamageToMobs(character, mapInstance, damages)

	magicAttackPacket := &response.MagicAttack{
		AttackInfo:  req.AttackInfo,
		CharacterId: character.GetID(),
		SkillLevel:  uint8(skillLevel),
	}

	mapInstance.Broadcast(magicAttackPacket, &entity.BroadcastOption{
		ExceptPlayerIDs:    []uint32{character.GetID()},
		ReferenceCharacter: character,
		RecipientFilter:    entity.BroadcastVisibleByReference,
	})

	// Phase 3: per-skill on_activated (finalization).
	h.callSkillHook(ctx, character, uint32(skillID), "on_activated")

	return nil
}

// callSkillHook delegates to Attack.callSkillHook so magic/ranged attacks share the same per-skill hook behavior.
func (h *MagicAttack) callSkillHook(ctx *core.ClientContext, character *entity.Character, skillID uint32, hook string) {
	attack := (&Attack{}).New(h.gs)
	attack.callSkillHook(ctx, character, skillID, hook)
}
