package entity

func (ch *Character) IsAdventurer() bool {
	return ch.Class < 1000
}

func (ch *Character) IsBeginner() bool {
	return ch.Class == 0
}

func (ch *Character) getClassAdvancementLevel() int {
	class := ch.Class

	if class == 0 {
		return 0
	}

	if class >= 100 {
		secondDigit := (class / 10) % 10
		thirdDigit := class % 10

		if secondDigit == 0 {
			return 1
		} else if thirdDigit == 0 {
			return 2
		} else if thirdDigit == 1 {
			return 3
		} else {
			return 4
		}
	}

	return 0
}

func (ch *Character) IsCannon() bool {
	return ch.Class == 1 || ch.Class == 501 || (ch.Class >= 530 && ch.Class <= 532)
}

func (ch *Character) IsMagician() bool {
	return ch.Class >= 200 && ch.Class < 300
}

func (ch *Character) ChangeClass(newClass uint16) {
	oldClass := ch.Class
	ch.Class = newClass

	ch.grantClassChangeSP(newClass)
	ch.initializeBaseSkills(newClass)

	if ch.Listener != nil {
		ch.Listener.OnClassChange(oldClass, newClass)
	}
}

func (ch *Character) grantClassChangeSP(newClass uint16) {
	if ch.IsBeginner() {
		return
	}

	ch.SkillPoint++

	if newClass >= 100 {
		thirdDigit := newClass % 10
		if thirdDigit >= 2 {
			ch.SkillPoint += 2
		}
	}

	if newClass%100 == 0 {
		minLevel := uint8(10)
		if newClass == 200 {
			minLevel = 8
		}

		if ch.level > minLevel {
			spToGrant := uint16(3 * (int(ch.level) - int(minLevel)))
			ch.SkillPoint += spToGrant
		}
	}
}
