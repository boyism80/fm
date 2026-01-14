package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/game/client"
	"github.com/boyism80/fm/game/entity"
	"github.com/boyism80/fm/game/wz"
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/protocol/request"
)

type Attack struct {
	gs     *GameServer
	opcode byte
}

func (Attack) New(gs *GameServer) *Attack {
	return &Attack{
		gs:     gs,
		opcode: 0x1B,
	}
}

func (h *Attack) GetOpcode() byte {
	return h.opcode
}

func (h *Attack) Handle(ctx *core.ClientContext, req *request.Attack) error {
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

	mapID := character.GetMap()
	mapInstance := h.gs.GetMap(mapID)
	if mapInstance == nil {
		log.Printf("Character is not in a map")
		return fmt.Errorf("character is not in a map")
	}

	if req.AttackInfo.Skill != 0 {
		if !h.validateAndConsumeSkill(character, req.AttackInfo.Skill) {
			return nil
		}
	}

	h.applyDamageToMobs(character, mapInstance, req.AttackInfo.Damages)
	character.Listener.OnAttack(mapID, character.GetID(), req.AttackInfo, 0)

	return nil
}

func (h *Attack) validateAndConsumeSkill(character *entity.Character, skillID uint32) bool {
	var wzSkill *wz.Skill
	if character.Context != nil {
		resources := character.Context.GetResources()
		if resources != nil {
			wzSkill = resources.GetSkill(skillID)
		}
	}

	if wzSkill == nil {
		log.Printf("Skill not found: %d", skillID)
		return false
	}

	skillLevel := character.GetTotalSkillLevel(skillID)
	if skillLevel <= 0 {
		log.Printf("Character does not have skill %d or skill level is 0", skillID)
		return false
	}

	levelData := wzSkill.GetLevelData(skillLevel)
	if levelData == nil {
		return false
	}

	if levelData.Cooldown > 0 {
		if character.IsSkillCooling(skillID) {
			log.Printf("Skill %d is on cooldown", skillID)
			return false
		}
		character.AddCooldown(skillID, levelData.Cooldown)
	}

	if levelData.MPCon > 0 {
		mpCon := uint16(levelData.MPCon)
		if !character.ConsumeMP(mpCon) {
			log.Printf("Not enough MP for skill %d (required: %d, current: %d)", skillID, mpCon, character.Mp)
			return false
		}
	}

	return true
}

func (h *Attack) applyDamageToMobs(character *entity.Character, mapInstance *entity.Map, damages []dto.AttackPair) {
	for _, damage := range damages {
		mob := mapInstance.GetMob(damage.OID)
		if mob == nil {
			log.Printf("Mob not found for OID: %d", damage.OID)
			continue
		}

		if character.Admin {
			mob.Damage(mob.Hp, character)
		} else {
			for _, damagePair := range damage.DamagePairs {
				mob.Damage(uint16(damagePair.Damage), character)
			}
		}
	}
}
