package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/game/client"
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/game/entity"
	"github.com/boyism80/fm/game/wz"
	"github.com/boyism80/fm/protocol/request"
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

	skillEntry := character.Skills.Get(skillID)
	isFirstPoint := false
	if skillEntry == nil {
		isFirstPoint = true
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

		skillEntry = entity.NewSkillEntry(character, wzSkill, 0, masterLevel)
	}

	maxLevel := skillEntry.MasterLevel
	if maxLevel == 0 && skillEntry.Wz.MaxLevel > 0 {
		maxLevel = skillEntry.Wz.MaxLevel
	}
	if maxLevel == 0 {
		log.Printf("Skill %d has no maxLevel for character %d", skillID, character.GetID())
		return nil
	}

	if skillEntry.MasterLevel > skillEntry.Wz.MaxLevel {
		log.Printf("Skill %d MasterLevel %d exceeds MaxLevel %d for character %d, clamping", skillID, skillEntry.MasterLevel, skillEntry.Wz.MaxLevel, character.GetID())
		skillEntry.MasterLevel = skillEntry.Wz.MaxLevel
		maxLevel = skillEntry.MasterLevel
	}

	if skillEntry.Level() >= maxLevel {
		log.Printf("Skill %d is already at max level %d for character %d", skillID, maxLevel, character.GetID())
		return nil
	}

	character.SkillPoint = character.SkillPoint - 1
	if isFirstPoint {
		skillEntry.SetLevel(1)
		character.Skills.Bind(skillID, skillEntry)
	} else {
		skillEntry.SetLevel(skillEntry.Level() + 1)
	}

	stats := map[constant.Stat]int32{
		constant.STAT_AVAILABLE_SP: int32(character.SkillPoint),
	}
	character.Listener.OnUpdateStats(character, stats, false)

	return nil
}
