package core

import (
	"fmt"
	"github.com/boyism80/fm/core/clock"
	"sync"
	"time"
)

type Timer struct {
	ID         uint64
	Logic      func() error
	Callback   func(bool, error)
	CancelChan chan struct{}
	CreatedAt  time.Time
}

func (t *Timer) Cancel() {
	select {
	case t.CancelChan <- struct{}{}:

	default:

	}
}

type RepeatingTimer struct {
	ID         uint64
	Interval   time.Duration
	Logic      func() error
	Callback   func(bool, error)
	CancelChan chan struct{}
	CreatedAt  time.Time
	IsRunning  bool
}

func (rt *RepeatingTimer) Cancel() {
	select {
	case rt.CancelChan <- struct{}{}:

	default:

	}
}

type TimerManager struct {
	timers          map[uint64]*Timer
	timerMu         sync.RWMutex
	nextTimerID     uint64
	repeatingTimers map[uint64]*RepeatingTimer
	repeatingMu     sync.RWMutex
	nextRepeatingID uint64
}

func NewTimerManager() *TimerManager {
	return &TimerManager{
		timers:          make(map[uint64]*Timer),
		repeatingTimers: make(map[uint64]*RepeatingTimer),
	}
}

func (tm *TimerManager) Schedule(duration time.Duration, logic func() error, callback func(bool, error)) *Timer {
	tm.timerMu.Lock()
	defer tm.timerMu.Unlock()

	timer := &Timer{
		ID:         tm.nextTimerID,
		Logic:      logic,
		Callback:   callback,
		CancelChan: make(chan struct{}),
		CreatedAt:  clock.Now(),
	}
	tm.nextTimerID++

	tm.timers[timer.ID] = timer

	return timer
}

func (tm *TimerManager) SetRepeatingTimer(interval time.Duration, logic func() error, callback func(bool, error)) *RepeatingTimer {
	tm.repeatingMu.Lock()
	defer tm.repeatingMu.Unlock()

	repeatingTimer := &RepeatingTimer{
		ID:         tm.nextRepeatingID,
		Interval:   interval,
		Logic:      logic,
		Callback:   callback,
		CancelChan: make(chan struct{}),
		CreatedAt:  clock.Now(),
		IsRunning:  true,
	}
	tm.nextRepeatingID++

	tm.repeatingTimers[repeatingTimer.ID] = repeatingTimer

	return repeatingTimer
}

func (tm *TimerManager) GetTimerCount() int {
	tm.timerMu.RLock()
	defer tm.timerMu.RUnlock()
	return len(tm.timers)
}

func (tm *TimerManager) GetRepeatingTimerCount() int {
	tm.repeatingMu.RLock()
	defer tm.repeatingMu.RUnlock()
	return len(tm.repeatingTimers)
}

func (tm *TimerManager) CancelRepeatingTimer(timerID uint64) error {
	tm.repeatingMu.Lock()
	defer tm.repeatingMu.Unlock()

	if timer, exists := tm.repeatingTimers[timerID]; exists {
		timer.Cancel()
		return nil
	}
	return fmt.Errorf("repeating timer %d not found", timerID)
}

func (tm *TimerManager) RemoveTimer(timerID uint64) {
	tm.timerMu.Lock()
	delete(tm.timers, timerID)
	tm.timerMu.Unlock()
}

func (tm *TimerManager) RemoveRepeatingTimer(timerID uint64) {
	tm.repeatingMu.Lock()
	delete(tm.repeatingTimers, timerID)
	tm.repeatingMu.Unlock()
}

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
