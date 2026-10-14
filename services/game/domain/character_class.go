package domain

import "github.com/boyism80/fm/services/game/constant"

func (ch *Character) ChangeClass(newClass constant.ClassType) {
	oldClass := ch.Class
	ch.Class = newClass

	ch.grantClassChangeSP(newClass)

	ch.Listener.OnClassChange(ch, oldClass, newClass)
}

func (ch *Character) grantClassChangeSP(newClass constant.ClassType) {
	if ch.Class == constant.ClassBeginner {
		return
	}

	ch.Points.SP++

	if newClass%10 >= 2 {
		ch.Points.SP += 2
	}

	if newClass%100 == 0 {
		minLevel := uint8(10)
		if newClass == constant.ClassMagician {
			minLevel = 8
		}

		if ch.level > minLevel {
			spToGrant := uint16(3 * (int(ch.level) - int(minLevel)))
			ch.Points.SP += spToGrant
		}
	}
}
