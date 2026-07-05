package entity

import (
	"strconv"
	"time"

	"github.com/boyism80/fm/core/clock"
)

func parseQuestEventTime(raw string) (time.Time, bool) {
	if len(raw) != 10 {
		return time.Time{}, false
	}
	year, err := strconv.Atoi(raw[0:4])
	if err != nil {
		return time.Time{}, false
	}
	month, err := strconv.Atoi(raw[4:6])
	if err != nil {
		return time.Time{}, false
	}
	day, err := strconv.Atoi(raw[6:8])
	if err != nil {
		return time.Time{}, false
	}
	hour, err := strconv.Atoi(raw[8:10])
	if err != nil {
		return time.Time{}, false
	}
	if month < 1 || month > 12 || day < 1 || day > 31 || hour < 0 || hour > 23 {
		return time.Time{}, false
	}
	return time.Date(year, time.Month(month), day, hour, 0, 0, 0, time.Local), true
}

func questSameCalendarDay(a, b time.Time) bool {
	ay, am, ad := a.In(time.Local).Date()
	by, bm, bd := b.In(time.Local).Date()
	return ay == by && am == bm && ad == bd
}

func questRequirementNow() time.Time {
	return clock.Now()
}
