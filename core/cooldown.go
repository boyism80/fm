package core

import "time"

func AllowCooldown(last *time.Time, interval time.Duration, now time.Time) bool {
	if last == nil {
		return false
	}
	if last.IsZero() || now.Sub(*last) >= interval {
		*last = now
		return true
	}
	return false
}
