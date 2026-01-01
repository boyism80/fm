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
	gameServer *GameServer
	opcode     byte
}

func (DistributeSP) New(gameServer *GameServer) *DistributeSP {
	return &DistributeSP{
		gameServer: gameServer,
		opcode:     0x49,
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
	skillBookIndex := character.GetSkillBookIndexForSkill(skillID)

	if skillBookIndex < 0 || skillBookIndex >= len(character.SkillPoint) {
		return nil
	}

	remainingSP := character.SkillPoint[skillBookIndex]
	if remainingSP == 0 {
		return nil
	}

	if character.SkillsMap == nil {
		character.SkillsMap = make(map[uint32]*entity.SkillEntry)
	}

	skillEntry, exists := character.SkillsMap[skillID]
	if !exists {
		var wzSkill *wz.Skill
		if character.Context != nil {
			resources := character.Context.GetResources()
			if resources != nil {
				wzSkill = resources.GetSkill(skillID)
			}
		}

		if wzSkill == nil {
			return nil
		}

		masterLevel := 0
		if wzSkill.MasterLevel > 0 {
			masterLevel = wzSkill.MasterLevel
		} else if wzSkill.MaxLevel > 0 {
			masterLevel = wzSkill.MaxLevel
		} else {
			return nil
		}

		skillEntry = &entity.SkillEntry{
			Skill:       wzSkill,
			SkillLevel:  0,
			MasterLevel: masterLevel,
			Expiration:  time.Time{},
		}
		character.SkillsMap[skillID] = skillEntry
	}

	if skillEntry.Skill == nil {
		return nil
	}

	maxLevel := skillEntry.MasterLevel
	if maxLevel == 0 && skillEntry.Skill.MaxLevel > 0 {
		maxLevel = skillEntry.Skill.MaxLevel
	}
	if maxLevel == 0 {
		return nil
	}

	if skillEntry.MasterLevel > skillEntry.Skill.MaxLevel {
		skillEntry.MasterLevel = skillEntry.Skill.MaxLevel
		maxLevel = skillEntry.MasterLevel
	}

	if skillEntry.SkillLevel >= maxLevel {
		return nil
	}

	character.SkillPoint[skillBookIndex]--
	skillEntry.SkillLevel++

	if character.Listener != nil {
		stats := map[constant.Stat]int32{
			constant.STAT_AVAILABLE_SP: int32(character.SkillPoint[skillBookIndex]),
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
