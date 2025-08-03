package timer

import (
	"time"
)

// Timer represents a one-shot timer that executes logic on a specific logic thread
type Timer struct {
	ID         uint64
	Logic      func() error
	Callback   func(bool, error)
	CancelChan chan struct{}
	CreatedAt  time.Time
}

// Cancel cancels a timer if it hasn't been executed yet
func (t *Timer) Cancel() {
	select {
	case t.CancelChan <- struct{}{}:
		// Successfully cancelled
	default:
		// Timer already executed or cancelled
	}
}
