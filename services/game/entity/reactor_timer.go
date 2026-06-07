package entity

import "time"

const (
	reactorTimerStateRevertKey  = "reactor:stateRevert"
	reactorTimerDelayedHitKey   = "reactor:delayedHit"
	reactorTimerItemActivateKey = "reactor:itemActivate"
	reactorTimerResetStateKey   = "reactor:resetState"
)

func (r *Reactor) ClearReactorTimers() {
	if r == nil {
		return
	}
	r.RemoveTimer(reactorTimerStateRevertKey)
	r.RemoveTimer(reactorTimerDelayedHitKey)
	r.RemoveTimer(reactorTimerItemActivateKey)
	r.RemoveTimer(reactorTimerResetStateKey)
}

func (r *Reactor) CancelStateRevert() {
	if r == nil {
		return
	}
	r.RemoveTimer(reactorTimerStateRevertKey)
}

func (r *Reactor) ScheduleStateRevert(oldState byte, newState byte, delay time.Duration) bool {
	if r == nil || delay <= 0 {
		return false
	}
	r.CancelStateRevert()
	return r.AddTimer(reactorTimerStateRevertKey, delay, false, func() {
		if r.State == oldState {
			r.ForceHitState(newState)
		}
	})
}

func (r *Reactor) CancelDelayedHit() {
	if r == nil {
		return
	}
	r.RemoveTimer(reactorTimerDelayedHitKey)
}

func (r *Reactor) ScheduleDelayedHit(delay time.Duration, hit func()) bool {
	if r == nil || delay <= 0 || hit == nil {
		return false
	}
	r.CancelDelayedHit()
	return r.AddTimer(reactorTimerDelayedHitKey, delay, false, hit)
}

func (r *Reactor) CancelItemActivation() {
	if r == nil {
		return
	}
	r.TimerActive = false
	r.RemoveTimer(reactorTimerItemActivateKey)
}

func (r *Reactor) ScheduleItemActivation(delay time.Duration, activate func()) bool {
	if r == nil || delay <= 0 || activate == nil || r.TimerActive {
		return false
	}
	r.TimerActive = true
	r.RemoveTimer(reactorTimerItemActivateKey)
	if !r.AddTimer(reactorTimerItemActivateKey, delay, false, func() {
		r.TimerActive = false
		activate()
	}) {
		r.TimerActive = false
		return false
	}
	return true
}

func (r *Reactor) CancelResetState() {
	if r == nil {
		return
	}
	r.RemoveTimer(reactorTimerResetStateKey)
}

func (r *Reactor) ScheduleResetState(delay time.Duration) bool {
	if r == nil || delay <= 0 {
		return false
	}
	r.CancelResetState()
	return r.AddTimer(reactorTimerResetStateKey, delay, false, func() {
		r.ForceHitState(0)
	})
}

func (r *Reactor) ScheduleResetStateFromSpawn() bool {
	if r == nil || r.Spawn == nil {
		return false
	}
	return r.ScheduleResetState(r.Spawn.RespawnDelay())
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
