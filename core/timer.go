package core

import (
	"fmt"
	"sync"
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

// RepeatingTimer represents a repeating timer that executes logic periodically
type RepeatingTimer struct {
	ID         uint64
	Interval   time.Duration
	Logic      func() error
	Callback   func(bool, error)
	CancelChan chan struct{}
	CreatedAt  time.Time
	IsRunning  bool
}

// Cancel cancels a repeating timer
func (rt *RepeatingTimer) Cancel() {
	select {
	case rt.CancelChan <- struct{}{}:
		// Successfully cancelled
	default:
		// Timer already cancelled
	}
}

// TimerManager manages timers for a logic thread
type TimerManager struct {
	timers          map[uint64]*Timer
	timerMu         sync.RWMutex
	nextTimerID     uint64
	repeatingTimers map[uint64]*RepeatingTimer
	repeatingMu     sync.RWMutex
	nextRepeatingID uint64
}

// NewTimerManager creates a new timer manager
func NewTimerManager() *TimerManager {
	return &TimerManager{
		timers:          make(map[uint64]*Timer),
		repeatingTimers: make(map[uint64]*RepeatingTimer),
	}
}

// Schedule schedules a one-shot timer
func (tm *TimerManager) Schedule(duration time.Duration, logic func() error, callback func(bool, error)) *Timer {
	tm.timerMu.Lock()
	defer tm.timerMu.Unlock()

	timer := &Timer{
		ID:         tm.nextTimerID,
		Logic:      logic,
		Callback:   callback,
		CancelChan: make(chan struct{}),
		CreatedAt:  time.Now(),
	}
	tm.nextTimerID++

	tm.timers[timer.ID] = timer

	return timer
}

// SetRepeatingTimer sets a repeating timer
func (tm *TimerManager) SetRepeatingTimer(interval time.Duration, logic func() error, callback func(bool, error)) *RepeatingTimer {
	tm.repeatingMu.Lock()
	defer tm.repeatingMu.Unlock()

	repeatingTimer := &RepeatingTimer{
		ID:         tm.nextRepeatingID,
		Interval:   interval,
		Logic:      logic,
		Callback:   callback,
		CancelChan: make(chan struct{}),
		CreatedAt:  time.Now(),
		IsRunning:  true,
	}
	tm.nextRepeatingID++

	tm.repeatingTimers[repeatingTimer.ID] = repeatingTimer

	return repeatingTimer
}

// GetTimerCount returns the number of active timers
func (tm *TimerManager) GetTimerCount() int {
	tm.timerMu.RLock()
	defer tm.timerMu.RUnlock()
	return len(tm.timers)
}

// GetRepeatingTimerCount returns the number of active repeating timers
func (tm *TimerManager) GetRepeatingTimerCount() int {
	tm.repeatingMu.RLock()
	defer tm.repeatingMu.RUnlock()
	return len(tm.repeatingTimers)
}

// CancelRepeatingTimer cancels a repeating timer by ID
func (tm *TimerManager) CancelRepeatingTimer(timerID uint64) error {
	tm.repeatingMu.Lock()
	defer tm.repeatingMu.Unlock()

	if timer, exists := tm.repeatingTimers[timerID]; exists {
		timer.Cancel()
		return nil
	}
	return fmt.Errorf("repeating timer %d not found", timerID)
}

// RemoveTimer removes a timer from the manager
func (tm *TimerManager) RemoveTimer(timerID uint64) {
	tm.timerMu.Lock()
	delete(tm.timers, timerID)
	tm.timerMu.Unlock()
}

// RemoveRepeatingTimer removes a repeating timer from the manager
func (tm *TimerManager) RemoveRepeatingTimer(timerID uint64) {
	tm.repeatingMu.Lock()
	delete(tm.repeatingTimers, timerID)
	tm.repeatingMu.Unlock()
}

// CancelAllTimers cancels all active timers
func (tm *TimerManager) CancelAllTimers() {
	tm.timerMu.Lock()
	for _, timer := range tm.timers {
		timer.Cancel()
	}
	tm.timerMu.Unlock()

	tm.repeatingMu.Lock()
	for _, timer := range tm.repeatingTimers {
		timer.Cancel()
	}
	tm.repeatingMu.Unlock()
}
