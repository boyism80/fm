package entity

import (
	"time"

	"github.com/boyism80/fm/game/wz"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/types"
)

func (ch *Character) getSkillBookIndexByLevel(level uint8) int {
	if ch.IsBeginner() {
		return -1
	}

	isMagician := ch.IsMagician()
	minLevel := uint8(10)
	if isMagician {
		minLevel = 8
	}

	if level < minLevel {
		return -1
	}

	if level <= 30 {
		return 0
	} else if level <= 70 {
		return 1
	} else if level <= 120 {
		return 2
	} else {
		return 3
	}
}

func (ch *Character) GetSkillBookIndex() int {
	class := ch.Class

	if class >= 100 {
		secondDigit := (class / 10) % 10
		thirdDigit := class % 10

		if secondDigit == 0 {
			return 0
		} else if thirdDigit == 0 {
			return 1
		} else if thirdDigit == 1 {
			return 2
		} else if thirdDigit == 2 {
			return 3
		}
	}

	return 0
}

func (ch *Character) GetSkillBookIndexForSkill(skillID uint32) int {
	classID := skillID / 10000

	if classID >= 100 {
		secondDigit := (classID / 10) % 10
		thirdDigit := classID % 10

		if secondDigit == 0 {
			return 0
		} else if thirdDigit == 0 {
			return 1
		} else if thirdDigit == 1 {
			return 2
		} else if thirdDigit == 2 {
			return 3
		}
	}

	return 0
}

func (ch *Character) RemainingSkillPoints() uint16 {
	return ch.SkillPoint
}

func (ch *Character) GetTotalSkillLevel(skillID uint32) int {
	if ch.Skills == nil {
		return 0
	}
	skillEntry, exists := ch.Skills[skillID]
	if !exists || skillEntry == nil {
		return 0
	}
	return skillEntry.SkillLevel
}

func (ch *Character) initializeBaseSkills(newClass uint16) {
	if ch.Context == nil {
		return
	}

	resources := ch.Context.GetResources()
	if resources == nil {
		return
	}

	advancementLevel := ch.getClassAdvancementLevel()
	if advancementLevel < 3 {
		return
	}

	classID := uint32(newClass)
	skillIDStart := classID * 10000
	skillIDEnd := skillIDStart + 9999

	if ch.Skills == nil {
		ch.Skills = make(map[uint32]*SkillEntry)
	}

	for skillID := skillIDStart; skillID <= skillIDEnd; skillID++ {
		wzSkill := resources.GetSkill(skillID)
		if wzSkill == nil {
			continue
		}

		if wzSkill.Invisible {
			continue
		}

		if !ch.isFourthClassSkill(skillID, wzSkill) {
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

		existingEntry, exists := ch.Skills[skillID]
		if exists && existingEntry != nil {
			if existingEntry.SkillLevel > 0 || existingEntry.MasterLevel > 0 {
				continue
			}
		}

		skillEntry := &SkillEntry{
			Wz:          wzSkill,
			SkillLevel:  0,
			MasterLevel: masterLevel,
			Expiration:  time.Time{},
			Owner:       ch,
		}
		ch.Skills[skillID] = skillEntry

		if ch.Listener != nil {
			ch.Send(&response.UpdateSkills{
				SkillID:     skillID,
				Level:       0,
				MasterLevel: int32(skillEntry.MasterLevel),
			}, types.SEND_POLICY_ENCRYPT)
		}
	}
}

func (ch *Character) isFourthClassSkill(skillID uint32, wzSkill *wz.Skill) bool {
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
