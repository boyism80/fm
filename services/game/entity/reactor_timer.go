package entity

import "time"

const (
	reactorTimerStateRevertKey  = "reactor:stateRevert"
	reactorTimerItemActivateKey = "reactor:itemActivate"
	reactorTimerResetStateKey   = "reactor:resetState"
)

func (r *Reactor) ClearReactorTimers() {
	if r == nil {
		return
	}
	r.RemoveTimer(reactorTimerStateRevertKey)
	r.RemoveTimer(reactorTimerItemActivateKey)
	r.RemoveTimer(reactorTimerResetStateKey)
}

func (r *Reactor) ScheduleStateRevert(oldState byte, newState byte, delay time.Duration) bool {
	if r == nil || delay <= 0 {
		return false
	}
	return r.AddTimer(reactorTimerStateRevertKey, delay, false, func() {
		if r.State == oldState {
			r.ForceHitState(newState)
		}
	})
}

func (r *Reactor) ScheduleItemActivation(delay time.Duration, activate func()) bool {
	if r == nil || delay <= 0 || activate == nil || r.TimerActive {
		return false
	}
	r.TimerActive = true
	if !r.AddTimer(reactorTimerItemActivateKey, delay, false, func() {
		r.TimerActive = false
		activate()
	}) {
		r.TimerActive = false
		return false
	}
	return true
}

func (r *Reactor) ScheduleResetState(delay time.Duration) bool {
	if r == nil || delay <= 0 {
		return false
	}
	return r.AddTimer(reactorTimerResetStateKey, delay, false, func() {
		r.ForceHitState(0)
	})
}

func (r *Reactor) StateTimeOut(state byte) time.Duration {
	if r == nil || r.Wz == nil {
		return 0
	}
	event := r.Wz.States[state]
	if event == nil || event.TimeOut <= 0 {
		return 0
	}
	return time.Duration(event.TimeOut) * time.Millisecond
}
