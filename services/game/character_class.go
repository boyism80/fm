package game

import "github.com/boyism80/fm/maple/constant"

func (ch *Character) ChangeClass(newClass constant.ClassType) {
	oldClass := ch.class
	ch.class = newClass

	ch.grantClassChangePoints(newClass)

	ch.listener.OnClassChange(ch, oldClass, newClass, ch.points.sp, ch.points.ap)
	ch.NotifyLevelClassChanged()
}

func (ch *Character) grantClassChangePoints(newClass constant.ClassType) {
	if newClass == constant.ClassBeginner || newClass == constant.ClassNoblesse || newClass == constant.ClassLegend {
		return
	}

	ch.points.sp++

	if newClass%10 >= 2 {
		ch.points.sp += 2
	}

	if newClass%10 >= 1 && ch.level >= 70 {
		ch.points.ap += 5
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
