package game

import (
	"time"

	"github.com/boyism80/fm/core/clock"

	"github.com/asynkron/protoactor-go/actor"
)

type ObjectTimer struct {
	timer      *time.Timer
	interval   time.Duration
	Repeat     bool
	Callback   func()
	nextFireAt time.Time
	remaining  time.Duration
}

func (obj *ObjectCore) initTimers() {
	obj.timers = make(map[string]*ObjectTimer)
}

func (obj *ObjectCore) AddTimer(key string, interval time.Duration, repeat bool, callback func()) bool {
	if _, exists := obj.timers[key]; exists {
		obj.RemoveTimer(key)
	}
	if obj.gameWorld == nil || obj.self == nil {
		return false
	}
	mapInstance := obj.onMap
	if mapInstance == nil {
		return false
	}
	pid := mapInstance.LogicActorPID()
	if pid == nil {
		return false
	}
	entry := &ObjectTimer{
		interval: interval,
		Repeat:   repeat,
		Callback: callback,
	}
	obj.timers[key] = entry
	obj.scheduleTimer(key, entry, interval, pid)
	return true
}

func (obj *ObjectCore) scheduleTimer(key string, entry *ObjectTimer, delay time.Duration, pid *actor.PID) {
	if entry == nil || obj.gameWorld == nil || obj.self == nil {
		return
	}
	k := key
	gen := obj.timerGen.Load()
	entry.nextFireAt = clock.Now().Add(delay)
	entry.timer = time.AfterFunc(delay, func() {
		if obj.timerGen.Load() != gen {
			return
		}
		if gw := obj.gameWorld; gw != nil {
			gw.Systems().Get[*SchedulerSystem]().RunObjectTimer(pid, obj.self, k)
		}
	})
}

func (obj *ObjectCore) RescheduleTimer(key string) bool {
	entry := obj.timers[key]
	if entry == nil || !entry.Repeat {
		return false
	}
	mapInstance := obj.onMap
	if mapInstance == nil {
		return false
	}
	pid := mapInstance.LogicActorPID()
	if pid == nil {
		return false
	}
	if entry.timer != nil {
		entry.timer.Stop()
		entry.timer = nil
	}
	obj.scheduleTimer(key, entry, entry.interval, pid)
	return true
}

func (obj *ObjectCore) RemoveTimer(key string) bool {
	entry := obj.timers[key]
	if entry == nil {
		return false
	}
	if entry.timer != nil {
		entry.timer.Stop()
	}
	delete(obj.timers, key)
	return true
}

func (obj *ObjectCore) GetTimerEntry(key string) *ObjectTimer {
	return obj.timers[key]
}

func (obj *ObjectCore) ClearTimers() {
	obj.timerGen.Add(1)
	for key, entry := range obj.timers {
		if entry != nil && entry.timer != nil {
			entry.timer.Stop()
		}
		delete(obj.timers, key)
	}
}

func (obj *ObjectCore) SuspendTimers() {
	now := clock.Now()
	for _, entry := range obj.timers {
		if entry == nil || entry.timer == nil {
			continue
		}
		entry.remaining = entry.nextFireAt.Sub(now)
		if entry.remaining < 0 {
			entry.remaining = 0
		}
		entry.timer.Stop()
		entry.timer = nil
	}
}

func (obj *ObjectCore) ResumeTimers(pid *actor.PID) {
	if pid == nil || obj.gameWorld == nil || obj.self == nil {
		return
	}
	for key, entry := range obj.timers {
		if entry == nil || entry.timer != nil {
			continue
		}
		duration := entry.remaining
		if duration <= 0 && entry.Repeat {
			duration = entry.interval
		}
		entry.remaining = 0
		obj.scheduleTimer(key, entry, duration, pid)
	}
}
