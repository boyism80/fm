package entity

import (
	"time"

	"github.com/asynkron/protoactor-go/actor"
)

type ObjectTimer struct {
	Timer      *time.Timer
	Interval   time.Duration
	Repeat     bool
	Callback   func()
	NextFireAt time.Time
	Remaining  time.Duration
}

func (obj *ObjectCore) initTimers() {
	obj.timers = make(map[string]*ObjectTimer)
}

func (obj *ObjectCore) AddTimer(key string, interval time.Duration, repeat bool, callback func()) bool {
	if _, exists := obj.timers[key]; exists {
		obj.RemoveTimer(key)
	}
	if obj.GameWorld == nil || obj.self == nil {
		return false
	}
	mapInstance := obj.Map
	if mapInstance == nil {
		return false
	}
	pid := mapInstance.GetActorPID()
	if pid == nil {
		return false
	}
	entry := &ObjectTimer{
		Interval: interval,
		Repeat:   repeat,
		Callback: callback,
	}
	obj.timers[key] = entry
	obj.armTimer(key, entry, interval, pid)
	return true
}

func (obj *ObjectCore) armTimer(key string, entry *ObjectTimer, delay time.Duration, pid *actor.PID) {
	if entry == nil || obj.GameWorld == nil || obj.self == nil {
		return
	}
	k := key
	entry.NextFireAt = time.Now().Add(delay)
	entry.Timer = time.AfterFunc(delay, func() {
		if gw := obj.GameWorld; gw != nil {
			gw.GetSchedulerSystem().RunObjectTimer(pid, obj.self, k)
		}
	})
}

func (obj *ObjectCore) RescheduleTimer(key string) bool {
	entry := obj.timers[key]
	if entry == nil || !entry.Repeat {
		return false
	}
	mapInstance := obj.Map
	if mapInstance == nil {
		return false
	}
	pid := mapInstance.GetActorPID()
	if pid == nil {
		return false
	}
	if entry.Timer != nil {
		entry.Timer.Stop()
		entry.Timer = nil
	}
	obj.armTimer(key, entry, entry.Interval, pid)
	return true
}

func (obj *ObjectCore) RemoveTimer(key string) bool {
	entry := obj.timers[key]
	if entry == nil {
		return false
	}
	if entry.Timer != nil {
		entry.Timer.Stop()
	}
	delete(obj.timers, key)
	return true
}

func (obj *ObjectCore) GetTimerEntry(key string) *ObjectTimer {
	return obj.timers[key]
}

func (obj *ObjectCore) ClearTimers() {
	for key, entry := range obj.timers {
		if entry != nil && entry.Timer != nil {
			entry.Timer.Stop()
		}
		delete(obj.timers, key)
	}
}

func (obj *ObjectCore) SuspendTimers() {
	now := time.Now()
	for _, entry := range obj.timers {
		if entry == nil || entry.Timer == nil {
			continue
		}
		entry.Remaining = entry.NextFireAt.Sub(now)
		if entry.Remaining < 0 {
			entry.Remaining = 0
		}
		entry.Timer.Stop()
		entry.Timer = nil
	}
}

func (obj *ObjectCore) ResumeTimers(pid *actor.PID) {
	if pid == nil || obj.GameWorld == nil || obj.self == nil {
		return
	}
	for key, entry := range obj.timers {
		if entry == nil || entry.Timer != nil {
			continue
		}
		duration := entry.Remaining
		if duration <= 0 {
			duration = entry.Interval
		}
		entry.Remaining = 0
		obj.armTimer(key, entry, duration, pid)
	}
}
