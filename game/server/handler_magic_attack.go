package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/game/client"
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/game/entity"
	"github.com/boyism80/fm/game/wz"
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/types"
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

	mapID := character.GetMap()
	mapInstance := h.gs.GetMap(mapID)
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

	h.applyDamageToMobs(character, mapInstance, req.AttackInfo.Damages)

	magicAttackPacket := &response.MagicAttack{
		AttackInfo:  req.AttackInfo,
		CharacterId: character.GetID(),
		SkillLevel:  uint8(skillLevel),
	}

	mapInstance.Broadcast(character, magicAttackPacket, types.SEND_POLICY_ENCRYPT, character.GetID())

	return nil
}

func (h *MagicAttack) applyDamageToMobs(character *entity.Character, mapInstance *entity.Map, damages []dto.AttackPair) {
	for _, damage := range damages {
		mob := mapInstance.GetMob(damage.OID)
		if mob == nil {
			log.Printf("Mob not found for OID: %d", damage.OID)
			continue
		}

		if character.HasRoleAtLeast(constant.RoleAdmin) {
			mob.Damage(mob.Hp, character)
		} else {
			for _, damagePair := range damage.DamagePairs {
				mob.Damage(uint16(damagePair.Damage), character)
			}
		}
	}
}
