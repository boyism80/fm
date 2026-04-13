package entity

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
	entry := ch.Skills.Get(skillID)
	if entry == nil {
		return 0
	}
	return entry.Level()
}
