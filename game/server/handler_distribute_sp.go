package server

import (
	"fmt"
	"log"
	"time"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/game/client"
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/game/entity"
	"github.com/boyism80/fm/game/wz"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/types"
)

type DistributeSP struct {
	gs     *GameServer
	opcode byte
}

func (DistributeSP) New(gs *GameServer) *DistributeSP {
	return &DistributeSP{
		gs:     gs,
		opcode: 0x49,
	}
}

func (h *DistributeSP) GetOpcode() byte {
	return h.opcode
}

func (h *DistributeSP) Handle(ctx *core.ClientContext, req *request.DistributeSP) error {
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

	skillID := req.SkillID

	if character.SkillPoint == 0 {
		log.Printf("Character %d has no SP to distribute for skill %d", character.GetID(), skillID)
		return nil
	}

	skills := character.Skills

	skillEntry, exists := skills[skillID]
	if !exists {
		var wzSkill *wz.Skill
		if character.Context != nil {
			resources := character.Context.GetResources()
			if resources != nil {
				wzSkill = resources.GetSkill(skillID)
			}
		}

		if wzSkill == nil {
			log.Printf("Skill %d not found in WZ for character %d", skillID, character.GetID())
			return nil
		}

		masterLevel := 0
		if wzSkill.MasterLevel > 0 {
			masterLevel = wzSkill.MasterLevel
		} else if wzSkill.MaxLevel > 0 {
			masterLevel = wzSkill.MaxLevel
		} else {
			log.Printf("Skill %d has no masterLevel or maxLevel for character %d", skillID, character.GetID())
			return nil
		}

		skillEntry = &entity.SkillEntry{
			Skill:       wzSkill,
			SkillLevel:  0,
			MasterLevel: masterLevel,
			Expiration:  time.Time{},
			Owner:       character,
		}
		skills[skillID] = skillEntry
	}

	if skillEntry.Skill == nil {
		log.Printf("SkillEntry for skill %d has nil Skill for character %d", skillID, character.GetID())
		return nil
	}

	maxLevel := skillEntry.MasterLevel
	if maxLevel == 0 && skillEntry.Skill.MaxLevel > 0 {
		maxLevel = skillEntry.Skill.MaxLevel
	}
	if maxLevel == 0 {
		log.Printf("Skill %d has no maxLevel for character %d", skillID, character.GetID())
		return nil
	}

	if skillEntry.MasterLevel > skillEntry.Skill.MaxLevel {
		log.Printf("Skill %d MasterLevel %d exceeds MaxLevel %d for character %d, clamping", skillID, skillEntry.MasterLevel, skillEntry.Skill.MaxLevel, character.GetID())
		skillEntry.MasterLevel = skillEntry.Skill.MaxLevel
		maxLevel = skillEntry.MasterLevel
	}

	if skillEntry.SkillLevel >= maxLevel {
		log.Printf("Skill %d is already at max level %d for character %d", skillID, maxLevel, character.GetID())
		return nil
	}

	character.SkillPoint = character.SkillPoint - 1
	skillEntry.SkillLevel++

	if character.Listener != nil {
		stats := map[constant.Stat]int32{
			constant.STAT_AVAILABLE_SP: int32(character.SkillPoint),
		}
		character.Listener.OnUpdateStats(stats, false)
	}

	character.Send(&response.UpdateSkills{
		SkillID:     skillID,
		Level:       int32(skillEntry.SkillLevel),
		MasterLevel: int32(skillEntry.MasterLevel),
	}, types.SEND_POLICY_ENCRYPT)

	return nil
}
