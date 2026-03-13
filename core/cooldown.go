package core

import "time"

// AllowCooldown returns true if at least interval has passed since *last (or *last is zero).
// When true, *last is set to now. Used for rate-limiting heal-over-time and similar.
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
