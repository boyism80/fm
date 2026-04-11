package entity

// getClassParent returns the parent class code in the job tree. ok is false at roots.
func getClassParent(class uint16) (parent uint16, ok bool) {
	if class == 0 {
		return 0, false
	}
	if class%10 != 0 {
		return class - 1, true
	}
	switch {
	case class == 100 || class == 200 || class == 300 || class == 400 || class == 500 || class == 800 || class == 900:
		return 0, false
	case class == 1000 || class == 2000:
		return 0, false
	case class >= 1100 && class <= 1500 && class%100 == 0:
		return 1000, true
	case class >= 2100 && class%100 == 0:
		return 2000, true
	case class < 1000:
		return (class / 100) * 100, true
	default:
		return (class/10 - 1) * 10, true
	}
}

// ClassOf returns true if the character's class is the given class or any advancement of it in the job tree.
func (ch *Character) ClassOf(classCode uint16) bool {
	c := ch.Class
	for c != 0 {
		if c == classCode {
			return true
		}
		var ok bool
		c, ok = getClassParent(c)
		if !ok {
			break
		}
	}
	return c == classCode
}

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

	ch.Listener.OnClassChange(ch, oldClass, newClass)
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
