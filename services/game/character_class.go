package game

import "github.com/boyism80/fm/maple/constant"

func (ch *Character) ChangeClass(newClass constant.ClassType) {
	oldClass := ch.class
	ch.class = newClass

	ch.grantClassChangeSP(newClass)

	ch.listener.OnClassChange(ch, oldClass, newClass, ch.points.sp)
	ch.NotifyLevelClassChanged()
}

func (ch *Character) grantClassChangeSP(newClass constant.ClassType) {
	if newClass == constant.ClassBeginner || newClass == constant.ClassNoblesse || newClass == constant.ClassLegend {
		return
	}

	ch.points.sp++

	if newClass%10 >= 2 {
		ch.points.sp += 2
	}

	if newClass%100 == 0 {
		minLevel := uint8(10)
		if newClass == constant.ClassMagician {
			minLevel = 8
		}

		if ch.level > minLevel {
			spToGrant := uint16(3 * (int(ch.level) - int(minLevel)))
			ch.points.sp += spToGrant
		}
	}
}
