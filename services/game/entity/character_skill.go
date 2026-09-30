package entity

import (
	"log"
)

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
	entry := ch.Skills.Get(skillID)
	if entry == nil {
		return 0
	}
	return entry.Level()
}

func (ch *Character) UseAttackSkill(skillID uint32) bool {
	if ch.GameWorld == nil {
		return false
	}
	resources := ch.GameWorld.GetResources()
	if resources == nil {
		return false
	}
	wzSkill := resources.GetSkill(skillID)
	if wzSkill == nil {
		log.Printf("Skill not found: %d", skillID)
		return false
	}
	skillLevel := ch.GetTotalSkillLevel(skillID)
	if skillLevel <= 0 {
		log.Printf("Character does not have skill %d or skill level is 0", skillID)
		return false
	}
	levelData := wzSkill.GetLevelData(skillLevel)
	if levelData == nil {
		return false
	}
	if levelData.Cooldown > 0 {
		skillEntry := ch.Skills.Get(skillID)
		if skillEntry == nil || skillEntry.IsCooling() {
			return false
		}
		skillEntry.StartCooldown(levelData.Cooldown)
	}
	return true
}
