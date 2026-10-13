package entity

import (
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/wz"
)

func (ch *Character) ChangeClass(newClass constant.ClassType) {
	oldClass := ch.Class.ID
	ch.Class = wz.NewClass(newClass)

	ch.grantClassChangeSP()

	ch.Listener.OnClassChange(ch, uint16(oldClass), uint16(newClass))
}

func (ch *Character) grantClassChangeSP() {
	if ch.Class.ID == constant.ClassBeginner {
		return
	}

	ch.Points.SP++

	if ch.Class.Advancement >= 4 {
		ch.Points.SP += 2
	}

	if ch.Class.Advancement <= 1 {
		minLevel := uint8(10)
		if ch.Class.ID == constant.ClassMagician {
			minLevel = 8
		}

		if ch.level > minLevel {
			spToGrant := uint16(3 * (int(ch.level) - int(minLevel)))
			ch.Points.SP += spToGrant
		}
	}
}
