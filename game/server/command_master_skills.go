package server

import (
	"fmt"
	"log"
	"time"

	"github.com/boyism80/fm/game/client"
	"github.com/boyism80/fm/game/entity"
	"github.com/boyism80/fm/game/wz"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/types"
)

type MasterSkills struct {
	gs *GameServer
}

func (*MasterSkills) New(gs *GameServer) *MasterSkills {
	return &MasterSkills{
		gs: gs,
	}
}

func (h *MasterSkills) GetCommandName() string {
	return "스킬마스터"
}

func (h *MasterSkills) GetUsage() string {
	return "- 모든 배운 스킬의 레벨을 최대치로 올림"
}

func (h *MasterSkills) Handle(gameClient *client.GameClient, args ...string) error {
	character := gameClient.GetCharacter()
	if character == nil {
		return fmt.Errorf("character not found")
	}

	resources := h.gs.GetResources()
	if resources == nil {
		return fmt.Errorf("resources not available")
	}

	if character.Skills == nil {
		character.Skills = make(map[uint32]*entity.SkillEntry)
	}

	class := character.Class
	if class < 100 {
		return fmt.Errorf("character class %d is too low", class)
	}

	classIDs := getJobAdvancementClasses(class)
	addedCount := 0
	updatedCount := 0

	for _, classID := range classIDs {
		skillIDStart := uint32(classID) * 10000
		skillIDEnd := skillIDStart + 9999

		for skillID := skillIDStart; skillID <= skillIDEnd; skillID++ {
			wzSkill := resources.GetSkill(skillID)
			if wzSkill == nil {
				continue
			}

			masterLevel := 0
			if wzSkill.MasterLevel > 0 {
				masterLevel = wzSkill.MasterLevel
			} else if wzSkill.MaxLevel > 0 {
				masterLevel = wzSkill.MaxLevel
			} else {
				continue
			}

			skillEntry, exists := character.Skills[skillID]
			if !exists || skillEntry == nil {
				skillEntry = &entity.SkillEntry{
					Skill:       wzSkill,
					SkillLevel:  wzSkill.MaxLevel,
					MasterLevel: masterLevel,
					Expiration:  time.Time{},
					Owner:       character,
				}
				character.Skills[skillID] = skillEntry
				addedCount++
			} else {
				if skillEntry.Skill == nil {
					skillEntry.Skill = wzSkill
				}
				skillEntry.MasterLevel = masterLevel
				skillEntry.SkillLevel = wzSkill.MaxLevel
			}

			updatedCount++
			character.Send(&response.UpdateSkills{
				SkillID:     skillID,
				Level:       int32(skillEntry.SkillLevel),
				MasterLevel: int32(skillEntry.MasterLevel),
			}, types.SEND_POLICY_ENCRYPT)
		}
	}

	log.Printf("Command: Added %d new skills and mastered %d total skills for character %d", addedCount, updatedCount, character.GetID())
	return nil
}

func getJobAdvancementClasses(class uint16) []uint16 {
	if class < 100 {
		return []uint16{class}
	}

	classes := make([]uint16, 0, 4)

	baseClass := (class / 100) * 100
	secondClass := (class / 10) * 10
	thirdClass := secondClass + 1

	classes = append(classes, baseClass)

	if secondClass != baseClass {
		classes = append(classes, secondClass)
	}

	if thirdClass != secondClass && thirdClass <= class {
		classes = append(classes, thirdClass)
	}

	if class != thirdClass && class != secondClass && class != baseClass {
		classes = append(classes, class)
	}

	return classes
}

func isFourthClassSkill(skillID uint32, wzSkill *wz.Skill) bool {
	classID := skillID / 10000

	if classID == 2312 {
		return true
	}

	if (wzSkill.MaxLevel <= 15 && !wzSkill.Invisible && wzSkill.MasterLevel <= 0) ||
		skillID == 3220010 || skillID == 3120011 || skillID == 33120010 || skillID == 32120009 ||
		skillID == 5321006 || skillID == 21120011 || skillID == 22181004 || skillID == 4340010 {
		return false
	}

	if classID >= 2212 && classID < 3000 {
		return (classID % 10) >= 7
	}

	if classID >= 430 && classID <= 434 {
		return (classID%10) == 4 || wzSkill.MasterLevel > 0
	}

	return (classID%10) == 2 && skillID < 90000000
}
